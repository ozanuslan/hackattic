package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

/*
Your task is pretty simple, but can be a launch pad into very interesting territory.

Send a request to the challenge endpoint. You'll receive a JSON with a single attribute - a pretty random string.

You have about 15 seconds to generate two different files - each of them including the string you just received - which when hashed with MD5 produce the same hash.

Send them back, base64-encoded, to the solution endpoint to get your reward.
*/

type Input struct {
	Include string `json:"include"`
}

type Output struct {
	Files []string `json:"files"`
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	stdin, err := io.ReadAll(os.Stdin)
	e(err)

	var input Input
	err = json.Unmarshal(stdin, &input)
	e(err)

	collidingHex1 := "4dc968ff0ee35c209572d4777b721587d36fa7b21bdc56b74a3dc0783e7b9518afbfa200a8284bf36e8e4b55b35f427593d849676da0d1555d8360fb5f07fea2"
	collidingHex2 := "4dc968ff0ee35c209572d4777b721587d36fa7b21bdc56b74a3dc0783e7b9518afbfa202a8284bf36e8e4b55b35f427593d849676da0d1d55d8360fb5f07fea2"

	collidingBytes1, _ := hex.DecodeString(collidingHex1)
	collidingBytes2, _ := hex.DecodeString(collidingHex2)

	file1 := append(collidingBytes1, []byte(input.Include)...)
	file2 := append(collidingBytes2, []byte(input.Include)...)

	encoded1 := base64.StdEncoding.EncodeToString(file1)
	encoded2 := base64.StdEncoding.EncodeToString(file2)

	fmt.Fprintf(os.Stderr, "include: %s | base64 file: %s\n", input.Include, encoded1)

	output := Output{Files: []string{encoded1, encoded2}}
	out, err := json.Marshal(output)
	e(err)

	fmt.Println(string(out))
}
