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
	if err := tt.Click("›"); err != nil {
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

func TestCascaderChoose(t *testing.T) {
	path := []string{}
	nodes := []extra.CascaderNode{
		{Label: "浙江", Children: []extra.CascaderNode{{Label: "杭州"}, {Label: "宁波"}}},
		{Label: "广东", Children: []extra.CascaderNode{{Label: "广州"}}},
	}
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		extra.Cascader(c, &path, nodes,
			extra.CascaderPlaceholder("请选择地区"),
			extra.CascaderOnChange(func(p []string) { t.Logf("选中: %v", p) }),
		)
	}
	tt := ui.NewTester(view, 400, 240)
	tt.Frame()
	if !tt.HasText("请选择地区") {
		t.Fatal("触发器占位文字未渲染")
	}
	// 打开浮层 → 点第一级 → 点第二级叶子。
	if err := tt.Click("请选择地区"); err != nil {
		t.Fatalf("点击触发器失败: %v", err)
	}
	tt.Frame()
	if !tt.HasText("浙江") {
		t.Fatal("浮层第一级未出现")
	}
	if err := tt.Click("浙江"); err != nil {
		t.Fatalf("点击第一级失败: %v", err)
	}
	tt.Frame()
	if !tt.HasText("宁波") {
		t.Fatal("浮层第二级未出现")
	}
	if err := tt.Click("杭州"); err != nil {
		t.Fatalf("点击叶子失败: %v", err)
	}
	tt.Frame()
	if len(path) != 2 || path[0] != "浙江" || path[1] != "杭州" {
		t.Errorf("path = %v, want [浙江 杭州]", path)
	}
	if !tt.HasText("浙江 / 杭州") {
		t.Error("触发器未显示选中路径")
	}
}

func TestCascaderDisabled(t *testing.T) {
	path := []string{}
	nodes := []extra.CascaderNode{{Label: "浙江"}}
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		extra.Cascader(c, &path, nodes, extra.CascaderDisabled(true))
	}
	tt := ui.NewTester(view, 400, 200)
	tt.Frame()
	if err := tt.Click("请选择"); err != nil {
		t.Fatalf("点击失败: %v", err)
	}
	tt.Frame()
	if tt.HasText("浙江") {
		t.Error("禁用后浮层仍然打开")
	}
}

func TestTreeSelectChoose(t *testing.T) {
	selected := ""
	nodes := []extra.TreeNode{
		{Label: "北京", Key: "bj"},
		{Label: "研发部", Key: "rd", Children: []extra.TreeNode{
			{Label: "前端组", Key: "fe"},
		}},
	}
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		extra.TreeSelect(c, &selected, nodes, extra.TreeSelectPlaceholder("选择部门"))
	}
	tt := ui.NewTester(view, 400, 240)
	tt.Frame()
	if err := tt.Click("选择部门"); err != nil {
		t.Fatalf("点击触发器失败: %v", err)
	}
	tt.Frame()
	if !tt.HasText("北京") {
		t.Fatal("浮层树未出现")
	}
	// 顶层叶子直接选中。
	if err := tt.Click("北京"); err != nil {
		t.Fatalf("点击节点失败: %v", err)
	}
	tt.Frame()
	if selected != "bj" {
		t.Errorf("selected = %q, want bj", selected)
	}
	if !tt.HasText("北京") {
		t.Error("触发器未显示选中节点")
	}
}

func TestCascaderTreeSelectThemes(t *testing.T) {
	for _, tk := range []*theme.AntToken{
		theme.Light(),
		theme.Dark(),
		theme.Light(theme.WithDensity(theme.Compact)),
		theme.Dark(theme.WithDensity(theme.Compact)),
	} {
		token := tk
		path := []string{}
		selected := ""
		nodes := []extra.CascaderNode{
			{Label: "浙江", Children: []extra.CascaderNode{{Label: "杭州"}}},
		}
		tree := []extra.TreeNode{{Label: "北京", Key: "bj"}}
		view := func(c *ui.Context) {
			theme.Apply(c, token)
			extra.Cascader(c, &path, nodes)
			extra.TreeSelect(c, &selected, tree, extra.TreeSelectPlaceholder("选择节点"))
		}
		tt := ui.NewTester(view, 640, 160)
		tt.Frame()
		if tt.FramePending() {
			t.Errorf("密度 %v 下持续请求帧", token.Density)
		}
		// 打开浮层并选中，验证交互在 4 主题下都正常。
		if err := tt.Click("请选择"); err != nil {
			t.Fatalf("密度 %v: 点击 Cascader 失败: %v", token.Density, err)
		}
		tt.Frame()
		if err := tt.Click("浙江"); err != nil {
			t.Fatalf("密度 %v: 点击选项失败: %v", token.Density, err)
		}
		tt.Frame()
		if err := tt.Click("杭州"); err != nil {
			t.Fatalf("密度 %v: 点击叶子失败: %v", token.Density, err)
		}
		tt.Frame()
		if len(path) != 2 {
			t.Errorf("密度 %v: path = %v", token.Density, path)
		}
		if err := tt.Click("选择节点"); err != nil {
			t.Fatalf("密度 %v: 点击 TreeSelect 失败: %v", token.Density, err)
		}
		tt.Frame()
		if err := tt.Click("北京"); err != nil {
			t.Fatalf("密度 %v: 点击树节点失败: %v", token.Density, err)
		}
		tt.Frame()
		if selected != "bj" {
			t.Errorf("密度 %v: selected = %q", token.Density, selected)
		}
	}
}
