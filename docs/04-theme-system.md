# 04 · 主题系统设计

> 对应需求 5：所有控件共用主题，可运行时切换主题、颜色、间距等。

## 设计目标

1. **单一真源**：一个主题对象驱动全部换肤控件与缺口控件，杜绝硬编码。
2. **运行时切换**：用户可在运行时切换亮/暗、品牌色、密度，无需重启；支持按窗口、按子树。
3. **Ant Design 语义**：用 Ant Design 的设计_token_命名与色板，让 familiarity 直接迁移。
4. **跟随系统**：默认跟随系统外观、强调色、高对比度、文字缩放（复用 MyGo 的 `WithPreferences`）。
5. **可扩展**：允许业务追加自定义 token（如品牌专属色）而不破坏内建控件。

## 与 MyGo Theme 的关系

MyGo 的 `ui.Theme` 是颜色与度量的**物理载体**，本库在其上定义 Ant Design 语义的**逻辑层**：

```
AntToken（语义层，antd 风格命名）
        │  映射函数
        ▼
ui.Theme（物理层，MyGo 原生）
        │  c.SetTheme(&t) / ui.Themed(c, t)
        ▼
全部控件（只从 c.Theme() 取色与度量）
```

**关键原则**：控件永远只读 `c.Theme()`，不直接读 `AntToken`。这样任何主题切换只要最终落到 `SetTheme`，全部控件自动跟随；局部主题用 `ui.Themed` 也自动生效（含菜单/弹层/Tooltip）。

## Token 映射表（Ant Design → MyGo Theme 字段）

### 颜色

| Ant Design token | 语义 | MyGo Theme 字段 | 说明 |
| --- | --- | --- | --- |
| `colorPrimary` | 品牌主色 | `Accent` | 用 `WithAccent` 自动生成整套色阶 |
| `colorPrimaryHover` | 主色悬浮 | `AccentHover` | 同上 |
| `colorPrimaryActive` | 主色按下 | `AccentPressed` | 同上 |
| `colorText` | 主文本 | `Text` | |
| `colorTextSecondary` | 次要文本 | `TextMuted` | |
| `colorTextLightSolid` | 主色上的文本 | `AccentText` | |
| `colorBgContainer` | 组件容器背景 | `Surface` | antd 组件底色多为白/深灰 |
| `colorBgContainerHover` | 容器悬浮 | `SurfaceHover` | |
| `colorBgContainerPressed` | 容器按下 | `SurfacePressed` | |
| `colorBgLayout` | 页面底色 | `Background` | |
| `colorBorder` | 描边 | `Border` | |
| `colorError` | 错误 | `Danger` | |
| `colorWarning` | 警告 | `Warning` | |
| `colorSuccess` | 成功 | `Success` | |
| `colorInfo` | 信息 | 局部用 `Accent` 或自定义 Element 色 | MyGo 无独立 Info 槽 |
| `colorTextSelectionBg` | 文本选中 | `Selection` | |
| `colorPrimaryBorder` / focus | 焦点环 | `Focus` + `FocusRingWidth/Offset` | |
| `colorBgElevated`（弹层） | Tooltip/Toast 反色 | `Inverse` / `InverseText` | 注意 MyGo 语义是「反色」，antd 弹层是浮起色，需按预设区分 |
| `colorFillAlter`（弱填充） | 表头/斑马纹 | 用 `Surface` + `Alpha` 局部表达 | |
| 非活动窗口主色 | — | `InactiveAccent` | MyGo 特有，设为 `Accent` 可保持各窗口一致 |

> 说明：Ant Design 的 `colorPrimaryBg`、`colorFillQuaternary` 等细粒度 token 在 MyGo 中无独立槽位，统一通过 `Surface` 派生或在控件内部用 `ui.Color.Mix/Alpha` 表达，并在 `AntToken` 中保留字段以便未来扩展。

### 度量

| Ant Design token | MyGo 字段 / 方法 | 换算 |
| --- | --- | --- |
| `borderRadius` | `Theme.Radius` | 直接映射（antd 默认 6，与 MyGo 一致） |
| `fontSize` | `Theme.FontSize` | 直接映射 |
| `fontFamily` | `Theme.Font` | 留空=系统字体 |
| `sizeUnit` / `controlHeight` | `Theme.Spacing` + `t.Space(n)` / `t.Rem(n)` | 见 [05 紧凑密度](05-compact-density.md) |
| `marginXS/SM/M/L/XL` | `t.Space(n)` | n = 1/2/3/4/6（Spacing=4 时即 4/8/12/16/24） |
| `lineWidth` | 链式 `Border(1, t.Border)` | 固定 1 |
| 焦点环宽度/偏移 | `FocusRingWidth` / `FocusRingOffset` | 直接映射 |
| 运动时长/曲线 | `ui.Motion` / `AnimateTo` / transition | 由 MyGo 动效系统表达 |

## 主题对象（设计草案，非实现代码）

```go
// theme 包对外的语义对象，字段对齐 Ant Design token 的常用子集。
type AntToken struct {
    // 颜色
    Primary, Success, Warning, Error, Info Color
    Text, TextSecondary, TextInverse      Color
    BgContainer, BgLayout                 Color
    Border                                Color
    // 度量
    Radius   float32
    FontSize float32
    Spacing  float32   // MyGo 的密度单位；见 05 文档
    Density  Density   // Comfortable | Compact
    // 派生色阶由 WithAccent 自动生成，不在此重复声明
}
```

对外 API 草案：

```go
theme.Light(opts ...Option) *AntToken   // antd 默认亮色
theme.Dark(opts ...Option)  *AntToken   // antd 暗色算法近似
theme.WithPrimary(hex string) Option    // 品牌色，自动生成色阶
theme.WithDensity(d Density) Option     // 密度切换
theme.Apply(c *ui.Context, t *AntToken) // 转成 ui.Theme 并 c.SetTheme
theme.Sub(c *ui.Context, t *AntToken)   // 局部主题：包一层 ui.Themed
```

`Apply` 内部流程：
1. 以 `ui.LightTheme()` / `ui.DarkTheme()` 为基底；
2. 写入 Token 映射表中的颜色与度量；
3. 对 Primary 调用 `WithAccent` 生成 `AccentHover/Pressed/Selection/Focus/AccentText`；
4. 按 Density 调整 `Spacing`/`FontSize`/`Radius`（见 05）；
5. 默认叠加 `WithPreferences(c.Preferences())` 以跟随系统；
6. `c.SetTheme(&built)`。

## 运行时切换机制

| 场景 | 机制 | 效果 |
| --- | --- | --- |
| **全局换主题** | 用户操作 → 更新 app 状态 → `theme.Apply(c, next)` | 下一帧全窗口生效；MyGo 每帧重建 view，无需手动刷新控件 |
| **切换品牌色** | `WithPrimary(hex)` 重新生成色阶 → `Apply` | 主按钮、选中态、焦点环同步变色 |
| **切换密度** | `WithDensity(Compact)` → `Apply` | `Spacing` 变化，全部控件 padding/gap/行高同步 |
| **单个子树换主题** | `theme.Sub(c, t)`（内部 `ui.Themed`） | 只影响该子树，弹层/菜单/Tooltip 同样跟随；可做亮暗并排预览 |
| **跟随系统外观** | 默认不覆盖 `Dark` 字段，由 MyGo 默认主题随系统；或用 `WithPreferences` | 系统切换亮暗时控件自动跟随 |

> 注意：MyGo 的 `SetTheme` 改变的是「调用之后构建的控件」，因此主题切换按钮本身若想立刻变外观，应把主题选择放在 app state，由下一帧统一重建。这是即时模式的正常心智，不是限制。

## 预设主题

| 预设 | 基底 | 特点 |
| --- | --- | --- |
| `PresetLight` | antd 默认亮色 | 主色 `#1677ff`，圆角 6，Spacing 4 |
| `PresetDark` | antd 暗色算法近似 | 主色 `#4096ff`，深底 `#141414` 系 |
| `PresetCompactLight/Dark` | 上面两套 + `Density=Compact` | 用于上位机 / IDE（见 05） |
| `PresetHighContrast` | `WithPreferences` 的高对比分支 | 描边加深、焦点环不透明 |

## 一致性校验（防回归）

1. **单测**：遍历控件 gallery，断言所有颜色来自 Theme（静态检查：源码中不得出现裸 `ui.Hex("#...")`，图标数据除外）。
2. **视觉回归**：用 MyGo 的「无窗口渲染成 PNG」能力，对每套预设截图对比。
3. **切换测试**：运行时切换主题/密度后截图，确认无遗漏控件仍是旧色。
