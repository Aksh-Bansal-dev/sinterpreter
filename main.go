package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	mul = "MUL"
	div = "DIV"
	add = "ADD"
	sub = "SUB"
	num = "NUM"
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
	case add:
		return lVal.(int) + rVal.(int)
	case sub:
		return lVal.(int) - rVal.(int)
	case mul:
		return lVal.(int) * rVal.(int)
	case div:
		return lVal.(int) / rVal.(int)
	}
	return nil
}

type NumNode struct {
	token Token
}

func (node NumNode) getType() string {
	return "numnode"
}
func (node NumNode) eval() interface{} {
	val, err := strconv.Atoi(node.token.val)
	if err != nil {
		log.Fatal("Invalid number at", node.token.pos)
	}
	return val
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
		} else if c >= '0' && c <= '9' {
			tokens = append(tokens, Token{num, i, getNum(&i, s)})
			i--
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
	return res
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
	t := p.tokens[p.cur]
	if t.tt == num {
		p.next()
		return NumNode{token: t}
	}
	log.Fatal("Expected number at", t.pos)
	return nil
}

func (p *Parser) factor() Node {
	left := p.primary()
	for p.match([]string{mul, div}) {
		op := p.prev()
		right := p.primary()
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

func (p *Parser) parse() Node {
	return p.term()
}

// Interpreter

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
