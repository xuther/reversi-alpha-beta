package main

import "github.com/xuther/reversi-alpha-beta/board"

func main() {
	b := board.InitializeBoard(10)
	b.Board[4][4] = 2
	b.DrawBoard()
}
