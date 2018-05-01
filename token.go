package parser

import (
	"io"
)

type TokensReader struct {
	tokens []Token
}

func NewTokensReader(tokens []Token) TokenReader {
	return &TokensReader{tokens: tokens}
}

func (tr *TokensReader) Scan() (Token, error) {
	if len(tr.tokens) > 0 {
		t := tr.tokens[0]
		tr.tokens = tr.tokens[1:]
		return t, nil
	}

	return Token{Type: EOF}, io.EOF
}

func (tr *TokensReader) Peek(n int) ([]Token, error) {
	if len(tr.tokens) >= n {
		return tr.tokens[:n], nil
	}

	return nil, io.EOF
}

func (tr *TokensReader) Discard(n int) {
	if len(tr.tokens) >= n {
		tr.tokens = tr.tokens[n:]
	}
}
