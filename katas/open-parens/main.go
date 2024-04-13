package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

/*
Every line you read in from STDIN will contain a bunch of parentheses.

Your task is to determine if they are properly nested – i.e. if every opening parenthesis has a closing one.

Sample input
(())
()))
(()((())))
(()(()(()))

Sample output
yes
no
yes
no
*/

func main() {
	rdr := bufio.NewReader(os.Stdin)

	for {
		switch line, err := rdr.ReadBytes('\n'); err {
		case nil:
			line = line[:len(line)-1]
			if len(line) < 1 {
				fmt.Println()
				continue
			}
			handleLine(line)
		case io.EOF:
			if len(line) >= 1 {
				handleLine(line)
			}
			os.Exit(0)
		default:
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	}
}

func handleLine(line []byte) {
	isProper := isProperlyNested(line)
	if isProper {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}
}

type Stack []byte

func (s *Stack) Push(e byte) {
	*s = append(*s, e)
}

func (s *Stack) Pop() byte {
	e := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return e
}

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Len() int {
	return len(*s)
}

func isProperlyNested(line []byte) bool {
	stk := make(Stack, 0)
	for _, paren := range line {
		switch paren {
		case '(':
			stk.Push(paren)
		case ')':
			if stk.IsEmpty() {
				return false
			}
			if stk.Pop() != '(' {
				return false
			}
		default:
			panic(fmt.Errorf("err: expected parenthesis char, got %c", paren))
		}
	}
	return stk.Len() == 0
}
