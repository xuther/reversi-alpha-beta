package ab

import "log"

type Node interface {
	GetUtility() int //get the utility
	Branch() []Node  //get the children
	GetMinMax() int  //-1 if minimizing +1 if maximizing
	GetNodeID() string
}

const debug = true

func Search(n Node, depth int, alpha int, beta int) (int, []string) {
	log.Printf("Expanding %v. Which is maxmizing %v", n.GetNodeID(), n.GetMinMax())
	chillins := n.Branch()
	log.Printf("%v has %v children", n.GetNodeID(), len(chillins))
	toPrint := ""
	for i := range chillins {
		toPrint += chillins[i].GetNodeID() + ", "
	}
	log.Printf(toPrint)
	if depth == 0 || len(chillins) == 0 {
		log.Printf("%v is leaf node. Returning %v", n.GetNodeID(), n.GetUtility())
		return n.GetUtility(), []string{n.GetNodeID()}
	}

	if n.GetMinMax() == 1 {
		val := -100000
		chain := []string{}
		for i := range chillins {
			log.Printf("Expanding from %v", n.GetNodeID())
			challenger, tempChain := Search(chillins[i], depth-1, alpha, beta)
			if challenger > val {
				log.Printf("Returned value from %v is new concrete value", chillins[i].GetNodeID())
				val = challenger
				chain = tempChain
			}
			if val > alpha {
				log.Printf("Returned value from %v is new alpha", chillins[i].GetNodeID())
				alpha = val
			}
			if beta <= alpha {
				log.Printf("Trimming all other chillins on node %v", n.GetNodeID())
				break
			}
		}
		log.Printf("%v returning %v", n.GetNodeID(), val)
		return val, append(chain, n.GetNodeID())
	} else {
		val := 100000
		chain := []string{}
		for i := range chillins {
			challenger, tempChain := Search(chillins[i], depth-1, alpha, beta)
			if challenger < val {
				log.Printf("Returned value from %v is new concrete value", chillins[i].GetNodeID())
				val = challenger
				chain = tempChain
			}
			if val < beta {
				log.Printf("Returned value from %v is new beta", chillins[i].GetNodeID())
				beta = val
			}
			if beta <= alpha {
				log.Printf("Trimming all other chillins on node %v", n.GetNodeID())
				break
			}
		}
		log.Printf("%v returning %v", n.GetNodeID(), val)
		return val, append(chain, n.GetNodeID())

	}
}
