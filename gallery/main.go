// Gallery 应用入口：注册控件示例并打开窗口。
package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	"github.com/lishuangquan1987/mygo-antdesign/samples/button"
	"github.com/lishuangquan1987/mygo-antdesign/samples/display"
	"github.com/lishuangquan1987/mygo-antdesign/samples/divider"
	extrademo "github.com/lishuangquan1987/mygo-antdesign/samples/extrademo"
	"github.com/lishuangquan1987/mygo-antdesign/samples/feedback"
	"github.com/lishuangquan1987/mygo-antdesign/samples/ide"
	"github.com/lishuangquan1987/mygo-antdesign/samples/input"
	"github.com/lishuangquan1987/mygo-antdesign/samples/layout"
	"github.com/lishuangquan1987/mygo-antdesign/samples/nav"
	"github.com/lishuangquan1987/mygo-antdesign/samples/selectdemo"
	"github.com/lishuangquan1987/mygo-antdesign/samples/table"
	"github.com/lishuangquan1987/mygo-antdesign/samples/typography"
)

// version 由发布流水线注入：go build -ldflags "-X main.version=v0.1.0"。
var version = "dev"

func main() {
	log.Printf("mygo-antdesign Gallery %s", version)

	// 注册控件示例：新增控件时在此追加一行。
	samples.Register(button.Demo())
	samples.Register(typography.Demo())
	samples.Register(divider.Demo())
	samples.Register(layout.Demo())
	samples.Register(input.Demo())
	samples.Register(selectdemo.Demo())
	samples.Register(display.Demo())
	samples.Register(table.Demo())
	samples.Register(feedback.Demo())
	samples.Register(nav.Demo())
	samples.Register(extrademo.Demo())
	samples.Register(ide.Demo())

	g := NewGallery()
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:   "mygo-antdesign Gallery",
			Content: ui.View(g.View),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
