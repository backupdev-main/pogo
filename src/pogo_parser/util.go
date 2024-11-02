package pogo_parser

import (
	"fmt"
	"pogo/src/semantic"
	"pogo/src/token"
)

func (p *Parser) next() {
	p.curr = p.lexer.Scan()
}

func (p *Parser) expect(typ token.Type) error {
	//fmt.Println("Processing token", string(p.curr.Lit), "as ", token.TokMap.Id(typ))
	// fmt.Println("This is the context", p.lexer.s)
	if p.curr.Type != typ {
		return p.error(fmt.Sprintf("expected %v, got %v", token.TokMap.Id(typ), token.TokMap.Id(p.curr.Type)))
	}
	p.next()
	return nil
}

// Error helper function
func (p *Parser) error(msg string) error {
	return fmt.Errorf("line %d: %s", p.curr.Line, msg)
}

func (p *Parser) isStatementStart() (bool, error) {
	statementStarts := map[token.Type]struct{}{
		token.TokMap.Type("kwdWhile"): {},
		token.TokMap.Type("kwdIf"):    {},
		token.TokMap.Type("kwdPrint"): {},
		token.TokMap.Type("id"):       {},
	}

	if _, exists := statementStarts[p.curr.Type]; exists {
		return true, nil
	}
	return false, nil
}

func (p *Parser) addVariablesToSymbolTable(semType semantic.Type, currentVars []string) error {

	for _, varName := range currentVars {
		if err := p.symbolTable.AddVariable(varName, semType, p.curr.Line, p.curr.Column); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) returnSemanticType(currType string) (semantic.Type, error) {
	var semType semantic.Type
	switch string(currType) {
	case "int":
		semType = semantic.TypeInt
	case "float":
		semType = semantic.TypeFloat
	default:
		return semantic.TypeError, fmt.Errorf("line %d: unsupported type: %s", p.curr.Line, string(currType))
	}
	return semType, nil
}

func (p *Parser) getType(tok *token.Token) (semantic.Type, error) {
	switch tok.Type {
	case token.TokMap.Type("intLit"):
		p.next()
		return semantic.TypeInt, nil
	case token.TokMap.Type("floatLit"):
		p.next()
		return semantic.TypeFloat, nil
	case token.TokMap.Type("id"):
		p.next()
		return p.symbolTable.GetType(string(tok.Lit))
	default:
		return semantic.TypeError, fmt.Errorf("expected number after %s", p.curr.Lit)
	}
}

func (p *Parser) isAssignable(varType, exprType semantic.Type) bool {
	if varType == exprType {
		return true
	}
	// Allow int -> float conversion
	if varType == semantic.TypeFloat && exprType == semantic.TypeInt {
		return true
	}
	return false
}
