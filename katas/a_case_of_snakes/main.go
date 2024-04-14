package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

/*
The best thing about standards is that there's so many to choose from!

This kata is about converting hungarian notation to plain snake case.

For every line of input – which contain some totes real-world variable names – print out the variable in a simplified snake-case form. Strip the type prefix, flip the uppercase letters and smack an underscore or two wherever if makes sense. That's it! Deceivingly simple.

Fun facts
A study found snake_case to be more efficient to read than CamelCase.

Sample input
szWindowContents
iAirflowParameter
fMixtureRatio

Sample output
window_contents
airflow_parameter
mixture_ratio
*/

func main() {
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(stdin), "\n")
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if len(line) < 1 {
			fmt.Println()
			continue
		}
		snake := hungarianToSnake(line)
		fmt.Println(snake)
	}
}

func hungarianToSnake(hungarian string) string {
	varName := filterOutPrefix(hungarian)
	if len(varName) < 1 {
		return ""
	}

	var snake strings.Builder
	snake.WriteRune(unicode.ToLower(rune(varName[0])))
	for _, c := range varName[1:] {
		if unicode.IsUpper(c) {
			snake.WriteRune('_')
		}
		snake.WriteRune(unicode.ToLower(c))
	}

	return snake.String()
}

/*
Hungarian Notation type prefixes from: https://en.wikipedia.org/wiki/Hungarian_notation#Examples

bBusy : boolean
chInitial : char
cApples : count of items
dwLightYears : double word (Systems)
fBusy : flag (or float)
nSize : integer (Systems) or count (Apps)
iSize : integer (Systems) or index (Apps)
fpPrice : floating-point
decPrice : decimal
dbPi : double (Systems)
pFoo : pointer
rgStudents : array, or range
szLastName : zero-terminated string
u16Identifier : unsigned 16-bit integer (Systems)
u32Identifier : unsigned 32-bit integer (Systems)
stTime : clock time structure
fnFunction : function name
*/
var prefixes = []string{
	"b",
	"f",
	"c",
	"p",
	"w",
	"d",
	"ch",
	"dw",
	"rg",
	"sz",
	"st",
	"fn",
	"fp",
	"dec",

	"i8",
	"i16",
	"i32",
	"i64",
	"i128",
	"u8",
	"u16",
	"u32",
	"u64",
	"u128",
}

func filterOutPrefix(hungarian string) string {
	for _, prefix := range prefixes {
		if len(prefix) >= len(hungarian) {
			continue
		}
		if prefix == hungarian[:len(prefix)] && unicode.IsUpper(rune(hungarian[len(prefix)])) {
			return hungarian[len(prefix):]
		}
	}
	return hungarian
}
