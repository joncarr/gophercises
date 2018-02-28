package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	winPath = "C:\\Code\\Go\\src\\github.com\\joncarr\\gophercises\\quiz\\"
)

func main() {

	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	_ = timeLimit
	flag.Parse()

	file, err := os.Open(winPath + *csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open CSV file: %s", *csvFilename))
	}
	defer file.Close()

	rdr := csv.NewReader(file)
	lines, err := rdr.ReadAll()
	if err != nil {
		exit("Failed to parse the provied CSV file.")
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem %d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nTIME EXPIRED!!!\nYou answered %d out of %d correctly!", correct, len(problems))
			return
		case answer := <-answerCh:
			if p.a == answer {
				correct++
			}
		}
	}

	fmt.Printf("You answered %d out of %d correctly!", correct, len(problems))

}

type problem struct {
	q string
	a string
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: strings.TrimSpace(line[0]),
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
