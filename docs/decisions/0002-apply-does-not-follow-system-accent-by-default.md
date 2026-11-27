# 0002 · Apply 默认不跟随系统强调色

日期：2026-11-26
状态：已采纳（2026-11-27 补充实现状态）

## 背景

`docs/04-theme-system.md` 初稿提出「Apply 默认叠加 `WithPreferences(c.Preferences())` 以跟随系统」。但 MyGo 的 `Theme.WithPreferences` 会把主题的 Accent 换成桌面强调色并派生整套色阶——这与 Ant Design 明确的品牌主色（`colorPrimary`，如 `#1677ff`）直接冲突：叠加后 `WithPrimary` 设置的品牌色会被系统色覆盖，控件不再呈现 antd 观感。

## 决策

`theme.Apply` **默认不跟随**系统强调色 / 高对比度 / 文字缩放，保证 Ant Design 品牌色忠实呈现。

当前实现（`theme/theme.go` + `theme/bridge.go`）：

- `Apply(c, t)` 只做 `c.SetTheme(t.toTheme())`；
- `toTheme()` 按 `isDark` 选择 `ui.LightTheme()` / `ui.DarkTheme()` 基底，写入 Token 映射，再 `WithAccent(Primary)` 派生成套色阶；
- **不会**调用 `WithPreferences`，也没有 `WithFollowSystem()` Option。

需要跟随系统的应用，可以在 `theme.Apply` 之后自行调用 `c.SetTheme(c.Theme().WithPreferences(c.Preferences()))`（注意这会覆盖品牌色，属于有意的应用级取舍）。

## 后续计划（未实现）

- 提供 `theme.WithFollowSystem()` Option：在品牌色之上叠加系统偏好，并明确「强调色以谁为准」的优先级（建议：品牌色优先，仅高对比/文字缩放跟随）。
- 系统亮暗切换的自动跟随（依赖 MyGo 的 dark mode 通知）。

## 后果

- 正向：默认输出即标准 antd 观感；品牌色与系统色不冲突。
- 负向：无障碍偏好（高对比、文字缩放）默认不生效，需要应用主动开启；系统亮暗切换暂不自动跟随（记录为 TODO）。
