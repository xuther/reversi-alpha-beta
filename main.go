package main

import (
	"bufio"
	"log"
	"os"

	tm "github.com/buger/goterm"
	"github.com/xuther/reversi-alpha-beta/ai"
	"github.com/xuther/reversi-alpha-beta/board"
)

func main() {
	f, _ := os.Create("outputFile.txt")
	log.SetOutput(f)
	defer f.Close()

	boardC()
}

func boardC() {
	b := board.InitializeBoard(4)
	tm.Clear()
	b.Turn = 1

	val := ai.GetNextMove(&b)
	log.Printf("Next Turn: %+v", val)

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

}

func boardA() {
	b := board.InitializeBoard(4)
	b.Board[1][2] = board.White
	b.Board[3][1] = board.White
	b.Board[2][1] = board.White
	b.Board[2][0] = board.White
	b.Board[3][0] = board.Black
	tm.Clear()
	b.DrawBoard()
	b.Turn = 1
	b.TurnCount = 6

	val := ai.GetNextMove(&b)
	log.Printf("Next Turn: %+v", val)

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

func boardB() {
	b := board.InitializeBoard(4)
	b.Board[0][3] = board.Black
	b.Board[1][3] = board.White
	b.Board[2][2] = board.White
	tm.Clear()
	b.DrawBoard()
	b.Turn = 1
	b.TurnCount = 6

	val := ai.GetNextMove(&b)
	log.Printf("Next Turn: %+v", val)

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

}
