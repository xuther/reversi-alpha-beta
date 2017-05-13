package board

import (
	"log"
	"testing"

	"github.com/xuther/reversi-alpha-beta/ab"
)

func TestAI(t *testing.T) {
	b := InitializeBoard(4)
	b.Board[1][2] = White
	b.Board[3][1] = White
	b.Board[2][1] = White
	b.Board[2][0] = White
	b.Board[3][0] = Black
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
	b := InitializeBoard(4)
	b.Turn = 1

	Root := ReversiNode{
		b:             b,
		maximizingFor: 1,
	}

	utility, path := ab.Search(Root, 5, -1000, 1000)

	log.Printf("Utility: %v, Path: %v", utility, path)
	log.Printf("Final Board:")

}
func TestAIC(t *testing.T) {
	b := InitializeBoard(8)
	b.Turn = 1

	Root := ReversiNode{
		b:             b,
		maximizingFor: 1,
	}

	utility, path := ab.Search(Root, 9, -1000, 1000)

	log.Printf("Utility: %v, Path: %v", utility, path)
	log.Printf("Final Board:")
}
