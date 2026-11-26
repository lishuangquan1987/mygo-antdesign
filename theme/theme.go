package theme

import "github.com/egoist/mygo/ui"

// Apply 把 token 转为 ui.Theme 并设为整个窗口的主题：窗口根、根背景
// 与之后构建的控件全部生效。在 view 内每帧调用是安全的，也是用户切换
// 主题的标准入口：
//
//	func (app *App) view(c *ui.Context) {
//	    theme.Apply(c, app.token) // app.token 可在运行时替换
//	    // ... 之后的控件全部拿到 antd 主题
//	}
//
// 默认不跟随桌面强调色与无障碍偏好（见 ADR 0002）。
func Apply(c *ui.Context, t *AntToken) {
	c.SetTheme(t.toTheme())
}

// Sub 返回一个局部主题容器：Children 内的控件用 t 构建，弹出层、菜单、
// Tooltip 同样跟随；容器之外的控件保持窗口主题。适合亮/暗并排预览、
// 独立主题面板等场景：
//
//	row := ui.Row(c).Gap(8).Fill()
//	left := theme.Sub(c, theme.Light())   // 亮色列
//	right := theme.Sub(c, theme.Dark())   // 暗色列
func Sub(c *ui.Context, t *AntToken) ui.Element {
	return ui.Themed(c, t.toTheme())
}
