package ant_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestButtonRendersAllVariants(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		col := ui.Column(c).Gap(8)
		col.Children(func() {
			widgets.Button(c, "默认")
			widgets.Button(c, "主要", widgets.Primary())
			widgets.Button(c, "虚线", widgets.Dashed())
			widgets.Button(c, "文字", widgets.Text())
			widgets.Button(c, "链接", widgets.Link())
			widgets.Button(c, "删除", widgets.Primary(), widgets.Danger())
		})
	}
	tt := ui.NewTester(view, 320, 400)
	tt.Frame()
	for _, want := range []string{"默认", "主要", "虚线", "文字", "链接", "删除"} {
		if !tt.HasText(want) {
			t.Errorf("变体按钮 %q 未渲染", want)
		}
	}
}

func TestButtonOnClick(t *testing.T) {
	n := 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Button(c, "计数", widgets.OnClick(func() { n++ }))
	}
	tt := ui.NewTester(view, 240, 120)
	if err := tt.Click("计数"); err != nil {
		t.Fatalf("点击失败: %v", err)
	}
	if n != 1 {
		t.Errorf("OnClick 触发 %d 次, want 1（一次点击上升沿一次）", n)
	}
	if err := tt.Click("计数"); err != nil {
		t.Fatalf("第二次点击失败: %v", err)
	}
	if n != 2 {
		t.Errorf("OnClick 累计 %d 次, want 2", n)
	}
}

func TestButtonDisabled(t *testing.T) {
	n := 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Button(c, "禁用", widgets.Disabled(true), widgets.OnClick(func() { n++ }))
	}
	tt := ui.NewTester(view, 240, 120)
	_ = tt.Click("禁用")
	if n != 0 {
		t.Errorf("禁用按钮仍触发了 OnClick")
	}
}

func TestButtonLoadingBlocksClick(t *testing.T) {
	n := 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Button(c, "加载", widgets.Loading(true), widgets.OnClick(func() { n++ }))
	}
	tt := ui.NewTester(view, 240, 120)
	tt.Frame()
	if !tt.HasText("加载") {
		t.Fatal("loading 按钮未渲染")
	}
	_ = tt.Click("加载")
	if n != 0 {
		t.Errorf("loading 按钮不应触发 OnClick")
	}
}

func TestButtonFollowsThemeAndDensity(t *testing.T) {
	// 标准 / 紧凑两种密度、亮 / 暗两种主题下渲染均不应 panic，
	// 且主题切换后新主题生效（Apply 的闭环已在 theme 包验证）。
	for _, tk := range []*theme.AntToken{
		theme.Light(),
		theme.Dark(),
		theme.Light(theme.WithDensity(theme.Compact)),
		theme.Dark(theme.WithDensity(theme.Compact), theme.WithPrimary(ui.Hex("#7c3aed"))),
	} {
		token := tk
		view := func(c *ui.Context) {
			theme.Apply(c, token)
			widgets.Button(c, "保存", widgets.Primary(), widgets.Size(widgets.Large))
			widgets.Button(c, "取消", widgets.Size(widgets.Small))
		}
		tt := ui.NewTester(view, 320, 200)
		tt.Frame()
		if !tt.HasText("保存") || !tt.HasText("取消") {
			t.Errorf("密度 %v 下按钮未渲染完整", token.Density)
		}
	}
}
