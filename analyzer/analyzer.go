package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
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
						checkLogCall(pass, call)
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func checkLogCall(pass *analysis.Pass, call *ast.CallExpr) {
	if len(call.Args) == 0 {
		return
	}

	arg := call.Args[0]
	lit, ok := arg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return
	}

	msg, _ := strconv.Unquote(lit.Value)

	// Rule 1: the first letter is lowercase
	if first, _ := utf8.DecodeRuneInString(msg); unicode.IsUpper(first) {
		pass.Reportf(arg.Pos(), "log message should start with a lowercase letter")
	}

	// Rule 2: English
	for _, r := range msg {
		if !unicode.IsLetter(r) && r > unicode.MaxASCII {
			pass.Reportf(arg.Pos(), "log message should be in English")
			break
		}
	}

	// Rule 3: Special Characters and Emojis
	for _, r := range msg {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r)) {
			pass.Reportf(arg.Pos(), "log message should not contain special symbols or emojis")
			break
		}
	}

	// Rule 4: Sensitive Data
	keywords := []string{"password", "token", "api_key"}
	lower := strings.ToLower(msg)
	for _, k := range keywords {
		if strings.Contains(lower, k) {
			pass.Reportf(arg.Pos(), "log message should not contain sensitive data")
			break
		}
	}
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
