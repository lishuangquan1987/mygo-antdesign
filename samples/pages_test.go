package sample_test

import (
	"testing"

	sample "github.com/lishuangquan1987/mygo-antdesign/samples"
)

// pagesBuilt 防止重复注册（registry 是包级全局）。
var pagesBuilt bool

// buildAll 注册全部分组（与 gallery/main.go 一致）。
func buildAll(t *testing.T) {
	t.Helper()
	if pagesBuilt {
		return
	}
	for _, d := range allDemos() {
		sample.Register(d)
	}
	pagesBuilt = true
}

// TestPagesOnePerControl 校验「每控件一页」：
//   - 页名不重复（Pages 聚合后每页唯一）；
//   - 每页至少 1 个示例片段、简介非空；
//   - 关键控件的页存在且片段数正确。
func TestPagesOnePerControl(t *testing.T) {
	buildAll(t)
	pages := sample.Pages()
	if len(pages) < 90 {
		t.Fatalf("控件页数量 = %d, want ≥ 90（checklist 共 89 个包装控件 + IDE 示例）", len(pages))
	}
	seen := map[string]bool{}
	for _, p := range pages {
		if p.Name == "" {
			t.Fatal("存在空页名")
		}
		if seen[p.Name] {
			t.Errorf("控件页重复：%s", p.Name)
		}
		seen[p.Name] = true
		if len(p.Samples) == 0 {
			t.Errorf("控件页 %s 没有示例片段", p.Name)
		}
		if p.Summary == "" {
			t.Errorf("控件页 %s 缺少简介", p.Name)
		}
		if p.Category == "" {
			t.Errorf("控件页 %s 缺少分类", p.Name)
		}
	}
	for _, name := range []string{
		"Button", "Input", "TextArea", "Select", "Switch", "Table", "Tabs",
		"Tag", "Badge", "Card", "Empty", "Result", "Skeleton", "Descriptions",
		"Timeline", "Collapse", "List", "Drawer", "Modal", "Alert", "Tooltip",
		"Pagination", "Steps", "Breadcrumb", "Sidebar", "Tree", "Upload",
		"Transfer", "Cascader", "TreeSelect", "Carousel", "QRCode", "Form",
		"ColorPicker", "Mentions",
		"Toolbar", "GridView", "Markdown", "CodeEditor", "ChartPlot",
		"RichTextEditor", "RemoteImage", "CheckboxGroup", "ColorWell",
		"Splitter", "FloatButton", "Affix", "Watermark", "Tour",
		"Notification", "Popconfirm", "ContextMenu", "Popover", "Confirm",
		"Message", "Progress", "Spin", "Segmented", "Toggle",
		"RangeSlider", "Radio", "SearchInput", "NumberInput", "Hyperlink",
		"DecimalInput", "SearchField", "Stepper", "TokenField", "EditableText",
		"Calendar", "Mentions", "Title", "Paragraph", "Code", "IDE 风格",
		"DateInput", "TimeInput", "Autocomplete", "Rating", "Sidebar",
	} {
		if !seen[name] {
			t.Errorf("缺少控件页 %s", name)
		}
	}
}

// TestPagesSampleCounts 校验多片段控件页的片段聚合。
func TestPagesSampleCounts(t *testing.T) {
	buildAll(t)
	counts := map[string]int{}
	for _, p := range sample.Pages() {
		counts[p.Name] = len(p.Samples)
	}
	// Button：基础用法/按钮类型/尺寸/状态 4 个片段。
	if counts["Button"] != 4 {
		t.Errorf("Button 片段数 = %d, want 4", counts["Button"])
	}
	// Input：基础输入框 + 尺寸与状态 2 个片段。
	if counts["Input"] != 2 {
		t.Errorf("Input 片段数 = %d, want 2", counts["Input"])
	}
	// 单片段控件页。
	for _, name := range []string{"Empty", "Result", "Meter", "Anchor", "QRCode", "Table"} {
		if counts[name] != 1 {
			t.Errorf("%s 片段数 = %d, want 1", name, counts[name])
		}
	}
	// 片段总数（去重）= 注册表片段总数。
	total := 0
	for _, d := range sample.All() {
		total += len(d.Samples)
	}
	if sample.FragmentCount() != total {
		t.Errorf("FragmentCount = %d, want %d", sample.FragmentCount(), total)
	}
}
