package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pct/pgpackage/internal/apply"
	"github.com/pct/pgpackage/internal/diff"
	"github.com/pct/pgpackage/internal/introspect"
	"github.com/pct/pgpackage/internal/packagefmt"
	"github.com/pct/pgpackage/internal/parser"
	"github.com/pct/pgpackage/internal/projectxml"
	"github.com/pct/pgpackage/internal/render"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		printUsage()
		return errors.New("missing command")
	}

	switch args[0] {
	case "build":
		return runBuild(ctx, args[1:])
	case "plan":
		return runPlan(ctx, args[1:])
	case "apply":
		return runApply(ctx, args[1:])
	case "help", "--help", "-h":
		printUsage()
		return nil
	default:
		printUsage()
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runBuild(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	projectPath := fs.String("project", "", "Path to a .pgpackage project file")
	outputPath := fs.String("output", "", "Output file or directory")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *projectPath == "" || *outputPath == "" {
		return errors.New("build requires --project and --output")
	}

	project, rawXML, err := projectxml.Load(*projectPath)
	if err != nil {
		return err
	}

	desired, err := parser.BuildDesiredModel(project)
	if err != nil {
		return err
	}

	manifest := packagefmt.NewManifest(project, desired)
	targetOutput, err := resolvePackageOutput(*outputPath, project.PackageID)
	if err != nil {
		return err
	}

	if err := packagefmt.Write(targetOutput, manifest, desired, rawXML, project); err != nil {
		return err
	}

	fmt.Println(targetOutput)
	_ = ctx
	return nil
}

func runPlan(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("plan", flag.ContinueOnError)
	packagePath := fs.String("package", "", "Path to a .pgpkg package")
	connectionString := fs.String("connection", "", "PostgreSQL connection string")
	format := fs.String("format", "text", "Output format: text or json")
	scriptPath := fs.String("script", "", "Optional path to write the SQL preview")
	allowDrop := fs.Bool("allow-drop", false, "Allow destructive operations to be rendered as executable SQL")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *packagePath == "" || *connectionString == "" {
		return errors.New("plan requires --package and --connection")
	}

	pkg, err := packagefmt.Read(*packagePath)
	if err != nil {
		return err
	}

	actual, err := introspect.LoadActualModel(
		ctx,
		*connectionString,
		pkg.Project.Target.OwnedSchemaNames(),
		pkg.Project.Target.ExtensionNames(),
		pkg.Project.PostgresVersion,
	)
	if err != nil {
		return err
	}

	plan := diff.BuildPlan(pkg.Project, pkg.Model, actual, diff.Options{AllowDrop: *allowDrop})
	if *scriptPath != "" {
		if err := os.WriteFile(*scriptPath, []byte(render.SQL(plan)), 0o644); err != nil {
			return err
		}
	}

	switch strings.ToLower(*format) {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(plan)
	case "text":
		fmt.Print(render.Text(plan))
		return nil
	default:
		return fmt.Errorf("unsupported format %q", *format)
	}
}

func runApply(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("apply", flag.ContinueOnError)
	packagePath := fs.String("package", "", "Path to a .pgpkg package")
	connectionString := fs.String("connection", "", "PostgreSQL connection string")
	allowDrop := fs.Bool("allow-drop", false, "Allow destructive operations")
	force := fs.Bool("force", false, "Bypass destructive operation protection")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *packagePath == "" || *connectionString == "" {
		return errors.New("apply requires --package and --connection")
	}

	pkg, err := packagefmt.Read(*packagePath)
	if err != nil {
		return err
	}

	actual, err := introspect.LoadActualModel(
		ctx,
		*connectionString,
		pkg.Project.Target.OwnedSchemaNames(),
		pkg.Project.Target.ExtensionNames(),
		pkg.Project.PostgresVersion,
	)
	if err != nil {
		return err
	}

	plan := diff.BuildPlan(pkg.Project, pkg.Model, actual, diff.Options{AllowDrop: *allowDrop || *force})
	if err := apply.Execute(ctx, *connectionString, pkg.Project, plan, apply.Options{AllowDrop: *allowDrop, Force: *force}); err != nil {
		return err
	}

	fmt.Println("Applied package.")
	return nil
}

func resolvePackageOutput(outputPath, packageID string) (string, error) {
	outputPath = filepath.Clean(outputPath)
	if strings.HasSuffix(outputPath, ".pgpkg") {
		if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
			return "", err
		}
		return outputPath, nil
	}

	if err := os.MkdirAll(outputPath, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(outputPath, packageID+".pgpkg"), nil
}

func printUsage() {
	fmt.Println("pgpackage build --project <file.pgpackage> --output <dir-or-file>")
	fmt.Println("pgpackage plan --package <file.pgpkg> --connection <postgres-uri> [--format text|json] [--script <file>] [--allow-drop]")
	fmt.Println("pgpackage apply --package <file.pgpkg> --connection <postgres-uri> [--allow-drop] [--force]")
}
