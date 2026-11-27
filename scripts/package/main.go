// Command package 把 Gallery 可执行文件与文档打包成发布归档（zip 或 tar.gz）。
//
// 纯 Go 实现（不依赖 zip/tar 等系统命令），因此 Linux CI 与本地 Windows/
// macOS 行为一致。发布工作流用法：
//
//	GOOS= GOARCH= go run ./scripts/package \
//	  -bin dist/mygo-ant-design-gallery \
//	  -out dist/mygo-ant-design-gallery-0.3.0-linux-amd64.tar.gz \
//	  -app mygo-ant-design-gallery -version 0.3.0 \
//	  -os linux -arch amd64 -mode tar.gz
package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	var (
		bin     = flag.String("bin", "", "可执行文件路径")
		out     = flag.String("out", "", "输出归档路径（含扩展名）")
		app     = flag.String("app", "mygo-ant-design-gallery", "应用名（归档内目录与可执行文件名）")
		version = flag.String("version", "dev", "版本号")
		goos    = flag.String("os", "", "目标系统（写进归档名与说明）")
		goarch  = flag.String("arch", "", "目标架构（写进归档名与说明）")
		mode    = flag.String("mode", "zip", "归档格式：zip 或 tar.gz")
		readme  = flag.String("readme", "README.md", "要一并打包的 README 路径（可空）")
		license = flag.String("license", "LICENSE", "要一并打包的 LICENSE 路径（可空）")
	)
	flag.Parse()

	if *bin == "" || *out == "" || *goos == "" || *goarch == "" {
		fatalf("-bin、-out、-os、-arch 均为必填")
	}
	if *mode != "zip" && *mode != "tar.gz" {
		fatalf("-mode 仅支持 zip 或 tar.gz，收到 %q", *mode)
	}

	name := fmt.Sprintf("%s-%s-%s-%s", *app, *version, *goos, *goarch)
	dir := filepath.Join(filepath.Dir(*out), name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fatalf("创建目录失败：%v", err)
	}

	exe := *app
	if *goos == "windows" {
		exe += ".exe"
	}
	if err := copyFile(*bin, filepath.Join(dir, exe), 0o755); err != nil {
		fatalf("复制可执行文件失败：%v", err)
	}
	for _, doc := range []struct{ src, dst string }{
		{*readme, "README.md"},
		{*license, "LICENSE"},
	} {
		if doc.src == "" {
			continue
		}
		if _, err := os.Stat(doc.src); err != nil {
			continue // 文档缺失不阻断发布
		}
		if err := copyFile(doc.src, filepath.Join(dir, doc.dst), 0o644); err != nil {
			fatalf("复制 %s 失败：%v", doc.dst, err)
		}
	}

	runText := fmt.Sprintf(`mygo-antdesign Gallery %s
Ant Design 风格的 MyGo 原生 UI 控件库 · 示例程序

运行方式：
  Windows: 双击 %s.exe
  macOS:   chmod +x %s && ./%s
  Linux:   chmod +x %s && ./%s

说明：
  - 纯 Go 构建，无 cgo 依赖；需要桌面环境（Windows 10+/macOS 11+/X11 或 Wayland）。
  - 仓库：https://github.com/lishuangquan1987/mygo-antdesign
  - 文档：https://github.com/lishuangquan1987/mygo-antdesign/tree/main/docs
`, *version, *app, *app, *app, *app, *app)
	if err := os.WriteFile(filepath.Join(dir, "RUN.txt"), []byte(runText), 0o644); err != nil {
		fatalf("写入 RUN.txt 失败：%v", err)
	}

	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		fatalf("创建输出目录失败：%v", err)
	}
	var err error
	if *mode == "zip" {
		err = writeZip(dir, *out)
	} else {
		err = writeTarGz(dir, *out)
	}
	if err != nil {
		fatalf("打包失败：%v", err)
	}
	fmt.Printf("已生成 %s\n", *out)
}

func writeZip(dir, out string) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	defer zw.Close()

	base := filepath.Base(dir)
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		hdr, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		hdr.Name = filepath.ToSlash(filepath.Join(base, rel))
		hdr.Method = zip.Deflate
		w, err := zw.CreateHeader(hdr)
		if err != nil {
			return err
		}
		return copyInto(w, path)
	})
}

func writeTarGz(dir, out string) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	gz := gzip.NewWriter(f)
	defer gz.Close()
	tw := tar.NewWriter(gz)
	defer tw.Close()

	base := filepath.Base(dir)
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		name := filepath.ToSlash(filepath.Join(base, rel))
		if info.IsDir() {
			return tw.WriteHeader(&tar.Header{
				Name:     name + "/",
				Mode:     0o755,
				Typeflag: tar.TypeDir,
			})
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = name
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		return copyInto(tw, path)
	})
}

func copyInto(w io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(w, f)
	return err
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Chmod(dst, mode)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "package: "+format+"\n", args...)
	os.Exit(1)
}
