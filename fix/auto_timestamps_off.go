package fix

import (
	"strings"

	"github.com/gobuffalo/plush/ast"
	"github.com/gobuffalo/plush/parser"
)

// AutoTimestampsOff adds a t.Timestamps() statement to fizz migrations
// when they still use the implicit auto-timestamp old fizz feature.
func AutoTimestampsOff(content string) (string, error) {
	var p *ast.Program
	var err error
	if p, err = parser.Parse("<% " + content + "%>"); err != nil {
		return "", err
	}

	var pt *ast.Program
	if pt, err = parser.Parse("<% t.Timestamps() %>"); err != nil {
		return "", err
	}
	ts := pt.Statements[0].(*ast.ExpressionStatement)

	for _, s := range p.Statements {
		stmt := s.(*ast.ExpressionStatement)
		if function, ok := stmt.Expression.(*ast.CallExpression); ok {
			if function.Function.TokenLiteral() == "create_table" {
				args := function.Arguments
				enableTimestamps := true
				if len(args) > 1 {
					if v, ok := args[1].(*ast.HashLiteral); ok {
						if strings.Contains(v.String(), `"timestamps": false`) {
							enableTimestamps = false
						}
					}
				}
				for _, bs := range function.Block.Statements {
					bstmt := bs.(*ast.ExpressionStatement)
					if f, ok := bstmt.Expression.(*ast.CallExpression); ok {
						fs := f.Function.String()
						if fs == "t.DisableTimestamps" || fs == "t.Timestamps" {
							enableTimestamps = false
						}
					}
				}
				if enableTimestamps {
					function.Block.Statements = append(function.Block.Statements, ts)
				}
			}
		}
	}

	return p.String(), nil
}
