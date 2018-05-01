package parser

import (
	"io"
	"log"
	"reflect"
	"testing"
)

func assertTokens(expected []Token, actual []Token) bool {
	for i, e := range expected {
		if e.Type != actual[i].Type {
			return false
		}

		switch e.Type {
		case IDENT:
			if e.Value != actual[i].Value {
				return false
			}
		case INT:
			if e.IntValue != actual[i].IntValue {
				return false
			}
		}
	}
	if len(expected) != len(actual) {
		return false
	}

	return true
}

func assertWhereClause(t *testing.T, expected WhereClause, actual WhereClause, i int) {
	if reflect.TypeOf(expected.Cond) != reflect.TypeOf(actual.Cond) {
		t.Fatalf("tokens %d: Expected %v, but actual %v", i, reflect.TypeOf(expected.Cond), reflect.TypeOf(actual.Cond))
	}
	if v, ok := expected.Cond.(*ValueExpr); ok {
		assertValueExpr(t, v, actual.Cond.(*ValueExpr))
	}
	if v, ok := expected.Cond.(*BooleanTerm); ok {
		assertBooleanTerm(t, v, actual.Cond.(*BooleanTerm))
	}
}

func assertValueExpr(t *testing.T, expected *ValueExpr, actual *ValueExpr) {
	if v, ok := expected.LeftValue.(Tokens); ok {
		for i, token := range v {
			assertToken(t, token, actual.LeftValue.(Tokens)[i])
		}
	}

	if v, ok := expected.RightValue.(Tokens); ok {
		for i, token := range v {
			assertToken(t, token, actual.RightValue.(Tokens)[i])
		}
	}
}

func assertBooleanTerm(t *testing.T, expected *BooleanTerm, actual *BooleanTerm) {
	if expected.Boolean.Type != actual.Boolean.Type {
		t.Fatalf("Expected %v but actual %v", expected.Boolean.Type, actual.Boolean.Type)
	}

	if v, ok := expected.Left.(*ValueExpr); ok {
		assertValueExpr(t, v, actual.Left.(*ValueExpr))
	}
	if v, ok := expected.Left.(*BooleanTerm); ok {
		assertBooleanTerm(t, v, actual.Left.(*BooleanTerm))
	}

	if v, ok := expected.Right.(*ValueExpr); ok {
		assertValueExpr(t, v, actual.Right.(*ValueExpr))
	}
	if v, ok := expected.Right.(*BooleanTerm); ok {
		assertBooleanTerm(t, v, actual.Right.(*BooleanTerm))
	}
}

func assertToken(t *testing.T, expected Token, actual Token) {
	switch expected.Type {
	case IDENT:
		if actual.Type != IDENT {
			t.Fatalf("Expected Type ident but got %v", actual.Type)
		}
		if expected.Value != actual.Value {
			t.Fatalf("Expected %s but actual %s", expected.Value, actual.Value)
		}
	case INT:
		if actual.Type != INT {
			t.Fatalf("Expected Type int but got %v", actual.Type)
		}
		if expected.IntValue != actual.IntValue {
			t.Fatalf("Expected %d but actual %d", expected.IntValue, actual.IntValue)
		}
	}
}

func TestParser_Parse(t *testing.T) {
	t.Run("Select", func(t *testing.T) {
		cases := []struct {
			Tokens []Token
			Select Select
		}{
			{
				Tokens: []Token{ // select * from test
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
					{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
					{Type: IDENT, Value: "test", Position: Position{Line: 1, Offset: 14, Column: 4}},
					{Type: EOF, Position: Position{Line: 1, Offset: 18}},
				},
				Select: Select{Table: TableExpression{From: FromClause{Table: []Token{{Type: IDENT, Value: "test"}}}}, SelectList: []Token{{Type: ASTERISK}}},
			},
			{
				Tokens: []Token{ // select * from `test`
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
					{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
					{Type: IDENT, Value: "`test`", Position: Position{Line: 1, Offset: 14, Column: 6}},
					{Type: EOF, Position: Position{Line: 1, Offset: 20}},
				},
				Select: Select{Table: TableExpression{From: FromClause{Table: []Token{{Type: IDENT, Value: "`test`"}}}}, SelectList: []Token{{Type: ASTERISK}}},
			},
			{
				Tokens: []Token{ // SELECT id,age from users
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: IDENT, Value: "id", Position: Position{Line: 1, Offset: 7, Column: 2}},
					{Type: COMMA, Position: Position{Line: 1, Offset: 9, Column: 1}},
					{Type: IDENT, Value: "age", Position: Position{Line: 1, Offset: 10, Column: 3}},
					{Type: FROM, Position: Position{Line: 1, Offset: 14, Column: 4}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 19, Column: 5}},
					{Type: EOF, Position: Position{Line: 1, Offset: 24}},
				},
				Select: Select{Table: TableExpression{From: FromClause{Table: []Token{{Type: IDENT, Value: "users"}}}}, SelectList: []Token{{Type: IDENT, Value: "id"}, {Type: IDENT, Value: "age"}}},
			},
			{
				Tokens: []Token{ // select * from users where id = 1
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
					{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 14, Column: 5}},
					{Type: WHERE, Position: Position{Line: 1, Offset: 20, Column: 5}},
					{Type: IDENT, Value: "id", Position: Position{Line: 1, Offset: 26, Column: 2}},
					{Type: EQUAL, Position: Position{Line: 1, Offset: 29, Column: 1}},
					{Type: INT, IntValue: 1, Position: Position{Line: 1, Offset: 31, Column: 1}},
					{Type: EOF, Position: Position{Line: 1, Offset: 27}},
				},
				Select: Select{
					Table: TableExpression{
						From: FromClause{
							Table: []Token{{Type: IDENT, Value: "users"}},
						},
						Where: WhereClause{Cond: &ValueExpr{
							Operator:   Token{Type: EQUAL},
							LeftValue:  Tokens([]Token{{Type: IDENT, Value: "id"}}),
							RightValue: Tokens([]Token{{Type: INT, IntValue: 1}}),
						}},
					},
					SelectList: []Token{{Type: ASTERISK}}},
			},
			{
				Tokens: []Token{ // select * from users where id = 1 and age = 20
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
					{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 14, Column: 5}},
					{Type: WHERE, Position: Position{Line: 1, Offset: 20, Column: 5}},
					{Type: IDENT, Value: "id", Position: Position{Line: 1, Offset: 26, Column: 2}},
					{Type: EQUAL, Position: Position{Line: 1, Offset: 29, Column: 1}},
					{Type: INT, IntValue: 1, Position: Position{Line: 1, Offset: 31, Column: 1}},
					{Type: AND, Position: Position{Line: 1, Offset: 33, Column: 3}},
					{Type: IDENT, Value: "age", Position: Position{Line: 1, Offset: 37, Column: 3}},
					{Type: EQUAL, Position: Position{Line: 1, Offset: 41, Column: 1}},
					{Type: INT, IntValue: 20, Position: Position{Line: 1, Offset: 43, Column: 2}},
					{Type: EOF, Position: Position{Line: 1, Offset: 45}},
				},
				Select: Select{
					Table: TableExpression{
						From: FromClause{
							Table: []Token{{Type: IDENT, Value: "users"}},
						},
						Where: WhereClause{
							Cond: &BooleanTerm{
								Boolean: Token{Type: AND},
								Left: &ValueExpr{
									Operator:   Token{Type: EQUAL},
									LeftValue:  []Token{{Type: IDENT, Value: "id"}},
									RightValue: []Token{{Type: INT, IntValue: 1}},
								},
								Right: &ValueExpr{
									Operator:   Token{Type: EQUAL},
									LeftValue:  []Token{{Type: IDENT, Value: "age"}},
									RightValue: []Token{{Type: INT, IntValue: 20}},
								},
							},
						},
					},
					SelectList: []Token{{Type: ASTERISK}},
				},
			},
		}

		parser := Parser{}
		for i, c := range cases {
			tr := NewTokensReader(c.Tokens)
			p, err := parser.Parse(tr)
			if err != nil && err != io.EOF {
				t.Fatalf("Failed parse tokens %d (%v): %v", i, c.Tokens, err)
			}
			s, ok := p.(*Select)
			if ok == false {
				t.Fatalf("tokens %d, Expected Select but not", i)
			}

			if !assertTokens(s.Table.From.Table, c.Select.Table.From.Table) {
				t.Fatalf("tokens %d, Expected from %v but actual %v", i, c.Select.Table.From.Table, s.Table.From.Table)
			}
			if !assertTokens(s.SelectList, c.Select.SelectList) {
				t.Fatalf("tokens %d, Expected colums %v but actual %v", i, c.Select.SelectList, s.SelectList)
			}
			assertWhereClause(t, c.Select.Table.Where, s.Table.Where, i)
		}
	})
}

func TestParser_parseExpression(t *testing.T) {
	//var tokens = []Token{{Type: IDENT, Value: "A"}, {Type: EQUAL}, {Type: IDENT, Value: "B"}, {Type: AND}, {Type: IDENT, Value: "C"}, {Type: EQUAL}, {Type: IDENT, Value: "D"}}
	//var tokens = []Token{{Type: IDENT, Value: "A"}, {Type: AND}, {Type: IDENT, Value: "B"}, {Type: AND}, {Type: LPAREN}, {Type: IDENT, Value: "C"}, {Type: OR}, {Type: IDENT, Value: "D"}, {Type: RPAREN}}
	//var tokens = []Token{{Type: IDENT, Value: "A"}, {Type: LSS}, {Type: IDENT, Value: "A'"}, {Type: AND}, {Type: IDENT, Value: "B"}, {Type: EQUAL}, {Type: IDENT, Value: "B'"}}

	//var tokens = []Token{{Type: IDENT, Value: "A"}, {Type: AND}, {Type: LPAREN}, {Type: IDENT, Value: "B"}, {Type: AND}, {Type: IDENT, Value: "C"}, {Type: RPAREN}} // invalid
	//var tokens = []Token{{Type: LPAREN}, {Type: IDENT, Value: "B"}, {Type: OR}, {Type: IDENT, Value: "C"}, {Type: RPAREN}, {Type: AND}, {Type: IDENT, Value: "D"}}
	//var tokens = []Token{{Type: IDENT, Value: "A"}, {Type: AND}, {Type: LPAREN}, {Type: LPAREN}, {Type: IDENT, Value: "B"}, {Type: OR}, {Type: IDENT, Value: "C"}, {Type: RPAREN}, {Type: AND}, {Type: IDENT, Value: "D"}, {Type: RPAREN}}
	var tokens = []Token{{Type: IDENT, Value: "id"}, {Type: EQUAL}, {Type: INT, IntValue: 1}}

	parser := Parser{}
	res, _ := parser.parseSearchCondition(NewTokensReader(tokens))
	log.Print(res.(*ValueExpr).LeftValue)
}
