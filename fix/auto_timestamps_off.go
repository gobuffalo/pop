package fix

import (
	"fmt"

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
				for _, a := range args {
					fmt.Println(a)
				}
				for _, bs := range function.Block.Statements {
					bstmt := bs.(*ast.ExpressionStatement)
					if f, ok := bstmt.Expression.(*ast.CallExpression); ok {
						fmt.Printf("%T\n", f)
					}
				}
				function.Block.Statements = append(function.Block.Statements, ts)
			}
		}
	}

	fmt.Println(p.Statements)

	return p.String(), nil
}
