package board

import (
	"log"
	"testing"
	"time"

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

	start := time.Now()
	utility, path := ab.Search(Root, 3, -1000, 1000)
	nonParallel := time.Since(start)

	start = time.Now()
	putility, ppath := ab.StartSearchMulti(Root, 3, -1000, 1000)
	parallel := time.Since(start)

	log.Printf("Non Parallel Time: %s", nonParallel)
	log.Printf("Utility: %v, Path: %v", utility, path)
	log.Printf("Parallel Time: %s", parallel)
	log.Printf("Utility: %v, Path: %v", putility, ppath)
}

func TestAIB(t *testing.T) {
	b := InitializeBoard(4)
	b.Turn = 1

	Root := ReversiNode{
		b:             b,
		maximizingFor: 1,
	}

	start := time.Now()
	utility, path := ab.Search(Root, 5, -1000, 1000)
	nonParallel := time.Since(start)

	start = time.Now()
	putility, ppath := ab.StartSearchMulti(Root, 5, -1000, 1000)
	parallel := time.Since(start)

	log.Printf("Utility: %v, Path: %v", utility, path)
	log.Printf("Non Parallel Time: %s", nonParallel)
	log.Printf("Parallel Time: %s", parallel)
	log.Printf("Utility: %v, Path: %v", putility, ppath)

}
func TestAIC(t *testing.T) {
	b := InitializeBoard(8)
	b.Turn = 1

	Root := ReversiNode{
		b:             b,
		maximizingFor: 1,
	}

	start := time.Now()
	utility, path := ab.Search(Root, 9, -1000, 1000)
	nonParallel := time.Since(start)

	start = time.Now()
	putility, ppath := ab.StartSearchMulti(Root, 9, -1000, 1000)
	parallel := time.Since(start)
	time.Sleep(3 * time.Second)

	log.Printf("Utility: %v, Path: %v", utility, path)
	log.Printf("Non Parallel Time: %s", nonParallel)
	log.Printf("Parallel Time: %s", parallel)
	log.Printf("Utility: %v, Path: %v", putility, ppath)
}
