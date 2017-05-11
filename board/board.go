package board

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tm "github.com/buger/goterm"
)

type Board struct {
	Size        int
	Board       [][]int
	Turn        int
	TurnCount   int
	initialized bool
}

const White = 2
const Black = 1
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
	if size%2 != 0 {
		fmt.Printf("Board must be an even size")
		return Board{}
	}

	b := make([][]int, size)
	for i := range b {
		b[i] = make([]int, size)
	}
	return Board{
		Size:        size,
		Board:       b,
		Turn:        1,
		TurnCount:   0,
		initialized: true,
	}
}

func (b *Board) drawPrompt() {
	ix := 0
	iy := b.Size + 4

	tm.MoveCursor(iy, ix)

	bPrompt := "Black's Move: input move coords x y"
	wPrompt := "Whites's Move: input move coords x y"

	movePrompts := []string{bPrompt, wPrompt}
	if b.TurnCount < 4 {
		tm.Printf("First moves must be in the central squares of the board.\n")
	}
	tm.Printf(movePrompts[b.Turn-1])
	tm.Flush()
}

func (b *Board) redraw() {
	tm.Clear()
	b.DrawBoard()
	b.drawPrompt()
}

func (b *Board) StartGame() {
	if !b.initialized {
		log.Fatal("Board not net initialized")
	}
	reader := bufio.NewReader(os.Stdin)

	invalidFormat := "Invalid format, input move coords x y"
	invalidOpeningMove := "Invalid move, moves must be in the central four squares"
	invalidMove := "Invalid move, moves must flip at least one stone"

	//first four turns must me taken in the middle squares
	var x int
	var y int
	var err error
	b.redraw()

	hasFailed := false
	for {
		if b.TurnCount > 4 && !b.CanPlay(b.Turn) {
			b.nextTurn()
			if !b.CanPlay(b.Turn) {
				//the game is over!
				tm.MoveCursorUp(1)
				tm.ResetLine("")
				tm.Println("Game over!")
				//calculate winner
				tm.Flush()
				return
			} else {
				tm.Println("No valid moves for player, skipping turn")
				tm.Flush()
			}
		}

		move, _ := reader.ReadString('\n')
		x, y, err = validateInput(move)
		if err != nil {
			if hasFailed {
				tm.ResetLine("")
				tm.MoveCursorUp(1)
				tm.ResetLine("")
				tm.MoveCursorUp(2)
			} else {
				tm.MoveCursorUp(1)
			}
			tm.ResetLine("")
			tm.Println(invalidFormat)
			tm.Flush()
			hasFailed = true
			continue
		}
		if !b.validMove(x, y, b.Turn) {
			if hasFailed {
				tm.ResetLine("")
				tm.MoveCursorUp(1)
				tm.ResetLine("")
				tm.MoveCursorUp(2)
			} else {
				tm.MoveCursorUp(1)
			}
			tm.ResetLine("")
			if b.TurnCount < 4 {
				tm.Println(invalidOpeningMove)
			} else {
				tm.Println(invalidMove)
			}
			tm.Flush()
			hasFailed = true
			continue
		}
		hasFailed = false
		b.PlacePiece(x, y, b.Turn)
		b.nextTurn()
		b.redraw()
	}
}

func (b *Board) nextTurn() {
	b.Turn = ((b.Turn % 2) + 1)
	b.TurnCount++
}

func validateInput(move string) (int, int, error) {

	move = strings.TrimSpace(move)
	if len(move) > 3 {
		return 0, 0, errors.New("Invalid Input")
	}
	//check to make sure the move is in the middle of the board
	moves := strings.Split(move, " ")

	if len(moves) != 2 {
		return 0, 0, errors.New("Invalid Input")
	}

	x, err := strconv.Atoi(moves[0])
	if err != nil {
		return 0, 0, errors.New("Invalid Input")
	}
	y, err := strconv.Atoi(moves[1])
	if err != nil {
		return 0, 0, errors.New("Invalid Input")
	}

	return x, y, nil

}

func (b *Board) validMove(x int, y int, player int) bool {
	if b.TurnCount < 4 {
		midpoint := b.Size / 2
		if x < midpoint-1 || x > midpoint || y < midpoint-1 || x > midpoint {
			return false
		}
		if b.Board[y][x] != Clear {
			return false
		}
		//it's a valid move
		return true
	}

	if b.Board[y][x] != Clear {
		return false
	}

	//need to validate that the move results in a piece being captured
	for k := 0; k < 8; k++ {
		depth, ok := b.checkNext(0, k, x, y, player, false, true)
		if ok && depth > 0 {
			//it's a valid move
			return true
		}
	}
	return false

}

//PlacePice places a piece in that 0-indexed location on the board
//
func (b *Board) PlacePiece(x int, y int, player int) error {

	//make sure that it's a valid player
	if player != White && player != Black {
		return errors.New("Invalid Player")
	}
	if x < 0 || x > b.Size || y < 0 || y > b.Size {
		return errors.New("Invalid position")
	}
	if b.Board[y][x] != Clear {
		return errors.New("Position already taken")
	}

	b.Board[y][x] = player

	for i := 0; i < 8; i++ {
		b.checkNext(0, i, x, y, player, true, true)
	}
	b.DrawBoard()

	return nil

}

/*
CanPlay returns if a player is able to place a stone to flip an opponents piece
*/
func (b *Board) CanPlay(player int) bool {
	//search each field on board, check if has adjacent tile of another color, if so, follow it to see if it has another of the players stones on the other side.

	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			if b.Board[i][j] != player {
				continue
			}

			for k := 0; k < 8; k++ {
				depth, ok := b.checkNext(0, k, i, j, player, false, false)
				if ok && depth > 0 {
					return true
				}
			}
		}
	}
	return false
}

//we need to check in lines from the piece being played in all 8 directions until we hit either
// a) a clear space
// b) a wall
// c) a piece of the same player
//the function that will recursively check the next piece in the line

func (b *Board) checkNext(depth int, dir int, curX int, curY int, player int, flip bool, real bool) (int, bool) {
	if player != White && player != Black {
		return depth, false
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
		return depth, false
	}

	//we're off the board
	if nextX >= b.Size || nextX < 0 || nextY >= b.Size || nextY < 0 {
		return depth, false
	}

	//we hit a clear space
	if b.Board[nextY][nextX] == Clear {
		if flip || real {
			return depth, false
		}
		//we're just checking to see if we can play, we can, since there's a blank space
		return depth, true
	}

	if b.Board[nextY][nextX] == player {
		if !flip && !real {
			return depth, false

		}

		return depth, true
	}
	depth += 1

	d, val := b.checkNext(depth, dir, nextX, nextY, player, flip, real)

	if val {
		if flip {
			b.Board[nextY][nextX] = player //capture the pieces
		}
		return d, true
	}

	return d, false //nothing happened
}

func (b *Board) DrawBoard() {
	tm.MoveCursor(1, 1)
	vals := make([]interface{}, b.Size+1)

	str := ""
	for i := 0; i < b.Size; i++ {
		str += "%v\t"
		vals[i] = i - 1
	}

	str += "%v\n"
	vals[0] = " "
	vals[b.Size] = b.Size - 1

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
