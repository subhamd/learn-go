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
	quizProblemsFileName := flag.String("csvfn", "problems.csv", "File name of a CSV file in the "+
			"format 'Question,Answer'")
	expiryTime := flag.Int("limit", 30, "Time limit after which the quiz expires if a question is " +
		"not answered")
	flag.Parse()

	file, err := os.Open(*quizProblemsFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *quizProblemsFileName))
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		exit("Failed to read the CSV file.")
	}

	problems := parseLinesFromCSV(lines)
	timer := time.NewTimer(time.Duration(*expiryTime) * time.Second)

	correct := 0
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)
		ansChannel := make(chan string, 1)
		go func(channel chan string) {
			var answer string
			fmt.Scanf("%s\n", &answer)
			channel <- answer
		}(ansChannel)

		select {
		case <- timer.C:
			fmt.Printf("\nYou scored %d out of %d\n", correct, len(problems))
			return
		case answer := <- ansChannel:
			if answer == problem.answer {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d\n", correct, len(problems))
}

func parseLinesFromCSV(lines [][]string) []Problem {
	problems := make([]Problem, len(lines))
	for i, line := range lines {
		problems[i] = Problem{
			question: line[0],
			answer: strings.TrimSpace(line[1]),
		}
	}
	return problems
}

type Problem struct {
	question string
	answer string
}

func exit(msg string)  {
	fmt.Println(msg)
	os.Exit(1)
}



