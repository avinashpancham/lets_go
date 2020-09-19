package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "questions.csv", "csvfile with questions and answers")
	timerTime := flag.Int("time", 30, "timer for questions")
	flag.Parse()

	// Open file
	file, err := os.Open(*csvFilename)

	if err != nil {
		fmt.Printf("Failed to read %s due to: %s\n", *csvFilename, err)
	}

	// Open reader for csv and stdin
	csvReader := csv.NewReader(file)
	stdReader := bufio.NewReader(os.Stdin)

	// Parse all the questions and answers
	lines, err := csvReader.ReadAll()
	if err != nil {
		fmt.Printf("Failed to read file due to: %s\n", err)
	}
	parsedLines := parseLines(lines)

	score := 0
	timer := time.NewTimer(time.Duration(*timerTime) * time.Second)
	for i, p := range parsedLines {
		fmt.Printf("Problem %d: %s\n", i, p.q)

		answerCh := make(chan string)
		go askQuestion(answerCh, stdReader)

		select {
		case <-timer.C:
			fmt.Println("time is over")
			return
		case userAnswer := <-answerCh:
			// Up score
			if strings.TrimSuffix(userAnswer, "\n") == p.a {
				score++
			}
		}
	}

	fmt.Printf("Your score was: %d\n", score)
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: line[1],
		}
	}
	return ret
}

func askQuestion(ch chan string, reader *bufio.Reader) {
	userAnswer, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed to read file due to: %s\n", err)
	}
	ch <- userAnswer

}

type problem struct {
	q string
	a string
}
