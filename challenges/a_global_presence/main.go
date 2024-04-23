package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

/*
Connect to the problem endpoint, write down your presence_token.

You now have 30 seconds to perform at least 7 requests to the URL https://hackattic.com/_/presence/$presence_token.

Here's the catch: every request has to come from a different country.

For every request, the response will return a list of countries that have checked in so far, in the form of a simple string, like so: PL,NL,DE,CA,RU. After you've accumulated at least 7, send an empty JSON to the solution endpoint to mark the challenge as solved.

Go!
*/

type Input struct {
	PresenceToken string `json:"presence_token"`
}

type Output struct {
}

func p(err error) {
	if err != nil {
		panic(err)
	}
}

const baseUrl = "https://hackattic.com/_/presence/"
const requiredCountries = 7

var requestUrl string
var countryCount = 0

func main() {
	start := time.Now()

	stdin, err := io.ReadAll(os.Stdin)
	p(err)

	var input Input
	err = json.Unmarshal(stdin, &input)
	p(err)

	requestUrl = baseUrl + input.PresenceToken

	proxyFileContent, err := os.ReadFile("proxies.txt")
	if err == os.ErrNotExist {
		fmt.Fprintln(os.Stderr, "proxies.txt file not found")
		fmt.Fprintln(os.Stderr, "please check the proxies.txt.example for the details")
		os.Exit(1)
	}
	p(err)

	proxies := strings.Split(string(proxyFileContent), "\n")
	rand.Shuffle(len(proxies), func(i, j int) { proxies[i], proxies[j] = proxies[j], proxies[i] })

	proxyCh := make(chan string)
	updateCh := make(chan int)
	done := make(chan interface{})

	workerCount := 100
	for i := 0; i < workerCount; i++ {
		go func() {
			workRequestor(i, proxyCh, updateCh, done)
		}()
	}

	go workCountryUpdator(updateCh, done)

	go func() {
		isDone := false
		for _, proxy := range proxies {
			select {
			case <-done:
				isDone = true
				break
			default:
				proxyCh <- proxy
			}
		}
		close(proxyCh)
		if !isDone {
			done <- nil
		}
	}()

	<-done
	close(done)

	if countryCount < requiredCountries {
		log.Fatalln("proxies exhausted but could not reach required country count")
	}
	log.Println("finished! took:", time.Now().Sub(start).String())

	var output Output
	out, err := json.Marshal(output)
	p(err)
	fmt.Println(string(out))
}

func workRequestor(id int, proxyCh <-chan string, countryUpdateCh chan<- int, done <-chan interface{}) {
	for {
		select {
		case proxy := <-proxyCh:
			start := time.Now()
			body, err := sendRequest(proxy)
			elapsed := time.Now().Sub(start)
			if err != nil {
				log.Println("worker:", id, "proxy:", proxy, "took:", elapsed.String(), "err:", err)
				continue
			}
			log.Println("worker:", id, "proxy:", proxy, "took:", elapsed.String(), "body:", body)

			countries := strings.Split(body, ",")
			countryUpdateCh <- len(countries)
		case <-done:
			return
		}
	}
}

func sendRequest(proxy string) (string, error) {
	proxyUrl := &url.URL{Scheme: "http", Host: proxy}
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	client.Timeout = 5 * time.Second

	resp, err := client.Get(requestUrl)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func workCountryUpdator(countryUpdateCh <-chan int, done chan interface{}) {
	for {
		select {
		case <-done:
			return
		case newCount := <-countryUpdateCh:
			if newCount > countryCount {
				countryCount = newCount
				log.Println("country_update:", countryCount)
			}
			if countryCount >= requiredCountries {
				done <- nil
				break
			}
		}
	}
}
