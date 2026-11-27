package ant_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestTagVariants(t *testing.T) {
	closed := 0
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		row := ui.Row(c).Gap(8).Wrap()
		row.Children(func() {
			widgets.Tag(c, "默认")
			widgets.Tag(c, "成功", widgets.TagSuccess())
			widgets.Tag(c, "警告", widgets.TagWarning())
			widgets.Tag(c, "错误", widgets.TagError())
			widgets.Tag(c, "品牌", widgets.TagPrimary())
			widgets.Tag(c, "可关闭", widgets.TagClosable(func() { closed++ }))
		})
	}
	tt := ui.NewTester(view, 480, 160)
	tt.Frame()
	for _, want := range []string{"默认", "成功", "警告", "错误", "品牌", "可关闭"} {
		if !tt.HasText(want) {
			t.Errorf("标签 %q 未渲染", want)
		}
	}
	if err := tt.Click("关闭"); err != nil {
		t.Fatalf("点击关闭失败: %v", err)
	}
	if closed != 1 {
		t.Errorf("TagClosable 触发 %d 次, want 1", closed)
	}
}

func TestBadgeAndAvatar(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		row := ui.Row(c).Gap(16).AlignItems(ui.Center)
		row.Children(func() {
			widgets.Badge(c, 5, func() { widgets.Button(c, "消息") })
			widgets.Badge(c, 120, func() { widgets.Button(c, "通知") })
			widgets.BadgeDot(c, func() { widgets.AvatarEl(c, "李") })
			widgets.AvatarEl(c, "Wang")
		})
	}
	tt := ui.NewTester(view, 480, 160)
	tt.Frame()
	if !tt.HasText("消息") || !tt.HasText("通知") {
		t.Error("Badge 内容未渲染")
	}
	if !tt.HasText("5") || !tt.HasText("99+") {
		t.Error("徽标数字未渲染（5 与 99+）")
	}
}

func TestCardLayout(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Card(c, "卡片标题",
			func() {
				widgets.Paragraph(c, "卡片内容第一行")
				widgets.Paragraph(c, "卡片内容第二行", widgets.Secondary())
			},
			widgets.CardExtra(func(c *ui.Context) {
				widgets.Button(c, "更多", widgets.Text())
			}),
			widgets.CardActions(func(c *ui.Context) {
				widgets.Button(c, "取消")
				widgets.Button(c, "确定", widgets.Primary())
			}))
	}
	tt := ui.NewTester(view, 400, 260)
	tt.Frame()
	for _, want := range []string{"卡片标题", "卡片内容第一行", "更多", "取消", "确定"} {
		if !tt.HasText(want) {
			t.Errorf("卡片元素 %q 未渲染", want)
		}
	}
}

func TestDisplayFollowsThemeAndDensity(t *testing.T) {
	for _, tk := range []*theme.AntToken{
		theme.Light(),
		theme.Dark(),
		theme.Light(theme.WithDensity(theme.Compact)),
		theme.Dark(theme.WithDensity(theme.Compact), theme.WithPrimary(ui.Hex("#7c3aed"))),
	} {
		token := tk
		view := func(c *ui.Context) {
			theme.Apply(c, token)
			col := ui.Column(c).Gap(8)
			col.Children(func() {
				widgets.Tag(c, "标签", widgets.TagPrimary())
				widgets.Tag(c, "青色", widgets.TagCyan())
				widgets.Tag(c, "金色", widgets.TagGold())
				widgets.Badge(c, 3, func() { widgets.AvatarEl(c, "王") })
				widgets.Card(c, "标题", func() { widgets.Paragraph(c, "内容") })
			})
		}
		tt := ui.NewTester(view, 360, 320)
		tt.Frame()
		if tt.FramePending() {
			t.Errorf("密度 %v 下数据展示控件持续请求帧", token.Density)
		}
	}
}
