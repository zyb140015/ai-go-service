package service

import (
	"testing"
	"time"
)

func TestFilterUsageLoginLogsByTenant(t *testing.T) {
	t.Parallel()

	items := []DesktopLoginLogRecord{{TenantCode: "A"}, {TenantCode: "B"}, {TenantCode: "A"}}
	filtered := filterUsageLoginLogsByTenant(items, "A")
	if len(filtered) != 2 {
		t.Fatalf("expected 2 items, got %d", len(filtered))
	}
}

func TestFilterUsageOperationLogsByTenant(t *testing.T) {
	t.Parallel()

	items := []DesktopOperationLogRecord{{TenantCode: "A"}, {TenantCode: "B"}}
	filtered := filterUsageOperationLogsByTenant(items, "B")
	if len(filtered) != 1 || filtered[0].TenantCode != "B" {
		t.Fatalf("unexpected filtered result: %#v", filtered)
	}
}

func TestBuildModuleRanksSortsDescending(t *testing.T) {
	t.Parallel()

	ranks := buildModuleRanks([]DesktopOperationLogRecord{{Module: "模块A"}, {Module: "模块B"}, {Module: "模块A"}}, true)
	if len(ranks) == 0 || ranks[0].Name != "模块A" || ranks[0].Count != 2 {
		t.Fatalf("unexpected ranks: %#v", ranks)
	}
}

func TestBuildUsageTrendBuildsRequestedRange(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.March, 25, 12, 0, 0, 0, time.UTC)
	trend := buildUsageTrend([]DesktopLoginLogRecord{{Time: now}}, []DesktopOperationLogRecord{{Time: now}}, "week", now)
	if len(trend) != 7 {
		t.Fatalf("expected 7 points, got %d", len(trend))
	}
	if trend[len(trend)-1].LoginCount != 1 || trend[len(trend)-1].UsageCount != 1 {
		t.Fatalf("unexpected last trend point: %#v", trend[len(trend)-1])
	}
}
