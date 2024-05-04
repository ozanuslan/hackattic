package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

/*
Your task is to grab a Redis snapshot from our API and extract some [meta-]data about the contents.

Every time you send a request to the problem endpoint, you'll get back a new, neatly-crafted Redis snapshot with a few fun things inside:

a few totally boring keys with equally boring values
one emoji key - that's right, there's no running from it
one key with an expiry time set
Data will be spread out over a few databases.

To solve the challenge, extract the expiry timestamp of that one key, extract the value of the emoji key (as in, the value the key is set to, not the actual emoji codepoint) and the number of non-empty databases inside the snapshot. Oh, and we'll ask that you also check the type of one of the keys inside and send that back. You'll get the name of the key in the problem JSON.

Minor inconvenience, though - looks like the header may have been... tampered with by a truly demonic, envious entity. Nothing too serious, though.
*/

type Input struct {
	Rdb          string `json:"rdb"`
	Requirements struct {
		CheckTypeOf string `json:"check_type_of"`
	} `json:"requirements"`
}

type Output struct {
	DbCount          int    `json:"db_count"`
	EmojiKeyValue    string `json:"emoji_key_value"`
	ExpiryMillis     int64  `json:"expiry_millis"`
	CheckTypeOf      string
	CheckTypeOfValue string
}

func (o *Output) Marshal() ([]byte, error) {
	interMap := map[string]interface{}{
		"db_count":        o.DbCount,
		"emoji_key_value": o.EmojiKeyValue,
		"expiry_millis":   o.ExpiryMillis,
		o.CheckTypeOf:     o.CheckTypeOfValue,
	}
	json, err := json.Marshal(interMap)
	if err != nil {
		return nil, err
	}
	return json, nil
}

func p(err error) {
	if err != nil {
		panic(err)
	}
}

const REDIS_MAGIC_STRING = "REDIS"

func main() {
	if len(os.Args) < 2 {
		panic(fmt.Errorf("err: must have at least 1 argument"))
	}
	switch os.Args[1] {
	case "decode":
		decodeInput()
	case "read-redis":
		readRedis()
	default:
		panic(fmt.Errorf("err: unknown cmd %s", os.Args[0]))
	}
}

func decodeInput() {
	stdin, err := io.ReadAll(os.Stdin)
	p(err)

	var input Input
	err = json.Unmarshal(stdin, &input)
	p(err)

	decodedRdb, err := base64.StdEncoding.DecodeString(input.Rdb)
	p(err)

	repairedRdb := append([]byte(REDIS_MAGIC_STRING), decodedRdb[5:]...)

	fmt.Print(string(repairedRdb))
}

func readRedis() {
	stdin, err := io.ReadAll(os.Stdin)
	p(err)

	var input Input
	err = json.Unmarshal(stdin, &input)
	p(err)

	client := newRedis(0)

	ctx := context.Background()
	keyspace, err := client.Info(ctx, "keyspace").Result()
	client.Close()
	p(err)
	keyspace = strings.TrimSpace(strings.ReplaceAll(keyspace, "# Keyspace", ""))
	keyspaceLines := strings.Split(keyspace, "\n")

	dbIndexRegexp := regexp.MustCompile(`\d+`)
	dbs := make([]int, 0)
	for _, keyspaceLine := range keyspaceLines {
		match := dbIndexRegexp.FindString(keyspaceLine)
		if match == "" {
			log.Println("keyspace line did not match, skipping:", keyspaceLine)
			continue
		}
		idx, err := strconv.Atoi(match)
		p(err)
		dbs = append(dbs, idx)
	}

	var expiryMillis int64
	var emojiKeyValue string
	var checkTypeOfValue string
	dbCount := len(keyspaceLines)
	log.Println("found db_count:", dbCount)

	emojiRegex := regexp.MustCompile(`[\x{1F600}-\x{1F64F}]|[\x{1F300}-\x{1F5FF}]|[\x{1F680}-\x{1F6FF}]|[\x{1F900}-\x{1F9FF}]|[\x{2600}-\x{26FF}]|[\x{2700}-\x{27BF}]`)
	for _, db := range dbs {
		client := newRedis(db)
		p(err)

		keys, err := client.Keys(ctx, "*").Result()
		p(err)
		log.Println("db:", db, "| keys:", keys)

		for _, key := range keys {
			if emojiKeyValue == "" && emojiRegex.MatchString(key) { // Emoji key
				emojiKeyValue, err = client.Get(ctx, key).Result()
				p(err)
				log.Println("found emoji_key_value key:", key, "| value:", emojiKeyValue)
			} else if checkTypeOfValue == "" && key == input.Requirements.CheckTypeOf {
				checkTypeOfValue, err = client.Type(ctx, key).Result()
				p(err)
				log.Println("found check_type_of key:", key, "| value:", checkTypeOfValue)
			} else if expiryMillis == 0 {
				ttl, err := client.PTTL(ctx, key).Result()
				p(err)
				if ttl.Milliseconds() > 0 {
					expiryMillis = ttl.Milliseconds()
					log.Println("found expiry_millis key: ", key, "| value:", expiryMillis)
				} else {
					log.Println("key:", key, "| ttl:", ttl)
				}
			}
		}

		client.Close()
	}

	if expiryMillis == 0 {
		envTs := os.Getenv("EXPIRY_TIMESTAMP")
		log.Println("expiry_millis not found, so falling back to env var EXPIRY_TIMESTAMP:", envTs)
		expiryMillis, _ = strconv.ParseInt(envTs, 10, 64)
	}

	output := Output{
		DbCount:          dbCount,
		ExpiryMillis:     expiryMillis,
		EmojiKeyValue:    emojiKeyValue,
		CheckTypeOf:      input.Requirements.CheckTypeOf,
		CheckTypeOfValue: checkTypeOfValue,
	}
	out, err := output.Marshal()
	p(err)

	fmt.Println(string(out))
}

func newRedis(db int) *redis.Client {
	host, port, user, password := os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), os.Getenv("REDIS_USER"), os.Getenv("REDIS_PASSWORD")
	return redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Username: user,
		Password: password,
		DB:       db,
	})
}
