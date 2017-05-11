package main

import (
	"log"

	"github.com/xuther/reversi-alpha-beta/board"
)

func main() {
	b := board.InitializeBoard(10)
	b.Board[4][4] = 2
	b.DrawBoard()

	b.PlacePiece(4, 5, board.White)
	b.PlacePiece(4, 3, board.White)
	b.PlacePiece(3, 4, board.White)
	b.PlacePiece(5, 5, board.White)
	b.PlacePiece(6, 6, board.Black)
	b.PlacePiece(3, 3, board.Black)
	b.PlacePiece(5, 3, board.Black)
	b.PlacePiece(3, 5, board.Black)
	b.PlacePiece(7, 7, board.Black)
	b.PlacePiece(8, 8, board.Black)
	b.PlacePiece(9, 9, board.Black)
	b.PlacePiece(2, 2, board.White)

	log.Printf("Black can play: %v", b.CanPlay(board.Black))
	log.Printf("White can play: %v", b.CanPlay(board.White))
	log.Printf("Piece Placed")

	b = board.InitializeBoard(11)
	b.StartGame()
}
