package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ai-go-service/internal/service"
	"github.com/jackc/pgx/v5"
)

// DesktopDataRepository provides PostgreSQL-backed sample data for desktop pages.
type DesktopDataRepository struct {
	pool *Pool
}

// NewDesktopDataRepository creates a repository when a PostgreSQL pool is available.
func NewDesktopDataRepository(pool *Pool) *DesktopDataRepository {
	if pool == nil || pool.pool == nil {
		return nil
	}

	return &DesktopDataRepository{pool: pool}
}

// EnsureSchemaAndSeed creates required tables and inserts sample rows when empty.
func (repository *DesktopDataRepository) EnsureSchemaAndSeed(ctx context.Context) error {
	if repository == nil || repository.pool == nil || repository.pool.pool == nil {
		return nil
	}

	statements := []string{
		`CREATE TABLE IF NOT EXISTS app_messages (
			id BIGSERIAL PRIMARY KEY,
			type TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			published_at TIMESTAMPTZ NOT NULL,
			author TEXT NOT NULL,
			is_read BOOLEAN NOT NULL DEFAULT FALSE
		)`,
		`CREATE TABLE IF NOT EXISTS app_monitor_events (
			id BIGSERIAL PRIMARY KEY,
			metric TEXT NOT NULL,
			metric_value TEXT NOT NULL,
			is_alarm BOOLEAN NOT NULL DEFAULT FALSE,
			level TEXT NOT NULL,
			occurred_at TIMESTAMPTZ NOT NULL,
			description TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending'
		)`,
		`CREATE TABLE IF NOT EXISTS app_monitor_collection_state (
			id SMALLINT PRIMARY KEY,
			last_collected_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS app_announcements (
			id BIGSERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			announcement_type TEXT NOT NULL,
			status TEXT NOT NULL,
			published_at TIMESTAMPTZ NOT NULL,
			published_by TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			created_by TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS app_login_logs (
			id BIGSERIAL PRIMARY KEY,
			log_id TEXT NOT NULL,
			category TEXT NOT NULL,
			user_id TEXT NOT NULL,
			name TEXT NOT NULL,
			status TEXT NOT NULL,
			occurred_at TIMESTAMPTZ NOT NULL,
			ip TEXT NOT NULL,
			address TEXT NOT NULL,
			browser TEXT NOT NULL,
			description TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS app_operation_logs (
			id BIGSERIAL PRIMARY KEY,
			log_id TEXT NOT NULL,
			module TEXT NOT NULL,
			category TEXT NOT NULL,
			user_id TEXT NOT NULL,
			name TEXT NOT NULL,
			status TEXT NOT NULL,
			occurred_at TIMESTAMPTZ NOT NULL,
			ip TEXT NOT NULL,
			browser TEXT NOT NULL,
			description TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS app_dicts (
			id BIGSERIAL PRIMARY KEY,
			dict_id TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			modified_at TIMESTAMPTZ NOT NULL,
			modified_by TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			created_by TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS app_dict_items (
			id BIGSERIAL PRIMARY KEY,
			dict_ref_id BIGINT NOT NULL REFERENCES app_dicts(id) ON DELETE CASCADE,
			item_index INT NOT NULL,
			label TEXT NOT NULL,
			key_val TEXT NOT NULL,
			style_type TEXT NOT NULL,
			modified_at TIMESTAMPTZ NOT NULL,
			modified_by TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			created_by TEXT NOT NULL,
			badge_color TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS app_system_settings (
			id INT PRIMARY KEY,
			website_title TEXT NOT NULL,
			system_logo TEXT NOT NULL,
			theme TEXT NOT NULL,
			icp TEXT NOT NULL,
			copyright TEXT NOT NULL,
			require_strong_password BOOLEAN NOT NULL,
			login_fail_limit INT NOT NULL,
			login_lock_minutes INT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS app_departments (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			user_count INT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS app_users (
			id BIGSERIAL PRIMARY KEY,
			uid TEXT NOT NULL,
			name TEXT NOT NULL,
			department TEXT NOT NULL,
			phone TEXT NOT NULL,
			role_name TEXT NOT NULL,
			status TEXT NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			updated_by TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS app_roles (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			role_key TEXT NOT NULL,
			display_order INT NOT NULL,
			status BOOLEAN NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS app_role_data_permissions (
			role_id BIGINT PRIMARY KEY REFERENCES app_roles(id) ON DELETE CASCADE,
			scope TEXT NOT NULL,
			departments_json TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS app_orgs (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			code TEXT NOT NULL,
			parent_name TEXT NOT NULL,
			sort_order INT NOT NULL,
			category TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			created_by TEXT NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			updated_by TEXT NOT NULL,
			department TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS app_tenants (
			id BIGSERIAL PRIMARY KEY,
			code TEXT NOT NULL,
			name TEXT NOT NULL,
			status TEXT NOT NULL,
			period TEXT NOT NULL,
			admin_name TEXT NOT NULL,
			admin_phone TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS app_menu_templates (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			tenant_count INT NOT NULL,
			description TEXT NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			updated_by TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			created_by TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS app_desktop_menus (
			id BIGSERIAL PRIMARY KEY,
			parent_id BIGINT NOT NULL,
			name TEXT NOT NULL,
			menu_level INT NOT NULL,
			sort_order INT NOT NULL,
			menu_type TEXT NOT NULL,
			icon TEXT NOT NULL,
			status TEXT NOT NULL,
			path TEXT NOT NULL,
			permission TEXT NOT NULL
		)`,
	}

	alterStatements := []string{
		`ALTER TABLE app_login_logs ADD COLUMN IF NOT EXISTS tenant_code TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE app_operation_logs ADD COLUMN IF NOT EXISTS tenant_code TEXT NOT NULL DEFAULT ''`,
	}

	for _, statement := range statements {
		if _, err := repository.pool.pool.Exec(ctx, statement); err != nil {
			return fmt.Errorf("ensure desktop sample tables: %w", err)
		}
	}
	for _, statement := range alterStatements {
		if _, err := repository.pool.pool.Exec(ctx, statement); err != nil {
			return fmt.Errorf("ensure desktop sample table columns: %w", err)
		}
	}

	if err := repository.seedMessages(ctx); err != nil {
		return err
	}
	if err := repository.seedMonitorEvents(ctx); err != nil {
		return err
	}
	if err := repository.seedAnnouncements(ctx); err != nil {
		return err
	}
	if err := repository.seedLoginLogs(ctx); err != nil {
		return err
	}
	if err := repository.seedOperationLogs(ctx); err != nil {
		return err
	}
	if err := repository.seedDicts(ctx); err != nil {
		return err
	}
	if err := repository.seedSystemSettings(ctx); err != nil {
		return err
	}
	if err := repository.seedUsersAndRoles(ctx); err != nil {
		return err
	}
	if err := repository.seedOrgsAndTenants(ctx); err != nil {
		return err
	}
	if err := repository.backfillUsageLogTenants(ctx); err != nil {
		return err
	}
	if err := repository.seedMenuTemplates(ctx); err != nil {
		return err
	}
	if err := repository.seedDesktopMenus(ctx); err != nil {
		return err
	}

	return nil
}

func (repository *DesktopDataRepository) ListDesktopMenus(ctx context.Context, name string, menuType string, status string) ([]service.DesktopMenuRecord, error) {
	rows, err := repository.pool.pool.Query(ctx, `SELECT id, parent_id, name, menu_level, sort_order, menu_type, icon, status, path, permission FROM app_desktop_menus WHERE ($1='' OR name ILIKE '%'||$1||'%') AND ($2='' OR menu_type=$2) AND ($3='' OR status=$3) ORDER BY menu_level ASC, sort_order ASC, id ASC`, strings.TrimSpace(name), strings.TrimSpace(menuType), strings.TrimSpace(status))
	if err != nil {
		return nil, fmt.Errorf("query desktop menus: %w", err)
	}
	defer rows.Close()
	items := make([]service.DesktopMenuRecord, 0)
	for rows.Next() {
		var item service.DesktopMenuRecord
		if err := rows.Scan(&item.ID, &item.ParentID, &item.Name, &item.Level, &item.Sort, &item.Type, &item.Icon, &item.Status, &item.Path, &item.Permission); err != nil {
			return nil, fmt.Errorf("scan desktop menus: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (repository *DesktopDataRepository) GetSystemStatsSource(ctx context.Context) (service.DesktopSystemStatsSource, error) {
	if repository == nil || repository.pool == nil || repository.pool.pool == nil {
		return service.DesktopSystemStatsSource{}, nil
	}
	monitorRows, err := repository.pool.pool.Query(ctx, `SELECT id, metric, metric_value, is_alarm, level, occurred_at, description, status FROM app_monitor_events ORDER BY occurred_at DESC`)
	if err != nil {
		return service.DesktopSystemStatsSource{}, fmt.Errorf("query stats monitor events: %w", err)
	}
	defer monitorRows.Close()
	monitorEvents := make([]service.DesktopMonitorRecord, 0)
	for monitorRows.Next() {
		var item service.DesktopMonitorRecord
		if err := monitorRows.Scan(&item.ID, &item.Metric, &item.Value, &item.Alarm, &item.Level, &item.OccurredAt, &item.Description, &item.Status); err != nil {
			return service.DesktopSystemStatsSource{}, fmt.Errorf("scan stats monitor event: %w", err)
		}
		monitorEvents = append(monitorEvents, item)
	}
	if err := monitorRows.Err(); err != nil {
		return service.DesktopSystemStatsSource{}, err
	}

	tenantRows, err := repository.pool.pool.Query(ctx, `SELECT id, code, name, status, period, admin_name, admin_phone FROM app_tenants ORDER BY id ASC`)
	if err != nil {
		return service.DesktopSystemStatsSource{}, fmt.Errorf("query stats tenants: %w", err)
	}
	defer tenantRows.Close()
	tenants := make([]service.DesktopTenantRecord, 0)
	for tenantRows.Next() {
		var item service.DesktopTenantRecord
		if err := tenantRows.Scan(&item.ID, &item.Code, &item.Name, &item.Status, &item.Period, &item.Admin, &item.Phone); err != nil {
			return service.DesktopSystemStatsSource{}, fmt.Errorf("scan stats tenant: %w", err)
		}
		tenants = append(tenants, item)
	}
	if err := tenantRows.Err(); err != nil {
		return service.DesktopSystemStatsSource{}, err
	}

	userRows, err := repository.pool.pool.Query(ctx, `SELECT id, uid, name, department, phone, role_name, status, updated_at, updated_by FROM app_users ORDER BY id ASC`)
	if err != nil {
		return service.DesktopSystemStatsSource{}, fmt.Errorf("query stats users: %w", err)
	}
	defer userRows.Close()
	users := make([]service.DesktopUserRecord, 0)
	for userRows.Next() {
		var item service.DesktopUserRecord
		if err := userRows.Scan(&item.ID, &item.UID, &item.Name, &item.Department, &item.Phone, &item.Role, &item.Status, &item.UpdatedAt, &item.UpdatedBy); err != nil {
			return service.DesktopSystemStatsSource{}, fmt.Errorf("scan stats user: %w", err)
		}
		users = append(users, item)
	}
	if err := userRows.Err(); err != nil {
		return service.DesktopSystemStatsSource{}, err
	}

	return service.DesktopSystemStatsSource{MonitorEvents: monitorEvents, Tenants: tenants, Users: users}, nil
}

// GetUsageStatsSource returns the raw records required for usage statistics.
func (repository *DesktopDataRepository) GetUsageStatsSource(ctx context.Context) (service.DesktopUsageStatsSource, error) {
	if repository == nil || repository.pool == nil || repository.pool.pool == nil {
		return service.DesktopUsageStatsSource{}, nil
	}

	loginRows, err := repository.pool.pool.Query(ctx, `SELECT id, log_id, tenant_code, category, user_id, name, status, occurred_at, ip, address, browser, description FROM app_login_logs ORDER BY occurred_at DESC, id DESC`)
	if err != nil {
		return service.DesktopUsageStatsSource{}, fmt.Errorf("query usage login logs: %w", err)
	}
	defer loginRows.Close()
	loginLogs := make([]service.DesktopLoginLogRecord, 0)
	for loginRows.Next() {
		var item service.DesktopLoginLogRecord
		if err := loginRows.Scan(&item.ID, &item.LogID, &item.TenantCode, &item.Category, &item.UserID, &item.Name, &item.Status, &item.Time, &item.IP, &item.Address, &item.Browser, &item.Desc); err != nil {
			return service.DesktopUsageStatsSource{}, fmt.Errorf("scan usage login log: %w", err)
		}
		loginLogs = append(loginLogs, item)
	}
	if err := loginRows.Err(); err != nil {
		return service.DesktopUsageStatsSource{}, err
	}

	operationRows, err := repository.pool.pool.Query(ctx, `SELECT id, log_id, tenant_code, module, category, user_id, name, status, occurred_at, ip, browser, description FROM app_operation_logs ORDER BY occurred_at DESC, id DESC`)
	if err != nil {
		return service.DesktopUsageStatsSource{}, fmt.Errorf("query usage operation logs: %w", err)
	}
	defer operationRows.Close()
	operationLogs := make([]service.DesktopOperationLogRecord, 0)
	for operationRows.Next() {
		var item service.DesktopOperationLogRecord
		if err := operationRows.Scan(&item.ID, &item.LogID, &item.TenantCode, &item.Module, &item.Category, &item.UserID, &item.Name, &item.Status, &item.Time, &item.IP, &item.Browser, &item.Desc); err != nil {
			return service.DesktopUsageStatsSource{}, fmt.Errorf("scan usage operation log: %w", err)
		}
		operationLogs = append(operationLogs, item)
	}
	if err := operationRows.Err(); err != nil {
		return service.DesktopUsageStatsSource{}, err
	}

	tenantRows, err := repository.pool.pool.Query(ctx, `SELECT code, name FROM app_tenants ORDER BY id ASC`)
	if err != nil {
		return service.DesktopUsageStatsSource{}, fmt.Errorf("query usage tenants: %w", err)
	}
	defer tenantRows.Close()
	tenants := make([]service.DesktopTenantRecord, 0)
	for tenantRows.Next() {
		var item service.DesktopTenantRecord
		if err := tenantRows.Scan(&item.Code, &item.Name); err != nil {
			return service.DesktopUsageStatsSource{}, fmt.Errorf("scan usage tenant: %w", err)
		}
		tenants = append(tenants, item)
	}
	if err := tenantRows.Err(); err != nil {
		return service.DesktopUsageStatsSource{}, err
	}

	return service.DesktopUsageStatsSource{LoginLogs: loginLogs, OperationLogs: operationLogs, Tenants: tenants}, nil
}

func (repository *DesktopDataRepository) CreateDesktopMenu(ctx context.Context, input service.DesktopMenuRecord) (service.DesktopMenuRecord, error) {
	var item service.DesktopMenuRecord
	err := repository.pool.pool.QueryRow(ctx, `INSERT INTO app_desktop_menus (parent_id, name, menu_level, sort_order, menu_type, icon, status, path, permission) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id, parent_id, name, menu_level, sort_order, menu_type, icon, status, path, permission`, input.ParentID, input.Name, input.Level, input.Sort, input.Type, input.Icon, input.Status, input.Path, input.Permission).Scan(&item.ID, &item.ParentID, &item.Name, &item.Level, &item.Sort, &item.Type, &item.Icon, &item.Status, &item.Path, &item.Permission)
	if err != nil {
		return service.DesktopMenuRecord{}, fmt.Errorf("create desktop menu: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) UpdateDesktopMenu(ctx context.Context, input service.DesktopMenuRecord) (service.DesktopMenuRecord, error) {
	var item service.DesktopMenuRecord
	err := repository.pool.pool.QueryRow(ctx, `UPDATE app_desktop_menus SET parent_id=$2, name=$3, menu_level=$4, sort_order=$5, menu_type=$6, icon=$7, status=$8, path=$9, permission=$10 WHERE id=$1 RETURNING id, parent_id, name, menu_level, sort_order, menu_type, icon, status, path, permission`, input.ID, input.ParentID, input.Name, input.Level, input.Sort, input.Type, input.Icon, input.Status, input.Path, input.Permission).Scan(&item.ID, &item.ParentID, &item.Name, &item.Level, &item.Sort, &item.Type, &item.Icon, &item.Status, &item.Path, &item.Permission)
	if err != nil {
		return service.DesktopMenuRecord{}, fmt.Errorf("update desktop menu: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) DeleteDesktopMenu(ctx context.Context, id int64) error {
	commandTag, err := repository.pool.pool.Exec(ctx, `DELETE FROM app_desktop_menus WHERE id=$1 OR parent_id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete desktop menu: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("delete desktop menu: %w", pgx.ErrNoRows)
	}
	return nil
}

func (repository *DesktopDataRepository) ListMenuTemplates(ctx context.Context, name string, description string, page int, pageSize int) ([]service.DesktopMenuTemplateRecord, int, error) {
	offset := (page - 1) * pageSize
	var total int
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_menu_templates WHERE ($1='' OR name ILIKE '%'||$1||'%') AND ($2='' OR description ILIKE '%'||$2||'%')`, strings.TrimSpace(name), strings.TrimSpace(description)).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count menu templates: %w", err)
	}
	rows, err := repository.pool.pool.Query(ctx, `SELECT id, name, tenant_count, description, updated_at, updated_by, created_at, created_by FROM app_menu_templates WHERE ($1='' OR name ILIKE '%'||$1||'%') AND ($2='' OR description ILIKE '%'||$2||'%') ORDER BY id ASC LIMIT $3 OFFSET $4`, strings.TrimSpace(name), strings.TrimSpace(description), pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query menu templates: %w", err)
	}
	defer rows.Close()
	items := make([]service.DesktopMenuTemplateRecord, 0)
	for rows.Next() {
		var item service.DesktopMenuTemplateRecord
		if err := rows.Scan(&item.ID, &item.Name, &item.TenantCount, &item.Description, &item.UpdatedAt, &item.UpdatedBy, &item.CreatedAt, &item.CreatedBy); err != nil {
			return nil, 0, fmt.Errorf("scan menu templates: %w", err)
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (repository *DesktopDataRepository) CreateMenuTemplate(ctx context.Context, input service.DesktopMenuTemplateRecord) (service.DesktopMenuTemplateRecord, error) {
	var item service.DesktopMenuTemplateRecord
	err := repository.pool.pool.QueryRow(ctx, `INSERT INTO app_menu_templates (name, tenant_count, description, updated_at, updated_by, created_at, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, name, tenant_count, description, updated_at, updated_by, created_at, created_by`, input.Name, input.TenantCount, input.Description, input.UpdatedAt, input.UpdatedBy, input.CreatedAt, input.CreatedBy).Scan(&item.ID, &item.Name, &item.TenantCount, &item.Description, &item.UpdatedAt, &item.UpdatedBy, &item.CreatedAt, &item.CreatedBy)
	if err != nil {
		return service.DesktopMenuTemplateRecord{}, fmt.Errorf("create menu template: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) UpdateMenuTemplate(ctx context.Context, input service.DesktopMenuTemplateRecord) (service.DesktopMenuTemplateRecord, error) {
	var item service.DesktopMenuTemplateRecord
	err := repository.pool.pool.QueryRow(ctx, `UPDATE app_menu_templates SET name=$2, tenant_count=$3, description=$4, updated_at=$5, updated_by=$6 WHERE id=$1 RETURNING id, name, tenant_count, description, updated_at, updated_by, created_at, created_by`, input.ID, input.Name, input.TenantCount, input.Description, input.UpdatedAt, input.UpdatedBy).Scan(&item.ID, &item.Name, &item.TenantCount, &item.Description, &item.UpdatedAt, &item.UpdatedBy, &item.CreatedAt, &item.CreatedBy)
	if err != nil {
		return service.DesktopMenuTemplateRecord{}, fmt.Errorf("update menu template: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) DeleteMenuTemplate(ctx context.Context, id int64) error {
	commandTag, err := repository.pool.pool.Exec(ctx, `DELETE FROM app_menu_templates WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete menu template: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("delete menu template: %w", pgx.ErrNoRows)
	}
	return nil
}

func (repository *DesktopDataRepository) ListOrgs(ctx context.Context, department string, name string, code string, page int, pageSize int) ([]service.DesktopOrgRecord, int, error) {
	offset := (page - 1) * pageSize
	var total int
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_orgs WHERE ($1='' OR department=$1) AND ($2='' OR name ILIKE '%'||$2||'%') AND ($3='' OR code ILIKE '%'||$3||'%')`, strings.TrimSpace(department), strings.TrimSpace(name), strings.TrimSpace(code)).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count orgs: %w", err)
	}
	rows, err := repository.pool.pool.Query(ctx, `SELECT id, name, code, parent_name, sort_order, category, created_at, created_by, updated_at, updated_by FROM app_orgs WHERE ($1='' OR department=$1) AND ($2='' OR name ILIKE '%'||$2||'%') AND ($3='' OR code ILIKE '%'||$3||'%') ORDER BY sort_order ASC, id ASC LIMIT $4 OFFSET $5`, strings.TrimSpace(department), strings.TrimSpace(name), strings.TrimSpace(code), pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query orgs: %w", err)
	}
	defer rows.Close()
	items := make([]service.DesktopOrgRecord, 0)
	for rows.Next() {
		var item service.DesktopOrgRecord
		if err := rows.Scan(&item.ID, &item.Name, &item.Code, &item.ParentName, &item.Sort, &item.Category, &item.CreatedAt, &item.CreatedBy, &item.UpdatedAt, &item.UpdatedBy); err != nil {
			return nil, 0, fmt.Errorf("scan orgs: %w", err)
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (repository *DesktopDataRepository) CreateOrg(ctx context.Context, input service.DesktopOrgRecord) (service.DesktopOrgRecord, error) {
	var item service.DesktopOrgRecord
	err := repository.pool.pool.QueryRow(ctx, `INSERT INTO app_orgs (name, code, parent_name, sort_order, category, created_at, created_by, updated_at, updated_by, department) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id, name, code, parent_name, sort_order, category, created_at, created_by, updated_at, updated_by`, input.Name, input.Code, input.ParentName, input.Sort, input.Category, input.CreatedAt, input.CreatedBy, input.UpdatedAt, input.UpdatedBy, inferDepartmentFromOrg(input)).Scan(&item.ID, &item.Name, &item.Code, &item.ParentName, &item.Sort, &item.Category, &item.CreatedAt, &item.CreatedBy, &item.UpdatedAt, &item.UpdatedBy)
	if err != nil {
		return service.DesktopOrgRecord{}, fmt.Errorf("create org: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) UpdateOrg(ctx context.Context, input service.DesktopOrgRecord) (service.DesktopOrgRecord, error) {
	var item service.DesktopOrgRecord
	err := repository.pool.pool.QueryRow(ctx, `UPDATE app_orgs SET name=$2, code=$3, parent_name=$4, sort_order=$5, category=$6, updated_at=$7, updated_by=$8, department=$9 WHERE id=$1 RETURNING id, name, code, parent_name, sort_order, category, created_at, created_by, updated_at, updated_by`, input.ID, input.Name, input.Code, input.ParentName, input.Sort, input.Category, input.UpdatedAt, input.UpdatedBy, inferDepartmentFromOrg(input)).Scan(&item.ID, &item.Name, &item.Code, &item.ParentName, &item.Sort, &item.Category, &item.CreatedAt, &item.CreatedBy, &item.UpdatedAt, &item.UpdatedBy)
	if err != nil {
		return service.DesktopOrgRecord{}, fmt.Errorf("update org: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) DeleteOrg(ctx context.Context, id int64) error {
	commandTag, err := repository.pool.pool.Exec(ctx, `DELETE FROM app_orgs WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete org: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("delete org: %w", pgx.ErrNoRows)
	}
	return nil
}

func (repository *DesktopDataRepository) ListTenants(ctx context.Context, name string, code string, status string, page int, pageSize int) ([]service.DesktopTenantRecord, int, error) {
	offset := (page - 1) * pageSize
	var total int
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_tenants WHERE ($1='' OR name ILIKE '%'||$1||'%') AND ($2='' OR code ILIKE '%'||$2||'%') AND ($3='' OR status=$3)`, strings.TrimSpace(name), strings.TrimSpace(code), strings.TrimSpace(status)).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count tenants: %w", err)
	}
	rows, err := repository.pool.pool.Query(ctx, `SELECT id, code, name, status, period, admin_name, admin_phone FROM app_tenants WHERE ($1='' OR name ILIKE '%'||$1||'%') AND ($2='' OR code ILIKE '%'||$2||'%') AND ($3='' OR status=$3) ORDER BY id ASC LIMIT $4 OFFSET $5`, strings.TrimSpace(name), strings.TrimSpace(code), strings.TrimSpace(status), pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query tenants: %w", err)
	}
	defer rows.Close()
	items := make([]service.DesktopTenantRecord, 0)
	for rows.Next() {
		var item service.DesktopTenantRecord
		if err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.Status, &item.Period, &item.Admin, &item.Phone); err != nil {
			return nil, 0, fmt.Errorf("scan tenants: %w", err)
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (repository *DesktopDataRepository) CreateTenant(ctx context.Context, input service.DesktopTenantRecord) (service.DesktopTenantRecord, error) {
	var item service.DesktopTenantRecord
	err := repository.pool.pool.QueryRow(ctx, `INSERT INTO app_tenants (code, name, status, period, admin_name, admin_phone) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, code, name, status, period, admin_name, admin_phone`, input.Code, input.Name, input.Status, input.Period, input.Admin, input.Phone).Scan(&item.ID, &item.Code, &item.Name, &item.Status, &item.Period, &item.Admin, &item.Phone)
	if err != nil {
		return service.DesktopTenantRecord{}, fmt.Errorf("create tenant: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) UpdateTenant(ctx context.Context, input service.DesktopTenantRecord) (service.DesktopTenantRecord, error) {
	var item service.DesktopTenantRecord
	err := repository.pool.pool.QueryRow(ctx, `UPDATE app_tenants SET code=$2, name=$3, status=$4, period=$5, admin_name=$6, admin_phone=$7 WHERE id=$1 RETURNING id, code, name, status, period, admin_name, admin_phone`, input.ID, input.Code, input.Name, input.Status, input.Period, input.Admin, input.Phone).Scan(&item.ID, &item.Code, &item.Name, &item.Status, &item.Period, &item.Admin, &item.Phone)
	if err != nil {
		return service.DesktopTenantRecord{}, fmt.Errorf("update tenant: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) DeleteTenant(ctx context.Context, id int64) error {
	commandTag, err := repository.pool.pool.Exec(ctx, `DELETE FROM app_tenants WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete tenant: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("delete tenant: %w", pgx.ErrNoRows)
	}
	return nil
}

func inferDepartmentFromOrg(input service.DesktopOrgRecord) string {
	if strings.Contains(input.Name, "生产") {
		return "生产部"
	}
	if strings.Contains(input.Name, "营销") {
		return "营销部"
	}
	if strings.Contains(input.Name, "财务") {
		return "财务部"
	}
	return "人力资源部"
}

func (repository *DesktopDataRepository) CreateAnnouncement(ctx context.Context, input service.DesktopAnnouncementRecord) (service.DesktopAnnouncementRecord, error) {
	var item service.DesktopAnnouncementRecord
	err := repository.pool.pool.QueryRow(ctx, `INSERT INTO app_announcements (title, announcement_type, status, published_at, published_by, created_at, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, title, announcement_type, status, published_at, published_by, created_at, created_by`, input.Title, input.Type, input.Status, input.PublishedAt, input.PublishedBy, input.CreatedAt, input.CreatedBy).Scan(&item.ID, &item.Title, &item.Type, &item.Status, &item.PublishedAt, &item.PublishedBy, &item.CreatedAt, &item.CreatedBy)
	if err != nil {
		return service.DesktopAnnouncementRecord{}, fmt.Errorf("create announcement: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) UpdateAnnouncement(ctx context.Context, input service.DesktopAnnouncementRecord) (service.DesktopAnnouncementRecord, error) {
	var item service.DesktopAnnouncementRecord
	err := repository.pool.pool.QueryRow(ctx, `UPDATE app_announcements SET title=$2, announcement_type=$3, status=$4, published_at=$5, published_by=$6 WHERE id=$1 RETURNING id, title, announcement_type, status, published_at, published_by, created_at, created_by`, input.ID, input.Title, input.Type, input.Status, input.PublishedAt, input.PublishedBy).Scan(&item.ID, &item.Title, &item.Type, &item.Status, &item.PublishedAt, &item.PublishedBy, &item.CreatedAt, &item.CreatedBy)
	if err != nil {
		return service.DesktopAnnouncementRecord{}, fmt.Errorf("update announcement: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) DeleteAnnouncement(ctx context.Context, id int64) error {
	commandTag, err := repository.pool.pool.Exec(ctx, `DELETE FROM app_announcements WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete announcement: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("delete announcement: %w", pgx.ErrNoRows)
	}
	return nil
}

func (repository *DesktopDataRepository) PublishAnnouncement(ctx context.Context, id int64) error {
	commandTag, err := repository.pool.pool.Exec(ctx, `UPDATE app_announcements SET status = '已发布', published_at = NOW(), published_by = '桌面端' WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("publish announcement: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("publish announcement: %w", pgx.ErrNoRows)
	}
	return nil
}

func (repository *DesktopDataRepository) ListUserDepartments(ctx context.Context) ([]service.DesktopDepartmentRecord, error) {
	rows, err := repository.pool.pool.Query(ctx, `SELECT department AS name, COUNT(*)::int AS user_count FROM app_users GROUP BY department ORDER BY department ASC`)
	if err != nil {
		return nil, fmt.Errorf("query user departments: %w", err)
	}
	defer rows.Close()
	items := make([]service.DesktopDepartmentRecord, 0)
	for rows.Next() {
		var item service.DesktopDepartmentRecord
		if err := rows.Scan(&item.Name, &item.Count); err != nil {
			return nil, fmt.Errorf("scan user departments: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (repository *DesktopDataRepository) ListOrgDepartments(ctx context.Context) ([]service.DesktopDepartmentRecord, error) {
	rows, err := repository.pool.pool.Query(ctx, `SELECT department AS name, COUNT(*)::int AS user_count FROM app_orgs GROUP BY department ORDER BY department ASC`)
	if err != nil {
		return nil, fmt.Errorf("query org departments: %w", err)
	}
	defer rows.Close()
	items := make([]service.DesktopDepartmentRecord, 0)
	for rows.Next() {
		var item service.DesktopDepartmentRecord
		if err := rows.Scan(&item.Name, &item.Count); err != nil {
			return nil, fmt.Errorf("scan org departments: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (repository *DesktopDataRepository) ListUsers(ctx context.Context, department string, uid string, name string, role string, page int, pageSize int) ([]service.DesktopUserRecord, int, error) {
	offset := (page - 1) * pageSize
	var total int
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_users WHERE ($1='' OR department=$1) AND ($2='' OR uid ILIKE '%'||$2||'%') AND ($3='' OR name ILIKE '%'||$3||'%') AND ($4='' OR role_name=$4)`, strings.TrimSpace(department), strings.TrimSpace(uid), strings.TrimSpace(name), strings.TrimSpace(role)).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	rows, err := repository.pool.pool.Query(ctx, `SELECT id, uid, name, department, phone, role_name, status, updated_at, updated_by FROM app_users WHERE ($1='' OR department=$1) AND ($2='' OR uid ILIKE '%'||$2||'%') AND ($3='' OR name ILIKE '%'||$3||'%') AND ($4='' OR role_name=$4) ORDER BY id ASC LIMIT $5 OFFSET $6`, strings.TrimSpace(department), strings.TrimSpace(uid), strings.TrimSpace(name), strings.TrimSpace(role), pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()
	items := make([]service.DesktopUserRecord, 0)
	for rows.Next() {
		var item service.DesktopUserRecord
		if err := rows.Scan(&item.ID, &item.UID, &item.Name, &item.Department, &item.Phone, &item.Role, &item.Status, &item.UpdatedAt, &item.UpdatedBy); err != nil {
			return nil, 0, fmt.Errorf("scan users: %w", err)
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (repository *DesktopDataRepository) CreateUser(ctx context.Context, input service.DesktopUserRecord) (service.DesktopUserRecord, error) {
	var item service.DesktopUserRecord
	err := repository.pool.pool.QueryRow(ctx, `INSERT INTO app_users (uid, name, department, phone, role_name, status, updated_at, updated_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, uid, name, department, phone, role_name, status, updated_at, updated_by`, input.UID, input.Name, input.Department, input.Phone, input.Role, input.Status, input.UpdatedAt, input.UpdatedBy).Scan(&item.ID, &item.UID, &item.Name, &item.Department, &item.Phone, &item.Role, &item.Status, &item.UpdatedAt, &item.UpdatedBy)
	if err != nil {
		return service.DesktopUserRecord{}, fmt.Errorf("create user: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) UpdateUser(ctx context.Context, input service.DesktopUserRecord) (service.DesktopUserRecord, error) {
	var item service.DesktopUserRecord
	err := repository.pool.pool.QueryRow(ctx, `UPDATE app_users SET uid=$2, name=$3, department=$4, phone=$5, role_name=$6, status=$7, updated_at=$8, updated_by=$9 WHERE id=$1 RETURNING id, uid, name, department, phone, role_name, status, updated_at, updated_by`, input.ID, input.UID, input.Name, input.Department, input.Phone, input.Role, input.Status, input.UpdatedAt, input.UpdatedBy).Scan(&item.ID, &item.UID, &item.Name, &item.Department, &item.Phone, &item.Role, &item.Status, &item.UpdatedAt, &item.UpdatedBy)
	if err != nil {
		return service.DesktopUserRecord{}, fmt.Errorf("update user: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) ResetUserPassword(ctx context.Context, id int64, updatedBy string, updatedAt time.Time) error {
	commandTag, err := repository.pool.pool.Exec(ctx, `UPDATE app_users SET updated_at=$2, updated_by=$3 WHERE id=$1`, id, updatedAt, updatedBy)
	if err != nil {
		return fmt.Errorf("reset user password: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("reset user password: %w", pgx.ErrNoRows)
	}
	return nil
}

func (repository *DesktopDataRepository) DeleteUser(ctx context.Context, id int64) error {
	commandTag, err := repository.pool.pool.Exec(ctx, `DELETE FROM app_users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("delete user: %w", pgx.ErrNoRows)
	}
	return nil
}

func (repository *DesktopDataRepository) ListRoles(ctx context.Context, name string, key string, status string, page int, pageSize int) ([]service.DesktopRoleRecord, int, error) {
	offset := (page - 1) * pageSize
	statusFilter := strings.TrimSpace(status)
	var total int
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_roles WHERE ($1='' OR name ILIKE '%'||$1||'%') AND ($2='' OR role_key ILIKE '%'||$2||'%') AND ($3='' OR status = ($3='正常' OR $3='active'))`, strings.TrimSpace(name), strings.TrimSpace(key), statusFilter).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count roles: %w", err)
	}
	rows, err := repository.pool.pool.Query(ctx, `SELECT id, name, role_key, display_order, status, created_at FROM app_roles WHERE ($1='' OR name ILIKE '%'||$1||'%') AND ($2='' OR role_key ILIKE '%'||$2||'%') AND ($3='' OR status = ($3='正常' OR $3='active')) ORDER BY display_order ASC, id ASC LIMIT $4 OFFSET $5`, strings.TrimSpace(name), strings.TrimSpace(key), statusFilter, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query roles: %w", err)
	}
	defer rows.Close()
	items := make([]service.DesktopRoleRecord, 0)
	for rows.Next() {
		var item service.DesktopRoleRecord
		if err := rows.Scan(&item.ID, &item.Name, &item.Key, &item.Order, &item.Status, &item.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan roles: %w", err)
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (repository *DesktopDataRepository) CreateRole(ctx context.Context, input service.DesktopRoleRecord) (service.DesktopRoleRecord, error) {
	var item service.DesktopRoleRecord
	err := repository.pool.pool.QueryRow(ctx, `INSERT INTO app_roles (name, role_key, display_order, status, created_at) VALUES ($1,$2,$3,$4,$5) RETURNING id, name, role_key, display_order, status, created_at`, input.Name, input.Key, input.Order, input.Status, input.CreatedAt).Scan(&item.ID, &item.Name, &item.Key, &item.Order, &item.Status, &item.CreatedAt)
	if err != nil {
		return service.DesktopRoleRecord{}, fmt.Errorf("create role: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) UpdateRole(ctx context.Context, input service.DesktopRoleRecord) (service.DesktopRoleRecord, error) {
	var item service.DesktopRoleRecord
	err := repository.pool.pool.QueryRow(ctx, `UPDATE app_roles SET name=$2, role_key=$3, display_order=$4, status=$5 WHERE id=$1 RETURNING id, name, role_key, display_order, status, created_at`, input.ID, input.Name, input.Key, input.Order, input.Status).Scan(&item.ID, &item.Name, &item.Key, &item.Order, &item.Status, &item.CreatedAt)
	if err != nil {
		return service.DesktopRoleRecord{}, fmt.Errorf("update role: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) DeleteRole(ctx context.Context, id int64) error {
	commandTag, err := repository.pool.pool.Exec(ctx, `DELETE FROM app_roles WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete role: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("delete role: %w", pgx.ErrNoRows)
	}
	return nil
}

func (repository *DesktopDataRepository) GetRoleDataPermission(ctx context.Context, roleID int64) (service.DesktopRoleDataPermissionRecord, error) {
	var item service.DesktopRoleDataPermissionRecord
	err := repository.pool.pool.QueryRow(ctx, `SELECT role_id, scope, departments_json FROM app_role_data_permissions WHERE role_id=$1`, roleID).Scan(&item.RoleID, &item.Scope, &item.Departments)
	if err != nil {
		if err == pgx.ErrNoRows {
			return service.DesktopRoleDataPermissionRecord{RoleID: roleID, Scope: "all", Departments: "[]"}, nil
		}
		return service.DesktopRoleDataPermissionRecord{}, fmt.Errorf("get role data permission: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) SaveRoleDataPermission(ctx context.Context, input service.DesktopRoleDataPermissionRecord) (service.DesktopRoleDataPermissionRecord, error) {
	_, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_role_data_permissions (role_id, scope, departments_json) VALUES ($1,$2,$3) ON CONFLICT (role_id) DO UPDATE SET scope=EXCLUDED.scope, departments_json=EXCLUDED.departments_json`, input.RoleID, input.Scope, input.Departments)
	if err != nil {
		return service.DesktopRoleDataPermissionRecord{}, fmt.Errorf("save role data permission: %w", err)
	}
	return input, nil
}

// GetSystemSettings returns one settings row.
func (repository *DesktopDataRepository) GetSystemSettings(ctx context.Context) (service.DesktopSystemSettingsRecord, error) {
	if repository == nil || repository.pool == nil || repository.pool.pool == nil {
		return service.DesktopSystemSettingsRecord{}, nil
	}
	var item service.DesktopSystemSettingsRecord
	err := repository.pool.pool.QueryRow(ctx, `SELECT website_title, system_logo, theme, icp, copyright, require_strong_password, login_fail_limit, login_lock_minutes FROM app_system_settings WHERE id = 1`).Scan(&item.WebsiteTitle, &item.SystemLogo, &item.Theme, &item.ICP, &item.Copyright, &item.RequireStrongPassword, &item.LoginFailLimit, &item.LoginLockMinutes)
	if err != nil {
		if err == pgx.ErrNoRows {
			return service.DesktopSystemSettingsRecord{}, nil
		}
		return service.DesktopSystemSettingsRecord{}, fmt.Errorf("query system settings: %w", err)
	}
	return item, nil
}

// SaveSystemSettings upserts one settings row.
func (repository *DesktopDataRepository) SaveSystemSettings(ctx context.Context, input service.DesktopSystemSettingsRecord) (service.DesktopSystemSettingsRecord, error) {
	if repository == nil || repository.pool == nil || repository.pool.pool == nil {
		return input, nil
	}
	_, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_system_settings (id, website_title, system_logo, theme, icp, copyright, require_strong_password, login_fail_limit, login_lock_minutes) VALUES (1,$1,$2,$3,$4,$5,$6,$7,$8)
	ON CONFLICT (id) DO UPDATE SET website_title=EXCLUDED.website_title, system_logo=EXCLUDED.system_logo, theme=EXCLUDED.theme, icp=EXCLUDED.icp, copyright=EXCLUDED.copyright, require_strong_password=EXCLUDED.require_strong_password, login_fail_limit=EXCLUDED.login_fail_limit, login_lock_minutes=EXCLUDED.login_lock_minutes`,
		input.WebsiteTitle, input.SystemLogo, input.Theme, input.ICP, input.Copyright, input.RequireStrongPassword, input.LoginFailLimit, input.LoginLockMinutes)
	if err != nil {
		return service.DesktopSystemSettingsRecord{}, fmt.Errorf("save system settings: %w", err)
	}
	return input, nil
}

// ListDicts returns dictionaries and nested items.
func (repository *DesktopDataRepository) ListDicts(ctx context.Context, name string, dictID string) ([]service.DesktopDictRecord, error) {
	if repository == nil || repository.pool == nil || repository.pool.pool == nil {
		return nil, nil
	}
	rows, err := repository.pool.pool.Query(ctx, `SELECT id, dict_id, name, description, modified_at, modified_by, created_at, created_by FROM app_dicts WHERE ($1 = '' OR name ILIKE '%' || $1 || '%') AND ($2 = '' OR dict_id ILIKE '%' || $2 || '%') ORDER BY id ASC`, strings.TrimSpace(name), strings.TrimSpace(dictID))
	if err != nil {
		return nil, fmt.Errorf("query dicts: %w", err)
	}
	defer rows.Close()
	items := make([]service.DesktopDictRecord, 0)
	for rows.Next() {
		var item service.DesktopDictRecord
		if err := rows.Scan(&item.ID, &item.DictID, &item.Name, &item.Description, &item.ModifiedAt, &item.ModifiedBy, &item.CreatedAt, &item.CreatedBy); err != nil {
			return nil, fmt.Errorf("scan dicts: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate dicts: %w", err)
	}
	for index, item := range items {
		subRows, err := repository.pool.pool.Query(ctx, `SELECT id, item_index, label, key_val, style_type, modified_at, modified_by, created_at, created_by, badge_color FROM app_dict_items WHERE dict_ref_id = $1 ORDER BY item_index ASC, id ASC`, item.ID)
		if err != nil {
			return nil, fmt.Errorf("query dict items: %w", err)
		}
		children := make([]service.DesktopDictItemRecord, 0)
		for subRows.Next() {
			var child service.DesktopDictItemRecord
			if err := subRows.Scan(&child.ID, &child.Index, &child.Label, &child.KeyVal, &child.StyleType, &child.ModifiedAt, &child.ModifiedBy, &child.CreatedAt, &child.CreatedBy, &child.BadgeColor); err != nil {
				subRows.Close()
				return nil, fmt.Errorf("scan dict items: %w", err)
			}
			children = append(children, child)
		}
		subRows.Close()
		items[index].Items = children
	}
	return items, nil
}

func (repository *DesktopDataRepository) CreateDict(ctx context.Context, input service.DesktopDictRecord) (service.DesktopDictRecord, error) {
	var item service.DesktopDictRecord
	err := repository.pool.pool.QueryRow(ctx, `INSERT INTO app_dicts (dict_id, name, description, modified_at, modified_by, created_at, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, dict_id, name, description, modified_at, modified_by, created_at, created_by`, input.DictID, input.Name, input.Description, input.ModifiedAt, input.ModifiedBy, input.CreatedAt, input.CreatedBy).Scan(&item.ID, &item.DictID, &item.Name, &item.Description, &item.ModifiedAt, &item.ModifiedBy, &item.CreatedAt, &item.CreatedBy)
	if err != nil {
		return service.DesktopDictRecord{}, fmt.Errorf("create dict: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) UpdateDict(ctx context.Context, input service.DesktopDictRecord) (service.DesktopDictRecord, error) {
	var item service.DesktopDictRecord
	err := repository.pool.pool.QueryRow(ctx, `UPDATE app_dicts SET dict_id=$2, name=$3, description=$4, modified_at=$5, modified_by=$6 WHERE id=$1 RETURNING id, dict_id, name, description, modified_at, modified_by, created_at, created_by`, input.ID, input.DictID, input.Name, input.Description, input.ModifiedAt, input.ModifiedBy).Scan(&item.ID, &item.DictID, &item.Name, &item.Description, &item.ModifiedAt, &item.ModifiedBy, &item.CreatedAt, &item.CreatedBy)
	if err != nil {
		return service.DesktopDictRecord{}, fmt.Errorf("update dict: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) DeleteDict(ctx context.Context, id int64) error {
	commandTag, err := repository.pool.pool.Exec(ctx, `DELETE FROM app_dicts WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete dict: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("delete dict: %w", pgx.ErrNoRows)
	}
	return nil
}

func (repository *DesktopDataRepository) CreateDictItem(ctx context.Context, dictID int64, input service.DesktopDictItemRecord) (service.DesktopDictItemRecord, error) {
	var item service.DesktopDictItemRecord
	err := repository.pool.pool.QueryRow(ctx, `INSERT INTO app_dict_items (dict_ref_id, item_index, label, key_val, style_type, modified_at, modified_by, created_at, created_by, badge_color) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id, item_index, label, key_val, style_type, modified_at, modified_by, created_at, created_by, badge_color`, dictID, input.Index, input.Label, input.KeyVal, input.StyleType, input.ModifiedAt, input.ModifiedBy, input.CreatedAt, input.CreatedBy, input.BadgeColor).Scan(&item.ID, &item.Index, &item.Label, &item.KeyVal, &item.StyleType, &item.ModifiedAt, &item.ModifiedBy, &item.CreatedAt, &item.CreatedBy, &item.BadgeColor)
	if err != nil {
		return service.DesktopDictItemRecord{}, fmt.Errorf("create dict item: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) UpdateDictItem(ctx context.Context, input service.DesktopDictItemRecord) (service.DesktopDictItemRecord, error) {
	var item service.DesktopDictItemRecord
	err := repository.pool.pool.QueryRow(ctx, `UPDATE app_dict_items SET item_index=$2, label=$3, key_val=$4, style_type=$5, modified_at=$6, modified_by=$7, badge_color=$8 WHERE id=$1 RETURNING id, item_index, label, key_val, style_type, modified_at, modified_by, created_at, created_by, badge_color`, input.ID, input.Index, input.Label, input.KeyVal, input.StyleType, input.ModifiedAt, input.ModifiedBy, input.BadgeColor).Scan(&item.ID, &item.Index, &item.Label, &item.KeyVal, &item.StyleType, &item.ModifiedAt, &item.ModifiedBy, &item.CreatedAt, &item.CreatedBy, &item.BadgeColor)
	if err != nil {
		return service.DesktopDictItemRecord{}, fmt.Errorf("update dict item: %w", err)
	}
	return item, nil
}

func (repository *DesktopDataRepository) DeleteDictItem(ctx context.Context, id int64) error {
	commandTag, err := repository.pool.pool.Exec(ctx, `DELETE FROM app_dict_items WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete dict item: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("delete dict item: %w", pgx.ErrNoRows)
	}
	return nil
}

// ListLoginLogs returns login logs ordered by occurrence time.
func (repository *DesktopDataRepository) ListLoginLogs(ctx context.Context, name string, status string, startDate string, endDate string, page int, pageSize int) ([]service.DesktopLoginLogRecord, int, error) {
	if repository == nil || repository.pool == nil || repository.pool.pool == nil {
		return nil, 0, nil
	}
	offset := (page - 1) * pageSize
	var total int
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_login_logs WHERE ($1 = '' OR name ILIKE '%' || $1 || '%') AND ($2 = '' OR status = $2) AND ($3 = '' OR occurred_at::date >= $3::date) AND ($4 = '' OR occurred_at::date <= $4::date)`, strings.TrimSpace(name), strings.TrimSpace(status), strings.TrimSpace(startDate), strings.TrimSpace(endDate)).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count login logs: %w", err)
	}
	rows, err := repository.pool.pool.Query(ctx, `SELECT id, log_id, category, user_id, name, status, occurred_at, ip, address, browser, description FROM app_login_logs WHERE ($1 = '' OR name ILIKE '%' || $1 || '%') AND ($2 = '' OR status = $2) AND ($3 = '' OR occurred_at::date >= $3::date) AND ($4 = '' OR occurred_at::date <= $4::date) ORDER BY occurred_at DESC, id DESC LIMIT $5 OFFSET $6`, strings.TrimSpace(name), strings.TrimSpace(status), strings.TrimSpace(startDate), strings.TrimSpace(endDate), pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query login logs: %w", err)
	}
	defer rows.Close()
	items := make([]service.DesktopLoginLogRecord, 0)
	for rows.Next() {
		var item service.DesktopLoginLogRecord
		if err := rows.Scan(&item.ID, &item.LogID, &item.Category, &item.UserID, &item.Name, &item.Status, &item.Time, &item.IP, &item.Address, &item.Browser, &item.Desc); err != nil {
			return nil, 0, fmt.Errorf("scan login logs: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate login logs: %w", err)
	}
	return items, total, nil
}

// ListOperationLogs returns operation logs ordered by occurrence time.
func (repository *DesktopDataRepository) ListOperationLogs(ctx context.Context, name string, status string, startDate string, endDate string, page int, pageSize int) ([]service.DesktopOperationLogRecord, int, error) {
	if repository == nil || repository.pool == nil || repository.pool.pool == nil {
		return nil, 0, nil
	}
	offset := (page - 1) * pageSize
	var total int
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_operation_logs WHERE ($1 = '' OR name ILIKE '%' || $1 || '%') AND ($2 = '' OR status = $2) AND ($3 = '' OR occurred_at::date >= $3::date) AND ($4 = '' OR occurred_at::date <= $4::date)`, strings.TrimSpace(name), strings.TrimSpace(status), strings.TrimSpace(startDate), strings.TrimSpace(endDate)).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count operation logs: %w", err)
	}
	rows, err := repository.pool.pool.Query(ctx, `SELECT id, log_id, module, category, user_id, name, status, occurred_at, ip, browser, description FROM app_operation_logs WHERE ($1 = '' OR name ILIKE '%' || $1 || '%') AND ($2 = '' OR status = $2) AND ($3 = '' OR occurred_at::date >= $3::date) AND ($4 = '' OR occurred_at::date <= $4::date) ORDER BY occurred_at DESC, id DESC LIMIT $5 OFFSET $6`, strings.TrimSpace(name), strings.TrimSpace(status), strings.TrimSpace(startDate), strings.TrimSpace(endDate), pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query operation logs: %w", err)
	}
	defer rows.Close()
	items := make([]service.DesktopOperationLogRecord, 0)
	for rows.Next() {
		var item service.DesktopOperationLogRecord
		if err := rows.Scan(&item.ID, &item.LogID, &item.Module, &item.Category, &item.UserID, &item.Name, &item.Status, &item.Time, &item.IP, &item.Browser, &item.Desc); err != nil {
			return nil, 0, fmt.Errorf("scan operation logs: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate operation logs: %w", err)
	}
	return items, total, nil
}

// ListAnnouncements returns announcement rows ordered by creation time.
func (repository *DesktopDataRepository) ListAnnouncements(ctx context.Context, title string, announcementType string, status string, page int, pageSize int) ([]service.DesktopAnnouncementRecord, int, error) {
	if repository == nil || repository.pool == nil || repository.pool.pool == nil {
		return nil, 0, nil
	}
	offset := (page - 1) * pageSize
	var total int
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_announcements WHERE ($1 = '' OR title ILIKE '%' || $1 || '%') AND ($2 = '' OR announcement_type = $2) AND ($3 = '' OR status = $3)`, strings.TrimSpace(title), strings.TrimSpace(announcementType), strings.TrimSpace(status)).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count announcements: %w", err)
	}
	rows, err := repository.pool.pool.Query(ctx, `SELECT id, title, announcement_type, status, published_at, published_by, created_at, created_by FROM app_announcements WHERE ($1 = '' OR title ILIKE '%' || $1 || '%') AND ($2 = '' OR announcement_type = $2) AND ($3 = '' OR status = $3) ORDER BY created_at DESC, id DESC LIMIT $4 OFFSET $5`, strings.TrimSpace(title), strings.TrimSpace(announcementType), strings.TrimSpace(status), pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query announcements: %w", err)
	}
	defer rows.Close()
	items := make([]service.DesktopAnnouncementRecord, 0)
	for rows.Next() {
		var item service.DesktopAnnouncementRecord
		if err := rows.Scan(&item.ID, &item.Title, &item.Type, &item.Status, &item.PublishedAt, &item.PublishedBy, &item.CreatedAt, &item.CreatedBy); err != nil {
			return nil, 0, fmt.Errorf("scan announcements: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate announcements: %w", err)
	}
	return items, total, nil
}

// ListMessages returns message rows ordered by publish time.
func (repository *DesktopDataRepository) ListMessages(ctx context.Context, messageType string, keyword string, page int, pageSize int) ([]service.DesktopMessageRecord, int, error) {
	if repository == nil || repository.pool == nil || repository.pool.pool == nil {
		return nil, 0, nil
	}
	offset := (page - 1) * pageSize
	var total int
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_messages WHERE ($1 = '' OR type = $1) AND ($2 = '' OR title ILIKE '%' || $2 || '%' OR description ILIKE '%' || $2 || '%')`, strings.TrimSpace(messageType), strings.TrimSpace(keyword)).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count desktop messages: %w", err)
	}

	query := `
		SELECT id, type, title, description, published_at, author, is_read
		FROM app_messages
		WHERE ($1 = '' OR type = $1)
		  AND ($2 = '' OR title ILIKE '%' || $2 || '%' OR description ILIKE '%' || $2 || '%')
		ORDER BY published_at DESC, id DESC
		LIMIT $3 OFFSET $4`

	rows, err := repository.pool.pool.Query(ctx, query, strings.TrimSpace(messageType), strings.TrimSpace(keyword), pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query desktop messages: %w", err)
	}
	defer rows.Close()

	items := make([]service.DesktopMessageRecord, 0)
	for rows.Next() {
		var item service.DesktopMessageRecord
		if err := rows.Scan(&item.ID, &item.Type, &item.Title, &item.Description, &item.PublishedAt, &item.Author, &item.Read); err != nil {
			return nil, 0, fmt.Errorf("scan desktop message: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate desktop messages: %w", err)
	}

	return items, total, nil
}

// ListMonitorEvents returns monitor events ordered by occurrence time.
func (repository *DesktopDataRepository) ListMonitorEvents(ctx context.Context, metric string, level string, page int, pageSize int) ([]service.DesktopMonitorRecord, int, error) {
	if repository == nil || repository.pool == nil || repository.pool.pool == nil {
		return nil, 0, nil
	}
	offset := (page - 1) * pageSize
	var total int
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_monitor_events WHERE ($1 = '' OR metric ILIKE '%' || $1 || '%') AND ($2 = '' OR level = $2)`, strings.TrimSpace(metric), strings.TrimSpace(level)).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count monitor events: %w", err)
	}

	query := `
		SELECT id, metric, metric_value, is_alarm, level, occurred_at, description, status
		FROM app_monitor_events
		WHERE ($1 = '' OR metric ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR level = $2)
		ORDER BY occurred_at DESC, id DESC
		LIMIT $3 OFFSET $4`

	rows, err := repository.pool.pool.Query(ctx, query, strings.TrimSpace(metric), strings.TrimSpace(level), pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query monitor events: %w", err)
	}
	defer rows.Close()

	items := make([]service.DesktopMonitorRecord, 0)
	for rows.Next() {
		var item service.DesktopMonitorRecord
		if err := rows.Scan(&item.ID, &item.Metric, &item.Value, &item.Alarm, &item.Level, &item.OccurredAt, &item.Description, &item.Status); err != nil {
			return nil, 0, fmt.Errorf("scan monitor event: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate monitor events: %w", err)
	}

	return items, total, nil
}

// GetMonitorSummary returns aggregate monitor metrics.
func (repository *DesktopDataRepository) GetMonitorSummary(ctx context.Context) (service.DesktopMonitorSummaryRecord, error) {
	if repository == nil || repository.pool == nil || repository.pool.pool == nil {
		return service.DesktopMonitorSummaryRecord{}, nil
	}

	query := `
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE is_alarm = TRUE) AS active_alarm_count,
			COALESCE(
				CASE MIN(
					CASE level
						WHEN '一级告警' THEN 1
						WHEN '二级告警' THEN 2
						WHEN '三级告警' THEN 3
						ELSE 4
					END
				)
					WHEN 1 THEN '一级告警'
					WHEN 2 THEN '二级告警'
					WHEN 3 THEN '三级告警'
					ELSE '无告警'
				END,
			'无告警') AS highest_level,
			COALESCE((SELECT last_collected_at FROM app_monitor_collection_state WHERE id = 1), MAX(occurred_at), NOW()) AS last_collected_at
		FROM app_monitor_events`

	var summary service.DesktopMonitorSummaryRecord
	if err := repository.pool.pool.QueryRow(ctx, query).Scan(&summary.Total, &summary.ActiveAlarmCount, &summary.HighestLevel, &summary.LastCollectedAt); err != nil {
		return service.DesktopMonitorSummaryRecord{}, fmt.Errorf("query monitor summary: %w", err)
	}

	return summary, nil
}

// CollectMonitor records one manual monitor collection timestamp.
func (repository *DesktopDataRepository) CollectMonitor(ctx context.Context, collectedAt time.Time) error {
	if repository == nil || repository.pool == nil || repository.pool.pool == nil {
		return nil
	}

	_, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_monitor_collection_state (id, last_collected_at) VALUES (1, $1) ON CONFLICT (id) DO UPDATE SET last_collected_at = EXCLUDED.last_collected_at`, collectedAt)
	if err != nil {
		return fmt.Errorf("collect monitor: %w", err)
	}

	return nil
}

// MarkMessageRead marks one stored message as read.
func (repository *DesktopDataRepository) MarkMessageRead(ctx context.Context, id int64) error {
	if repository == nil || repository.pool == nil || repository.pool.pool == nil {
		return nil
	}

	commandTag, err := repository.pool.pool.Exec(ctx, `UPDATE app_messages SET is_read = TRUE WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("mark desktop message read: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("mark desktop message read: %w", pgx.ErrNoRows)
	}
	return nil
}

// MarkMessagesRead marks multiple stored messages as read.
func (repository *DesktopDataRepository) MarkMessagesRead(ctx context.Context, ids []int64) error {
	if repository == nil || repository.pool == nil || repository.pool.pool == nil {
		return nil
	}
	if len(ids) == 0 {
		return nil
	}

	commandTag, err := repository.pool.pool.Exec(ctx, `UPDATE app_messages SET is_read = TRUE WHERE id = ANY($1)`, ids)
	if err != nil {
		return fmt.Errorf("mark desktop messages read: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("mark desktop messages read: %w", pgx.ErrNoRows)
	}
	return nil
}

// UpdateMonitorStatus updates one stored monitor status.
func (repository *DesktopDataRepository) UpdateMonitorStatus(ctx context.Context, id int64, status string) error {
	if repository == nil || repository.pool == nil || repository.pool.pool == nil {
		return nil
	}

	commandTag, err := repository.pool.pool.Exec(ctx, `UPDATE app_monitor_events SET status = $2, is_alarm = CASE WHEN $2 = 'ignored' THEN FALSE ELSE is_alarm END WHERE id = $1`, id, status)
	if err != nil {
		return fmt.Errorf("update monitor status: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("update monitor status: %w", pgx.ErrNoRows)
	}
	return nil
}

func (repository *DesktopDataRepository) seedMessages(ctx context.Context) error {
	var unreadCount int64
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_messages WHERE is_read = FALSE`).Scan(&unreadCount); err != nil {
		return fmt.Errorf("count unread desktop messages: %w", err)
	}
	const targetUnreadMessages = 10
	if unreadCount >= targetUnreadMessages {
		return nil
	}

	seedRows := make([]service.DesktopMessageRecord, 0, targetUnreadMessages)
	for index := int64(0); index < targetUnreadMessages-unreadCount; index++ {
		seedRows = append(seedRows, service.DesktopMessageRecord{
			Type:        "通知消息",
			Title:       fmt.Sprintf("桌面端待处理消息提醒 %02d", index+1),
			Description: fmt.Sprintf("第 %02d 条未读消息，提醒及时查看工作台中的待办事项与系统通知。", index+1),
			PublishedAt: time.Now().Add(-time.Duration(index+1) * time.Hour),
			Author:      "系统管理员",
			Read:        false,
		})
	}

	for _, item := range seedRows {
		if _, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_messages (type, title, description, published_at, author, is_read) VALUES ($1, $2, $3, $4, $5, $6)`, item.Type, item.Title, item.Description, item.PublishedAt, item.Author, item.Read); err != nil {
			return fmt.Errorf("seed desktop messages: %w", err)
		}
	}

	return nil
}

func (repository *DesktopDataRepository) seedMonitorEvents(ctx context.Context) error {
	var count int64
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_monitor_events`).Scan(&count); err != nil {
		return fmt.Errorf("count monitor events: %w", err)
	}
	if count > 0 {
		return nil
	}

	seedRows := []service.DesktopMonitorRecord{
		{Metric: "CPU使用率", Value: "67%", Alarm: true, Level: "一级告警", OccurredAt: time.Now().Add(-20 * time.Minute), Description: "推理节点 CPU 使用率持续高于 65%，建议扩容或错峰处理任务。", Status: "pending"},
		{Metric: "GPU显存使用率", Value: "82%", Alarm: true, Level: "二级告警", OccurredAt: time.Now().Add(-45 * time.Minute), Description: "GPU 显存已接近阈值，请关注模型装载情况。", Status: "processing"},
		{Metric: "消息投递延迟", Value: "240ms", Alarm: false, Level: "三级告警", OccurredAt: time.Now().Add(-90 * time.Minute), Description: "消息队列延迟已恢复到可接受范围。", Status: "resolved"},
		{Metric: "磁盘使用率", Value: "58%", Alarm: false, Level: "三级告警", OccurredAt: time.Now().Add(-3 * time.Hour), Description: "本地缓存目录使用率稳定。", Status: "resolved"},
	}

	for _, item := range seedRows {
		if _, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_monitor_events (metric, metric_value, is_alarm, level, occurred_at, description, status) VALUES ($1, $2, $3, $4, $5, $6, $7)`, item.Metric, item.Value, item.Alarm, item.Level, item.OccurredAt, item.Description, item.Status); err != nil {
			return fmt.Errorf("seed monitor events: %w", err)
		}
	}
	if _, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_monitor_collection_state (id, last_collected_at) VALUES (1, NOW()) ON CONFLICT (id) DO NOTHING`); err != nil {
		return fmt.Errorf("seed monitor collection state: %w", err)
	}

	return nil
}

func (repository *DesktopDataRepository) seedAnnouncements(ctx context.Context) error {
	var publishedCount int64
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_announcements WHERE status = '已发布'`).Scan(&publishedCount); err != nil {
		return fmt.Errorf("count published announcements: %w", err)
	}
	const targetPublishedAnnouncements = 10
	if publishedCount >= targetPublishedAnnouncements {
		return nil
	}
	seedRows := make([]service.DesktopAnnouncementRecord, 0, targetPublishedAnnouncements)
	for index := int64(0); index < targetPublishedAnnouncements-publishedCount; index++ {
		publishedAt := time.Now().Add(-time.Duration(index+1) * 2 * time.Hour)
		seedRows = append(seedRows, service.DesktopAnnouncementRecord{
			Title:       fmt.Sprintf("桌面端公告通知 %02d", index+1),
			Type:        "公告",
			Status:      "已发布",
			PublishedAt: publishedAt,
			PublishedBy: "系统管理员",
			CreatedAt:   publishedAt.Add(-30 * time.Minute),
			CreatedBy:   "系统管理员",
		})
	}
	for _, item := range seedRows {
		if _, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_announcements (title, announcement_type, status, published_at, published_by, created_at, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7)`, item.Title, item.Type, item.Status, item.PublishedAt, item.PublishedBy, item.CreatedAt, item.CreatedBy); err != nil {
			return fmt.Errorf("seed announcements: %w", err)
		}
	}
	return nil
}

func (repository *DesktopDataRepository) seedLoginLogs(ctx context.Context) error {
	rows := []service.DesktopLoginLogRecord{{LogID: "2345641204845120", TenantCode: "FQJT", Category: "登录", UserID: "100001", Name: "张三", Status: "成功", Time: time.Now().Add(-2 * time.Hour), IP: "10.111.123.131", Address: "广东省深圳市福田区", Browser: "Chrome 11", Desc: "登录成功"}, {LogID: "2345641204845121", TenantCode: "HLGJT", Category: "登录", UserID: "100002", Name: "李四", Status: "失败", Time: time.Now().Add(-90 * time.Minute), IP: "10.111.123.132", Address: "上海市浦东新区", Browser: "Chrome 11", Desc: "密码错误，登录失败"}, {LogID: "2345641204845122", TenantCode: "", Category: "登出", UserID: "100003", Name: "王五", Status: "成功", Time: time.Now().Add(-30 * time.Minute), IP: "10.111.123.133", Address: "北京市朝阳区", Browser: "Chrome 11", Desc: "登出成功"}}
	for _, item := range rows {
		if _, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_login_logs (log_id, tenant_code, category, user_id, name, status, occurred_at, ip, address, browser, description) SELECT $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11 WHERE NOT EXISTS (SELECT 1 FROM app_login_logs WHERE log_id = $1)`, item.LogID, item.TenantCode, item.Category, item.UserID, item.Name, item.Status, item.Time, item.IP, item.Address, item.Browser, item.Desc); err != nil {
			return fmt.Errorf("seed login logs: %w", err)
		}
	}
	if err := repository.seedTenantUsageLoginLogs(ctx); err != nil {
		return err
	}
	return nil
}

func (repository *DesktopDataRepository) seedOperationLogs(ctx context.Context) error {
	rows := []service.DesktopOperationLogRecord{{LogID: "3345641204845120", TenantCode: "FQJT", Module: "租户管理", Category: "修改", UserID: "100001", Name: "张三", Status: "成功", Time: time.Now().Add(-4 * time.Hour), IP: "10.111.123.131", Browser: "Chrome 11", Desc: "修改租户成功"}, {LogID: "3345641204845121", TenantCode: "HLGJT", Module: "菜单管理", Category: "删除", UserID: "100002", Name: "李四", Status: "失败", Time: time.Now().Add(-2 * time.Hour), IP: "10.111.123.132", Browser: "Chrome 11", Desc: "删除菜单失败"}, {LogID: "3345641204845122", TenantCode: "", Module: "用户管理", Category: "新增", UserID: "100003", Name: "王五", Status: "成功", Time: time.Now().Add(-1 * time.Hour), IP: "10.111.123.133", Browser: "Chrome 11", Desc: "新增用户成功"}}
	for _, item := range rows {
		if _, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_operation_logs (log_id, tenant_code, module, category, user_id, name, status, occurred_at, ip, browser, description) SELECT $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11 WHERE NOT EXISTS (SELECT 1 FROM app_operation_logs WHERE log_id = $1)`, item.LogID, item.TenantCode, item.Module, item.Category, item.UserID, item.Name, item.Status, item.Time, item.IP, item.Browser, item.Desc); err != nil {
			return fmt.Errorf("seed operation logs: %w", err)
		}
	}
	if err := repository.seedTenantUsageOperationLogs(ctx); err != nil {
		return err
	}
	return nil
}

func (repository *DesktopDataRepository) seedTenantUsageLoginLogs(ctx context.Context) error {
	tenantRows, err := repository.pool.pool.Query(ctx, `SELECT code, name FROM app_tenants WHERE code <> '' ORDER BY id ASC`)
	if err != nil {
		return fmt.Errorf("query tenants for login log seed: %w", err)
	}
	defer tenantRows.Close()

	index := 0
	for tenantRows.Next() {
		var tenantCode string
		var tenantName string
		if err := tenantRows.Scan(&tenantCode, &tenantName); err != nil {
			return fmt.Errorf("scan tenant for login log seed: %w", err)
		}
		var count int64
		if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_login_logs WHERE tenant_code = $1`, tenantCode).Scan(&count); err != nil {
			return fmt.Errorf("count tenant login logs: %w", err)
		}
		if count > 0 {
			index++
			continue
		}
		logID := fmt.Sprintf("tenant-login-seed-%s", tenantCode)
		occurredAt := time.Now().Add(-time.Duration(index+1) * 75 * time.Minute)
		if _, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_login_logs (log_id, tenant_code, category, user_id, name, status, occurred_at, ip, address, browser, description) SELECT $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11 WHERE NOT EXISTS (SELECT 1 FROM app_login_logs WHERE log_id = $1)`, logID, tenantCode, "登录", fmt.Sprintf("tenant-%02d", index+1), tenantName, "成功", occurredAt, fmt.Sprintf("10.10.10.%d", index+10), "北京市朝阳区", "Chrome 11", fmt.Sprintf("%s 登录成功", tenantName)); err != nil {
			return fmt.Errorf("seed tenant login log: %w", err)
		}
		index++
	}
	if err := tenantRows.Err(); err != nil {
		return fmt.Errorf("iterate tenants for login log seed: %w", err)
	}
	return nil
}

func (repository *DesktopDataRepository) seedTenantUsageOperationLogs(ctx context.Context) error {
	tenantRows, err := repository.pool.pool.Query(ctx, `SELECT code, name FROM app_tenants WHERE code <> '' ORDER BY id ASC`)
	if err != nil {
		return fmt.Errorf("query tenants for operation log seed: %w", err)
	}
	defer tenantRows.Close()

	modules := []string{"用户管理", "租户管理", "菜单管理", "系统监控"}
	index := 0
	for tenantRows.Next() {
		var tenantCode string
		var tenantName string
		if err := tenantRows.Scan(&tenantCode, &tenantName); err != nil {
			return fmt.Errorf("scan tenant for operation log seed: %w", err)
		}
		var count int64
		if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_operation_logs WHERE tenant_code = $1`, tenantCode).Scan(&count); err != nil {
			return fmt.Errorf("count tenant operation logs: %w", err)
		}
		if count > 0 {
			index++
			continue
		}
		logID := fmt.Sprintf("tenant-operation-seed-%s", tenantCode)
		occurredAt := time.Now().Add(-time.Duration(index+1) * 90 * time.Minute)
		module := modules[index%len(modules)]
		if _, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_operation_logs (log_id, tenant_code, module, category, user_id, name, status, occurred_at, ip, browser, description) SELECT $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11 WHERE NOT EXISTS (SELECT 1 FROM app_operation_logs WHERE log_id = $1)`, logID, tenantCode, module, "查看", fmt.Sprintf("tenant-%02d", index+1), tenantName, "成功", occurredAt, fmt.Sprintf("10.10.20.%d", index+10), "Chrome 11", fmt.Sprintf("%s 查看%s", tenantName, module)); err != nil {
			return fmt.Errorf("seed tenant operation log: %w", err)
		}
		index++
	}
	if err := tenantRows.Err(); err != nil {
		return fmt.Errorf("iterate tenants for operation log seed: %w", err)
	}
	return nil
}

func (repository *DesktopDataRepository) backfillUsageLogTenants(ctx context.Context) error {
	adminTenantRows, err := repository.pool.pool.Query(ctx, `SELECT admin_name, code FROM app_tenants WHERE admin_name <> '' AND code <> ''`)
	if err != nil {
		return fmt.Errorf("query tenant admins for usage backfill: %w", err)
	}
	defer adminTenantRows.Close()

	adminToCode := make(map[string]string)
	for adminTenantRows.Next() {
		var adminName string
		var tenantCode string
		if err := adminTenantRows.Scan(&adminName, &tenantCode); err != nil {
			return fmt.Errorf("scan tenant admin for usage backfill: %w", err)
		}
		adminToCode[strings.TrimSpace(adminName)] = strings.TrimSpace(tenantCode)
	}
	if err := adminTenantRows.Err(); err != nil {
		return fmt.Errorf("iterate tenant admins for usage backfill: %w", err)
	}

	for adminName, tenantCode := range adminToCode {
		if _, err := repository.pool.pool.Exec(ctx, `UPDATE app_login_logs SET tenant_code = $1 WHERE tenant_code = '' AND name = $2`, tenantCode, adminName); err != nil {
			return fmt.Errorf("backfill login log tenant code: %w", err)
		}
		if _, err := repository.pool.pool.Exec(ctx, `UPDATE app_operation_logs SET tenant_code = $1 WHERE tenant_code = '' AND name = $2`, tenantCode, adminName); err != nil {
			return fmt.Errorf("backfill operation log tenant code: %w", err)
		}
	}

	return nil
}

func (repository *DesktopDataRepository) seedDicts(ctx context.Context) error {
	var count int64
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_dicts`).Scan(&count); err != nil {
		return fmt.Errorf("count dicts: %w", err)
	}
	if count > 0 {
		return nil
	}
	seedRows := []service.DesktopDictRecord{{DictID: "001", Name: "岗位", Description: "岗位字典", ModifiedAt: time.Now().Add(-2 * time.Hour), ModifiedBy: "张三", CreatedAt: time.Now().Add(-24 * time.Hour), CreatedBy: "张三", Items: []service.DesktopDictItemRecord{{Index: 1, Label: "前端工程师", KeyVal: "FE", StyleType: "主要", ModifiedAt: time.Now().Add(-2 * time.Hour), ModifiedBy: "张三", CreatedAt: time.Now().Add(-24 * time.Hour), CreatedBy: "张三", BadgeColor: "bg-orange-50 text-orange-500"}, {Index: 2, Label: "后端工程师", KeyVal: "BE", StyleType: "次要", ModifiedAt: time.Now().Add(-90 * time.Minute), ModifiedBy: "李四", CreatedAt: time.Now().Add(-20 * time.Hour), CreatedBy: "李四", BadgeColor: "bg-emerald-50 text-emerald-500"}}}, {DictID: "002", Name: "职称", Description: "职称字典", ModifiedAt: time.Now().Add(-3 * time.Hour), ModifiedBy: "张三", CreatedAt: time.Now().Add(-48 * time.Hour), CreatedBy: "张三", Items: []service.DesktopDictItemRecord{{Index: 1, Label: "高级职称", KeyVal: "1", StyleType: "主要", ModifiedAt: time.Now().Add(-3 * time.Hour), ModifiedBy: "张三", CreatedAt: time.Now().Add(-48 * time.Hour), CreatedBy: "张三", BadgeColor: "bg-orange-50 text-orange-500"}, {Index: 2, Label: "中级职称", KeyVal: "2", StyleType: "次要", ModifiedAt: time.Now().Add(-2 * time.Hour), ModifiedBy: "李四", CreatedAt: time.Now().Add(-30 * time.Hour), CreatedBy: "李四", BadgeColor: "bg-emerald-50 text-emerald-500"}, {Index: 3, Label: "低级职称", KeyVal: "3", StyleType: "次要", ModifiedAt: time.Now().Add(-1 * time.Hour), ModifiedBy: "王五", CreatedAt: time.Now().Add(-26 * time.Hour), CreatedBy: "王五", BadgeColor: "bg-slate-100 dark:bg-slate-800 text-slate-500 dark:text-slate-400"}}}}
	for _, item := range seedRows {
		var dictID int64
		if err := repository.pool.pool.QueryRow(ctx, `INSERT INTO app_dicts (dict_id, name, description, modified_at, modified_by, created_at, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`, item.DictID, item.Name, item.Description, item.ModifiedAt, item.ModifiedBy, item.CreatedAt, item.CreatedBy).Scan(&dictID); err != nil {
			return fmt.Errorf("seed dicts: %w", err)
		}
		for _, child := range item.Items {
			if _, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_dict_items (dict_ref_id, item_index, label, key_val, style_type, modified_at, modified_by, created_at, created_by, badge_color) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, dictID, child.Index, child.Label, child.KeyVal, child.StyleType, child.ModifiedAt, child.ModifiedBy, child.CreatedAt, child.CreatedBy, child.BadgeColor); err != nil {
				return fmt.Errorf("seed dict items: %w", err)
			}
		}
	}
	return nil
}

func (repository *DesktopDataRepository) seedSystemSettings(ctx context.Context) error {
	var count int64
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_system_settings`).Scan(&count); err != nil {
		return fmt.Errorf("count system settings: %w", err)
	}
	if count > 0 {
		return nil
	}
	_, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_system_settings (id, website_title, system_logo, theme, icp, copyright, require_strong_password, login_fail_limit, login_lock_minutes) VALUES (1,$1,$2,$3,$4,$5,$6,$7,$8)`,
		"多租户后台管理系统", "", "emerald", "京ICP备10000000号-1", "Copyright © 2023 MTBM SYSTEM. All Rights Reserved.", true, 5, 30)
	if err != nil {
		return fmt.Errorf("seed system settings: %w", err)
	}
	return nil
}

func (repository *DesktopDataRepository) seedUsersAndRoles(ctx context.Context) error {
	var count int64
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_users`).Scan(&count); err != nil {
		return fmt.Errorf("count users: %w", err)
	}
	if count == 0 {
		departments := []service.DesktopDepartmentRecord{{Name: "人力资源部", Count: 12}, {Name: "财务部", Count: 23}, {Name: "生产部", Count: 32}, {Name: "营销部", Count: 23}}
		for _, item := range departments {
			if _, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_departments (name, user_count) VALUES ($1,$2)`, item.Name, item.Count); err != nil {
				return fmt.Errorf("seed departments: %w", err)
			}
		}
		users := []service.DesktopUserRecord{{UID: "1001", Name: "冯政", Department: "生产部", Phone: "13386911277", Role: "超管", Status: "开启", UpdatedAt: time.Now().Add(-2 * time.Hour), UpdatedBy: "张三"}, {UID: "1002", Name: "孙子面", Department: "生产部", Phone: "15366978328", Role: "管理员", Status: "关闭", UpdatedAt: time.Now().Add(-90 * time.Minute), UpdatedBy: "张三"}, {UID: "1003", Name: "钱继初", Department: "财务部", Phone: "18555418491", Role: "财务", Status: "关闭", UpdatedAt: time.Now().Add(-70 * time.Minute), UpdatedBy: "张三"}, {UID: "1004", Name: "李建华", Department: "营销部", Phone: "18716223293", Role: "运营", Status: "开启", UpdatedAt: time.Now().Add(-30 * time.Minute), UpdatedBy: "张三"}}
		for _, item := range users {
			if _, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_users (uid, name, department, phone, role_name, status, updated_at, updated_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, item.UID, item.Name, item.Department, item.Phone, item.Role, item.Status, item.UpdatedAt, item.UpdatedBy); err != nil {
				return fmt.Errorf("seed users: %w", err)
			}
		}
	}
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_roles`).Scan(&count); err != nil {
		return fmt.Errorf("count roles: %w", err)
	}
	if count == 0 {
		roles := []service.DesktopRoleRecord{{Name: "超级管理员", Key: "admin", Order: 1, Status: true, CreatedAt: time.Now().Add(-24 * time.Hour)}, {Name: "高管", Key: "manager", Order: 2, Status: true, CreatedAt: time.Now().Add(-20 * time.Hour)}, {Name: "部门负责人", Key: "dept_leader", Order: 3, Status: true, CreatedAt: time.Now().Add(-18 * time.Hour)}, {Name: "普通员工", Key: "common", Order: 4, Status: true, CreatedAt: time.Now().Add(-12 * time.Hour)}, {Name: "闲置角色", Key: "idle", Order: 7, Status: false, CreatedAt: time.Now().Add(-6 * time.Hour)}}
		for _, item := range roles {
			if _, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_roles (name, role_key, display_order, status, created_at) VALUES ($1,$2,$3,$4,$5)`, item.Name, item.Key, item.Order, item.Status, item.CreatedAt); err != nil {
				return fmt.Errorf("seed roles: %w", err)
			}
		}
	}
	if _, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_role_data_permissions (role_id, scope, departments_json) VALUES (1,'all','[]') ON CONFLICT (role_id) DO NOTHING`); err != nil {
		return fmt.Errorf("seed role data permission: %w", err)
	}
	return nil
}

func (repository *DesktopDataRepository) seedOrgsAndTenants(ctx context.Context) error {
	var count int64
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_orgs`).Scan(&count); err != nil {
		return fmt.Errorf("count orgs: %w", err)
	}
	if count == 0 {
		orgs := []service.DesktopOrgRecord{{Name: "生产一部", Code: "32", ParentName: "东方科技", Sort: 1, Category: "部门", CreatedAt: time.Now().Add(-24 * time.Hour), CreatedBy: "张三", UpdatedAt: time.Now().Add(-2 * time.Hour), UpdatedBy: "张三"}, {Name: "生产二部", Code: "33", ParentName: "东方科技", Sort: 2, Category: "部门", CreatedAt: time.Now().Add(-23 * time.Hour), CreatedBy: "张三", UpdatedAt: time.Now().Add(-90 * time.Minute), UpdatedBy: "张三"}, {Name: "营销中心", Code: "88", ParentName: "东方科技", Sort: 3, Category: "部门", CreatedAt: time.Now().Add(-22 * time.Hour), CreatedBy: "李四", UpdatedAt: time.Now().Add(-70 * time.Minute), UpdatedBy: "李四"}}
		departments := []string{"生产部", "生产部", "营销部"}
		for index, item := range orgs {
			if _, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_orgs (name, code, parent_name, sort_order, category, created_at, created_by, updated_at, updated_by, department) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, item.Name, item.Code, item.ParentName, item.Sort, item.Category, item.CreatedAt, item.CreatedBy, item.UpdatedAt, item.UpdatedBy, departments[index]); err != nil {
				return fmt.Errorf("seed orgs: %w", err)
			}
		}
	}
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_tenants`).Scan(&count); err != nil {
		return fmt.Errorf("count tenants: %w", err)
	}
	if count == 0 {
		tenants := []service.DesktopTenantRecord{{Code: "FQJT", Name: "智迪互动(北京)广告有限公司", Status: "开启", Period: "2023-01-01 - 2023-12-31", Admin: "张三", Phone: "131 1242 1320"}, {Code: "HIMGJT", Name: "北京长友物业管理有限公司", Status: "关闭", Period: "2023-01-01 - 2023-12-31", Admin: "张三", Phone: "131 1242 1320"}, {Code: "HLGJT", Name: "陕西沙龙传媒有限公司", Status: "关闭", Period: "2023-01-01 - 2023-12-31", Admin: "李四", Phone: "131 1242 1320"}}
		for _, item := range tenants {
			if _, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_tenants (code, name, status, period, admin_name, admin_phone) VALUES ($1,$2,$3,$4,$5,$6)`, item.Code, item.Name, item.Status, item.Period, item.Admin, item.Phone); err != nil {
				return fmt.Errorf("seed tenants: %w", err)
			}
		}
	}
	return nil
}

func (repository *DesktopDataRepository) seedMenuTemplates(ctx context.Context) error {
	var count int64
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_menu_templates`).Scan(&count); err != nil {
		return fmt.Errorf("count menu templates: %w", err)
	}
	if count > 0 {
		return nil
	}
	rows := []service.DesktopMenuTemplateRecord{{Name: "试用版", TenantCount: 32, Description: "试用版", UpdatedAt: time.Now().Add(-2 * time.Hour), UpdatedBy: "张三", CreatedAt: time.Now().Add(-24 * time.Hour), CreatedBy: "张三"}, {Name: "基础版", TenantCount: 32, Description: "基础版", UpdatedAt: time.Now().Add(-90 * time.Minute), UpdatedBy: "张三", CreatedAt: time.Now().Add(-23 * time.Hour), CreatedBy: "张三"}, {Name: "VIP版", TenantCount: 32, Description: "VIP版", UpdatedAt: time.Now().Add(-30 * time.Minute), UpdatedBy: "杜三", CreatedAt: time.Now().Add(-22 * time.Hour), CreatedBy: "张三"}}
	for _, item := range rows {
		if _, err := repository.pool.pool.Exec(ctx, `INSERT INTO app_menu_templates (name, tenant_count, description, updated_at, updated_by, created_at, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7)`, item.Name, item.TenantCount, item.Description, item.UpdatedAt, item.UpdatedBy, item.CreatedAt, item.CreatedBy); err != nil {
			return fmt.Errorf("seed menu templates: %w", err)
		}
	}
	return nil
}

func (repository *DesktopDataRepository) seedDesktopMenus(ctx context.Context) error {
	var count int64
	if err := repository.pool.pool.QueryRow(ctx, `SELECT COUNT(*) FROM app_desktop_menus`).Scan(&count); err != nil {
		return fmt.Errorf("count desktop menus: %w", err)
	}
	if count > 0 {
		return nil
	}
	rows := []service.DesktopMenuRecord{{ParentID: 0, Name: "首页", Level: 1, Sort: 1, Type: "目录", Icon: "Home", Status: "可见", Path: "/", Permission: "desktop:home:view"}, {ParentID: 0, Name: "租户配置", Level: 1, Sort: 2, Type: "目录", Icon: "Settings", Status: "不可见", Path: "/tenant", Permission: "desktop:tenant:view"}, {ParentID: 2, Name: "租户管理", Level: 2, Sort: 1, Type: "菜单", Icon: "Settings", Status: "不可见", Path: "/tenant/management", Permission: "desktop:tenant:list"}, {ParentID: 2, Name: "菜单模板", Level: 2, Sort: 2, Type: "菜单", Icon: "Menu", Status: "不可见", Path: "/tenant/menu-template", Permission: "desktop:tenant:template"}, {ParentID: 2, Name: "组织管理", Level: 2, Sort: 3, Type: "菜单", Icon: "Layers", Status: "不可见", Path: "/orgs", Permission: "desktop:org:list"}}
	var createdIDs []int64
	for _, item := range rows {
		record, err := repository.CreateDesktopMenu(ctx, item)
		if err != nil {
			return err
		}
		createdIDs = append(createdIDs, record.ID)
	}
	if len(createdIDs) >= 5 {
		buttons := []service.DesktopMenuRecord{{ParentID: createdIDs[4], Name: "查询/查看", Level: 3, Sort: 1, Type: "按钮", Icon: "Layers", Status: "不可见", Path: "/orgs", Permission: "desktop:org:view"}, {ParentID: createdIDs[4], Name: "新增", Level: 3, Sort: 2, Type: "按钮", Icon: "Layers", Status: "可见", Path: "/orgs", Permission: "desktop:org:create"}, {ParentID: createdIDs[4], Name: "修改", Level: 3, Sort: 3, Type: "按钮", Icon: "Layers", Status: "可见", Path: "/orgs", Permission: "desktop:org:update"}, {ParentID: createdIDs[4], Name: "删除", Level: 3, Sort: 4, Type: "按钮", Icon: "Layers", Status: "可见", Path: "/orgs", Permission: "desktop:org:delete"}}
		for _, item := range buttons {
			if _, err := repository.CreateDesktopMenu(ctx, item); err != nil {
				return err
			}
		}
	}
	return nil
}

var _ = pgx.ErrNoRows
