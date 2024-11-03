// Package osexitchecker реализует собственный анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакета main.
package osexitchecker

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "exitcheck",                         // имя анализатора. Нужно указать валидный Go-идентификатор, так как он может использоваться в параметрах командной строки, URL и т. д.
	Doc:  "check for os.Exit in main package", // текст с описанием работы анализатора. Этот текст будет отображаться по команде help, поэтому его нужно сделать многострочным и описать в нём все флаги анализатора.
	Run:  run,                                 // функция, которая отвечает за анализ исходного кода.
}

func run(pass *analysis.Pass) (interface{}, error) {

	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)

			if ok {
				// Проверяем, является ли вызов функцией os.Exit
				if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					if ident, ok := selExpr.X.(*ast.Ident); ok && ident.Name == "os" && selExpr.Sel.Name == "Exit" {
						pass.Reportf(callExpr.Pos(), "os.Exit is not allowed in main package")
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
