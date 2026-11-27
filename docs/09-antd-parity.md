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
| Drawer | 宽 320、头 40、体 16 | 宽 378（antd 默认）、头 56、体 24、标题 16px；打开时 220ms 从右侧滑入 |
| List | 行内边距 8/12、无悬停反馈 | 行内边距 12/16、选中/悬停浅底、1px 行分隔线 |
| Descriptions | 标签 12px、内边距 12 | 标签/值 14px、16px 内边距、8px 圆角带边框 |
| Timeline | 圆点 8px、连线 1px、内容间距 12 | 圆点 10px、连线 2px、内容下间距 20、标题/描述 14px |
| Switch | 原生 40×22 | 自绘 44×22 轨道 + 18px 滑块（2px 内边距），选中主色、可点击标签切换 |
| Select | 透传原生 Select，内边距/箭头不可控 | 自绘触发器：32px 高、11px 水平内边距、6px 圆角、右侧箭头；下拉面板 32px 选项、选中/悬停浅底、键盘导航 |
| Tag | 高度 20、内边距 8、关闭按钮无悬停反馈 | 高度 22、内边距 7、关闭按钮 16px + 悬停浅底 + 4px 间距 |
| Badge | 角标 16×16、字号 12 | 角标 20×20、圆角 10、字号 12、水平内边距 6 |
| Gallery 窗口 | 默认 800×600 | 默认 1280×800（最小 960×600），标题带版本号 |
| Gallery 字号 | 无全局字号切换 | 顶部「字号 小/中/大」（12/14/16），工具条支持换行 |
| Modal | 直接出现 | 打开时 200ms 从 0.92 缩放 + 淡入（antd Modal 的开场） |
| Transfer | 行首圆形符号、行高 24、表头白底 | 方形复选框、32px 行高 + 悬停浅底、40px 浅灰表头、32×32 描边箭头按钮 |
| Upload | 48px 拖放区、12px 提示文字 | 96px 拖放区、40px 上传图标、16px 主文案 + 14px 提示、文件行悬停/移除反馈 |
| Cascader / TreeSelect 触发器 | 文本 `▾`、无悬停反馈 | antd chevron 图标、悬停主色描边、占位文字弱化 |
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

## 3.1 第二轮 Gallery 打磨（antd 视觉复审）

对照 antd v5 Segmented / Statistic / Card / 文档站 demo 区块规范复审 Gallery 外壳：

| 区域 | 打磨前 | 打磨后 |
| --- | --- | --- |
| Segmented | 容器带 1px 边框、项高 24px、选中项无边框、字号 11.5px | 容器无边框 + 8px 圆角（antd 无描边承托底），项高 28px、水平内边距 12px，选中项容器底 + 1px 描边 + 阴影 + 500 字重，未选中项用次要文字色 |
| 首页主视觉 | 主标题/副标题/描述统一 16px 间距，副标题用主文字色，GitHub 入口为裸链接 | 标题区 8px 紧凑节奏、描述上方 8px，行动区上移 24px 独立呼吸空间；副标题改次要文字色 18px |
| 统计卡 | 内边距 24/16、数字与标题间距 4px、圆角 8px | 纵向内边距 24px、间距 8px、圆角 10px（antd borderRadiusLG） |
| 特性卡 | 480px 基准宽 → 1280 窗口下 3+1 折行不齐 | 384px 基准宽 → 2×2 等宽网格，窄窗口逐级折行 |
| 快速开始 | 纯文本 muted 色代码、无行号无高亮 | 只读 CodeEditor（语法高亮 + 行号），「Go」标签改浅灰 pill，代码区浅灰承托 |
| 示例卡预览区 | 灰底（与页面同色，卡片内无层次） | 容器白底 + 底部 1px 分隔（对齐 antd demo 白底预览），卡片间距始终 24px |
| 示例卡序号 | 10px 次要灰 | 12px/600 次要色，标题区左间距 12px |
| 代码区头 | 4px/16px 内边距、「Go」裸文字 | 8px/24px 内边距、「Go」浅灰 pill；代码区浅灰承托 |
| 侧栏搜索框 | 灰底 + 1px 边框双重描边 | 仅灰底无边框（对齐 antd 文档站侧栏搜索） |
| 侧栏分类标题 | 9.5px | 10.5px/600，计数 10px |

控件本体未发现新的规范违背：颜色/圆角/高度/字号/焦点环仍由 Token 驱动
（§1 基准），`go run ./internal/check` 样式铁律通过；本轮改动仅触及
Gallery 外壳（`gallery/app.go`），逐示例基线不受影响，Gallery 级基线与
README 截图已重新生成。

### 3.2 第三轮：侧栏与工具条细节（用户复审反馈）

| 区域 | 打磨前 | 打磨后 |
| --- | --- | --- |
| 侧栏搜索框 | 手搓 Row 承托 + 裸 ui.TextInput（原语默认边框未清空，灰底药丸里嵌套带边框输入框，双重框） | 复用库自身 Input：正规 antd 框架（容器底描边、聚焦主色描边 + 2px 焦点环），前缀搜索图标、后缀 allowClear 清除按钮 |
| 导航项 | 32px 高 + 左侧主色竖条（非 antd Menu 形制） | 对齐 antd Menu：40px 项高、8px 圆角、16px 文字缩进、去掉竖条，补 hover 浅填充态 |
| 品牌色圆点 | 选中态 2px 黑圈，生硬 | antd ColorPicker 形制：白色对勾 + 1px 描边环 |
| Button 等控件配色 | 复审 ButtonColors：primary/default/dashed/text/link 五型与 hover/pressed 色阶均对齐 antd v5（default hover 文字与描边变主色、text hover 浅填充、link 主色文字），未发现违背 | 无需改动 |

## 4. 仍存在的差异（有意保留 / 后续）

| 项 | 说明 |
| --- | --- |
| Drawer / Modal 动画 | 打开动效已实现（Drawer 220ms 滑入、Modal 200ms 缩放淡入，尊重 ReduceMotion）；关闭仍为即时（MyGo `DialogBase` 关闭即卸载面板，无退场钩子） |
| Tooltip 箭头 | 反色气泡已对齐；箭头在无头 Tester 中无法稳定绘制（原生 tooltip 层不参与 `Tester.Image()`），未纳入视觉核对 |
| Dropdown 触发按钮 | 使用原生 `ui.MenuButton`（尺寸/字号/描边无法完全覆盖），下拉项样式由原生菜单绘制 |
| Select / Combobox | Select 已自绘对齐 antd；Combobox 现基于 `ui.ComboboxBase` 自绘输入框（32px/11px 内边距）与 32px 候选面板；AutoComplete 仍透传原生 |
| Gallery 控件尺寸切换 | ✅ 已实现——顶部「尺寸 小/中/大」通过 `theme.WithSpacing/WithRadius/WithFontSize` 全局缩放（示例内部的三档尺寸演示仍保留） |
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
