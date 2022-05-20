package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type flashcard struct {
	term       string
	definition string
	hardness   int
}

type deck []flashcard

var deckOfFlashcards = deck{}
var log []string

func checkTerm(term string) bool {
	for _, card := range deckOfFlashcards {
		if card.term == term {
			return true
		}
	}
	return false
}

func checkDefinition(definition string) bool {
	for _, card := range deckOfFlashcards {
		if card.definition == definition {
			return true
		}
	}
	return false
}

func lineScanner() string {
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		os.Exit(1)
	}
	// remove the delimiter from the string
	input = strings.TrimSuffix(input, "\n")
	log = append(log, input)
	return input
}

var termRefill = false
var definitionRefill = false
var thisCardNum = 1

func printAndLog(str string) {
	fmt.Printf(str)
	log = append(log, str)
}

func populateTerm() string {
	if !termRefill {
		fmt.Printf("The card:\n")
	}
	term := lineScanner()
	if checkTerm(term) {
		printAndLog(fmt.Sprintf("The card \"%s\" already exists. Try again:\n", term))
		termRefill = true
		return populateTerm()
	} else {
		return term
	}
}

func populateDefinition() string {
	if !definitionRefill {
		fmt.Printf("The definition of the card:\n")
	}
	definition := lineScanner()
	if checkDefinition(definition) {
		printAndLog(fmt.Sprintf("The definition \"%s\" already exists. Try again:\n", definition))
		definitionRefill = true
		return populateDefinition()
	} else {
		thisCardNum++
		return definition
	}
}

func populateCard() flashcard {
	var card flashcard
	card.term = populateTerm()
	card.definition = populateDefinition()
	return card
}

func defInOtherCard(def string) bool {
	for _, card := range deckOfFlashcards {
		if card.definition == def {
			return true
		}
	}
	return false
}

func getCardTerm(def string) string {
	for _, card := range deckOfFlashcards {
		if card.definition == def {
			return card.term
		}
	}
	return ""
}

func addCard() {
	card := populateCard()
	deckOfFlashcards = append(deckOfFlashcards, card)
	printAndLog(fmt.Sprintf("The pair (\"%s\":\"%s\") has been added.\n", card.term, card.definition))
}

func removeCard() {
	printAndLog(fmt.Sprintf("Which card?\n"))
	term := lineScanner()
	for i, card := range deckOfFlashcards {
		if card.term == term {
			deckOfFlashcards = append(deckOfFlashcards[:i], deckOfFlashcards[i+1:]...)
			printAndLog(fmt.Sprintf("The card has been removed.\n"))
			return
		}
	}
	printAndLog(fmt.Sprintf("Can't remove \"%s\": there is no such card.\n", term))
}

func importDeck(fileName string) {
	if fileName == "" {
		fmt.Println("File name:")
		fileName = lineScanner()
	}
	file, err := os.Open(fileName)
	if err != nil {
		printAndLog(fmt.Sprintf("File not found.\n"))
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	scanner := bufio.NewScanner(file)
	cardsImported := 0
	for scanner.Scan() {
		itLine := scanner.Text()
		line := strings.Split(itLine, ":")
		hardness := 0
		if len(line) == 3 {
			hardness, err = strconv.Atoi(line[2])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		card := flashcard{line[0], line[1], hardness}
		if !checkTerm(card.term) { // if the term is not in the deck, add it
			deckOfFlashcards = append(deckOfFlashcards, card)
		} else { // if the term is in the deck, update the definition
			for _, thisCard := range deckOfFlashcards {
				if thisCard.term == card.term {
					thisCard.definition = card.definition
				}
			}
		}
		cardsImported++
	}
	printAndLog(fmt.Sprintf("%d cards have been loaded.\n", cardsImported))
}

func exportDeck(fileName string) {
	if fileName == "" {
		fmt.Println("File name:")
		fileName = lineScanner()
	}
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	for _, card := range deckOfFlashcards {
		_, err := file.WriteString(card.term + ":" + card.definition + "\n")
		if err != nil {
			fmt.Println(err)
		}
	}
	printAndLog(fmt.Sprintf("%v cards have been saved.\n", len(deckOfFlashcards)))
}

func ask() {
	printAndLog(fmt.Sprintf("How many times to ask?\n"))
	timesToAsk, err := strconv.Atoi(lineScanner())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	askQuiz(timesToAsk)
}

func updateHardness(term string) {
	for i, thisCard := range deckOfFlashcards {
		if thisCard.term == term {
			deckOfFlashcards[i].hardness = deckOfFlashcards[i].hardness + 1
		}
	}
}

func askQuiz(num int) {
	// ask this number of questions from the quiz deck, if the index is out of range, loop over the deck again
	quizDeckIndex := 0
	for i := 0; i < num; i++ {
		if len(deckOfFlashcards) == 0 {
			return
		} else {
			quizDeckIndex = i % len(deckOfFlashcards)
		}
		card := deckOfFlashcards[quizDeckIndex]
		printAndLog(fmt.Sprintf("Print the definition of \"%s\":\n", card.term))
		userInput := lineScanner()
		if userInput == card.definition {
			printAndLog(fmt.Sprintf("Correct!\n"))
		} else if defInOtherCard(userInput) {
			// update hardness of the card in deckOfFlashcards
			updateHardness(card.term)
			// card.hardness += 1
			anotherTerm := getCardTerm(userInput)
			printAndLog(fmt.Sprintf("Wrong. The right answer is \"%s\", but your definition is correct for \"%s\".\n", card.definition, anotherTerm))
		} else {
			updateHardness(card.term)
			// card.hardness += 1
			printAndLog(fmt.Sprintf("Wrong. The right answer is \"%s\".\n", card.definition))
		}
		quizDeckIndex++
	}
}

func maxHardness() int {
	max := 0
	for _, card := range deckOfFlashcards {
		if card.hardness > max {
			max = card.hardness
		}
	}
	return max
}

func hardest() {
	// for the deckOfFlashcards, iterate through and make a slice of cards with the highest hardness
	var hardestCards []string
	max := maxHardness()
	if max != 0 {
		for _, card := range deckOfFlashcards {
			if card.hardness == max {
				hardestCards = append(hardestCards, "\""+card.term+"\"")
			}
		}
	}
	if len(hardestCards) == 0 {
		printAndLog(fmt.Sprintf("There are no cards with errors.\n"))
		return
	}
	if len(hardestCards) == 1 {
		printAndLog(fmt.Sprintf("The hardest card is %s.  You have %d errors answering it.\n", hardestCards[0], max))
		return
	}
	s := strings.Join(hardestCards, ", ")
	printAndLog(fmt.Sprintf("The hardest cards are: %s. You have %d errors answering them.", s, max))
}

func resetHardness() {
	for i := range deckOfFlashcards {
		deckOfFlashcards[i].hardness = 0
	}
	printAndLog(fmt.Sprintf("Card statistics have been reset.\n"))
}

func saveLogToFile() {
	fmt.Println("File name:")
	fileName := lineScanner()
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	for _, line := range log {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			fmt.Println(err)
		}
	}
	printAndLog(fmt.Sprintf("The log has been saved.\n"))
}

var fileNameExport = ""
var fileName = ""

func checkArgs(activity string) {
	if strings.HasPrefix(activity, "--import_from=") {
		fileName = strings.TrimPrefix(activity, "--import_from=")
		importDeck(fileName)
	}
	if strings.HasPrefix(activity, "--export_to=") {
		fileNameExport = strings.TrimPrefix(activity, "--export_to=")
	}
}

func main() {
	if len(os.Args) == 2 {
		checkArgs(os.Args[1])
	}
	if len(os.Args) == 3 {
		checkArgs(os.Args[1])
		checkArgs(os.Args[2])
	}
	for {
		printAndLog(fmt.Sprintf("Input the action (add, remove, import, export, ask, exit, log, hardest card, reset stats):\n"))
		activity := lineScanner()
		checkArgs(activity)
		switch activity {
		case "exit":
			if fileNameExport != "" {
				exportDeck(fileNameExport)
			} else {
				printAndLog(fmt.Sprintf("bye bye\n"))
			}
			os.Exit(0)
		case "add":
			addCard()
		case "remove":
			removeCard()
		case "import":
			importDeck("")
		case "export":
			exportDeck("")
		case "ask":
			ask()
		case "hardest card":
			hardest()
		case "reset stats":
			resetHardness()
		case "log":
			saveLogToFile()
		}
	}
}
