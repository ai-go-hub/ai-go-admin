package driver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalFullPathAcceptsPathInsideRoot(t *testing.T) {
	root := t.TempDir()
	local := NewLocal(root, "/static/upload")

	got, err := local.FullPath("/static/upload/2026/09/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, filepath.FromSlash("2026/09/test.txt"))
	if got != want {
		t.Errorf("FullPath() = %q, want %q", got, want)
	}
}

func TestLocalSaveCreatesFileInsideRoot(t *testing.T) {
	root := t.TempDir()
	local := NewLocal(root, "/static/upload")

	if err := local.Save(strings.NewReader("content"), "/static/upload/2026/09/test.txt"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash("2026/09/test.txt")))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "content" {
		t.Fatalf("saved content = %q, want content", data)
	}
}
func TestLocalFullPathRejectsPathOutsideRoot(t *testing.T) {
	root := t.TempDir()
	local := NewLocal(root, "/static/upload")

	invalid := []string{
		"",
		".",
	}

	for _, path := range invalid {
		if got, err := local.FullPath(path); err == nil {
			t.Errorf("FullPath(%q) = %q, want error", path, got)
		}
	}
}

func TestLocalSaveRejectsTraversalPath(t *testing.T) {
	root := t.TempDir()
	local := NewLocal(root, "/static/upload")

	traversal := []string{
		"../outside.txt",
		"../../outside.txt",
		"/static/upload/../outside.txt",
		"/static/upload/../../outside.txt",
	}

	for _, path := range traversal {
		if err := local.Save(strings.NewReader("content"), path); err == nil {
			t.Errorf("Save(%q) succeeded, want error", path)
		}
	}
}

func TestLocalDeleteRejectsPathOutsideRoot(t *testing.T) {
	root := t.TempDir()
	parent := filepath.Dir(root)
	outside := filepath.Join(parent, "local-delete-outside.txt")
	if err := os.WriteFile(outside, []byte("outside"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Remove(outside)
	})

	local := NewLocal(root, "/static/upload")
	if err := local.Delete("/static/upload/../" + filepath.Base(outside)); err == nil {
		t.Fatal("Delete outside root succeeded, want error")
	}

	data, err := os.ReadFile(outside)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "outside" {
		t.Fatalf("outside file content = %q, want outside", data)
	}
}

func TestLocalDeleteRejectsSymlink(t *testing.T) {
	root := t.TempDir()
	parent := filepath.Dir(root)
	outside := filepath.Join(parent, "local-delete-symlink-outside.txt")
	if err := os.WriteFile(outside, []byte("outside"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Remove(outside)
	})

	link := filepath.Join(root, "link.txt")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("current environment cannot create symlink: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(link)
	})

	local := NewLocal(root, "/static/upload")
	if err := local.Delete("/static/upload/link.txt"); err == nil {
		t.Fatal("Delete symlink succeeded, want error")
	}

	data, err := os.ReadFile(outside)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "outside" {
		t.Fatalf("outside file content = %q, want outside", data)
	}
}
