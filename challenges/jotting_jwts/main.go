package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

/*
You'll need to write a simple app capable of receiving POST requests, validating JWT tokens and storing some trivial data between requests.

Grab a jwt_secret from the problem endpoint. Configure your app to use it for validating all incoming JWTs. POST your app_url to the solution endpoint.

What happens now is as follows:

our server will send a few requests to your app
they will all be POST requests with a JWT token as body (path will always be /)
the token's payload will contain a key named append set to a string
starting with an empty string, append whatever your receive for every valid request
after some time, you'll receive a token without the append key inside - when this happens, respond with a simple JSON object with the solution key set to whatever you got after appending all the strings received
grab points if constructed string matches what we were expecting
Once you grab a jwt_secret, you have 5 seconds to submit an app_url to the solution endpoint.
*/

type Input struct {
	JwtSecret string `json:"jwt_secret"`
}

type Output struct {
	Solution string `json:"solution"`
}

func p(err error) {
	if err != nil {
		panic(err)
	}
}

var jwtSecret string
var appendedString string

func main() {
	stdin, err := io.ReadAll(os.Stdin)
	p(err)

	var input Input
	err = json.Unmarshal(stdin, &input)
	p(err)
	jwtSecret = input.JwtSecret

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		body, err := io.ReadAll(r.Body)
		p(err)

		log.Println("Token:", string(body))
		token, err := jwt.Parse(string(body), func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil {
			return
		}

		var append string
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if a, ok := claims["append"]; ok {
				append, ok = a.(string)
				if !ok {
					log.Printf("Claims[\"append\"] is not ok: %v\n", a)
				}
			}
		}

		if append != "" {
			appendedString += append
			log.Println("Append:", append)
		} else {
			output := Output{Solution: appendedString}
			out, err := json.Marshal(output)
			p(err)
			log.Println("Result:", string(out))
			io.WriteString(w, string(out))
		}
	})

	log.Println("Starting server...")
	err = http.ListenAndServe(":1337", nil)
	p(err)
}
