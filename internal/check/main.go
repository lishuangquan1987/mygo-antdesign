// 命令 check 是样式铁律的静态检查工具（CI 门禁）：
//
//	go run ./internal/check
//
// 扫描 widgets/ 与 extra/ 的非测试 Go 文件，禁止：
//   - 裸颜色：ui.Hex( / ui.RGB( / ui.RGBA(（颜色必须来自主题）
//   - 裸间距：Padding/Gap/Margin/Size( 后直接跟数字字面量
//     （间距必须走 t.Space / t.Rem）
//
// 豁免：theme/（Token 真源）、icon/（SVG 数据）、*_test.go、goldens。
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// 扫描的目录；缺失的目录（如 extra/ 尚未创建）自动跳过。
var targets = []string{"widgets", "extra"}

// 禁止的裸颜色构造。
var rawColorFns = map[string]bool{
	"ui.Hex": true, "ui.RGB": true, "ui.RGBA": true, "ui.Oklch": true,
}

// 间距方法与最小参数个数：Padding(v ...)、Size(w,h) 全查；
// Gap 单参数方法由裸数字规则覆盖。
var spacingMethods = map[string]bool{
	"Padding": true, "PaddingX": true, "PaddingY": true,
	"PaddingStart": true, "PaddingEnd": true,
	"Margin": true, "MarginX": true, "MarginY": true,
	"Gap": true, "GapX": true, "GapY": true,
	"Size": true, "Width": true, "Height": true,
	"Radius": true,
}

// 间距的合法取值调用前缀：t.Space / t.Rem / internal.* / 内部变量均放行，
// 只拦截「方法后直接跟非零数字字面量」这一类明显违规。0 豁免——
// Padding(0, …) 的 0 是「无边距」语义，与 1px 描边同属视觉基准例外。
func isRawSpacing(call *ast.CallExpr) bool {
	for _, arg := range call.Args {
		var lit *ast.BasicLit
		switch a := arg.(type) {
		case *ast.BasicLit:
			lit = a
		case *ast.UnaryExpr: // 负数字面量
			if l, ok := a.X.(*ast.BasicLit); ok {
				lit = l
			}
		}
		if lit != nil && (lit.Kind == token.FLOAT || lit.Kind == token.INT) {
			if lit.Value == "0" || lit.Value == "0.0" || lit.Value == "0." {
				continue // 0 = 无边距，豁免
			}
			return true
		}
	}
	return false
}

func selectorName(e ast.Expr) string {
	if s, ok := e.(*ast.SelectorExpr); ok {
		if x, ok := s.X.(*ast.Ident); ok {
			return x.Name + "." + s.Sel.Name
		}
	}
	return ""
}

func checkFile(path string) []string {
	var issues []string
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return []string{fmt.Sprintf("%s: 解析失败: %v", path, err)}
	}
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			// 裸颜色
			if sel := selectorName(call.Fun); rawColorFns[sel] {
				pos := fset.Position(call.Pos())
				issues = append(issues, fmt.Sprintf("%s:%d 裸颜色 %s(…)，颜色必须来自主题",
					pos.Filename, pos.Line, sel))
			}
			// 裸间距
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok && spacingMethods[sel.Sel.Name] {
				if isRawSpacing(call) {
					pos := fset.Position(call.Pos())
					issues = append(issues, fmt.Sprintf("%s:%d 裸间距 %s(%v)，间距必须走 t.Space/t.Rem",
						pos.Filename, pos.Line, sel.Sel.Name, call.Args))
				}
			}
			return true
		})
	}
	return issues
}

func main() {
	var all []string
	for _, dir := range targets {
		if _, err := os.Stat(dir); err != nil {
			continue // extra/ 等尚不存在
		}
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			all = append(all, checkFile(path)...)
			return nil
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}
	if len(all) > 0 {
		fmt.Println("样式铁律违规：")
		for _, s := range all {
			fmt.Println("  " + s)
		}
		os.Exit(1)
	}
	fmt.Println("样式铁律检查通过（widgets/extra 无裸颜色与裸间距）")
}
