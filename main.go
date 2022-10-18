package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	mul   = "MUL"
	div   = "DIV"
	add   = "ADD"
	sub   = "SUB"
	num   = "NUM"
	lPara = "L_PARA"
	rPara = "R_PARA"
	equal = "EQUALS"
	// TODO
	boolT              = "BOOL"
	not                = "NOT"
	notEqual           = "NOT_EQUALS"
	lessThanOrEqual    = "LESS_THAN_OR_EQUALS"
	greaterThanOrEqual = "GREATER_THAN_OR_EQUALS"
	lessThan           = "GREATER"
	greaterThan        = "LESSER"
)

type Token struct {
	tt  string
	pos int
	val string
}

type Node interface {
	eval() interface{}
	getType() string
}

type BinOp struct {
	left  Node
	op    Token
	right Node
}

func (node BinOp) getType() string {
	return "binop"
}

func (node BinOp) eval() interface{} {
	lVal := node.left.eval()
	rVal := node.right.eval()
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
func (node UnaryOp) eval() interface{} {
	rVal := node.right.eval()
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
func (node PrimaryNode) eval() interface{} {
	if node.token.tt == num {
		val, err := strconv.Atoi(node.token.val)
		if err != nil {
			log.Fatal("Invalid number at", node.token.pos)
		}
		return val
	} else if node.token.tt == boolT {
		val := true
		if node.token.val == "false" {
			val = false
		}
		return val
	}
	log.Fatal("invalid value")
	return nil
}

type Parser struct {
	cur    int
	tokens []Token
}

// Lexer
func lexer(s string) []Token {
	tokens := []Token{}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '+' {
			tokens = append(tokens, Token{add, i, ""})
		} else if c == '-' {
			tokens = append(tokens, Token{sub, i, ""})
		} else if c == '*' {
			tokens = append(tokens, Token{mul, i, ""})
		} else if c == '/' {
			tokens = append(tokens, Token{div, i, ""})
		} else if c == ')' {
			tokens = append(tokens, Token{rPara, i, ""})
		} else if c == '(' {
			tokens = append(tokens, Token{lPara, i, ""})
		} else if c >= '0' && c <= '9' {
			tokens = append(tokens, Token{num, i, getNum(&i, s)})
		} else if getKeyword(&i, s, []string{"true"}) != "" {
			tokens = append(tokens, Token{boolT, i, "true"})
		} else if getKeyword(&i, s, []string{"false"}) != "" {
			tokens = append(tokens, Token{boolT, i, "false"})
		} else if getKeyword(&i, s, []string{">="}) != "" {
			tokens = append(tokens, Token{greaterThanOrEqual, i, ""})
		} else if getKeyword(&i, s, []string{"<="}) != "" {
			tokens = append(tokens, Token{lessThanOrEqual, i, ""})
		} else if getKeyword(&i, s, []string{"<"}) != "" {
			tokens = append(tokens, Token{lessThan, i, ""})
		} else if getKeyword(&i, s, []string{">"}) != "" {
			tokens = append(tokens, Token{greaterThan, i, ""})
		} else if getKeyword(&i, s, []string{"=="}) != "" {
			tokens = append(tokens, Token{equal, i, ""})
		} else if getKeyword(&i, s, []string{"!="}) != "" {
			tokens = append(tokens, Token{notEqual, i, ""})
		} else if c == '!' {
			tokens = append(tokens, Token{not, i, ""})
		} else if c == ' ' {
			continue
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

// Parser
func (p *Parser) next() Token {
	if p.cur >= len(p.tokens) {
		log.Fatal("EOF reached!")
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
func (p *Parser) consume(reqTokens []string, errMsg string) bool {
	t := p.next()
	for _, reqToken := range reqTokens {
		if reqToken == t.tt {
			return true
		}
	}
	log.Fatal(errMsg)
	return false
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
	if p.match([]string{num, boolT}) {
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

func (p *Parser) parse() Node {
	return p.comparision()
}

func main() {
	log.SetFlags(log.Lshortfile)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("♪ ")
		inp, _, _ := reader.ReadLine()
		tokens := lexer(string(inp))
		// fmt.Println(tokens)
		parser := Parser{0, tokens}
		ast := parser.parse()
		// fmt.Println(ast)
		res := ast.eval()
		fmt.Println(res)
	}

}
