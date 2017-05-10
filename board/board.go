package board

import (
	"fmt"
	"strconv"

	tm "github.com/buger/goterm"
)

type Board struct {
	Size  int
	Board [][]int
	Turn  bool
}

const White = 1
const Black = 2
const Clear = 0

const N = 0
const NE = 1
const E = 2
const SE = 3
const S = 4
const SW = 5
const W = 6
const NW = 7

func InitializeBoard(size int) Board {
	b := make([][]int, size)
	for i := range b {
		b[i] = make([]int, size)
	}
	return Board{
		Size:  size,
		Board: b,
		Turn:  false,
	}
}

//PlacePice places a piece in that 0-indexed location on the board
//
func (b *Board) PlacePiece(x int, y int, player int) {

	//make sure that it's a valid player
	if player != White && player != Black {
		return
	}
	//we need to check in lines from the piece being played in all 8 directions until we hit either
	// a) a clear space
	// b) a wall
	// c) anoter pice by same player
}

//the function that will recursively check the next piece in the line

func (b *Board) checkNext(dir int, curX int, curY int, turn int) bool {
	if turn != White && turn != Black {
		return false
	}

	nextY := curY
	nextX := curX

	switch dir {
	case N:
		nextY--
	case NE:
		nextY--
		nextX++
	case E:
		nextX++
	case SE:
		nextY++
		nextX++
	case S:
		nextY++
	case SW:
		nextY++
		nextX--
	case W:
		nextX--
	case NW:
		nextX--
		nextY--
	default:
		return false
	}

	//we're off the board
	if nextX >= b.Size || nextX < 0 || nextY >= b.Size || nextY < 0 {
		return false
	}

	//we hit a clear space
	if b.Board[nextY][nextX] == Clear {
		return false
	}
	if b.Board[nextY][nextX] == turn {
		return true
	}

	val := b.checkNext(dir, nextX, nextY, turn)

	if val {
		b.Board[nextY][nextX] = turn //capture the pieces
		return true
	}

	return false //nothing happened
}

func (b *Board) DrawBoard() {
	tm.Clear()
	tm.MoveCursor(1, 1)
	vals := make([]interface{}, b.Size+1)

	str := ""
	for i := 0; i < b.Size; i++ {
		str += "%v\t"
		vals[i] = i - 1
	}

	str += "%v\n"
	vals[0] = " "
	vals[b.Size] = b.Size

	board := tm.NewTable(0, 1, 1, ' ', 0)

	fmt.Fprintf(board, str, vals...)
	for i := 0; i < b.Size; i++ {
		fmt.Fprintf(board, b.genBoardStrings(i))
	}

	tm.Println(board)
	tm.Flush()
}

func (b *Board) genBoardStrings(row int) string {
	toReturn := ""
	toReturn += strconv.Itoa(row)

	for i := 0; i < b.Size; i++ {
		cur := ""
		switch b.Board[row][i] {
		case Clear:
			cur = "\t-"
		case White:
			cur = "\tW"
		case Black:
			cur = "\tB"
		}
		toReturn += cur
	}
	toReturn += "\n"
	fmt.Printf("%v", toReturn)
	return toReturn
}
