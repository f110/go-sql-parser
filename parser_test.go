package parser

import (
	"io"
	"log"
	"reflect"
	"testing"
)

func assertSelectList(t *testing.T, expected SelectList, actual SelectList, i int) {
	if len(expected) != len(actual) {
		t.Fatalf("tokens %d: Expected length %d, but actual %d", i, len(expected), len(actual))
	}

	for k, c := range expected {
		if c.Column != actual[k].Column {
			t.Fatalf("tokens %d: Expected column %s, but got %s", i, c.Column, actual[k].Column)
		}
		if c.Alias != "" && c.Alias != actual[k].Alias {
			t.Fatalf("tokens %d: Expected alias %s, but got %s", i, c.Alias, actual[k].Alias)
		}
	}
}

func assertFromClause(t *testing.T, expected FromClause, actual FromClause, i int) {
	if len(expected.Table) != len(actual.Table) {
		t.Fatalf("tokens %d: Expected length %d, but actual %d", i, len(expected.Table), len(actual.Table))
	}

	assertJoinedTable(t, expected.Join, actual.Join, i)
}

func assertJoinedTable(t *testing.T, expected JoinedTable, actual JoinedTable, i int) {
	if len(expected.Table) != len(actual.Table) {
		t.Fatalf("tokens %d: Expected length %d, but actual %d", i, len(expected.Table), len(actual.Table))
	}

	assertTokens(t, expected.Type, actual.Type)
	assertExpr(t, expected.Cond, actual.Cond, i)
}

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
		parser := Parser{}
		for i, c := range TestSelectQuery {
			t.Logf("testing: %s", c.Query)
			tr := NewTokensReader(c.Tokens)
			p, err := parser.Parse(tr)
			if err != nil && err != io.EOF {
				t.Fatalf("Failed parse tokens %d (%v): %v", i, c.Tokens, err)
			}
			s, ok := p.(*Select)
			if ok == false {
				t.Fatalf("tokens %d, Expected Select but not", i)
			}

			ast := c.Ast.(Select)
			assertSelectList(t, ast.SelectList, s.SelectList, i)
			assertFromClause(t, ast.Table.From, s.Table.From, i)
			assertWhereClause(t, ast.Table.Where, s.Table.Where, i)
			assertGroupByClause(t, ast.Table.GroupBy, s.Table.GroupBy, i)
			assertHavingClause(t, ast.Table.Having, s.Table.Having, i)
			assertOrderByClause(t, ast.OrderBy, s.OrderBy, i)
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
