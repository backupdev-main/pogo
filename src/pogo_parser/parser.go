package pogo_parser

import (
	"fmt"
	"pogo/src/lexer"
	"pogo/src/semantic"
	"pogo/src/token"
)

const (
// INT_TYPE   = "int"
// FLOAT_TYPE = "float"
)

type Parser struct {
	lexer        *lexer.Lexer
	curr         *token.Token
	symbolTable  *semantic.SymbolTable
	semanticCube *semantic.SemanticCube
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:        l,
		symbolTable:  semantic.NewSymbolTable(),
		semanticCube: semantic.NewSemanticCube(),
	}
	p.next() // prime first token
	return p
}

func (p *Parser) next() {
	p.curr = p.lexer.Scan()
}

// This is a basic error handling method for our parser
func (p *Parser) error(msg string) error {
	return fmt.Errorf("line %d: %s", p.curr.Line, msg)
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

func (p *Parser) ParseProgram() error {
	if err := p.parseProgramName(); err != nil {
		return err
	}

	if err := p.parseVarDeclarationSection(); err != nil {
		return err
	}

	if err := p.parseFunctionListOpt(); err != nil {
		return err
	}

	if err := p.parseMainSection(); err != nil {
		return err
	}

	return nil
}

func (p *Parser) parseProgramName() error {
	if err := p.expect(token.TokMap.Type("kwdProgram")); err != nil {
		return err
	}
	if err := p.expect(token.TokMap.Type("id")); err != nil {
		return err
	}
	return p.expect(token.TokMap.Type("terminator"))
}

func (p *Parser) parseVarDeclarationSection() error {
	if p.curr.Type != token.TokMap.Type("kwdVars") {
		if p.curr.Type != token.TokMap.Type("kwdFunc") &&
			p.curr.Type != token.TokMap.Type("kwdBegin") {
			return fmt.Errorf("line %d: unexpected token '%s', expected 'var', 'func', or 'begin'",
				p.curr.Line, p.curr.Lit)
		}
		return nil
	}

	for p.curr.Type == token.TokMap.Type("kwdVars") {
		if err := p.parseVarDeclaration(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) parseVarDeclaration() error {
	if err := p.expect(token.TokMap.Type("kwdVars")); err != nil {
		return err
	}

	if err := p.parseVarList(); err != nil {
		return err
	}

	// process varType in a later process

	if err := p.expect(token.TokMap.Type("terminator")); err != nil {
		return err
	}

	return p.parseVarDeclarationSection()
}

func (p *Parser) parseVarList() error {
	currentVars := make([]string, 0)
	initialVar := string(p.curr.Lit)
	currentVars = append(currentVars, initialVar)
	if err := p.expect(token.TokMap.Type("id")); err != nil {
		return err
	}
	// Append id into currentVars

	for p.curr.Type == token.TokMap.Type("repeatTerminator") {
		p.next()
		currentVars = append(currentVars, string(p.curr.Lit))
		if err := p.expect(token.TokMap.Type("id")); err != nil {
			return err
		}
	}

	if err := p.expect(token.TokMap.Type("typeAssignOp")); err != nil {
		return err
	}

	currType := string(p.curr.Lit)
	semType, err := p.returnSemanticType(currType)
	if err != nil {
		return err
	}

	if err := p.addVariablesToSymbolTable(semType, currentVars); err != nil {
		return err
	}

	_, err = p.parseType()

	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) parseType() (string, error) {
	currType := p.curr.Lit
	if err := p.expect(token.TokMap.Type("type")); err != nil {
		return "", err
	}
	return string(currType), nil
}

func (p *Parser) parseFunctionListOpt() error {
	if p.curr.Type != token.TokMap.Type("kwdFunc") {
		return nil
	}

	if p.curr.Type != token.TokMap.Type("kwdFunc") {
		if p.curr.Type != token.TokMap.Type("kwdBegin") {
			return fmt.Errorf("line %d: unexpected token '%s', expected 'var', 'func', or 'begin'",
				p.curr.Line, p.curr.Lit)
		}
		return nil
	}

	if err := p.parseFunctionList(); err != nil {
		return err
	}

	return nil
}

func (p *Parser) parseFunctionList() error {

	for p.curr.Type == token.TokMap.Type("kwdFunc") {
		if err := p.parseFunction(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) parseFunction() error {
	if err := p.expect(token.TokMap.Type("kwdFunc")); err != nil {
		return err
	}

	functionId := p.curr.Lit

	if err := p.expect(token.TokMap.Type("id")); err != nil {
		return err
	}

	if err := p.expect(token.TokMap.Type("openParan")); err != nil {
		return err
	}

	if params, err := p.parseParameterList(); err != nil {
		return err
	} else if err := p.symbolTable.AddFunction(string(functionId), params, p.curr.Line, p.curr.Column); err != nil {
		return err
	}

	if err := p.expect(token.TokMap.Type("closeParan")); err != nil {
		return err
	}

	if err := p.parseBlock(); err != nil {
		return err
	}

	if err := p.expect(token.TokMap.Type("terminator")); err != nil {
		return err
	}

	return p.parseFunctionList()
}

func (p *Parser) parseParameterList() ([]semantic.Variable, error) {
	if p.curr.Type == token.TokMap.Type("closeParan") {
		return []semantic.Variable{}, nil
	}

	currentParams := make([]semantic.Variable, 0)
	currId := p.curr
	currType, err := p.parseParameter()
	if err != nil {
		return []semantic.Variable{}, err
	}
	semType, err := p.returnSemanticType(currType)

	currentParams = append(currentParams, semantic.Variable{
		Name:   string(currId.Lit),
		Type:   semType,
		Line:   currId.Line,
		Column: currId.Column,
	})

	if err != nil {
		return []semantic.Variable{}, err
	}

	for p.curr.Type == token.TokMap.Type("repeatTerminator") {
		p.next() // consume the repeat terminator
		currId := p.curr
		currType, err := p.parseParameter()
		semType, err := p.returnSemanticType(currType)

		if err != nil {
			return []semantic.Variable{}, err
		}

		currentParams = append(currentParams, semantic.Variable{
			Name:   string(currId.Lit),
			Type:   semType,
			Line:   currId.Line,
			Column: currId.Column,
		})

		if err != nil {
			return []semantic.Variable{}, err
		}
	}

	return currentParams, nil
}

func (p *Parser) parseParameter() (string, error) {
	if err := p.expect(token.TokMap.Type("id")); err != nil {
		return "", err
	}

	if err := p.expect(token.TokMap.Type("typeAssignOp")); err != nil {
		return "", err
	}

	currType, err := p.parseType()
	if err != nil {
		return "", err
	}

	return currType, nil
}

func (p *Parser) parseBlock() error {
	if err := p.expect(token.TokMap.Type("openBrace")); err != nil {
		return err
	}

	if err := p.parseStatementList(); err != nil {
		return err
	}

	return p.expect(token.TokMap.Type("closeBrace"))
}

// Parsing Statements
func (p *Parser) parseStatementList() error {
	for {
		validStatementStart, err := p.isStatementStart()
		if err != nil {
			return err
		}

		if !validStatementStart {
			break // No more statements to parse
		}

		if err := p.parseStatement(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseStatement() error {
	validStatementStart, err := p.isStatementStart()

	if !validStatementStart {
		return nil
	}

	if err != nil {
		return err
	}

	switch p.curr.Type {
	case token.TokMap.Type("kwdIf"):
		return p.parseIfStatement()
	case token.TokMap.Type("kwdWhile"):
		return p.parseWhileStatement()
	case token.TokMap.Type("kwdPrint"):
		return p.parsePrintStatement()
	case token.TokMap.Type("id"):
		// Functionality of function calls to be added

		idToken := p.curr
		p.next()
		nextToken := p.curr
		if nextToken.Type == token.TokMap.Type("openParan") {
			if err := p.symbolTable.ValidateFunctionCall(string(idToken.Lit), idToken.Line); err != nil {
				return err
			}
			return p.parseFunctionCall(idToken)
		} else if nextToken.Type == token.TokMap.Type("assignOp") {
			if err := p.symbolTable.ValidateVarAssignment(string(idToken.Lit), idToken.Line); err != nil {
				return err
			}
			return p.parseAssignment(idToken, nextToken)
		} else {
			return fmt.Errorf("Expected either = or (, got %v at line %d, column %d", token.TokMap.Id(nextToken.Type), nextToken.Line, nextToken.Column)
		}
	}
	return nil
}

func (p *Parser) parseAssignment(id, nextToken *token.Token) error {
	if nextToken.Type != token.TokMap.Type("assignOp") {
		return fmt.Errorf("Expected =, got %v at line %d, column %d", nextToken.Type, nextToken.Line, nextToken.Column)
	}
	p.next()

	exprType, err := p.parseExpression()
	varType, err := p.symbolTable.GetType(string(id.Lit))

	if err != nil {
		return err
	}
	// fmt.Println("The expression type is ", exprType, "and the tok", string(id.Lit))

	resultType := p.semanticCube.GetResultType(varType, exprType, "=")
	if resultType == semantic.TypeError {
		return fmt.Errorf("line %d: cannot assign value of type %v to variable '%s' of type %v",
			id.Line, exprType, id.Lit, varType)
	}

	if err != nil {
		return err
	}

	if err := p.expect(token.TokMap.Type("terminator")); err != nil {
		return err
	}

	return nil
}

func (p *Parser) parseWhileStatement() error {
	if err := p.expect(token.TokMap.Type("kwdWhile")); err != nil {
		return err
	}

	if err := p.expect(token.TokMap.Type("openParan")); err != nil {
		return err
	}

	_, err := p.parseExpression()
	if err != nil {
		return err
	}

	if err := p.expect(token.TokMap.Type("closeParan")); err != nil {
		return err
	}

	return p.parseBlock()
}

func (p *Parser) parseIfStatement() error {
	if err := p.expect(token.TokMap.Type("kwdIf")); err != nil {
		return err
	}

	if err := p.expect(token.TokMap.Type("openParan")); err != nil {
		return err
	}

	_, err := p.parseExpression()
	if err != nil {
		return err
	}

	if err := p.expect(token.TokMap.Type("closeParan")); err != nil {
		return err
	}

	if err := p.parseBlock(); err != nil {
		return err
	}

	if p.curr.Type == token.TokMap.Type("kwdElse") {
		p.next()
		return p.parseBlock()
	}

	return nil
}

func (p *Parser) parsePrintStatement() error {
	if err := p.expect(token.TokMap.Type("kwdPrint")); err != nil {
		return err
	}

	if err := p.expect(token.TokMap.Type("openParan")); err != nil {
		return err
	}

	if err := p.parsePrintList(); err != nil {
		return err
	}

	return p.expect(token.TokMap.Type("closeParan"))
}

func (p *Parser) parseFunctionCall(id *token.Token) error {
	// Still not ready
	if err := p.expect(token.TokMap.Type("openParan")); err != nil {
		return err
	}

	if err := p.parseArgumentList(); err != nil {
		return err
	}

	if err := p.expect(token.TokMap.Type("closeParan")); err != nil {
		return err
	}

	return nil
}

func (p *Parser) parseArgumentList() error {
	if p.curr.Type == token.TokMap.Type("closeParan") {
		return nil
	}

	_, err := p.parseExpression()
	if err != nil {
		return err
	}

	for p.curr.Type == token.TokMap.Type("repeatTerminator") {
		p.next()
		_, err := p.parseExpression()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parsePrintList() error {
	if err := p.parsePrintItem(); err != nil {
		return err
	}

	for p.curr.Type == token.TokMap.Type("repeatTerminator") {
		p.next() // consume separator
		if err := p.parsePrintItem(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) parsePrintItem() error {
	if p.curr.Type == token.TokMap.Type("stringLit") {
		p.next()
		return nil
	}

	_, err := p.parseExpression()

	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) parseExpression() (semantic.Type, error) {
	leftType, err := p.parseExp()
	if err != nil {
		return semantic.TypeError, err
	}

	if p.curr.Type == token.TokMap.Type("relOp") {
		operator := string(p.curr.Lit)
		p.next()

		rightType, err := p.parseExp()
		if err != nil {
			return semantic.TypeError, err
		}

		// Check compatibility using semantic cube
		resultType := p.semanticCube.GetResultType(leftType, rightType, operator)
		if resultType == semantic.TypeError {
			return semantic.TypeError, fmt.Errorf("line %d: invalid operation %v %s %v",
				p.curr.Line, leftType, operator, rightType)
		}

		return resultType, nil
	}

	return leftType, nil
}

func (p *Parser) parseExp() (semantic.Type, error) {
	leftType, err := p.parseTerm()
	if err != nil {
		return semantic.TypeError, err
	}

	if p.curr.Type == token.TokMap.Type("expressionOp") {
		operator := string(p.curr.Lit)
		p.next()

		rightType, err := p.parseExp()
		if err != nil {
			return semantic.TypeError, err
		}

		// Check compatibility using semantic cube
		resultType := p.semanticCube.GetResultType(leftType, rightType, operator)
		if resultType == semantic.TypeError {
			return semantic.TypeError, fmt.Errorf("line %d: invalid operation %v %s %v",
				p.curr.Line, leftType, operator, rightType)
		}

		return resultType, nil
	}
	return leftType, nil
}

func (p *Parser) parseTerm() (semantic.Type, error) {
	leftType, err := p.parseFactor()
	if err != nil {
		return semantic.TypeError, err
	}

	if p.curr.Type == token.TokMap.Type("termOp") {
		operator := string(p.curr.Lit)
		p.next()

		rightType, err := p.parseTerm()
		if err != nil {
			return semantic.TypeError, err
		}

		// Check compatibility using semantic cube
		resultType := p.semanticCube.GetResultType(leftType, rightType, operator)
		if resultType == semantic.TypeError {
			return semantic.TypeError, fmt.Errorf("line %d: invalid operation %v %s %v",
				p.curr.Line, leftType, operator, rightType)
		}

		return resultType, nil
	}

	return leftType, nil
}

func (p *Parser) parseFactor() (semantic.Type, error) {
	switch p.curr.Type {
	case token.TokMap.Type("openParan"):
		p.next()
		exprType, err := p.parseExpression()
		if err != nil {
			return semantic.TypeError, err
		}
		if err := p.expect(token.TokMap.Type("closeParan")); err != nil {
			return semantic.TypeError, err
		}
		return exprType, nil
	case token.TokMap.Type("closeParan"):
		p.next()
		switch p.curr.Type {
		case token.TokMap.Type("id"), token.TokMap.Type("intLit"), token.TokMap.Type("floatLit"):
			return p.getType(p.curr)
		default:
			return semantic.TypeError, fmt.Errorf("expected ID, IntLit, or FloatLit after expressionOp")
		}
	case token.TokMap.Type("expressionOp"):
		if string(p.curr.Lit) != "+" && string(p.curr.Lit) != "-" {
			return semantic.TypeError, fmt.Errorf("unexpected operator: %s", p.curr.Lit)
		}
		p.next()

		switch p.curr.Type {
		case token.TokMap.Type("intLit"), token.TokMap.Type("floatLit"):
			return p.getType(p.curr)
		default:
			return semantic.TypeError, fmt.Errorf("expected number after %s", p.curr.Lit)
		}
	case token.TokMap.Type("id"), token.TokMap.Type("intLit"), token.TokMap.Type("floatLit"):
		tok := p.curr
		return p.getType(tok)
	default:
		return semantic.TypeError, fmt.Errorf("unexpected token in factor: %v",
			token.TokMap.Id(p.curr.Type))
	}
}

func (p *Parser) parseMainSection() error {
	if err := p.expect(token.TokMap.Type("kwdBegin")); err != nil {
		return err
	}

	if err := p.parseStatementList(); err != nil {
		return err
	}

	if err := p.expect(token.TokMap.Type("kwdEnd")); err != nil {
		return err
	}

	return nil
}
