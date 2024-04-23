package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Hallicopter/go-dtmf/dtmf"
)

/*
And now for something old-school cool...

You'll need to download a small .wav file with a DTMF-encoded sequence inside. Play it back, you'll hear it. The task is simple: decode the sequence and send it back as the solution.

The sequence only uses the 0-9 digits as well as * and #.

Good luck!
*/

type Input struct {
	WavUrl string `json:"wav_url"`
}

type Output struct {
	Sequence string `json:"sequence"`
}

func p(err error) {
	if err != nil {
		panic(err)
	}
}

const sampleRate = 4000
const wiggleRoom = 2

// PS: This solutions seems flaky. The underlygin DTMF lib is not working 100% stable. Re-run if need be.

func main() {
	stdin, err := io.ReadAll(os.Stdin)
	p(err)

	var input Input
	err = json.Unmarshal(stdin, &input)
	p(err)

	resp, err := http.Get(input.WavUrl)
	p(err)

	wavBytes, err := io.ReadAll(resp.Body)
	p(err)

	solution, err := dtmf.DecodeDTMFFromBytes(wavBytes, sampleRate, wiggleRoom)
	p(err)

	output := Output{Sequence: solution}
	out, err := json.Marshal(output)
	p(err)

	fmt.Println(string(out))
}
