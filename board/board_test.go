package board

import (
	"bufio"
	"os"
	"testing"
)

func TestGeneration(t *testing.T) {
	b := InitializeBoard(10)
	b.DrawBoard()

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
