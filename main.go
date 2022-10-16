package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
	val string
}

type BinOp struct {
	left  *BinOp
	root  Token
	right *BinOp
}
type NumNode struct {
	token Token
}

type Parser struct {
	cur    int
	tokens []Token
}

func lexer(s string) []Token {
	tokens := []Token{}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '+' {
			tokens = append(tokens, Token{add, ""})
		} else if c == '-' {
			tokens = append(tokens, Token{sub, ""})
		} else if c == '*' {
			tokens = append(tokens, Token{mul, ""})
		} else if c == '/' {
			tokens = append(tokens, Token{div, ""})
		} else if c >= '0' && c <= '9' {
			tokens = append(tokens, Token{num, getNum(&i, s)})
			i--
		}
	}
	return tokens
}
func getNum(i *int, s string) string {
	res := ""
	for ; *i < len(s) && s[*i] <= '9' && s[*i] >= '0'; *i++ {
		res = string(s[*i]) + res
	}
	return res
}

func (p *Parser) next() Token {
	if p.cur >= len(p.tokens) {
		log.Fatal("EOF reached!")
	}
	t := p.tokens[p.cur]
	p.cur++
	return t
}
func (p *Parser) match(reqTokens []Token, errMsg string) bool {
	t := p.next()
	for _, reqToken := range reqTokens {
		if reqToken == t {
			return true
		}
	}
	log.Fatal(errMsg)
	return false
}

func (p *Parser) factor() BinOp {
	t := p.tokens[p.cur]
	if t.tt == num {
		p.next()
		return
	}
}

func (p *Parser) parse() BinOp {

}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("♪ ")
		inp, _, _ := reader.ReadLine()
		tokens := lexer(string(inp))
		fmt.Println(tokens)
	}

}
