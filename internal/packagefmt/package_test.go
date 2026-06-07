package packagefmt

import (
	"path/filepath"
	"testing"

	"github.com/pct/pgpackage/internal/parser"
	"github.com/pct/pgpackage/internal/projectxml"
)

func TestWriteReadRoundTrip(t *testing.T) {
	project, raw, err := projectxml.Load("../../testdata/sample/sample.pgpackage")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	model, err := parser.BuildDesiredModel(project)
	if err != nil {
		t.Fatalf("BuildDesiredModel() error = %v", err)
	}
	path := filepath.Join(t.TempDir(), "sample.pgpkg")
	if err := Write(path, NewManifest(project, model), model, raw, project); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	pkg, err := Read(path)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if pkg.Manifest.PackageID != "SampleProject" {
		t.Fatalf("unexpected package id %q", pkg.Manifest.PackageID)
	}
	if len(pkg.Model.Tables) != 1 {
		t.Fatalf("expected 1 table, got %d", len(pkg.Model.Tables))
	}
}
