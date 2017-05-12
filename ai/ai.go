package ai

import (
	"log"

	"github.com/xuther/reversi-alpha-beta/board"
)

type Node struct {
	Board       board.Board
	Alpha       int
	Beta        int
	Depth       int
	Move        []int
	Max         bool // true means max false means min
	Range       []int
	Parent      *Node
	Children    []Node
	CurBestMove Move
}

type Move struct {
	Utility int
	Move    []int
	child   *Node
}

const debug = true

var maxUtility = 1000
var minUtility = -1000
var maxDepth = 2

func spawnChildNode(parent *Node, move []int) Node {
	log.Printf("Spawning child with move %v", move)
	log.Printf("Parent: %v", parent.Board)
	//tempBoard, make the move being considered so that the board here matches the move being considered
	temp := parent.Board.Copy()
	temp.PlacePiece(move[0], move[1], parent.Board.Turn, false)
	temp.NextTurn()

	max := !parent.Max

	log.Printf("%v", temp.Turn)

	//check if the next player can play, else switch turns again
	if !temp.CanPlay(temp.Turn) {
		if debug {
			log.Printf("No valid moves for player, switching turn")
		}
		max = !max //swap max again
		temp.NextTurn()
	}

	tempMove := Move{}
	//we're looking to minimize here
	if !max {
		tempMove.Utility = 10000
	} else {
		//we're looking to Maximize here
		tempMove.Utility = -10000
	}

	return Node{
		Board:       temp,
		Alpha:       maxUtility,
		Beta:        minUtility,
		Depth:       parent.Depth + 1,
		Move:        move,
		Max:         max,
		Range:       parent.Range,
		Parent:      parent,
		Children:    []Node{},
		CurBestMove: tempMove,
	}
}

var maximizingFor = 0

//assume maximizing for current player
func GetNextMove(b *board.Board) []int {
	maxUtility = b.Size * 2
	minUtility = b.Size * 2 * -1

	tempMove := Move{}
	tempMove.Utility = -10000
	Root := Node{
		Board:       b.Copy(),
		Depth:       0,
		Alpha:       maxUtility,
		Beta:        minUtility,
		Move:        []int{},
		Max:         true,
		Range:       []int{minUtility, maxUtility},
		Parent:      nil,
		Children:    []Node{},
		CurBestMove: tempMove,
	}

	maximizingFor = b.Turn

	//start dfs
	moves := Root.Board.GetAllPossibleMoves(Root.Board.Turn)
	if len(moves) < 1 {
		log.Printf("invalid moves")
		return []int{-1, -1}
	}
	//there's only one move, might as well return it now
	if len(moves) == 1 {
		//	return moves[0]
	}

	for i := range moves {
		Root.Children = append(Root.Children, spawnChildNode(&Root, moves[i]))
	}

	for i := range Root.Children {
		if i > 0 {
			break

		}
		contender := Search(&Root.Children[i])

		//This is where we would do our pruning, based on the alpha/beta

		if contender > Root.CurBestMove.Utility {
			log.Printf("New best maximizing utility")
			Root.CurBestMove = Move{
				Utility: contender,
				Move:    Root.Children[i].Move,
			}
			if debug {
				Root.CurBestMove.child = &Root.Children[i]
			}
		}

	}
	if debug {
		log.Printf("Best move found: %+v", Root.CurBestMove)
		log.Printf("Move sequence:")
		log.Printf("%v", Root.CurBestMove.Move)
		Root.CurBestMove.child.PrintTree()

	}
	return Root.CurBestMove.Move
}

func (n *Node) PrintTree() {
	if len(n.Children) == 0 {
		n.Board.PrintBoard()
		return
	}
	log.Printf("%v", n.Move)
	n.CurBestMove.child.PrintTree()
	return
}

func Search(n *Node) int {
	if debug {
		log.Printf("Searching on move %v, depth %v, player %v, max %v", n.Move, n.Depth, n.Board.Turn, n.Max)
	}

	moves := n.Board.GetAllPossibleMoves(n.Board.Turn)

	if n.Depth >= maxDepth || len(moves) == 0 {
		//evaluate and return

		//basic eval function, number of stones captured is the utility
		score := n.Board.CalculateScore()
		utility := score[maximizingFor] - score[maximizingFor%2+1]
		if debug {
			log.Printf("Hit base score : %v", utility)
		}
		return utility

		/*DEBUG
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		return r1.Intn(1000) - 500
		*/
	}

	for i := range moves {
		n.Children = append(n.Children, spawnChildNode(n, moves[i]))
	}
	for i := range n.Children {
		contender := Search(&n.Children[i])
		if debug {
			log.Printf("child returned with utilty: %v", contender)

		}

		//This is where we would do our pruning, based on the alpha/beta

		if n.Max {
			if contender > n.CurBestMove.Utility {
				log.Printf("Child is new best maximizing utility")
				n.Beta = contender
				//check if our beta is >= our parents alpha.
				n.CurBestMove = Move{
					Utility: contender,
					Move:    n.Children[i].Move,
				}
				if debug {
					n.CurBestMove.child = &n.Children[i]
				}
				if n.Beta >= n.Parent.Alpha {
					//no Point in continuing
					if debug {
						log.Printf("Pruning")
					}
					return contender
				}
			}
		} else {
			if contender < n.CurBestMove.Utility {
				log.Printf("Child is new best minimizing utility")
				n.Alpha = contender
				n.CurBestMove = Move{
					Utility: contender,
					Move:    n.Children[i].Move,
				}
				if debug {
					n.CurBestMove.child = &n.Children[i]
				}
				//check if our alpha is <= our parents beta
				if n.Alpha <= n.Parent.Beta {
					//no point in continuing
					if debug {
						log.Printf("Pruning")
					}
					return contender
				}
			}
		}
	}
	return n.CurBestMove.Utility
}
