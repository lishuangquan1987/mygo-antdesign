package sample_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
)

// controlManifest 记录逐示例基线的生成平台（字体栈不同会导致跨平台像素差异）。
type controlManifest struct {
	GOOS   string `json:"goos"`
	GOARCH string `json:"goarch"`
}

// controlFile 把示例名转成安全文件名。
func controlFile(group string, index int, title string) string {
	name := fmt.Sprintf("%s-%02d-%s", group, index+1, title)
	repl := strings.NewReplacer("/", "_", `\`, "_", ":", "_", "*", "_", "?", "_", " ", "_")
	return repl.Replace(name) + "-light.png"
}

func twoDigits(n int) string {
	if n < 10 {
		return "0" + string(rune('0'+n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}

// animated 记录含加载动画（ui.Spinner）的示例：逐帧像素不稳定，不纳入基线。
var animated = map[string]bool{
	"Button/状态": true,
	"反馈/加载中":    true,
}

// TestControlScreenshots 生成/比对「逐示例」截图基线（docs/07 §5.2 / F6）：
// 每个示例以 480×360 亮色标准密度渲染，与 samples/golden/controls/ 逐字节
// 比对；首次运行生成基线。跨平台只渲染不比对（基线带 manifest 平台信息）。
func TestControlScreenshots(t *testing.T) {
	dir := filepath.Join("golden", "controls")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(dir, "manifest.json")
	var m controlManifest
	if data, err := os.ReadFile(manifestPath); err == nil {
		_ = json.Unmarshal(data, &m)
	}
	otherPlatform := m.GOOS != "" && (m.GOOS != runtime.GOOS || m.GOARCH != runtime.GOARCH)

	render := func(view func(*ui.Context)) []byte {
		tt := ui.NewTester(view, 480, 360)
		tt.Frame()
		img := tt.Image()
		if img == nil {
			t.Fatal("无渲染帧")
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			t.Fatalf("编码 PNG 失败: %v", err)
		}
		return buf.Bytes()
	}

	wrote := false
	for _, demo := range allDemos() {
		for i, s := range demo.Samples {
			s := s
			view := func(c *ui.Context) {
				theme.Apply(c, theme.Light())
				s.View(c)
			}
			path := filepath.Join(dir, controlFile(demo.Name, i, s.Title))
			if animated[demo.Name+"/"+s.Title] {
				// 含加载动画（ui.Spinner）的示例逐帧像素不稳定，不纳入逐字节基线。
				_ = os.Remove(path)
				continue
			}
			got := render(view)
			if otherPlatform {
				continue
			}
			want, err := os.ReadFile(path)
			if err != nil {
				if err := os.WriteFile(path, got, 0o644); err != nil {
					t.Fatal(err)
				}
				wrote = true
				continue
			}
			if !bytes.Equal(want, got) {
				t.Errorf("逐示例截图与基线不一致：%s（如为预期改动，删除该文件重新生成）", path)
			}
		}
	}
	if wrote {
		data, err := json.MarshalIndent(controlManifest{GOOS: runtime.GOOS, GOARCH: runtime.GOARCH}, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(manifestPath, append(data, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
