package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*
The "Fizz-Buzz test" is an interview question designed to help filter out the 99.5% of programming job candidates who can't seem to program their way out of a wet paper bag.

Read in a single line with two numbers, N and M. Then simply print all the numbers from N to M. But for multiples of three print “Fizz” instead of the number and for the multiples of five print “Buzz”. For numbers which are multiples of both three and five print “FizzBuzz”.

About the legend
If you've never been graced with the FizzBuzz test and the havoc it causes in the real world, you have some catching up to do! See here, here, here and here.

Sample input
8 16

Sample output
8
Fizz
Buzz
11
Fizz
13
14
FizzBuzz
16
*/

func main() {
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	nm := strings.Split(strings.TrimSpace(string(stdin)), " ")
	if len(nm) < 2 {
		panic(fmt.Errorf("err: range must have 2 elements, got %d", len(nm)))
	}

	n, err := strconv.Atoi(nm[0])
	if err != nil {
		panic(err)
	}

	m, err := strconv.Atoi(nm[1])
	if err != nil {
		panic(err)
	}

	fizzBuzz(n, m)
}

func fizzBuzz(n, m int) {
	for i := n; i <= m; i++ {
		switch true {
		case i%3 == 0 && i%5 == 0:
			fmt.Println("FizzBuzz")
		case i%3 == 0:
			fmt.Println("Fizz")
		case i%5 == 0:
			fmt.Println("Buzz")
		default:
			fmt.Println(i)
		}
	}
}
