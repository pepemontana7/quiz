package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type QA struct {
	left     string
	operand  string
	right    string
	answer   string
	question string
}

func split(s string, separators []rune) ([]string, string) {
	var rs rune
	f := func(r rune) bool {
		for _, s := range separators {
			if r == s {
				rs = r
				return true
			}
		}
		return false
	}

	return strings.FieldsFunc(s, f), string(rs)

}

var total int
var correct int

func main() {

	timePtr := flag.Int("timeout", 30, "Default timeout  (Required)")
	shufPtr := flag.Bool("shuffle", false, "Shuffle the quiz q's  (Required)")
	quizPtr := flag.String("quiz-file", "problems.csv", "Quiz File(Required)")
	flag.Parse()
	fmt.Println("===== Go Quiz =====")
	fmt.Print("Press Enter Key to start... ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	var timeoutSeconds = time.Duration(*timePtr) * time.Second
	var (
		// timeout limits the amount of time the program has.
		timeout = time.After(timeoutSeconds)
		// complete is used to report processing is done.
		complete = make(chan error)
	)

	go runQuiz(*quizPtr, complete, *shufPtr)
ControlLoop:
	for {
		select {

		case <-timeout:

			// We have taken too much time. Kill the app hard.
			fmt.Println("Timeout - Killing Program")
			fmt.Println("Correct: ", correct, " Total: ", total)

			os.Exit(1)

		case <-complete:

			// Everything completed within the time given.
			//fmt.Printf("Task Completed: Error[%s]", err)
			break ControlLoop
		}
	}

	// Program finished.
	fmt.Println("Quiz Ended")
}

func runQuiz(f string, complete chan<- error, shuf bool) {

	var err error

	// Defer the send on the channel so it happens
	// regardless of how this function terminates.
	defer func() {

		// Signal the goroutine we have shutdown.
		complete <- err
	}()

	problems, err := os.Open(f)
	if err != nil {

		fmt.Println("Error:", err)
		return
	}
	defer problems.Close()

	var qs []QA

	// Read File into a Variable
	lines, err := csv.NewReader(problems).ReadAll()
	if err != nil {
		panic(err)
	}

	separators := []rune{'+', '-', '\\', '*'}
	total = len(lines)
	correct = 0
	var list []int
	if shuf {
		list = rand.Perm(total)
	}
	// Loop through lines & turn into object
	for i, _ := range lines {
		//q := strings.Split(line[0], "+")
		var l []string

		if shuf {
			l = lines[list[i]]

		} else {
			l = lines[i]
		}
		q, op := split(l[0], separators)

		//fmt.Println(q, op)
		data := QA{
			left:     q[0],
			right:    q[1],
			operand:  op,
			question: l[0],
			answer:   l[1],
		}

		qs = append(qs, data)

		fmt.Print("Question ", i+1, ": ", data.question, "= ? ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		if strings.ToLower(strings.TrimSpace(scanner.Text())) == strings.ToLower(data.answer) {
			correct = correct + 1
		}
		//ans, _ := reader.ReadString('\n')
		//fmt.Println(ans, data.answer, ans == data.answer)

	}

	fmt.Println("Correct: ", correct, " Total: ", total)
}
