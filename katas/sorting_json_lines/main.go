package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
Read in the lines, parse the (slightly uwieldy) object, and return the names in sorted order, according to the balance value.

One trick though - if the entry contains an "extra" object, the balance value inside takes precedence.

And please format large numbers using a , as a thousands-separator.

Yeah, someone really didn't think that JSON structure through, right? Devs nowadays, I tell you...

Patchy-patchy, patchy-fixy.

Sample input
{"Bentley.G":{"balance":2134,"account_no":233831255"}}
{"Barclay.E":{"balance":1123,"account_no":312333321}}
{"Alton.K":{"balance":9315,"account_no":203123613,"extra":{"balance":131}}}
{"Bancroft.M":{"balance": 233,"account_no":287655771101,"extra":{"balance":98}}

Sample output
Bancroft.M: 98
Alton.K: 131
Barclay.E: 1,123
Bentley.G: 2,134
*/

type AccountInfo struct {
	Balance   int   `json:"balance"`
	AccountNo int64 `json:"account_no"`
}

type AccountJson map[string]AccountInfo

type Account struct {
	Name    string
	Balance int
}

func main() {
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(stdin), "\n")
	var filtered []string
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if len(l) > 0 {
			filtered = append(filtered, l)
		}
	}

	var accounts []Account
	for _, line := range filtered {
		var accountJson AccountJson
		err := json.Unmarshal([]byte(line), &accountJson)
		if err != nil {
			panic(err)
		}
		var account Account
		for name, info := range accountJson {
			if name == "extra" {
				account.Balance = info.Balance
				continue
			}
			account.Name = name
			if account.Balance == 0 {
				account.Balance = info.Balance
			}
		}
		accounts = append(accounts, account)
	}

	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].Balance < accounts[j].Balance
	})

	for _, acc := range accounts {
		formattedBalance := Format(int64(acc.Balance))
		fmt.Println(acc.Name + ": " + formattedBalance)
	}
}

func Format(n int64) string {
	in := strconv.FormatInt(n, 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}
