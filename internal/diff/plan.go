package diff

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pct/pgpackage/internal/model"
	"github.com/pct/pgpackage/internal/projectxml"
)

type Options struct {
	AllowDrop bool
}

type Plan struct {
	Summary    Summary     `json:"summary"`
	Operations []Operation `json:"operations"`
}

type Summary struct {
	Supported      bool `json:"supported"`
	Destructive    bool `json:"destructive"`
	OperationCount int  `json:"operationCount"`
}

type Operation struct {
	Kind       string `json:"kind"`
	ObjectType string `json:"objectType"`
	ObjectKey  string `json:"objectKey"`
	Risk       string `json:"risk"`
	SQL        string `json:"sql"`
}

func BuildPlan(project *projectxml.Project, desired, actual *model.SchemaModel, options Options) Plan {
	var ops []Operation

	diffByName(
		project,
		desired.Schemas,
		actual.Schemas,
		func(item model.SchemaDef) string { return item.Name },
		func(item model.SchemaDef) Operation {
			return op("create-schema", "schema", item.Name, "safe", ensureSemicolon(item.SQL))
		},
		func(item model.SchemaDef) Operation {
			return op("drop-schema", "schema", item.Name, "destructive", fmt.Sprintf("DROP SCHEMA %s CASCADE;", quoteQName(item.Name)))
		},
		func(a, b model.SchemaDef) bool { return true },
		&ops,
		options,
	)

	diffByName(
		project,
		desired.Extensions,
		actual.Extensions,
		func(item model.ExtensionDef) string { return item.Name },
		func(item model.ExtensionDef) Operation {
			return op("create-extension", "extension", item.Name, "safe", ensureSemicolon(item.SQL))
		},
		func(item model.ExtensionDef) Operation {
			return op("drop-extension", "extension", item.Name, "destructive", fmt.Sprintf("DROP EXTENSION %s;", quoteQName(item.Name)))
		},
		func(a, b model.ExtensionDef) bool { return a.Version == b.Version && a.Name == b.Name },
		&ops,
		options,
	)

	diffByName(
		project,
		desired.Tables,
		actual.Tables,
		func(item model.TableDef) string { return model.QualifiedName(item.Schema, item.Name) },
		func(item model.TableDef) Operation {
			return op("create-table", "table", model.QualifiedName(item.Schema, item.Name), "safe", ensureSemicolon(item.SQL))
		},
		func(item model.TableDef) Operation {
			return op("drop-table", "table", model.QualifiedName(item.Schema, item.Name), "destructive", fmt.Sprintf("DROP TABLE %s CASCADE;", quoteQName(item.Schema, item.Name)))
		},
		func(a, b model.TableDef) bool { return a.SQL == b.SQL },
		&ops,
		options,
	)

	diffByName(
		project,
		desired.Indexes,
		actual.Indexes,
		func(item model.IndexDef) string { return model.QualifiedName(item.Schema, item.Name) },
		func(item model.IndexDef) Operation {
			return op("create-index", "index", model.QualifiedName(item.Schema, item.Name), "safe", ensureSemicolon(item.SQL))
		},
		func(item model.IndexDef) Operation {
			return op("drop-index", "index", model.QualifiedName(item.Schema, item.Name), "destructive", fmt.Sprintf("DROP INDEX %s;", quoteQName(item.Schema, item.Name)))
		},
		func(a, b model.IndexDef) bool { return a.SQL == b.SQL },
		&ops,
		options,
	)

	diffByName(
		project,
		desired.Views,
		actual.Views,
		func(item model.ViewDef) string { return model.QualifiedName(item.Schema, item.Name) },
		func(item model.ViewDef) Operation {
			return op("create-view", "view", model.QualifiedName(item.Schema, item.Name), "safe", ensureSemicolon(item.SQL))
		},
		func(item model.ViewDef) Operation {
			return op("drop-view", "view", model.QualifiedName(item.Schema, item.Name), "destructive", fmt.Sprintf("DROP VIEW %s CASCADE;", quoteQName(item.Schema, item.Name)))
		},
		func(a, b model.ViewDef) bool { return a.SQL == b.SQL },
		&ops,
		options,
	)

	diffByName(
		project,
		desired.Routines,
		actual.Routines,
		func(item model.RoutineDef) string { return model.RoutineKey(item.Schema, item.Name, item.IdentityArgs) },
		func(item model.RoutineDef) Operation {
			return op("create-routine", item.Kind, model.RoutineKey(item.Schema, item.Name, item.IdentityArgs), "safe", ensureSemicolon(item.SQL))
		},
		func(item model.RoutineDef) Operation {
			return op("drop-routine", item.Kind, model.RoutineKey(item.Schema, item.Name, item.IdentityArgs), "destructive", fmt.Sprintf("DROP %s %s(%s);", strings.ToUpper(item.Kind), quoteQName(item.Schema, item.Name), item.IdentityArgs))
		},
		func(a, b model.RoutineDef) bool { return a.SQL == b.SQL },
		&ops,
		options,
	)

	diffByName(
		project,
		desired.Enums,
		actual.Enums,
		func(item model.EnumDef) string { return model.QualifiedName(item.Schema, item.Name) },
		func(item model.EnumDef) Operation {
			return op("create-enum", "enum", model.QualifiedName(item.Schema, item.Name), "safe", ensureSemicolon(item.SQL))
		},
		func(item model.EnumDef) Operation {
			return op("drop-enum", "enum", model.QualifiedName(item.Schema, item.Name), "destructive", fmt.Sprintf("DROP TYPE %s;", quoteQName(item.Schema, item.Name)))
		},
		func(a, b model.EnumDef) bool { return a.SQL == b.SQL },
		&ops,
		options,
	)

	diffByName(
		project,
		desired.Domains,
		actual.Domains,
		func(item model.DomainDef) string { return model.QualifiedName(item.Schema, item.Name) },
		func(item model.DomainDef) Operation {
			return op("create-domain", "domain", model.QualifiedName(item.Schema, item.Name), "safe", ensureSemicolon(item.SQL))
		},
		func(item model.DomainDef) Operation {
			return op("drop-domain", "domain", model.QualifiedName(item.Schema, item.Name), "destructive", fmt.Sprintf("DROP DOMAIN %s;", quoteQName(item.Schema, item.Name)))
		},
		func(a, b model.DomainDef) bool { return a.SQL == b.SQL },
		&ops,
		options,
	)

	diffByName(
		project,
		desired.Sequences,
		actual.Sequences,
		func(item model.SequenceDef) string { return model.QualifiedName(item.Schema, item.Name) },
		func(item model.SequenceDef) Operation {
			return op("create-sequence", "sequence", model.QualifiedName(item.Schema, item.Name), "safe", ensureSemicolon(item.SQL))
		},
		func(item model.SequenceDef) Operation {
			return op("drop-sequence", "sequence", model.QualifiedName(item.Schema, item.Name), "destructive", fmt.Sprintf("DROP SEQUENCE %s;", quoteQName(item.Schema, item.Name)))
		},
		func(a, b model.SequenceDef) bool { return a.SQL == b.SQL },
		&ops,
		options,
	)

	diffByName(
		project,
		desired.Comments,
		actual.Comments,
		func(item model.CommentDef) string { return item.ObjectType + ":" + item.ObjectKey },
		func(item model.CommentDef) Operation {
			return op("set-comment", item.ObjectType, item.ObjectKey, "safe", ensureSemicolon(item.SQL))
		},
		func(item model.CommentDef) Operation {
			return op("clear-comment", item.ObjectType, item.ObjectKey, "safe", clearCommentSQL(item))
		},
		func(a, b model.CommentDef) bool { return a.Comment == b.Comment },
		&ops,
		options,
	)

	sort.Slice(ops, func(i, j int) bool {
		if weight(ops[i].Kind) == weight(ops[j].Kind) {
			return ops[i].ObjectKey < ops[j].ObjectKey
		}
		return weight(ops[i].Kind) < weight(ops[j].Kind)
	})

	destructive := false
	for _, operation := range ops {
		if operation.Risk == "destructive" {
			destructive = true
			break
		}
	}

	return Plan{
		Summary: Summary{
			Supported:      true,
			Destructive:    destructive,
			OperationCount: len(ops),
		},
		Operations: ops,
	}
}

func diffByName[T any](project *projectxml.Project, desired, actual []T, key func(T) string, createOp func(T) Operation, dropOp func(T) Operation, equal func(T, T) bool, ops *[]Operation, options Options) {
	desiredMap := make(map[string]T, len(desired))
	for _, item := range desired {
		desiredMap[key(item)] = item
	}
	actualMap := make(map[string]T, len(actual))
	for _, item := range actual {
		actualMap[key(item)] = item
	}

	for name, desiredItem := range desiredMap {
		actualItem, exists := actualMap[name]
		if !exists {
			*ops = append(*ops, createOp(desiredItem))
			continue
		}
		if !equal(desiredItem, actualItem) {
			*ops = append(*ops, recreateOps(createOp(desiredItem), dropOp(actualItem), project, options)...)
		}
	}

	for name, actualItem := range actualMap {
		if _, exists := desiredMap[name]; exists {
			continue
		}
		if project.Target.Plan.AllowDrop || options.AllowDrop {
			*ops = append(*ops, dropOp(actualItem))
		} else {
			operation := dropOp(actualItem)
			operation.Kind = "blocked-" + operation.Kind
			operation.SQL = "-- " + operation.SQL
			*ops = append(*ops, operation)
		}
	}
}

func recreateOps(create Operation, drop Operation, project *projectxml.Project, options Options) []Operation {
	if create.ObjectType == "view" {
		create.Kind = "replace-view"
		create.SQL = ensureOrReplace(create.SQL, "CREATE VIEW", "CREATE OR REPLACE VIEW")
		return []Operation{create}
	}
	if create.ObjectType == "function" || create.ObjectType == "procedure" {
		needle := "CREATE " + strings.ToUpper(create.ObjectType)
		replacement := "CREATE OR REPLACE " + strings.ToUpper(create.ObjectType)
		create.Kind = "replace-routine"
		create.SQL = ensureOrReplace(create.SQL, needle, replacement)
		return []Operation{create}
	}

	if project.Target.Plan.AllowDrop || options.AllowDrop {
		create.Risk = "destructive"
		return []Operation{drop, create}
	}
	drop.Kind = "blocked-" + drop.Kind
	drop.SQL = "-- " + drop.SQL
	create.Kind = "blocked-recreate-" + create.Kind
	create.Risk = "destructive"
	create.SQL = "-- requires destructive recreate\n-- " + create.SQL
	return []Operation{drop, create}
}

func ensureSemicolon(sql string) string {
	sql = strings.TrimSpace(sql)
	if strings.HasSuffix(sql, ";") {
		return sql
	}
	return sql + ";"
}

func ensureOrReplace(sql, needle, replacement string) string {
	sql = ensureSemicolon(sql)
	up := strings.ToUpper(sql)
	if strings.HasPrefix(up, replacement) {
		return sql
	}
	if strings.HasPrefix(up, needle) {
		return replacement + sql[len(needle):]
	}
	return sql
}

func clearCommentSQL(comment model.CommentDef) string {
	switch comment.ObjectType {
	case "table":
		return fmt.Sprintf("COMMENT ON TABLE %s IS NULL;", quoteQName(strings.Split(comment.ObjectKey, ".")...))
	case "column":
		parts := strings.Split(comment.ObjectKey, ".")
		if len(parts) == 3 {
			return fmt.Sprintf("COMMENT ON COLUMN %s.%s.%s IS NULL;", quoteQName(parts[0]), quoteQName(parts[1]), quoteQName(parts[2]))
		}
	}
	return "-- unsupported comment clear;"
}

func op(kind, objectType, objectKey, risk, sql string) Operation {
	return Operation{Kind: kind, ObjectType: objectType, ObjectKey: objectKey, Risk: risk, SQL: sql}
}

func quoteQName(parts ...string) string {
	quoted := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		quoted = append(quoted, `"`+strings.ReplaceAll(part, `"`, `""`)+`"`)
	}
	return strings.Join(quoted, ".")
}

func weight(kind string) int {
	switch {
	case strings.Contains(kind, "create-schema"):
		return 10
	case strings.Contains(kind, "create-extension"):
		return 20
	case strings.Contains(kind, "create-enum"), strings.Contains(kind, "create-domain"), strings.Contains(kind, "create-sequence"):
		return 30
	case strings.Contains(kind, "drop-table"):
		return 40
	case strings.Contains(kind, "create-table"):
		return 50
	case strings.Contains(kind, "replace-view"), strings.Contains(kind, "create-view"):
		return 60
	case strings.Contains(kind, "replace-routine"), strings.Contains(kind, "create-routine"):
		return 70
	case strings.Contains(kind, "create-index"):
		return 80
	case strings.Contains(kind, "set-comment"), strings.Contains(kind, "clear-comment"):
		return 90
	default:
		return 100
	}
}
