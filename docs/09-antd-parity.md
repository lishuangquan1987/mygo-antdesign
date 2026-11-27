# 09 · 与 Ant Design v5 的细节对齐（UI Polish 审计）

> 目标：控件不只是「看起来像」，而是把 antd v5 的尺寸、间距、圆角、状态色、
> 层级与交互细节逐项校准；Gallery 作为展示窗口也要达到「可对外演示」的完成度。
>
> 本文件记录本轮打磨的范围、量化依据与仍存在的差异。所有视觉改动都通过
> `ui.NewTester` 无窗口渲染截图核对（`samples/golden/` 与 `docs/images/`），
> 不使用 computer use。

## 1. 基准 Token（antd v5 默认）

| Token | antd v5 默认 | 本库实现 | 位置 |
| --- | --- | --- | --- |
| `colorPrimary` | `#1677ff` | 同 | `theme/` |
| `borderRadius` | 6 | 6（紧凑 4） | `theme.WithDensity` |
| `controlHeight` / `SM` / `LG` | 32 / 24 / 40 | `Space(8)` / `Space(6)` / `Space(10)` | `internal/metrics.go` |
| `fontSize` | 14 | 14（紧凑 13） | `theme` |
| `paddingContentHorizontal` | 16 | `Space(4)` | `internal.ControlPaddingX` |
| `colorBorder` | `#d9d9d9` | `Theme.Border` | `theme` |
| `colorBgContainer` | `#ffffff` | `Theme.Surface` | `theme` |
| `colorFillAlter`（表头） | `#fafafa` | `Background.Mix(Surface, 0.5)` | `widgets/table.go` |
| 焦点环 | 1px 主色边 + 2px `rgba(5,145,255,.1)` 阴影 | 同（`Shadow(0,0,0,2,Accent·12%)`） | `widgets/input.go` |

## 2. 本轮已校准的细节

| 控件 | 问题（打磨前） | 处理 |
| --- | --- | --- |
| Table | 直接透传原生表格：无外框/表头底色/单元格内边距，行高仅 24px | 自绘 antd 表格：8px 圆角外框、`#fafafa` 表头 + 600 字重、16px 单元格内边距、1px 行分隔线、行悬停/选中底色、排序 caret 与 Cmd/Shift 多选 |
| Pagination | 20px 小方块、无边框、无翻页箭头 | 32×32 页码、6px 圆角、选中主色底白字、悬停主色描边、首尾 `‹`/`›` 箭头 |
| Steps | 24px 图标、1px 连线、13px 标签 | 32px 图标、2px 连线、14px 标签、左右 16px 留白避免末项贴边 |
| Empty | 描述 12.6px、元素间距 12px | 64px 插画、14px 描述、8px 间距 |
| Result | 56px 图标、20px 标题 | 72px 图标、24px 标题、14px 副标题、24px 操作区间距 |
| FormItem | 透传原生 `ui.Field`（必填星号为纯文本、间距不可控） | 纵向标签 + 红色 `*`、24px 项距、12px 帮助/错误文字 |
| Collapse | 32px 头部、13px 标题、白底 | 44px 头部、`#fafafa` 底、14px/500 标题、16px 内容内边距、8px 圆角外框 |
| Checkbox | 原生控件内部间距 4px | 自绘一行：8px 间距、标签文字可点击切换 |
| Input / TextArea / NumberInput | 内层原语自带边框+内边距，与外层容器叠加出现双边框与右端「圆角帽」；聚焦为 2px 实心边（会跳动） | 清除内层边框/内边距/圆角；聚焦改为 1px 主色边 + 2px 半透明焦点环 |
| Tabs | 透传原生页签，间距/字号/下划线不可控 | 自绘：32px 页签间距、选中主色 + 2px 下划线、悬停主色、整行 1px 分隔线 |
| Progress（线形） | 直接把 0–100 传给 0–1 的原语，68% 显示近似满格；条高 4px | 换算为 0–1；8px 条高；百分比文本 14px 正文色 |
| Progress（环形） | 直径 64px、环宽 4px，偏小 | 直径 120px、环宽 8px、中心 20px 文本（antd 默认 circle） |
| Alert | 彩色圆形底 + 白色符号（antd 4.x 观感），圆角 6 | 语义色 16px 图标、无圆底、8px 圆角、8/12 内边距 |
| Card | 内边距 12/16、圆角 6、无阴影 | 24px 内边距、8px 圆角、1px 描边 + 轻阴影、头部 16/24 内边距 |
| Modal | 宽 360、圆角 8、头/体/脚内边距偏小 | 宽 520（`MaxWidth` 640）、头 20/24、体 24、脚 16/24、面板阴影 |
| Drawer（前序） | 无滑入动画（已记录） | 保持简化，见 §4 |

## 3. Gallery 打磨

| 区域 | 打磨前 | 打磨后 |
| --- | --- | --- |
| 整体骨架 | 顶部控制条 + 扁平列表 + 卡片堆叠 | 左侧品牌/搜索/分类目录 + 右侧面包屑/标题/工具条/卡片流 |
| 侧栏 | 无品牌区，分类标题与条目样式相同，计数含义不明 | 品牌 + 版本、带图标搜索框与清除按钮、分类标题（`N 控件`）+ 条目（`N 例`）、选中项主色指示条、底部统计（分组/示例/MIT/版本） |
| 顶部工具条 | 亮/暗、密度、色块、对比挤在一行，窄窗口下「品牌色」竖排折行 | 面包屑 + 标题 + Tag + 简介独立成行；工具条独占一行，分组标签 + 分段控制器（antd Segmented 外观），并加入「无/亮暗对比/密度对比」 |
| 卡片 | 标题 + 描述 + 预览 + 纯文本代码 | 序号 + 标题 + 描述 + 复制代码；预览区浅灰承托（最小高度 96）；代码区独立头部（Go 语言标识 + 展开/收起），代码用 Scroll 承载避免溢出卡片 |
| 空状态 | 一行文本 | `ant.Empty` + 引导文案（搜索无结果时） |
| 截图 | 960×640 旧版 | `docs/images/gallery.png`（1280×800 亮色）、`gallery-dark.png`（暗色）、`quickstart.png`（最小示例） |

## 4. 仍存在的差异（有意保留 / 后续）

| 项 | 说明 |
| --- | --- |
| Drawer / Modal 动画 | 直接显示/隐藏，未接入 antd 的滑入与缩放曲线（MyGo `DialogBase` 无动画参数） |
| Tooltip 箭头 | 反色气泡已对齐，但未绘制箭头（依赖原语能力） |
| Dropdown 触发按钮 | 使用原生 `ui.MenuButton`（尺寸/字号/描边无法完全覆盖），下拉项样式由原生菜单绘制 |
| Switch 尺寸 | 40×22（antd 44×22），原生 Switch 尺寸固定 |
| Select / Combobox 内边距 | 透传原生 `ui.Select`，仅统一高度/字号/圆角，未做到 11px 水平内边距 |
| Gallery 全局尺寸切换 | 控件尺寸由各示例显式指定，无统一注入点（见 docs/07） |
| 系统亮暗自动跟随 | MyGo `ui` 未暴露系统外观 API（见 ADR 0002） |

## 5. 验证方式

```bash
# 全部测试（含 Gallery 与逐示例截图基线）
go test -count=1 ./...

# 重新生成截图基线
go test ./gallery ./samples -count=1 -run 'TestGalleryBaselines|TestControlScreenshots'

# 重新生成 README 截图（1280×800 亮色/暗色 + 最小示例）
WRITE_DOCS_SCREENSHOTS=1 go test ./gallery -run TestWriteDocsScreenshots

# 单组控件人工核对
go run ./scripts/visualcheck -group 表格 -theme light -out .visualcheck
```
