package theme

import "github.com/egoist/mygo/ui"

// antd v5 默认色板（docs/04-theme-system.md 的映射表）。
var (
	antdBlue       = ui.Hex("#1677ff") // colorPrimary（亮）
	antdBlueDark   = ui.Hex("#1668dc") // colorPrimary（暗）
	antdGreen      = ui.Hex("#52c41a") // colorSuccess（亮）
	antdGreenDark  = ui.Hex("#49aa19") // colorSuccess（暗）
	antdGold       = ui.Hex("#faad14") // colorWarning（亮）
	antdGoldDark   = ui.Hex("#d89614") // colorWarning（暗）
	antdRed        = ui.Hex("#ff4d4f") // colorError（亮）
	antdRedDark    = ui.Hex("#dc4446") // colorError（暗）
	antdGrayBorder = ui.Hex("#d9d9d9") // colorBorder（亮）
	antdLayout     = ui.Hex("#f5f5f5") // colorBgLayout（亮）
	antdDarkBg     = ui.Hex("#000000") // colorBgLayout（暗）
	antdDarkPanel  = ui.Hex("#141414") // colorBgContainer（暗）
	white          = ui.RGB(255, 255, 255)
	black          = ui.RGB(0, 0, 0)
)

// Light 返回 Ant Design 亮色主题（标准密度）。
func Light(opts ...Option) *AntToken {
	t := &AntToken{
		Primary:       antdBlue,
		Success:       antdGreen,
		Warning:       antdGold,
		Error:         antdRed,
		Info:          antdBlue,
		Text:          black.Alpha(0.88), // antd colorText
		TextSecondary: black.Alpha(0.65), // antd colorTextSecondary
		BgContainer:   white,             // antd colorBgContainer
		BgLayout:      antdLayout,        // antd colorBgLayout
		Border:        antdGrayBorder,    // antd colorBorder
		Radius:        6,
		FontSize:      14,
		Spacing:       4,
		Density:       Comfortable,
	}
	apply(t, opts)
	return t
}

// Dark 返回 Ant Design 暗色主题（标准密度），色值对齐 antd dark 算法。
func Dark(opts ...Option) *AntToken {
	t := &AntToken{
		Primary:       antdBlueDark,
		Success:       antdGreenDark,
		Warning:       antdGoldDark,
		Error:         antdRedDark,
		Info:          antdBlueDark,
		Text:          white.Alpha(0.85), // antd colorText（dark）
		TextSecondary: white.Alpha(0.65), // antd colorTextSecondary（dark）
		BgContainer:   antdDarkPanel,
		BgLayout:      antdDarkBg,
		Border:        white.Alpha(0.15), // antd colorBorder（dark ≈ rgba(255,255,255,0.15)）
		Radius:        6,
		FontSize:      14,
		Spacing:       4,
		Density:       Comfortable,
		isDark:        true,
	}
	apply(t, opts)
	return t
}

// apply 依序套用 Option。WithDensity 会重置度量，因此密度选项应放在
// 字号/圆角等覆盖之前传入；这里通过"先整体后个别"的顺序保证直觉行为：
// 传入的 Option 按顺序执行，WithDensity 放前面即可被后续覆盖。
func apply(t *AntToken, opts []Option) {
	for _, o := range opts {
		o(t)
	}
}
