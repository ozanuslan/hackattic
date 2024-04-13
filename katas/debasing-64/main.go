package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

/*
Base64 is sort of a running gag around here. It's so widely used for the classic challenges that this one thing is certain – base64 or die.

Your task is to read in a Base64-encoded string, decode it and print it out. Easy money.

Sample input
bGF0ZS1hdC1uaWdodA==
d2l0aC10aGUtcmlzaW5nLWFwZQ==
dGhlLXJ1dGhsZXNzLXNldmVu

Sample output
late-at-night
with-the-rising-ape
the-ruthless-seven
*/

var reverseLookup [128]byte // byte lookup table containing ASCII codes for reversing base64 encoded bytes

func main() {
	reverseLookup = [128]byte{
		80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, /* 0 - 15 */
		80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, /* 16 - 31 */
		80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 62, 80, 80, 80, 63, /* 32 - 47 */
		52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 80, 80, 80, 64, 80, 80, /* 48 - 63 */
		80, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, /* 64 - 79 */
		15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 80, 80, 80, 80, 80, /* 80 - 96 */
		80, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, /* 87 - 111 */
		41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 80, 80, 80, 80, 80, /* 112 - 127 */
	}

	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	for _, line := range bytes.Split(stdin, []byte("\n")) {
		fmt.Println(string(base64Decode(line)))
	}
}

func base64Decode(encoded []byte) []byte {
	var decoded bytes.Buffer

	// algorithm source: http://www.sunshine2k.de/articles/coding/base64/understanding_base64.html#ch43
	var c1, c2, c3 byte
	for i := 0; i < len(encoded); i += 4 {
		b1 := reverseLookup[encoded[i]]
		b2 := reverseLookup[encoded[i+1]]
		b3 := reverseLookup[encoded[i+2]]
		b4 := reverseLookup[encoded[i+3]]

		if encoded[i+3] == '=' {
			if encoded[i+2] == '=' {
				c1 = b1<<2 | (b2&0xF0)>>4
				decoded.WriteByte(c1)
			} else {
				c1 = b1<<2 | (b2&0xF0)>>4
				c2 = (b2&0x0F)<<4 | (b3&0x3C)>>2
				decoded.WriteByte(c1)
				decoded.WriteByte(c2)
			}
		} else {
			c1 = b1<<2 | (b2&0xF0)>>4
			c2 = (b2&0x0F)<<4 | (b3&0x3C)>>2
			c3 = (b3&0x03)<<6 | (b4 & 0x3F)
			decoded.WriteByte(c1)
			decoded.WriteByte(c2)
			decoded.WriteByte(c3)
		}
	}

	return decoded.Bytes()
}
