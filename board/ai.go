package board

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/xuther/reversi-alpha-beta/ab"
)

//here's where we hook into the AI using the AB package

type ReversiNode struct {
	b             Board
	madeTurn      int
	movex         int //this makes up part of the id
	movey         int //this makes up part of the id
	maximizingFor int
}

func GetNextMove(b Board) []int {
	log.Printf("Getting next turn for %v", b.Turn)

	Root := ReversiNode{
		b:             b,
		maximizingFor: b.Turn,
	}

	// set search depth based on number of available moves and free spaces

	utility, path := ab.Search(Root, 8, -1000, 1000)

	log.Printf("NextTurn: Utility: %v, Path: %v", utility, path)

	vals := strings.Split(path[len(path)-2], ",") //the last one in the path corresponds to the root node
	x, _ := strconv.Atoi(vals[1])
	y, _ := strconv.Atoi(vals[2])
	return []int{x, y}
}

func (r ReversiNode) GetUtility() int {
	score := r.b.CalculateScore()
	return score[r.maximizingFor] - score[(r.maximizingFor%2)+1]
}

func (r ReversiNode) Branch() []ab.Node {
	validMoves := r.b.GetAllPossibleMoves()
	toReturn := []ab.Node{}

	for i := range validMoves {
		newB := r.b.Copy()
		newB.PlacePiece(validMoves[i][0], validMoves[i][1], newB.Turn, false)
		newB.NextTurn()
		toAdd := ReversiNode{
			b:             newB,
			madeTurn:      r.b.Turn,
			movex:         validMoves[i][0],
			movey:         validMoves[i][1],
			maximizingFor: r.maximizingFor,
		}

		toReturn = append(toReturn, toAdd)
	}

	return toReturn
}

func (r ReversiNode) GetMinMax() int {
	if r.maximizingFor == r.b.Turn {
		return 1
	} else {
		return -1
	}
}

func (r ReversiNode) PrintNode() {
	stringsToPrint := make([]string, r.b.Size)
	for i := 0; i < r.b.Size; i++ {
		for j := 0; j < r.b.Size; j++ {
			stringsToPrint[j] += fmt.Sprintf("%v  ", r.b.Board[j][i])
		}
	}
	for _, str := range stringsToPrint {
		log.Printf("")
		log.Printf(str)
	}
}

func (r ReversiNode) GetNodeID() string {
	Turn := ""
	if r.madeTurn == 1 {
		Turn = "B"
	} else {
		Turn = "W"
	}
	return fmt.Sprintf("%v,%v,%v", Turn, r.movex, r.movey)
}
