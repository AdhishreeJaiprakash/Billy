package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	instString = `This tool will help you split your bills. Follow the instructions below to use Billy:
1. Start with entering the total bill amount.
2. Then enter the names of the people participating in the bill in the format: <name1>,<name2>,<name3>...
3. Proceed to enter each entry following the format:
   <iterm-name>:<price>:<person1>,<person2>,...
4. If you'd like to delete any entry, enter 'remove'/'r'.
5. If you'd like to list all people, enter 'list people'/'lp'.
6. If you'd like to list all the entries Billy has recorded, enter 'list entry'/'le'.
7. When you'd like Billy to work its magic, enter 'done'/'d'.
8. To view instructions again, enter 'print instructions'/'p'.
9. To quit, enter 'quit'/'q'.

Note: Common expenses like tax will be scaled according to each person's expenditure.`
)

//Returns string of stdin
func readFromStdin() string {
	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read from stdin.")
		return ""
	}
	return strings.TrimSuffix(str, "\n") //remove newline char
}

//Checks if input is request to quit billy
func isQuit(s string) bool {
	if s == "quit" || s == "q" {
		return true
	}
	return false
}

//init func
func initialize() {
	_PEOPLE = make(map[string]*people)
	_ITEMS = make(map[string]*item)
}

//destroyer
func destroy() {
	_PEOPLE = nil
	_ITEMS = nil
}

//Returs map of people participating in this billy instance
func setPeople(p string) {
	peeps := strings.Split(p, ",")
	if len(peeps)==0 {
		//input not in expected format
		fmt.Println("Names not in expected format. Try again...Want: <name1>,<name2>,..")
		return 
	}
	for _, name := range peeps {
		name = strings.ReplaceAll(name," ","")
		_PEOPLE[name]=&people{
			name: name,
		}	
	}
}

func main() {
	//print welcome message and defer last message
	fmt.Println("\nWelcome to Billy!\n")
	defer func() {
		fmt.Println("\nGoodbye!\n")
		destroy()
	}()

	//print instructions
	fmt.Println(instString + "\n")

	//initialize global variables
	initialize()

	//ask for some basic info
	var total float64
	var err error
	for total <= 0 {
		fmt.Printf("Enter total: ")
		totalStr := readFromStdin()

		//proceed if not quit
		if isQuit(totalStr) {
			return
		}

		//some data validity
		if len(totalStr) == 0 {
			fmt.Println("Missing total. Cannot proceed. Retry..\n")
			continue
		}

		total, err = strconv.ParseFloat(totalStr, 64)
		if err != nil {
			fmt.Printf("Failed to read total. Need a number. Data provided: %s. Retry..\n\n", totalStr)
			continue
		}
		if total <= 0 {
			fmt.Printf("Invalid number entered. Want: >0, Got: %d\n\n", total)
			continue
		}
	}

	_TOTAL = total

	//request for details
	for len(_PEOPLE) <= 0 {
		fmt.Printf("Enter people: ")
		people := readFromStdin()
		if isQuit(people) {
			return
		}
		setPeople(people)
	}

	//read until quit
	fmt.Println("Proceed to make entries.\n")
	var input string
	input = readFromStdin()
	for !isQuit(input) {
		quit := handleAction(input)
		if quit {
			return
		}

		//read more
		input = readFromStdin()
	}
}
