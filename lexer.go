package parser

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type Lexer struct {
	ctx *lexerCtx
}

type lexerCtx struct {
	r       *bufio.Reader
	line    int
	start   int
	cur     int
	prevCur int
}

func (ctx *lexerCtx) Position() Position {
	return Position{Line: ctx.line, Offset: ctx.prevCur, Column: ctx.cur - ctx.start}
}

func (ctx *lexerCtx) Token(typ TokenType, value interface{}) Token {
	var intValue int
	var stringValue string

	switch x := value.(type) {
	case int:
		intValue = x
	case string:
		stringValue = x
	}
	return Token{Type: typ, Position: ctx.Position(), Value: stringValue, IntValue: intValue}
}

func (ctx *lexerCtx) peek() (rune, error) {
	b, err := ctx.r.Peek(1)
	if err != nil {
		return 0, err
	}
	if isLineBreak(rune(b[0])) {
		ctx.line++
	}
	return rune(b[0]), err
}

func (ctx *lexerCtx) discard() {
	ctx.cur++
	ctx.r.Discard(1)
}

func (ctx *lexerCtx) skipWhiteSpace() {
	for {
		r, err := ctx.peek()
		if err == io.EOF {
			break
		}
		if isWhiteSpace(r) {
			ctx.discard()
		} else {
			break
		}
	}
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{ctx: &lexerCtx{r: bufio.NewReader(r), line: 1, cur: 0}}
}

func (lexer *Lexer) Scan() (Token, error) {
	var r rune
	for {
		b, err := lexer.ctx.peek()
		if err == io.EOF {
			return Token{Type: EOF}, err
		}
		if err != nil {
			return Token{}, err
		}
		r = b

		if isWhiteSpace(r) {
			lexer.ctx.discard()
			continue
		} else if isLineBreak(r) {
			lexer.ctx.discard()
			continue
		} else {
			lexer.ctx.prevCur = lexer.ctx.cur
			lexer.ctx.start = lexer.ctx.cur
			break
		}
	}

	r = unicode.ToLower(r)
	var typ TokenType
	var value interface{}
	switch r {
	case ',':
		lexer.ctx.discard()
		typ = COMMA
	case '.':
		lexer.ctx.discard()
		typ = PERIOD
	case '*':
		lexer.ctx.discard()
		typ = ASTERISK
	case '=':
		lexer.ctx.discard()
		typ = EQUAL
	case '<':
		lexer.ctx.discard()
		typ = LSS
	case '>':
		lexer.ctx.discard()
		typ = GTR
	case '(':
		lexer.ctx.discard()
		typ = LPAREN
	case ')':
		lexer.ctx.discard()
		typ = RPAREN
	default:
		typ, value = lexer.scanStatement(lexer.ctx, r)
	}

	token := lexer.ctx.Token(typ, value)
	return token, nil
}

func (lexer *Lexer) scanStatement(ctx *lexerCtx, s rune) (TokenType, interface{}) {
	statement := make([]rune, 0, 0)
	statement = append(statement, s)
	ctx.discard()
	for {
		r, err := ctx.peek()
		if err == io.EOF {
			break
		}
		if isWhiteSpace(r) {
			break
		} else if isComma(r) || isPeriod(r) || isParen(r) {
			break
		} else {
			statement = append(statement, r)
			ctx.discard()
			continue
		}
	}

	state := strings.ToLower(string(statement))
	switch state {
	case "select":
		return SELECT, nil
	case "insert":
		return INSERT, nil
	case "update":
		return UPDATE, nil
	case "create":
		return CREATE, nil
	case "alter":
		return ALTER, nil
	case "from":
		return FROM, nil
	case "as":
		return AS, nil
	case "where":
		return WHERE, nil
	case "set":
		return SET, nil
	case "into":
		return INTO, nil
	case "order", "group":
		ctx.skipWhiteSpace()
		b, err := ctx.r.Peek(2)
		if err == io.EOF {
			ctx.cur += 2
			ctx.r.Discard(2)
			return ILLEGAL, nil
		}
		if err != nil {
			ctx.cur += 2
			ctx.r.Discard(2)
			return ILLEGAL, nil
		}
		ctx.cur += 2
		ctx.r.Discard(2)
		if string(b) == "by" {
			if state == "order" {
				return ORDERBY, nil
			} else {
				return GROUPBY, nil
			}
		}
	case "primary":
		ctx.skipWhiteSpace()
		b, err := ctx.r.Peek(3)
		if err == io.EOF {
			ctx.cur += 3
			ctx.r.Discard(3)
			return ILLEGAL, nil
		}
		if string(b) == "key" {
			ctx.cur += 3
			ctx.r.Discard(3)
			return PRIMARYKEY, nil
		}
	case "desc":
		return DESC, nil
	case "asc":
		return ASC, nil
	case "having":
		return HAVING, nil
	case "left":
		return LEFT, nil
	case "right":
		return RIGHT, nil
	case "inner":
		return INNER, nil
	case "outer":
		return OUTER, nil
	case "join":
		return JOIN, nil
	case "on":
		return ON, nil
	case "values":
		return VALUES, nil
	case "database":
		return DATABASE, nil
	case "table":
		return TABLE, nil
	case "assertion":
		return ASSERTION, nil
	case "check":
		return CHECK, nil
	case "and":
		return AND, nil
	case "or":
		return OR, nil
	case "not":
		return NOT, nil
	case "exist":
		return EXIST, nil
	case "add":
		return ADD, nil
	case "column":
		return COLUMN, nil
	case "serial":
		return SERIAL, nil
	case "references":
		return REFERENCE, nil
	case "default":
		return DEFAULT, nil
	case "varchar":
		return VARCHAR, nil
	case "int":
		return INTEGER, nil
	case "unique":
		return UNIQUE, nil
	case "null":
		return NULL, nil
	default:
		if i, err := strconv.Atoi(state); err == nil {
			return INT, i
		}
		return IDENT, state
	}

	return ILLEGAL, nil
}

func isWhiteSpace(r rune) bool {
	if r == ' ' {
		return true
	}
	return false
}

func isLineBreak(r rune) bool {
	if r == '\n' {
		return true
	}
	return false
}

func isComma(r rune) bool {
	if r == ',' {
		return true
	}
	return false
}

func isPeriod(r rune) bool {
	if r == '.' {
		return true
	}
	return false
}

func isParen(r rune) bool {
	if r == '(' || r == ')' {
		return true
	}
	return false
}
