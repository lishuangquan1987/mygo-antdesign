package sample_test

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	sample "github.com/lishuangquan1987/mygo-antdesign/samples"
	"github.com/lishuangquan1987/mygo-antdesign/samples/button"
	"github.com/lishuangquan1987/mygo-antdesign/samples/display"
	"github.com/lishuangquan1987/mygo-antdesign/samples/divider"
	"github.com/lishuangquan1987/mygo-antdesign/samples/extrademo"
	"github.com/lishuangquan1987/mygo-antdesign/samples/feedback"
	"github.com/lishuangquan1987/mygo-antdesign/samples/ide"
	"github.com/lishuangquan1987/mygo-antdesign/samples/input"
	"github.com/lishuangquan1987/mygo-antdesign/samples/layout"
	"github.com/lishuangquan1987/mygo-antdesign/samples/nav"
	"github.com/lishuangquan1987/mygo-antdesign/samples/selectdemo"
	"github.com/lishuangquan1987/mygo-antdesign/samples/table"
	"github.com/lishuangquan1987/mygo-antdesign/samples/typography"
	"github.com/lishuangquan1987/mygo-antdesign/theme"
)

// allDemos 返回 Gallery 注册的全部示例分组（与 gallery/main.go 保持一致）。
func allDemos() []sample.WidgetDemo {
	return []sample.WidgetDemo{
		button.Demo(), typography.Demo(), divider.Demo(), layout.Demo(),
		input.Demo(), selectdemo.Demo(), display.Demo(), table.Demo(),
		feedback.Demo(), nav.Demo(), extrademo.Demo(), ide.Demo(),
	}
}

// TestAllDemoSamplesRender 校验每个示例分组：
//   - 至少一个示例；
//   - 代码片段由 // sample: 标记成功提取（非空、非错误提示）；
//   - 在亮/暗主题下都能渲染出帧。
func TestAllDemoSamplesRender(t *testing.T) {
	for _, d := range allDemos() {
		if d.Name == "" || d.Summary == "" {
			t.Errorf("示例分组缺少 Name/Summary：%+v", d)
		}
		if len(d.Samples) == 0 {
			t.Fatalf("%s 没有示例", d.Name)
		}
		for _, s := range d.Samples {
			s := s
			if s.Title == "" || s.Description == "" {
				t.Errorf("%s 的示例缺少标题/描述", d.Name)
			}
			if strings.TrimSpace(s.Code) == "" ||
				strings.HasPrefix(s.Code, "// 未找到") ||
				strings.HasPrefix(s.Code, "// 标记段") {
				t.Errorf("%s/%s 代码片段提取失败：%.60q", d.Name, s.Title, s.Code)
				continue
			}
			if s.View == nil {
				t.Errorf("%s/%s 缺少 View", d.Name, s.Title)
				continue
			}
			for _, tk := range []*theme.AntToken{theme.Light(), theme.Dark()} {
				token := tk
				view := func(c *ui.Context) {
					theme.Apply(c, token)
					s.View(c)
				}
				tt := ui.NewTester(view, 720, 520)
				tt.Frame()
				if tt.Image() == nil {
					t.Errorf("%s/%s 在 %v 主题下无渲染帧", d.Name, s.Title, token.Density)
				}
			}
		}
	}
}
