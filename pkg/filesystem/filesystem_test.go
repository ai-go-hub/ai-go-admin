package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTrimExt(t *testing.T) {
	cases := []struct {
		name string
		path string
		want string
	}{
		{"普通文件", "a.txt", "a"},
		{"多级路径", "dir/deep/a.txt", "a"},
		{"多扩展名", "a.tar.gz", "a.tar"},
		{"无扩展名", "noext", "noext"},
		{"结尾点", "a.b.c", "a.b"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := TrimExt(c.path); got != c.want {
				t.Errorf("TrimExt(%q) = %q, want %q", c.path, got, c.want)
			}
		})
	}
}

func TestExists(t *testing.T) {
	t.Run("存在的文件", func(t *testing.T) {
		f := filepath.Join(t.TempDir(), "a.txt")
		if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		ok, err := Exists(f)
		if err != nil || !ok {
			t.Errorf("Exists(%q) = %v, %v, want true, nil", f, ok, err)
		}
	})
	t.Run("存在的目录", func(t *testing.T) {
		d := t.TempDir()
		ok, err := Exists(d)
		if err != nil || !ok {
			t.Errorf("Exists(%q) = %v, %v, want true, nil", d, ok, err)
		}
	})
	t.Run("不存在的路径", func(t *testing.T) {
		p := filepath.Join(t.TempDir(), "nope")
		ok, err := Exists(p)
		if err != nil || ok {
			t.Errorf("Exists(%q) = %v, %v, want false, nil", p, ok, err)
		}
	})
}

func TestExtension(t *testing.T) {
	cases := []struct {
		name     string
		filename string
		want     string
	}{
		{"普通", "a.txt", "txt"},
		{"大写转小写", "A.PNG", "png"},
		{"多扩展名取最后一段", "a.tar.gz", "gz"},
		{"无扩展名", "noext", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Extension(c.filename); got != c.want {
				t.Errorf("Extension(%q) = %q, want %q", c.filename, got, c.want)
			}
		})
	}
}

func TestIsImageExtension(t *testing.T) {
	for _, ext := range []string{"jpg", "jpeg", "png", "gif", "webp", "bmp", "svg", "avif"} {
		if !IsImageExtension(ext) {
			t.Errorf("IsImageExtension(%q) = false, want true", ext)
		}
	}
	if !IsImageExtension("JPG") {
		t.Error("IsImageExtension 应忽略大小写")
	}
	for _, ext := range []string{"txt", "pdf", "go", "md", ""} {
		if IsImageExtension(ext) {
			t.Errorf("IsImageExtension(%q) = true, want false", ext)
		}
	}
}

func TestFormatBytes(t *testing.T) {
	cases := []struct {
		name  string
		bytes int64
		want  string
	}{
		{"零", 0, "0 B"},
		{"1023 字节", 1023, "1023 B"},
		{"1 KB", 1024, "1 KB"},
		{"1.5 KB", 1536, "1.5 KB"},
		{"2 KB 去尾零", 2048, "2 KB"},
		{"1 MB", 1048576, "1 MB"},
		{"5.25 MB", 5505024, "5.25 MB"},
		{"1 GB", 1073741824, "1 GB"},
		{"1 TB", 1099511627776, "1 TB"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := FormatBytes(c.bytes); got != c.want {
				t.Errorf("FormatBytes(%d) = %q, want %q", c.bytes, got, c.want)
			}
		})
	}
}

func TestReadDir(t *testing.T) {
	t.Run("递归读取文件", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, "a", "b"), 0o755); err != nil {
			t.Fatal(err)
		}
		files := []string{
			filepath.Join(dir, "root.txt"),
			filepath.Join(dir, "a", "x.txt"),
			filepath.Join(dir, "a", "b", "y.txt"),
		}
		for _, f := range files {
			if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		paths, err := ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		if len(paths) != len(files) {
			t.Fatalf("ReadDir 返回 %d 个文件，want %d", len(paths), len(files))
		}
		set := make(map[string]bool, len(paths))
		for _, p := range paths {
			set[p] = true
		}
		for _, f := range files {
			if !set[f] {
				t.Errorf("ReadDir 缺少文件 %q", f)
			}
		}
	})
	t.Run("空目录", func(t *testing.T) {
		paths, err := ReadDir(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		if len(paths) != 0 {
			t.Errorf("空目录应返回 0 个文件，got %d", len(paths))
		}
	})
	t.Run("不存在的目录", func(t *testing.T) {
		if _, err := ReadDir(filepath.Join(t.TempDir(), "nope")); err == nil {
			t.Error("不存在的目录应返回错误")
		}
	})
}

func TestAbsInProject(t *testing.T) {
	t.Chdir(t.TempDir())
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	abs, _ := filepath.Abs(".")

	cases := []struct {
		name   string
		rel    string
		wantOK bool
		want   string
	}{
		{"普通相对路径", "a/b.txt", true, filepath.Join(root, "a", "b.txt")},
		{"正斜杠", "a/b/c.go", true, filepath.Join(root, "a", "b", "c.go")},
		{"含 .. 的合法路径", "a/../b.txt", true, filepath.Join(root, "b.txt")},
		{"项目根目录内的绝对路径", filepath.Join(root, "a", "b.txt"), true, filepath.Join(root, "a", "b.txt")},
		{"点", ".", true, root},
		{"项目根目录本身", abs, true, abs},
		{"项目根目录外的绝对路径", filepath.Join(filepath.Dir(root), "x.txt"), false, ""},
		{"越出项目根目录", "../x", false, ""},
		{"多层越界", "../../x", false, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := AbsInProject(c.rel)
			if ok != c.wantOK {
				t.Errorf("AbsInProject(%q) ok = %v, want %v", c.rel, ok, c.wantOK)
				return
			}
			if c.wantOK && got != c.want {
				t.Errorf("AbsInProject(%q) = %q, want %q", c.rel, got, c.want)
			}
		})
	}
}

func TestRemoveFileWithRetry(t *testing.T) {
	t.Run("删除存在的文件", func(t *testing.T) {
		f := filepath.Join(t.TempDir(), "a.txt")
		if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if !RemoveFileWithRetry(f) {
			t.Error("RemoveFileWithRetry 应返回 true")
		}
		if _, err := os.Stat(f); !os.IsNotExist(err) {
			t.Error("文件应已被删除")
		}
	})
	t.Run("删除不存在的文件", func(t *testing.T) {
		if RemoveFileWithRetry(filepath.Join(t.TempDir(), "nope.txt")) {
			t.Error("RemoveFileWithRetry 对不存在文件应返回 false")
		}
	})
}

func TestRemoveDirWithRetry(t *testing.T) {
	t.Run("删除空目录", func(t *testing.T) {
		d := filepath.Join(t.TempDir(), "d")
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
		if !RemoveDirWithRetry(d) {
			t.Error("RemoveDirWithRetry 应返回 true")
		}
		if _, err := os.Stat(d); !os.IsNotExist(err) {
			t.Error("目录应已被删除")
		}
	})
	t.Run("删除非空目录", func(t *testing.T) {
		d := t.TempDir()
		if err := os.WriteFile(filepath.Join(d, "f"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if RemoveDirWithRetry(d) {
			t.Error("RemoveDirWithRetry 对非空目录应返回 false")
		}
	})
	t.Run("删除不存在的目录", func(t *testing.T) {
		if RemoveDirWithRetry(filepath.Join(t.TempDir(), "nope")) {
			t.Error("RemoveDirWithRetry 对不存在目录应返回 false")
		}
	})
}

func TestRemoveEmptyParents(t *testing.T) {
	t.Chdir(t.TempDir())
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("整链删除空目录", func(t *testing.T) {
		f := filepath.Join(root, "a", "b", "c", "file.txt")
		if err := os.MkdirAll(filepath.Dir(f), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(f); err != nil {
			t.Fatal(err)
		}
		RemoveEmptyParents(f)
		for _, p := range []string{
			filepath.Join(root, "a", "b", "c"),
			filepath.Join(root, "a", "b"),
			filepath.Join(root, "a"),
		} {
			if _, err := os.Stat(p); !os.IsNotExist(err) {
				t.Errorf("空目录 %q 应被删除", p)
			}
		}
		if _, err := os.Stat(root); err != nil {
			t.Errorf("项目根目录不应被删除: %v", err)
		}
	})

	t.Run("遇非空目录停止", func(t *testing.T) {
		f := filepath.Join(root, "x", "y", "file.txt")
		if err := os.MkdirAll(filepath.Dir(f), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "x", "keep.txt"), []byte("k"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(f); err != nil {
			t.Fatal(err)
		}
		RemoveEmptyParents(f)
		if _, err := os.Stat(filepath.Join(root, "x", "y")); !os.IsNotExist(err) {
			t.Errorf("空目录 %q 应被删除", filepath.Join(root, "x", "y"))
		}
		if _, err := os.Stat(filepath.Join(root, "x")); err != nil {
			t.Errorf("非空目录不应被删除: %v", err)
		}
	})
}

func TestRemoveAllWithRetry(t *testing.T) {
	t.Chdir(t.TempDir())
	root := t.TempDir()

	t.Run("删除 root 之下的目录树", func(t *testing.T) {
		dir := filepath.Join(root, "a", "b")
		if err := os.MkdirAll(filepath.Join(dir, "c"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "c", "f.txt"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if !RemoveAllWithRetry(dir, root) {
			t.Error("RemoveAllWithRetry 应返回 true")
		}
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			t.Error("目录树应已被删除")
		}
		if _, err := os.Stat(root); err != nil {
			t.Error("root 不应被删除")
		}
	})

	t.Run("path 为 root 本身应拒绝", func(t *testing.T) {
		if RemoveAllWithRetry(root, root) {
			t.Error("path 等于 root 时应返回 false")
		}
	})

	t.Run("path 越出 root 应拒绝", func(t *testing.T) {
		outside := filepath.Join(filepath.Dir(root), "other")
		if RemoveAllWithRetry(outside, root) {
			t.Error("path 不在 root 之下时应返回 false")
		}
	})

	t.Run("path 不存在但位于 root 之下", func(t *testing.T) {
		if !RemoveAllWithRetry(filepath.Join(root, "nope"), root) {
			t.Error("root 之下的不存在路径应返回 true")
		}
	})

	t.Run("相对 root 与相对 path", func(t *testing.T) {
		if err := os.MkdirAll(filepath.Join("views", "a"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join("views", "a", "f.txt"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if !RemoveAllWithRetry("views/a", "views") {
			t.Error("相对 root 应返回 true")
		}
		if _, err := os.Stat(filepath.Join("views", "a")); !os.IsNotExist(err) {
			t.Error("views/a 应被删除")
		}
		if _, err := os.Stat("views"); err != nil {
			t.Error("root(views) 不应被删除")
		}
	})

	t.Run("root 为空时取当前项目根目录", func(t *testing.T) {
		if err := os.MkdirAll(filepath.Join("er", "x"), 0o755); err != nil {
			t.Fatal(err)
		}
		if !RemoveAllWithRetry("er/x", "") {
			t.Error("root 为空应基于当前工作目录，返回 true")
		}
		if _, err := os.Stat(filepath.Join("er", "x")); !os.IsNotExist(err) {
			t.Error("er/x 应被删除")
		}
	})
}

func TestRemoveFileWithCleanup(t *testing.T) {
	t.Chdir(t.TempDir())

	f := filepath.Join("a", "b", "c.txt")
	if err := os.MkdirAll(filepath.Dir(f), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	RemoveFileWithCleanup("a/b/c.txt")

	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Error("文件应被删除")
	}
	// 空目录链 a、b 应被一并清理
	if _, err := os.Stat(filepath.Join("a")); !os.IsNotExist(err) {
		t.Error("空目录链应被清理")
	}
}

func TestRemoveFileWithCleanupAbs(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)

	f := filepath.Join("a", "b", "c.txt")
	if err := os.MkdirAll(filepath.Dir(f), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	RemoveFileWithCleanup(filepath.Join(root, "a/b/c.txt"))

	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Error("文件应被删除")
	}
	// 空目录链 a、b 应被一并清理
	if _, err := os.Stat(filepath.Join("a")); !os.IsNotExist(err) {
		t.Error("空目录链应被清理")
	}
}
