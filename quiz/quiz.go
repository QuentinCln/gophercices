package main

// https://gophercises.com/exercises/quiz
 
import (
	"encoding/csv"
	"io"
	"os"
	"flag"
	"fmt"
	"time"
)

type line struct {
	question string
	answer string
}

type quiz struct {
	lines []line
	goodAnswer int
}

func result(completed bool, goodAnswer int, totalQuestion int) {
	if completed {
		fmt.Printf("Score: %d / %d\n", goodAnswer, totalQuestion)
	} else {
		fmt.Printf("Out of time !\nScore: %d / %d\n", goodAnswer, totalQuestion)
	}
}

func readLine(c chan string) {
	var answer string
	fmt.Scanf("%s\n", &answer)
	c <- answer
}

func game(quiz quiz, timer time.Timer) (bool, int, int) {
	readAnswerCh := make(chan string)
	for i, line := range quiz.lines {
		fmt.Printf("Q#%d: %s\n", i+1, line.question)
		go readLine(readAnswerCh)
		select {
		case result := <- readAnswerCh:
			if (result == line.answer) {
				quiz.goodAnswer += 1
				fmt.Println("Correct !") 
			} else {
				fmt.Println("Incorrect !")
			}
		case  <- timer.C:
			return false, quiz.goodAnswer, len(quiz.lines)
		}
	}
	return true, quiz.goodAnswer, len(quiz.lines)
}

func getQuiz(filename string) quiz {
	csvFile, err := getFile(filename);
	if (err != nil) {
		fmt.Println("Error on reading file ~")
		os.Exit(1)
	}
	quizLines, err := csv.NewReader(csvFile).ReadAll()
	if (err != nil) {
		fmt.Println("Error on parsing CSV File")
		os.Exit(1)
	}
	
	records := make([]line, len(quizLines))
	for i, l := range quizLines {
		records[i] = line {
			question: l[0],
			answer: l[1],
		}
	}
	quiz := quiz {
		lines: records,
	}
	return quiz
}

func getFile(filename string) (io.Reader, error) {
	return os.Open(filename)
}

func getArguments() (string, time.Timer) {
	fileName := flag.String("filename", "problem.csv", "CSV File that conatins quiz questions")
	timerSecond := flag.Int("timer", 42, "Integer")
	timer := time.NewTimer(time.Duration(*timerSecond) * time.Second)
	flag.Parse()
	return *fileName, *timer
}

func main() {
	filename, timer := getArguments()
	quiz := getQuiz(filename)
	result(game(quiz, timer))
	os.Exit(0)
}