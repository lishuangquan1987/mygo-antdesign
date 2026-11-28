package extra

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
)

func pageOf(title string) func(*ui.Context) {
	return func(c *ui.Context) { ui.Text(c, title) }
}

// TestCarouselArrows 左右箭头切换页面。
func TestCarouselArrows(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		Carousel(c, []func(*ui.Context){pageOf("第一页"), pageOf("第二页"), pageOf("第三页")},
			CarouselArrows(true))
	}
	tt := ui.NewTester(view, 480, 260)
	tt.Frame()
	if !tt.HasText("第一页") {
		t.Fatal("初始应显示第一页")
	}
	if err := tt.Click("下一页"); err != nil {
		t.Fatalf("点下一页失败: %v", err)
	}
	tt.Frame()
	if !tt.HasText("第二页") {
		t.Fatal("点下一页后应显示第二页")
	}
	if err := tt.Click("上一页"); err != nil {
		t.Fatalf("点上一页失败: %v", err)
	}
	tt.Frame()
	if !tt.HasText("第一页") {
		t.Fatal("点上一页后应回到第一页")
	}
}

// TestCarouselAutoplay 自动播放：时间推进后切到下一页。
func TestCarouselAutoplay(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		Carousel(c, []func(*ui.Context){pageOf("甲"), pageOf("乙")},
			CarouselAutoplay(time.Second))
	}
	tt := ui.NewTester(view, 400, 220)
	tt.Frame()
	if !tt.HasText("甲") {
		t.Fatal("初始应为甲")
	}
	tt.Advance(1200 * time.Millisecond)
	tt.Frame()
	if !tt.HasText("乙") {
		t.Fatal("1.2s 后应自动切到乙")
	}
}

// TestUploadExtensions 扩展名过滤与数量上限。
func TestUploadExtensions(t *testing.T) {
	cfg := &uploadConfig{exts: []string{".png", ".PDF"}, maxCount: 2}
	if !allowedExt(cfg, "/tmp/a.PNG") {
		t.Error(".PNG 应匹配 .png（不区分大小写）")
	}
	if !allowedExt(cfg, "/tmp/b.pdf") {
		t.Error(".pdf 应匹配 .PDF")
	}
	if allowedExt(cfg, "/tmp/c.txt") {
		t.Error(".txt 不应通过")
	}
	files := []UploadFile{}
	addPaths(cfg, &files, []string{"/tmp/a.png", "/tmp/b.pdf", "/tmp/c.png", "/tmp/d.png"})
	if len(files) != 2 {
		t.Fatalf("maxCount=2 应只收 2 个，实收 %d", len(files))
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name, ".txt") {
			t.Errorf("不应收 %s", f.Name)
		}
	}
	// 去重：同一路径重复传入不新增。
	addPaths(cfg, &files, []string{"/tmp/a.png"})
	if len(files) != 2 {
		t.Errorf("重复路径不应新增，现有 %d", len(files))
	}
}
