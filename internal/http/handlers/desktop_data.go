package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"ai-go-service/internal/http/response"
	"ai-go-service/internal/service"
)

// DesktopDataUseCase defines the desktop data use cases exposed by the desktop BFF.
type DesktopDataUseCase interface {
	ListMessages(ctx context.Context, messageType string, keyword string, page int, pageSize int) (service.DesktopMessagesResponse, error)
	ListMonitor(ctx context.Context, metric string, level string, page int, pageSize int) (service.DesktopMonitorResponse, error)
	ListAnnouncements(ctx context.Context, title string, announcementType string, status string, page int, pageSize int) (service.DesktopAnnouncementsResponse, error)
	ListLoginLogs(ctx context.Context, name string, status string, startDate string, endDate string, page int, pageSize int) (service.DesktopLoginLogsResponse, error)
	ListOperationLogs(ctx context.Context, name string, status string, startDate string, endDate string, page int, pageSize int) (service.DesktopOperationLogsResponse, error)
	ListDicts(ctx context.Context, name string, dictID string) (service.DesktopDictsResponse, error)
	CreateDict(ctx context.Context, input service.DesktopDict) (service.DesktopDict, error)
	UpdateDict(ctx context.Context, input service.DesktopDict) (service.DesktopDict, error)
	DeleteDict(ctx context.Context, id int64) error
	CreateDictItem(ctx context.Context, dictID int64, input service.DesktopDictItem) (service.DesktopDictItem, error)
	UpdateDictItem(ctx context.Context, input service.DesktopDictItem) (service.DesktopDictItem, error)
	DeleteDictItem(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, department string, uid string, name string, role string, page int, pageSize int) (service.DesktopUsersResponse, error)
	CreateUser(ctx context.Context, input service.DesktopUserListItem) (service.DesktopUserListItem, error)
	UpdateUser(ctx context.Context, input service.DesktopUserListItem) (service.DesktopUserListItem, error)
	DeleteUser(ctx context.Context, id int64) error
	ResetUserPassword(ctx context.Context, id int64) error
	ListRoles(ctx context.Context, name string, key string, status string, page int, pageSize int) (service.DesktopRolesResponse, error)
	CreateRole(ctx context.Context, input service.DesktopRole) (service.DesktopRole, error)
	UpdateRole(ctx context.Context, input service.DesktopRole) (service.DesktopRole, error)
	DeleteRole(ctx context.Context, id int64) error
	GetRoleDataPermission(ctx context.Context, roleID int64) (service.DesktopRoleDataPermission, error)
	SaveRoleDataPermission(ctx context.Context, input service.DesktopRoleDataPermission) (service.DesktopRoleDataPermission, error)
	ListOrgs(ctx context.Context, department string, name string, code string, page int, pageSize int) (service.DesktopOrgsResponse, error)
	CreateOrg(ctx context.Context, input service.DesktopOrg) (service.DesktopOrg, error)
	UpdateOrg(ctx context.Context, input service.DesktopOrg) (service.DesktopOrg, error)
	DeleteOrg(ctx context.Context, id int64) error
	ListTenants(ctx context.Context, name string, code string, status string, page int, pageSize int) (service.DesktopTenantsResponse, error)
	CreateTenant(ctx context.Context, input service.DesktopTenant) (service.DesktopTenant, error)
	UpdateTenant(ctx context.Context, input service.DesktopTenant) (service.DesktopTenant, error)
	DeleteTenant(ctx context.Context, id int64) error
	ListMenuTemplates(ctx context.Context, name string, description string, page int, pageSize int) (service.DesktopMenuTemplatesResponse, error)
	CreateMenuTemplate(ctx context.Context, input service.DesktopMenuTemplate) (service.DesktopMenuTemplate, error)
	UpdateMenuTemplate(ctx context.Context, input service.DesktopMenuTemplate) (service.DesktopMenuTemplate, error)
	DeleteMenuTemplate(ctx context.Context, id int64) error
	GetSystemStats(ctx context.Context, rangeKey string) (service.DesktopSystemStatsResponse, error)
	GetUsageStats(ctx context.Context, rangeKey string, tenantCode string) (service.DesktopUsageStatsResponse, error)
	GetSystemSettings(ctx context.Context) (service.DesktopSystemSettings, error)
	SaveSystemSettings(ctx context.Context, input service.DesktopSystemSettings) (service.DesktopSystemSettings, error)
	CreateAnnouncement(ctx context.Context, input service.DesktopAnnouncement) (service.DesktopAnnouncement, error)
	UpdateAnnouncement(ctx context.Context, input service.DesktopAnnouncement) (service.DesktopAnnouncement, error)
	DeleteAnnouncement(ctx context.Context, id int64) error
	PublishAnnouncement(ctx context.Context, id int64) error
	MarkMessageRead(ctx context.Context, id int64) error
	MarkMessagesRead(ctx context.Context, ids []int64) error
	UpdateMonitorStatus(ctx context.Context, id int64, status string) error
}

// DesktopSystemStatsHandler returns aggregated system statistics.
func DesktopSystemStatsHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		result, err := dataService.GetSystemStats(r.Context(), r.URL.Query().Get("range"))
		if err != nil {
			response.Error(w, http.StatusInternalServerError, response.CodeInternal, "failed to load system statistics")
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopUsageStatsHandler returns aggregated usage statistics.
func DesktopUsageStatsHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		result, err := dataService.GetUsageStats(r.Context(), r.URL.Query().Get("range"), r.URL.Query().Get("tenant"))
		if err != nil {
			response.Error(w, http.StatusInternalServerError, response.CodeInternal, "failed to load usage statistics")
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

type desktopOrgRequest = service.DesktopOrg
type desktopTenantRequest = service.DesktopTenant
type desktopUserRequest = service.DesktopUserListItem
type desktopRoleRequest = service.DesktopRole
type desktopRoleDataPermissionRequest = service.DesktopRoleDataPermission
type desktopDictRequest = service.DesktopDict
type desktopDictItemRequest struct {
	DictID int64 `json:"dictId"`
	service.DesktopDictItem
}

// DesktopMenuTemplatesHandler returns an HTTP handler that lists menu templates.
func DesktopMenuTemplatesHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		page, pageSize := parsePagination(r)
		result, err := dataService.ListMenuTemplates(r.Context(), r.URL.Query().Get("name"), r.URL.Query().Get("description"), page, pageSize)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, response.CodeInternal, "failed to load desktop menu templates")
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopOrgsHandler returns an HTTP handler that lists organizations.
func DesktopOrgsHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		page, pageSize := parsePagination(r)
		result, err := dataService.ListOrgs(r.Context(), r.URL.Query().Get("department"), r.URL.Query().Get("name"), r.URL.Query().Get("code"), page, pageSize)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, response.CodeInternal, "failed to load desktop orgs")
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopTenantsHandler returns an HTTP handler that lists tenants.
func DesktopTenantsHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		page, pageSize := parsePagination(r)
		result, err := dataService.ListTenants(r.Context(), r.URL.Query().Get("name"), r.URL.Query().Get("code"), r.URL.Query().Get("status"), page, pageSize)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, response.CodeInternal, "failed to load desktop tenants")
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopCreateOrgHandler creates an organization.
func DesktopCreateOrgHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return desktopOrgMutationHandler(dataService, true)
}

// DesktopUpdateOrgHandler updates an organization.
func DesktopUpdateOrgHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return desktopOrgMutationHandler(dataService, false)
}

func desktopOrgMutationHandler(dataService DesktopDataUseCase, create bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopOrgRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		var result service.DesktopOrg
		var err error
		if create {
			result, err = dataService.CreateOrg(r.Context(), request)
		} else {
			result, err = dataService.UpdateOrg(r.Context(), request)
		}
		if err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopDeleteOrgHandler deletes an organization.
func DesktopDeleteOrgHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopMessageActionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		if err := dataService.DeleteOrg(r.Context(), request.ID); err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

// DesktopCreateTenantHandler creates a tenant.
func DesktopCreateTenantHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return desktopTenantMutationHandler(dataService, true)
}

// DesktopUpdateTenantHandler updates a tenant.
func DesktopUpdateTenantHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return desktopTenantMutationHandler(dataService, false)
}

func desktopTenantMutationHandler(dataService DesktopDataUseCase, create bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopTenantRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		var result service.DesktopTenant
		var err error
		if create {
			result, err = dataService.CreateTenant(r.Context(), request)
		} else {
			result, err = dataService.UpdateTenant(r.Context(), request)
		}
		if err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopDeleteTenantHandler deletes a tenant.
func DesktopDeleteTenantHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopMessageActionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		if err := dataService.DeleteTenant(r.Context(), request.ID); err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

// DesktopUsersHandler returns an HTTP handler that lists desktop users.
func DesktopUsersHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		page, pageSize := parsePagination(r)
		result, err := dataService.ListUsers(r.Context(), r.URL.Query().Get("department"), r.URL.Query().Get("uid"), r.URL.Query().Get("name"), r.URL.Query().Get("role"), page, pageSize)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, response.CodeInternal, "failed to load desktop users")
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

func DesktopCreateUserHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return desktopUserMutationHandler(dataService, true)
}
func DesktopUpdateUserHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return desktopUserMutationHandler(dataService, false)
}
func desktopUserMutationHandler(dataService DesktopDataUseCase, create bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopUserRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		var result service.DesktopUserListItem
		var err error
		if create {
			result, err = dataService.CreateUser(r.Context(), request)
		} else {
			result, err = dataService.UpdateUser(r.Context(), request)
		}
		if err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}
func DesktopDeleteUserHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopMessageActionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		if err := dataService.DeleteUser(r.Context(), request.ID); err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

func DesktopResetUserPasswordHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopMessageActionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		if err := dataService.ResetUserPassword(r.Context(), request.ID); err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

// DesktopRolesHandler returns an HTTP handler that lists desktop roles.
func DesktopRolesHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		page, pageSize := parsePagination(r)
		result, err := dataService.ListRoles(r.Context(), r.URL.Query().Get("name"), r.URL.Query().Get("key"), r.URL.Query().Get("status"), page, pageSize)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, response.CodeInternal, "failed to load desktop roles")
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

func DesktopCreateRoleHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return desktopRoleMutationHandler(dataService, true)
}
func DesktopUpdateRoleHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return desktopRoleMutationHandler(dataService, false)
}
func desktopRoleMutationHandler(dataService DesktopDataUseCase, create bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		var result service.DesktopRole
		var err error
		if create {
			result, err = dataService.CreateRole(r.Context(), request)
		} else {
			result, err = dataService.UpdateRole(r.Context(), request)
		}
		if err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}
func DesktopDeleteRoleHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopMessageActionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		if err := dataService.DeleteRole(r.Context(), request.ID); err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

func DesktopRoleDataPermissionHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		roleID, _ := strconv.ParseInt(r.URL.Query().Get("roleId"), 10, 64)
		result, err := dataService.GetRoleDataPermission(r.Context(), roleID)
		if err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

func DesktopSaveRoleDataPermissionHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopRoleDataPermissionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		result, err := dataService.SaveRoleDataPermission(r.Context(), request)
		if err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopSystemSettingsHandler returns the desktop system settings.
func DesktopSystemSettingsHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		result, err := dataService.GetSystemSettings(r.Context())
		if err != nil {
			response.Error(w, http.StatusInternalServerError, response.CodeInternal, "failed to load desktop system settings")
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopSaveSystemSettingsHandler saves the desktop system settings.
func DesktopSaveSystemSettingsHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request service.DesktopSystemSettings
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		result, err := dataService.SaveSystemSettings(r.Context(), request)
		if err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopDictsHandler returns an HTTP handler that lists desktop dictionaries.
func DesktopDictsHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		result, err := dataService.ListDicts(r.Context(), r.URL.Query().Get("name"), r.URL.Query().Get("dictId"))
		if err != nil {
			response.Error(w, http.StatusInternalServerError, response.CodeInternal, "failed to load desktop dicts")
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

func DesktopCreateDictHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return desktopDictMutationHandler(dataService, true)
}
func DesktopUpdateDictHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return desktopDictMutationHandler(dataService, false)
}
func desktopDictMutationHandler(dataService DesktopDataUseCase, create bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopDictRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		var result service.DesktopDict
		var err error
		if create {
			result, err = dataService.CreateDict(r.Context(), request)
		} else {
			result, err = dataService.UpdateDict(r.Context(), request)
		}
		if err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}
func DesktopDeleteDictHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopMessageActionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		if err := dataService.DeleteDict(r.Context(), request.ID); err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}
func DesktopCreateDictItemHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return desktopDictItemMutationHandler(dataService, true)
}
func DesktopUpdateDictItemHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return desktopDictItemMutationHandler(dataService, false)
}
func desktopDictItemMutationHandler(dataService DesktopDataUseCase, create bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopDictItemRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		var result service.DesktopDictItem
		var err error
		if create {
			result, err = dataService.CreateDictItem(r.Context(), request.DictID, request.DesktopDictItem)
		} else {
			result, err = dataService.UpdateDictItem(r.Context(), request.DesktopDictItem)
		}
		if err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}
func DesktopDeleteDictItemHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopMessageActionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		if err := dataService.DeleteDictItem(r.Context(), request.ID); err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

// DesktopLoginLogsHandler returns an HTTP handler that lists desktop login logs.
func DesktopLoginLogsHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		page, pageSize := parsePagination(r)
		result, err := dataService.ListLoginLogs(r.Context(), r.URL.Query().Get("name"), r.URL.Query().Get("status"), r.URL.Query().Get("startDate"), r.URL.Query().Get("endDate"), page, pageSize)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, response.CodeInternal, "failed to load desktop login logs")
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopOperationLogsHandler returns an HTTP handler that lists desktop operation logs.
func DesktopOperationLogsHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		page, pageSize := parsePagination(r)
		result, err := dataService.ListOperationLogs(r.Context(), r.URL.Query().Get("name"), r.URL.Query().Get("status"), r.URL.Query().Get("startDate"), r.URL.Query().Get("endDate"), page, pageSize)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, response.CodeInternal, "failed to load desktop operation logs")
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopAnnouncementsHandler returns an HTTP handler that lists desktop announcements.
func DesktopAnnouncementsHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		page, pageSize := parsePagination(r)
		result, err := dataService.ListAnnouncements(r.Context(), r.URL.Query().Get("title"), r.URL.Query().Get("type"), r.URL.Query().Get("status"), page, pageSize)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, response.CodeInternal, "failed to load desktop announcements")
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

type desktopMessageActionRequest struct {
	ID int64 `json:"id"`
}

type desktopMessagesBatchActionRequest struct {
	IDs []int64 `json:"ids"`
}

type desktopMonitorActionRequest struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

type desktopMenuTemplateRequest = service.DesktopMenuTemplate
type desktopAnnouncementRequest = service.DesktopAnnouncement

// DesktopMessagesHandler returns an HTTP handler that lists desktop messages.
func DesktopMessagesHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		page, pageSize := parsePagination(r)
		result, err := dataService.ListMessages(r.Context(), r.URL.Query().Get("type"), r.URL.Query().Get("q"), page, pageSize)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, response.CodeInternal, "failed to load desktop messages")
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopCreateMenuTemplateHandler creates a menu template.
func DesktopCreateMenuTemplateHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return desktopMenuTemplateMutationHandler(dataService, true)
}

// DesktopUpdateMenuTemplateHandler updates a menu template.
func DesktopUpdateMenuTemplateHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return desktopMenuTemplateMutationHandler(dataService, false)
}

func desktopMenuTemplateMutationHandler(dataService DesktopDataUseCase, create bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopMenuTemplateRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		var (
			result service.DesktopMenuTemplate
			err    error
		)
		if create {
			result, err = dataService.CreateMenuTemplate(r.Context(), request)
		} else {
			result, err = dataService.UpdateMenuTemplate(r.Context(), request)
		}
		if err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopDeleteMenuTemplateHandler deletes a menu template.
func DesktopDeleteMenuTemplateHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopMessageActionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		if err := dataService.DeleteMenuTemplate(r.Context(), request.ID); err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

// DesktopCreateAnnouncementHandler creates an announcement.
func DesktopCreateAnnouncementHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return desktopAnnouncementMutationHandler(dataService, true)
}

// DesktopUpdateAnnouncementHandler updates an announcement.
func DesktopUpdateAnnouncementHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return desktopAnnouncementMutationHandler(dataService, false)
}

func desktopAnnouncementMutationHandler(dataService DesktopDataUseCase, create bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopAnnouncementRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		var (
			result service.DesktopAnnouncement
			err    error
		)
		if create {
			result, err = dataService.CreateAnnouncement(r.Context(), request)
		} else {
			result, err = dataService.UpdateAnnouncement(r.Context(), request)
		}
		if err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopDeleteAnnouncementHandler deletes an announcement.
func DesktopDeleteAnnouncementHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopMessageActionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		if err := dataService.DeleteAnnouncement(r.Context(), request.ID); err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

// DesktopPublishAnnouncementHandler publishes an announcement.
func DesktopPublishAnnouncementHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopMessageActionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		if err := dataService.PublishAnnouncement(r.Context(), request.ID); err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

// DesktopMonitorHandler returns an HTTP handler that lists desktop monitor data.
func DesktopMonitorHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		page, pageSize := parsePagination(r)
		result, err := dataService.ListMonitor(r.Context(), r.URL.Query().Get("metric"), r.URL.Query().Get("level"), page, pageSize)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, response.CodeInternal, "failed to load desktop monitor data")
			return
		}
		response.JSON(w, http.StatusOK, result)
	}
}

// DesktopMarkMessageReadHandler marks a desktop message as read.
func DesktopMarkMessageReadHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopMessageActionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		if err := dataService.MarkMessageRead(r.Context(), request.ID); err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

// DesktopMarkMessagesReadHandler marks multiple desktop messages as read.
func DesktopMarkMessagesReadHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopMessagesBatchActionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		if err := dataService.MarkMessagesRead(r.Context(), request.IDs); err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

// DesktopUpdateMonitorStatusHandler updates a desktop monitor item status.
func DesktopUpdateMonitorStatusHandler(dataService DesktopDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dataService == nil {
			response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "desktop data service is not ready")
			return
		}
		if _, err := parseDesktopBearerToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
			return
		}
		var request desktopMonitorActionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, "request body must be valid JSON")
			return
		}
		if err := dataService.UpdateMonitorStatus(r.Context(), request.ID, request.Status); err != nil {
			handleDesktopDataError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

func handleDesktopDataError(w http.ResponseWriter, err error) {
	if errors.Is(err, service.ErrInvalidInput) {
		response.Error(w, http.StatusBadRequest, response.CodeInvalidRequest, err.Error())
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		response.Error(w, http.StatusNotFound, response.CodeNotFound, "resource not found")
		return
	}
	response.Error(w, http.StatusInternalServerError, response.CodeInternal, "internal server error")
}

func parsePagination(r *http.Request) (int, int) {
	page := 1
	pageSize := 10
	if value := r.URL.Query().Get("page"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if value := r.URL.Query().Get("pageSize"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}
	return page, pageSize
}
