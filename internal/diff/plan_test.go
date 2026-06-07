package diff

import (
	"testing"

	"github.com/pct/pgpackage/internal/model"
	"github.com/pct/pgpackage/internal/projectxml"
)

func TestBuildPlanCreateAndDrop(t *testing.T) {
	project := &projectxml.Project{
		Target: projectxml.TargetConfig{
			Plan: projectxml.PlanConfig{AllowDrop: false},
		},
	}
	desired := &model.SchemaModel{
		Tables: []model.TableDef{{Schema: "app", Name: "widgets", SQL: "CREATE TABLE app.widgets (id uuid)"}},
	}
	actual := &model.SchemaModel{
		Tables: []model.TableDef{{Schema: "app", Name: "legacy", SQL: "CREATE TABLE app.legacy (id uuid)"}},
	}
	plan := BuildPlan(project, desired, actual, Options{})
	if len(plan.Operations) != 2 {
		t.Fatalf("expected 2 operations, got %d", len(plan.Operations))
	}
	if plan.Operations[0].Kind != "blocked-drop-table" && plan.Operations[1].Kind != "blocked-drop-table" {
		t.Fatalf("expected blocked drop operation, got %#v", plan.Operations)
	}
}
