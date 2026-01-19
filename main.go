package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

type Question struct {
	Question string
	Answer   string
}

type App struct {
	CSV       string
	TimeLimit int
}

func NewApp() *App {
	app := &App{}
	pflag.StringVar(&app.CSV, "csv", "problems.csv", "a csv file in the format of 'question,answer'")
	pflag.IntVar(&app.TimeLimit, "time", 30, "time limit in seconds")
	pflag.Usage = func() {
		fmt.Println("Usage: quiz [options]")
		fmt.Println("Options:")
		fmt.Println("  --csv <file>     CSV file path (default: problems.csv)")
		fmt.Println("  --time <sec>     Time limit in seconds (default: 30)")
		fmt.Println("  -h, --help       Show this help message")
	}

	pflag.Parse()

	return app
}

func (a *App) Run() {
	lines, err := ReadCSV(a.CSV)
	if err != nil {
		ExitWithMessage(err.Error())
	}
	questions, err := ParseLines(lines)
	if err != nil {
		ExitWithMessage(err.Error())
	}
	timer := time.NewTimer(time.Duration(a.TimeLimit) * time.Second)
	answerCh := make(chan string)

	var correctCount int = 0
	for i, question := range questions {
		fmt.Printf("Problem #%d: %s = ", i+1, question.Question)
		go func() {
			var answer string
			fmt.Scanln(&answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("You scored %d out of %d", correctCount, len(questions))
			return
		case answer := <-answerCh:
			if answer == question.Answer {
				correctCount++
			}
		}
	}
	fmt.Printf("You scored %d out of %d", correctCount, len(questions))
}

func ReadCSV(Path string) (lines [][]string, err error) {
	file, err := os.Open(Path)
	if err != nil {
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	lines, err = reader.ReadAll()
	if err != nil {
		return
	}
	return
}

func ParseLines(Lines [][]string) (questions []Question, err error) {
	questions = make([]Question, len(Lines))
	for i, line := range Lines {
		if len(line) != 2 {
			err = fmt.Errorf("invalid line: %v", line)
			return
		}
		questions[i] = Question{
			Question: strings.TrimSpace(line[0]),
			Answer:   strings.TrimSpace(line[1]),
		}
	}
	return
}

func ExitWithMessage(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func main() {
	app := NewApp()
	app.Run()
}
