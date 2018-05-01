package parser

import (
	"io"
	"log"
	"reflect"
	"testing"
)

func assertGroupByClause(t *testing.T, expected GroupByClause, actual GroupByClause, i int) {
	if len(expected) != len(actual) {
		t.Fatalf("tokens %d: Expected length %d, but actual %d", i, len(expected), len(actual))
	}

	for i, e := range expected {
		if e.Value != actual[i].Value {
			t.Fatalf("tokens %d: Expected %s, but got %s", i, e.Value, actual[i].Value)
		}
	}
}

func assertHavingClause(t *testing.T, expected HavingClause, actual HavingClause, i int) {
	assertExpr(t, expected.Cond, actual.Cond, i)
}

func assertOrderByClause(t *testing.T, expected OrderByClause, actual OrderByClause, i int) {
	if len(expected) != len(actual) {
		t.Fatalf("tokens %d: Expected length %d, but actual %d", i, len(expected), len(actual))
	}
	for i, e := range expected {
		if e.Key.Value != actual[i].Key.Value {
			t.Fatalf("tokens %d: Expected key %s, but got %s", i, e.Key.Value, actual[i].Key.Value)
		}
		if e.Order.Type != ILLEGAL && e.Order.Type != actual[i].Order.Type {
			t.Fatalf("tokens %d: Expected order %v, but actual %v", i, e.Order.Type, actual[i].Key.Type)
		}
	}
}

func assertExpr(t *testing.T, expected Expr, actual Expr, i int) {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		t.Fatalf("tokens %d: Expected %v, but actual %v", i, reflect.TypeOf(expected), reflect.TypeOf(actual))
	}
	if v, ok := expected.(*ValueExpr); ok {
		assertValueExpr(t, v, actual.(*ValueExpr))
	}
	if v, ok := expected.(*BooleanTerm); ok {
		assertBooleanTerm(t, v, actual.(*BooleanTerm))
	}
}

func assertWhereClause(t *testing.T, expected WhereClause, actual WhereClause, i int) {
	assertExpr(t, expected.Cond, actual.Cond, i)
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

func assertTokens(t *testing.T, expected Tokens, actual Tokens) {
	for i, e := range expected {
		assertToken(t, e, actual[i])
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
							LeftValue:  []Token{{Type: IDENT, Value: "id"}},
							RightValue: []Token{{Type: INT, IntValue: 1}},
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
			{
				Tokens: []Token{ // select * from users where id > 10 and age > 20
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
					{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 14, Column: 5}},
					{Type: WHERE, Position: Position{Line: 1, Offset: 20, Column: 5}},
					{Type: IDENT, Value: "id", Position: Position{Line: 1, Offset: 26, Column: 2}},
					{Type: GTR, Position: Position{Line: 1, Offset: 29, Column: 1}},
					{Type: INT, IntValue: 10, Position: Position{Line: 1, Offset: 31, Column: 2}},
					{Type: AND, Position: Position{Line: 1, Offset: 34, Column: 3}},
					{Type: IDENT, Value: "age", Position: Position{Line: 1, Offset: 38, Column: 3}},
					{Type: GTR, Position: Position{Line: 1, Offset: 42, Column: 1}},
					{Type: INT, IntValue: 20, Position: Position{Line: 1, Offset: 44, Column: 2}},
					{Type: EOF, Position: Position{Line: 1, Offset: 46}},
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
									Operator:   Token{Type: GTR},
									LeftValue:  []Token{{Type: IDENT, Value: "id"}},
									RightValue: []Token{{Type: INT, IntValue: 1}},
								},
								Right: &ValueExpr{
									Operator:   Token{Type: GTR},
									LeftValue:  []Token{{Type: IDENT, Value: "age"}},
									RightValue: []Token{{Type: INT, IntValue: 20}},
								},
							},
						},
					},
					SelectList: []Token{{Type: ASTERISK}},
				},
			},
			{
				Tokens: []Token{ // select * from users order by created_date desc, rank
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
					{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 14, Column: 5}},
					{Type: ORDERBY, Position: Position{Line: 1, Offset: 20, Column: 8}},
					{Type: IDENT, Value: "created_date", Position: Position{Line: 1, Offset: 29, Column: 12}},
					{Type: DESC, Position: Position{Line: 1, Offset: 42, Column: 4}},
					{Type: COMMA, Position: Position{Line: 1, Offset: 46, Column: 1}},
					{Type: IDENT, Value: "rank", Position: Position{Line: 1, Offset: 48, Column: 4}},
					{Type: EOF, Position: Position{Line: 1, Offset: 52}},
				},
				Select: Select{
					Table: TableExpression{
						From: FromClause{
							Table: []Token{{Type: IDENT, Value: "users"}},
						},
						Where: WhereClause{},
					},
					OrderBy: []*SortSpecification{{Key: Token{Type: IDENT, Value: "created_date"}, Order: Token{Type: DESC}}, {Key: Token{Type: IDENT, Value: "rank"}}},
				},
			},
			{
				Tokens: []Token{ // select * from users group by group_id
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
					{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 14, Column: 5}},
					{Type: GROUPBY, Position: Position{Line: 1, Offset: 20, Column: 8}},
					{Type: IDENT, Value: "group_id", Position: Position{Line: 1, Offset: 29, Column: 8}},
					{Type: EOF, Position: Position{Line: 1, Offset: 37}},
				},
				Select: Select{
					Table: TableExpression{
						From: FromClause{
							Table: []Token{{Type: IDENT, Value: "users"}},
						},
						Where:   WhereClause{},
						GroupBy: GroupByClause([]Token{{Type: IDENT, Value: "group_id"}}),
					},
				},
			},
			{
				Tokens: []Token{ // select * from users group by group_id having group_id > 10
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
					{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 14, Column: 5}},
					{Type: GROUPBY, Position: Position{Line: 1, Offset: 20, Column: 8}},
					{Type: IDENT, Value: "group_id", Position: Position{Line: 1, Offset: 29, Column: 8}},
					{Type: HAVING, Position: Position{Line: 1, Offset: 38, Column: 6}},
					{Type: IDENT, Value: "group_id", Position: Position{Line: 1, Offset: 45, Column: 8}},
					{Type: GTR, Position: Position{Line: 1, Offset: 54, Column: 1}},
					{Type: INT, Position: Position{Line: 1, Offset: 56, Column: 2}},
					{Type: EOF, Position: Position{Line: 1, Offset: 58}},
				},
				Select: Select{
					Table: TableExpression{
						From: FromClause{
							Table: []Token{{Type: IDENT, Value: "users"}},
						},
						Where:   WhereClause{},
						GroupBy: GroupByClause([]Token{{Type: IDENT, Value: "group_id"}}),
						Having: HavingClause{
							Cond: &ValueExpr{
								Operator:   Token{Type: GTR},
								LeftValue:  []Token{{Type: IDENT, Value: "group_d"}},
								RightValue: []Token{{Type: INT, IntValue: 10}},
							},
						},
					},
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

			assertTokens(t, Tokens(c.Select.Table.From.Table), Tokens(s.Table.From.Table))
			assertTokens(t, Tokens(c.Select.SelectList), Tokens(s.SelectList))
			assertWhereClause(t, c.Select.Table.Where, s.Table.Where, i)
			assertGroupByClause(t, c.Select.Table.GroupBy, s.Table.GroupBy, i)
			assertHavingClause(t, c.Select.Table.Having, s.Table.Having, i)
			assertOrderByClause(t, c.Select.OrderBy, s.OrderBy, i)
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
	var tokens = []Token{{Type: IDENT, Value: "id"}, {Type: EQUAL}, {Type: INT, IntValue: 1}, {Type: AND}, {Type: IDENT, Value: "age"}, {Type: EQUAL}, {Type: INT, IntValue: 20}}

	parser := Parser{}
	res, _ := parser.parseSearchCondition(NewTokensReader(tokens))
	log.Print(res.(*BooleanTerm))
}
