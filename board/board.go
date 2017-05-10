package board

type Board struct {
	Size  int
	Board [][]byte
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

func InitializeBoard(int size) Board {
	b := make([][]byte, size)
	for i := range b {
		b[i] = make([]byte, size)
	}
	return Board{
		Size:  size,
		Board: b,
		Turn:  false,
	}
}

//PlacePice places a piece in that 0-indexed location on the board
//
func (b *Board) PlacePiece(int x, int y, int player) {

	//make sure that it's a valid player
	if player != White || player != Black {
		return
	}
	//we need to check in lines from the piece being played in all 8 directions until we hit either
	// a) a clear space
	// b) a wall
	// c) anoter pice by same player


//the function that will recursively check the next piece in the line
func (b *Board) checkNext(int direction, int curX, int curY) bool {
	nextY := curY
	nextX := curX
	switch direction {
	case N:
		nextY--
	case NE: 
		nextY--
		nextX++
	case E:
		nextX++

	}
}

func (b *Board) DrawBoard() {

}
