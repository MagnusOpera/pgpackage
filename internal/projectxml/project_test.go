package projectxml

import "testing"

func TestLoadAndResolve(t *testing.T) {
	project, _, err := Load("../../testdata/sample/sample.pgpackage")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	files, err := project.ResolveFiles()
	if err != nil {
		t.Fatalf("ResolveFiles() error = %v", err)
	}
	if len(files) != 9 {
		t.Fatalf("expected 9 sql files, got %d", len(files))
	}
}
