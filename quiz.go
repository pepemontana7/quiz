package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
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

/*var timeoutSeconds = 30 * time.Second

var wg sync.WaitGroup
var (
	// sigChan receives os signals.
	sigChan = make(chan os.Signal, 1)

	// timeout limits the amount of time the program has.
	timeout = time.After(timeoutSeconds)

	// complete is used to report processing is done.
	complete = make(chan error)

	// shutdown provides system wide notification.
	shutdown = make(chan struct{})
)*/

func main() {

	//timePtr := flag.Int("timeout", timeoutSeconds, "Default timeout  (Required)")
	quizPtr := flag.String("quiz-file", "problems.csv", "Quiz File(Required)")
	flag.Parse()
	runQuiz(*quizPtr)

}

func runQuiz(f string) {

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
	total := len(lines)
	correct := 0
	// Loop through lines & turn into object
	for _, line := range lines {
		//q := strings.Split(line[0], "+")
		q, op := split(line[0], separators)

		//fmt.Println(q, op)
		data := QA{
			left:     q[0],
			right:    q[1],
			operand:  op,
			question: line[0],
			answer:   line[1],
		}

		//go func(data QA){
		//	fmt.Println("Question: ", data.question, " = ?  ")
		//}
		qs = append(qs, data)
		//fmt.Println(data.left + " " + data.operand + " " + data.right + " = " + data.answer)
		//reader := bufio.NewReader(os.Stdin)
		fmt.Print("Question: ", data.question, "= ? ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		//fmt.Println(scanner.Text())

		if strings.ToLower(strings.TrimSpace(scanner.Text())) == strings.ToLower(data.answer) {
			correct = correct + 1
		}
		//ans, _ := reader.ReadString('\n')
		//fmt.Println(ans, data.answer, ans == data.answer)

	}

	fmt.Println("Correct: ", correct, " Total: ", total)
}
