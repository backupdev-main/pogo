package ast

import (
	"pogo/src/token"
)

type Attrib interface{}

type Program struct {
	VarDecl     []Statement
	Functions   []Statement
	MainSection []Statement
}

type Node interface {
	TokenLiteral() string
}

type Expression interface {
	Node
	expressionNode()
}

type Statement interface {
	Node
	statementNode()
}

type AssigmentStatement struct {
	Token *token.Token
	Left  Identifier
	Right Expression
}

type FunctionStatement struct {
	Token      *token.Token
	Name       string
	Parameters []Parameters
	Body       *BlockStatement
}

type Parameters struct {
	Param string
	Type  string
}

type IfStatement struct {
	Token       *token.Token
	Condition   Expression
	Block       *BlockStatement
	Alternative *BlockStatement
}

type WhileStatement struct {
	Token     *token.Token
	Condition Expression
	Block     *BlockStatement
}

type BlockStatement struct {
	Token      *token.Token
	Statements []Statement
}

type PrintStatement struct {
	Token       *token.Token
	Expressions []Expression
}

type Identifier struct {
	Token *token.Token
	Value string
}

type IntegerLiteral struct {
	Token *token.Token
	Value string
}

type StringLiteral struct {
	Token *token.Token
	Value string
}

type FloatLiteral struct {
	Token *token.Token
	Value string
}

type FunctionCall struct {
	Token  *token.Token
	Name   string
	Params []Expression
	Type   string
}

type Declare struct {
	Token   *token.Token
	VarList []Identifier
	Type    string
}
