package ai

import (
	"log"
	"testing"

	"github.com/xuther/reversi-alpha-beta/ab"
	"github.com/xuther/reversi-alpha-beta/board"
)

func TestAI(t *testing.T) {
	b := board.InitializeBoard(4)
	b.Board[1][2] = board.White
	b.Board[3][1] = board.White
	b.Board[2][1] = board.White
	b.Board[2][0] = board.White
	b.Board[3][0] = board.Black
	b.Turn = 1
	b.TurnCount = 4

	Root := ReversiNode{
		b:             b,
		maximizingFor: 1,
	}

	utility, path := ab.Search(Root, 3, -1000, 1000)

	log.Printf("Utility: %v, Path: %v", utility, path)
}

func TestAIB(t *testing.T) {
	b := board.InitializeBoard(4)
	b.Turn = 1

	Root := ReversiNode{
		b:             b,
		maximizingFor: 1,
	}

	utility, path := ab.Search(Root, 15, -1000, 1000)

	log.Printf("Utility: %v, Path: %v", utility, path)
	log.Printf("Final Board:")

}
func TestAIC(t *testing.T) {
	b := board.InitializeBoard(6)
	b.Turn = 1

	Root := ReversiNode{
		b:             b,
		maximizingFor: 1,
	}

	utility, path := ab.Search(Root, 15, -1000, 1000)

	log.Printf("Utility: %v, Path: %v", utility, path)
	log.Printf("Final Board:")

}
