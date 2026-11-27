# 更新日志

本文件记录 mygo-antdesign 的版本变更。格式参考 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### 计划中

- 逐控件示例与逐控件截图基线（当前为 12 个分组 / 66 个片段 + 3 张 Gallery 级基线）。
- 系统亮暗自动跟随（`theme.WithFollowSystem` 已覆盖强调色/高对比/文字缩放，亮暗跟随仍为 TODO，见 ADR 0002）。
- Gallery 全局尺寸切换（各控件尺寸选项无统一注入点）与 `ui/editor` 代码视图。

## [0.3.0] - 2026-11-27

### 新增

- **antd 缺口控件（P1）**：`Empty`、`Result`、`Skeleton`、`Statistic`、`Descriptions`、`Timeline`、`Spin`。
- **antd 缺口控件（P2）**：`Tooltip`、`Popconfirm`、`Collapse`、`Form` / `FormItem` / `FormGroup`、`List`、`Sidebar` 系列、`Splitter` 系列、`AutoComplete`、`Drawer`、`ImageEl`、`ColorPicker` / `ColorPanel`、`ProgressCircle`、`Toggle` / `ToggleGroup`、`Affix`。
- **antd 缺口控件（P3）**：`Watermark`、`Mentions`、`Tour`、`Notification`、`FloatButton`、`Anchor`、`BorderBeam`、`Masonry`。
- **原生控件补全包装**：`Hyperlink`、`DecimalInput`、`SearchField`、`Stepper`、`Calendar`、`TokenField`、`EditableText`、`Meter`、`Collapsible`、`Accordion` / `AccordionItem`、`Popover`、`ContextMenu`。
- **主题**：`theme.WithFollowSystem()`（跟随系统强调色/高对比/文字缩放，默认关闭，见 ADR 0002）。
- **Gallery**：侧栏分类分组 + 搜索 + 每组示例数、代码复制按钮、亮暗对比与密度对比、按平台生效的截图基线。
- **示例**：新增 `samples/ide`（紧凑型 IDE 布局）；`scripts/visualcheck` 可把任意示例分组渲染为 PNG 供人工核对。

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

### 变更

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
