package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	_ "github.com/lib/pq"
)

/*
This is a good place to start your h^ adventure.

Your task is simple. Connect to the problem endpoint and grab a PostgreSQL database dump. It's going to be base64 encoded, because bytes.

We want you to restore that backup on a Postgres server of your choosing and extract the social security numbers of people marked as alive. The others are of no concern to us. Things should be clearer once you get to know the structure of the DB.

Submit the correct list of SSNs to the solution endpoint and grab your reward.
*/

type Input struct {
	Dump string `json:"dump"`
}

type Output struct {
	AliveSSNs []string `json:"alive_ssns"`
}

func main() {
	if len(os.Args) < 2 {
		panic(fmt.Errorf("err: must have at least 1 argument"))
	}
	switch os.Args[1] {
	case "decode":
		decodeInput()
	case "get-ssns":
		getSSNs()
	default:
		panic(fmt.Errorf("err: unknown cmd %s", os.Args[0]))
	}
}

func decodeInput() {
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	var input Input
	err = json.Unmarshal(stdin, &input)
	if err != nil {
		panic(err)
	}

	decoded, err := base64.StdEncoding.DecodeString(input.Dump)
	if err != nil {
		panic(err)
	}

	rdr := bytes.NewReader(decoded)
	gr, err := gzip.NewReader(rdr)
	if err != nil {
		panic(err)
	}
	defer gr.Close()

	buf, err := io.ReadAll(gr)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buf))
}

func getSSNs() {
	user, host := os.Getenv("PGUSER"), os.Getenv("PGHOST")
	connStr := fmt.Sprintf("host=%s user=%s sslmode=disable", host, user)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var ssns []string
	rows, err := db.Query("SELECT ssn FROM public.criminal_records WHERE status = 'alive'")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var ssn string
		if err := rows.Scan(&ssn); err != nil {
			panic(err)
		}
		ssns = append(ssns, ssn)
	}

	output := Output{AliveSSNs: ssns}
	out, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}
