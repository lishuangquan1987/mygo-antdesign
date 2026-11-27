package ant

import (
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/chart"
	"github.com/egoist/mygo/ui/editor"
	"github.com/egoist/mygo/ui/markdown"
	"github.com/egoist/mygo/ui/netimage"
	"github.com/egoist/mygo/ui/richtext"
	"github.com/egoist/mygo/ui/video"
)

// 本文件适配 MyGo 的 UI 子包（ui/markdown、ui/editor、ui/chart、ui/richtext、
// ui/netimage、ui/video），让它们与 ant 控件共享同一套主题圆角/字号。
// 子包保留自身交互与渲染，不重新实现。

// MarkdownOptions 是 markdown.Options 的别名（见 MyGo 文档）。
type MarkdownOptions = markdown.Options

// Markdown 渲染 Markdown 文档（MyGo ui/markdown）：锚点、代码块、
// 复制为 Markdown 等交互由原语提供，主题取自当前窗口。
func Markdown(c *ui.Context, src string, o *MarkdownOptions) ui.Element {
	return markdown.View(c, src, o)
}

// CodeEditorState 是 editor.State 的别名。
type CodeEditorState = editor.State

// CodeEditorBuffer 是 editor.Buffer 的别名。
type CodeEditorBuffer = editor.Buffer

// CodeEditor 代码编辑器（MyGo ui/editor）：语法高亮、多光标、查找替换、
// 折叠与补全由原语提供，外层按主题圆角裁切。
func CodeEditor(c *ui.Context, st *CodeEditorState) ui.Element {
	return editor.View(c, st).Radius(c.Theme().Radius)
}

// EditorDiffState 是 editor.DiffState 的别名。
type EditorDiffState = editor.DiffState

// DiffView 差异视图（MyGo ui/editor）。
func DiffView(c *ui.Context, st *EditorDiffState) ui.Element {
	return editor.DiffView(c, st).Radius(c.Theme().Radius)
}

// ChartAxes 是 chart.Axes 的别名。
type ChartAxes = chart.Axes

// ChartTick 是 chart.Tick 的别名。
type ChartTick = chart.Tick

// ChartPlot 坐标图（MyGo ui/chart）：轴刻度与绘图区由原语排布，
// build 负责绘制网格线与数据标记。
func ChartPlot(c *ui.Context, axes ChartAxes, build func(area ui.Element)) ui.Element {
	return chart.Plot(c, axes, build).Radius(c.Theme().Radius)
}

// RemoteImageOptions 是 netimage.Options 的别名。
type RemoteImageOptions = netimage.Options

// RemoteImage 网络图片（MyGo ui/netimage）：占位图 → 加载完成淡入，
// 支持 Fit / MaxSide / NoFade 等选项。
func RemoteImage(c *ui.Context, url string, o ...RemoteImageOptions) ui.Element {
	return netimage.Image(c, url, o...).Radius(c.Theme().Radius)
}

// RichTextDocument 是 richtext.Document 的别名。
type RichTextDocument = richtext.Document

// RichTextEditorOptions 是 richtext.EditorOptions 的别名。
type RichTextEditorOptions = richtext.EditorOptions

// RichTextEditor 富文本编辑器（MyGo ui/richtext）：Markdown 快捷键、
// 选区格式化、任务列表等由原语提供。
func RichTextEditor(c *ui.Context, doc *RichTextDocument, options ...RichTextEditorOptions) ui.Element {
	return richtext.Editor(c, doc, options...).Radius(c.Theme().Radius)
}

// OpenVideo 从文件打开视频（MyGo ui/video），交给 ant.Video 显示：
//
//	player, err := ant.OpenVideo("waves.mp4")
//	ant.Video(c, player).FillWidth()
func OpenVideo(path string) (*ui.VideoPlayer, error) {
	return video.Open(path)
}
