package theme

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestLightTokenValues(t *testing.T) {
	tk := Light()
	if tk.Primary != ui.Hex("#1677ff") {
		t.Errorf("Primary = %v, want antd #1677ff", tk.Primary)
	}
	if tk.Error != ui.Hex("#ff4d4f") {
		t.Errorf("Error = %v, want #ff4d4f", tk.Error)
	}
	if tk.Text != black.Alpha(0.88) {
		t.Errorf("Text = %v, want rgba(0,0,0,0.88)", tk.Text)
	}
	if tk.Density != Comfortable || tk.Spacing != 4 || tk.Radius != 6 || tk.FontSize != 14 {
		t.Errorf("metrics = %v/%v/%v/%v, want comfortable/4/6/14",
			tk.Density, tk.Spacing, tk.Radius, tk.FontSize)
	}
	if tk.isDark {
		t.Error("Light().isDark = true")
	}
}

func TestDarkTokenValues(t *testing.T) {
	tk := Dark()
	if tk.Primary != ui.Hex("#1668dc") {
		t.Errorf("Primary = %v, want antd dark #1668dc", tk.Primary)
	}
	if tk.BgContainer != ui.Hex("#141414") {
		t.Errorf("BgContainer = %v, want #141414", tk.BgContainer)
	}
	if !tk.isDark {
		t.Error("Dark().isDark = false")
	}
}

func TestWithDensity(t *testing.T) {
	tk := Light(WithDensity(Compact))
	if tk.Spacing != 3 || tk.Radius != 4 || tk.FontSize != 13 {
		t.Errorf("compact metrics = %v/%v/%v, want 3/4/13", tk.Spacing, tk.Radius, tk.FontSize)
	}
	if tk.Primary != ui.Hex("#1677ff") {
		t.Errorf("WithDensity 不应影响颜色，Primary = %v", tk.Primary)
	}
}

func TestOptionsOrder(t *testing.T) {
	// WithDensity 在前、覆盖在后：后面的 Option 应生效。
	tk := Light(WithDensity(Compact), WithFontSize(12))
	if tk.FontSize != 12 || tk.Spacing != 3 {
		t.Errorf("FontSize = %v, Spacing = %v, want 12 / 3", tk.FontSize, tk.Spacing)
	}
}

func TestBridge(t *testing.T) {
	tk := Light()
	th := tk.toTheme()

	if th.Background != tk.BgLayout {
		t.Errorf("Background = %v, want BgLayout %v", th.Background, tk.BgLayout)
	}
	if th.Surface != tk.BgContainer {
		t.Errorf("Surface = %v, want BgContainer %v", th.Surface, tk.BgContainer)
	}
	if th.Accent != tk.Primary {
		t.Errorf("Accent = %v, want Primary %v", th.Accent, tk.Primary)
	}
	if th.Danger != tk.Error || th.Warning != tk.Warning || th.Success != tk.Success {
		t.Errorf("semantic colors not mapped")
	}
	if th.Text != tk.Text || th.TextMuted != tk.TextSecondary {
		t.Errorf("text colors not mapped")
	}
	if th.Spacing != tk.Spacing || th.Radius != tk.Radius || th.FontSize != tk.FontSize {
		t.Errorf("metrics not mapped: %v/%v/%v", th.Spacing, th.Radius, th.FontSize)
	}

	// WithAccent 应派生成套色阶。
	if th.AccentHover == th.Accent {
		t.Errorf("AccentHover 未派生（与 Accent 相同）")
	}
	if th.Selection == ui.Transparent || th.Focus == ui.Transparent {
		t.Errorf("Selection/Focus 未派生")
	}

	// 暗色基底：Dark 标记应传递。
	if dth := Dark().toTheme(); !dth.Dark {
		t.Errorf("Dark().toTheme().Dark = false")
	}
}

func TestApplySetsWindowTheme(t *testing.T) {
	tk := Light(WithPrimary(ui.Hex("#7c3aed")))
	view := func(c *ui.Context) {
		Apply(c, tk)
		if got := c.Theme().Accent; got != ui.Hex("#7c3aed") {
			t.Errorf("Apply 后 Accent = %v, want #7c3aed", got)
		}
	}
	ui.NewTester(view, 320, 240).Frame()
}

func TestSubThemeSwitchesInside(t *testing.T) {
	// Sub 内部的控件主题应切换为暗色，外部保持亮色。
	var inside, outside ui.Theme
	view := func(c *ui.Context) {
		Apply(c, Light())
		outside = *c.Theme()
		Sub(c, Dark(WithPrimary(ui.Hex("#7c3aed")))).Children(func() {
			inside = *c.Theme()
		})
	}
	ui.NewTester(view, 320, 240).Frame()
	if !inside.Dark {
		t.Errorf("Sub 内部主题应为暗色")
	}
	if inside.Accent != ui.Hex("#7c3aed") {
		t.Errorf("Sub 内部 Accent = %v, want #7c3aed", inside.Accent)
	}
	if outside.Dark {
		t.Errorf("Sub 外部应保持亮色主题")
	}
}
func TestWithFollowSystem(t *testing.T) {
	var got *ui.Theme
	view := func(c *ui.Context) {
		Apply(c, Light(WithFollowSystem()))
		got = c.Theme()
	}
	tt := ui.NewTester(view, 320, 200)
	tt.SetPreferences(ui.Preferences{HighContrast: true, TextScale: 1.5})
	tt.Frame()
	if got == nil {
		t.Fatal("未渲染")
	}
	if got.FontSize <= 14 {
		t.Errorf("followSystem 下 FontSize = %v, want > 14（文字缩放未生效）", got.FontSize)
	}
}

func TestDefaultDoesNotFollowSystem(t *testing.T) {
	var got *ui.Theme
	view := func(c *ui.Context) {
		Apply(c, Light())
		got = c.Theme()
	}
	tt := ui.NewTester(view, 320, 200)
	tt.SetPreferences(ui.Preferences{TextScale: 1.5, Accent: ui.RGB(255, 0, 255)})
	tt.Frame()
	if got == nil {
		t.Fatal("未渲染")
	}
	if got.FontSize != 14 {
		t.Errorf("默认 FontSize = %v, want 14（不应跟随系统文字缩放）", got.FontSize)
	}
	if got.Accent != ui.Hex("#1677ff") {
		t.Errorf("默认 Accent = %v, want #1677ff（不应被系统强调色覆盖）", got.Accent)
	}
}

func TestWithFontSize(t *testing.T) {
	var got *ui.Theme
	view := func(c *ui.Context) {
		Apply(c, Light(WithFontSize(16)))
		got = c.Theme()
	}
	tt := ui.NewTester(view, 200, 120)
	tt.Frame()
	if got == nil {
		t.Fatal("未渲染")
	}
	if got.FontSize != 16 {
		t.Errorf("FontSize = %v, want 16", got.FontSize)
	}
}
