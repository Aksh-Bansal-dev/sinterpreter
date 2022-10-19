package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	mul                = "MUL"
	div                = "DIV"
	add                = "ADD"
	sub                = "SUB"
	num                = "NUM"
	lPara              = "L_PARA"
	rPara              = "R_PARA"
	equal              = "EQUALS"
	boolT              = "BOOL"
	not                = "NOT"
	notEqual           = "NOT_EQUALS"
	lessThanOrEqual    = "LESS_THAN_OR_EQUALS"
	greaterThanOrEqual = "GREATER_THAN_OR_EQUALS"
	lessThan           = "GREATER"
	greaterThan        = "LESSER"
	semicolon          = "SEMICOLON"
	printT             = "PRINT"
	varT               = "VAR"
	identifierT        = "IDENTIFIER"
	assignT            = "ASSIGN"
)

type Token struct {
	tt   string
	col  int
	line int
	val  string
}

type Node interface {
	eval(env map[string]interface{}) interface{}
	getType() string
}
type VarDecl struct {
	identifier Token
	expr       Node
}

func (node VarDecl) getType() string {
	return "var_decl"
}
func (node VarDecl) eval(env map[string]interface{}) interface{} {
	env[node.identifier.val] = node.expr.eval(env)
	return nil
}

type PrintStmt struct {
	expr Node
}

func (node PrintStmt) getType() string {
	return "print_stmt"
}
func (node PrintStmt) eval(env map[string]interface{}) interface{} {
	fmt.Println(node.expr.eval(env))
	return nil
}

type BinOp struct {
	left  Node
	op    Token
	right Node
}

func (node BinOp) getType() string {
	return "binop"
}

func (node BinOp) eval(env map[string]interface{}) interface{} {
	lVal := node.left.eval(env)
	rVal := node.right.eval(env)
	if lVal == nil || rVal == nil {
		log.Fatal("Invalid operation")
	}
	switch node.op.tt {
	case equal:
		return lVal == rVal
	case notEqual:
		return lVal != rVal
	case lessThan:
		return lVal.(int) < rVal.(int)
	case lessThanOrEqual:
		return lVal.(int) <= rVal.(int)
	case greaterThan:
		return lVal.(int) > rVal.(int)
	case greaterThanOrEqual:
		return lVal.(int) >= rVal.(int)
	case add:
		return lVal.(int) + rVal.(int)
	case sub:
		return lVal.(int) - rVal.(int)
	case mul:
		return lVal.(int) * rVal.(int)
	case div:
		if rVal.(int) == 0 {
			log.Fatal("Division by 0")
		}
		return lVal.(int) / rVal.(int)
	}
	return nil
}

type UnaryOp struct {
	token Token
	right Node
}

func (node UnaryOp) getType() string {
	return "unary_op"
}
func (node UnaryOp) eval(env map[string]interface{}) interface{} {
	rVal := node.right.eval(env)
	if node.token.tt == sub {
		return -(rVal.(int))
	} else if node.token.tt == not {
		return !(rVal.(bool))
	}
	return rVal
}

type PrimaryNode struct {
	token Token
}

func (node PrimaryNode) getType() string {
	return "PrimaryNode"
}
func (node PrimaryNode) eval(env map[string]interface{}) interface{} {
	if node.token.tt == num {
		val, err := strconv.Atoi(node.token.val)
		if err != nil {
			log.Fatal("Invalid number at", node.token.col)
		}
		return val
	} else if node.token.tt == boolT {
		val := true
		if node.token.val == "false" {
			val = false
		}
		return val
	} else if node.token.tt == identifierT {
		if _, ok := env[node.token.val]; ok {
			return env[node.token.val]
		}
		log.Fatal("Undefined variable ", node.token.val)
	}
	log.Fatal("invalid value")
	return nil
}

type Parser struct {
	cur    int
	tokens []Token
}

// Lexer
func lexer(s string, line int) []Token {
	tokens := []Token{}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '+' {
			tokens = append(tokens, Token{add, i, line, ""})
		} else if c == '-' {
			tokens = append(tokens, Token{sub, i, line, ""})
		} else if c == '*' {
			tokens = append(tokens, Token{mul, i, line, ""})
		} else if c == '/' {
			tokens = append(tokens, Token{div, i, line, ""})
		} else if c == ')' {
			tokens = append(tokens, Token{rPara, i, line, ""})
		} else if c == '(' {
			tokens = append(tokens, Token{lPara, i, line, ""})
		} else if c >= '0' && c <= '9' {
			tokens = append(tokens, Token{num, i, line, getNum(&i, s)})
		} else if getKeyword(&i, s, []string{"true"}) != "" {
			tokens = append(tokens, Token{boolT, i, line, "true"})
		} else if getKeyword(&i, s, []string{"print"}) != "" {
			tokens = append(tokens, Token{printT, i, line, ""})
		} else if getKeyword(&i, s, []string{"false"}) != "" {
			tokens = append(tokens, Token{boolT, i, line, "false"})
		} else if getKeyword(&i, s, []string{"var"}) != "" {
			tokens = append(tokens, Token{varT, i, line, ""})
		} else if getKeyword(&i, s, []string{">="}) != "" {
			tokens = append(tokens, Token{greaterThanOrEqual, i, line, ""})
		} else if getKeyword(&i, s, []string{"<="}) != "" {
			tokens = append(tokens, Token{lessThanOrEqual, i, line, ""})
		} else if getKeyword(&i, s, []string{"<"}) != "" {
			tokens = append(tokens, Token{lessThan, i, line, ""})
		} else if getKeyword(&i, s, []string{">"}) != "" {
			tokens = append(tokens, Token{greaterThan, i, line, ""})
		} else if getKeyword(&i, s, []string{"=="}) != "" {
			tokens = append(tokens, Token{equal, i, line, ""})
		} else if getKeyword(&i, s, []string{"!="}) != "" {
			tokens = append(tokens, Token{notEqual, i, line, ""})
		} else if getKeyword(&i, s, []string{"="}) != "" {
			tokens = append(tokens, Token{assignT, i, line, ""})
		} else if c == '!' {
			tokens = append(tokens, Token{not, i, line, ""})
		} else if c == ';' {
			tokens = append(tokens, Token{semicolon, i, line, ""})
		} else if c == ' ' {
			continue
		} else if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			tokens = append(tokens, Token{identifierT, i, line, getIdentifier(&i, s)})
		} else {
			log.Fatal("Invalid token")
		}

	}
	return tokens
}
func getNum(i *int, s string) string {
	res := ""
	for ; *i < len(s) && s[*i] <= '9' && s[*i] >= '0'; *i++ {
		res += string(s[*i])
	}
	*i--
	return res
}
func getKeyword(idx *int, s string, keywords []string) string {
	i := *idx
	for _, keyword := range keywords {
		if len(keyword)+i <= len(s) && keyword == s[i:i+len(keyword)] {
			*idx = i + len(keyword) - 1
			return keyword
		}
	}
	return ""
}
func getIdentifier(i *int, s string) string {
	res := ""
	for ; *i < len(s) &&
		s[*i] != ' ' &&
		(s[*i] <= 'z' && s[*i] >= 'a') ||
		(s[*i] <= 'Z' && s[*i] >= 'A') ||
		(s[*i] <= '9' && s[*i] >= '0'); *i++ {
		res += string(s[*i])
	}
	*i--
	return res
}

// Parser
func (p *Parser) next() Token {
	if p.cur >= len(p.tokens) {
		log.Fatal("EOF reached")
	}
	t := p.tokens[p.cur]
	p.cur++
	return t
}
func (p *Parser) prev() Token {
	if p.cur <= 0 {
		log.Fatal("Negative index")
	}
	return p.tokens[p.cur-1]
}
func (p *Parser) consume(reqTokens []string, errMsg string) Token {
	if p.cur >= len(p.tokens) {
		log.Fatal(errMsg)
	}
	t := p.next()
	for _, reqToken := range reqTokens {
		if reqToken == t.tt {
			return t
		}
	}
	log.Fatal(errMsg)
	return t
}
func (p *Parser) isEnd() bool {
	return p.cur == len(p.tokens)
}
func (p *Parser) match(reqTokens []string) bool {
	if p.isEnd() {
		return false
	}
	t := p.tokens[p.cur]
	for _, reqToken := range reqTokens {
		if reqToken == t.tt {
			p.next()
			return true
		}
	}
	return false
}

func (p *Parser) primary() Node {
	if p.match([]string{num, boolT, identifierT}) {
		return PrimaryNode{token: p.prev()}
	} else if p.match([]string{lPara}) {
		innerNode := p.term()
		p.consume([]string{rPara}, fmt.Sprintf("Expected ) at %d", p.cur))
		return innerNode
	}
	log.Fatal("Expected number at ", p.cur)
	return nil
}
func (p *Parser) unary() Node {
	for p.match([]string{sub, not}) {
		op := p.prev()
		right := p.primary()
		uNode := UnaryOp{op, right}
		return uNode
	}
	return p.primary()
}

func (p *Parser) factor() Node {
	left := p.unary()
	for p.match([]string{mul, div}) {
		op := p.prev()
		right := p.unary()
		temp := BinOp{left: left, op: op, right: right}
		left = temp
	}
	return left
}

func (p *Parser) term() Node {
	left := p.factor()
	for p.match([]string{add, sub}) {
		op := p.prev()
		right := p.factor()
		temp := BinOp{left: left, op: op, right: right}
		left = temp
	}
	return left
}

func (p *Parser) comparision() Node {
	left := p.term()
	for p.match([]string{equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual}) {
		op := p.prev()
		right := p.comparision()
		temp := BinOp{left: left, op: op, right: right}
		left = temp
	}
	return left
}

func (p *Parser) printStmt() Node {
	expr := p.comparision()
	p.consume([]string{semicolon}, "expected ; after expression")
	return PrintStmt{expr}
}

func (p *Parser) varDecl() Node {
	identifier := p.consume([]string{identifierT}, "Expected a variable name")
	p.consume([]string{assignT}, "expected = after identifier")
	return VarDecl{identifier, p.exprStmt()}
}

func (p *Parser) exprStmt() Node {
	expr := p.comparision()
	p.consume([]string{semicolon}, "expected ; after expression")
	return expr
}

func (p *Parser) statement() Node {
	if p.match([]string{printT}) {
		return p.printStmt()
	} else if p.match([]string{varT}) {
		return p.varDecl()
	}
	return p.exprStmt()
}

func (p *Parser) parse() Node {
	return p.statement()
}

func main() {
	log.SetFlags(log.Lshortfile)
	filePath := os.Args[1]
	if filePath[len(filePath)-3:] != ".si" {
		log.Fatal("Invalid file extension")
	}
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("File doesn't exists")
	}
	reader := bufio.NewReader(file)
	lineNum := 1
	env := map[string]interface{}{}
	for {
		inp, _, err := reader.ReadLine()
		if err != nil {
			return
		}
		tokens := lexer(string(inp), lineNum)
		// fmt.Println(tokens)
		parser := Parser{0, tokens}
		ast := parser.parse()
		// fmt.Println(ast)
		ast.eval(env)
		// fmt.Println(res)
		lineNum++
	}

}
