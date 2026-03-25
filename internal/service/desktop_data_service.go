package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// DesktopMessageRecord is the service-layer message record shape.
type DesktopMessageRecord struct {
	ID          int64
	Type        string
	Title       string
	Description string
	PublishedAt time.Time
	Author      string
	Read        bool
}

// DesktopMonitorRecord is the service-layer monitor record shape.
type DesktopMonitorRecord struct {
	ID          int64
	Metric      string
	Value       string
	Alarm       bool
	Level       string
	OccurredAt  time.Time
	Description string
	Status      string
}

// DesktopMonitorSummaryRecord aggregates monitor statistics.
type DesktopMonitorSummaryRecord struct {
	Total            int64
	ActiveAlarmCount int64
	HighestLevel     string
	LastCollectedAt  time.Time
}

// DesktopAnnouncementRecord is the service-layer announcement record shape.
type DesktopAnnouncementRecord struct {
	ID          int64
	Title       string
	Type        string
	Status      string
	PublishedAt time.Time
	PublishedBy string
	CreatedAt   time.Time
	CreatedBy   string
}

// DesktopLoginLogRecord is the service-layer login log record shape.
type DesktopLoginLogRecord struct {
	ID         int64
	LogID      string
	TenantCode string
	Category   string
	UserID     string
	Name       string
	Status     string
	Time       time.Time
	IP         string
	Address    string
	Browser    string
	Desc       string
}

// DesktopOperationLogRecord is the service-layer operation log record shape.
type DesktopOperationLogRecord struct {
	ID         int64
	LogID      string
	TenantCode string
	Module     string
	Category   string
	UserID     string
	Name       string
	Status     string
	Time       time.Time
	IP         string
	Browser    string
	Desc       string
}

// DesktopDictRecord is the service-layer dictionary record shape.
type DesktopDictRecord struct {
	ID          int64
	DictID      string
	Name        string
	Description string
	ModifiedAt  time.Time
	ModifiedBy  string
	CreatedAt   time.Time
	CreatedBy   string
	Items       []DesktopDictItemRecord
}

// DesktopDictItemRecord is the service-layer dictionary item record shape.
type DesktopDictItemRecord struct {
	ID         int64
	Index      int
	Label      string
	KeyVal     string
	StyleType  string
	ModifiedAt time.Time
	ModifiedBy string
	CreatedAt  time.Time
	CreatedBy  string
	BadgeColor string
}

// DesktopSystemSettingsRecord is the service-layer system settings record shape.
type DesktopSystemSettingsRecord struct {
	WebsiteTitle          string
	SystemLogo            string
	Theme                 string
	ICP                   string
	Copyright             string
	RequireStrongPassword bool
	LoginFailLimit        int
	LoginLockMinutes      int
}

// DesktopDepartmentRecord is the service-layer department summary shape.
type DesktopDepartmentRecord struct {
	Name  string
	Count int
}

// DesktopUserRecord is the service-layer user record shape.
type DesktopUserRecord struct {
	ID         int64
	UID        string
	Name       string
	Department string
	Phone      string
	Role       string
	Status     string
	UpdatedAt  time.Time
	UpdatedBy  string
}

// DesktopRoleRecord is the service-layer role record shape.
type DesktopRoleRecord struct {
	ID        int64
	Name      string
	Key       string
	Order     int
	Status    bool
	CreatedAt time.Time
}

// DesktopRoleDataPermissionRecord is the service-layer role data permission shape.
type DesktopRoleDataPermissionRecord struct {
	RoleID      int64
	Scope       string
	Departments string
}

// DesktopOrgRecord is the service-layer organization record shape.
type DesktopOrgRecord struct {
	ID         int64
	Name       string
	Code       string
	ParentName string
	Sort       int
	Category   string
	CreatedAt  time.Time
	CreatedBy  string
	UpdatedAt  time.Time
	UpdatedBy  string
}

// DesktopTenantRecord is the service-layer tenant record shape.
type DesktopTenantRecord struct {
	ID     int64
	Code   string
	Name   string
	Status string
	Period string
	Admin  string
	Phone  string
}

// DesktopMenuTemplateRecord is the service-layer menu template record shape.
type DesktopMenuTemplateRecord struct {
	ID          int64
	Name        string
	TenantCount int
	Description string
	UpdatedAt   time.Time
	UpdatedBy   string
	CreatedAt   time.Time
	CreatedBy   string
}

// DesktopMenuRecord is the service-layer menu record shape.
type DesktopMenuRecord struct {
	ID         int64
	ParentID   int64
	Name       string
	Level      int
	Sort       int
	Type       string
	Icon       string
	Status     string
	Path       string
	Permission string
}

// DesktopSystemStatsSource collects the raw entities needed to build system statistics.
type DesktopSystemStatsSource struct {
	MonitorEvents []DesktopMonitorRecord
	Tenants       []DesktopTenantRecord
	Users         []DesktopUserRecord
}

// DesktopUsageStatsSource collects the raw entities needed to build usage statistics.
type DesktopUsageStatsSource struct {
	LoginLogs     []DesktopLoginLogRecord
	OperationLogs []DesktopOperationLogRecord
	Tenants       []DesktopTenantRecord
}

// DesktopDataRepository defines the persistence behavior needed for desktop data pages.
type DesktopDataRepository interface {
	EnsureSchemaAndSeed(ctx context.Context) error
	ListMessages(ctx context.Context, messageType string, keyword string, page int, pageSize int) ([]DesktopMessageRecord, int, error)
	ListMonitorEvents(ctx context.Context, metric string, level string, page int, pageSize int) ([]DesktopMonitorRecord, int, error)
	ListAnnouncements(ctx context.Context, title string, announcementType string, status string, page int, pageSize int) ([]DesktopAnnouncementRecord, int, error)
	CreateAnnouncement(ctx context.Context, input DesktopAnnouncementRecord) (DesktopAnnouncementRecord, error)
	UpdateAnnouncement(ctx context.Context, input DesktopAnnouncementRecord) (DesktopAnnouncementRecord, error)
	DeleteAnnouncement(ctx context.Context, id int64) error
	PublishAnnouncement(ctx context.Context, id int64) error
	ListLoginLogs(ctx context.Context, name string, status string, startDate string, endDate string, page int, pageSize int) ([]DesktopLoginLogRecord, int, error)
	ListOperationLogs(ctx context.Context, name string, status string, startDate string, endDate string, page int, pageSize int) ([]DesktopOperationLogRecord, int, error)
	ListDicts(ctx context.Context, name string, dictID string) ([]DesktopDictRecord, error)
	CreateDict(ctx context.Context, input DesktopDictRecord) (DesktopDictRecord, error)
	UpdateDict(ctx context.Context, input DesktopDictRecord) (DesktopDictRecord, error)
	DeleteDict(ctx context.Context, id int64) error
	CreateDictItem(ctx context.Context, dictID int64, input DesktopDictItemRecord) (DesktopDictItemRecord, error)
	UpdateDictItem(ctx context.Context, input DesktopDictItemRecord) (DesktopDictItemRecord, error)
	DeleteDictItem(ctx context.Context, id int64) error
	GetSystemSettings(ctx context.Context) (DesktopSystemSettingsRecord, error)
	SaveSystemSettings(ctx context.Context, input DesktopSystemSettingsRecord) (DesktopSystemSettingsRecord, error)
	ListUserDepartments(ctx context.Context) ([]DesktopDepartmentRecord, error)
	ListOrgDepartments(ctx context.Context) ([]DesktopDepartmentRecord, error)
	ListUsers(ctx context.Context, department string, uid string, name string, role string, page int, pageSize int) ([]DesktopUserRecord, int, error)
	CreateUser(ctx context.Context, input DesktopUserRecord) (DesktopUserRecord, error)
	UpdateUser(ctx context.Context, input DesktopUserRecord) (DesktopUserRecord, error)
	DeleteUser(ctx context.Context, id int64) error
	ResetUserPassword(ctx context.Context, id int64, updatedBy string, updatedAt time.Time) error
	ListRoles(ctx context.Context, name string, key string, status string, page int, pageSize int) ([]DesktopRoleRecord, int, error)
	CreateRole(ctx context.Context, input DesktopRoleRecord) (DesktopRoleRecord, error)
	UpdateRole(ctx context.Context, input DesktopRoleRecord) (DesktopRoleRecord, error)
	DeleteRole(ctx context.Context, id int64) error
	GetRoleDataPermission(ctx context.Context, roleID int64) (DesktopRoleDataPermissionRecord, error)
	SaveRoleDataPermission(ctx context.Context, input DesktopRoleDataPermissionRecord) (DesktopRoleDataPermissionRecord, error)
	ListOrgs(ctx context.Context, department string, name string, code string, page int, pageSize int) ([]DesktopOrgRecord, int, error)
	CreateOrg(ctx context.Context, input DesktopOrgRecord) (DesktopOrgRecord, error)
	UpdateOrg(ctx context.Context, input DesktopOrgRecord) (DesktopOrgRecord, error)
	DeleteOrg(ctx context.Context, id int64) error
	ListTenants(ctx context.Context, name string, code string, status string, page int, pageSize int) ([]DesktopTenantRecord, int, error)
	CreateTenant(ctx context.Context, input DesktopTenantRecord) (DesktopTenantRecord, error)
	UpdateTenant(ctx context.Context, input DesktopTenantRecord) (DesktopTenantRecord, error)
	DeleteTenant(ctx context.Context, id int64) error
	ListMenuTemplates(ctx context.Context, name string, description string, page int, pageSize int) ([]DesktopMenuTemplateRecord, int, error)
	CreateMenuTemplate(ctx context.Context, input DesktopMenuTemplateRecord) (DesktopMenuTemplateRecord, error)
	UpdateMenuTemplate(ctx context.Context, input DesktopMenuTemplateRecord) (DesktopMenuTemplateRecord, error)
	DeleteMenuTemplate(ctx context.Context, id int64) error
	ListDesktopMenus(ctx context.Context, name string, menuType string, status string) ([]DesktopMenuRecord, error)
	CreateDesktopMenu(ctx context.Context, input DesktopMenuRecord) (DesktopMenuRecord, error)
	UpdateDesktopMenu(ctx context.Context, input DesktopMenuRecord) (DesktopMenuRecord, error)
	DeleteDesktopMenu(ctx context.Context, id int64) error
	GetSystemStatsSource(ctx context.Context) (DesktopSystemStatsSource, error)
	GetUsageStatsSource(ctx context.Context) (DesktopUsageStatsSource, error)
	GetMonitorSummary(ctx context.Context) (DesktopMonitorSummaryRecord, error)
	CollectMonitor(ctx context.Context, collectedAt time.Time) error
	MarkMessageRead(ctx context.Context, id int64) error
	MarkMessagesRead(ctx context.Context, ids []int64) error
	UpdateMonitorStatus(ctx context.Context, id int64, status string) error
}

// DesktopDataService serves desktop messages and monitor data.
type DesktopDataService struct {
	repository         DesktopDataRepository
	mu                 sync.RWMutex
	messages           []DesktopMessage
	monitor            []DesktopMonitorItem
	monitorCollectedAt time.Time
}

// DesktopMessage is the stable desktop-facing message DTO.
type DesktopMessage struct {
	ID          int64  `json:"id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PublishedAt string `json:"publishedAt"`
	Author      string `json:"author"`
	Read        bool   `json:"read"`
}

// DesktopMessagesResponse wraps desktop messages with simple metadata.
type DesktopMessagesResponse struct {
	Items    []DesktopMessage `json:"items"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}

// DesktopMonitorSummary is the stable desktop-facing monitor summary DTO.
type DesktopMonitorSummary struct {
	Total            int64  `json:"total"`
	ActiveAlarmCount int64  `json:"activeAlarmCount"`
	HighestLevel     string `json:"highestLevel"`
	LastCollectedAt  string `json:"lastCollectedAt"`
}

// DesktopMonitorItem is the stable desktop-facing monitor row DTO.
type DesktopMonitorItem struct {
	ID          int64  `json:"id"`
	Metric      string `json:"metric"`
	Value       string `json:"value"`
	Alarm       bool   `json:"alarm"`
	Level       string `json:"level"`
	OccurredAt  string `json:"occurredAt"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// DesktopMonitorResponse wraps monitor rows and summary.
type DesktopMonitorResponse struct {
	Summary  DesktopMonitorSummary `json:"summary"`
	Items    []DesktopMonitorItem  `json:"items"`
	Total    int                   `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"pageSize"`
}

// DesktopAnnouncement is the stable desktop-facing announcement DTO.
type DesktopAnnouncement struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Publish string `json:"publish"`
	Create  string `json:"create"`
}

// DesktopAnnouncementsResponse wraps announcements with paging metadata.
type DesktopAnnouncementsResponse struct {
	Items    []DesktopAnnouncement `json:"items"`
	Total    int                   `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"pageSize"`
}

// DesktopLoginLog is the stable desktop-facing login log DTO.
type DesktopLoginLog struct {
	ID       int64  `json:"id"`
	LogID    string `json:"logId"`
	Category string `json:"category"`
	UserID   string `json:"userId"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Time     string `json:"time"`
	IP       string `json:"ip"`
	Address  string `json:"address"`
	Browser  string `json:"browser"`
	Desc     string `json:"desc"`
}

// DesktopLoginLogsResponse wraps login logs with paging metadata.
type DesktopLoginLogsResponse struct {
	Items    []DesktopLoginLog `json:"items"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}

// DesktopOperationLog is the stable desktop-facing operation log DTO.
type DesktopOperationLog struct {
	ID       int64  `json:"id"`
	LogID    string `json:"logId"`
	Module   string `json:"module"`
	Category string `json:"category"`
	UserID   string `json:"userId"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Time     string `json:"time"`
	IP       string `json:"ip"`
	Browser  string `json:"browser"`
	Desc     string `json:"desc"`
}

// DesktopOperationLogsResponse wraps operation logs with paging metadata.
type DesktopOperationLogsResponse struct {
	Items    []DesktopOperationLog `json:"items"`
	Total    int                   `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"pageSize"`
}

// DesktopDictItem is the stable desktop-facing dictionary item DTO.
type DesktopDictItem struct {
	ID         int64  `json:"id"`
	Index      int    `json:"index"`
	Label      string `json:"label"`
	KeyVal     string `json:"keyVal"`
	StyleType  string `json:"styleType"`
	Modified   string `json:"modified"`
	Created    string `json:"created"`
	BadgeColor string `json:"badgeColor"`
}

// DesktopDict is the stable desktop-facing dictionary DTO.
type DesktopDict struct {
	ID       int64             `json:"id"`
	DictID   string            `json:"dictId"`
	Name     string            `json:"name"`
	Desc     string            `json:"desc"`
	Modified string            `json:"modified"`
	Created  string            `json:"created"`
	Items    []DesktopDictItem `json:"items"`
}

// DesktopDictsResponse wraps dictionaries.
type DesktopDictsResponse struct {
	Items []DesktopDict `json:"items"`
	Total int           `json:"total"`
}

// DesktopDepartment is the stable desktop-facing department DTO.
type DesktopDepartment struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// DesktopDepartmentsResponse wraps departments.
type DesktopDepartmentsResponse struct {
	Items []DesktopDepartment `json:"items"`
}

// DesktopUserListItem is the stable desktop-facing user-list DTO.
type DesktopUserListItem struct {
	ID     int64  `json:"id"`
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Dept   string `json:"dept"`
	Phone  string `json:"phone"`
	Role   string `json:"role"`
	Status string `json:"status"`
	Date   string `json:"date"`
}

// DesktopUsersResponse wraps users with paging metadata.
type DesktopUsersResponse struct {
	Departments []DesktopDepartment   `json:"departments"`
	Items       []DesktopUserListItem `json:"items"`
	Total       int                   `json:"total"`
	Page        int                   `json:"page"`
	PageSize    int                   `json:"pageSize"`
}

// DesktopRole is the stable desktop-facing role DTO.
type DesktopRole struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Key    string `json:"key"`
	Order  int    `json:"order"`
	Status bool   `json:"status"`
	Date   string `json:"date"`
}

// DesktopRolesResponse wraps roles with paging metadata.
type DesktopRolesResponse struct {
	Items    []DesktopRole `json:"items"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
}

// DesktopRoleDataPermission is the stable desktop-facing role data permission DTO.
type DesktopRoleDataPermission struct {
	RoleID      int64    `json:"roleId"`
	Scope       string   `json:"scope"`
	Departments []string `json:"departments"`
}

// DesktopOrg is the stable desktop-facing organization DTO.
type DesktopOrg struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Parent1  string `json:"parent1"`
	Order    int    `json:"order"`
	Parent2  string `json:"parent2"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}

// DesktopOrgsResponse wraps organizations.
type DesktopOrgsResponse struct {
	Departments []DesktopDepartment `json:"departments"`
	Items       []DesktopOrg        `json:"items"`
	Total       int                 `json:"total"`
	Page        int                 `json:"page"`
	PageSize    int                 `json:"pageSize"`
}

// DesktopTenant is the stable desktop-facing tenant DTO.
type DesktopTenant struct {
	ID     int64  `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Period string `json:"period"`
	Admin  string `json:"admin"`
	Phone  string `json:"phone"`
}

// DesktopTenantsResponse wraps tenants.
type DesktopTenantsResponse struct {
	Items    []DesktopTenant `json:"items"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
}

// DesktopMenuTemplate is the stable desktop-facing menu-template DTO.
type DesktopMenuTemplate struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Tenants  int    `json:"tenants"`
	Desc     string `json:"desc"`
	Modified string `json:"modified"`
	Created  string `json:"created"`
}

// DesktopMenuTemplatesResponse wraps menu templates.
type DesktopMenuTemplatesResponse struct {
	Items    []DesktopMenuTemplate `json:"items"`
	Total    int                   `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"pageSize"`
}

// DesktopManagedMenu is the stable desktop-facing menu management DTO.
type DesktopManagedMenu struct {
	ID       int64  `json:"id"`
	ParentID int64  `json:"parentId"`
	Name     string `json:"name"`
	Level    string `json:"level"`
	Order    int    `json:"order"`
	Type     string `json:"type"`
	Icon     string `json:"icon"`
	Status   string `json:"status"`
	Path     string `json:"path"`
	Perm     string `json:"perm"`
	Depth    int    `json:"depth"`
	HasSub   bool   `json:"hasSub"`
}

// DesktopManagedMenusResponse wraps menu management rows.
type DesktopManagedMenusResponse struct {
	Items []DesktopManagedMenu `json:"items"`
}

// DesktopSystemStatsAlertCard is the system stats summary card DTO.
type DesktopSystemStatsAlertCard struct {
	Label  string `json:"label"`
	Value  int64  `json:"value"`
	Change string `json:"change"`
	Period string `json:"period"`
}

// DesktopSystemStatsTrendPoint is the alert trend point DTO.
type DesktopSystemStatsTrendPoint struct {
	Label string `json:"label"`
	Value int64  `json:"value"`
}

// DesktopSystemStatsTenantBucket is the tenant bucket DTO.
type DesktopSystemStatsTenantBucket struct {
	Count int      `json:"count"`
	Names []string `json:"names"`
}

// DesktopSystemStatsRoleStat is the role distribution DTO.
type DesktopSystemStatsRoleStat struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

// DesktopSystemStatsResponse is the system stats page DTO.
type DesktopSystemStatsResponse struct {
	AlertCards        []DesktopSystemStatsAlertCard  `json:"alertCards"`
	Trend             []DesktopSystemStatsTrendPoint `json:"trend"`
	TenantTotal       int                            `json:"tenantTotal"`
	TenantOpenCount   int                            `json:"tenantOpenCount"`
	TenantClosedCount int                            `json:"tenantClosedCount"`
	TenantNormal      DesktopSystemStatsTenantBucket `json:"tenantNormal"`
	TenantExpiring    DesktopSystemStatsTenantBucket `json:"tenantExpiring"`
	TenantExpired     DesktopSystemStatsTenantBucket `json:"tenantExpired"`
	EmployeeTotal     int                            `json:"employeeTotal"`
	RoleStats         []DesktopSystemStatsRoleStat   `json:"roleStats"`
	DepartmentStats   []DesktopDepartment            `json:"departmentStats"`
}

// DesktopUsageStatsTenantOption is one tenant selector option for usage stats.
type DesktopUsageStatsTenantOption struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// DesktopUsageStatsTrendPoint is one daily usage trend point.
type DesktopUsageStatsTrendPoint struct {
	Label      string `json:"label"`
	LoginCount int    `json:"loginCount"`
	UsageCount int    `json:"usageCount"`
}

// DesktopUsageStatsModuleRank is one module ranking row.
type DesktopUsageStatsModuleRank struct {
	Rank   int    `json:"rank"`
	Name   string `json:"name"`
	Count  int    `json:"count"`
	Change string `json:"change"`
}

// DesktopUsageStatsResponse is the usage stats page DTO.
type DesktopUsageStatsResponse struct {
	Tenants  []DesktopUsageStatsTenantOption `json:"tenants"`
	Trend    []DesktopUsageStatsTrendPoint   `json:"trend"`
	TopMost  []DesktopUsageStatsModuleRank   `json:"topMost"`
	TopLeast []DesktopUsageStatsModuleRank   `json:"topLeast"`
}

// DesktopSystemSettings is the stable desktop-facing system settings DTO.
type DesktopSystemSettings struct {
	WebsiteTitle          string `json:"websiteTitle"`
	SystemLogo            string `json:"systemLogo"`
	Theme                 string `json:"theme"`
	ICP                   string `json:"icp"`
	Copyright             string `json:"copyright"`
	RequireStrongPassword bool   `json:"requireStrongPassword"`
	LoginFailLimit        int    `json:"loginFailLimit"`
	LoginLockMinutes      int    `json:"loginLockMinutes"`
}

// NewDesktopDataService creates a desktop data service.
func NewDesktopDataService(repository DesktopDataRepository) *DesktopDataService {
	return &DesktopDataService{
		repository:         repository,
		messages:           fallbackMessages(),
		monitor:            fallbackMonitorItems(),
		monitorCollectedAt: time.Now(),
	}
}

// EnsureSeedData provisions required tables and inserts sample rows.
func (service *DesktopDataService) EnsureSeedData(ctx context.Context) error {
	if service == nil || service.repository == nil {
		return nil
	}

	return service.repository.EnsureSchemaAndSeed(ctx)
}

// ListMessages returns seeded or database-backed desktop messages.
func (service *DesktopDataService) ListMessages(ctx context.Context, messageType string, keyword string, page int, pageSize int) (DesktopMessagesResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	if service == nil || service.repository == nil {
		service.mu.RLock()
		defer service.mu.RUnlock()
		items := filterMessages(service.messages, messageType, keyword)
		pagedItems := paginateMessages(items, page, pageSize)
		return DesktopMessagesResponse{Items: pagedItems, Total: len(items), Page: page, PageSize: pageSize}, nil
	}

	items, total, err := service.repository.ListMessages(ctx, messageType, keyword, page, pageSize)
	if err != nil {
		return DesktopMessagesResponse{}, err
	}

	result := make([]DesktopMessage, 0, len(items))
	for _, item := range items {
		result = append(result, DesktopMessage{
			ID:          item.ID,
			Type:        item.Type,
			Title:       item.Title,
			Description: item.Description,
			PublishedAt: item.PublishedAt.Format(time.DateTime),
			Author:      item.Author,
			Read:        item.Read,
		})
	}

	if len(result) == 0 {
		service.mu.RLock()
		defer service.mu.RUnlock()
		fallback := filterMessages(service.messages, messageType, keyword)
		pagedFallback := paginateMessages(fallback, page, pageSize)
		return DesktopMessagesResponse{Items: pagedFallback, Total: len(fallback), Page: page, PageSize: pageSize}, nil
	}

	return DesktopMessagesResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

// ListMonitor returns monitor summary and rows.
func (service *DesktopDataService) ListMonitor(ctx context.Context, metric string, level string, page int, pageSize int) (DesktopMonitorResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	if service == nil || service.repository == nil {
		service.mu.RLock()
		defer service.mu.RUnlock()
		items := filterMonitorItems(service.monitor, metric, level)
		pagedItems := paginateMonitorItems(items, page, pageSize)
		summary := buildFallbackMonitorSummary(items)
		if !service.monitorCollectedAt.IsZero() {
			summary.LastCollectedAt = service.monitorCollectedAt.Format(time.DateTime)
		}
		return DesktopMonitorResponse{Summary: summary, Items: pagedItems, Total: len(items), Page: page, PageSize: pageSize}, nil
	}

	summary, err := service.repository.GetMonitorSummary(ctx)
	if err != nil {
		return DesktopMonitorResponse{}, err
	}

	items, total, err := service.repository.ListMonitorEvents(ctx, metric, level, page, pageSize)
	if err != nil {
		return DesktopMonitorResponse{}, err
	}

	result := make([]DesktopMonitorItem, 0, len(items))
	for _, item := range items {
		result = append(result, DesktopMonitorItem{
			ID:          item.ID,
			Metric:      item.Metric,
			Value:       item.Value,
			Alarm:       item.Alarm,
			Level:       item.Level,
			OccurredAt:  item.OccurredAt.Format(time.DateTime),
			Description: item.Description,
			Status:      item.Status,
		})
	}

	if len(result) == 0 {
		service.mu.RLock()
		defer service.mu.RUnlock()
		fallback := filterMonitorItems(service.monitor, metric, level)
		pagedFallback := paginateMonitorItems(fallback, page, pageSize)
		return DesktopMonitorResponse{Summary: buildFallbackMonitorSummary(fallback), Items: pagedFallback, Total: len(fallback), Page: page, PageSize: pageSize}, nil
	}

	return DesktopMonitorResponse{
		Summary: DesktopMonitorSummary{
			Total:            summary.Total,
			ActiveAlarmCount: summary.ActiveAlarmCount,
			HighestLevel:     summary.HighestLevel,
			LastCollectedAt:  summary.LastCollectedAt.Format(time.DateTime),
		},
		Items:    result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// CollectMonitor records one manual monitor collection timestamp.
func (service *DesktopDataService) CollectMonitor(ctx context.Context) error {
	if service == nil {
		return ErrUnavailable
	}

	collectedAt := time.Now()
	service.mu.Lock()
	service.monitorCollectedAt = collectedAt
	service.mu.Unlock()

	if service.repository == nil {
		return nil
	}

	return service.repository.CollectMonitor(ctx, collectedAt)
}

// ListAnnouncements returns announcements for the desktop page.
func (service *DesktopDataService) ListAnnouncements(ctx context.Context, title string, announcementType string, status string, page int, pageSize int) (DesktopAnnouncementsResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	if service == nil || service.repository == nil {
		return DesktopAnnouncementsResponse{Items: fallbackAnnouncements(), Total: len(fallbackAnnouncements()), Page: page, PageSize: pageSize}, nil
	}
	items, total, err := service.repository.ListAnnouncements(ctx, title, announcementType, status, page, pageSize)
	if err != nil {
		return DesktopAnnouncementsResponse{}, err
	}
	result := make([]DesktopAnnouncement, 0, len(items))
	for _, item := range items {
		result = append(result, DesktopAnnouncement{
			ID:      item.ID,
			Title:   item.Title,
			Type:    item.Type,
			Status:  item.Status,
			Publish: item.PublishedAt.Format(time.DateTime) + "/" + item.PublishedBy,
			Create:  item.CreatedAt.Format(time.DateTime) + "/" + item.CreatedBy,
		})
	}
	if len(result) == 0 {
		fallback := fallbackAnnouncements()
		return DesktopAnnouncementsResponse{Items: fallback, Total: len(fallback), Page: page, PageSize: pageSize}, nil
	}
	return DesktopAnnouncementsResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

// CreateAnnouncement creates one announcement.
func (service *DesktopDataService) CreateAnnouncement(ctx context.Context, input DesktopAnnouncement) (DesktopAnnouncement, error) {
	if input.Title == "" || input.Type == "" {
		return DesktopAnnouncement{}, fmt.Errorf("title and type are required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.CreateAnnouncement(ctx, DesktopAnnouncementRecord{Title: input.Title, Type: input.Type, Status: input.Status, PublishedAt: time.Now(), PublishedBy: "桌面端", CreatedAt: time.Now(), CreatedBy: "桌面端"})
	if err != nil {
		return DesktopAnnouncement{}, err
	}
	return DesktopAnnouncement{ID: record.ID, Title: record.Title, Type: record.Type, Status: record.Status, Publish: record.PublishedAt.Format(time.DateTime) + "/" + record.PublishedBy, Create: record.CreatedAt.Format(time.DateTime) + "/" + record.CreatedBy}, nil
}

// UpdateAnnouncement updates one announcement.
func (service *DesktopDataService) UpdateAnnouncement(ctx context.Context, input DesktopAnnouncement) (DesktopAnnouncement, error) {
	if input.ID <= 0 || input.Title == "" || input.Type == "" {
		return DesktopAnnouncement{}, fmt.Errorf("id, title and type are required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.UpdateAnnouncement(ctx, DesktopAnnouncementRecord{ID: input.ID, Title: input.Title, Type: input.Type, Status: input.Status, PublishedAt: time.Now(), PublishedBy: "桌面端", CreatedAt: time.Now(), CreatedBy: "桌面端"})
	if err != nil {
		return DesktopAnnouncement{}, err
	}
	return DesktopAnnouncement{ID: record.ID, Title: record.Title, Type: record.Type, Status: record.Status, Publish: record.PublishedAt.Format(time.DateTime) + "/" + record.PublishedBy, Create: record.CreatedAt.Format(time.DateTime) + "/" + record.CreatedBy}, nil
}

// DeleteAnnouncement removes one announcement.
func (service *DesktopDataService) DeleteAnnouncement(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return nil
	}
	return service.repository.DeleteAnnouncement(ctx, id)
}

// PublishAnnouncement marks one announcement as published.
func (service *DesktopDataService) PublishAnnouncement(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return nil
	}
	return service.repository.PublishAnnouncement(ctx, id)
}

// ListLoginLogs returns login logs for the desktop page.
func (service *DesktopDataService) ListLoginLogs(ctx context.Context, name string, status string, startDate string, endDate string, page int, pageSize int) (DesktopLoginLogsResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	if service == nil || service.repository == nil {
		fallback := fallbackLoginLogs()
		filtered := filterLoginLogs(fallback, name, status, startDate, endDate)
		return DesktopLoginLogsResponse{Items: paginateLoginLogs(filtered, page, pageSize), Total: len(filtered), Page: page, PageSize: pageSize}, nil
	}
	items, total, err := service.repository.ListLoginLogs(ctx, name, status, startDate, endDate, page, pageSize)
	if err != nil {
		return DesktopLoginLogsResponse{}, err
	}
	result := make([]DesktopLoginLog, 0, len(items))
	for _, item := range items {
		result = append(result, DesktopLoginLog{ID: item.ID, LogID: item.LogID, Category: item.Category, UserID: item.UserID, Name: item.Name, Status: item.Status, Time: item.Time.Format(time.DateTime), IP: item.IP, Address: item.Address, Browser: item.Browser, Desc: item.Desc})
	}
	if len(result) == 0 {
		fallback := filterLoginLogs(fallbackLoginLogs(), name, status, startDate, endDate)
		return DesktopLoginLogsResponse{Items: paginateLoginLogs(fallback, page, pageSize), Total: len(fallback), Page: page, PageSize: pageSize}, nil
	}
	return DesktopLoginLogsResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

// ListOperationLogs returns operation logs for the desktop page.
func (service *DesktopDataService) ListOperationLogs(ctx context.Context, name string, status string, startDate string, endDate string, page int, pageSize int) (DesktopOperationLogsResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	if service == nil || service.repository == nil {
		fallback := fallbackOperationLogs()
		filtered := filterOperationLogs(fallback, name, status, startDate, endDate)
		return DesktopOperationLogsResponse{Items: paginateOperationLogs(filtered, page, pageSize), Total: len(filtered), Page: page, PageSize: pageSize}, nil
	}
	items, total, err := service.repository.ListOperationLogs(ctx, name, status, startDate, endDate, page, pageSize)
	if err != nil {
		return DesktopOperationLogsResponse{}, err
	}
	result := make([]DesktopOperationLog, 0, len(items))
	for _, item := range items {
		result = append(result, DesktopOperationLog{ID: item.ID, LogID: item.LogID, Module: item.Module, Category: item.Category, UserID: item.UserID, Name: item.Name, Status: item.Status, Time: item.Time.Format(time.DateTime), IP: item.IP, Browser: item.Browser, Desc: item.Desc})
	}
	if len(result) == 0 {
		fallback := filterOperationLogs(fallbackOperationLogs(), name, status, startDate, endDate)
		return DesktopOperationLogsResponse{Items: paginateOperationLogs(fallback, page, pageSize), Total: len(fallback), Page: page, PageSize: pageSize}, nil
	}
	return DesktopOperationLogsResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

// ListDicts returns dictionaries for the desktop page.
func (service *DesktopDataService) ListDicts(ctx context.Context, name string, dictID string) (DesktopDictsResponse, error) {
	if service == nil || service.repository == nil {
		fallback := filterDicts(fallbackDicts(), name, dictID)
		return DesktopDictsResponse{Items: fallback, Total: len(fallback)}, nil
	}
	items, err := service.repository.ListDicts(ctx, name, dictID)
	if err != nil {
		return DesktopDictsResponse{}, err
	}
	result := make([]DesktopDict, 0, len(items))
	for _, item := range items {
		dictItems := make([]DesktopDictItem, 0, len(item.Items))
		for _, subItem := range item.Items {
			dictItems = append(dictItems, DesktopDictItem{
				ID: subItem.ID, Index: subItem.Index, Label: subItem.Label, KeyVal: subItem.KeyVal, StyleType: subItem.StyleType,
				Modified:   subItem.ModifiedAt.Format(time.DateTime) + "/" + subItem.ModifiedBy,
				Created:    subItem.CreatedAt.Format(time.DateTime) + "/" + subItem.CreatedBy,
				BadgeColor: subItem.BadgeColor,
			})
		}
		result = append(result, DesktopDict{
			ID: item.ID, DictID: item.DictID, Name: item.Name, Desc: item.Description,
			Modified: item.ModifiedAt.Format(time.DateTime) + "/" + item.ModifiedBy,
			Created:  item.CreatedAt.Format(time.DateTime) + "/" + item.CreatedBy,
			Items:    dictItems,
		})
	}
	if len(result) == 0 {
		fallback := filterDicts(fallbackDicts(), name, dictID)
		return DesktopDictsResponse{Items: fallback, Total: len(fallback)}, nil
	}
	return DesktopDictsResponse{Items: result, Total: len(result)}, nil
}

// CreateDict creates one dictionary.
func (service *DesktopDataService) CreateDict(ctx context.Context, input DesktopDict) (DesktopDict, error) {
	if input.DictID == "" || input.Name == "" {
		return DesktopDict{}, fmt.Errorf("dictId and name are required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.CreateDict(ctx, DesktopDictRecord{DictID: input.DictID, Name: input.Name, Description: input.Desc, ModifiedAt: time.Now(), ModifiedBy: "桌面端", CreatedAt: time.Now(), CreatedBy: "桌面端"})
	if err != nil {
		return DesktopDict{}, err
	}
	return DesktopDict{ID: record.ID, DictID: record.DictID, Name: record.Name, Desc: record.Description, Modified: record.ModifiedAt.Format(time.DateTime) + "/" + record.ModifiedBy, Created: record.CreatedAt.Format(time.DateTime) + "/" + record.CreatedBy, Items: []DesktopDictItem{}}, nil
}

// UpdateDict updates one dictionary.
func (service *DesktopDataService) UpdateDict(ctx context.Context, input DesktopDict) (DesktopDict, error) {
	if input.ID <= 0 || input.DictID == "" || input.Name == "" {
		return DesktopDict{}, fmt.Errorf("id, dictId and name are required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.UpdateDict(ctx, DesktopDictRecord{ID: input.ID, DictID: input.DictID, Name: input.Name, Description: input.Desc, ModifiedAt: time.Now(), ModifiedBy: "桌面端"})
	if err != nil {
		return DesktopDict{}, err
	}
	return DesktopDict{ID: record.ID, DictID: record.DictID, Name: record.Name, Desc: record.Description, Modified: record.ModifiedAt.Format(time.DateTime) + "/" + record.ModifiedBy, Created: record.CreatedAt.Format(time.DateTime) + "/" + record.CreatedBy, Items: input.Items}, nil
}

// DeleteDict removes one dictionary.
func (service *DesktopDataService) DeleteDict(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return nil
	}
	return service.repository.DeleteDict(ctx, id)
}

// CreateDictItem creates one dictionary item.
func (service *DesktopDataService) CreateDictItem(ctx context.Context, dictID int64, input DesktopDictItem) (DesktopDictItem, error) {
	if dictID <= 0 || input.Label == "" || input.KeyVal == "" {
		return DesktopDictItem{}, fmt.Errorf("dict id, label and keyVal are required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.CreateDictItem(ctx, dictID, DesktopDictItemRecord{Index: input.Index, Label: input.Label, KeyVal: input.KeyVal, StyleType: input.StyleType, ModifiedAt: time.Now(), ModifiedBy: "桌面端", CreatedAt: time.Now(), CreatedBy: "桌面端", BadgeColor: input.BadgeColor})
	if err != nil {
		return DesktopDictItem{}, err
	}
	return DesktopDictItem{ID: record.ID, Index: record.Index, Label: record.Label, KeyVal: record.KeyVal, StyleType: record.StyleType, Modified: record.ModifiedAt.Format(time.DateTime) + "/" + record.ModifiedBy, Created: record.CreatedAt.Format(time.DateTime) + "/" + record.CreatedBy, BadgeColor: record.BadgeColor}, nil
}

// UpdateDictItem updates one dictionary item.
func (service *DesktopDataService) UpdateDictItem(ctx context.Context, input DesktopDictItem) (DesktopDictItem, error) {
	if input.ID <= 0 || input.Label == "" || input.KeyVal == "" {
		return DesktopDictItem{}, fmt.Errorf("id, label and keyVal are required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.UpdateDictItem(ctx, DesktopDictItemRecord{ID: input.ID, Index: input.Index, Label: input.Label, KeyVal: input.KeyVal, StyleType: input.StyleType, ModifiedAt: time.Now(), ModifiedBy: "桌面端", BadgeColor: input.BadgeColor})
	if err != nil {
		return DesktopDictItem{}, err
	}
	return DesktopDictItem{ID: record.ID, Index: record.Index, Label: record.Label, KeyVal: record.KeyVal, StyleType: record.StyleType, Modified: record.ModifiedAt.Format(time.DateTime) + "/" + record.ModifiedBy, Created: record.CreatedAt.Format(time.DateTime) + "/" + record.CreatedBy, BadgeColor: record.BadgeColor}, nil
}

// DeleteDictItem removes one dictionary item.
func (service *DesktopDataService) DeleteDictItem(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return nil
	}
	return service.repository.DeleteDictItem(ctx, id)
}

// GetSystemSettings returns desktop system settings.
func (service *DesktopDataService) GetSystemSettings(ctx context.Context) (DesktopSystemSettings, error) {
	if service == nil || service.repository == nil {
		return fallbackSystemSettings(), nil
	}
	record, err := service.repository.GetSystemSettings(ctx)
	if err != nil {
		return DesktopSystemSettings{}, err
	}
	return mapSystemSettings(record), nil
}

// SaveSystemSettings stores desktop system settings.
func (service *DesktopDataService) SaveSystemSettings(ctx context.Context, input DesktopSystemSettings) (DesktopSystemSettings, error) {
	if input.WebsiteTitle == "" {
		return DesktopSystemSettings{}, fmt.Errorf("website title is required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.SaveSystemSettings(ctx, DesktopSystemSettingsRecord{
		WebsiteTitle:          input.WebsiteTitle,
		SystemLogo:            input.SystemLogo,
		Theme:                 input.Theme,
		ICP:                   input.ICP,
		Copyright:             input.Copyright,
		RequireStrongPassword: input.RequireStrongPassword,
		LoginFailLimit:        input.LoginFailLimit,
		LoginLockMinutes:      input.LoginLockMinutes,
	})
	if err != nil {
		return DesktopSystemSettings{}, err
	}
	return mapSystemSettings(record), nil
}

// ListUsers returns user data and department summary.
func (service *DesktopDataService) ListUsers(ctx context.Context, department string, uid string, name string, role string, page int, pageSize int) (DesktopUsersResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	if service == nil || service.repository == nil {
		users := filterUsers(fallbackUsers(), department, uid, name, role)
		return DesktopUsersResponse{Departments: fallbackDepartments(), Items: paginateUsers(users, page, pageSize), Total: len(users), Page: page, PageSize: pageSize}, nil
	}
	departments, err := service.repository.ListUserDepartments(ctx)
	if err != nil {
		return DesktopUsersResponse{}, err
	}
	items, total, err := service.repository.ListUsers(ctx, department, uid, name, role, page, pageSize)
	if err != nil {
		return DesktopUsersResponse{}, err
	}
	result := make([]DesktopUserListItem, 0, len(items))
	for _, item := range items {
		result = append(result, DesktopUserListItem{ID: item.ID, UID: item.UID, Name: item.Name, Dept: item.Department, Phone: item.Phone, Role: item.Role, Status: item.Status, Date: item.UpdatedAt.Format(time.DateTime) + "/" + item.UpdatedBy})
	}
	deptDTOs := make([]DesktopDepartment, 0, len(departments))
	for _, item := range departments {
		deptDTOs = append(deptDTOs, DesktopDepartment{Name: item.Name, Count: item.Count})
	}
	if len(result) == 0 {
		fallback := filterUsers(fallbackUsers(), department, uid, name, role)
		return DesktopUsersResponse{Departments: fallbackDepartments(), Items: paginateUsers(fallback, page, pageSize), Total: len(fallback), Page: page, PageSize: pageSize}, nil
	}
	return DesktopUsersResponse{Departments: deptDTOs, Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

// CreateUser creates one user.
func (service *DesktopDataService) CreateUser(ctx context.Context, input DesktopUserListItem) (DesktopUserListItem, error) {
	if input.UID == "" || input.Name == "" {
		return DesktopUserListItem{}, fmt.Errorf("uid and name are required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.CreateUser(ctx, DesktopUserRecord{UID: input.UID, Name: input.Name, Department: input.Dept, Phone: input.Phone, Role: input.Role, Status: input.Status, UpdatedAt: time.Now(), UpdatedBy: "桌面端"})
	if err != nil {
		return DesktopUserListItem{}, err
	}
	return DesktopUserListItem{ID: record.ID, UID: record.UID, Name: record.Name, Dept: record.Department, Phone: record.Phone, Role: record.Role, Status: record.Status, Date: record.UpdatedAt.Format(time.DateTime) + "/" + record.UpdatedBy}, nil
}

// UpdateUser updates one user.
func (service *DesktopDataService) UpdateUser(ctx context.Context, input DesktopUserListItem) (DesktopUserListItem, error) {
	if input.ID <= 0 || input.UID == "" || input.Name == "" {
		return DesktopUserListItem{}, fmt.Errorf("id, uid and name are required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.UpdateUser(ctx, DesktopUserRecord{ID: input.ID, UID: input.UID, Name: input.Name, Department: input.Dept, Phone: input.Phone, Role: input.Role, Status: input.Status, UpdatedAt: time.Now(), UpdatedBy: "桌面端"})
	if err != nil {
		return DesktopUserListItem{}, err
	}
	return DesktopUserListItem{ID: record.ID, UID: record.UID, Name: record.Name, Dept: record.Department, Phone: record.Phone, Role: record.Role, Status: record.Status, Date: record.UpdatedAt.Format(time.DateTime) + "/" + record.UpdatedBy}, nil
}

// DeleteUser removes one user.
func (service *DesktopDataService) DeleteUser(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return nil
	}
	return service.repository.DeleteUser(ctx, id)
}

// ResetUserPassword performs a simulated password reset audit update.
func (service *DesktopDataService) ResetUserPassword(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return nil
	}
	return service.repository.ResetUserPassword(ctx, id, "桌面端重置密码", time.Now())
}

// ListRoles returns role data.
func (service *DesktopDataService) ListRoles(ctx context.Context, name string, key string, status string, page int, pageSize int) (DesktopRolesResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	if service == nil || service.repository == nil {
		roles := filterRoles(fallbackRoles(), name, key, status)
		return DesktopRolesResponse{Items: paginateRoles(roles, page, pageSize), Total: len(roles), Page: page, PageSize: pageSize}, nil
	}
	items, total, err := service.repository.ListRoles(ctx, name, key, status, page, pageSize)
	if err != nil {
		return DesktopRolesResponse{}, err
	}
	result := make([]DesktopRole, 0, len(items))
	for _, item := range items {
		result = append(result, DesktopRole{ID: item.ID, Name: item.Name, Key: item.Key, Order: item.Order, Status: item.Status, Date: item.CreatedAt.Format(time.DateTime)})
	}
	if len(result) == 0 {
		fallback := filterRoles(fallbackRoles(), name, key, status)
		return DesktopRolesResponse{Items: paginateRoles(fallback, page, pageSize), Total: len(fallback), Page: page, PageSize: pageSize}, nil
	}
	return DesktopRolesResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

// CreateRole creates one role.
func (service *DesktopDataService) CreateRole(ctx context.Context, input DesktopRole) (DesktopRole, error) {
	if input.Name == "" || input.Key == "" {
		return DesktopRole{}, fmt.Errorf("name and key are required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.CreateRole(ctx, DesktopRoleRecord{Name: input.Name, Key: input.Key, Order: input.Order, Status: input.Status, CreatedAt: time.Now()})
	if err != nil {
		return DesktopRole{}, err
	}
	return DesktopRole{ID: record.ID, Name: record.Name, Key: record.Key, Order: record.Order, Status: record.Status, Date: record.CreatedAt.Format(time.DateTime)}, nil
}

// UpdateRole updates one role.
func (service *DesktopDataService) UpdateRole(ctx context.Context, input DesktopRole) (DesktopRole, error) {
	if input.ID <= 0 || input.Name == "" || input.Key == "" {
		return DesktopRole{}, fmt.Errorf("id, name and key are required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.UpdateRole(ctx, DesktopRoleRecord{ID: input.ID, Name: input.Name, Key: input.Key, Order: input.Order, Status: input.Status, CreatedAt: time.Now()})
	if err != nil {
		return DesktopRole{}, err
	}
	return DesktopRole{ID: record.ID, Name: record.Name, Key: record.Key, Order: record.Order, Status: record.Status, Date: record.CreatedAt.Format(time.DateTime)}, nil
}

// DeleteRole removes one role.
func (service *DesktopDataService) DeleteRole(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return nil
	}
	return service.repository.DeleteRole(ctx, id)
}

// GetRoleDataPermission returns one role's data permission.
func (service *DesktopDataService) GetRoleDataPermission(ctx context.Context, roleID int64) (DesktopRoleDataPermission, error) {
	if roleID <= 0 {
		return DesktopRoleDataPermission{}, fmt.Errorf("roleId is required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return DesktopRoleDataPermission{RoleID: roleID, Scope: "all", Departments: []string{}}, nil
	}
	record, err := service.repository.GetRoleDataPermission(ctx, roleID)
	if err != nil {
		return DesktopRoleDataPermission{}, err
	}
	departments := []string{}
	if record.Departments != "" {
		_ = json.Unmarshal([]byte(record.Departments), &departments)
	}
	if record.Scope == "" {
		record.Scope = "all"
	}
	return DesktopRoleDataPermission{RoleID: roleID, Scope: record.Scope, Departments: departments}, nil
}

// SaveRoleDataPermission saves one role's data permission.
func (service *DesktopDataService) SaveRoleDataPermission(ctx context.Context, input DesktopRoleDataPermission) (DesktopRoleDataPermission, error) {
	if input.RoleID <= 0 {
		return DesktopRoleDataPermission{}, fmt.Errorf("roleId is required: %w", ErrInvalidInput)
	}
	if input.Scope == "" {
		input.Scope = "all"
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	payload, _ := json.Marshal(input.Departments)
	record, err := service.repository.SaveRoleDataPermission(ctx, DesktopRoleDataPermissionRecord{RoleID: input.RoleID, Scope: input.Scope, Departments: string(payload)})
	if err != nil {
		return DesktopRoleDataPermission{}, err
	}
	departments := []string{}
	if record.Departments != "" {
		_ = json.Unmarshal([]byte(record.Departments), &departments)
	}
	return DesktopRoleDataPermission{RoleID: record.RoleID, Scope: record.Scope, Departments: departments}, nil
}

// ListOrgs returns organization data.
func (service *DesktopDataService) ListOrgs(ctx context.Context, department string, name string, code string, page int, pageSize int) (DesktopOrgsResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	if service == nil || service.repository == nil {
		items := filterOrgs(fallbackOrgs(), department, name, code)
		return DesktopOrgsResponse{Departments: fallbackDepartments(), Items: paginateOrgs(items, page, pageSize), Total: len(items), Page: page, PageSize: pageSize}, nil
	}
	departments, err := service.repository.ListOrgDepartments(ctx)
	if err != nil {
		return DesktopOrgsResponse{}, err
	}
	items, total, err := service.repository.ListOrgs(ctx, department, name, code, page, pageSize)
	if err != nil {
		return DesktopOrgsResponse{}, err
	}
	result := make([]DesktopOrg, 0, len(items))
	for _, item := range items {
		result = append(result, DesktopOrg{ID: item.ID, Name: item.Name, Code: item.Code, Parent1: item.ParentName, Order: item.Sort, Parent2: item.Category, Created: item.CreatedAt.Format(time.DateTime) + "/" + item.CreatedBy, Modified: item.UpdatedAt.Format(time.DateTime) + "/" + item.UpdatedBy})
	}
	deptDTOs := make([]DesktopDepartment, 0, len(departments))
	for _, item := range departments {
		deptDTOs = append(deptDTOs, DesktopDepartment{Name: item.Name, Count: item.Count})
	}
	if len(result) == 0 {
		fallback := filterOrgs(fallbackOrgs(), department, name, code)
		return DesktopOrgsResponse{Departments: fallbackDepartments(), Items: paginateOrgs(fallback, page, pageSize), Total: len(fallback), Page: page, PageSize: pageSize}, nil
	}
	return DesktopOrgsResponse{Departments: deptDTOs, Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

// CreateOrg creates one organization.
func (service *DesktopDataService) CreateOrg(ctx context.Context, input DesktopOrg) (DesktopOrg, error) {
	if input.Name == "" || input.Code == "" {
		return DesktopOrg{}, fmt.Errorf("name and code are required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.CreateOrg(ctx, DesktopOrgRecord{Name: input.Name, Code: input.Code, ParentName: input.Parent1, Sort: input.Order, Category: input.Parent2, CreatedAt: time.Now(), CreatedBy: "桌面端", UpdatedAt: time.Now(), UpdatedBy: "桌面端"})
	if err != nil {
		return DesktopOrg{}, err
	}
	return DesktopOrg{ID: record.ID, Name: record.Name, Code: record.Code, Parent1: record.ParentName, Order: record.Sort, Parent2: record.Category, Created: record.CreatedAt.Format(time.DateTime) + "/" + record.CreatedBy, Modified: record.UpdatedAt.Format(time.DateTime) + "/" + record.UpdatedBy}, nil
}

// UpdateOrg updates one organization.
func (service *DesktopDataService) UpdateOrg(ctx context.Context, input DesktopOrg) (DesktopOrg, error) {
	if input.ID <= 0 || input.Name == "" || input.Code == "" {
		return DesktopOrg{}, fmt.Errorf("id, name and code are required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.UpdateOrg(ctx, DesktopOrgRecord{ID: input.ID, Name: input.Name, Code: input.Code, ParentName: input.Parent1, Sort: input.Order, Category: input.Parent2, UpdatedAt: time.Now(), UpdatedBy: "桌面端"})
	if err != nil {
		return DesktopOrg{}, err
	}
	return DesktopOrg{ID: record.ID, Name: record.Name, Code: record.Code, Parent1: record.ParentName, Order: record.Sort, Parent2: record.Category, Created: record.CreatedAt.Format(time.DateTime) + "/" + record.CreatedBy, Modified: record.UpdatedAt.Format(time.DateTime) + "/" + record.UpdatedBy}, nil
}

// DeleteOrg removes one organization.
func (service *DesktopDataService) DeleteOrg(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return nil
	}
	return service.repository.DeleteOrg(ctx, id)
}

// ListTenants returns tenant data.
func (service *DesktopDataService) ListTenants(ctx context.Context, name string, code string, status string, page int, pageSize int) (DesktopTenantsResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	if service == nil || service.repository == nil {
		items := filterTenants(fallbackTenants(), name, code, status)
		return DesktopTenantsResponse{Items: paginateTenants(items, page, pageSize), Total: len(items), Page: page, PageSize: pageSize}, nil
	}
	items, total, err := service.repository.ListTenants(ctx, name, code, status, page, pageSize)
	if err != nil {
		return DesktopTenantsResponse{}, err
	}
	result := make([]DesktopTenant, 0, len(items))
	for _, item := range items {
		result = append(result, DesktopTenant{ID: item.ID, Code: item.Code, Name: item.Name, Status: item.Status, Period: item.Period, Admin: item.Admin, Phone: item.Phone})
	}
	if len(result) == 0 {
		fallback := filterTenants(fallbackTenants(), name, code, status)
		return DesktopTenantsResponse{Items: paginateTenants(fallback, page, pageSize), Total: len(fallback), Page: page, PageSize: pageSize}, nil
	}
	return DesktopTenantsResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

// CreateTenant creates one tenant.
func (service *DesktopDataService) CreateTenant(ctx context.Context, input DesktopTenant) (DesktopTenant, error) {
	if input.Name == "" || input.Code == "" {
		return DesktopTenant{}, fmt.Errorf("name and code are required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.CreateTenant(ctx, DesktopTenantRecord{Code: input.Code, Name: input.Name, Status: input.Status, Period: input.Period, Admin: input.Admin, Phone: input.Phone})
	if err != nil {
		return DesktopTenant{}, err
	}
	return DesktopTenant{ID: record.ID, Code: record.Code, Name: record.Name, Status: record.Status, Period: record.Period, Admin: record.Admin, Phone: record.Phone}, nil
}

// UpdateTenant updates one tenant.
func (service *DesktopDataService) UpdateTenant(ctx context.Context, input DesktopTenant) (DesktopTenant, error) {
	if input.ID <= 0 || input.Name == "" || input.Code == "" {
		return DesktopTenant{}, fmt.Errorf("id, name and code are required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.UpdateTenant(ctx, DesktopTenantRecord{ID: input.ID, Code: input.Code, Name: input.Name, Status: input.Status, Period: input.Period, Admin: input.Admin, Phone: input.Phone})
	if err != nil {
		return DesktopTenant{}, err
	}
	return DesktopTenant{ID: record.ID, Code: record.Code, Name: record.Name, Status: record.Status, Period: record.Period, Admin: record.Admin, Phone: record.Phone}, nil
}

// DeleteTenant removes one tenant.
func (service *DesktopDataService) DeleteTenant(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return nil
	}
	return service.repository.DeleteTenant(ctx, id)
}

// ListMenuTemplates returns menu template data.
func (service *DesktopDataService) ListMenuTemplates(ctx context.Context, name string, description string, page int, pageSize int) (DesktopMenuTemplatesResponse, error) {
	page, pageSize = normalizePagination(page, pageSize)
	if service == nil || service.repository == nil {
		items := filterMenuTemplates(fallbackMenuTemplates(), name, description)
		return DesktopMenuTemplatesResponse{Items: paginateMenuTemplates(items, page, pageSize), Total: len(items), Page: page, PageSize: pageSize}, nil
	}
	items, total, err := service.repository.ListMenuTemplates(ctx, name, description, page, pageSize)
	if err != nil {
		return DesktopMenuTemplatesResponse{}, err
	}
	result := make([]DesktopMenuTemplate, 0, len(items))
	for _, item := range items {
		result = append(result, DesktopMenuTemplate{ID: item.ID, Name: item.Name, Tenants: item.TenantCount, Desc: item.Description, Modified: item.UpdatedAt.Format(time.DateTime) + "/" + item.UpdatedBy, Created: item.CreatedAt.Format(time.DateTime) + "/" + item.CreatedBy})
	}
	if len(result) == 0 {
		fallback := filterMenuTemplates(fallbackMenuTemplates(), name, description)
		return DesktopMenuTemplatesResponse{Items: paginateMenuTemplates(fallback, page, pageSize), Total: len(fallback), Page: page, PageSize: pageSize}, nil
	}
	return DesktopMenuTemplatesResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

// ListDesktopMenus returns desktop menu management rows.
func (service *DesktopDataService) ListDesktopMenus(ctx context.Context, name string, menuType string, status string) (DesktopManagedMenusResponse, error) {
	if service == nil || service.repository == nil {
		fallback := filterManagedMenus(fallbackManagedMenus(), name, menuType, status)
		return DesktopManagedMenusResponse{Items: fallback}, nil
	}
	items, err := service.repository.ListDesktopMenus(ctx, name, menuType, status)
	if err != nil {
		return DesktopManagedMenusResponse{}, err
	}
	result := mapManagedMenus(items)
	if len(result) == 0 {
		fallback := filterManagedMenus(fallbackManagedMenus(), name, menuType, status)
		return DesktopManagedMenusResponse{Items: fallback}, nil
	}
	return DesktopManagedMenusResponse{Items: result}, nil
}

// CreateDesktopMenu creates one desktop menu.
func (service *DesktopDataService) CreateDesktopMenu(ctx context.Context, input DesktopManagedMenu) (DesktopManagedMenu, error) {
	if input.Name == "" {
		return DesktopManagedMenu{}, fmt.Errorf("name is required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.CreateDesktopMenu(ctx, DesktopMenuRecord{ParentID: input.ParentID, Name: input.Name, Level: parseMenuLevel(input.Level), Sort: input.Order, Type: input.Type, Icon: input.Icon, Status: input.Status, Path: input.Path, Permission: input.Perm})
	if err != nil {
		return DesktopManagedMenu{}, err
	}
	return mapManagedMenus([]DesktopMenuRecord{record})[0], nil
}

// UpdateDesktopMenu updates one desktop menu.
func (service *DesktopDataService) UpdateDesktopMenu(ctx context.Context, input DesktopManagedMenu) (DesktopManagedMenu, error) {
	if input.ID <= 0 || input.Name == "" {
		return DesktopManagedMenu{}, fmt.Errorf("id and name are required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.UpdateDesktopMenu(ctx, DesktopMenuRecord{ID: input.ID, ParentID: input.ParentID, Name: input.Name, Level: parseMenuLevel(input.Level), Sort: input.Order, Type: input.Type, Icon: input.Icon, Status: input.Status, Path: input.Path, Permission: input.Perm})
	if err != nil {
		return DesktopManagedMenu{}, err
	}
	return mapManagedMenus([]DesktopMenuRecord{record})[0], nil
}

// DeleteDesktopMenu removes one desktop menu.
func (service *DesktopDataService) DeleteDesktopMenu(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return nil
	}
	return service.repository.DeleteDesktopMenu(ctx, id)
}

// CreateMenuTemplate creates one menu template.
func (service *DesktopDataService) CreateMenuTemplate(ctx context.Context, input DesktopMenuTemplate) (DesktopMenuTemplate, error) {
	if input.Name == "" {
		return DesktopMenuTemplate{}, fmt.Errorf("name is required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.CreateMenuTemplate(ctx, DesktopMenuTemplateRecord{Name: input.Name, TenantCount: input.Tenants, Description: input.Desc, UpdatedAt: time.Now(), UpdatedBy: "桌面端", CreatedAt: time.Now(), CreatedBy: "桌面端"})
	if err != nil {
		return DesktopMenuTemplate{}, err
	}
	return DesktopMenuTemplate{ID: record.ID, Name: record.Name, Tenants: record.TenantCount, Desc: record.Description, Modified: record.UpdatedAt.Format(time.DateTime) + "/" + record.UpdatedBy, Created: record.CreatedAt.Format(time.DateTime) + "/" + record.CreatedBy}, nil
}

// UpdateMenuTemplate updates one menu template.
func (service *DesktopDataService) UpdateMenuTemplate(ctx context.Context, input DesktopMenuTemplate) (DesktopMenuTemplate, error) {
	if input.ID <= 0 || input.Name == "" {
		return DesktopMenuTemplate{}, fmt.Errorf("id and name are required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return input, nil
	}
	record, err := service.repository.UpdateMenuTemplate(ctx, DesktopMenuTemplateRecord{ID: input.ID, Name: input.Name, TenantCount: input.Tenants, Description: input.Desc, UpdatedAt: time.Now(), UpdatedBy: "桌面端", CreatedAt: time.Now(), CreatedBy: "桌面端"})
	if err != nil {
		return DesktopMenuTemplate{}, err
	}
	return DesktopMenuTemplate{ID: record.ID, Name: record.Name, Tenants: record.TenantCount, Desc: record.Description, Modified: record.UpdatedAt.Format(time.DateTime) + "/" + record.UpdatedBy, Created: record.CreatedAt.Format(time.DateTime) + "/" + record.CreatedBy}, nil
}

// DeleteMenuTemplate removes one menu template.
func (service *DesktopDataService) DeleteMenuTemplate(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	if service == nil || service.repository == nil {
		return nil
	}
	return service.repository.DeleteMenuTemplate(ctx, id)
}

// MarkMessageRead marks one desktop message as read.
func (service *DesktopDataService) MarkMessageRead(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("message id must be positive: %w", ErrInvalidInput)
	}

	if service == nil {
		return ErrUnavailable
	}

	service.mu.Lock()
	updated := false
	for index, item := range service.messages {
		if item.ID == id {
			service.messages[index].Read = true
			updated = true
			break
		}
	}
	service.mu.Unlock()

	if service.repository == nil {
		if !updated {
			return ErrNotFound
		}
		return nil
	}

	if err := service.repository.MarkMessageRead(ctx, id); err != nil && !updated {
		return err
	}

	if !updated {
		return nil
	}

	return nil
}

// MarkMessagesRead marks a batch of desktop messages as read.
func (service *DesktopDataService) MarkMessagesRead(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return fmt.Errorf("message ids are required: %w", ErrInvalidInput)
	}
	if service == nil {
		return ErrUnavailable
	}

	service.mu.Lock()
	updated := false
	idSet := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		if id > 0 {
			idSet[id] = struct{}{}
		}
	}
	for index, item := range service.messages {
		if _, ok := idSet[item.ID]; ok {
			service.messages[index].Read = true
			updated = true
		}
	}
	service.mu.Unlock()

	if service.repository == nil {
		if !updated {
			return ErrNotFound
		}
		return nil
	}

	if err := service.repository.MarkMessagesRead(ctx, ids); err != nil && !updated {
		return err
	}
	return nil
}

// UpdateMonitorStatus updates one desktop monitor item status.
func (service *DesktopDataService) UpdateMonitorStatus(ctx context.Context, id int64, status string) error {
	if id <= 0 {
		return fmt.Errorf("monitor id must be positive: %w", ErrInvalidInput)
	}
	if status != "ignored" && status != "processing" {
		return fmt.Errorf("status must be ignored or processing: %w", ErrInvalidInput)
	}

	if service == nil {
		return ErrUnavailable
	}

	service.mu.Lock()
	updated := false
	for index, item := range service.monitor {
		if item.ID == id {
			service.monitor[index].Status = status
			if status == "ignored" {
				service.monitor[index].Alarm = false
			}
			updated = true
			break
		}
	}
	service.mu.Unlock()

	if service.repository == nil {
		if !updated {
			return ErrNotFound
		}
		return nil
	}

	if err := service.repository.UpdateMonitorStatus(ctx, id, status); err != nil && !updated {
		return err
	}

	if !updated {
		return nil
	}

	return nil
}

func fallbackMessages() []DesktopMessage {
	return []DesktopMessage{
		{ID: 1, Type: "通知消息", Title: "欢迎使用桌面管理台", Description: "当前为无数据库模式，已展示内置样例消息。", PublishedAt: time.Now().Add(-2 * time.Hour).Format(time.DateTime), Author: "系统管理员", Read: false},
		{ID: 2, Type: "系统公告", Title: "消息中心已接入", Description: "后续配置数据库后会自动切换到 PostgreSQL 数据源。", PublishedAt: time.Now().Add(-8 * time.Hour).Format(time.DateTime), Author: "研发团队", Read: true},
	}
}

func fallbackMonitorItems() []DesktopMonitorItem {
	return []DesktopMonitorItem{
		{ID: 1, Metric: "CPU使用率", Value: "67%", Alarm: true, Level: "一级告警", OccurredAt: time.Now().Add(-20 * time.Minute).Format(time.DateTime), Description: "当前为无数据库模式，展示内置监控样例。", Status: "pending"},
		{ID: 2, Metric: "GPU显存使用率", Value: "82%", Alarm: true, Level: "二级告警", OccurredAt: time.Now().Add(-50 * time.Minute).Format(time.DateTime), Description: "可配置 PostgreSQL 后自动切换为数据库数据。", Status: "processing"},
	}
}

func buildFallbackMonitorSummary(items []DesktopMonitorItem) DesktopMonitorSummary {
	lastCollectedAt := time.Now().Format(time.DateTime)
	if len(items) > 0 {
		lastCollectedAt = items[0].OccurredAt
	}
	activeAlarmCount := int64(0)
	highestLevel := "无告警"
	for _, item := range items {
		if item.Alarm {
			activeAlarmCount++
		}
		if highestLevel == "无告警" || item.Level == "一级告警" {
			highestLevel = item.Level
		}
	}

	return DesktopMonitorSummary{Total: int64(len(items)), ActiveAlarmCount: activeAlarmCount, HighestLevel: highestLevel, LastCollectedAt: lastCollectedAt}
}

func fallbackAnnouncements() []DesktopAnnouncement {
	return []DesktopAnnouncement{
		{ID: 1, Title: "桌面端版本升级通知", Type: "公告", Status: "已发布", Publish: time.Now().Add(-6*time.Hour).Format(time.DateTime) + "/张三", Create: time.Now().Add(-7*time.Hour).Format(time.DateTime) + "/张三"},
		{ID: 2, Title: "节假日值班安排通知", Type: "通知", Status: "待发布", Publish: time.Now().Add(-2*time.Hour).Format(time.DateTime) + "/李四", Create: time.Now().Add(-3*time.Hour).Format(time.DateTime) + "/李四"},
		{ID: 3, Title: "系统巡检草稿", Type: "通知", Status: "草稿", Publish: time.Now().Add(-1*time.Hour).Format(time.DateTime) + "/王五", Create: time.Now().Add(-90*time.Minute).Format(time.DateTime) + "/王五"},
	}
}

func fallbackLoginLogs() []DesktopLoginLog {
	return []DesktopLoginLog{{ID: 1, LogID: "LOG-10001", Category: "登录", UserID: "100001", Name: "张三", Status: "成功", Time: time.Now().Add(-2 * time.Hour).Format(time.DateTime), IP: "10.111.123.131", Address: "广东省深圳市福田区", Browser: "Chrome 11", Desc: "登录成功"}, {ID: 2, LogID: "LOG-10002", Category: "登录", UserID: "100002", Name: "李四", Status: "失败", Time: time.Now().Add(-90 * time.Minute).Format(time.DateTime), IP: "10.111.123.132", Address: "上海市浦东新区", Browser: "Chrome 11", Desc: "密码错误，登录失败"}}
}

func fallbackOperationLogs() []DesktopOperationLog {
	return []DesktopOperationLog{{ID: 1, LogID: "OP-10001", Module: "租户管理", Category: "修改", UserID: "100001", Name: "张三", Status: "成功", Time: time.Now().Add(-3 * time.Hour).Format(time.DateTime), IP: "10.111.123.131", Browser: "Chrome 11", Desc: "修改租户配额成功"}, {ID: 2, LogID: "OP-10002", Module: "菜单管理", Category: "删除", UserID: "100002", Name: "李四", Status: "失败", Time: time.Now().Add(-2 * time.Hour).Format(time.DateTime), IP: "10.111.123.132", Browser: "Chrome 11", Desc: "删除菜单失败，权限不足"}}
}

func fallbackDicts() []DesktopDict {
	return []DesktopDict{
		{ID: 1, DictID: "001", Name: "岗位", Desc: "岗位字典", Modified: time.Now().Add(-2*time.Hour).Format(time.DateTime) + "/张三", Created: time.Now().Add(-24*time.Hour).Format(time.DateTime) + "/张三", Items: []DesktopDictItem{{ID: 11, Index: 1, Label: "前端工程师", KeyVal: "FE", StyleType: "主要", Modified: time.Now().Add(-2*time.Hour).Format(time.DateTime) + "/张三", Created: time.Now().Add(-24*time.Hour).Format(time.DateTime) + "/张三", BadgeColor: "bg-orange-50 text-orange-500"}, {ID: 12, Index: 2, Label: "后端工程师", KeyVal: "BE", StyleType: "次要", Modified: time.Now().Add(-90*time.Minute).Format(time.DateTime) + "/李四", Created: time.Now().Add(-20*time.Hour).Format(time.DateTime) + "/李四", BadgeColor: "bg-emerald-50 text-emerald-500"}}},
		{ID: 2, DictID: "002", Name: "职称", Desc: "职称字典", Modified: time.Now().Add(-3*time.Hour).Format(time.DateTime) + "/张三", Created: time.Now().Add(-48*time.Hour).Format(time.DateTime) + "/张三", Items: []DesktopDictItem{{ID: 21, Index: 1, Label: "高级职称", KeyVal: "1", StyleType: "主要", Modified: time.Now().Add(-3*time.Hour).Format(time.DateTime) + "/张三", Created: time.Now().Add(-48*time.Hour).Format(time.DateTime) + "/张三", BadgeColor: "bg-orange-50 text-orange-500"}, {ID: 22, Index: 2, Label: "中级职称", KeyVal: "2", StyleType: "次要", Modified: time.Now().Add(-2*time.Hour).Format(time.DateTime) + "/李四", Created: time.Now().Add(-30*time.Hour).Format(time.DateTime) + "/李四", BadgeColor: "bg-emerald-50 text-emerald-500"}, {ID: 23, Index: 3, Label: "低级职称", KeyVal: "3", StyleType: "次要", Modified: time.Now().Add(-1*time.Hour).Format(time.DateTime) + "/王五", Created: time.Now().Add(-26*time.Hour).Format(time.DateTime) + "/王五", BadgeColor: "bg-slate-100 dark:bg-slate-800 text-slate-500 dark:text-slate-400"}}},
	}
}

func filterDicts(items []DesktopDict, name string, dictID string) []DesktopDict {
	result := make([]DesktopDict, 0, len(items))
	for _, item := range items {
		if name != "" && !containsText(item.Name, name) {
			continue
		}
		if dictID != "" && !containsText(item.DictID, dictID) {
			continue
		}
		result = append(result, item)
	}
	return result
}

func fallbackSystemSettings() DesktopSystemSettings {
	return DesktopSystemSettings{
		WebsiteTitle:          "多租户后台管理系统",
		SystemLogo:            "",
		Theme:                 "emerald",
		ICP:                   "京ICP备10000000号-1",
		Copyright:             "Copyright © 2023 MTBM SYSTEM. All Rights Reserved.",
		RequireStrongPassword: true,
		LoginFailLimit:        5,
		LoginLockMinutes:      30,
	}
}

// GetSystemStats aggregates current system statistics from existing monitor, tenant, and user data.
func (service *DesktopDataService) GetSystemStats(ctx context.Context, rangeKey string) (DesktopSystemStatsResponse, error) {
	if service == nil || service.repository == nil {
		return buildFallbackSystemStats(), nil
	}

	source, err := service.repository.GetSystemStatsSource(ctx)
	if err != nil {
		return DesktopSystemStatsResponse{}, err
	}

	if len(source.MonitorEvents) == 0 && len(source.Tenants) == 0 && len(source.Users) == 0 {
		return buildFallbackSystemStats(), nil
	}

	now := time.Now()
	todayStart := startOfDay(now)
	yesterdayStart := todayStart.AddDate(0, 0, -1)
	weekStart := startOfWeek(now)
	prevWeekStart := weekStart.AddDate(0, 0, -7)
	monthStart := startOfMonth(now)
	prevMonthStart := monthStart.AddDate(0, -1, 0)

	todayCount := countMonitorsBetween(source.MonitorEvents, todayStart, now)
	yesterdayCount := countMonitorsBetween(source.MonitorEvents, yesterdayStart, todayStart)
	weekCount := countMonitorsBetween(source.MonitorEvents, weekStart, now)
	prevWeekCount := countMonitorsBetween(source.MonitorEvents, prevWeekStart, weekStart)
	monthCount := countMonitorsBetween(source.MonitorEvents, monthStart, now)
	prevMonthCount := countMonitorsBetween(source.MonitorEvents, prevMonthStart, monthStart)

	tenantNormal, tenantExpiring, tenantExpired, openCount, closedCount := classifyTenants(source.Tenants, now)
	departmentStats := aggregateUserDepartments(source.Users)
	roleStats := aggregateUserRoles(source.Users)
	trend := buildMonitorTrend(source.MonitorEvents, rangeKey, now)

	return DesktopSystemStatsResponse{
		AlertCards: []DesktopSystemStatsAlertCard{
			{Label: "本日告警总数", Value: todayCount, Change: formatPercentChange(todayCount, yesterdayCount), Period: "较昨日"},
			{Label: "本周告警数量", Value: weekCount, Change: formatPercentChange(weekCount, prevWeekCount), Period: "较上周"},
			{Label: "本月告警数量", Value: monthCount, Change: formatPercentChange(monthCount, prevMonthCount), Period: "较上月"},
		},
		Trend:             trend,
		TenantTotal:       len(source.Tenants),
		TenantOpenCount:   openCount,
		TenantClosedCount: closedCount,
		TenantNormal:      tenantNormal,
		TenantExpiring:    tenantExpiring,
		TenantExpired:     tenantExpired,
		EmployeeTotal:     len(source.Users),
		RoleStats:         roleStats,
		DepartmentStats:   departmentStats,
	}, nil
}

func buildFallbackSystemStats() DesktopSystemStatsResponse {
	return DesktopSystemStatsResponse{
		AlertCards:        []DesktopSystemStatsAlertCard{{Label: "本日告警总数", Value: 2, Change: "0.00%", Period: "较昨日"}, {Label: "本周告警数量", Value: 4, Change: "0.00%", Period: "较上周"}, {Label: "本月告警数量", Value: 4, Change: "0.00%", Period: "较上月"}},
		Trend:             []DesktopSystemStatsTrendPoint{{Label: "03-20", Value: 1}, {Label: "03-21", Value: 0}, {Label: "03-22", Value: 2}, {Label: "03-23", Value: 1}, {Label: "03-24", Value: 3}, {Label: "03-25", Value: 2}},
		TenantTotal:       4,
		TenantOpenCount:   1,
		TenantClosedCount: 3,
		TenantNormal:      DesktopSystemStatsTenantBucket{Count: 1, Names: []string{"新租户测试公司"}},
		TenantExpiring:    DesktopSystemStatsTenantBucket{Count: 0, Names: []string{}},
		TenantExpired:     DesktopSystemStatsTenantBucket{Count: 3, Names: []string{"智迪互动(北京)广告有限公司", "北京长友物业管理有限公司", "陕西沙龙传媒有限公司"}},
		EmployeeTotal:     5,
		RoleStats:         []DesktopSystemStatsRoleStat{{Label: "超管", Count: 1}, {Label: "管理员", Count: 2}, {Label: "财务", Count: 1}, {Label: "运营", Count: 1}},
		DepartmentStats:   []DesktopDepartment{{Name: "生产部", Count: 3}, {Name: "财务部", Count: 1}, {Name: "营销部", Count: 1}},
	}
}

// GetUsageStats aggregates usage statistics from existing login logs and operation logs.
func (service *DesktopDataService) GetUsageStats(ctx context.Context, rangeKey string, tenantCode string) (DesktopUsageStatsResponse, error) {
	if service == nil || service.repository == nil {
		return buildFallbackUsageStats(), nil
	}

	source, err := service.repository.GetUsageStatsSource(ctx)
	if err != nil {
		return DesktopUsageStatsResponse{}, err
	}

	if len(source.LoginLogs) == 0 && len(source.OperationLogs) == 0 {
		return buildFallbackUsageStats(), nil
	}

	filteredLoginLogs := filterUsageLoginLogsByTenant(source.LoginLogs, tenantCode)
	filteredOperationLogs := filterUsageOperationLogsByTenant(source.OperationLogs, tenantCode)
	if tenantCode != "" && tenantCode != "all" && len(filteredLoginLogs) == 0 && len(filteredOperationLogs) == 0 {
		return DesktopUsageStatsResponse{
			Tenants:  mapUsageTenants(source.Tenants),
			Trend:    []DesktopUsageStatsTrendPoint{},
			TopMost:  []DesktopUsageStatsModuleRank{},
			TopLeast: []DesktopUsageStatsModuleRank{},
		}, nil
	}

	now := time.Now()
	return DesktopUsageStatsResponse{
		Tenants:  mapUsageTenants(source.Tenants),
		Trend:    buildUsageTrend(filteredLoginLogs, filteredOperationLogs, rangeKey, now),
		TopMost:  buildModuleRanks(filteredOperationLogs, true),
		TopLeast: buildModuleRanks(filteredOperationLogs, false),
	}, nil
}

func buildFallbackUsageStats() DesktopUsageStatsResponse {
	return DesktopUsageStatsResponse{
		Tenants: []DesktopUsageStatsTenantOption{{Code: "all", Name: "全部租户"}},
		Trend: []DesktopUsageStatsTrendPoint{
			{Label: "03-19", LoginCount: 1, UsageCount: 1},
			{Label: "03-20", LoginCount: 2, UsageCount: 1},
			{Label: "03-21", LoginCount: 1, UsageCount: 2},
			{Label: "03-22", LoginCount: 3, UsageCount: 2},
			{Label: "03-23", LoginCount: 1, UsageCount: 1},
			{Label: "03-24", LoginCount: 2, UsageCount: 3},
			{Label: "03-25", LoginCount: 2, UsageCount: 2},
		},
		TopMost:  []DesktopUsageStatsModuleRank{{Rank: 1, Name: "用户管理", Count: 18, Change: "0.00%"}, {Rank: 2, Name: "租户管理", Count: 12, Change: "0.00%"}, {Rank: 3, Name: "菜单管理", Count: 10, Change: "0.00%"}},
		TopLeast: []DesktopUsageStatsModuleRank{{Rank: 1, Name: "系统监控", Count: 2, Change: "0.00%"}, {Rank: 2, Name: "公告管理", Count: 3, Change: "0.00%"}, {Rank: 3, Name: "字典管理", Count: 4, Change: "0.00%"}},
	}
}

func startOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

func startOfWeek(value time.Time) time.Time {
	weekdayOffset := int(value.Weekday())
	if weekdayOffset == 0 {
		weekdayOffset = 7
	}
	return startOfDay(value).AddDate(0, 0, -(weekdayOffset - 1))
}

func startOfMonth(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, value.Location())
}

func countMonitorsBetween(items []DesktopMonitorRecord, start time.Time, end time.Time) int64 {
	count := int64(0)
	for _, item := range items {
		if (item.OccurredAt.Equal(start) || item.OccurredAt.After(start)) && item.OccurredAt.Before(end) {
			count++
		}
	}
	return count
}

func formatPercentChange(current int64, previous int64) string {
	if previous <= 0 {
		if current <= 0 {
			return "0.00%"
		}
		return "100.00%"
	}
	return fmt.Sprintf("%.2f%%", (float64(current-previous)/float64(previous))*100)
}

func classifyTenants(items []DesktopTenantRecord, now time.Time) (DesktopSystemStatsTenantBucket, DesktopSystemStatsTenantBucket, DesktopSystemStatsTenantBucket, int, int) {
	const expiringWindowDays = 30
	normal := DesktopSystemStatsTenantBucket{Names: []string{}}
	expiring := DesktopSystemStatsTenantBucket{Names: []string{}}
	expired := DesktopSystemStatsTenantBucket{Names: []string{}}
	openCount := 0
	closedCount := 0
	for _, item := range items {
		if item.Status == "开启" {
			openCount++
		} else {
			closedCount++
		}
		endDate := parseTenantEndDate(item.Period)
		switch {
		case endDate.IsZero() || endDate.Before(now):
			expired.Count++
			expired.Names = append(expired.Names, item.Name)
		case endDate.Before(now.AddDate(0, 0, expiringWindowDays)):
			expiring.Count++
			expiring.Names = append(expiring.Names, item.Name)
		default:
			normal.Count++
			normal.Names = append(normal.Names, item.Name)
		}
	}
	return normal, expiring, expired, openCount, closedCount
}

func parseTenantEndDate(period string) time.Time {
	parts := strings.Split(period, " - ")
	if len(parts) != 2 {
		return time.Time{}
	}
	result, err := time.Parse("2006-01-02", strings.TrimSpace(parts[1]))
	if err != nil {
		return time.Time{}
	}
	return result
}

func aggregateUserDepartments(items []DesktopUserRecord) []DesktopDepartment {
	counts := map[string]int{}
	for _, item := range items {
		counts[item.Department]++
	}
	result := make([]DesktopDepartment, 0, len(counts))
	for name, count := range counts {
		result = append(result, DesktopDepartment{Name: name, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Name < result[j].Name
		}
		return result[i].Count > result[j].Count
	})
	if len(result) > 5 {
		result = result[:5]
	}
	return result
}

func aggregateUserRoles(items []DesktopUserRecord) []DesktopSystemStatsRoleStat {
	counts := map[string]int{}
	for _, item := range items {
		counts[item.Role]++
	}
	result := make([]DesktopSystemStatsRoleStat, 0, len(counts))
	for label, count := range counts {
		result = append(result, DesktopSystemStatsRoleStat{Label: label, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Label < result[j].Label
		}
		return result[i].Count > result[j].Count
	})
	return result
}

func buildMonitorTrend(items []DesktopMonitorRecord, rangeKey string, now time.Time) []DesktopSystemStatsTrendPoint {
	days := 7
	switch rangeKey {
	case "today":
		days = 1
	case "month":
		days = 30
	}
	counts := map[string]int64{}
	for _, item := range items {
		key := item.OccurredAt.Format("01-02")
		counts[key]++
	}
	result := make([]DesktopSystemStatsTrendPoint, 0, days)
	for index := days - 1; index >= 0; index-- {
		day := now.AddDate(0, 0, -index)
		key := day.Format("01-02")
		result = append(result, DesktopSystemStatsTrendPoint{Label: key, Value: counts[key]})
	}
	return result
}

func mapUsageTenants(items []DesktopTenantRecord) []DesktopUsageStatsTenantOption {
	result := make([]DesktopUsageStatsTenantOption, 0, len(items)+1)
	result = append(result, DesktopUsageStatsTenantOption{Code: "all", Name: "全部租户"})
	for _, item := range items {
		result = append(result, DesktopUsageStatsTenantOption{Code: item.Code, Name: item.Name})
	}
	return result
}

func buildUsageTrend(loginLogs []DesktopLoginLogRecord, operationLogs []DesktopOperationLogRecord, rangeKey string, now time.Time) []DesktopUsageStatsTrendPoint {
	days := 7
	switch rangeKey {
	case "today":
		days = 1
	case "month":
		days = 30
	}

	loginCounts := map[string]int{}
	for _, item := range loginLogs {
		loginCounts[item.Time.Format("01-02")]++
	}
	usageCounts := map[string]int{}
	for _, item := range operationLogs {
		usageCounts[item.Time.Format("01-02")]++
	}

	result := make([]DesktopUsageStatsTrendPoint, 0, days)
	for index := days - 1; index >= 0; index-- {
		day := now.AddDate(0, 0, -index)
		key := day.Format("01-02")
		result = append(result, DesktopUsageStatsTrendPoint{Label: key, LoginCount: loginCounts[key], UsageCount: usageCounts[key]})
	}
	return result
}

func buildModuleRanks(items []DesktopOperationLogRecord, descending bool) []DesktopUsageStatsModuleRank {
	type moduleCounter struct {
		Name  string
		Count int
	}

	counts := map[string]int{}
	for _, item := range items {
		moduleName := strings.TrimSpace(item.Module)
		if moduleName == "" {
			continue
		}
		counts[moduleName]++
	}

	modules := make([]moduleCounter, 0, len(counts))
	for name, count := range counts {
		modules = append(modules, moduleCounter{Name: name, Count: count})
	}

	sort.Slice(modules, func(i, j int) bool {
		if modules[i].Count == modules[j].Count {
			return modules[i].Name < modules[j].Name
		}
		if descending {
			return modules[i].Count > modules[j].Count
		}
		return modules[i].Count < modules[j].Count
	})

	if len(modules) > 5 {
		modules = modules[:5]
	}

	result := make([]DesktopUsageStatsModuleRank, 0, len(modules))
	for index, item := range modules {
		result = append(result, DesktopUsageStatsModuleRank{Rank: index + 1, Name: item.Name, Count: item.Count, Change: "0.00%"})
	}
	return result
}

func filterUsageLoginLogsByTenant(items []DesktopLoginLogRecord, tenantCode string) []DesktopLoginLogRecord {
	if tenantCode == "" || tenantCode == "all" {
		return items
	}
	result := make([]DesktopLoginLogRecord, 0, len(items))
	for _, item := range items {
		if item.TenantCode == tenantCode {
			result = append(result, item)
		}
	}
	return result
}

func filterUsageOperationLogsByTenant(items []DesktopOperationLogRecord, tenantCode string) []DesktopOperationLogRecord {
	if tenantCode == "" || tenantCode == "all" {
		return items
	}
	result := make([]DesktopOperationLogRecord, 0, len(items))
	for _, item := range items {
		if item.TenantCode == tenantCode {
			result = append(result, item)
		}
	}
	return result
}

func fallbackDepartments() []DesktopDepartment {
	return []DesktopDepartment{{Name: "人力资源部", Count: 12}, {Name: "财务部", Count: 23}, {Name: "生产部", Count: 32}, {Name: "营销部", Count: 23}, {Name: "安全部", Count: 23}, {Name: "党群部", Count: 23}, {Name: "保卫部", Count: 23}, {Name: "后勤部", Count: 23}, {Name: "行政部", Count: 23}}
}

func fallbackUsers() []DesktopUserListItem {
	return []DesktopUserListItem{{ID: 1, UID: "1001", Name: "冯政", Dept: "生产部", Phone: "13386911277", Role: "超管", Status: "开启", Date: "2023-05-26 12:12:00/张三"}, {ID: 2, UID: "1002", Name: "孙子面", Dept: "生产部", Phone: "15366978328", Role: "管理员", Status: "关闭", Date: "2023-05-26 12:12:00/张三"}, {ID: 3, UID: "1003", Name: "钱继初", Dept: "财务部", Phone: "18555418491", Role: "财务", Status: "关闭", Date: "2023-05-26 12:12:00/张三"}, {ID: 4, UID: "1004", Name: "李建华", Dept: "营销部", Phone: "18716223293", Role: "运营", Status: "开启", Date: "2023-05-26 12:12:00/张三"}}
}

func fallbackRoles() []DesktopRole {
	return []DesktopRole{{ID: 1, Name: "超级管理员", Key: "admin", Order: 1, Status: true, Date: "2023-05-26 12:12:00"}, {ID: 2, Name: "高管", Key: "manager", Order: 2, Status: true, Date: "2023-05-26 12:12:00"}, {ID: 3, Name: "部门负责人", Key: "dept_leader", Order: 3, Status: true, Date: "2023-05-26 12:12:00"}, {ID: 4, Name: "普通员工", Key: "common", Order: 4, Status: true, Date: "2023-05-26 12:12:00"}}
}

func fallbackOrgs() []DesktopOrg {
	return []DesktopOrg{{ID: 1, Name: "生产一部", Code: "32", Parent1: "东方科技", Order: 1, Parent2: "部门", Created: "2023-05-26 12:12:00/张三", Modified: "2023-05-26 12:12:00/张三"}, {ID: 2, Name: "生产二部", Code: "33", Parent1: "东方科技", Order: 2, Parent2: "部门", Created: "2023-05-26 12:12:00/张三", Modified: "2023-05-26 12:12:00/张三"}}
}

func fallbackTenants() []DesktopTenant {
	return []DesktopTenant{{ID: 1, Code: "FQJT", Name: "智迪互动(北京)广告有限公司", Status: "开启", Period: "2023-01-01 - 2023-12-31", Admin: "张三", Phone: "131 1242 1320"}, {ID: 2, Code: "HIMGJT", Name: "北京长友物业管理有限公司", Status: "关闭", Period: "2023-01-01 - 2023-12-31", Admin: "张三", Phone: "131 1242 1320"}}
}

func fallbackMenuTemplates() []DesktopMenuTemplate {
	return []DesktopMenuTemplate{{ID: 1, Name: "试用版", Tenants: 32, Desc: "试用版", Modified: "2023-05-26 12:12:00/张三", Created: "2023-05-26 12:12:00/张三"}, {ID: 2, Name: "基础版", Tenants: 32, Desc: "基础版", Modified: "2023-05-26 12:12:00/张三", Created: "2023-05-26 12:12:00/张三"}, {ID: 3, Name: "VIP版", Tenants: 32, Desc: "VIP版", Modified: "2023-05-26 12:12:00/张三", Created: "2023-05-26 12:12:00/张三"}}
}

func filterUsers(items []DesktopUserListItem, department string, uid string, name string, role string) []DesktopUserListItem {
	result := make([]DesktopUserListItem, 0, len(items))
	for _, item := range items {
		if department != "" && item.Dept != department {
			continue
		}
		if uid != "" && !containsText(item.UID, uid) {
			continue
		}
		if name != "" && !containsText(item.Name, name) {
			continue
		}
		if role != "" && item.Role != role {
			continue
		}
		result = append(result, item)
	}
	return result
}

func filterRoles(items []DesktopRole, name string, key string, status string) []DesktopRole {
	result := make([]DesktopRole, 0, len(items))
	for _, item := range items {
		if name != "" && !containsText(item.Name, name) {
			continue
		}
		if key != "" && !containsText(item.Key, key) {
			continue
		}
		if status != "" {
			statusValue := status == "正常" || status == "active"
			if item.Status != statusValue {
				continue
			}
		}
		result = append(result, item)
	}
	return result
}

func paginateUsers(items []DesktopUserListItem, page int, pageSize int) []DesktopUserListItem {
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []DesktopUserListItem{}
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func paginateRoles(items []DesktopRole, page int, pageSize int) []DesktopRole {
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []DesktopRole{}
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func filterOrgs(items []DesktopOrg, department string, name string, code string) []DesktopOrg {
	result := make([]DesktopOrg, 0, len(items))
	for _, item := range items {
		if department != "" && item.Parent1 != "东方科技" && false {
		}
		if name != "" && !containsText(item.Name, name) {
			continue
		}
		if code != "" && !containsText(item.Code, code) {
			continue
		}
		result = append(result, item)
	}
	return result
}

func filterTenants(items []DesktopTenant, name string, code string, status string) []DesktopTenant {
	result := make([]DesktopTenant, 0, len(items))
	for _, item := range items {
		if name != "" && !containsText(item.Name, name) {
			continue
		}
		if code != "" && !containsText(item.Code, code) {
			continue
		}
		if status != "" && item.Status != status {
			continue
		}
		result = append(result, item)
	}
	return result
}

func filterMenuTemplates(items []DesktopMenuTemplate, name string, description string) []DesktopMenuTemplate {
	result := make([]DesktopMenuTemplate, 0, len(items))
	for _, item := range items {
		if name != "" && !containsText(item.Name, name) {
			continue
		}
		if description != "" && !containsText(item.Desc, description) {
			continue
		}
		result = append(result, item)
	}
	return result
}

func paginateOrgs(items []DesktopOrg, page int, pageSize int) []DesktopOrg {
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []DesktopOrg{}
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func paginateTenants(items []DesktopTenant, page int, pageSize int) []DesktopTenant {
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []DesktopTenant{}
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func paginateMenuTemplates(items []DesktopMenuTemplate, page int, pageSize int) []DesktopMenuTemplate {
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []DesktopMenuTemplate{}
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func fallbackManagedMenus() []DesktopManagedMenu {
	return []DesktopManagedMenu{{ID: 1, ParentID: 0, Name: "首页", HasSub: false, Depth: 0, Level: "1级", Order: 1, Type: "目录", Status: "可见", Path: "/", Perm: "desktop:home:view", Icon: "Home"}, {ID: 2, ParentID: 0, Name: "租户配置", HasSub: true, Depth: 0, Level: "1级", Order: 2, Type: "目录", Status: "不可见", Path: "/tenant", Perm: "desktop:tenant:view", Icon: "Settings"}, {ID: 3, ParentID: 2, Name: "租户管理", HasSub: false, Depth: 1, Level: "2级", Order: 1, Type: "菜单", Status: "不可见", Path: "/tenant/management", Perm: "desktop:tenant:list", Icon: "Settings"}, {ID: 4, ParentID: 2, Name: "菜单模板", HasSub: false, Depth: 1, Level: "2级", Order: 2, Type: "菜单", Status: "不可见", Path: "/tenant/menu-template", Perm: "desktop:tenant:template", Icon: "Menu"}, {ID: 5, ParentID: 2, Name: "组织管理", HasSub: true, Depth: 1, Level: "2级", Order: 3, Type: "菜单", Status: "不可见", Path: "/orgs", Perm: "desktop:org:list", Icon: "Layers"}, {ID: 6, ParentID: 5, Name: "查询/查看", HasSub: false, Depth: 2, Level: "3级", Order: 1, Type: "按钮", Status: "不可见", Path: "/orgs", Perm: "desktop:org:view", Icon: "Layers"}}
}

func filterManagedMenus(items []DesktopManagedMenu, name string, menuType string, status string) []DesktopManagedMenu {
	result := make([]DesktopManagedMenu, 0, len(items))
	for _, item := range items {
		if name != "" && !containsText(item.Name, name) {
			continue
		}
		if menuType != "" && item.Type != menuType {
			continue
		}
		if status != "" && item.Status != status {
			continue
		}
		result = append(result, item)
	}
	return result
}

func parseMenuLevel(level string) int {
	if strings.Contains(level, "3") {
		return 3
	}
	if strings.Contains(level, "2") {
		return 2
	}
	return 1
}

func mapManagedMenus(records []DesktopMenuRecord) []DesktopManagedMenu {
	childCount := map[int64]int{}
	for _, item := range records {
		childCount[item.ParentID]++
	}
	result := make([]DesktopManagedMenu, 0, len(records))
	for _, item := range records {
		result = append(result, DesktopManagedMenu{ID: item.ID, ParentID: item.ParentID, Name: item.Name, Level: fmt.Sprintf("%d级", item.Level), Order: item.Sort, Type: item.Type, Icon: item.Icon, Status: item.Status, Path: item.Path, Perm: item.Permission, Depth: max(item.Level-1, 0), HasSub: childCount[item.ID] > 0})
	}
	return result
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func mapSystemSettings(record DesktopSystemSettingsRecord) DesktopSystemSettings {
	return DesktopSystemSettings{
		WebsiteTitle:          record.WebsiteTitle,
		SystemLogo:            record.SystemLogo,
		Theme:                 record.Theme,
		ICP:                   record.ICP,
		Copyright:             record.Copyright,
		RequireStrongPassword: record.RequireStrongPassword,
		LoginFailLimit:        record.LoginFailLimit,
		LoginLockMinutes:      record.LoginLockMinutes,
	}
}

func filterLoginLogs(items []DesktopLoginLog, name string, status string, startDate string, endDate string) []DesktopLoginLog {
	result := make([]DesktopLoginLog, 0, len(items))
	for _, item := range items {
		if name != "" && !containsText(item.Name, name) {
			continue
		}
		if status != "" && item.Status != status {
			continue
		}
		if !matchDateRange(item.Time, startDate, endDate) {
			continue
		}
		result = append(result, item)
	}
	return result
}

func filterOperationLogs(items []DesktopOperationLog, name string, status string, startDate string, endDate string) []DesktopOperationLog {
	result := make([]DesktopOperationLog, 0, len(items))
	for _, item := range items {
		if name != "" && !containsText(item.Name, name) {
			continue
		}
		if status != "" && item.Status != status {
			continue
		}
		if !matchDateRange(item.Time, startDate, endDate) {
			continue
		}
		result = append(result, item)
	}
	return result
}

func matchDateRange(value string, startDate string, endDate string) bool {
	if startDate == "" && endDate == "" {
		return true
	}
	if len(value) < 10 {
		return false
	}
	date := value[:10]
	if startDate != "" && date < startDate {
		return false
	}
	if endDate != "" && date > endDate {
		return false
	}
	return true
}

func paginateLoginLogs(items []DesktopLoginLog, page int, pageSize int) []DesktopLoginLog {
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []DesktopLoginLog{}
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func paginateOperationLogs(items []DesktopOperationLog, page int, pageSize int) []DesktopOperationLog {
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []DesktopOperationLog{}
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func filterMessages(items []DesktopMessage, messageType string, keyword string) []DesktopMessage {
	filtered := make([]DesktopMessage, 0, len(items))
	for _, item := range items {
		if messageType != "" && item.Type != messageType {
			continue
		}
		if keyword != "" && !containsText(item.Title, keyword) && !containsText(item.Description, keyword) {
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered
}

func filterMonitorItems(items []DesktopMonitorItem, metric string, level string) []DesktopMonitorItem {
	filtered := make([]DesktopMonitorItem, 0, len(items))
	for _, item := range items {
		if metric != "" && !containsText(item.Metric, metric) {
			continue
		}
		if level != "" && item.Level != level {
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered
}

func containsText(value string, query string) bool {
	if query == "" {
		return true
	}
	return strings.Contains(strings.ToLower(value), strings.ToLower(query))
}

func normalizePagination(page int, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func paginateMessages(items []DesktopMessage, page int, pageSize int) []DesktopMessage {
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []DesktopMessage{}
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func paginateMonitorItems(items []DesktopMonitorItem, page int, pageSize int) []DesktopMonitorItem {
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []DesktopMonitorItem{}
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}
