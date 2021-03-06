package parser

import (
	"errors"
	"io"
)

var (
	ErrInvalidQuery = errors.New("parser: invalid query")
)

// Select
// Query: SELECT * FROM test
//		  SELECT id, name FROM test
// <query specification> ::= SELECT [ <set quantifier> ] <select list> <table expression>
// <set quantifier> ::= DISTINCT | ALL
// <select list> ::= <asterisk> | <select sublist> [ { <comma> <select sublist> }... ]
// <select sublist> ::= <derived column> | <qualified asterisk>
// <derived column> ::= <value expression> [ <as clause> ]
// <as clause> ::= [ AS ] <column name>
// <qualified asterisk> ::= <asterisked identifier chain> <period> <asterisk> |	<all fields reference>
// <table expression> ::= <from clause> [ <where clause> ] [ <group by clause> ] [ <having clause> ]
// <from clause> ::= FROM <table reference list>
// <table reference list> ::= <table reference> [ { <comma> <table reference> }... ]
// <where clause> ::= WHERE <search condition>
// <search condition> ::= <boolean value expression>
// <boolean value expression> ::= <boolean term> | <boolean value expression> OR <boolean term>
// <boolean term> ::= <boolean factor> | <boolean term> AND <boolean factor>
// <boolean factor> ::= [ NOT ] <boolean test>
// <boolean test> ::= <boolean primary> [ IS [ NOT ] <truth value> ]
// <boolean primary> ::= <predicate> | <parenthesized boolean value expression> | <nonparenthesized value expression primary>
// <predicate> ::=
//		<comparison predicate>
//	|	<between predicate>
//	|	<in predicate>
//	|	<like predicate>
//	|	<null predicate>
//	|	<quantified comparison predicate>
//	|	<exists predicate>
//	|	<unique predicate>
//	|	<match predicate>
//	|	<overlaps predicate>
//	|	<similar predicate>
//	|	<distinct predicate>
//	|	<type predicate>
// <comparison predicate> ::= <row value expression> <comp op> <row value expression>
// <row value expression> ::= <row value special case> | <row value constructor>
// <row value constructor> ::=
//		<row value constructor element>
//	|	[ ROW ] <left paren> <row value constructor element list> <right paren>
//	|	<row subquery>
// <row value constructor element> ::= <value expression>
// <between predicate> ::=
//		<row value expression> [ NOT ] BETWEEN [ ASYMMETRIC | SYMMETRIC ]
//		<row value expression> AND <row value expression>

type Tokens []Token

type Query interface{}

type Select struct {
	SelectList SelectList
	Table      TableExpression
	OrderBy    OrderByClause
}

type OrderByClause []*SortSpecification

type SortSpecification struct {
	Key   Token
	Order Token
}

type SelectList []SelectExpr

type SelectExpr struct {
	Column   string
	Alias    string
	Asterisk bool
}

type TableExpression struct {
	From    FromClause
	Where   WhereClause
	GroupBy GroupByClause
	Having  HavingClause
}

type FromClause struct {
	Table TableList
	Join  JoinedTable
}

type JoinedTable struct {
	Table TableList
	Type  Tokens
	Cond  Expr
}

type WhereClause struct {
	Cond Expr
}

type GroupByClause ColumnList

type HavingClause struct {
	Cond Expr
}

type ColumnList Tokens

type TableList []TableReference

type TableReference struct {
	Name  string
	Alias string
}

type Expr interface{}

type SearchCondition struct{}

type BooleanTerm struct {
	Boolean Token
	Left    Expr
	Right   Expr
}

const (
	ComparisonOperatorEqual = iota
	ComparisonOperatorLessThan
	ComparisonOperatorGreaterThan
)

type ComparisonOperator int

type ComparisonExpr struct {
	Operator   ComparisonOperator
	LeftValue  ValueExpr
	RightValue ValueExpr
}

const (
	ValueTypeInt = iota
	ValueTypeString
	ValueTypeParameter
	ValueTypeDynamicParameter // ?
)

type ValueType int

type ValueExpr struct {
	Type        ValueType
	IntValue    int
	StringValue string
	Identifiers []string
}

type RawValue struct {
	Token Token
}

type Parser struct{}

type TokenReader interface {
	Scan() (Token, error)
	Peek(n int) ([]Token, error)
	Discard(n int)
}

func (p *Parser) Parse(tokens TokenReader) (Query, error) {
	t, err := tokens.Peek(1)
	if err != nil {
		return nil, ErrInvalidQuery
	}

	switch t[0].Type {
	case SELECT:
		return p.parseSelect(tokens)
	}

	return nil, ErrInvalidQuery
}

func (p *Parser) parseSelect(tokens TokenReader) (Query, error) {
	if t, err := tokens.Peek(1); err != nil || t[0].Type != SELECT {
		return nil, ErrInvalidQuery
	} else {
		tokens.Discard(1)
	}
	query := NewSelect()

	selectList, err := p.parseSelectList(tokens)
	if err != nil && err != io.EOF {
		return nil, err
	}
	query.SelectList = selectList

	tableExpression, err := p.parseTableExpression(tokens)
	if err != nil && err != io.EOF {
		return nil, err
	}
	query.Table = tableExpression
	if err == io.EOF {
		return query, nil
	}

	orderByClause, err := p.parseOrderByClause(tokens)
	if err != nil && err != io.EOF {
		return nil, err
	}
	query.OrderBy = orderByClause

	return query, nil
}

func (p *Parser) parseOrderByClause(tokens TokenReader) (OrderByClause, error) {
	if t, err := tokens.Peek(1); err != nil || t[0].Type != ORDERBY {
		return OrderByClause{}, nil
	} else {
		tokens.Discard(1)
	}

	res := make([]*SortSpecification, 0)
	left := make(Tokens, 0, 2)
	for {
		t, err := tokens.Peek(1)
		if err != nil && err != io.EOF {
			return OrderByClause{}, nil
		}
		if err == io.EOF {
			res = append(res, p.parseSortSpecification(left))
			break
		}

		tokens.Discard(1)
		switch t[0].Type {
		case COMMA:
			res = append(res, p.parseSortSpecification(left))
			left = make(Tokens, 0, 2)
			continue
		default:
			left = append(left, t[0])
		}
	}

	return OrderByClause(res), nil
}

func (p *Parser) parseSortSpecification(t []Token) *SortSpecification {
	s := &SortSpecification{Key: t[0]}
	if len(t) > 1 {
		s.Order = t[1]
	}

	return s
}

func (p *Parser) parseSelectList(tokens TokenReader) (SelectList, error) {
	res := make(SelectList, 0)

	left := make(Tokens, 0)
	for {
		t, err := tokens.Peek(1)
		if err != nil {
			return nil, err
		}
		if err == io.EOF {
			break
		}
		if t[0].Type == FROM {
			e, err := p.parseSelectExpr(left)
			if err != nil {
				return nil, err
			}
			res = append(res, e)
			return res, nil
		}

		tokens.Discard(1)
		switch t[0].Type {
		case COMMA:
			e, err := p.parseSelectExpr(left)
			if err != nil {
				return nil, err
			}

			res = append(res, e)
			left = make(Tokens, 0)
		default:
			left = append(left, t[0])
		}
	}

	return res, nil
}

func (p *Parser) parseSelectExpr(tokens Tokens) (SelectExpr, error) {
	res := SelectExpr{}
	if len(tokens) > 1 && tokens[1].Type == AS && tokens[2].Type == IDENT {
		res.Alias = tokens[2].Value
	}

	switch tokens[0].Type {
	case IDENT:
		res.Column = tokens[0].Value
	case ASTERISK:
		res.Asterisk = true
	}

	return res, nil
}

func (p *Parser) parseTableExpression(tokens TokenReader) (TableExpression, error) {
	fromClause, err := p.parseFromClause(tokens)
	if err != nil {
		return TableExpression{}, err
	}
	tableExpression := &TableExpression{From: fromClause}

	if t, err := tokens.Peek(1); err == nil && t[0].Type == WHERE {
		whereClause, err := p.parseWhereClause(tokens)
		if err != nil {
			return TableExpression{}, err
		}
		tableExpression.Where = whereClause
	}

	if t, err := tokens.Peek(1); err == nil && t[0].Type == GROUPBY {
		groupByClause, err := p.parseGroupByClause(tokens)
		if err != nil {
			return TableExpression{}, err
		}
		tableExpression.GroupBy = groupByClause
	}

	if t, err := tokens.Peek(1); err == nil && t[0].Type == HAVING {
		havingClause, err := p.parseHavingClause(tokens)
		if err != nil {
			return TableExpression{}, err
		}
		tableExpression.Having = havingClause
	}

	return *tableExpression, nil
}

func (p *Parser) parseFromClause(tokens TokenReader) (FromClause, error) {
	if t, err := tokens.Peek(1); err != nil || t[0].Type != FROM {
		return FromClause{}, ErrInvalidQuery
	} else {
		tokens.Discard(1)
	}

	tableList := make(TableList, 0)
TableList:
	for {
		t, err := tokens.Peek(1)
		if err != nil && err != io.EOF {
			return FromClause{}, err
		}
		if err == io.EOF {
			break
		}

		switch t[0].Type {
		case WHERE, ORDERBY, GROUPBY, LEFT, RIGHT:
			break TableList
		}

		tokens.Discard(1)
		switch t[0].Type {
		case IDENT:
			tableList = append(tableList, TableReference{Name: t[0].Value})
		}
	}

	joinedTableList, err := p.parseJoinedTable(tokens)
	if err != nil && err != io.EOF {
		return FromClause{}, ErrInvalidQuery
	}

	return FromClause{Table: tableList, Join: joinedTableList}, nil
}

func (p *Parser) parseJoinedTable(tokens TokenReader) (JoinedTable, error) {
	if t, err := tokens.Peek(1); err != nil {
		return JoinedTable{}, err
	} else {
		switch t[0].Type {
		case LEFT, RIGHT:
		default:
			return JoinedTable{}, err
		}
	}

	joined := JoinedTable{}
	left := make(Tokens, 0)
JoinedTable:
	for {
		t, err := tokens.Peek(1)
		if err != nil {
			return JoinedTable{}, err
		}

		switch t[0].Type {
		case LEFT, RIGHT, OUTER:
			left = append(left, t[0])
			tokens.Discard(1)
		case JOIN:
			joined.Type = left
			tokens.Discard(1)
			left = make(Tokens, 0)
		case ON:
			if len(left) > 0 {
				joined.Table = append(joined.Table, TableReference{Name: left[0].Value})
			} else {
				return JoinedTable{}, ErrInvalidQuery
			}
			tokens.Discard(1)
			e, err := p.parseSearchCondition(tokens)
			if err != nil {
				return JoinedTable{}, err
			}
			joined.Cond = e
			break JoinedTable
		case IDENT:
			left = append(left, t[0])
			tokens.Discard(1)
		}
	}

	return joined, nil
}

func (p *Parser) parseGroupByClause(tokens TokenReader) (GroupByClause, error) {
	if t, err := tokens.Peek(1); err != nil || t[0].Type != GROUPBY {
		return GroupByClause([]Token{}), nil
	} else {
		tokens.Discard(1)
	}

	res := make(GroupByClause, 0)
	for {
		t, err := tokens.Peek(1)
		if err != nil && err != io.EOF {
			return GroupByClause([]Token{}), nil
		}
		if err == io.EOF {
			break
		}
		if t[0].Type == HAVING {
			break
		}

		tokens.Discard(1)
		switch t[0].Type {
		case COMMA, EOF:
			continue
		default:
			res = append(res, t[0])
		}
	}

	return res, nil
}

func (p *Parser) parseHavingClause(tokens TokenReader) (HavingClause, error) {
	if t, err := tokens.Peek(1); err != nil || t[0].Type != HAVING {
		return HavingClause{}, nil
	} else {
		tokens.Discard(1)
	}

	expr, err := p.parseSearchCondition(tokens)
	if err != nil {
		return HavingClause{}, ErrInvalidQuery
	}

	return HavingClause{Cond: expr}, nil
}

func (p *Parser) parseWhereClause(tokens TokenReader) (WhereClause, error) {
	if t, err := tokens.Peek(1); err != nil || t[0].Type != WHERE {
		return WhereClause{}, ErrInvalidQuery
	} else {
		tokens.Discard(1)
	}

	expr, err := p.parseSearchCondition(tokens)
	if err != nil {
		return WhereClause{}, ErrInvalidQuery
	}

	return WhereClause{Cond: expr}, nil
}

func (p *Parser) parseSearchCondition(tokens TokenReader) (Expr, error) {
	depth := 0

	left := make(Tokens, 0)
	for {
		t, err := tokens.Peek(1)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}

		switch t[0].Type {
		case LPAREN:
			depth++
			if depth < 0 {
				tokens.Discard(1)
				continue
			}
		case RPAREN:
			depth--
			if depth < 0 {
				tokens.Discard(1)
				continue
			}
		case AND, OR:
			if depth != 0 {
				break
			}

			l, err := p.parseSearchCondition(NewTokensReader(removeRedundantParen(left)))
			tokens.Discard(1)

			least := make(Tokens, 0)
		Skip:
			for {
				t, err := tokens.Peek(1)
				if err == io.EOF {
					break
				}

				switch t[0].Type {
				case ORDERBY, GROUPBY:
					break Skip
				}
				tokens.Discard(1)
				least = append(least, t[0])
			}

			r, err := p.parseSearchCondition(NewTokensReader(removeRedundantParen(least)))
			_ = err
			b := &BooleanTerm{
				Boolean: t[0],
				Left:    l,
				Right:   r,
			}
			return b, nil
		}

		left = append(left, t[0])
		tokens.Discard(1)
	}

	return p.parseBooleanValueExpression(NewTokensReader(left))
}

func (p *Parser) parseBooleanValueExpression(tokens TokenReader) (Expr, error) {
	left := make(Tokens, 0)
	for {
		t, err := tokens.Peek(1)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}

		switch t[0].Type {
		case EQUAL, LSS, GTR:
			l, err := p.parseBooleanValueExpression(NewTokensReader(left))
			if err != nil {
				return &ComparisonExpr{}, err
			}
			tokens.Discard(1)
			r, err := p.parseBooleanValueExpression(tokens)
			if err != nil {
				return &ComparisonExpr{}, err
			}
			o := ComparisonOperator(ComparisonOperatorEqual)
			switch t[0].Type {
			case LSS:
				o = ComparisonOperatorLessThan
			case GTR:
				o = ComparisonOperatorGreaterThan
			}

			v := &ComparisonExpr{
				Operator:   o,
				LeftValue:  p.parseValueExpr(l.(Tokens)),
				RightValue: p.parseValueExpr(r.(Tokens)),
			}
			return v, nil
		case EOF:
			tokens.Discard(1)
			return left, nil
		}

		left = append(left, t[0])
		tokens.Discard(1)
	}

	return left, nil
}

func (p *Parser) parseValueExpr(tokens Tokens) ValueExpr {
	if len(tokens) == 1 {
		switch tokens[0].Type {
		case IDENT:
			return ValueExpr{Type: ValueTypeString, StringValue: tokens[0].Value}
		case INT:
			return ValueExpr{Type: ValueTypeInt, IntValue: tokens[0].IntValue}
		case QUESTION:
			return ValueExpr{Type: ValueTypeDynamicParameter}
		}
	} else {
		if tokens[1].Type == PERIOD {
			identifies := make([]string, 0, len(tokens)+1)
			for i := 0; i < len(tokens); i += 2 {
				identifies = append(identifies, tokens[i].Value)
			}
			return ValueExpr{Type: ValueTypeParameter, Identifiers: identifies}
		}
	}
	return ValueExpr{}
}

func (p *Parser) parserPredicate(token Token) Expr {
	return &RawValue{Token: token}
}

func (p *Parser) parseBooleanTerm(tokens TokenReader) (Expr, error) {
	return BooleanTerm{}, nil
}

func removeRedundantParen(tokens Tokens) Tokens {
	if len(tokens) < 3 {
		return tokens
	}

	if tokens[0].Type == LPAREN && tokens[len(tokens)-1].Type == RPAREN {
		return tokens[1 : len(tokens)-1]
	}

	return tokens
}

func NewSelect() *Select {
	return &Select{}
}
