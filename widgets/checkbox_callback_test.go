package ant_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

// 验证 CheckWithLabel 回调是否恰好触发一次（点标签 / 点复选框两种路径）。
func TestCheckWithLabelSingleCallback(t *testing.T) {
	checked := false
	calls := 0

	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Checkbox(c, &checked, "同步", widgets.OnChange(func() { calls++ }))
	}
	tt := ui.NewTester(view, 240, 60)
	tt.Frame()

	// 路径 1：点击标签文字
	calls = 0
	if err := tt.Click("同步"); err != nil {
		t.Fatalf("点标签失败: %v", err)
	}
	tt.Frame()
	t.Logf("点标签: checked=%v calls=%d", checked, calls)
	if calls != 1 {
		t.Errorf("点标签触发回调 %d 次, want 1", calls)
	}

	// 路径 2：点击复选框本体（用坐标，标签文字第二次点击仍是标签）
	calls = 0
	if err := tt.Click("同步"); err != nil {
		t.Fatalf("第二次点击失败: %v", err)
	}
	tt.Frame()
	t.Logf("再点标签: checked=%v calls=%d", checked, calls)
	if calls != 1 {
		t.Errorf("第二次点击触发回调 %d 次, want 1", calls)
	}
}
