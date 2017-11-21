package parser

import (
	"strconv"

	"github.com/Zac-Garby/lang/ast"
	"github.com/Zac-Garby/lang/token"
)

// parseExpression parses an expression starting at the current
// token, and finishing just after the end of the node.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix, ok := p.prefixes[p.cur.Type]
	if !ok {
		p.unexpectedTokenErr(p.cur.Type)
		return nil
	}

	left := prefix()

	for !p.peekIs(token.Semi) && precedence < p.peekPrecedence() {
		infix, ok := p.infixes[p.peek.Type]
		if !ok {
			return left
		}

		p.next()
		left = infix(left)
	}

	return left
}

// Prefix Expression Parsers

func (p *Parser) parseID() ast.Expression {
	return &ast.Identifier{
		Tok:   p.cur,
		Value: p.cur.Literal,
	}
}

func (p *Parser) parseNum() ast.Expression {
	node := &ast.Number{
		Tok: p.cur,
	}

	val, err := strconv.ParseFloat(p.cur.Literal, 64)
	if err != nil {
		p.Errors = append(p.Errors, err)
		return nil
	}

	node.Value = val
}

func (p *Parser) parseBool() ast.Expression {
	return &ast.Boolean{
		Tok:   p.cur,
		Value: p.cur.Type == token.True,
	}
}

func (p *Parser) parseNil() ast.Expression {
	return &ast.Nil{
		Tok: p.cur,
	}
}

func (p *Parser) parseString() ast.Expression {
	return &ast.String{
		Tok:   p.cur,
		Value: p.cur.Literal,
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.next()

	if p.curIs(token.RightParen) {
		return &ast.Tuple{
			Tok: p.cur,
		}
	}

	var (
		node    = p.parseExpression(lowest)
		isTuple = false
	)

	if p.peekIs(token.Comma) {
		isTuple = true

		p.next()
		p.next()

		node = &ast.Tuple{
			Tok: node.Token(),
			Value: append(
				[]ast.Expression{node},
				p.parseExpressionList(token.RightParen)...,
			),
		}
	}

	if !isTuple && !p.expect(token.RightParen) {
		return nil
	}

	return node
}

func (p *Parser) parseList() ast.Expression {
	p.next()

	if p.curIs(token.RightSquare) {
		return &ast.List{
			Tok:      p.cur,
			Elements: []ast.Expression{},
		}
	}

	return &ast.Array{
		Tok:      p.cur,
		Elements: p.parseExpressionList(token.RightSquare),
	}
}

func (p *Parser) parseMap() ast.Expression {
	p.next()

	if !p.expect(token.LeftBrace) {
		return nil
	}

	if p.curIs(token.RightBrace) {
		return &ast.Map{
			Tok:   p.cur,
			Pairs: make(map[ast.Expression]ast.Expression),
		}
	}

	return &ast.Map{
		Tok:   p.cur,
		Pairs: p.parseExpressionPairs(token.RightBrace),
	}
}

func (p *Parser) parseSet() ast.Expression {
	p.next()

	if !p.expect(token.LeftBrace) {
		return nil
	}

	if p.curIs(token.RightBrace) {
		return &ast.List{
			Tok:      p.cur,
			Elements: []ast.Expression{},
		}
	}

	return &ast.Array{
		Tok:      p.cur,
		Elements: p.parseExpressionList(token.RightBrace),
	}
}

func (p *Parser) parseBlock() ast.Expression {
	node := &ast.Block{
		Tok:        p.cur,
		Statements: make([]ast.Statement, 0, 8),
	}

	p.next()

	for !p.curIs(token.RightBrace) && !p.curIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			node.Statements = append(node.Statements, stmt)
		}

		p.next()
	}

	return node
}

func (p *Parser) parsePrefix() ast.Expression {
	node := &ast.PrefixExpression{
		Tok:      p.cur,
		Operator: p.cur.Literal,
	}

	p.next()
	node.Right = p.parseExpression(prefix)

	return node
}

func (p *Parser) parseIfExpression() ast.Expression {
	node := &ast.IfExpression{
		Tok: p.cur,
	}

	p.next()
	node.Condition = p.parseExpression(lowest)

	if !p.expect(token.Then) {
		return nil
	}

	node.Consequence = p.parseExpression(lowest)

	if p.peekIs(token.Else) {
		p.next()

		if !p.expect(token.LeftBrace) {
			return nil
		}

		node.Alternative = p.parseExpression(lowest)
	}

	return node
}

func (p *Parser) parseMatchExpression() ast.Expression {
	node := &ast.MatchExpression{
		Tok: p.cur,
	}

	p.next()
	node.Input = p.parseExpression(lowest)

	for p.peekIs(token.If) {
		pair := ast.MatchBranch{}

		p.next()
		pair.Condition = p.parseExpression(lowest)

		if !p.expect(token.RightArrow) {
			return nil
		}

		pair.Body = p.parseExpression(lowest)
		node.Branches = append(node.Branches, pair)

		p.next()
	}

	return node
}

func (p *Parser) parseType() ast.Expression {
	node := &ast.Type{
		Tok: p.cur,
	}

	if !p.expect(token.LeftParen) {
		return nil
	}

	p.next()

	node.Parameters = p.parseParams(token.RightParen)

	return node
}

// Infix Expression Parsers

func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	node := &ast.InfixExpression{
		Tok:      p.cur,
		Operator: p.cur.Literal,
		Left:     left,
	}

	precedence := p.cur()
	p.next()
	node.Right = p.parseExpression(precedence)

	return node
}

func (p *Parser) parseIndex(left ast.Expression) ast.Expression {
	node := &ast.IndexExpression{
		Tok:  p.cur,
		Left: left,
	}

	p.next()
	node.Index = p.parseExpression(lowest)

	if !p.expect(token.RightSquare) {
		return nil
	}

	return node
}

func (p *Parser) parseFunctionCall(left ast.Expression) ast.Expression {
	node := &ast.FunctionCall{
		Tok:      p.cur,
		Function: left,
	}

	node.Arguments = p.parseExpressionList(token.RightParen)

	return node
}
