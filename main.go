package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	tm "github.com/buger/goterm"
	"github.com/xuther/reversi-alpha-beta/board"
)

func main() {
	f, _ := os.Create("outputFile.txt")
	log.SetOutput(f)
	defer f.Close()

	start()
}

func start() {
	tm.Clear()
	tm.Println("Select Game type:")
	tm.Println("0 - Play against another player (hotseat)")
	tm.Println("1 - Play against an AI (AI goes first)")
	tm.Println("2 - Play against an AI (AI goes Second)")
	tm.Println("3 - Watch an AI Game")
	tm.Flush()

	val := 0

	for {
		reader := bufio.NewReader(os.Stdin)
		selection, _ := reader.ReadString('\n')
		matches, err := regexp.MatchString("[0123]", strings.TrimSpace(selection))
		if err != nil {
			tm.Println("Invalid input")
			tm.Flush()
		} else if !matches {
			tm.Println("Invalid input")
			tm.Flush()
		} else if val, err = strconv.Atoi(strings.TrimSpace(selection)); err != nil {
			tm.Println("Invalid input")
			tm.Flush()
		} else {
			break
		}
	}

	b := board.InitializeBoard(8)

	b.StartGame(val)
}
