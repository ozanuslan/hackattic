package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

/*
Password hashing has come a long way.

The task is straightforward. You'll be given a password, some salt (watch out, it comes base64 encoded, because in this case salt - for extra high entropy - is basically just /dev/urandom bytes), and some algorithms-specific parameters.

Your job is to calculate the required SHA256, HMAC-SHA256, PBKDF-SHA256 and finally scrypt.

There's a secret step here, though you won't get points for it and the reward is englightenment itself: realize how each step uses the previous one on the way to the final result.
*/

type Input struct {
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Pbkdf2   struct {
		Hash   string `json:"hash"`
		Rounds int    `json:"rounds"`
	} `json:"pbkdf2"`
	Scrypt struct {
		N       int    `json:"N"`
		P       int    `json:"p"`
		R       int    `json:"r"`
		Buflen  int    `json:"buflen"`
		Control string `json:"_control"`
	} `json:"scrypt"`
}

type Output struct {
	Sha256 string `json:"sha256"`
	Hmac   string `json:"hmac"`
	Pbkdf2 string `json:"pbkdf2"`
	Scrypt string `json:"scrypt"`
}

func p(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	stdin, err := io.ReadAll(os.Stdin)
	p(err)

	var input Input
	err = json.Unmarshal(stdin, &input)
	p(err)

	salt, err := base64.StdEncoding.DecodeString(input.Salt)
	p(err)

	password := []byte(input.Password)
	pbkdf2Iter := input.Pbkdf2.Rounds
	scryptN := input.Scrypt.N
	scryptP := input.Scrypt.P
	scryptR := input.Scrypt.R
	scryptBuflen := input.Scrypt.Buflen

	fmt.Fprintf(os.Stderr, "pass: %s | b64(salt): %s\n", input.Password, input.Salt)

	outSha := sha256.Sum256(password)
	fmt.Fprintf(os.Stderr, "sha: %x\n", outSha)

	hmacHash := hmac.New(sha256.New, salt)
	_, err = hmacHash.Write(password)
	p(err)
	outHmac := hmacHash.Sum(nil)
	fmt.Fprintf(os.Stderr, "hmac: %x\n", outHmac)

	outPbkdf2 := pbkdf2.Key(password, salt, pbkdf2Iter, sha256.Size, sha256.New)
	fmt.Fprintf(os.Stderr, "pbkdf2: %x\n", outPbkdf2)

	outScrypt, err := scrypt.Key(password, salt, scryptN, scryptR, scryptP, scryptBuflen)
	p(err)
	fmt.Fprintf(os.Stderr, "scrypt: %x\n", outScrypt)

	output := Output{
		Sha256: hex.EncodeToString(outSha[:]),
		Hmac:   hex.EncodeToString(outHmac),
		Pbkdf2: hex.EncodeToString(outPbkdf2),
		Scrypt: hex.EncodeToString(outScrypt),
	}
	out, err := json.Marshal(output)
	p(err)

	fmt.Println(string(out))
}
