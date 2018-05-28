package parser

type TestQuery struct {
	Query  string
	Tokens []Token
	Ast    Query
}

var TestSelectQuery = []TestQuery{
	{ // # 0
		Query: "select * from test",
		Tokens: []Token{
			{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
			{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
			{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
			{Type: IDENT, Value: "test", Position: Position{Line: 1, Offset: 14, Column: 4}},
			{Type: EOF, Position: Position{Line: 1, Offset: 18}},
		},
		Ast: Select{
			SelectList: []SelectExpr{{Asterisk: true}},
			Table:      TableExpression{From: FromClause{Table: []TableReference{{Name: "test"}}}},
		},
	},
	{ // # 1
		Query: "select * from `test`",
		Tokens: []Token{
			{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
			{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
			{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
			{Type: IDENT, Value: "`test`", Position: Position{Line: 1, Offset: 14, Column: 6}},
			{Type: EOF, Position: Position{Line: 1, Offset: 20}},
		},
		Ast: Select{
			SelectList: []SelectExpr{{Asterisk: true}},
			Table:      TableExpression{From: FromClause{Table: []TableReference{{Name: "`test`"}}}},
		},
	},
	{ // # 2
		Query: "select    *    from    test",
		Tokens: []Token{
			{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
			{Type: ASTERISK, Position: Position{Line: 1, Offset: 10, Column: 1}},
			{Type: FROM, Position: Position{Line: 1, Offset: 15, Column: 4}},
			{Type: IDENT, Value: "test", Position: Position{Line: 1, Offset: 23, Column: 4}},
			{Type: EOF, Position: Position{Line: 1, Offset: 27}},
		},
		Ast: Select{
			SelectList: []SelectExpr{{Asterisk: true}},
			Table:      TableExpression{From: FromClause{Table: []TableReference{{Name: "test"}}}},
		},
	},
	{
		Query: "SELECT id,age from users",
		Tokens: []Token{
			{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
			{Type: IDENT, Value: "id", Position: Position{Line: 1, Offset: 7, Column: 2}},
			{Type: COMMA, Position: Position{Line: 1, Offset: 9, Column: 1}},
			{Type: IDENT, Value: "age", Position: Position{Line: 1, Offset: 10, Column: 3}},
			{Type: FROM, Position: Position{Line: 1, Offset: 14, Column: 4}},
			{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 19, Column: 5}},
			{Type: EOF, Position: Position{Line: 1, Offset: 24}},
		},
		Ast: Select{
			SelectList: []SelectExpr{{Column: "id"}, {Column: "age"}},
			Table: TableExpression{
				From: FromClause{Table: []TableReference{{Name: "users"}}},
			},
		},
	},
	{
		Query: "select * from users where id = 1",
		Tokens: []Token{
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
		Ast: Select{
			Table: TableExpression{
				From: FromClause{
					Table: []TableReference{{Name: "users"}},
				},
				Where: WhereClause{Cond: &ComparisonExpr{
					Operator:   ComparisonOperatorEqual,
					LeftValue:  ValueExpr{Type: ValueTypeString, StringValue: "id"},
					RightValue: ValueExpr{Type: ValueTypeInt, IntValue: 1},
				}},
			},
			SelectList: []SelectExpr{{Asterisk: true}},
		},
	},
	{
		Query: "select * from users where id = 1 and age = 20",
		Tokens: []Token{
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
		Ast: Select{
			Table: TableExpression{
				From: FromClause{
					Table: []TableReference{{Name: "users"}},
				},
				Where: WhereClause{
					Cond: &BooleanTerm{
						Boolean: Token{Type: AND},
						Left: &ComparisonExpr{
							Operator:   ComparisonOperatorEqual,
							LeftValue:  ValueExpr{Type: ValueTypeString, StringValue: "id"},
							RightValue: ValueExpr{Type: ValueTypeInt, IntValue: 1},
						},
						Right: &ComparisonExpr{
							Operator:   ComparisonOperatorEqual,
							LeftValue:  ValueExpr{Type: ValueTypeString, StringValue: "age"},
							RightValue: ValueExpr{Type: ValueTypeInt, IntValue: 20},
						},
					},
				},
			},
			SelectList: []SelectExpr{{Asterisk: true}},
		},
	},
	{
		Query: "select * from users where id > 10 and age > 20",
		Tokens: []Token{
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
		Ast: Select{
			Table: TableExpression{
				From: FromClause{
					Table: []TableReference{{Name: "users"}},
				},
				Where: WhereClause{
					Cond: &BooleanTerm{
						Boolean: Token{Type: AND},
						Left: &ComparisonExpr{
							Operator:   ComparisonOperatorGreaterThan,
							LeftValue:  ValueExpr{Type: ValueTypeString, StringValue: "id"},
							RightValue: ValueExpr{Type: ValueTypeInt, IntValue: 10},
						},
						Right: &ComparisonExpr{
							Operator:   ComparisonOperatorGreaterThan,
							LeftValue:  ValueExpr{Type: ValueTypeString, StringValue: "age"},
							RightValue: ValueExpr{Type: ValueTypeInt, IntValue: 20},
						},
					},
				},
			},
			SelectList: []SelectExpr{{Asterisk: true}},
		},
	},
	{
		Query: "select * from users order by created_date desc",
		Tokens: []Token{
			{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
			{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
			{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
			{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 14, Column: 5}},
			{Type: ORDERBY, Position: Position{Line: 1, Offset: 20, Column: 8}},
			{Type: IDENT, Value: "created_date", Position: Position{Line: 1, Offset: 29, Column: 12}},
			{Type: DESC, Position: Position{Line: 1, Offset: 42, Column: 4}},
			{Type: EOF, Position: Position{Line: 1, Offset: 46}},
		},
		Ast: Select{
			Table: TableExpression{
				From: FromClause{
					Table: []TableReference{{Name: "users"}},
				},
				Where: WhereClause{},
			},
			SelectList: []SelectExpr{{Asterisk: true}},
			OrderBy:    []*SortSpecification{{Key: Token{Type: IDENT, Value: "created_date"}, Order: Token{Type: DESC}}},
		},
	},
	{
		Query: "select * from users order by created_date desc, rank",
		Tokens: []Token{
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
		Ast: Select{
			Table: TableExpression{
				From: FromClause{
					Table: []TableReference{{Name: "users"}},
				},
				Where: WhereClause{},
			},
			SelectList: []SelectExpr{{Asterisk: true}},
			OrderBy:    []*SortSpecification{{Key: Token{Type: IDENT, Value: "created_date"}, Order: Token{Type: DESC}}, {Key: Token{Type: IDENT, Value: "rank"}}},
		},
	},
	{
		Query: "select * from users group by group_id",
		Tokens: []Token{
			{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
			{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
			{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
			{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 14, Column: 5}},
			{Type: GROUPBY, Position: Position{Line: 1, Offset: 20, Column: 8}},
			{Type: IDENT, Value: "group_id", Position: Position{Line: 1, Offset: 29, Column: 8}},
			{Type: EOF, Position: Position{Line: 1, Offset: 37}},
		},
		Ast: Select{
			Table: TableExpression{
				From: FromClause{
					Table: []TableReference{{Name: "users"}},
				},
				Where:   WhereClause{},
				GroupBy: GroupByClause([]Token{{Type: IDENT, Value: "group_id"}}),
			},
			SelectList: []SelectExpr{{Asterisk: true}},
		},
	},
	{
		Query: "select * from users group by group_id having group_id > 10",
		Tokens: []Token{
			{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
			{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
			{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
			{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 14, Column: 5}},
			{Type: GROUPBY, Position: Position{Line: 1, Offset: 20, Column: 8}},
			{Type: IDENT, Value: "group_id", Position: Position{Line: 1, Offset: 29, Column: 8}},
			{Type: HAVING, Position: Position{Line: 1, Offset: 38, Column: 6}},
			{Type: IDENT, Value: "group_id", Position: Position{Line: 1, Offset: 45, Column: 8}},
			{Type: GTR, Position: Position{Line: 1, Offset: 54, Column: 1}},
			{Type: INT, IntValue: 10, Position: Position{Line: 1, Offset: 56, Column: 2}},
			{Type: EOF, Position: Position{Line: 1, Offset: 58}},
		},
		Ast: Select{
			Table: TableExpression{
				From: FromClause{
					Table: []TableReference{{Name: "users"}},
				},
				Where:   WhereClause{},
				GroupBy: GroupByClause([]Token{{Type: IDENT, Value: "group_id"}}),
				Having: HavingClause{
					Cond: &ComparisonExpr{
						Operator:   ComparisonOperatorGreaterThan,
						LeftValue:  ValueExpr{Type: ValueTypeString, StringValue: "group_id"},
						RightValue: ValueExpr{Type: ValueTypeInt, IntValue: 10},
					},
				},
			},
			SelectList: []SelectExpr{{Asterisk: true}},
		},
	},
	{
		Query: "select * from users left outer join blog on users.id = blog.user_id",
		Tokens: []Token{
			{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
			{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
			{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
			{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 14, Column: 5}},
			{Type: LEFT, Position: Position{Line: 1, Offset: 20, Column: 4}},
			{Type: OUTER, Position: Position{Line: 1, Offset: 25, Column: 5}},
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
		Ast: Select{
			Table: TableExpression{
				From: FromClause{
					Table: []TableReference{{Name: "users"}},
					Join: JoinedTable{
						Type:  []Token{{Type: LEFT}, {Type: OUTER}},
						Table: []TableReference{{Name: "blog"}},
						Cond: &ComparisonExpr{
							Operator:   ComparisonOperatorEqual,
							LeftValue:  ValueExpr{Type: ValueTypeParameter, Identifiers: []string{"users", "id"}},
							RightValue: ValueExpr{Type: ValueTypeParameter, Identifiers: []string{"blog", "user_id"}},
						},
					},
				},
			},
			SelectList: []SelectExpr{{Asterisk: true}},
		},
	},
	{
		Query: "select * from users right outer join blog on users.id = blog.user_id",
		Tokens: []Token{
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
		Ast: Select{
			Table: TableExpression{
				From: FromClause{
					Table: []TableReference{{Name: "users"}},
					Join: JoinedTable{
						Type:  []Token{{Type: RIGHT}, {Type: OUTER}},
						Table: []TableReference{{Name: "blog"}},
						Cond: &ComparisonExpr{
							Operator:   ComparisonOperatorEqual,
							LeftValue:  ValueExpr{Type: ValueTypeParameter, Identifiers: []string{"users", "id"}},
							RightValue: ValueExpr{Type: ValueTypeParameter, Identifiers: []string{"blog", "user_id"}},
						},
					},
				},
			},
			SelectList: []SelectExpr{{Asterisk: true}},
		},
	},
	{ // # 13
		Query: "SELECT id as foo from users",
		Tokens: []Token{
			{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
			{Type: IDENT, Value: "id", Position: Position{Line: 1, Offset: 7, Column: 2}},
			{Type: AS, Position: Position{Line: 1, Offset: 10, Column: 2}},
			{Type: IDENT, Value: "foo", Position: Position{Line: 1, Offset: 13, Column: 3}},
			{Type: FROM, Position: Position{Line: 1, Offset: 17, Column: 4}},
			{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 22, Column: 5}},
			{Type: EOF, Position: Position{Line: 1, Offset: 27}},
		},
		Ast: Select{
			Table: TableExpression{
				From: FromClause{
					Table: []TableReference{{Name: "users"}},
				},
			},
			SelectList: []SelectExpr{{Column: "id", Alias: "foo"}},
		},
	},
	{ // # 14
		Query: "SELECT * FROM users WHERE name = ?",
		Tokens: []Token{
			{Type: SELECT, Position: Position{Line: 1, Offset: 0, Column: 6}},
			{Type: ASTERISK, Position: Position{Line: 1, Offset: 7, Column: 1}},
			{Type: FROM, Position: Position{Line: 1, Offset: 9, Column: 4}},
			{Type: IDENT, Value: "users", Position: Position{Line: 1, Offset: 14, Column: 5}},
			{Type: WHERE, Position: Position{Line: 1, Offset: 20, Column: 5}},
			{Type: IDENT, Value: "name", Position: Position{Line: 1, Offset: 26, Column: 4}},
			{Type: EQUAL, Position: Position{Line: 1, Offset: 31, Column: 1}},
			{Type: QUESTION, Position: Position{Line: 1, Offset: 33, Column: 1}},
			{Type: EOF, Position: Position{Line: 1, Offset: 35}},
		},
		Ast: Select{
			SelectList: []SelectExpr{{Asterisk: true}},
			Table: TableExpression{
				From: FromClause{
					Table: []TableReference{{Name: "users"}},
				},
				Where: WhereClause{
					Cond: &ComparisonExpr{
						Operator:   ComparisonOperatorEqual,
						LeftValue:  ValueExpr{Type: ValueTypeString, StringValue: "name"},
						RightValue: ValueExpr{Type: ValueTypeDynamicParameter},
					},
				},
			},
		},
	},
}
