# 07 · 控件示例（Samples）与 Gallery 设计

> 对应需求：为**每一个控件**写示例，告诉别人怎么用，并展示不同主题、不同 size 的效果——像 WPF 的 MaterialDesignInXamlToolkit 示例程序那样，让人一眼看到每个控件的全部样子。
>
> 本文区分「已实现」与「规划」。当前实现：13 个示例分组、80 个示例片段、Gallery 侧栏（分类 + 搜索 + 计数）+ 卡片流（代码复制）+ 控制条 + 3 张 Gallery 级基线 + 312 张逐示例基线（4 组合）。

## 1. 目标与体验参照

WPF 的 Material Design 示例程序之所以经典，是因为它做到了三件事，本库的 Gallery 也要做到：

1. **看得见**：每个控件都有独立展示区，列出它的所有变体、状态、尺寸，不用读文档就知道长什么样。
2. **学得会**：每个示例旁边就是可复制的真实代码，效果与代码一一对应。
3. **验得了**：示例本身就是视觉回归的载体，主题/尺寸一改，全部示例同步变化。

三个硬指标：

- **示例覆盖率**：`widgets/` 与 `extra/` 中每个控件至少 1 个示例页；高频控件覆盖全部变体。**当前按分类成组（一个分组覆盖多个控件），逐控件拆分是待办**。
- **四维展示**：主题（亮/暗/品牌色）× 密度（标准/紧凑）× 尺寸（大/中/小）× 状态（默认/禁用/加载/选中/错误）。**当前：主题与密度由 Gallery 全局切换；尺寸与状态写在示例内部，无全局尺寸开关**。
- **代码即示例**：展示的代码片段来自真实源码，不是手写的「示意代码」。

## 2. 为什么 MyGo 能把这件事做得更好

| 能力 | WPF / Web 做法 | MyGo 做法 |
| --- | --- | --- |
| 主题切换 | 切换全局资源字典 / CSS 变量 | `c.SetTheme(&t)`，下一帧全变 |
| **同屏对比亮/暗** | 需两个窗口或复杂逻辑 | `theme.Sub(c, t)`（内部 `ui.Themed`）——同一屏并排渲染亮色与暗色，弹层/Tooltip 自动跟随 |
| 截图回归 | 额外截图工具 | `ui.NewTester` 无窗口渲染 PNG，同一套示例直接出基线 |
| 代码展示 | 静态文本 | 真实源码 `go:embed` + 标记提取（`ui/editor` 只读编辑器可作为后续增强） |
| 即时刷新 | 手动触发渲染 | view = f(state)，切换器改 state 即可 |

**杀手锏是 `ui.Themed`**：可以在一个示例卡片里左右并排渲染同一控件的亮色版与暗色版（或标准与紧凑），这是传统框架需要额外窗口才能做到的。

## 3. 示例数据结构（`samples/registry.go`，已实现）

```go
// Sample 一个示例单元：标题 + 说明 + 可复制的真实代码 + 实时预览。
type Sample struct {
	Title       string
	Description string
	Code        string                // 由 go:embed + // sample: 标记提取
	View        func(c *ui.Context)   // 实时预览（无返回值，直接构建控件）
}

// WidgetDemo 一个控件的示例集合（或一组控件的示例分组）。
type WidgetDemo struct {
	Name     string   // "Button"
	Package  string   // "widgets" | "extra"
	Category Category // CatGeneral / CatInput / CatDisplay / CatNav / CatFeedback / CatExtra
	Summary  string
	Samples  []Sample
}

// 显式注册（Gallery 启动时），不做 init 自动发现。
func Register(d WidgetDemo)
func All() []WidgetDemo
```

> 与早期草案的差异：`View` 不返回 `ui.Element`（即时模式直接构建）；`WidgetDemo` 没有 `Variants` 字段；注册是 `gallery/main.go` 里显式调用 `samples.Register(button.Demo())`，不是每个控件包 `init()` 自动注册。`Category` 已在数据里定义，但 Gallery 侧栏目前是扁平列表，尚未按分类分组。

## 4. 示例覆盖矩阵（Sample DoD）与现状

一个控件的示例页应覆盖下表的行 × 列，至少做到「基础 + 主要变体 + 禁用态 + 三档尺寸」：

| 维度 | 必须展示的内容 | 现状 |
| --- | --- | --- |
| **用法** | 基础用法 → 进阶用法（事件/绑定）→ 组合用法 | ✅ 多数示例有 |
| **变体** | 控件自身的 type/variant（如 Button：primary/default/dashed/text/link/danger） | ✅ 高频控件覆盖 |
| **状态** | default / hover / pressed / **disabled** / **loading** / selected / error | ✅ 固定态展示 |
| **尺寸** | large / middle / small（`ant.Size(...)`） | ✅ 示例内展示 |
| **密度** | 标准 / 紧凑 | ✅ Gallery 全局切换 |
| **主题** | 亮 / 暗 / 自定义品牌色 | ✅ 顶部切换器 + 并排对比 |

以 Button 为例，当前示例页（`samples/button/demo.go`）包含：

- 基础按钮（点击计数，演示 `Clicked()` / `ant.OnClick`）
- 五种 type 变体并排
- 三种 size 并排
- 危险、禁用、加载态
- 并排对比由 Gallery 顶部「亮暗对比」开关提供（不是 Button 示例独有）

尚未覆盖：文字 + 图标按钮、与 Popconfirm 组合的危险操作确认、标准/紧凑并排（当前只做了亮/暗并排）。

### 4.1 当前示例分组（13 组 / 80 个片段）

| 分组 | 覆盖控件 | 示例数 |
| --- | --- | --- |
| `samples/button` | Button | 4 |
| `samples/typography` | Title / Paragraph / Code / TextMuted | 3 |
| `samples/divider` | Divider / VerticalDivider / DashedVerticalDivider | 4 |
| `samples/layout` | Space / Flex / Grid / Splitter / FloatButton / Affix / SplitterVertical | 8 |
| `samples/input` | Input / TextArea / SearchInput / NumberInputV / DatePicker / TimePicker / Form / AutoComplete / ColorPicker / Hyperlink·DecimalInput·SearchField·Stepper·TokenField·EditableText·Calendar / Mentions | 10 |
| `samples/selectdemo` | Checkbox / RadioGroup / Switch / Slider / RangeSlider / Rate / Segmented / Select / Combobox / Toggle / RangeSliderDouble / ToggleGroup | 9 |
| `samples/display` | Tag / Badge / Avatar / Card / Empty / Result / Skeleton / Statistic / Descriptions / Timeline / Collapse / List / Image / Collapsible·Accordion·Meter / Anchor·Masonry·BorderBeam | 12 |
| `samples/table` | Table | 1 |
| `samples/feedback` | Alert / Modal / Message / Notify / Progress / Confirm / Spin / Tooltip / Popconfirm / Drawer / Popover / ContextMenu / Watermark·Tour·Notification / 确认对话 | 10 |
| `samples/nav` | Breadcrumb / Tabs / Dropdown / Pagination / Steps / Sidebar / SidebarItemWithIcon | 6 |
| `samples/extrademo` | Tree / Upload / Transfer / Carousel / QRCode / Cascader / TreeSelect | 6 |
| `samples/native` | 原生补全：Toolbar / SplitButton / GridView / CheckboxGroup / ColorWell / Markdown | 6 |
| `samples/ide` | 紧凑型 IDE 布局（Sidebar / Tabs / 代码区 / 状态栏） | 1 |

## 5. Gallery 应用的界面结构

### 5.1 当前实现（`gallery/app.go`）

```
┌────────────────────────────────────────────────────────────────┐
│ TopBar: [亮|暗] [标准|紧凑] [品牌色■ ■ ■ ■] [亮暗对比]             │
├────────────┬───────────────────────────────────────────────────┤
│ 侧栏        │ 内容区（当前示例分组，可滚动）                       │
│  ▸ Button   │ Button                                            │
│  ▸ Typography│ 按钮：五种 type × danger × 三档尺寸…              │
│  ▸ Divider  │ ┌───────────────────────────────────────────────┐ │
│  ▸ Layout   │ │ 基础用法                                       │ │
│  ▸ Input    │ │ 说明文字…                                      │ │
│  ...        │ │ [保存]  ← 实时控件                             │ │
│  ▸ 缺口控件  │ │ // 真实源码片段（只读文本，最多 30 Space 高）     │ │
│            │ └───────────────────────────────────────────────┘ │
│            │ ... 其余示例卡片 ...                              │
└────────────┴───────────────────────────────────────────────────┘
```

- **左侧**：搜索框 + 按分类分组的控件列表，每个分组显示示例数量。
- **中部**：示例卡片流，每张卡片 = 标题 + 描述 + 实时预览 + 代码文本（`MaxHeight` 截断），右上角「复制代码」。
- **顶部控制条**：亮/暗、标准/紧凑、4 个品牌色块、亮暗对比与密度对比开关。
- **并排对比模式**：`theme.Sub` 左右两栏——亮暗对比（左亮右暗）或密度对比（左标准右紧凑），携带当前品牌色。

控制条只改 `Gallery` 的 app state，view 函数是 state 的函数，下一帧全部示例自动重绘——**不需要逐个刷新示例**。

### 5.2 实现状态

| 能力 | 状态 | 说明 |
| --- | --- | --- |
| 侧栏按 `Category` 分组 | ✅ | 6 个分类（通用/输入/展示/导航/反馈/缺口） |
| 侧栏搜索 / 过滤 | ✅ | 匹配控件名、简介与分类名 |
| 覆盖率显示（每组示例数） | ✅ | 行尾显示该分组示例数量 |
| 全局尺寸切换（大/中/小） | ◻ | 尺寸由各示例内部指定，无统一注入点 |
| 标准/紧凑并排 | ✅ | 顶部「密度对比」 |
| 亮/暗并排 | ✅ | 顶部「亮暗对比」 |
| 代码复制按钮 | ✅ | `c.WriteClipboard` |
| 代码高亮（`ui/editor` 只读） | ◻ | 当前用 `ui.Text` 等宽展示 |
| 逐示例截图 | ✅ | 312 张（4 组合 × 78 示例，`samples/golden/controls/`，含平台 manifest） |

### 5.3 本地构建与运行

```bash
# 直接运行（开发预览；需要桌面环境）
go run ./gallery

# 构建可执行文件（版本号通过 ldflags 注入，会显示在启动日志中）
go build -trimpath -ldflags "-X main.version=v0.3.0" -o bin/gallery ./gallery
./bin/gallery                 # Linux / macOS
bin\gallery.exe               # Windows（go build -o bin/gallery.exe ...）

# 可选：Make 封装
make run                      # go run ./gallery
make build                    # 输出 bin/gallery

# 无窗口渲染示例截图（人工核对控件效果，不依赖 computer use）
go run ./scripts/visualcheck -group 反馈 -theme dark -out .visualcheck

# 重新生成截图基线（Gallery 级 + 逐示例）
go test ./gallery ./samples -count=1 -run 'TestGalleryBaselines|TestControlScreenshots'
```

前置：Go 1.27.4+、桌面环境（Windows 10+/macOS 11+/X11 或 Wayland）；无需 cgo。发布时 CI 会为 5 个平台构建 Gallery 安装包并附加到对应 Release（见 `.github/workflows/release.yml`）。

## 6. 代码片段如何保证「真实不过时」

| 方案 | 做法 | 评价 |
| --- | --- | --- |
| A 手写字符串 | 在 Sample 里写一段代码文本 | 简单，但代码与实际实现会漂移 ❌ |
| **B `go:embed` 真实源码（已采用）** | 示例代码写在 `samples/<分组>/demo.go` 的真实 `.go` 文件里，用 `// sample:basic … // /sample` 标记片段，`go:embed` 读取后由 `samples.Extract` 按标记提取 | 示例即源码，永不漂移 ✅ |
| C AST 提取 | 用 `go/ast` 提取函数体 | 精确但复杂，维护成本高 |

实际写法：

```go
//go:embed demo.go
var demoSrc string

func Demo() samples.WidgetDemo {
	return samples.WidgetDemo{
		Name: "Button", Package: "widgets", Category: samples.CatGeneral,
		Samples: []samples.Sample{{
			Title: "基础用法",
			Code:  samples.Extract(demoSrc, "basic"), // // sample:basic … // /sample
			View:  Basic,
		}},
	}
}

// sample:basic
func Basic(c *ui.Context) { /* 真实可编译示例 */ }
// /sample
```

`Extract` 在缺少开始/结束标记或片段为空时返回注释提示而不是空串，便于发现遗漏。

## 7. 目录结构与注册表（实际）

```
samples/
├── registry.go              # WidgetDemo / Sample / Category / Register / All
├── extract.go               # // sample: 标记提取
├── button/demo.go           # 每个分组一个 demo.go（go:embed 自身）
├── typography/  divider/  layout/  input/  selectdemo/  ide/
├── display/  table/  feedback/  nav/  extrademo/
└── golden/                  # 截图基线（真实生成，非占位）
    ├── gallery-light-standard.png
    ├── gallery-dark-standard.png
    ├── gallery-light-compact.png
    └── gallery-dark-compact-compare.png
gallery/
├── main.go                  # 入口 + 显式注册 11 个 Demo()
├── app.go                   # 侧栏 + 卡片流 + 控制条 + 并排对比
└── app_test.go              # 渲染/交互测试 + 基线生成与比对
scripts/
└── visualcheck/main.go      # 把某个示例分组的每个示例按主题渲染成 PNG（人工核对）
```

**人工核对控件效果（不依赖 computer use）**：

```bash
# 渲染「数据展示」分组在亮色下的每个示例，输出到 .visualcheck/
go run ./scripts/visualcheck -group 数据展示 -theme light
# 暗色 / 紧凑密度
go run ./scripts/visualcheck -group 数据展示 -theme dark
go run ./scripts/visualcheck -group 数据展示 -theme compact
```

`visualcheck` 直接复用 Gallery 注册的示例 `View`，用 `ui.NewTester` 无窗口
渲染 PNG，所以看到的与 Gallery 中一致；配合 `samples/all_test.go`（校验全部
示例的代码片段提取与亮/暗渲染）可覆盖「效果 + 代码」两方面。

**新增示例的步骤**：

1. 在 `samples/<分组>/demo.go` 写示例函数，用 `// sample:<name>` … `// /sample` 包住需要展示的代码。
2. 在 `Demo()` 里追加 `samples.Sample{Title, Description, Code: samples.Extract(demoSrc, "<name>"), View: fn}`。
3. 若是新分组，在 `gallery/main.go` 里 `samples.Register(<pkg>.Demo())`。
4. 跑 `go test ./gallery -run TestGalleryBaselines`：基线不存在时生成、存在时比对；预期视觉变化则删除旧 PNG 重新生成。
5. 在 README / 本文档的分组表里补一行。

## 8. 截图生成与文档联动（实际）

1. `gallery/app_test.go` 的 `TestGalleryBaselines` 用 `ui.NewTester` 渲染 **3 种组合**：
   `light-standard`、`dark-standard`、`light-compact`、`dark-compact-compare`（960×640 / 960×900）。
2. 输出到 `samples/golden/gallery-<组合>.png`，与基线**逐字节**比对；不一致即失败（预期改动需删文件重建）。
3. **跨平台策略**：`samples/golden/manifest.json` 记录基线生成平台（`goos`/`goarch`）。在其它平台运行时只渲染、断言无空帧，然后跳过逐字节比对——不同系统的字体栈不同，像素必然不同。CI（Linux）因此不会因 Windows 基线而永久红。
4. 文档引用真实截图：`docs/images/gallery.png`（Gallery 首图）、`docs/images/quickstart.png`（README 最小示例）。
5. CI 中跑 `go test ./gallery`，示例或主题一变，截图 diff 立即暴露（仅基线平台）。

> 早期文档写的 `docs/screenshots/<widget>-<variant>-<theme>-<density>.png` 与「每个示例 × 组合」**尚未实现**；当前只有 Gallery 级 3 张基线，且基线带平台属性。并排对比模式包含加载动画（`ui.Spinner`），逐帧像素不稳定，由交互测试覆盖。彻底解决需要固定字体（`ui.RegisterFont`）或按平台提交多套基线。
>
> 重建基线：删除 `samples/golden/*.png` 与 `manifest.json` 后跑 `go test ./gallery -run TestGalleryBaselines`（会在当前平台生成新基线）。

## 9. 编写规范

1. 每个示例必须有 Title、Description、Code（真实源码）、View（实时预览）。
2. 代码必须可独立复制运行，不依赖示例框架内部状态。
3. 覆盖第 4 节矩阵（至少基础 + 主要变体 + 禁用 + 三档尺寸）。
4. 严守样式铁律（见 `AGENTS.md` 2.1）：示例本身也只读主题、只用 `t.Space` / `t.Rem`。
5. 在注册表登记，Gallery 侧栏可见。
6. 生成 3 组合截图并通过 diff。
7. 示例中的交互要有明确反馈（如点击计数变化），让人看出控件是「活的」。

## 10. 分阶段交付（状态）

| 阶段 | 交付 | 状态 |
| --- | --- | --- |
| M1 | Gallery 框架（侧栏/控制条/并排对比/代码视图）+ 基础控件示例 | ✅ |
| M2 | 随控件开发同步补齐全部换肤控件的示例 | ◐ 按分组覆盖，未逐控件 |
| M3 | 全控件示例的四维覆盖审计 + 紧凑并排 + 截图基线全量生成 | ◻ |
| M4 | 缺口控件示例 + 文档嵌图 + 发布 | ◐ 缺口示例已有 6 个（extrademo），文档嵌图仅 README |

## 11. 工作量补充估算

| 项 | 估算（人日） | 状态 |
| --- | --- | --- |
| Gallery 框架（侧栏、控制条、并排对比、代码视图、搜索） | 6 | ✅ 除搜索 |
| 示例注册表 + `go:embed` 片段提取工具 | 2 | ✅ |
| 全控件示例编写（约 70 个控件 × 0.3） | 21 | ◐ 13 组 / 80 片段（按 docs/02 §九 继续补齐） |
| 截图生成与文档嵌图联动 | 4 | ◐ Gallery 级 3 张 + README 2 图 |
| **小计** | **约 33 人日** | |

> 这部分工作并入 `06-development-plan.md` 的 WS-F，并计入每个控件的 DoD。示例不必等控件全部完成才写——**每个控件完成时立即配示例**，避免后期补写遗漏。
