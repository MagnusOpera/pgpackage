package apply

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/pct/pgpackage/internal/diff"
	"github.com/pct/pgpackage/internal/projectxml"
)

type Options struct {
	AllowDrop bool
	Force     bool
}

func Execute(ctx context.Context, connectionString string, project *projectxml.Project, plan diff.Plan, options Options) error {
	if plan.Summary.Destructive && !(project.Target.Plan.AllowDrop || options.AllowDrop || options.Force) {
		return fmt.Errorf("plan contains destructive operations; re-run with --allow-drop or --force")
	}

	conn, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	runInTx := func(tx pgx.Tx) error {
		if project.Target.Apply.LockTimeout != "" {
			if _, err := tx.Exec(ctx, fmt.Sprintf("SET lock_timeout = '%s'", project.Target.Apply.LockTimeout)); err != nil {
				return err
			}
		}
		if project.Target.Apply.StatementTimeout != "" {
			if _, err := tx.Exec(ctx, fmt.Sprintf("SET statement_timeout = '%s'", project.Target.Apply.StatementTimeout)); err != nil {
				return err
			}
		}
		for _, operation := range plan.Operations {
			sql := strings.TrimSpace(operation.SQL)
			if sql == "" || strings.HasPrefix(sql, "--") {
				continue
			}
			if _, err := tx.Exec(ctx, sql); err != nil {
				return fmt.Errorf("%s %s: %w", operation.Kind, operation.ObjectKey, err)
			}
		}
		return nil
	}

	runNoTx := func() error {
		if project.Target.Apply.LockTimeout != "" {
			if _, err := conn.Exec(ctx, fmt.Sprintf("SET lock_timeout = '%s'", project.Target.Apply.LockTimeout)); err != nil {
				return err
			}
		}
		if project.Target.Apply.StatementTimeout != "" {
			if _, err := conn.Exec(ctx, fmt.Sprintf("SET statement_timeout = '%s'", project.Target.Apply.StatementTimeout)); err != nil {
				return err
			}
		}
		for _, operation := range plan.Operations {
			sql := strings.TrimSpace(operation.SQL)
			if sql == "" || strings.HasPrefix(sql, "--") {
				continue
			}
			if _, err := conn.Exec(ctx, sql); err != nil {
				return fmt.Errorf("%s %s: %w", operation.Kind, operation.ObjectKey, err)
			}
		}
		return nil
	}

	if project.Target.Apply.UseTransaction {
		tx, err := conn.Begin(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback(ctx)
		if err := runInTx(tx); err != nil {
			return err
		}
		return tx.Commit(ctx)
	}

	return runNoTx()
}
