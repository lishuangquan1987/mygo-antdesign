# 更新日志

本文件记录 mygo-antdesign 的版本变更。格式参考 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

## [0.6.0] - 2026-11-28

### 新增

- **能力补全（对齐 antd v5）**：Table 自定义单元格 `TableRender` + 分页 `TablePagination` + 虚拟滚动 `TableVirtual`；Form 即时校验 `FormItemValue`/`FormItemRules`（必填/长度/正则/自定义）；`SelectMulti` 多选、`SelectSearch` 可搜索下拉；Tabs `TabsOnChange`/`TabsCard`；Pagination `Disabled`/`ShowTotal`/`PageSize`；Carousel `Autoplay`/`Arrows`/`Dots`；Upload `Extensions`/`MaxCount`；Transfer `Searchable`；`BadgeStatus` 状态点；Statistic `Formatter` 与千分位；Timeline `Pending` 待办节点；Message/Notify 类型与时长；Notification 类型图标；选择类控件统一占位 `ControlPlaceholder`。
- **Gallery 示例**：为新能力补齐 14 个示例（Tabs 卡片式页签、Pagination 选项、Select 多选与可搜索、Badge 状态徽标、Timeline 待办节点、Statistic 格式化、Form 即时校验、Notification 类型、Table 行多选与自定义单元格、Upload 类型与数量限制、Transfer 搜索、Carousel 箭头与自动播放），全部随页入库、代码即源码。
- **控件改造清单** `docs/10-checklist.md`：MyGo 原生控件（111 个构造函数，含 Base 底座结论）与 antd 缺口控件（35 个）的双清单打勾总览，全部完成；后续新增控件先加 `[ ]` 再打勾。

### 修复

- `Watermark` 把水印铺满整个窗口（含侧栏与工具条）：覆盖层改用宿主矩形，只在包裹内容内平铺。
- `FormItemRequired` 只画必填星号、不参与校验（标了必填却放行空值）：星号同时作为「必填」规则，空值即报错。
- `Pagination` 禁用态看不出禁用（选中页码仍是主色实底）：文字与选中底色一并置灰。
- `Select` 触发器右侧下拉箭头被单行文本挤出可视区：文本改 `Grow(1)`，箭头稳定贴右端。
- `UploadExtensions` 需要写「.png」才匹配：改为带点/不带点等价，均不区分大小写。
- Gallery 的 Form 示例用函数内局部变量存值，每帧重建导致输入被清空：改用 `ui.Local` 跨帧保存。
- `TableVirtual` 的 30 行门槛去掉：现在等价于 antd 的 `scroll.y`（任何行数都可限高滚动）。

### 变更

- **性能**：`icon.Icon` 增加解析缓存（mygo 的图片缓存以 `*SVG` 指针为键，此前每帧重复解析且栅格化缓存全程 miss）。
- **视觉统一**：Button loading 走 `ant.Spin`；Alert/Notification 图标改为 SVG + 矢量字形；Stepper/Statistic 趋势/搜索框/表格排序/级联箭头改用 icon 包 SVG。
- **一致性与可维护性**：`Confirm` 改用自绘 Modal（与 Modal 同一观感）；`sample.Register` 重名 panic、`All()` 返回副本、新增 `Sample.Page` 声明式页归属；Gallery 拆分为 app/sidebar/home/codepanel/compare，对比模式补 Scroll 兜底；快速开始与 README 逐字同源（测试守卫）。
- **Gallery 每控件一页**：`samples/pages.go` 把 13 组 / 84 片段聚合为 **99 个控件示例页**（页名对齐 antd 官网、MyGo 命名优先，同控件多片段合并、一片段多控件展开共享代码）；侧栏「开始」首页入口 + 每控件条目（含示例数）；新增首页（主视觉 / 特性卡片 / 统计 / 快速开始）；新增 `docs/images/home.png`。

## [0.5.0] - 2026-11-27

### 变更

- **UI 细节对齐 antd v5**：自绘 Table（表头底色/16px 内边距/行分隔线/悬停选中/排序 caret）、自绘 Tabs（32px 间距/主色下划线）、Input 双边框与焦点环修复、Progress 线形换算与环形 120px、Alert 语义图标、Card/Modal 间距与阴影、Gallery 整体重构（品牌区/搜索/分类目录/工具条/卡片/暗色截图）。详见 [docs/09](docs/09-antd-parity.md)。
- **第二批控件校准**：Pagination（32×32 页码 + 翻页箭头 + 悬停主色描边）、Steps（32px 图标/2px 连线）、Empty/Result（64/72px 插画、14/24px 文字）、FormItem（纵向标签 + 红色必填星号 + 24px 项距）、Collapse（44px 头 + `#fafafa` 底）、Checkbox（8px 间距 + 标签可点击）。
- **第三批校准**：Drawer（378px 宽/56px 头/24px 体）、List（12/16 内边距 + 悬停选中）、Descriptions（14px 标签 + 16px 内边距）、Timeline（10px 圆点/2px 连线）、Gallery 默认窗口 1280×800。
- **第四批校准**：Switch 自绘 44×22 轨道 + 18px 滑块（补鼠标悬停单测）；Tooltip 补充悬停交互测试（箭头受原生 tooltip 层限制，记录在 docs/09）。
- **第五批校准**：Tag（22 高/7 内边距/关闭按钮悬停反馈）、Badge（20×20 角标/12px 字号）、Gallery 新增「字号 小/中/大」全局切换与工具条换行；`theme.WithFontSize` 补充单测。
- **第六批校准**：Select 改为 `ui.SelectBase` 自绘（32px 高/11px 内边距/右侧箭头/32px 选项/选中浅底），保留键盘导航与原有点击交互测试。
- **第七批校准**：Drawer 打开时 220ms 从右侧滑入、Modal 打开时 200ms 缩放淡入（`AnimateTo`，尊重 ReduceMotion）；Drawer 测试补充动画帧断言与等待。
- **第八批校准（extra 缺口控件）**：Transfer（方形复选框/32px 行/40px 浅灰表头/箭头按钮）、Upload（96px 拖放区/40px 图标/16px 文案/文件行悬停）、Cascader/TreeSelect 触发器（chevron 图标 + 悬停主色描边）。
- **第九批校准**：Combobox 改为 `ui.ComboboxBase` 自绘（32px 输入框/11px 内边距/32px 候选面板），修复选中项未写回输入框的问题；AutoComplete 维持原生。
- **MyGo 全量覆盖审计**：docs/02 新增 §九 完整覆盖矩阵（111 个 `package ui` 构造函数逐项状态）与「同类型以 MyGo 命名」规则；新增 `widgets/mygo_names.go` 24+ 个 MyGo 命名入口（Breadcrumbs/Avatar/NumberInput/Progress/Rating/Spinner/RangeSlider/StepSlider/CheckboxGroup/Toolbar/GridView/SplitButton 等）；新增 `samples/native` 分组 5 个示例，示例总数 79、逐示例基线 308 张。
- **第十批：子包与 Buffer 输入适配**：`TextInputBuffer` / `TextAreaBuffer` 复用 Input/TextArea 外框；新增 `widgets/native_packages.go` 适配 `ui/markdown`（Markdown）、`ui/editor`（CodeEditor/DiffView）、`ui/chart`（ChartPlot）、`ui/richtext`（RichTextEditor）、`ui/netimage`（RemoteImage）、`ui/video`（OpenVideo）；`samples/native` 增加 Markdown 示例（共 6 个），示例总数 80、基线 312 张。
- **第十一批：子包示例补齐**：`samples/native` 增加 CodeEditor、ChartPlot、RichTextEditor、RemoteImage 示例（共 10 个）；RemoteImage 因网络依赖加入基线排除表；示例总数 84、基线 324 张。
- **第十二批：按文档计划收尾**：Gallery 新增「尺寸 小/中/大」全局切换，代码区改为 `ui/editor` 只读代码视图（语法高亮 + 行号）；主题新增 `theme.HighContrast()` 高对比预设与 `WithTextInverse`（桥接 `InverseText`）；docs/01、03、04、06、07 的状态与 TODO 同步更新。
- **README 截图更新**：`docs/images/gallery.png`（1280×800 亮色）、`gallery-dark.png`、`quickstart.png` 由 `WRITE_DOCS_SCREENSHOTS=1 go test ./gallery -run TestWriteDocsScreenshots` 重新生成。
- **Release 工作流说明**：`release` 事件使用 tag 内 workflow 版本；旧 tag 失败时可在 Actions 手动触发（main 最新 workflow）并填入该 tag 重传。

- **README 重写**为标准开源项目结构：徽章、特性、Gallery 本地构建/运行、安装、快速开始、控件覆盖表、主题与密度、开发测试、项目结构、文档索引、发布与贡献。
- **CI/Release 工作流重写**：Release 发布自动构建 5 平台 Gallery 安装包并上传附件；`workflow_dispatch` 支持填入已发布 tag 手动构建并上传到对应 Release；补充并发控制、Release 存在性校验与构建摘要。
- 新增 `Makefile`（`run`/`build`/`test`/`check`/`baselines`）与 `.gitattributes`（统一 LF、二进制标记）。
- Gallery 注册缺失的 IDE 示例分组；截图基线按 LF 源码重新生成。

### 计划中

- 系统亮暗自动跟随（MyGo `ui` 暂未暴露系统外观 API，见 ADR 0002）。
- Gallery 全局尺寸切换（各控件尺寸选项无统一注入点）与 `ui/editor` 代码视图。

## [0.3.0] - 2026-11-27

### 新增

- **antd 缺口控件（P1）**：`Empty`、`Result`、`Skeleton`、`Statistic`、`Descriptions`、`Timeline`、`Spin`。
- **antd 缺口控件（P2）**：`Tooltip`、`Popconfirm`、`Collapse`、`Form` / `FormItem` / `FormGroup`、`List`、`Sidebar` 系列、`Splitter` 系列、`AutoComplete`、`Drawer`、`ImageEl`、`ColorPicker` / `ColorPanel`、`ProgressCircle`、`Toggle` / `ToggleGroup`、`Affix`。
- **antd 缺口控件（P3）**：`Watermark`、`Mentions`、`Tour`、`Notification`、`FloatButton`、`Anchor`、`BorderBeam`、`Masonry`。
- **原生控件补全包装**：`Hyperlink`、`DecimalInput`、`SearchField`、`Stepper`、`Calendar`、`TokenField`、`EditableText`、`Meter`、`Collapsible`、`Accordion` / `AccordionItem`、`Popover`、`ContextMenu`。
- **主题**：`theme.WithFollowSystem()`（跟随系统强调色/高对比/文字缩放，默认关闭，见 ADR 0002）。
- **Gallery**：侧栏分类分组 + 搜索 + 每组示例数、代码复制按钮、亮暗对比与密度对比、按平台生效的截图基线。
- **示例**：新增 `samples/ide`（紧凑型 IDE 布局）；`scripts/visualcheck` 可把任意示例分组渲染为 PNG 供人工核对；每个控件均有示例（12 组 / 74 片段）。
- **逐示例截图基线**：`samples/golden_test.go` 按亮/暗 × 标准/紧凑 4 组合渲染 288 张逐示例 PNG（含平台 manifest），含动画的 2 个示例不纳入逐字节比对。

### 变更

- `samples/golden/` 改为 3 张稳定基线（亮标准、暗标准、亮紧凑）；并排对比含加载动画，逐帧像素不稳定，改由交互测试覆盖。
- 文档全面校对：版本、计数（106 个控件构造函数 / 12 个分组 / 66 个片段）、API 签名与已知限制。

### 修复

- 截图基线在 Linux CI 上因字体栈差异永久失败：改为带平台清单（`manifest.json`），跨平台只渲染不比对。
- `ant.Toggle` / `ant.ColorPicker` 在拉伸容器中被撑满宽度：改为内部 `Row` 包裹，恢复内容宽度。

## [0.2.0] - 2026-11-27

### 新增

- 缺口控件首批：`Tree`、`Upload`、`Transfer`、`Carousel`、`QRCode`、`Cascader`、`TreeSelect`（`extra/` 包）。
- Release 工作流：5 平台交叉编译安装包，release 后自动附加。
- 截图基线按平台生效（`samples/golden/manifest.json`）。

- README 快速开始与安装说明重写；文档补充 API 勘误表与状态列。

## [0.1.0] - 2026-11-26

### 新增

- 主题层：`AntToken` → `ui.Theme` 桥接，亮/暗、品牌色、标准/紧凑密度。
- 9 组原生控件换肤：Button、Typography、Divider、布局、输入、选择、数据展示、表格、反馈、导航。
- 44 个 antd 图标常量（`icon/`）。
- Gallery 示例应用骨架 + 全部设计文档（`docs/01`–`docs/08`、ADR）。

[Unreleased]: https://github.com/lishuangquan1987/mygo-antdesign/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/lishuangquan1987/mygo-antdesign/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/lishuangquan1987/mygo-antdesign/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/lishuangquan1987/mygo-antdesign/releases/tag/v0.1.0
