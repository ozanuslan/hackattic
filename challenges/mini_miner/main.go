package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

/*
With the Bitcoin thing going strong, I figured it would be interesting to do some simplified mining.

Connect to the problem endpoint. You'll receive a JSON with two attributes. One is block, which is in essence an object with a nonce value (initially empty) and a data key which contains some arbitrary data. The other attribute is difficulty - we'll get back to it in a moment.

Your goal is find a nonce value that will cause the SHA256 hash of the block object to begin with difficulty zero bits. E.g. a difficulty of 14 means that the SHA256 digest needs to start with at least 14 zero bits.

The hash should be calculated from a JSON-serialized block value without any whitespace. The keys needs to be in alphabetical order.

Let's illustrate this on a really simple case. For a block with an empty data array and a given difficulty of 8 (so the first 8 bits of the SHA256 hash need to be all 0), a nonce value of 45 is one perfectly valid solution:

SHA256('{"data":[],"nonce":45}') -> 00d696db4...cfb19ec2e0141

Keep in mind, difficulty is the number of 0 bits, not bytes.
*/

type Input struct {
	Difficulty int `json:"difficulty"`
	Block      struct {
		Nonce int           `json:"nonce"`
		Data  []interface{} `json:"data"`
	} `json:"block"`
}

type Output struct {
	Nonce int `json:"nonce"`
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}

var hashTemplate string
var zeroBytes int
var zeroBits int

func main() {
	stdin, err := io.ReadAll(os.Stdin)
	e(err)

	var input Input
	err = json.Unmarshal(stdin, &input)
	e(err)

	zeroBytes = input.Difficulty / 8
	zeroBits = input.Difficulty % 8

	data, err := json.Marshal(input.Block.Data)
	e(err)

	hashTemplate = "{\"data\":" + string(data) + ",\"nonce\":%d}"

	inCh := make(chan int)
	outCh := make(chan int)
	done := make(chan interface{})

	workerCount := runtime.NumCPU() * 4

	var wg sync.WaitGroup
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func() {
			work(inCh, outCh, done)
			wg.Done()
		}()
	}

	fmt.Fprintf(os.Stderr, "difficulty: %d | ", input.Difficulty)

	go func() {
		for i := 0; i <= (1<<31 - 1); i++ {
			select {
			case <-done:
				close(inCh)
				return
			default:
				inCh <- i
			}
		}
		close(inCh)
	}()

	nonce := <-outCh
	output := Output{Nonce: nonce}
	out, err := json.Marshal(output)
	e(err)

	fmt.Println(string(out))
	wg.Wait()
}

func work(in <-chan int, out chan<- int, done chan<- interface{}) {
	for nonce := range in {
		s, found := testNonce(nonce)
		if found {
			fmt.Fprintf(os.Stderr, "nonce: %d | sha: %s \n", nonce, s)
			out <- nonce
			close(done)
			return
		}
	}
}

func testNonce(n int) (string, bool) {
	digest := fmt.Sprintf(hashTemplate, n)
	sum := sha256.Sum256([]byte(digest))
	var i int
	for i = 0; i < zeroBytes; i++ {
		if sum[i]|0 != 0 {
			return "", false
		}
	}
	for targetBit := 0; targetBit < zeroBits; targetBit++ {
		if sum[i]&(1<<(7-targetBit)) != 0 {
			return "", false
		}
	}
	return fmt.Sprintf("%x", sum), true
}
