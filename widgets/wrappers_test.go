package ant_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestListRendersAndSelects(t *testing.T) {
	rows := []string{"苹果", "香蕉", "橙子"}
	state := &ui.ListState{}
	chosen := -1
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.List(c, state, len(rows), func(i int) {
			widgets.Paragraph(c, rows[i])
		}, widgets.ListSelected(&chosen))
	}
	tt := ui.NewTester(view, 400, 220)
	tt.Frame()
	for _, want := range rows {
		if !tt.HasText(want) {
			t.Errorf("列表行 %q 未渲染", want)
		}
	}
	if err := tt.Click("香蕉"); err != nil {
		t.Fatalf("点击列表行失败: %v", err)
	}
	tt.Frame()
	if chosen != 1 {
		t.Errorf("选中行 = %d, want 1", chosen)
	}
}

func TestSidebarChooses(t *testing.T) {
	chosen := "inbox"
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Sidebar(c, &chosen, func() {
			widgets.SidebarSection(c, "邮件", nil, func() {
				widgets.SidebarItem(c, "inbox", "收件箱")
				widgets.SidebarItem(c, "sent", "已发送")
			})
		})
	}
	tt := ui.NewTester(view, 640, 260)
	tt.Frame()
	if !tt.HasText("收件箱") || !tt.HasText("已发送") {
		t.Fatal("侧栏条目未渲染")
	}
	if err := tt.Click("已发送"); err != nil {
		t.Fatalf("点击侧栏条目失败: %v", err)
	}
	tt.Frame()
	if chosen != "sent" {
		t.Errorf("选中项 = %q, want sent", chosen)
	}
}

func TestSplitterRendersBothPanes(t *testing.T) {
	left := float32(160)
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Splitter(c, &left, func() {
			widgets.Paragraph(c, "左栏内容")
		}, func() {
			widgets.Paragraph(c, "右栏内容")
		}).Fill()
	}
	tt := ui.NewTester(view, 640, 320)
	tt.Frame()
	if !tt.HasText("左栏内容") || !tt.HasText("右栏内容") {
		t.Error("分隔面板两栏未渲染")
	}
}

func TestAutoCompleteSuggests(t *testing.T) {
	city := ""
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.AutoComplete(c, &city, []string{"北京", "上海", "广州"})
	}
	tt := ui.NewTester(view, 360, 240)
	tt.Frame()
	tt.ClickAt(100, 16)
	tt.Type("北")
	tt.Frame()
	if !tt.HasText("北京") {
		t.Error("自动补全建议未显示")
	}
	if city != "北" {
		t.Errorf("键入值 = %q, want 北", city)
	}
}

func TestWrappersFollowThemeAndDensity(t *testing.T) {
	for _, tk := range []*theme.AntToken{
		theme.Light(),
		theme.Dark(),
		theme.Light(theme.WithDensity(theme.Compact)),
		theme.Dark(theme.WithDensity(theme.Compact), theme.WithPrimary(ui.Hex("#7c3aed"))),
	} {
		token := tk
		state := &ui.ListState{}
		chosen := "a"
		view := func(c *ui.Context) {
			theme.Apply(c, token)
			col := ui.Column(c).Gap(8).FillWidth()
			col.Children(func() {
				widgets.List(c, state, 3, func(i int) {
					ui.Textf(c, "第 %d 行", i+1)
				})
				widgets.Sidebar(c, &chosen, func() {
					widgets.SidebarSection(c, "分组", nil, func() {
						widgets.SidebarItem(c, "a", "条目 A")
						widgets.SidebarItem(c, "b", "条目 B")
					})
				})
				widgets.AutoComplete(c, new(string), []string{"候选一"})
				widgets.SplitterPercent(c, new(float32), func() {
					widgets.Paragraph(c, "上")
				}, func() {
					widgets.Paragraph(c, "下")
				}).Height(120)
			})
		}
		tt := ui.NewTester(view, 560, 520)
		tt.Frame()
		if tt.FramePending() {
			t.Errorf("密度 %v 下包装控件持续请求帧", token.Density)
		}
	}
}
