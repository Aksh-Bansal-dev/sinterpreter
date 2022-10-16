package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("♪ ")
		inp, _, _ := reader.ReadLine()
		tokens := lexer(string(inp))
		// fmt.Printf("%+v", tokens)
		fmt.Println(tokens)
	}

}

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

func parser() {

}
