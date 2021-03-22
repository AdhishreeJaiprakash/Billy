package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const TOLERANCE = 1 //dollar tolerance accepted

//global variables
var (
	_PEOPLE map[string]*people
	_ITEMS  map[string]*item
	_TOTAL  float64
)

//each person taking part in the bill split
type people struct {
	name  string
	total float64
	items []string
}

//each item in the bill
type item struct {
	name  string
	price float64
	peeps []string
}

/*
These actions are based on instructions in billy.go. Actions supported:
  1. Quit - Exit now!
  2. NEW_ENTRY - Billy found a new entry
  3. LIST_ENTRY - List all entries recorded with Billy
  4. LIST_PEOPLE - List all people Billy has identified from entries made
  5. DONE - Big brain time!
  6. UNDO - Make changes to an entry
  7. PRINT - Print instructions
*/

type action uint8

const (
	UNKNOWN action = iota
	QUIT
	NEW_ENTRY
	LIST_ENTRY
	LIST_PEOPLE
	DONE
	REMOVE
	PRINT
)

/*
Determine the action to take from given input
*/
func getAction(s string) action {
	switch s {
	case "quit", "q":
		return QUIT
	case "le", "list entry":
		return LIST_ENTRY
	case "lp", "list people":
		return LIST_PEOPLE
	case "done", "d":
		return DONE
	case "remove", "r":
		return REMOVE
	case "print instructions", "p":
		return PRINT
	default:
		if strings.Contains(s, ":") {
			//new entry
			return NEW_ENTRY
		}
	}
	return UNKNOWN
}

//Returns slice of all people names
func getAllPeeps() []string {
	var s []string
	for pName, _ := range _PEOPLE {
		s = append(s, pName)
	}
	return s
}

/*
Update global maps given this new entry
*/
func handleNewEntry(s string) {
	//expected format: <item-name>:<price>:<p1>,<p2>,<p3>
	//extract data
	l := strings.Split(s, ":")
	if len(l) != 3 {
		fmt.Printf("Entry %s is not of expected format. Ignoring this entry.\nEnter 'print instructions'/'p' to see instructions\n", s)
		return
	}

	//get item's name
	if _, ok := _ITEMS[l[0]]; ok {
		fmt.Printf("Duplicate entry of item %s. Ignoring this entry.\n", l[0])
		return
	}

	//get price of item
	price, err := strconv.ParseFloat(l[1], 64)
	if err != nil {
		fmt.Printf("Failed to read item's (%s) price. Err: %v. Retry..\n", l[1])
		return
	}

	//get names of people splitting the item
	peeps := strings.Split(l[2], ",")
	if len(peeps) <= 0 {
		fmt.Printf("Failed to find people for item %s. Ignoring this entry. \nEnter 'print instructions'/'p' to see instructions\n", l[0])
		return
	}

	if strings.ReplaceAll(peeps[0], " ", "") == "all" {
		peeps = getAllPeeps()
	}

	itemPricePerPerson := price / float64(len(peeps))

	for _, pName := range peeps {
		p := _PEOPLE[pName]
		if p == nil { //unrecognized person
			fmt.Printf("Unrecognized person %s in entry. Ignoring this entry. Below are accepted names:\n", pName)
			handleListPeople()
			return
		}
		p.items = append(p.items, l[0])
		p.total = p.total + itemPricePerPerson
		_PEOPLE[pName] = p
	}

	//now update global tables
	_ITEMS[l[0]] = &item{
		name:  l[0],
		price: price,
		peeps: peeps,
	}
}

/*
Print item entries

Format:
<item-name> <item-price> <people-splitting-the-item>
*/
func handleListEntry() {
	i := 0

	//format string
	fmt.Printf("\nList of entries:\n")
	fmt.Printf("   | %-15s | price | split between\n", "Name")
	//draw a line
	dots := 50
	for j := 0; j < dots; j++ {
		fmt.Printf("-")
	}
	fmt.Printf("\n")
	for _, v := range _ITEMS {
		s := ""
		for _, p := range v.peeps {
			s += p + " "
		}
		fmt.Printf("%-2d | %-15s | %-4.2f | %-25s\n", i, v.name, v.price, s)
		i++
	}
	fmt.Printf("\n")
}

/*
Print the names of the people taking part in the bill split
*/
func handleListPeople() {
	s := ""
	i := 0

	fmt.Printf("\nList of people:\n")
	for k, _ := range _PEOPLE {
		s += strconv.Itoa(i) + ". " + k + "   "
		i++
	}
	fmt.Println(s, "\n")
}

/*
Final split calculations and print
*/
func runCalculations() (quit bool) {
	//confirm with the user before proceeding;
	//program will terminate after this
	fmt.Println("Enter 'yes'/'y' if everything looks good:")
	handleListEntry()

	confirmation := readFromStdin()
	if len(confirmation) == 0 {
		return false
	}

	if isQuit(confirmation) {
		return true
	}

	if !strings.Contains(confirmation, "y") {
		fmt.Println("Not proceeding with calculations.")
		return false
	}

	//start with including tax and other common expenses in
	//each person's expenditure
	spent := float64(0)
	for _, i := range _ITEMS {
		spent += i.price
	}
	common := _TOTAL - spent //tax

	//update each person's expense
	var total float64
	for k, v := range _PEOPLE {
		v.total = v.total + (v.total / spent * common)
		_PEOPLE[k] = v
		total += _PEOPLE[k].total
	}

	//make some space on stdout
	fmt.Println("\n")

	//print any disparity in calculation
	if total != _TOTAL {
		//complain only if the difference is greater than permitted tolerance
		if math.Abs(_TOTAL-total) > TOLERANCE {
			fmt.Printf("WARNING!! %4.2f not accounted for.\n", math.Abs(_TOTAL-total))
		}
	}

	//print each person's share
	fmt.Printf("Total amount: %4.2f\n", _TOTAL)
	fmt.Printf("Total tax: %4.2f\n", common)

	//format string
	fmt.Printf("%-10s | Charge | Items\n", "Name")
	//draw a line
	dots := 50
	for j := 0; j < dots; j++ {
		fmt.Printf("-")
	}
	fmt.Printf("\n")
	for _, v := range _PEOPLE {
		items := ""
		separator := ","
		for index, i := range v.items {
			if index == len(v.items)-1 {
				separator = ""
			}
			items += i + separator
		}
		fmt.Printf("%-10s | %-4.2f  | %-25s\n", v.name, v.total, items)
	}

	//terminate program
	return true
}

/*
Edit previously entered entry
*/
func handleDelete() (quit bool) {
	//show available entries
	fmt.Println("\nChoose entry from below list. Enter item name to edit entry.")
	handleListEntry()

	//change entry
	//read from stdin
	item := readFromStdin()
	if isQuit(item) {
		return true
	}

	if _, ok := _ITEMS[item]; !ok {
		fmt.Printf("Item %s not found in entries. Start again\n", item)
		return false
	}

	//before we remove the entry from _ITEMS,
	//update _PEOPLE
	itemEntry := _ITEMS[item]
	deduct := itemEntry.price / float64(len(itemEntry.peeps))
	for _, pName := range itemEntry.peeps {
		personEntry := _PEOPLE[pName]
		personEntry.total = personEntry.total - deduct
		//remove item from item list for this person
		for index, i := range personEntry.items {
			if i == item {
				personEntry.items = append(personEntry.items[:index], personEntry.items[index+1:]...)
				break
			}
		}
		_PEOPLE[pName] = personEntry
	}
	delete(_ITEMS, item)
	return false
}

func handleAction(s string) (quit bool) {
	a := getAction(s)
	switch a {
	case NEW_ENTRY:
		handleNewEntry(s)
	case LIST_ENTRY:
		handleListEntry()
	case LIST_PEOPLE:
		handleListPeople()
	case DONE:
		return runCalculations()
	case REMOVE:
		return handleDelete()
	case QUIT:
		return true
	case PRINT:
		fmt.Println(instString + "\n")
	default:
		fmt.Println("Unknown action! Repeating instructions.\n")
		fmt.Println(instString + "\n")
	}
	return false
}
