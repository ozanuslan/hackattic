package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

/*
For every line you read in, return the same line but in compressed format.

The compression algorithm is simple: if a character repeats itself more than twice in a row, replace it with <number_of_occurences><character>.

You now know the basis of basically 90% of all compression algorithms... maybe.

Sample input
aaaaaiiiixqvsm
rrdkuuuuyyyrrrrgghc
xhzzzccccvvsssqppc
jbiiiulllllvvvvtttttxxxxxs

Sample output
5a4ixqvsm
rdk4u3y4rghc
xh3z4cv3sqpc
jb3iu5l4v5t5xs
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
	compressed := compressLine(line)
	fmt.Println(compressed)
}

func compressLine(line []byte) string {
	var compressed bytes.Buffer
	last := line[0]
	cnt := 1
	for _, c := range line[1:] {
		if last == c {
			cnt++
			continue
		}
		if cnt <= 2 {
			for i := 0; i < cnt; i++ {
				compressed.WriteByte(last)
			}
		} else {
			compressed.WriteString(strconv.Itoa(cnt))
			compressed.WriteByte(last)
		}
		last = c
		cnt = 1
	}
	if cnt <= 2 {
		for i := 0; i < cnt; i++ {
			compressed.WriteByte(last)
		}
	} else {
		compressed.WriteString(strconv.Itoa(cnt))
		compressed.WriteByte(last)
	}
	return compressed.String()
}
