package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "checks log messages for style, language, and sensitive data",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// defining logger: log/slog/zap
			switch fun := call.Fun.(type) {
			case *ast.SelectorExpr:
				if id, ok := fun.X.(*ast.Ident); ok {
					if isLogger(id.Name) {
						if len(call.Args) == 0 {
							return false
						}

						// first arg
						arg := call.Args[0]
						lit, ok := arg.(*ast.BasicLit)
						if !ok || lit.Kind != token.STRING {
							return false
						}

						msg, _ := strconv.Unquote(lit.Value)

						if first, _ := utf8.DecodeRuneInString(msg); unicode.IsUpper(first) {
							pass.Reportf(arg.Pos(), "log message should start with a lowercase letter")
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

// package name check
func isLogger(pkg string) bool {
	switch pkg {
	case "log", "slog", "zap":
		return true
	default:
		return false
	}
}
