// Command visualcheck 把示例在指定主题下逐个渲染成 PNG，便于人工核对控件效果。
//
// 用法：
//
//	go run ./scripts/visualcheck -group 数据展示 -theme light
//	go run ./scripts/visualcheck -theme dark -out .visualcheck/dark
//
// 它复用 Gallery 注册的示例（View 函数），因此看到的就是 Gallery 里的效果；
// 与 computer use 无关，纯 ui.NewTester 无窗口渲染。
package main

import (
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/egoist/mygo/ui"

	sample "github.com/lishuangquan1987/mygo-antdesign/samples"
	"github.com/lishuangquan1987/mygo-antdesign/samples/button"
	"github.com/lishuangquan1987/mygo-antdesign/samples/display"
	"github.com/lishuangquan1987/mygo-antdesign/samples/divider"
	"github.com/lishuangquan1987/mygo-antdesign/samples/extrademo"
	"github.com/lishuangquan1987/mygo-antdesign/samples/feedback"
	"github.com/lishuangquan1987/mygo-antdesign/samples/input"
	"github.com/lishuangquan1987/mygo-antdesign/samples/layout"
	"github.com/lishuangquan1987/mygo-antdesign/samples/nav"
	"github.com/lishuangquan1987/mygo-antdesign/samples/selectdemo"
	"github.com/lishuangquan1987/mygo-antdesign/samples/table"
	"github.com/lishuangquan1987/mygo-antdesign/samples/typography"
	"github.com/lishuangquan1987/mygo-antdesign/theme"
)

func demos() []sample.WidgetDemo {
	return []sample.WidgetDemo{
		button.Demo(), typography.Demo(), divider.Demo(), layout.Demo(),
		input.Demo(), selectdemo.Demo(), display.Demo(), table.Demo(),
		feedback.Demo(), nav.Demo(), extrademo.Demo(),
	}
}

func token(dark, compact bool) *theme.AntToken {
	var opts []theme.Option
	if compact {
		opts = append(opts, theme.WithDensity(theme.Compact))
	}
	if dark {
		return theme.Dark(opts...)
	}
	return theme.Light(opts...)
}

func main() {
	out := flag.String("out", ".visualcheck", "输出目录")
	group := flag.String("group", "", "只渲染指定示例分组（按分组名匹配）")
	themeName := flag.String("theme", "light", "light | dark | compact")
	width := flag.Int("w", 800, "画布宽度")
	height := flag.Int("h", 600, "画布高度")
	flag.Parse()

	dark, compact := false, false
	switch *themeName {
	case "light":
	case "dark":
		dark = true
	case "compact":
		compact = true
	default:
		log.Fatalf("未知主题 %q（light|dark|compact）", *themeName)
	}
	tk := token(dark, compact)

	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatal(err)
	}
	n := 0
	for _, d := range demos() {
		if *group != "" && d.Name != *group {
			continue
		}
		for i, s := range d.Samples {
			s := s
			view := func(c *ui.Context) {
				theme.Apply(c, tk)
				s.View(c)
			}
			tt := ui.NewTester(view, *width, *height)
			tt.Frame()
			img := tt.Image()
			if img == nil {
				log.Fatalf("%s/%s：无渲染帧", d.Name, s.Title)
			}
			name := fmt.Sprintf("%s-%02d-%s-%s.png", d.Name, i+1, s.Title, *themeName)
			path := filepath.Join(*out, name)
			f, err := os.Create(path)
			if err != nil {
				log.Fatal(err)
			}
			if err := png.Encode(f, img); err != nil {
				log.Fatal(err)
			}
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
			n++
		}
	}
	log.Printf("已输出 %d 张截图到 %s", n, *out)
}
