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

// controlCombo 一种亮暗 × 密度组合。
type controlCombo struct {
	name    string
	dark    bool
	density theme.Density
}

// controlCombos 与 Gallery 级基线保持一致的 4 种组合。
var controlCombos = []controlCombo{
	{"light-standard", false, theme.Comfortable},
	{"dark-standard", true, theme.Comfortable},
	{"light-compact", false, theme.Compact},
	{"dark-compact", true, theme.Compact},
}

// controlFile 把示例名与组合转成安全文件名。
func controlFile(group string, index int, title, combo string) string {
	name := fmt.Sprintf("%s-%02d-%s-%s", group, index+1, title, combo)
	repl := strings.NewReplacer("/", "_", `\`, "_", ":", "_", "*", "_", "?", "_", " ", "_")
	return repl.Replace(name) + ".png"
}

// animated 记录含加载动画或网络依赖的示例：逐帧像素不稳定，不纳入基线。
var animated = map[string]bool{
	"Button/状态":             true,
	"反馈/加载中":                true,
	"原生补全/RemoteImage 网络图片": true,
}

// TestControlScreenshots 生成/比对「逐示例」截图基线（docs/07 §5.2 / F6）：
// 每个示例按 4 种亮暗 × 密度组合以 480×360 渲染，与
// samples/golden/controls/ 逐字节比对；首次运行生成基线。跨平台只渲染
// 不比对（基线带 manifest 平台信息）。
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

	render := func(s func(*ui.Context), cb controlCombo) []byte {
		view := func(c *ui.Context) {
			token := theme.Light(theme.WithDensity(cb.density))
			if cb.dark {
				token = theme.Dark(theme.WithDensity(cb.density))
			}
			theme.Apply(c, token)
			s(c)
		}
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
			if animated[demo.Name+"/"+s.Title] {
				// 含加载动画的示例逐帧像素不稳定，不纳入逐字节基线。
				continue
			}
			for _, cb := range controlCombos {
				path := filepath.Join(dir, controlFile(demo.Name, i, s.Title, cb.name))
				got := render(s.View, cb)
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
