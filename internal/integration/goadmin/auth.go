package goadmin

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type updatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type updateProfileRequest struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Sex      string `json:"sex"`
}

type avatarUploadResult struct {
	Avatar string `json:"avatar"`
}

// Login authenticates against go-admin using username and password.
func (client *Client) Login(ctx context.Context, username string, password string) (LoginEnvelope, error) {
	req, err := client.newRequest(ctx, http.MethodPost, "/login/password", loginRequest{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}, "")
	if err != nil {
		return LoginEnvelope{}, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return LoginEnvelope{}, fmt.Errorf("call go-admin login: %w", errors.Join(ErrUpstreamFailure, err))
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return LoginEnvelope{}, ErrUpstreamFailure
	}

	var payload LoginEnvelope
	if err := decodeJSON(resp, &payload); err != nil {
		return LoginEnvelope{}, fmt.Errorf("decode login response: %w", errors.Join(ErrInvalidResponse, err))
	}

	if payload.Code != http.StatusOK {
		message := businessMessage(payload.Msg, payload.Message)
		if strings.Contains(message, "用户名或密码错误") {
			return LoginEnvelope{}, fmt.Errorf("%s: %w", message, ErrUnauthorized)
		}

		return LoginEnvelope{}, &BusinessError{Code: payload.Code, Message: message}
	}

	if strings.TrimSpace(payload.Result.Token) == "" || strings.TrimSpace(payload.Result.RefreshToken) == "" {
		return LoginEnvelope{}, ErrInvalidResponse
	}

	return payload, nil
}

// Refresh exchanges a refresh token for a new token pair.
func (client *Client) Refresh(ctx context.Context, refreshToken string) (LoginResult, error) {
	req, err := client.newRequest(ctx, http.MethodPost, "/login/refresh", refreshRequest{
		RefreshToken: strings.TrimSpace(refreshToken),
	}, "")
	if err != nil {
		return LoginResult{}, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return LoginResult{}, fmt.Errorf("call go-admin refresh: %w", errors.Join(ErrUpstreamFailure, err))
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return LoginResult{}, ErrUpstreamFailure
	}

	var payload Envelope[LoginResult]
	if err := decodeJSON(resp, &payload); err != nil {
		return LoginResult{}, fmt.Errorf("decode refresh response: %w", errors.Join(ErrInvalidResponse, err))
	}

	if payload.Code != http.StatusOK {
		message := businessMessage(payload.Msg, payload.Message)
		if strings.Contains(message, "刷新令牌无效") || strings.Contains(message, "未登录") || strings.Contains(message, "无效") {
			return LoginResult{}, fmt.Errorf("%s: %w", message, ErrUnauthorized)
		}

		return LoginResult{}, &BusinessError{Code: payload.Code, Message: message}
	}

	if strings.TrimSpace(payload.Result.Token) == "" || strings.TrimSpace(payload.Result.RefreshToken) == "" {
		return LoginResult{}, ErrInvalidResponse
	}

	return payload.Result, nil
}

// GetCurrentUser resolves the current user through go-admin.
func (client *Client) GetCurrentUser(ctx context.Context, accessToken string) (ProfileResult, error) {
	req, err := client.newRequest(ctx, http.MethodGet, "/user/profile", nil, accessToken)
	if err != nil {
		return ProfileResult{}, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return ProfileResult{}, fmt.Errorf("call go-admin current user: %w", errors.Join(ErrUpstreamFailure, err))
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return ProfileResult{}, ErrUpstreamFailure
	}

	var payload Envelope[ProfileResult]
	if err := decodeJSON(resp, &payload); err != nil {
		return ProfileResult{}, fmt.Errorf("decode current-user response: %w", errors.Join(ErrInvalidResponse, err))
	}

	if payload.Code != http.StatusOK {
		message := businessMessage(payload.Msg, payload.Message)
		if strings.Contains(message, "登录信息无效") || strings.Contains(message, "未登录") {
			return ProfileResult{}, fmt.Errorf("%s: %w", message, ErrUnauthorized)
		}

		return ProfileResult{}, &BusinessError{Code: payload.Code, Message: message}
	}

	if payload.Result.ID <= 0 || strings.TrimSpace(payload.Result.Username) == "" {
		return ProfileResult{}, ErrInvalidResponse
	}

	return payload.Result, nil
}

// UpdateCurrentPassword updates the current user's password through go-admin.
func (client *Client) UpdateCurrentPassword(ctx context.Context, accessToken string, oldPassword string, newPassword string) error {
	req, err := client.newRequest(ctx, http.MethodPut, "/user/profile/password", updatePasswordRequest{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}, accessToken)
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call go-admin update password: %w", errors.Join(ErrUpstreamFailure, err))
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return ErrUpstreamFailure
	}

	var payload Envelope[map[string]any]
	if err := decodeJSON(resp, &payload); err != nil {
		return fmt.Errorf("decode update-password response: %w", errors.Join(ErrInvalidResponse, err))
	}

	if payload.Code != http.StatusOK {
		message := businessMessage(payload.Msg, payload.Message)
		if strings.Contains(message, "登录信息无效") || strings.Contains(message, "未登录") {
			return fmt.Errorf("%s: %w", message, ErrUnauthorized)
		}

		return &BusinessError{Code: payload.Code, Message: message}
	}

	return nil
}

// UpdateCurrentProfile updates the current user's profile through go-admin.
func (client *Client) UpdateCurrentProfile(ctx context.Context, accessToken string, username string, avatar string, sex string) (ProfileResult, error) {
	req, err := client.newRequest(ctx, http.MethodPut, "/user/profile", updateProfileRequest{
		Username: strings.TrimSpace(username),
		Avatar:   strings.TrimSpace(avatar),
		Sex:      strings.TrimSpace(sex),
	}, accessToken)
	if err != nil {
		return ProfileResult{}, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return ProfileResult{}, fmt.Errorf("call go-admin update profile: %w", errors.Join(ErrUpstreamFailure, err))
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return ProfileResult{}, ErrUpstreamFailure
	}

	var payload Envelope[ProfileResult]
	if err := decodeJSON(resp, &payload); err != nil {
		return ProfileResult{}, fmt.Errorf("decode update-profile response: %w", errors.Join(ErrInvalidResponse, err))
	}

	if payload.Code != http.StatusOK {
		message := businessMessage(payload.Msg, payload.Message)
		if strings.Contains(message, "登录信息无效") || strings.Contains(message, "未登录") {
			return ProfileResult{}, fmt.Errorf("%s: %w", message, ErrUnauthorized)
		}
		return ProfileResult{}, &BusinessError{Code: payload.Code, Message: message}
	}

	return payload.Result, nil
}

// UploadCurrentAvatar uploads the current user's avatar through go-admin.
func (client *Client) UploadCurrentAvatar(ctx context.Context, accessToken string, fileName string, fileContent []byte) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", fmt.Errorf("create multipart file: %w", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(fileContent)); err != nil {
		return "", fmt.Errorf("write multipart file: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, client.baseURL+"/user/profile/avatar", &body)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", bearerPrefix+strings.TrimSpace(accessToken))

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call go-admin upload avatar: %w", errors.Join(ErrUpstreamFailure, err))
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return "", ErrUpstreamFailure
	}

	var payload Envelope[avatarUploadResult]
	if err := decodeJSON(resp, &payload); err != nil {
		return "", fmt.Errorf("decode upload-avatar response: %w", errors.Join(ErrInvalidResponse, err))
	}
	if payload.Code != http.StatusOK {
		message := businessMessage(payload.Msg, payload.Message)
		if strings.Contains(message, "登录信息无效") || strings.Contains(message, "未登录") {
			return "", fmt.Errorf("%s: %w", message, ErrUnauthorized)
		}
		return "", &BusinessError{Code: payload.Code, Message: message}
	}
	return payload.Result.Avatar, nil
}
