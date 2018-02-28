package parser

import (
	"io"
	"strings"
	"testing"
)

func assertQuery(t *testing.T, query string, tokens []Token) {
	r := strings.NewReader(query)
	l := NewLexer(r)
	for i := 0; ; i++ {
		token, err := l.Scan()
		if err == io.EOF && token.Type != EOF {
			t.Error("Failed parse EOF")
		}
		if err == io.EOF {
			break
		}
		if token.Type != tokens[i].Type {
			t.Fatalf("Failed parse query %d (%s). expected %v got %v: %v", i, query, tokens[i].Type, token.Type, tokens)
		}
		if token.Position.Line != tokens[i].Position.Line {
			t.Fatalf("Failed parse query (%s). %s expected line is %d but actually %d", query, token.Type, tokens[i].Position.Line, token.Position.Line)
		}
		if token.Position.Offset != tokens[i].Position.Offset {
			t.Fatalf("Failed parse query (%s). %s expected Offset is %d but actually %d", query, token.Type, tokens[i].Position.Offset, token.Position.Offset)
		}
		if token.Position.Column != tokens[i].Position.Column {
			t.Fatalf("Failed parse query (%s). %s expected Column is %d but actually %d", query, token.Type, tokens[i].Position.Column, token.Position.Column)
		}
		if token.Type == IDENT && token.Value != tokens[i].Value {
			t.Fatalf("Failed parse query (%s). expected Value is \"%s\" but actually \"%s\"", query, tokens[i].Value, token.Value)
		}
	}
}

func TestLexer_Scan(t *testing.T) {
	t.Run("TestLexer_Scan_SELECT", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			Query string
			Token []Token
		}{
			{
				"select * from test",
				[]Token{
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
					{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
					{Type: IDENT, Value: "test", Position: Position{Line: 1, Offset: 14, Column: 4}},
					{Type: EOF, Position: Position{Line: 1, Offset: 18}},
				},
			},
			{
				"select * from `test`",
				[]Token{
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
					{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
					{Type: IDENT, Value: "`test`", Position: Position{Line: 1, Offset: 14, Column: 6}},
					{Type: EOF, Position: Position{Line: 1, Offset: 20}},
				},
			},
			{
				"select    *    from    test",
				[]Token{
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: ASTERISK, Position: Position{Line: 1, Offset: 10, Column: 1}},
					{Type: FROM, Position: Position{Line: 1, Offset: 15, Column: 4}},
					{Type: IDENT, Value: "test", Position: Position{Line: 1, Offset: 23, Column: 4}},
					{Type: EOF, Position: Position{Line: 1, Offset: 27}},
				},
			},
			{
				"SELECT id,age from users",
				[]Token{
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: IDENT, Value: "id", Position: Position{Line: 1, Offset: 7, Column: 2}},
					{Type: COMMA, Position: Position{Line: 1, Offset: 9, Column: 1}},
					{Type: IDENT, Value: "age", Position: Position{Line: 1, Offset: 10, Column: 3}},
					{Type: FROM, Position: Position{Line: 1, Offset: 14, Column: 4}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 19, Column: 5}},
					{Type: EOF, Position: Position{Line: 1, Offset: 24}},
				},
			},
			{
				"select * from users where id = 1",
				[]Token{
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
			},
			{
				"select * from users order by created_date desc",
				[]Token{
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
					{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 14, Column: 5}},
					{Type: ORDERBY, Position: Position{Line: 1, Offset: 20, Column: 8}},
					{Type: IDENT, Value: "created_date", Position: Position{Line: 1, Offset: 29, Column: 12}},
					{Type: DESC, Position: Position{Line: 1, Offset: 42, Column: 4}},
					{Type: EOF, Position: Position{Line: 1, Offset: 46}},
				},
			},
			{
				"select * from users group by group_id having group_id > 10",
				[]Token{
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
			},
			{
				"select * from users left inner join blog on users.id = blog.user_id",
				[]Token{
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
					{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 14, Column: 5}},
					{Type: LEFT, Position: Position{Line: 1, Offset: 20, Column: 4}},
					{Type: INNER, Position: Position{Line: 1, Offset: 25, Column: 5}},
					{Type: JOIN, Position: Position{Line: 1, Offset: 31, Column: 4}},
					{Type: IDENT, Value: "blog", Position: Position{Line: 1, Offset: 36, Column: 4}},
					{Type: ON, Position: Position{Line: 1, Offset: 41, Column: 2}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 44, Column: 5}},
					{Type: PERIOD, Position: Position{Line: 1, Offset: 49, Column: 1}},
					{Type: IDENT, Value: "id", Position: Position{Line: 1, Offset: 50, Column: 2}},
					{Type: EQUAL, Position: Position{Line: 1, Offset: 53, Column: 1}},
					{Type: IDENT, Value: "blog", Position: Position{Line: 1, Offset: 55, Column: 4}},
					{Type: PERIOD, Position: Position{Line: 1, Offset: 59, Column: 1}},
					{Type: IDENT, Value: "user_id", Position: Position{Line: 1, Offset: 60, Column: 7}},
					{Type: EOF, Position: Position{Line: 1, Offset: 67}},
				},
			},
			{
				"select * from users right outer join blog on users.id = blog.user_id",
				[]Token{
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
					{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 14, Column: 5}},
					{Type: RIGHT, Position: Position{Line: 1, Offset: 20, Column: 5}},
					{Type: OUTER, Position: Position{Line: 1, Offset: 26, Column: 5}},
					{Type: JOIN, Position: Position{Line: 1, Offset: 32, Column: 4}},
					{Type: IDENT, Value: "blog", Position: Position{Line: 1, Offset: 37, Column: 4}},
					{Type: ON, Position: Position{Line: 1, Offset: 42, Column: 2}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 45, Column: 5}},
					{Type: PERIOD, Position: Position{Line: 1, Offset: 50, Column: 1}},
					{Type: IDENT, Value: "id", Position: Position{Line: 1, Offset: 51, Column: 2}},
					{Type: EQUAL, Position: Position{Line: 1, Offset: 54, Column: 1}},
					{Type: IDENT, Value: "blog", Position: Position{Line: 1, Offset: 56, Column: 4}},
					{Type: PERIOD, Position: Position{Line: 1, Offset: 60, Column: 1}},
					{Type: IDENT, Value: "user_id", Position: Position{Line: 1, Offset: 61, Column: 7}},
					{Type: EOF, Position: Position{Line: 1, Offset: 68}},
				},
			},
			{
				"select * from users where id > 10 and age > 20",
				[]Token{
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
			},
			{
				"SELECT id as foo from users",
				[]Token{
					{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: IDENT, Value: "id", Position: Position{Line: 1, Offset: 7, Column: 2}},
					{Type: AS, Position: Position{Line: 1, Offset: 10, Column: 2}},
					{Type: IDENT, Value: "foo", Position: Position{Line: 1, Offset: 13, Column: 3}},
					{Type: FROM, Position: Position{Line: 1, Offset: 17, Column: 4}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 22, Column: 5}},
					{Type: EOF, Position: Position{Line: 1, Offset: 27}},
				},
			},
		}

		for _, c := range cases {
			assertQuery(t, c.Query, c.Token)
		}
	})

	t.Run("TestLexer_Scan_UPDATE", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			Query string
			Token []Token
		}{
			{
				"update users set name = \"test\" where id = 1",
				[]Token{
					{Type: UPDATE, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 7, Column: 5}},
					{Type: SET, Position: Position{Line: 1, Offset: 13, Column: 3}},
					{Type: IDENT, Value: "name", Position: Position{Line: 1, Offset: 17, Column: 4}},
					{Type: EQUAL, Position: Position{Line: 1, Offset: 22, Column: 1}},
					{Type: IDENT, Value: "\"test\"", Position: Position{Line: 1, Offset: 24, Column: 6}},
					{Type: WHERE, Position: Position{Line: 1, Offset: 31, Column: 5}},
					{Type: IDENT, Value: "id", Position: Position{Line: 1, Offset: 37, Column: 2}},
					{Type: EQUAL, Position: Position{Line: 1, Offset: 40, Column: 1}},
					{Type: INT, IntValue: 1, Position: Position{Line: 1, Offset: 42, Column: 1}},
					{Type: EOF, Position: Position{Line: 1, Offset: 43}},
				},
			},
		}

		for _, c := range cases {
			assertQuery(t, c.Query, c.Token)
		}
	})

	t.Run("TestLexer_Scan_CREATE", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			Query string
			Token []Token
		}{
			{
				"CREATE DATABASE users",
				[]Token{
					{Type: CREATE, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: DATABASE, Position: Position{Line: 1, Offset: 7, Column: 8}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 16, Column: 5}},
					{Type: EOF, Position: Position{Line: 1, Offset: 21}},
				},
			},
			{
				"CREATE table users (id int)",
				[]Token{
					{Type: CREATE, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: TABLE, Position: Position{Line: 1, Offset: 7, Column: 5}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 13, Column: 5}},
					{Type: LPAREN, Position: Position{Line: 1, Offset: 19, Column: 1}},
					{Type: IDENT, Value: "id", Position: Position{Line: 1, Offset: 20, Column: 2}},
					{Type: INTEGER, Position: Position{Line: 1, Offset: 23, Column: 3}},
					{Type: RPAREN, Position: Position{Line: 1, Offset: 26, Column: 1}},
					{Type: EOF, Position: Position{Line: 1, Offset: 27}},
				},
			},
			{
				"create assertion check_input check (not exist (select * from blog))",
				[]Token{
					{Type: CREATE, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: ASSERTION, Position: Position{Line: 1, Offset: 7, Column: 9}},
					{Type: IDENT, Value: "check_input", Position: Position{Line: 1, Offset: 17, Column: 11}},
					{Type: CHECK, Position: Position{Line: 1, Offset: 29, Column: 5}},
					{Type: LPAREN, Position: Position{Line: 1, Offset: 35, Column: 1}},
					{Type: NOT, Position: Position{Line: 1, Offset: 36, Column: 3}},
					{Type: EXIST, Position: Position{Line: 1, Offset: 40, Column: 5}},
					{Type: LPAREN, Position: Position{Line: 1, Offset: 46, Column: 1}},
					{Type: SELECT, Position: Position{Line: 1, Offset: 47, Column: 6}},
					{Type: ASTERISK, Position: Position{Line: 1, Offset: 54, Column: 1}},
					{Type: FROM, Position: Position{Line: 1, Offset: 56, Column: 4}},
					{Type: IDENT, Value: "blog", Position: Position{Line: 1, Offset: 61, Column: 4}},
					{Type: RPAREN, Position: Position{Line: 1, Offset: 65, Column: 1}},
					{Type: RPAREN, Position: Position{Line: 1, Offset: 66, Column: 1}},
					{Type: EOF, Position: Position{Line: 1, Offset: 67}},
				},
			},
			{
				`create table test (
id serial primary key,
name varchar(4) default "none",
user_id int references user(id),
blog_id int unique,
page_id int not null,
community_id int null,
file varchar(20) check(file = "foo"))`,
				[]Token{
					{Type: CREATE, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: TABLE, Position: Position{Line: 1, Offset: 7, Column: 5}},
					{Type: IDENT, Value: "test", Position: Position{Line: 1, Offset: 13, Column: 4}},
					{Type: LPAREN, Position: Position{Line: 1, Offset: 18, Column: 1}},
					{Type: IDENT, Value: "id", Position: Position{Line: 2, Offset: 20, Column: 2}},
					{Type: SERIAL, Position: Position{Line: 2, Offset: 23, Column: 6}},
					{Type: PRIMARYKEY, Position: Position{Line: 2, Offset: 30, Column: 11}},
					{Type: COMMA, Position: Position{Line: 2, Offset: 41, Column: 1}},
					{Type: IDENT, Value: "name", Position: Position{Line: 3, Offset: 43, Column: 4}},
					{Type: VARCHAR, Position: Position{Line: 3, Offset: 48, Column: 7}},
					{Type: LPAREN, Position: Position{Line: 3, Offset: 55, Column: 1}},
					{Type: INT, IntValue: 4, Position: Position{Line: 3, Offset: 56, Column: 1}},
					{Type: RPAREN, Position: Position{Line: 3, Offset: 57, Column: 1}},
					{Type: DEFAULT, Position: Position{Line: 3, Offset: 59, Column: 7}},
					{Type: IDENT, Value: "\"none\"", Position: Position{Line: 3, Offset: 67, Column: 6}},
					{Type: COMMA, Position: Position{Line: 3, Offset: 73, Column: 1}},
					{Type: IDENT, Value: "user_id", Position: Position{Line: 4, Offset: 75, Column: 7}},
					{Type: INTEGER, Position: Position{Line: 4, Offset: 83, Column: 3}},
					{Type: REFERENCE, Position: Position{Line: 4, Offset: 87, Column: 10}},
					{Type: IDENT, Value: "user", Position: Position{Line: 4, Offset: 98, Column: 4}},
					{Type: LPAREN, Position: Position{Line: 4, Offset: 102, Column: 1}},
					{Type: IDENT, Value: "id", Position: Position{Line: 4, Offset: 103, Column: 2}},
					{Type: RPAREN, Position: Position{Line: 4, Offset: 105, Column: 1}},
					{Type: COMMA, Position: Position{Line: 4, Offset: 106, Column: 1}},
					{Type: IDENT, Value: "blog_id", Position: Position{Line: 5, Offset: 108, Column: 7}},
					{Type: INTEGER, Position: Position{Line: 5, Offset: 116, Column: 3}},
					{Type: UNIQUE, Position: Position{Line: 5, Offset: 120, Column: 6}},
					{Type: COMMA, Position: Position{Line: 5, Offset: 126, Column: 1}},
					{Type: IDENT, Value: "page_id", Position: Position{Line: 6, Offset: 128, Column: 7}},
					{Type: INTEGER, Position: Position{Line: 6, Offset: 136, Column: 3}},
					{Type: NOT, Position: Position{Line: 6, Offset: 140, Column: 3}},
					{Type: NULL, Position: Position{Line: 6, Offset: 144, Column: 4}},
					{Type: COMMA, Position: Position{Line: 6, Offset: 148, Column: 1}},
					{Type: IDENT, Value: "community_id", Position: Position{Line: 7, Offset: 150, Column: 12}},
					{Type: INTEGER, Position: Position{Line: 7, Offset: 163, Column: 3}},
					{Type: NULL, Position: Position{Line: 7, Offset: 167, Column: 4}},
					{Type: COMMA, Position: Position{Line: 7, Offset: 171, Column: 1}},
					{Type: IDENT, Value: "file", Position: Position{Line: 8, Offset: 173, Column: 4}},
					{Type: VARCHAR, Position: Position{Line: 8, Offset: 178, Column: 7}},
					{Type: LPAREN, Position: Position{Line: 8, Offset: 185, Column: 1}},
					{Type: INT, IntValue: 20, Position: Position{Line: 8, Offset: 186, Column: 2}},
					{Type: RPAREN, Position: Position{Line: 8, Offset: 188, Column: 1}},
					{Type: CHECK, Position: Position{Line: 8, Offset: 190, Column: 5}},
					{Type: LPAREN, Position: Position{Line: 8, Offset: 195, Column: 1}},
					{Type: IDENT, Value: "file", Position: Position{Line: 8, Offset: 196, Column: 4}},
					{Type: EQUAL, Position: Position{Line: 8, Offset: 201, Column: 1}},
					{Type: IDENT, Value: "\"foo\"", Position: Position{Line: 8, Offset: 203, Column: 5}},
					{Type: RPAREN, Position: Position{Line: 8, Offset: 208, Column: 1}},
					{Type: RPAREN, Position: Position{Line: 8, Offset: 209, Column: 1}},
					{Type: EOF, Position: Position{Line: 8, Offset: 210}},
				},
			},
		}

		for _, c := range cases {
			assertQuery(t, c.Query, c.Token)
		}
	})

	t.Run("TestLexer_Scan_ALTER", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			Query string
			Token []Token
		}{
			{
				"alter table users add column name varchar(255)",
				[]Token{
					{Type: ALTER, Position: Position{Line: 1, Offset: 0, Column: 5}},
					{Type: TABLE, Position: Position{Line: 1, Offset: 6, Column: 5}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 12, Column: 5}},
					{Type: ADD, Position: Position{Line: 1, Offset: 18, Column: 3}},
					{Type: COLUMN, Position: Position{Line: 1, Offset: 22, Column: 6}},
					{Type: IDENT, Value: "name", Position: Position{Line: 1, Offset: 29, Column: 4}},
					{Type: VARCHAR, Position: Position{Line: 1, Offset: 34, Column: 7}},
					{Type: LPAREN, Position: Position{Line: 1, Offset: 41, Column: 1}},
					{Type: INT, IntValue: 255, Position: Position{Line: 1, Offset: 42, Column: 3}},
					{Type: RPAREN, Position: Position{Line: 1, Offset: 45, Column: 1}},
					{Type: EOF, Position: Position{Line: 1, Offset: 46}},
				},
			},
		}

		for _, c := range cases {
			assertQuery(t, c.Query, c.Token)
		}
	})

	t.Run("TestLexer_Scan_INSERT", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			Query string
			Token []Token
		}{
			{
				"insert into users values (1, \"test\")",
				[]Token{
					{Type: INSERT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: INTO, Position: Position{Line: 1, Offset: 7, Column: 4}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 12, Column: 5}},
					{Type: VALUES, Position: Position{Line: 1, Offset: 18, Column: 6}},
					{Type: LPAREN, Position: Position{Line: 1, Offset: 25, Column: 1}},
					{Type: INT, IntValue: 1, Position: Position{Line: 1, Offset: 26, Column: 1}},
					{Type: COMMA, Position: Position{Line: 1, Offset: 27, Column: 1}},
					{Type: IDENT, Value: "\"test\"", Position: Position{Line: 1, Offset: 29, Column: 6}},
					{Type: RPAREN, Position: Position{Line: 1, Offset: 35, Column: 1}},
					{Type: EOF, Position: Position{Line: 1, Offset: 36}},
				},
			},
			{
				"insert into users (id, name) values (1, \"test\")",
				[]Token{
					{Type: INSERT, Position: Position{Line: 1, Offset: 0, Column: 6}},
					{Type: INTO, Position: Position{Line: 1, Offset: 7, Column: 4}},
					{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 12, Column: 5}},
					{Type: LPAREN, Position: Position{Line: 1, Offset: 18, Column: 1}},
					{Type: IDENT, Value: "id", Position: Position{Line: 1, Offset: 19, Column: 2}},
					{Type: COMMA, Position: Position{Line: 1, Offset: 21, Column: 1}},
					{Type: IDENT, Value: "name", Position: Position{Line: 1, Offset: 23, Column: 4}},
					{Type: RPAREN, Position: Position{Line: 1, Offset: 27, Column: 1}},
					{Type: VALUES, Position: Position{Line: 1, Offset: 29, Column: 6}},
					{Type: LPAREN, Position: Position{Line: 1, Offset: 36, Column: 1}},
					{Type: INT, IntValue: 1, Position: Position{Line: 1, Offset: 37, Column: 1}},
					{Type: COMMA, Position: Position{Line: 1, Offset: 38, Column: 1}},
					{Type: IDENT, Value: "\"test\"", Position: Position{Line: 1, Offset: 40, Column: 6}},
					{Type: RPAREN, Position: Position{Line: 1, Offset: 46, Column: 1}},
					{Type: EOF, Position: Position{Line: 1, Offset: 47}},
				},
			},
		}

		for _, c := range cases {
			assertQuery(t, c.Query, c.Token)
		}
	})
}
