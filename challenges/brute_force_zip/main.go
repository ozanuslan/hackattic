package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/yeka/zip"
)

/*
Grab the zip_url from the problem endpoint, download the ZIP file. Inside, among other things that you can rummage through, is a file called secret.txt which contains the solution to this challenge. But the ZIP is password protected, and I'm not giving you the password.

The password is between 4-6 characters long, lowercase and numeric. ASCII only.

You'll probably need to brute-force your way to the secret.txt file. Oh, and you have 30 seconds until the problem expires.

Go! Use the force!
*/

type Input struct {
	ZipUrl string `json:"zip_url"`
}

type Output struct {
	Secret string `json:"secret"`
}

const zipFilename = "challenge.inter.zip"
const secretFilename = "secret.txt"

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

	access_token := os.Getenv("ACCESS_TOKEN")
	url := input.ZipUrl + "?access_token=" + access_token
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(zipFilename)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(body)
	if err != nil {
		panic(err)
	}
	f.Close()

	zipRdr, err := zip.OpenReader(zipFilename)
	if err != nil {
		panic(err)
	}
	defer zipRdr.Close()
	var secretFile *zip.File
	for _, file := range zipRdr.File {
		if file.FileInfo().Name() == secretFilename {
			secretFile = file
			break
		}
	}
	if secretFile == nil {
		panic(fmt.Errorf("err: zip archive does not contain the secret file"))
	}

	// fcrackzip --brute-force --charset a1 --length 4-6 --use-unzip challenge.inter.zip
	cmd := exec.Command("fcrackzip", "--brute-force", "--charset", "a1", "--length", "4-6", "--use-unzip", zipFilename)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	cmdOut := buf.String()
	cmdOut = strings.TrimSpace(cmdOut)
	passwordLinePrefix := "PASSWORD FOUND!!!!: pw == "
	zipPass := strings.ReplaceAll(cmdOut, passwordLinePrefix, "")

	secretFile.SetPassword(zipPass)
	secretRdr, err := secretFile.Open()
	if err != nil {
		panic(err)
	}
	defer secretRdr.Close()
	content, err := io.ReadAll(secretRdr)
	if err != nil {
		panic(err)
	}
	secret := string(content)
	secret = strings.TrimSpace(secret)

	output := Output{Secret: secret}
	out, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
