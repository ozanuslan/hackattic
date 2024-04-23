package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

/*
Hit the problem endpoint, grab a one-time token. The token expires in 60 minutes, so take your time.

Connect to our WebSocket server by appending your token to wss://hackattic.com/_/ws/$token.

Before we explain the challenge, know this: the server will send you ping! messages in random intervals of 700, 1500, 2000, 2500 or 3000 miliseconds.

So a short while after connecting you'll receive your first ping! message. This is your cue to send back the time after which it was sent. Obviously, due to network latency, the time measured on your end will be slightly different than the intented interval - it's your task to measure time and properly detect which interval was used by the server. Depending on your approach, you may even detect that some messages arrive slightly faster than expected. Something to think about.

Anyway, send back the detected interval (e.g. 700). You have until the next ping! message to do so. If you fail to send back a value before then, the next ping won't be a ping, but a message explaining how sorry the server feels about it, but it has to close the connection. If this happens, you'll have to reconnect and start over.

If you send back the proper interval, the server will confirm your answer immediately with a simple good! message. These messages don't influence when the ping! messages are sent, so ignore those when measuring intervals. Only look at ping! messages.

Keep the conversation up sufficiently long, and the server will reward you with a secret key.

Submit that secret key to the solution endpoint and grab your reward.
*/

type Input struct {
	Token string `json:"token"`
}

type Output struct {
	Secret string `json:"secret"`
}

func p(err error) {
	if err != nil {
		panic(err)
	}
}

const wsBaseUrl = "wss://hackattic.com/_/ws/"

func main() {
	stdin, err := io.ReadAll(os.Stdin)
	p(err)

	var input Input
	err = json.Unmarshal(stdin, &input)
	p(err)

	ws, _, err := websocket.DefaultDialer.Dial(wsBaseUrl+input.Token, nil)
	p(err)
	defer ws.Close()
	lastUpdatedTime := time.Now()
	log.Println("conn! dial successful.")

	var timeDrift int
	var lastInterval int
	var lastElapsed int
	var solution string
	exit := -1
	for exit < 0 {
		msgType, msg, err := ws.ReadMessage()
		p(err)
		readTime := time.Now()
		msgStr := string(msg)

		switch true {
		case strings.Contains(msgStr, "ping!"):
			lastElapsed = int(readTime.Sub(lastUpdatedTime).Milliseconds())
			lastUpdatedTime = readTime
			lastInterval = guessInterval(lastElapsed)
			log.Println("ping! time_elapsed:", lastElapsed, "| interval_guess:", lastInterval)
			err = ws.WriteMessage(1, []byte(strconv.Itoa(lastInterval)))
			p(err)
		case strings.Contains(msgStr, "good!"):
			timeDrift = lastElapsed - lastInterval
			lastUpdatedTime = lastUpdatedTime.Add(time.Duration(-timeDrift) * time.Millisecond)
			log.Println("good! time_drift:", timeDrift, "ms")
		case strings.Contains(msgStr, "congratulations!"):
			log.Println("congratulations! yippie, we did it.")
			msgPrefix := "congratulations! the solution to this challenge is "
			solution = strings.ReplaceAll(msgStr, msgPrefix, "")
			solution = strings.ReplaceAll(solution, "\"", "")
			exit = 0
		case strings.Contains(msgStr, "hello!"):
			continue
		case strings.Contains(msgStr, "ouch!"):
			log.Println("ouch! server is closing connectiong. too bad!")
			exit = 1
		case strings.Contains(msgStr, "expired challenge token"):
			log.Println("whoops! challenge token is expired.")
			exit = 1
		default:
			log.Println("unknown_message! msgType:", msgType, "msg:", msgStr)
			exit = 1
		}
	}

	log.Println("goodbye! shutting down.")
	if exit > 0 {
		os.Exit(exit)
	}

	output := Output{Secret: solution}
	out, err := json.Marshal(output)
	p(err)
	fmt.Println(string(out))
}

func guessInterval(elapsed int) int {
	intervals := []int{700, 1500, 2000, 2500, 3000}
	if elapsed <= intervals[0] {
		return intervals[0]
	}
	for i, _ := range intervals[:len(intervals)-1] {
		if elapsed >= intervals[i] && elapsed < intervals[i+1] {
			return intervals[i]
		}
	}
	return intervals[len(intervals)-1]
}
