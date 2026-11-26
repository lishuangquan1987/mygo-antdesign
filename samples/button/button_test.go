package button_test

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	"github.com/lishuangquan1987/mygo-antdesign/samples/button"
	"github.com/lishuangquan1987/mygo-antdesign/theme"
)

func TestDemoCompleteness(t *testing.T) {
	d := button.Demo()
	if d.Name != "Button" || d.Category != samples.CatGeneral {
		t.Errorf("Demo 元数据不完整: %v / %v", d.Name, d.Category)
	}
	if len(d.Samples) < 4 {
		t.Fatalf("示例数 = %d, want >= 4", len(d.Samples))
	}
	for _, s := range d.Samples {
		if s.Title == "" || s.Description == "" {
			t.Errorf("示例 %q 缺少标题或描述", s.Title)
		}
		if s.Code == "" || strings.Contains(s.Code, "// 未找到") || strings.Contains(s.Code, "// 标记段") {
			t.Errorf("示例 %q 的代码片段提取失败: %.40s…", s.Title, s.Code)
		}
		if !strings.Contains(s.Code, "widgets.Button(") {
			t.Errorf("示例 %q 的代码不是真实用法（缺少 widgets.Button 调用）", s.Title)
		}
		if s.View == nil {
			t.Errorf("示例 %q 缺少 View", s.Title)
		}
	}
}

func TestDemoViewsRender(t *testing.T) {
	d := button.Demo()
	for _, s := range d.Samples {
		v := s.View
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c, theme.Light())
			v(c)
		}, 480, 160)
		tt.Frame()
		if tt.FramePending() {
			t.Errorf("示例 %q 持续请求帧（可能存在每帧变化的状态）", s.Title)
		}
	}
}

func TestDemoInteractive(t *testing.T) {
	// 基础用法示例应可交互：点击「保存」出现「已保存」。
	d := button.Demo()
	basic := d.Samples[0]
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		basic.View(c)
	}, 320, 120)
	if err := tt.Click("保存"); err != nil {
		t.Fatalf("点击失败: %v", err)
	}
	if !tt.HasText("已保存") {
		t.Errorf("点击后未出现「已保存」反馈")
	}
}
