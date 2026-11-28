package sample

import "strings"

// ControlPage 每控件一页：把注册表中属于同一控件的示例片段聚合为
// 一个展示页（对齐 antd 官网「左侧控件目录、右侧该控件全部示例」的结构）。
type ControlPage struct {
	// Name 控件页名（对齐 antd 官网的英文控件名，如 "Button"）。
	Name string
	// Category 侧栏分组（沿用所在示例分组的分类）。
	Category Category
	// Package "widgets" 或 "extra"。
	Package string
	// Summary 一句话简介（pageSummary 缺省时沿用分组 Summary）。
	Summary string
	// Samples 该控件的全部示例片段（跨分组聚合，顺序保持注册顺序）。
	Samples []Sample
}

// CategoryOrder 侧栏分组的固定展示顺序（对齐 ant.design/components
// 组件总览的类别顺序：通用 → 布局 → 导航 → 数据录入 → 数据展示 → 反馈 → 其他）。
var CategoryOrder = []Category{
	CatGeneral, CatLayout, CatNav, CatInput, CatDisplay, CatFeedback, CatOther,
}

// pageCategory 控件页的分类覆写。分组（WidgetDemo）只有一个 Category，
// 但同一分组里的片段可能映射到不同分类的控件页（如布局组里的 FloatButton
// 属于「通用」、Affix 属于「其他」，反馈组里的 Tooltip/Popover/Tour 属于
// 「数据展示」）。此处按 antd 官网组件总览的归类逐页覆写。
var pageCategory = map[string]Category{
	"Divider":     CatLayout,
	"FloatButton": CatGeneral,
	"Hyperlink":   CatGeneral,
	"Anchor":      CatNav,
	"Masonry":     CatLayout,
	"Segmented":   CatDisplay,
	"Calendar":    CatDisplay,
	"Popover":     CatDisplay,
	"Tooltip":     CatDisplay,
	"Tour":        CatDisplay,
	"Tree":        CatDisplay,
	"Carousel":    CatDisplay,
	"QRCode":      CatDisplay,
	"Cascader":    CatInput,
	"TreeSelect":  CatInput,
	"Upload":      CatInput,
	"Transfer":    CatInput,
	"Rating":      CatInput,
	"Result":      CatFeedback,
	"Skeleton":    CatFeedback,
	"Confirm":     CatFeedback,
	"Affix":       CatOther,
	"BorderBeam":  CatOther,
	"ContextMenu": CatOther,
}

// pageZh 控件页的中文名（对齐 antd 官网「按钮 Button」式命名；
// 未列出的沿用英文名）。
var pageZh = map[string]string{
	"Button":                  "按钮",
	"FloatButton":             "悬浮按钮",
	"Typography":              "排版",
	"Title":                   "标题",
	"Paragraph":               "段落",
	"Code":                    "代码",
	"Divider":                 "分割线",
	"Space":                   "间距",
	"Flex":                    "弹性布局",
	"Grid":                    "栅格",
	"Splitter":                "分隔面板",
	"SplitterPercent":         "百分比分隔面板",
	"SplitterVerticalPercent": "垂直百分比分隔面板",
	"SplitterVertical":        "垂直分隔面板",
	"Anchor":                  "锚点",
	"Breadcrumb":              "面包屑",
	"Dropdown":                "下拉菜单",
	"Pagination":              "分页",
	"Steps":                   "步骤条",
	"Tabs":                    "标签页",
	"Sidebar":                 "导航菜单",
	"AutoComplete":            "自动完成",
	"Autocomplete":            "自动完成",
	"Cascader":                "级联选择",
	"Checkbox":                "多选框",
	"CheckboxGroup":           "复选框组",
	"ColorPicker":             "颜色选择器",
	"DateInput":               "日期选择框",
	"TimeInput":               "时间选择框",
	"Form":                    "表单",
	"Input":                   "输入框",
	"TextArea":                "文本域",
	"SearchInput":             "搜索框",
	"NumberInput":             "数字输入框",
	"Mentions":                "提及",
	"Radio":                   "单选框",
	"Rating":                  "评分",
	"Select":                  "选择器",
	"Slider":                  "滑动输入条",
	"RangeSlider":             "双端滑动输入条",
	"Switch":                  "开关",
	"Transfer":                "穿梭框",
	"TreeSelect":              "树选择",
	"Upload":                  "上传",
	"Avatar":                  "头像",
	"Badge":                   "徽标数",
	"Calendar":                "日历",
	"Card":                    "卡片",
	"Carousel":                "走马灯",
	"Collapse":                "折叠面板",
	"Collapsible":             "折叠块",
	"Accordion":               "手风琴",
	"Descriptions":            "描述列表",
	"Empty":                   "空状态",
	"Image":                   "图片",
	"RemoteImage":             "网络图片",
	"List":                    "列表",
	"Popover":                 "气泡卡片",
	"QRCode":                  "二维码",
	"Segmented":               "分段控制器",
	"Statistic":               "统计数值",
	"Table":                   "表格",
	"Tag":                     "标签",
	"Timeline":                "时间轴",
	"Tooltip":                 "文字提示",
	"Tour":                    "漫游式引导",
	"Tree":                    "树形控件",
	"Alert":                   "警告提示",
	"Drawer":                  "抽屉",
	"Message":                 "全局提示",
	"Modal":                   "对话框",
	"Notification":            "通知提醒框",
	"Popconfirm":              "气泡确认框",
	"Progress":                "进度条",
	"Result":                  "结果",
	"Skeleton":                "骨架屏",
	"Spin":                    "加载中",
	"Watermark":               "水印",
	"Affix":                   "固钉",
	"BorderBeam":              "边框流光",
	"Meter":                   "电平条",
	"Hyperlink":               "超链接",
	"DecimalInput":            "小数输入框",
	"SearchField":             "搜索字段",
	"Stepper":                 "步进器",
	"TokenField":              "标签输入",
	"EditableText":            "可编辑文本",
	"Toggle":                  "开关按钮",
	"ToggleGroup":             "开关按钮组",
	"GridView":                "网格视图",
	"ColorWell":               "取色块",
	"Markdown":                "文档渲染",
	"CodeEditor":              "代码编辑器",
	"ChartPlot":               "坐标图",
	"RichTextEditor":          "富文本编辑器",
	"Toolbar":                 "工具条",
	"SplitButton":             "分裂按钮",
	"ContextMenu":             "右键菜单",
	"Confirm":                 "确认对话",
	"IDE 风格":                  "IDE 布局",
}

// mygoOnly antd 没有对应类型组件的 MyGo 独有控件：组件总览卡片上
// 加「MyGo」角标（与 antd 给新组件标版本号同一形制）。
// 类型相同但命名不同的控件（如 NumberInput ↔ InputNumber、
// Rating ↔ Rate）不算独有，用 MyGo 名称即可。
var mygoOnly = map[string]bool{
	"DecimalInput":   true,
	"Meter":          true,
	"Collapsible":    true,
	"Toggle":         true,
	"ToggleGroup":    true,
	"GridView":       true,
	"ColorWell":      true,
	"Markdown":       true,
	"CodeEditor":     true,
	"ChartPlot":      true,
	"RichTextEditor": true,
	"Toolbar":        true,
	"ContextMenu":    true,
	"IDE 风格":         true,
}

// PageZh 返回控件页的中文名（未收录时返回英文名）。
func PageZh(name string) string {
	if zh, ok := pageZh[name]; ok && zh != "" {
		return zh
	}
	return name
}

// IsMyGoOnly 判断控件是否为 antd 没有对应类型的 MyGo 独有控件。
func IsMyGoOnly(name string) bool { return mygoOnly[name] }

// pageOverrides 片段 → 控件页名。key 为「分组名|片段标题」；
// 一个片段演示多个控件时展开为多页（各页共享同一 Code/View）。
// 未列出的片段归属到所在分组的 Name。
var pageOverrides = map[string][]string{
	// Typography
	// Divider（4 个片段同属一个控件）
	// 布局
	"布局|Space 间距容器": {"Space"},
	"布局|Flex 弹性布局":  {"Flex"},
	"布局|Grid 网格":    {"Grid"},
	"布局|分隔面板":       {"Splitter"},
	"布局|悬浮按钮":       {"FloatButton"},
	"布局|固钉":         {"Affix"},
	"布局|百分比与垂直分隔":   {"SplitterPercent", "SplitterVerticalPercent"},
	"布局|垂直分隔":       {"SplitterVertical"},
	// 输入框
	"输入框|基础输入框":    {"Input"},
	"输入框|多行文本":     {"TextArea"},
	"输入框|搜索框与数字输入": {"SearchInput", "NumberInput"},
	"输入框|尺寸与状态":    {"Input"},
	"输入框|日期与时间":    {"DateInput", "TimeInput"},
	"输入框|表单":       {"Form"},
	"输入框|自动补全":     {"Autocomplete"},
	"输入框|颜色选择":     {"ColorPicker"},
	"输入框|更多录入控件":   {"Hyperlink", "DecimalInput", "SearchField", "Stepper", "TokenField", "EditableText", "Calendar"},
	"输入框|提及":       {"Mentions"},
	// 选择控件
	"选择控件|复选框":    {"Checkbox"},
	"选择控件|单选组":    {"Radio"},
	"选择控件|开关":     {"Switch"},
	"选择控件|滑动条与评分": {"Slider", "Rating"},
	"选择控件|分段控制器":  {"Segmented"},
	"选择控件|下拉选择":   {"Select"},
	"选择控件|开关按钮":   {"Toggle"},
	"选择控件|双滑块范围":  {"RangeSlider"},
	"选择控件|开关按钮组":  {"ToggleGroup"},
	// 表格（片段用 Sample.Page 声明页名，不再需要 override）
	// 数据展示
	"数据展示|标签":          {"Tag"},
	"数据展示|徽标":          {"Badge"},
	"数据展示|头像":          {"Avatar"},
	"数据展示|卡片":          {"Card"},
	"数据展示|空状态与结果页":     {"Empty", "Result"},
	"数据展示|骨架屏与统计数值":    {"Skeleton", "Statistic"},
	"数据展示|描述列表与时间轴":    {"Descriptions", "Timeline"},
	"数据展示|折叠面板":        {"Collapse"},
	"数据展示|列表":          {"List"},
	"数据展示|图片":          {"Image"},
	"数据展示|折叠块与电平条":     {"Collapsible", "Accordion", "Meter"},
	"数据展示|锚点、瀑布流与边框流光": {"Anchor", "Masonry", "BorderBeam"},
	// 导航
	"导航|面包屑与页签": {"Breadcrumb", "Tabs"},
	"导航|下拉菜单":   {"Dropdown"},
	"导航|分页器":    {"Pagination"},
	"导航|步骤条":    {"Steps"},
	"导航|侧栏导航":   {"Sidebar"},
	"导航|带图标的侧栏": {"Sidebar"},
	// 反馈
	"反馈|提示条":       {"Alert"},
	"反馈|模态对话框":     {"Modal"},
	"反馈|轻提示与通知":    {"Message"},
	"反馈|进度条":       {"Progress"},
	"反馈|加载中":       {"Spin"},
	"反馈|文字提示与气泡确认": {"Tooltip", "Popconfirm"},
	"反馈|抽屉":        {"Drawer"},
	"反馈|气泡与右键菜单":   {"Popover", "ContextMenu"},
	"反馈|水印、引导与通知":  {"Watermark", "Tour", "Notification"},
	"反馈|确认对话":      {"Confirm"},
	// 缺口控件
	"缺口控件|树形控件":    {"Tree"},
	"缺口控件|级联选择":    {"Cascader"},
	"缺口控件|树形下拉选择":  {"TreeSelect"},
	"缺口控件|上传":      {"Upload"},
	"缺口控件|穿梭框":     {"Transfer"},
	"缺口控件|走马灯与二维码": {"Carousel", "QRCode"},
	// 原生补全
	"原生补全|Toolbar 工具条":        {"Toolbar"},
	"原生补全|SplitButton 分裂按钮":   {"SplitButton"},
	"原生补全|GridView 网格视图":      {"GridView"},
	"原生补全|CheckboxGroup 复选组":  {"CheckboxGroup"},
	"原生补全|ColorWell 取色块":      {"ColorWell"},
	"原生补全|Markdown 文档":        {"Markdown"},
	"原生补全|CodeEditor 代码编辑器":   {"CodeEditor"},
	"原生补全|ChartPlot 坐标图":      {"ChartPlot"},
	"原生补全|RichTextEditor 富文本": {"RichTextEditor"},
	"原生补全|RemoteImage 网络图片":   {"RemoteImage"},
}

// pageSummary 控件页的一句话简介（对齐 antd 官网口吻；未列出的沿用分组 Summary）。
var pageSummary = map[string]string{
	// 通用
	"Title":                   "提供 H1–H5 五级标题，字重与行高对齐 antd Typography。",
	"Paragraph":               "段落与文本组件：次要文字、删除线、加粗、下划线、标记。",
	"Code":                    "行内代码与代码块样式，等宽字体展示。",
	"Button":                  "按钮用于开始一个即时操作。五种类型 × danger × 三档尺寸，支持 Loading 与 Block。",
	"Divider":                 "分割线用于区隔内容，支持水平/垂直、实线/虚线与带文字。",
	"Space":                   "给一组子元素统一加间距，支持三档尺寸与纵向排列。",
	"Flex":                    "弹性布局：主轴/交叉轴对齐、间距与换行。",
	"Grid":                    "等宽网格布局，间距随密度缩放。",
	"Splitter":                "分隔面板：左右两栏，中间分隔条可拖拽调整宽度。",
	"SplitterPercent":         "按百分比左右分栏，两栏随窗口等比伸缩。",
	"SplitterVerticalPercent": "按百分比上下分栏。",
	"SplitterVertical":        "上下分栏，尺寸由调用方持有。",
	"FloatButton":             "固定在窗口右下角的悬浮按钮，用于全局动作。",
	"Affix":                   "固钉：把内容固定在滚动容器顶部。",
	"IDE 风格":                  "紧凑密度下的 IDE 布局示例：侧栏 + 页签 + 代码区 + 状态栏。",
	// 数据录入
	"Input":        "输入框：前缀/后缀、密码、字数限制、清除与三档尺寸。",
	"TextArea":     "多行文本输入，可设行数与最大长度。",
	"SearchInput":  "搜索框：回车提交触发搜索。",
	"NumberInput":  "数字输入框：范围与步进约束。",
	"DateInput":    "日期选择框，受控 *time.Time（antd DatePicker）。",
	"TimeInput":    "时间选择框，12/24 小时制随区域（antd TimePicker）。",
	"Form":         "表单：纵向标签 + 必填校验 + 帮助/错误文字。",
	"Autocomplete": "自动补全：保留键入内容并给出候选。",
	"ColorPicker":  "颜色选择器：色块展开面板 / 常驻面板。",
	"Mentions":     "提及：输入 @ 触发候选并插入。",
	"Hyperlink":    "超链接，主色文字。",
	"DecimalInput": "金额输入：字符串承载，避免浮点舍入。",
	"SearchField":  "原生搜索框：放大镜、清除与 Esc。",
	"Stepper":      "步进器：上下箭头调整数值。",
	"TokenField":   "标签输入：输入即成标签，附候选建议。",
	"EditableText": "原位编辑文本，点击进入编辑态。",
	"Calendar":     "日历：月历浏览与选择。",
	"Checkbox":     "复选框：8px 标签间距，标签文字可点击切换。",
	"Radio":        "单选组：互斥选择。",
	"Switch":       "开关：44×22 轨道 + 18px 滑块，选中主色。",
	"Slider":       "滑动条：单端取值，主题化轨道与滑块。",
	"Rating":       "评分：星级选择，选中色对齐 antd。",
	"Segmented":    "分段控制器：一组互斥选项。",
	"Select":       "下拉选择：自绘触发器与 32px 选项面板，支持键盘导航。",
	"Toggle":       "开关按钮：保持按下的状态按钮。",
	"RangeSlider":  "双端滑动条：区间取值。",
	"ToggleGroup":  "开关按钮组：一组互斥的状态按钮。",
	// 数据展示
	"Tag":          "标签：7 色系语义色 + 可关闭。",
	"Badge":        "徽标：计数与红点，右上角骑角展示；另有状态徽标。",
	"Avatar":       "头像：姓名首字回退，支持图片。",
	"Card":         "卡片：头部 extra 与底部 actions，8px 圆角轻阴影。",
	"Empty":        "空状态：默认插画 + 描述 + 操作区。",
	"Result":       "结果页：success/info/warning/error 四态。",
	"Skeleton":     "骨架屏：内容加载前的占位。",
	"Statistic":    "统计数值：大字号 + 涨跌趋势色。",
	"Descriptions": "描述列表：键值成对展示，支持多列与外框。",
	"Timeline":     "时间轴：节点色（蓝/绿/红/灰）+ 标题描述，支持末尾待办节点。",
	"Collapse":     "折叠面板：头部点击展开/收起，支持手风琴。",
	"List":         "列表：虚拟滚动 + 行分隔线 + 选中/多选。",
	"Image":        "图片：可指定宽度，点击放大预览。",
	"Collapsible":  "折叠块：单项展开/收起，带原生动画。",
	"Accordion":    "手风琴：互斥展开的分组面板。",
	"Meter":        "电平条：按阈值变色（正常/警告/危险）。",
	"Anchor":       "锚点目录：高亮当前项。",
	"Masonry":      "瀑布流：按列均衡分配。",
	"BorderBeam":   "边框流光：主题色渐变描边。",
	"Table":        "表格：列定义、表头排序（数值列按数值）、行多选、限高滚动。",
	// 导航
	"Breadcrumb": "面包屑：指示当前页面在层级结构中的位置。",
	"Tabs":       "页签：32px 间距 + 主色下划线，支持卡片式（card）样式。",
	"Dropdown":   "下拉菜单：声明式菜单项、分隔线与危险项。",
	"Pagination": "分页器：32×32 页码、翻页箭头、窗口折叠、总数文案与每页条数。",
	"Steps":      "步骤条：完成/进行/等待三态。",
	"Sidebar":    "侧栏导航：分区 + 条目 + 选中态，键盘可导航。",
	// 反馈
	"Alert":        "提示条：info/success/warning/error 四级别，可关闭。",
	"Modal":        "对话框：打开时缩放淡入，头部/内容/底部操作区。",
	"Message":      "轻提示：底部短暂出现的全局消息。",
	"Progress":     "进度条：线形 8px 条高 + 百分比文字。",
	"Spin":         "加载中：指示器 + 提示文字，可容器内居中。",
	"Tooltip":      "文字提示：悬停 0.6s 显示反色气泡。",
	"Popconfirm":   "气泡确认：点击触发器弹出确认/取消。",
	"Drawer":       "抽屉：右侧滑入面板，带遮罩与标题栏。",
	"Popover":      "气泡卡片：点击锚点弹出任意内容。",
	"ContextMenu":  "右键菜单：右键/菜单键触发。",
	"Watermark":    "水印：半透明文字平铺，不拦截输入。",
	"Tour":         "引导：右下角步骤卡片漫游。",
	"Notification": "通知：右上角卡片，标题 + 描述 + 关闭。",
	"Confirm":      "确认对话：基于系统 Alert 的确认回调。",
	// 缺口控件
	"Tree":       "树形控件：声明式节点树 + 展开状态。",
	"Cascader":   "级联选择：逐级下钻，叶子写回路径。",
	"TreeSelect": "树形下拉：下拉面板内嵌树。",
	"Upload":     "上传：拖放区 + 系统文件对话框 + 可移除列表。",
	"Transfer":   "穿梭框：双栏 + 方形复选框 + 箭头移动。",
	"Carousel":   "走马灯：指示点切页。",
	"QRCode":     "二维码：码色随主题（亮黑码白底/暗白码深底）。",
	// 原生补全
	"Toolbar":        "工具条：控件无独立面容，悬停高亮，支持溢出收纳。",
	"SplitButton":    "分裂按钮：主操作 + 箭头菜单（Dropdown.Button）。",
	"GridView":       "网格视图：等宽网格 + 虚拟滚动。",
	"CheckboxGroup":  "复选组：总复选框控制一组复选框。",
	"ColorWell":      "取色块：点击展开系统取色面板。",
	"Markdown":       "Markdown 渲染：标题、列表、行内代码与链接。",
	"CodeEditor":     "代码编辑器：语法高亮、行号、只读浏览。",
	"ChartPlot":      "坐标图：轴刻度 + 绘图区，自绘数据标记。",
	"RichTextEditor": "富文本编辑器：段落、加粗与列表。",
	"RemoteImage":    "网络图片：占位图加载完成后淡入。",
}

// pageKey 生成片段在 override 表里的 key。
func pageKey(group, title string) string {
	return group + "|" + title
}

// OverrideKeys 返回 pageOverrides 的全部 key（测试用于发现失效映射）。
func OverrideKeys() []string {
	out := make([]string, 0, len(pageOverrides))
	for k := range pageOverrides {
		out = append(out, k)
	}
	return out
}

// Pages 把注册表展开为「每控件一页」：同控件的片段聚合为一页，
// 一个片段演示多个控件时按 override 展开为多页（共享同一 Code/View）。
// 顺序按 CategoryOrder 与注册顺序；同页名的片段合并。
func Pages() []ControlPage {
	index := map[string]int{}
	out := make([]ControlPage, 0, 128)
	add := func(p ControlPage) {
		if i, ok := index[p.Name]; ok {
			out[i].Samples = append(out[i].Samples, p.Samples...)
			return
		}
		index[p.Name] = len(out)
		out = append(out, p)
	}
	for _, cat := range CategoryOrder {
		for _, d := range registry {
			if d.Category != cat {
				continue
			}
			for _, s := range d.Samples {
				// 目标页优先取 Sample.Page（声明式，推荐）；
				// 旧式多控件片段仍可用 pageOverrides 按「分组|标题」映射。
				names := []string{d.Name}
				if s.Page != "" {
					names = []string{s.Page}
				} else if over, ok := pageOverrides[pageKey(d.Name, s.Title)]; ok && len(over) > 0 {
					names = over
				}
				for _, name := range names {
					summary := d.Summary
					if s, ok := pageSummary[name]; ok {
						summary = s
					}
					category := d.Category
					if cat, ok := pageCategory[name]; ok {
						category = cat
					}
					add(ControlPage{
						Name:     name,
						Category: category,
						Package:  d.Package,
						Summary:  summary,
						Samples:  []Sample{s},
					})
				}
			}
		}
	}
	return out
}

// FragmentCount 返回去重后的示例片段总数（一个片段演示多个控件只计一次）。
func FragmentCount() int {
	n := 0
	for _, d := range registry {
		n += len(d.Samples)
	}
	return n
}

// MatchPage 判断控件页是否命中搜索词（控件名/中文名/简介/分类，大小写不敏感）。
func MatchPage(p ControlPage, query string) bool {
	q := strings.TrimSpace(strings.ToLower(query))
	if q == "" {
		return true
	}
	return strings.Contains(strings.ToLower(p.Name), q) ||
		strings.Contains(string(p.Category), q) ||
		strings.Contains(PageZh(p.Name), q) ||
		strings.Contains(strings.ToLower(p.Summary), q)
}
