package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

/*
Every line you receive on STDIN will be binary representation of a number. Only instead of zeros you'll get . and instead of ones you'll get #.

Numbers are all at most 16 bits.

Least significant bit comes last.

Simple enough!

It's a good place to start! A bit of strings, a bit of arrays, a bit of the most basic data types.

Sample input
#.#.#.###.#.##.#
##.##.......#..#
#..#####..#.#...
###..#....###.##
#..#..#.#....##
#......#.##....
#############.#.

Sample output
43949
55305
40744
58427
18755
16560
655305
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
			handleLine(line)
			os.Exit(0)
		default:
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	}
}

func handleLine(line []byte) {
	num, err := constructInt(line)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("%d\n", num)
}

func constructInt(chars []byte) (uint16, error) {
	var num uint16 = 0
	for i, c := range chars {
		switch c {
		case '.': // 0
		case '#': // 1
			num |= 1
		default:
			return 0, fmt.Errorf("unrecognized char %c", c)
		}
		if i != len(chars)-1 {
			num <<= 1
		}
	}
	return num, nil
}
