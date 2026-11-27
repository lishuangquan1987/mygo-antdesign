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
	"encoding/json"
	"flag"
	"fmt"
	"image"
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
	"github.com/lishuangquan1987/mygo-antdesign/samples/ide"
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
		feedback.Demo(), nav.Demo(), extrademo.Demo(), ide.Demo(),
	}
}

func token(dark, compact bool) *theme.AntToken {
	var opts []theme.Option
	// 静态截图没有键盘交互，按 antd :focus-visible 语义禁用焦点环，
	// 否则根容器常驻一圈主色描边，污染整批截图。
	opts = append(opts, theme.WithoutFocusRing())
	if compact {
		opts = append(opts, theme.WithDensity(theme.Compact))
	}
	if dark {
		return theme.Dark(opts...)
	}
	return theme.Light(opts...)
}

// ShotResult 是单张截图的自动校验结果，截图后立即计算。
type ShotResult struct {
	Group        string `json:"group"`
	Title        string `json:"title"`
	Theme        string `json:"theme"`
	File         string `json:"file"`
	Blank        bool   `json:"blank"`        // 画面近乎纯色，疑似渲染失败
	EdgeClip     bool   `json:"edgeClip"`     // 内容顶到画布边缘，疑似被裁切
	UniqueColors int    `json:"uniqueColors"` // 采样后的去重颜色数
	Notes        string `json:"notes,omitempty"`
}

// analyze 对渲染帧做自动校验：空白检测 + 边缘裁切检测。
func analyze(img image.Image) (blank, edgeClip bool, colors int) {
	b := img.Bounds()
	seen := map[string]int{}
	// 采样：每 4 像素取一点，足够判断空白与裁切。
	for y := b.Min.Y; y < b.Max.Y; y += 4 {
		for x := b.Min.X; x < b.Max.X; x += 4 {
			r, g, bl, _ := img.At(x, y).RGBA()
			seen[fmt.Sprintf("%04x%04x%04x", r>>8, g>>8, bl>>8)]++
		}
	}
	colors = len(seen)
	// 纯空白帧只含 1 种颜色；稀疏示例（如分割线）可能有 2~3 种，不算失败。
	blank = colors <= 1

	// 裁切检测：先取最外圈 2 像素的众数颜色作为底色参照，
	// 再统计与底色差异明显的像素占比（众数可避免角落异色环导致误报）。
	type rgb struct{ r, g, b uint32 }
	modeCount := map[rgb]int{}
	mode := rgb{}
	modeN := 0
	for x := b.Min.X; x < b.Max.X; x += 2 {
		for _, y := range []int{b.Min.Y, b.Min.Y + 1, b.Max.Y - 2, b.Max.Y - 1} {
			r, g, bl, _ := img.At(x, y).RGBA()
			k := rgb{r >> 8, g >> 8, bl >> 8}
			modeCount[k]++
			if modeCount[k] > modeN {
				mode, modeN = k, modeCount[k]
			}
		}
	}
	diff := func(r, g, bl uint32) int {
		dr := int(r>>8) - int(mode.r)
		dg := int(g>>8) - int(mode.g)
		db := int(bl>>8) - int(mode.b)
		if dr < 0 {
			dr = -dr
		}
		if dg < 0 {
			dg = -dg
		}
		if db < 0 {
			db = -db
		}
		return dr + dg + db
	}
	border, borderDiff := 0, 0
	for x := b.Min.X; x < b.Max.X; x += 2 {
		for _, y := range []int{b.Min.Y, b.Min.Y + 1, b.Max.Y - 2, b.Max.Y - 1} {
			border++
			r, g, bl, _ := img.At(x, y).RGBA()
			if diff(r, g, bl) > 90 {
				borderDiff++
			}
		}
	}
	for y := b.Min.Y; y < b.Max.Y; y += 2 {
		for _, x := range []int{b.Min.X, b.Min.X + 1, b.Max.X - 2, b.Max.X - 1} {
			border++
			r, g, bl, _ := img.At(x, y).RGBA()
			if diff(r, g, bl) > 90 {
				borderDiff++
			}
		}
	}
	edgeClip = border > 0 && float64(borderDiff)/float64(border) > 0.02
	return
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
	results := []ShotResult{}
	for _, d := range demos() {
		if *group != "" && d.Name != *group {
			continue
		}
		for i, s := range d.Samples {
			s := s
			view := func(c *ui.Context) {
				theme.Apply(c, tk)
				// 示例本身不带页面留白（Gallery 由外层补白），截图时包一层
				// 16px 内边距，避免内容顶到画布边缘、被误判为裁切。
				host := ui.Column(c).Fill().Padding(16)
				if s.Fill {
					// 布局型示例：撑满画布而非居中。
					host.Children(func() {
						s.View(c)
					})
					return
				}
				host.Center().Children(func() {
					s.View(c)
				})
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

			// 截图后立即全量校验
			res := ShotResult{Group: d.Name, Title: s.Title, Theme: *themeName, File: name}
			res.Blank, res.EdgeClip, res.UniqueColors = analyze(img)
			switch {
			case res.Blank:
				res.Notes = "画面近乎纯色，疑似渲染失败"
			case res.EdgeClip:
				res.Notes = "内容顶到画布边缘，疑似被裁切"
			}
			results = append(results, res)
			if res.Blank || res.EdgeClip {
				log.Printf("[校验不通过] %s/%s：%s", d.Name, s.Title, res.Notes)
			}
			n++
		}
	}

	vj, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	vpath := filepath.Join(*out, "validate.json")
	if err := os.WriteFile(vpath, vj, 0o644); err != nil {
		log.Fatal(err)
	}
	failed := 0
	for _, r := range results {
		if r.Blank || r.EdgeClip {
			failed++
		}
	}
	log.Printf("已输出 %d 张截图到 %s；自动校验：%d 张通过，%d 张不通过（明细见 %s）",
		n, *out, n-failed, failed, vpath)
}
