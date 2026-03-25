package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"ai-go-service/internal/http/response"
	"ai-go-service/internal/service"
)

// DesktopMenuUseCase defines the menu use cases exposed by the desktop BFF.
type DesktopMenuUseCase interface {
	CurrentMenus(ctx context.Context, accessToken string) (service.DesktopMenusResponse, error)
	ListDesktopMenus(ctx context.Context, name string, menuType string, status string) (service.DesktopManagedMenusResponse, error)
	CreateDesktopMenu(ctx context.Context, input service.DesktopManagedMenu) (service.DesktopManagedMenu, error)
	UpdateDesktopMenu(ctx context.Context, input service.DesktopManagedMenu) (service.DesktopManagedMenu, error)
	DeleteDesktopMenu(ctx context.Context, id int64) error
}

type desktopManagedMenuRequest = service.DesktopManagedMenu
type desktopManagedMenuDeleteRequest struct { ID int64 `json:"id"` }

// DesktopCurrentMenusHandler returns an HTTP handler that resolves desktop menus via go-admin.
func DesktopCurrentMenusHandler(menuService DesktopMenuUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if menuService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop menu service is not ready")
			return
		}

		accessToken, err := parseDesktopBearerToken(r)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}

		result, err := menuService.CurrentMenus(r.Context(), accessToken)
		if err != nil {
			handleDesktopAuthError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopManagedMenusHandler returns menu-management rows.
func DesktopManagedMenusHandler(menuService DesktopMenuUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if menuService == nil { response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop menu service is not ready"); return }
		if _, err := parseDesktopBearerToken(r); err != nil { response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error()); return }
		result, err := menuService.ListDesktopMenus(r.Context(), r.URL.Query().Get("name"), r.URL.Query().Get("type"), r.URL.Query().Get("status"))
		if err != nil { response.Error(w, http.StatusInternalServerError, response.CodeInternal, "failed to load desktop menus"); return }
		response.JSON(w, http.StatusOK, result)
	}
}

func DesktopCreateManagedMenuHandler(menuService DesktopMenuUseCase) http.HandlerFunc { return desktopManagedMenuMutationHandler(menuService, true) }
func DesktopUpdateManagedMenuHandler(menuService DesktopMenuUseCase) http.HandlerFunc { return desktopManagedMenuMutationHandler(menuService, false) }
func desktopManagedMenuMutationHandler(menuService DesktopMenuUseCase, create bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if menuService == nil { response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop menu service is not ready"); return }
		if _, err := parseDesktopBearerToken(r); err != nil { response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error()); return }
		var request desktopManagedMenuRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil { response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON"); return }
		var result service.DesktopManagedMenu
		var err error
		if create { result, err = menuService.CreateDesktopMenu(r.Context(), request) } else { result, err = menuService.UpdateDesktopMenu(r.Context(), request) }
		if err != nil { handleDesktopDataError(w, err); return }
		response.JSON(w, http.StatusOK, result)
	}
}

func DesktopDeleteManagedMenuHandler(menuService DesktopMenuUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if menuService == nil { response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop menu service is not ready"); return }
		if _, err := parseDesktopBearerToken(r); err != nil { response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error()); return }
		var request desktopManagedMenuDeleteRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil { response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON"); return }
		if err := menuService.DeleteDesktopMenu(r.Context(), request.ID); err != nil { handleDesktopDataError(w, err); return }
		response.JSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}
