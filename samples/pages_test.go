package sample_test

import (
	"strings"
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
	// Table：基础表格 + 自定义单元格与分页 2 个片段。
	if counts["Table"] != 2 {
		t.Errorf("Table 片段数 = %d, want 2", counts["Table"])
	}
	// 单片段控件页。
	for _, name := range []string{"Empty", "Result", "Meter", "Anchor", "QRCode"} {
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

// TestSampleCodeExtracted 校验每个示例的 Code 都是真实源码片段（不是
// Extract 的占位提示）：防止 sample: 标记写错后页面静默显示一句注释。
func TestSampleCodeExtracted(t *testing.T) {
	buildAll(t)
	for _, p := range sample.Pages() {
		for _, s := range p.Samples {
			if s.Code == "" {
				t.Errorf("%s / %s：代码为空", p.Name, s.Title)
				continue
			}
			if strings.HasPrefix(s.Code, "// 未找到标记段") ||
				strings.HasPrefix(s.Code, "// 标记段") {
				t.Errorf("%s / %s：%s", p.Name, s.Title, strings.SplitN(s.Code, "\n", 2)[0])
			}
		}
	}
}

// TestPageOverridesValid 校验 pages.go 的旧式映射没有失效条目：
// key 的「分组名|片段标题」必须真实存在，目标页名非空。
func TestPageOverridesValid(t *testing.T) {
	buildAll(t)
	exist := map[string]bool{}
	for _, d := range sample.All() {
		for _, s := range d.Samples {
			exist[d.Name+"|"+s.Title] = true
		}
	}
	for _, k := range sample.OverrideKeys() {
		if !exist[k] {
			t.Errorf("pageOverrides 失效 key：%q（分组或片段标题已改名）", k)
		}
	}
	for _, p := range sample.Pages() {
		if p.Name == "" {
			t.Error("存在空页名")
		}
	}
}

// TestRegisterDuplicatePanics 同名控件重复注册必须 panic（防止侧栏出现
// 两个同名页，也防止两个包误用同一个 Name）。
func TestRegisterDuplicatePanics(t *testing.T) {
	buildAll(t)
	d := sample.All()[0]
	defer func() {
		if recover() == nil {
			t.Error("重复注册同名控件应 panic")
		}
	}()
	sample.Register(d)
}
