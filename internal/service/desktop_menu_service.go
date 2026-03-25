package service

import (
	"context"
	"fmt"
	"strings"

	"ai-go-service/internal/integration/goadmin"
)

// GoAdminMenuClient defines the upstream menu behavior required by DesktopMenuService.
type GoAdminMenuClient interface {
	GetCurrentMenus(ctx context.Context, accessToken string) ([]goadmin.MenuResult, error)
}

// DesktopMenuService adapts go-admin menu responses to the desktop-facing API.
type DesktopMenuService struct {
	client GoAdminMenuClient
	repository DesktopMenuRepository
}

// DesktopMenuRepository defines persistence behavior for menu management.
type DesktopMenuRepository interface {
	ListDesktopMenus(ctx context.Context, name string, menuType string, status string) ([]DesktopMenuRecord, error)
	CreateDesktopMenu(ctx context.Context, input DesktopMenuRecord) (DesktopMenuRecord, error)
	UpdateDesktopMenu(ctx context.Context, input DesktopMenuRecord) (DesktopMenuRecord, error)
	DeleteDesktopMenu(ctx context.Context, id int64) error
}

// DesktopMenu is the stable desktop-facing menu DTO.
type DesktopMenu struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	ParentID      int64  `json:"parentId"`
	Path          string `json:"path"`
	Icon          string `json:"icon"`
	Sort          int    `json:"sort"`
	Level         int    `json:"level"`
	ComponentName string `json:"componentName"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

// DesktopMenusResponse is returned when the desktop loads current menus.
type DesktopMenusResponse struct {
	List []DesktopMenu `json:"list"`
}

// NewDesktopMenuService creates a desktop menu adapter service.
func NewDesktopMenuService(client GoAdminMenuClient, repositories ...DesktopMenuRepository) *DesktopMenuService {
	var repository DesktopMenuRepository
	if len(repositories) > 0 {
		repository = repositories[0]
	}
	return &DesktopMenuService{client: client, repository: repository}
}

// CurrentMenus resolves the current desktop user's menus through go-admin.
func (service *DesktopMenuService) CurrentMenus(ctx context.Context, accessToken string) (DesktopMenusResponse, error) {
	if service == nil || service.client == nil {
		return DesktopMenusResponse{}, ErrUnavailable
	}

	if strings.TrimSpace(accessToken) == "" {
		return DesktopMenusResponse{}, fmt.Errorf("access token is required: %w", ErrInvalidInput)
	}

	menus, err := service.client.GetCurrentMenus(ctx, accessToken)
	if err != nil {
		return DesktopMenusResponse{}, err
	}

	result := make([]DesktopMenu, 0, len(menus))
	for _, menu := range menus {
		result = append(result, DesktopMenu{
			ID:            menu.ID,
			Name:          menu.Name,
			ParentID:      menu.ParentID,
			Path:          menu.Path,
			Icon:          menu.WebIcon,
			Sort:          menu.Sort,
			Level:         menu.Level,
			ComponentName: menu.ComponentName,
			CreatedAt:     menu.CreatedAt,
			UpdatedAt:     menu.UpdatedAt,
		})
	}

	return DesktopMenusResponse{List: result}, nil
}

// ListDesktopMenus returns menu-management rows from ai-go-service persistence.
func (service *DesktopMenuService) ListDesktopMenus(ctx context.Context, name string, menuType string, status string) (DesktopManagedMenusResponse, error) {
	if service == nil || service.repository == nil { return DesktopManagedMenusResponse{Items: fallbackManagedMenus()}, nil }
	items, err := service.repository.ListDesktopMenus(ctx, name, menuType, status)
	if err != nil { return DesktopManagedMenusResponse{}, err }
	result := mapManagedMenus(items)
	if len(result) == 0 { return DesktopManagedMenusResponse{Items: fallbackManagedMenus()}, nil }
	return DesktopManagedMenusResponse{Items: result}, nil
}

// CreateDesktopMenu creates one managed menu.
func (service *DesktopMenuService) CreateDesktopMenu(ctx context.Context, input DesktopManagedMenu) (DesktopManagedMenu, error) {
	if service == nil || service.repository == nil { return input, nil }
	if strings.TrimSpace(input.Name) == "" { return DesktopManagedMenu{}, fmt.Errorf("name is required: %w", ErrInvalidInput) }
	record, err := service.repository.CreateDesktopMenu(ctx, DesktopMenuRecord{ParentID: input.ParentID, Name: input.Name, Level: parseMenuLevel(input.Level), Sort: input.Order, Type: input.Type, Icon: input.Icon, Status: input.Status, Path: input.Path, Permission: input.Perm})
	if err != nil { return DesktopManagedMenu{}, err }
	return mapManagedMenus([]DesktopMenuRecord{record})[0], nil
}

// UpdateDesktopMenu updates one managed menu.
func (service *DesktopMenuService) UpdateDesktopMenu(ctx context.Context, input DesktopManagedMenu) (DesktopManagedMenu, error) {
	if service == nil || service.repository == nil { return input, nil }
	if input.ID <= 0 || strings.TrimSpace(input.Name) == "" { return DesktopManagedMenu{}, fmt.Errorf("id and name are required: %w", ErrInvalidInput) }
	record, err := service.repository.UpdateDesktopMenu(ctx, DesktopMenuRecord{ID: input.ID, ParentID: input.ParentID, Name: input.Name, Level: parseMenuLevel(input.Level), Sort: input.Order, Type: input.Type, Icon: input.Icon, Status: input.Status, Path: input.Path, Permission: input.Perm})
	if err != nil { return DesktopManagedMenu{}, err }
	return mapManagedMenus([]DesktopMenuRecord{record})[0], nil
}

// DeleteDesktopMenu removes one managed menu.
func (service *DesktopMenuService) DeleteDesktopMenu(ctx context.Context, id int64) error {
	if service == nil || service.repository == nil { return nil }
	if id <= 0 { return fmt.Errorf("id is required: %w", ErrInvalidInput) }
	return service.repository.DeleteDesktopMenu(ctx, id)
}
