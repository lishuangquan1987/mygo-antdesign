// Gallery 应用入口：注册控件示例并打开窗口。
package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	"github.com/lishuangquan1987/mygo-antdesign/samples/button"
)

func main() {
	// 注册控件示例：新增控件时在此追加一行。
	samples.Register(button.Demo())

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
