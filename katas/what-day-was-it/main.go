package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
For every number you read from STDIN, add that many days to the epoch date of January 1st, 1970 and print out the name of the day of week on that date.

Sample input
100
0
128
2544

Sample output
Saturday
Thursday
Saturday
Sunday
*/

func main() {
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	split := strings.Split(string(stdin), "\n")
	days := make([]int, 0)
	for _, dayStr := range split {
		dayStr = strings.TrimSpace(dayStr)
		if len(dayStr) < 1 {
			continue
		}
		dayNum, err := strconv.Atoi(dayStr)
		if err != nil {
			panic(err)
		}
		days = append(days, dayNum)
	}

	for _, dayCount := range days {
		fmt.Println(findWeekDaySinceEpoch(dayCount))
	}
}

func findWeekDaySinceEpoch(days int) string {
	epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	targetTime := epoch.AddDate(0, 0, days)
	return targetTime.Weekday().String()
}
