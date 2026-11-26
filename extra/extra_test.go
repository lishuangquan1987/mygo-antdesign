package extra_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/extra"
	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestTreeSelectLeaf(t *testing.T) {
	selected := ""
	nodes := []extra.TreeNode{
		{Label: "src", Children: []extra.TreeNode{
			{Label: "main.go", Key: "main"},
			{Label: "go.mod"},
		}},
	}
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		extra.Tree(c, nodes, extra.OnTreeNodeSelect(func(key string) { selected = key }))
	}
	tt := ui.NewTester(view, 320, 200)
	tt.Frame()
	if !tt.HasText("src") {
		t.Fatal("树分支未渲染")
	}
	// 聚焦分支后用右方向键展开（点击行是选中语义，箭头才是展开）。
	if err := tt.Click("src"); err != nil {
		t.Fatalf("点击分支失败: %v", err)
	}
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	if !tt.HasText("main.go") {
		t.Fatal("展开后叶子未渲染")
	}
	if err := tt.Click("main.go"); err != nil {
		t.Fatalf("点击叶子失败: %v", err)
	}
	tt.Frame()
	if selected != "main" {
		t.Errorf("选中 = %q, want main", selected)
	}
}

func TestTransferMoveItems(t *testing.T) {
	right := []string{"b"}
	items := []extra.TransferItem{
		{Key: "a", Label: "选项 A"},
		{Key: "b", Label: "选项 B"},
		{Key: "c", Label: "选项 C"},
	}
	var saved []string
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		extra.Transfer(c, items, &right, extra.TransferOnChange(func(r []string) { saved = r }))
	}
	tt := ui.NewTester(view, 560, 240)
	tt.Frame()
	if !tt.HasText("选项 B") {
		t.Fatal("右列未渲染")
	}
	// 左列点选「选项 A」，移到右列。
	if err := tt.Click("选项 A"); err != nil {
		t.Fatalf("点击左列行失败: %v", err)
	}
	tt.Frame()
	if err := tt.Click("▶"); err != nil {
		t.Fatalf("点击移动按钮失败: %v", err)
	}
	tt.Frame()
	if len(saved) != 2 {
		t.Fatalf("右列 = %v, want [b a]（含 a）", saved)
	}
	hasA := false
	for _, k := range saved {
		if k == "a" {
			hasA = true
		}
	}
	if !hasA {
		t.Errorf("右列缺少 a：%v", saved)
	}
}

func TestCarouselSwitchPage(t *testing.T) {
	pages := []func(*ui.Context){
		func(c *ui.Context) { widgets.Title(c, "第一页", widgets.H3) },
		func(c *ui.Context) { widgets.Title(c, "第二页", widgets.H3) },
	}
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		extra.Carousel(c, pages)
	}
	tt := ui.NewTester(view, 360, 200)
	tt.Frame()
	if !tt.HasText("第一页") {
		t.Fatal("首页未渲染")
	}
	// 指示点无文字：页面区高 80 + 间距 8 → 点中心 y≈96；
	// 两点总宽 24 居中于 360，第二点中心 x≈188。
	tt.ClickAt(188, 96)
	tt.Frame()
	if !tt.HasText("第二页") {
		t.Error("点击指示点未切换到第二页")
	}
}

func TestQRCodeRender(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		col := ui.Column(c).Gap(8)
		col.Children(func() {
			extra.QRCode(c, "https://example.com/mygo-antdesign")
			extra.QRCode(c, "")
		})
	}
	tt := ui.NewTester(view, 320, 320)
	tt.Frame()
	if tt.FramePending() {
		t.Error("二维码持续请求帧")
	}
	if !tt.HasText("无内容") {
		t.Error("空内容占位未渲染")
	}
}

func TestUploadRender(t *testing.T) {
	var listed []string
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		extra.Upload(c, extra.UploadMultiple(true),
			extra.UploadOnChange(func(files []extra.UploadFile) {
				listed = nil
				for _, f := range files {
					listed = append(listed, f.Name)
				}
			}))
	}
	tt := ui.NewTester(view, 400, 240)
	tt.Frame()
	if tt.FramePending() {
		t.Error("上传控件持续请求帧")
	}
	if !tt.HasText("点击选择或拖拽文件到此处") {
		t.Error("拖放区提示未渲染")
	}
	// 无窗口环境不触发系统对话框点击；列表由拖放/对话框路径驱动，
	// 这里只验证渲染与状态稳定。
	_ = listed
}

func TestExtraFollowsThemeAndDensity(t *testing.T) {
	for _, tk := range []*theme.AntToken{
		theme.Light(),
		theme.Dark(),
		theme.Light(theme.WithDensity(theme.Compact)),
		theme.Dark(theme.WithDensity(theme.Compact), theme.WithPrimary(ui.Hex("#7c3aed"))),
	} {
		token := tk
		right := []string{}
		var chosen ui.Selection[int]
		_ = chosen
		view := func(c *ui.Context) {
			theme.Apply(c, token)
			col := ui.Column(c).Gap(8)
			col.Children(func() {
				extra.Tree(c, []extra.TreeNode{{Label: "根", Children: []extra.TreeNode{{Label: "叶"}}}})
				extra.Carousel(c, []func(*ui.Context){func(c *ui.Context) { widgets.Title(c, "页", widgets.H4) }})
				extra.QRCode(c, "hello")
				extra.Upload(c)
				extra.Transfer(c, []extra.TransferItem{{Key: "a", Label: "A"}}, &right)
			})
		}
		tt := ui.NewTester(view, 560, 620)
		tt.Frame()
		if tt.FramePending() {
			t.Errorf("密度 %v 下缺口控件持续请求帧", token.Density)
		}
	}
}
