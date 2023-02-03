package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	csvFilename := flag.String("csv", "problems.csv", "It is a  csv file for quiz in format of question,answer")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in second")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Printf("Unable to open file: %v\n", *csvFilename)
		os.Exit(1)

	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		fmt.Println("Problem in reading the content of csv file")
		os.Exit(1)
	}

	problems := ParseLines(lines)
	corr := 0

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

loop:
	for i, p := range problems {
		fmt.Printf("Problems #%d: %s = \n", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			break loop

		case answer := <-answerCh:
			if answer == p.a {
				corr++
			}

		}

	}
	fmt.Printf("You have got %d out of %d\n", corr, len(problems))
}

func ParseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i, p := range lines {
		problems[i] = problem{q: p[0], a: strings.TrimSpace(p[1])}
	}

	return problems
}

type problem struct {
	q string
	a string
}
