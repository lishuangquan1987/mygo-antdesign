package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// prepare 构造一个打包目录，返回其路径与输出目录。
func prepare(t *testing.T) (pkgDir, outDir string) {
	t.Helper()
	root := t.TempDir()
	pkgDir = filepath.Join(root, "mygo-ant-design-gallery-0.3.0-linux-amd64")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "mygo-ant-design-gallery"), []byte("binary"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "RUN.txt"), []byte("run"), 0o644); err != nil {
		t.Fatal(err)
	}
	return pkgDir, root
}

func TestWriteZip(t *testing.T) {
	pkgDir, outDir := prepare(t)
	out := filepath.Join(outDir, "app.zip")
	if err := writeZip(pkgDir, out); err != nil {
		t.Fatal(err)
	}
	r, err := zip.OpenReader(out)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	want := map[string]bool{
		"mygo-ant-design-gallery-0.3.0-linux-amd64/mygo-ant-design-gallery": false,
		"mygo-ant-design-gallery-0.3.0-linux-amd64/RUN.txt":                 false,
	}
	for _, f := range r.File {
		if _, ok := want[f.Name]; ok {
			want[f.Name] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("zip 缺少 %s", name)
		}
	}
}

func TestWriteTarGz(t *testing.T) {
	pkgDir, outDir := prepare(t)
	out := filepath.Join(outDir, "app.tar.gz")
	if err := writeTarGz(pkgDir, out); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	names := map[string]bool{}
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		names[hdr.Name] = true
	}
	for _, want := range []string{
		"mygo-ant-design-gallery-0.3.0-linux-amd64/mygo-ant-design-gallery",
		"mygo-ant-design-gallery-0.3.0-linux-amd64/RUN.txt",
	} {
		if !names[want] {
			t.Errorf("tar.gz 缺少 %s", want)
		}
	}
}
