# 04 · 主题系统设计

> 对应需求 5：所有控件共用主题，可运行时切换主题、颜色、间距等。
>
> 本文描述**当前已实现**的主题层（`theme/`），并标注尚未实现的能力。核对基线：MyGo v0.10.0、本库 v0.1.0（v0.2 开发中）。

## 设计目标

1. **单一真源**：一个 `AntToken` 驱动全部换肤控件与缺口控件，杜绝硬编码。
2. **运行时切换**：可在运行时切换亮/暗、品牌色、密度，无需重启；支持按窗口、按子树。
3. **Ant Design 语义**：用 Ant Design 的 design token 命名与色板，让 familiarity 直接迁移。
4. **可扩展**：允许业务通过 `Option` 追加/覆盖 token，而不破坏内建控件。
5. **可回归**：主题映射有单测、静态检查与截图基线兜底。

> ⚠️ **默认不跟随系统强调色/高对比度/文字缩放**（决策见 [ADR 0002](decisions/0002-apply-does-not-follow-system-accent-by-default.md)）。早期草案里「默认叠加 `WithPreferences`」与「`WithFollowSystem()` Option」在当前实现中**都不存在**，属于未实现计划。

## 与 MyGo Theme 的关系

MyGo 的 `ui.Theme` 是颜色与度量的**物理载体**，本库在其上定义 Ant Design 语义的**逻辑层**：

```
AntToken（语义层，antd 风格命名）
        │  toTheme() 映射
        ▼
ui.Theme（物理层，MyGo 原生）
        │  c.SetTheme(&t) / ui.Themed(c, t)
        ▼
全部控件（只从 c.Theme() 取色与度量）
```

**关键原则**：控件永远只读 `c.Theme()`，不直接读 `AntToken`。这样任何主题切换只要最终落到 `SetTheme`，全部控件自动跟随；局部主题用 `ui.Themed` 也自动生效（含菜单/弹层/Tooltip）。

## 主题对象（`theme/token.go`，已实现）

```go
// Density 控件密度。
type Density string

const (
	Comfortable Density = "comfortable" // Spacing 4 / Radius 6 / FontSize 14
	Compact     Density = "compact"     // Spacing 3 / Radius 4 / FontSize 13
)

// AntToken 是 Ant Design 语义 Token 的常用子集。
type AntToken struct {
	Primary       ui.Color // antd colorPrimary
	Success       ui.Color
	Warning       ui.Color
	Error         ui.Color
	Info          ui.Color // 无独立 Theme 槽位，桥接时跟随主色

	Text          ui.Color // colorText
	TextSecondary ui.Color // colorTextSecondary
	TextInverse   ui.Color // 预留：当前桥接未使用

	BgContainer ui.Color // colorBgContainer（按钮/输入框底色）
	BgLayout    ui.Color // colorBgLayout（页面底色）
	Border      ui.Color // colorBorder

	Radius   float32 // borderRadius
	FontSize float32 // fontSize
	Spacing  float32 // MyGo 的留白单位（4 标准 / 3 紧凑）
	Density  Density

	isDark bool // 私有：选择 ui.LightTheme()/ui.DarkTheme() 基底
}

type Option func(*AntToken)

func WithPrimary(c ui.Color) Option   // 品牌主色（色阶由 WithAccent 派生）
func WithDensity(d Density) Option    // 同时重置 Spacing/Radius/FontSize
func WithFontSize(px float32) Option
func WithRadius(px float32) Option
func WithSpacing(px float32) Option   // 高级用法，常规请用 WithDensity
```

> **Option 顺序**：`WithDensity` 会重置度量，因此应放在 `WithFontSize` / `WithRadius` / `WithSpacing` **之前**传入，否则后传的密度会覆盖前传的细调。

## Token 映射表（`theme/bridge.go`）

### 颜色

| Ant Design token | 语义 | MyGo `ui.Theme` 字段 | 桥接规则 |
| --- | --- | --- | --- |
| `colorPrimary` | 品牌主色 | `Accent` | 直接写入，再由 `WithAccent` 派生整套色阶 |
| `colorPrimaryHover` | 主色悬浮 | `AccentHover` | `WithAccent` 生成（亮：Mix 黑 12%；暗：Mix 白 20%） |
| `colorPrimaryActive` | 主色按下 | `AccentPressed` | `WithAccent` 生成（亮：Mix 黑 25%；暗：Mix 黑 15%） |
| `colorText` | 主文本 | `Text` | 直接写入 |
| `colorTextSecondary` | 次要文本 | `TextMuted` | 直接写入 |
| `colorTextLightSolid` | 主色上的文本 | `AccentText` | `WithAccent` 生成（浅色主色自动用深字） |
| `colorBgContainer` | 组件容器背景 | `Surface` | 直接写入（亮 `#ffffff` / 暗 `#141414`） |
| `colorBgContainerHover` | 容器悬浮 | `SurfaceHover` | `BgContainer.Mix(Text, 0.06)`（≈ antd `colorFillSecondary`） |
| `colorBgContainerPressed` | 容器按下 | `SurfacePressed` | `BgContainer.Mix(Text, 0.12)` |
| `colorBgLayout` | 页面底色 | `Background` | 亮 `#f5f5f5` / 暗 `#000000` |
| `colorBorder` | 描边 | `Border` | 亮 `#d9d9d9` / 暗 `rgba(255,255,255,0.15)` |
| `colorError` | 错误 | `Danger` | 亮 `#ff4d4f` / 暗 `#dc4446` |
| `colorWarning` | 警告 | `Warning` | 亮 `#faad14` / 暗 `#d89614` |
| `colorSuccess` | 成功 | `Success` | 亮 `#52c41a` / 暗 `#49aa19` |
| `colorInfo` | 信息 | — | MyGo 无独立 Info 槽；`AntToken.Info` 当前不参与桥接，需要时用主色 |
| `colorTextSelectionBg` | 文本选中 | `Selection` | `WithAccent` 生成（亮：主色 25% 透明） |
| `colorPrimaryBorder` / focus | 焦点环 | `Focus` + `FocusRingWidth/Offset` | `WithAccent` 生成 `Focus`；环宽/偏移沿用基底主题 |
| `colorBgElevated`（弹层） | Tooltip/Toast 反色 | `Inverse` / `InverseText` | ⚠️ MyGo 语义是「反色」，与 antd 浮起色不同；当前预设未使用 |
| `colorFillAlter`（弱填充） | 表头/斑马纹 | — | 无独立槽位，控件内用 `Surface` + `Mix`/`Alpha` 表达 |
| 非活动窗口主色 | — | `InactiveAccent` | 未设置，沿用 MyGo 默认派生 |

### 度量

| Ant Design token | MyGo 字段 / 方法 | 换算 |
| --- | --- | --- |
| `borderRadius` | `Theme.Radius` | 直接映射（antd 默认 6，与 MyGo 一致） |
| `fontSize` | `Theme.FontSize` | 直接映射（标准 14 / 紧凑 13） |
| `fontFamily` | `Theme.Font` | 当前留空（走系统字体栈），未开放选项 |
| `sizeUnit` / `controlHeight` | `Theme.Spacing` + `t.Space(n)` / `t.Rem(n)` | 见 [05 紧凑密度](05-compact-density.md) |
| `marginXS/SM/M/L/XL` | `t.Space(n)` | n = 1/2/3/4/6（Spacing=4 时即 4/8/12/16/24） |
| `lineWidth` | 链式 `Border(1, t.Border)` | 固定 1px（静态检查豁免） |
| 焦点环宽度/偏移 | `FocusRingWidth` / `FocusRingOffset` | 直接沿用 MyGo 基底 |
| 运动时长/曲线 | `ui.Animation` / `el.AnimateTo` / `ui.Motion` | 由 MyGo 动效系统表达 |

## 对外 API（`theme/theme.go`）

```go
// 构造预设（返回可继续链式/复用的 token）。
theme.Light(opts ...Option) *AntToken
theme.Dark(opts ...Option)  *AntToken

// 全局主题：每帧调用安全，是切换主题的标准入口。
theme.Apply(c *ui.Context, t *AntToken)          // 内部 c.SetTheme(t.toTheme())

// 局部主题：返回 ui.Themed 容器，只影响子树。
theme.Sub(c *ui.Context, t *AntToken) ui.Element // 内部 ui.Themed(c, t.toTheme())
```

`Apply` 内部流程（与 `bridge.go` 一致）：

1. 按 `isDark` 选择 `ui.LightTheme()` / `ui.DarkTheme()` 为基底；
2. 写入 Token 映射表中的颜色（`Accent` / `Danger` / `Warning` / `Success` / `Text` / `TextMuted` / `Background` / `Surface` / `Border`）；
3. 由 `BgContainer.Mix(Text, …)` 派生 `SurfaceHover` / `SurfacePressed`；
4. 写入 `Radius` / `FontSize` / `Spacing`；
5. 调用 `Theme.WithAccent(Primary)` 生成 `AccentHover` / `AccentPressed` / `Selection` / `Focus` / `AccentText`；
6. `c.SetTheme(&built)`。

> 系统偏好（强调色、高对比、文字缩放）**不在**上述流程内，见 ADR 0002。

## 运行时切换机制

| 场景 | 机制 | 效果 |
| --- | --- | --- |
| **全局换主题** | 用户操作 → 更新 app 状态 → `theme.Apply(c, next)` | 下一帧全窗口生效；MyGo 每帧重建 view，无需手动刷新控件 |
| **切换品牌色** | `WithPrimary(ui.Hex("#7c3aed"))` 重新生成色阶 → `Apply` | 主按钮、选中态、焦点环同步变色 |
| **切换密度** | `WithDensity(theme.Compact)` → `Apply` | `Spacing` / `Radius` / `FontSize` 变化，全部控件 padding/gap/行高同步 |
| **单个子树换主题** | `theme.Sub(c, t)`（内部 `ui.Themed`） | 只影响该子树，弹层/菜单/Tooltip 同样跟随；可做亮暗并排预览 |
| **跟随系统外观** | ❌ 未实现 | `Apply` 总是显式选择亮/暗基底；系统亮暗切换不会自动跟随（TODO） |

> 注意：MyGo 的 `SetTheme` 改变的是「调用之后构建的控件」，因此主题切换按钮本身若想立刻变外观，应把主题选择放在 app state，由下一帧统一重建。这是即时模式的正常心智，不是限制。

## 预设主题

| 预设 | 基底 | 数值 |
| --- | --- | --- |
| `theme.Light()` | antd 默认亮色 | 主色 `#1677ff`，圆角 6，Spacing 4，字号 14，容器白、页面 `#f5f5f5` |
| `theme.Dark()` | antd 暗色算法近似 | 主色 `#1668dc`，容器 `#141414`、页面 `#000000`，文字白 85% |
| `theme.Light(theme.WithDensity(theme.Compact))` | 亮色 + 紧凑 | Spacing 3、圆角 4、字号 13 |
| `theme.Dark(theme.WithDensity(theme.Compact))` | 暗色 + 紧凑 | 同上 |
| 自定义品牌色 | `WithPrimary(...)` | 任意 `ui.Color`，色阶自动派生 |

> 早期草案中的 `PresetLight` / `PresetDark` / `PresetHighContrast` 名称**不存在**；高对比（`WithPreferences`）留待后续版本，见 ADR 0002。

## 一致性校验（防回归）

1. **单测**：`theme/theme_test.go` 验证 Token→Theme 各字段映射、密度数值与 `Apply`/`Sub` 行为。
2. **静态检查**：`go run ./internal/check` 扫描 `widgets/`、`extra/` 的非测试代码，禁止裸 `ui.Hex` / `ui.RGB` / `ui.RGBA` / `ui.Oklch` 与裸间距字面量（0 与 1px 描边豁免）；`theme/`、`icon/`、测试文件不扫描。
3. **视觉回归**：`gallery` 包用 `ui.NewTester` 无窗口渲染 4 种组合（亮/暗 × 标准/紧凑 + 亮暗并排），与 `samples/golden/gallery-*.png` 逐字节比对。
4. **切换测试**：`TestGalleryInteractions` 运行时切换主题/密度，断言无 panic、示例区不丢失。

## 未实现 / TODO

- 跟随系统强调色与无障碍偏好（`WithFollowSystem` Option 尚未提供）。
- 系统亮暗切换的自动跟随（需 MyGo 的 dark mode 通知）。
- `TextInverse` 参与桥接（当前保留字段但未写入）。
- 高对比预设（`PresetHighContrast`）。
