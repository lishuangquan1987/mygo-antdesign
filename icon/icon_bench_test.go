package icon

import "testing"

// BenchmarkParseSVGCached 测缓存命中路径（Gallery 一帧数十次取图标）。
// 未缓存时每次都要完整解析 SVG 并生成新的 *SVG 指针（栅格化缓存随之失效）。
func BenchmarkParseSVGCached(b *testing.B) {
	if _, ok := parseSVG(Search); !ok {
		b.Fatal("Search 图标解析失败")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := parseSVG(Search); !ok {
			b.Fatal("缓存命中失败")
		}
	}
}
