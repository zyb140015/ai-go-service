package goadmin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// GetCurrentMenus resolves the current user's menus through go-admin.
func (client *Client) GetCurrentMenus(ctx context.Context, accessToken string) ([]MenuResult, error) {
	req, err := client.newRequest(ctx, http.MethodGet, "/menu/current?platform=desktop", nil, accessToken)
	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call go-admin current menus: %w", errors.Join(ErrUpstreamFailure, err))
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, ErrUpstreamFailure
	}

	var payload Envelope[CurrentMenusResult]
	if err := decodeJSON(resp, &payload); err != nil {
		return nil, fmt.Errorf("decode current-menus response: %w", errors.Join(ErrInvalidResponse, err))
	}

	if payload.Code != http.StatusOK {
		message := businessMessage(payload.Msg, payload.Message)
		if strings.Contains(message, "登录信息无效") || strings.Contains(message, "未登录") {
			return nil, fmt.Errorf("%s: %w", message, ErrUnauthorized)
		}

		return nil, &BusinessError{Code: payload.Code, Message: message}
	}

	return payload.Result.List, nil
}
