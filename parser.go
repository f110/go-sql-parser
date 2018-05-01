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
// <between predicate> ::=
//		<row value expression> [ NOT ] BETWEEN [ ASYMMETRIC | SYMMETRIC ]
//		<row value expression> AND <row value expression>

type Tokens []Token

type Query interface{}

type Select struct {
	SelectList SelectList
	Table      TableExpression
}

type SelectList []Token

type TableExpression struct {
	From    FromClause
	Where   WhereClause
	GroupBy GroupByClause
	Having  HavingClause
}

type FromClause struct {
	Table TableList
}

type WhereClause struct {
	Cond Expr
}

type GroupByClause struct{}

type HavingClause struct{}

type TableList []Token

type Expr interface{}

type SearchCondition struct{}

type BooleanTerm struct {
	Boolean Token
	Left    Expr
	Right   Expr
}

const (
	BooleanOperatorAnd BooleanOperator = iota
	BooleanOperatorOr
)

type BooleanOperator int

type ValueExpr struct {
	Operator   Token
	LeftValue  Expr
	RightValue Expr
}

type RawValue struct {
	Token Token
}

const (
	OperatorEqual ComparisonOperator = iota
	OperatorNotEqual
	OperatorLess
	OperatorGreater
	OperatorLessOrEqual
	OperatorGreatherOrEqual
)

type ComparisonOperator int

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

	return query, nil
}

func (p *Parser) parseSelectList(tokens TokenReader) (SelectList, error) {
	res := make(SelectList, 0)

	for {
		t, err := tokens.Peek(1)
		if err != nil {
			return nil, err
		}
		if err == io.EOF {
			break
		}
		if t[0].Type == FROM {
			return res, nil
		}

		tokens.Discard(1)
		switch t[0].Type {
		case ASTERISK:
			res = append(res, t[0])
		case IDENT:
			res = append(res, t[0])
		}
	}

	return res, nil
}

func (p *Parser) parseTableExpression(tokens TokenReader) (TableExpression, error) {
	fromClause, err := p.parseFromClause(tokens)
	if err != nil {
		return TableExpression{}, err
	}
	tableExpression := &TableExpression{From: fromClause}

	if _, err := tokens.Peek(1); err == nil {
		var whereClause WhereClause

		whereClause, err = p.parseWhereClause(tokens)
		if err != nil {
			return TableExpression{}, err
		}
		tableExpression.Where = whereClause
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
	for {
		t, err := tokens.Peek(1)
		if err != nil && err != io.EOF {
			return FromClause{}, err
		}
		if err == io.EOF {
			break
		}
		if t[0].Type == WHERE {
			break
		}

		tokens.Discard(1)
		switch t[0].Type {
		case IDENT:
			tableList = append(tableList, t[0])
		}
	}

	return FromClause{Table: tableList}, nil
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
			for {
				t, err := tokens.Scan()
				if err == io.EOF {
					break
				}
				least = append(least, t)
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
			tokens.Discard(1)
			r, err := p.parseBooleanValueExpression(tokens)
			_ = err
			v := &ValueExpr{
				Operator:   t[0],
				LeftValue:  l,
				RightValue: r,
			}
			return v, nil
		}

		left = append(left, t[0])
		tokens.Discard(1)
	}

	return left, nil
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
