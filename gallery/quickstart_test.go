package main

import (
	"os"
	"strings"
	"testing"
)

// TestQuickStartMatchesREADME 断言首页「快速开始」的代码与 README 的
// 代码块逐字一致（同一份示例，不能各写各的）。改 README 时同步
// gallery/quickstart.go。
func TestQuickStartMatchesREADME(t *testing.T) {
	data, err := os.ReadFile("../README.md")
	if err != nil {
		t.Fatalf("读取 README 失败: %v", err)
	}
	s := string(data)
	i := strings.Index(s, "## 快速开始")
	if i < 0 {
		t.Fatal("README 缺少「快速开始」章节")
	}
	j := strings.Index(s[i:], "```go")
	if j < 0 {
		t.Fatal("「快速开始」缺少 go 代码块")
	}
	j += i + len("```go")
	k := strings.Index(s[j:], "```")
	if k < 0 {
		t.Fatal("「快速开始」代码块未闭合")
	}
	want := strings.Trim(s[j:j+k], "\n")
	if want != quickStartCode {
		// 给出首个差异行，便于定位漂移。
		wl := strings.Split(want, "\n")
		gl := strings.Split(quickStartCode, "\n")
		for n := 0; n < len(wl) && n < len(gl); n++ {
			if wl[n] != gl[n] {
				t.Fatalf("快速开始与 README 不一致（第 %d 行）：\nREADME: %q\nGallery: %q", n+1, wl[n], gl[n])
			}
		}
		t.Fatalf("快速开始与 README 行数不一致：README %d 行，Gallery %d 行", len(wl), len(gl))
	}
}

// TestQuickStartTruncates 首页只展示前 14 行并给出「完整示例见 README」提示。
func TestQuickStartTruncates(t *testing.T) {
	out := truncateCode(quickStartCode, 14, "完整示例见 README 快速开始")
	if !strings.Contains(out, "完整示例见 README 快速开始") {
		t.Fatal("截断后缺少提示文字")
	}
	if n := len(strings.Split(out, "\n")); n > 16 {
		t.Fatalf("截断后 %d 行，超过 14 行 + 提示", n)
	}
}
