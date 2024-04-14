package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

/*
For every line you get from STDIN, add the... things on that line and print out the sum.

Hex, octal, bin and base-10 numbers should give you no problem. Should you stumble upon an ASCII character, take it's ASCII value.

Sample input
110 0x187 300 T / d
180 A 0x10e 0x18c N 95
423 0xac 417 0o20 q &
0x14e 0b10000 247 284 0o447 268

Sample output
1032
1084
1179
1444
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
	sum := sumOfThings(line)
	fmt.Println(sum)
}

func sumOfThings(line []byte) int {
	groups := strings.Split(string(line), " ")
	sum := 0
	for _, group := range groups {
		group = strings.TrimSpace(group)
		if len(group) == 0 {
			continue
		}
		var n int
		if len(group) == 1 {
			r, _ := utf8.DecodeRuneInString(group)
			if unicode.IsDigit(r) {
				n, _ = strconv.Atoi(group)
			} else {
				n = int(r) // ASCII value
			}
		} else {
			var num int64
			var err error
			switch group[:2] {
			case "0b": // binary
				num, err = strconv.ParseInt(group[2:], 2, 0)
			case "0o": // octal
				num, err = strconv.ParseInt(group[2:], 8, 0)
			case "0x": // hex
				num, err = strconv.ParseInt(group[2:], 16, 0)
			default: // decimal
				num, err = strconv.ParseInt(group, 10, 0)
			}
			if err != nil {
				panic(err)
			}
			n = int(num)
		}
		sum += n
	}
	return sum
}
