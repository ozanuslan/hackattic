package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"net/http"
	"os"

	"github.com/liyue201/goqr"
)

/*
Connect to the problem endpoint, grab the image with a QR code from the returned image_url. The code contains a hyphen-formatted, numeric code.

Your task is to parse the QR code and submit the resulting code.

That's it, it's practically free points!
*/

type Input struct {
	ImageUrl string `json:"image_url"`
}

type Output struct {
	Code string `json:"code"`
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

	resp, err := http.Get(input.ImageUrl)
	e(err)

	img, _, err := image.Decode(resp.Body)
	e(err)

	qrData, err := goqr.Recognize(img)
	e(err)
	if len(qrData) != 1 {
		panic(fmt.Errorf("err: expected to read 1 qrdata, got %d", len(qrData)))
	}

	output := Output{Code: string(qrData[0].Payload)}
	out, err := json.Marshal(output)
	e(err)

	fmt.Println(string(out))
}
