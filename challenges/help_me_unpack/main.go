package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

/*
The challenge is to receive bytes and extract some numbers from those bytes.

Connect to the problem endpoint, grab a base64-encoded pack of bytes, unpack the required values from it and send them back.

The pack contains, always in the following order:

a regular int (signed), to start off
an unsigned int
a short (signed) to make things interesting
a float because floating point is important
a double as well
another double but this time in big endian (network byte order)

In case you're wondering, we're using 4 byte ints, so everything is in the context of a 32-bit platform.

Extract those numbers from the byte string and send them back to the solution endpoint for your reward. See the solution section for a description of the expected JSON format.

Solution JSON structure:

int: the signed integer value
uint: the unsigned integer value
short: the decoded short value
float: surprisingly, the float value
double: the double value - shockingly
big_endian_double: you get the idea by now!

To make things easier, the response will usually include info about which value you got wrong and what was the expected value.
*/

type Input struct {
	Bytes string `json:"bytes"`
}

type Output struct {
	Int             int32   `json:"int"`
	Uint            uint32  `json:"uint"`
	Short           int16   `json:"short"`
	Float           float32 `json:"float"`
	Double          float64 `json:"double"`
	BigEndianDouble float64 `json:"big_endian_double"`
}

func main() {
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	var input Input
	err = json.Unmarshal(stdin, &input)
	if err != nil {
		panic(err)
	}

	decoded, err := base64.StdEncoding.DecodeString(input.Bytes)
	if err != nil {
		panic(err)
	}

	output, err := unpack(decoded)
	if err != nil {
		panic(err)
	}

	out, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}

func unpack(buffer []byte) (Output, error) {
	// Required buffer length: 32
	// Buffer layout: int32, uint32, int16, float32, float64, float64
	// Buffer       : 4      4       2(2)   4        8        8
	mustLen := 32
	if len(buffer) != mustLen {
		return Output{}, fmt.Errorf("err: buffer must have exactly %d bytes, got %d", mustLen, len(buffer))
	}

	output := Output{}
	buf := bytes.NewReader(buffer[:4])
	err := binary.Read(buf, binary.LittleEndian, &output.Int)
	if err != nil {
		return Output{}, fmt.Errorf("err: failed to read int32: %v", err)
	}
	buf = bytes.NewReader(buffer[4:8])
	err = binary.Read(buf, binary.LittleEndian, &output.Uint)
	if err != nil {
		return Output{}, fmt.Errorf("err: failed to read uint32: %v", err)
	}
	buf = bytes.NewReader(buffer[8:10])
	err = binary.Read(buf, binary.LittleEndian, &output.Short)
	if err != nil {
		return Output{}, fmt.Errorf("err: failed to read int16: %v", err)
	}
	buf = bytes.NewReader(buffer[12:16])
	err = binary.Read(buf, binary.LittleEndian, &output.Float)
	if err != nil {
		return Output{}, fmt.Errorf("err: failed to read float32: %v", err)
	}
	buf = bytes.NewReader(buffer[16:24])
	err = binary.Read(buf, binary.LittleEndian, &output.Double)
	if err != nil {
		return Output{}, fmt.Errorf("err: failed to read float64: %v", err)
	}
	buf = bytes.NewReader(buffer[24:32])
	err = binary.Read(buf, binary.BigEndian, &output.BigEndianDouble)
	if err != nil {
		return Output{}, fmt.Errorf("err: failed to read big endian float64: %v", err)
	}

	return output, nil
}
