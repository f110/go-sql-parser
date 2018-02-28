package parser

//go:generate stringer -type=TokenType

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	WS

	INT
	IDENT

	ASTERISK // *
	COMMA    // ,
	PERIOD   // .
	LPAREN   // (
	RPAREN   // )
	ADD      // +
	SUB      // -
	EQUAL    // =
	LSS      // <
	GTR      // >

	SELECT
	INSERT
	UPDATE
	DELETE
	CREATE
	ALTER
	DROP
	FROM
	AS
	SET
	INTO
	WHERE
	JOIN
	LEFT
	RIGHT
	OUTER
	INNER
	ON
	GROUPBY
	ORDERBY
	HAVING
	VALUES
	DESC
	ASC
	NULL
	PRIMARYKEY
	AND
	OR
	IF
	NOT
	EXIST
	COLUMN
	DEFAULT

	DATABASE
	TABLE
	ASSERTION
	INDEX
	CHECK
	REFERENCE
	UNIQUE

	INTEGER
	SERIAL
	VARCHAR
)

type Token struct {
	Type     TokenType
	Position Position
	Value    string
	IntValue int
}

type Position struct {
	Line   int
	Offset int
	Column int
}
