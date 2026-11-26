// 本文件覆盖 Typography / Divider / Layout 三个示例包的通用校验：
// 元数据完整、代码提取成功、示例可在亮暗两种主题下渲染。
package typography_test

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	"github.com/lishuangquan1987/mygo-antdesign/samples/divider"
	"github.com/lishuangquan1987/mygo-antdesign/samples/input"
	"github.com/lishuangquan1987/mygo-antdesign/samples/layout"
	selectdemo "github.com/lishuangquan1987/mygo-antdesign/samples/selectdemo"
	"github.com/lishuangquan1987/mygo-antdesign/samples/typography"
	"github.com/lishuangquan1987/mygo-antdesign/theme"
)

var demos = []struct {
	name string
	d    samples.WidgetDemo
}{
	{"Typography", typography.Demo()},
	{"Divider", divider.Demo()},
	{"Layout", layout.Demo()},
	{"Input", input.Demo()},
	{"Select", selectdemo.Demo()},
}

func TestDemoCompleteness(t *testing.T) {
	for _, dc := range demos {
		d := dc.d
		if d.Name == "" || d.Category == "" || d.Summary == "" {
			t.Errorf("%s: Demo 元数据不完整: %v / %v / %v", dc.name, d.Name, d.Category, d.Summary)
		}
		if len(d.Samples) < 3 {
			t.Fatalf("%s: 示例数 = %d, want >= 3", dc.name, len(d.Samples))
		}
		for _, s := range d.Samples {
			if s.Title == "" || s.Description == "" {
				t.Errorf("%s: 示例 %q 缺少标题或描述", dc.name, s.Title)
			}
			if s.Code == "" || strings.Contains(s.Code, "// 未找到") || strings.Contains(s.Code, "// 标记段") {
				t.Errorf("%s: 示例 %q 的代码片段提取失败: %.40s…", dc.name, s.Title, s.Code)
			}
			if !strings.Contains(s.Code, "widgets.") {
				t.Errorf("%s: 示例 %q 的代码不是真实用法", dc.name, s.Title)
			}
			if s.View == nil {
				t.Errorf("%s: 示例 %q 缺少 View", dc.name, s.Title)
			}
		}
	}
}

func TestDemoViewsRenderAllThemes(t *testing.T) {
	for _, dc := range demos {
		for _, s := range dc.d.Samples {
			for _, tk := range []*theme.AntToken{
				theme.Light(),
				theme.Dark(),
				theme.Light(theme.WithDensity(theme.Compact)),
			} {
				token := tk
				v := s.View
				tt := ui.NewTester(func(c *ui.Context) {
					theme.Apply(c, token)
					v(c)
				}, 480, 240)
				tt.Frame()
				if tt.FramePending() {
					t.Errorf("%s 示例 %q 持续请求帧", dc.name, s.Title)
				}
			}
		}
	}
}
