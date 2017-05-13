package ab

import (
	"log"
	"sync"
)

type update struct {
	Type    AlphaBeta
	Val     int
	ChildID int
	chain   []string
}

type AlphaBeta uint8

const (
	ALPHA AlphaBeta = iota
	BETA
	REAL
)

func StartSearchMulti(n Node, depth int, alpha int, beta int) (int, []string) {
	log.Printf("Expanding %v. Which is maxmizing %v", n.GetNodeID(), n.GetMinMax())
	chillins := n.Branch()
	log.Printf("%v has %v children", n.GetNodeID(), len(chillins))
	toPrint := ""
	for i := range chillins {
		toPrint += chillins[i].GetNodeID() + ", "
	}
	log.Printf(toPrint)
	if depth == 0 || len(chillins) == 0 {
		id := n.GetNodeID()
		util := n.GetUtility()
		log.Printf("%v is leaf node. Returning %v", id, util)
		return n.GetUtility(), []string{n.GetNodeID()}
	}

	val := 0
	if n.GetMinMax() == 1 {
		val = -100000
	} else {
		val = 100000
	}

	chain := []string{}
	curBestChild := -1

	childRecvChannel := make(chan update, 50)
	childSendChannels := make([]chan update, len(chillins))
	finished := make(chan int)

	mywg := sync.WaitGroup{}

	for i := range chillins {
		log.Printf("Start child search processes from %v", n.GetNodeID())
		childSendChannels[i] = make(chan update, 10)
		mywg.Add(1)
		go SearchMulti(chillins[i], depth-1, alpha, beta, childRecvChannel, childSendChannels[i], &mywg, i)
	}
	go func() {
		mywg.Wait()
		log.Printf("ROOT sending done signal")
		finished <- 0
	}()

	for {
		log.Printf("ROOT Waiting")
		alphaChanged := false
		betaChanged := false
		select {
		case challenger := <-childRecvChannel:
			if n.GetMinMax() == 1 {
				if challenger.Val > val || challenger.ChildID == curBestChild {
					log.Printf("New best local chain is %v", challenger.chain)
					val = challenger.Val
					chain = challenger.chain
					curBestChild = challenger.ChildID
				}
				if val > alpha {
					log.Printf("new alpha chain %v found with value %v", chain, val)
					alpha = val
					alphaChanged = true
				}
			} else {
				if challenger.Val < val || challenger.ChildID == curBestChild {
					log.Printf("New best local chain is %v", challenger.chain)
					val = challenger.Val
					chain = challenger.chain
					curBestChild = challenger.ChildID
				}
				if val < beta {
					log.Printf("new alpha chain %v found with value %v", chain, val)
					beta = val
					alphaChanged = true
				}

			}
		case <-finished:
			//drain our child channel
			goOn := true
			for goOn == true { //now we check if we've received any new information, essentially draining the channel
				select {
				case challenger := <-childRecvChannel:
					if n.GetMinMax() == 1 {
						if challenger.Val > val || challenger.ChildID == curBestChild {
							log.Printf("New best local chain is %v", challenger.chain)
							val = challenger.Val
							chain = challenger.chain
							curBestChild = challenger.ChildID
						}
						if val > alpha {
							log.Printf("new alpha chain %v found with value %v", chain, val)
							alpha = val
							alphaChanged = true
						}
					} else {
						if challenger.Val < val || challenger.ChildID == curBestChild {
							log.Printf("New best local chain is %v", challenger.chain)
							val = challenger.Val
							chain = challenger.chain
							curBestChild = challenger.ChildID
						}
						if val < beta {
							log.Printf("new alpha chain %v found with value %v", chain, val)
							beta = val
							alphaChanged = true
						}

					}
				default:
					goOn = false

				}
			}

			log.Printf("%v returning %v", n.GetNodeID(), val)
			return val, chain
		}

		if beta <= alpha {
			log.Printf("Trimming all other chillins on node %v", n.GetNodeID())
			break
		}

		//update our children
		if n.GetMinMax() == 1 && alphaChanged {
			log.Printf("Node %v updating children alpha to %v", n.GetNodeID(), alpha)
			for i := range childSendChannels {
				childSendChannels[i] <- update{Type: ALPHA, Val: alpha}
			}
		}
		if n.GetMinMax() == -1 && betaChanged {
			log.Printf("Node %v updating children beta to %v", n.GetNodeID(), beta)
			for i := range childSendChannels {
				childSendChannels[i] <- update{Type: BETA, Val: beta}
			}
		}
	}
	return val, chain
}

func SearchMulti(n Node, depth int, alpha int, beta int, send chan<- update, receive <-chan update, wg *sync.WaitGroup, childID int) {
	log.Printf("Expanding %v. Which is maxmizing %v", n.GetNodeID(), n.GetMinMax())
	chillins := n.Branch()

	log.Printf("%v has %v children", n.GetNodeID(), len(chillins))
	toPrint := ""

	for i := range chillins {
		toPrint += chillins[i].GetNodeID() + ", "
	}
	log.Printf(toPrint)

	if depth == 0 || len(chillins) == 0 {
		id := n.GetNodeID()
		util := n.GetUtility()
		log.Printf("%v is leaf node. Returning %v", id, util)
		send <- update{Type: REAL, Val: util, chain: []string{id}, ChildID: childID}
		wg.Done()
		return
	}

	val := 0
	chain := []string{}
	curBestChild := -1

	if depth == 1 {
		if n.GetMinMax() == 1 {
			val = -100000
		} else {
			val = 100000
		}
		for i := range chillins {

			log.Printf("Expanding from %v", n.GetNodeID())
			challenger, tempChain := Search(chillins[i], depth-1, alpha, beta)

			if n.GetMinMax() == 1 {
				if challenger > val {
					log.Printf("%v Returned value from %v is new concrete value", n.GetNodeID(), chillins[i].GetNodeID())
					val = challenger
					chain = tempChain
				}
				if val > alpha {
					log.Printf("%v Returned value from %v is new alpha", n.GetNodeID(), chillins[i].GetNodeID())
					alpha = val
					send <- update{Type: ALPHA, Val: alpha, chain: append(chain, n.GetNodeID()), ChildID: childID}
				}

			} else {
				if challenger < val {
					log.Printf("%v Returned value from %v is new concrete value", n.GetNodeID(), chillins[i].GetNodeID())
					val = challenger
					chain = tempChain
				}
				if val < beta {
					log.Printf("%v Returned value from %v is new beta", n.GetNodeID(), chillins[i].GetNodeID())
					beta = val
					send <- update{Type: BETA, Val: beta, chain: append(chain, n.GetNodeID()), ChildID: childID}
				}
			}

			goOn := true
			for goOn == true { //now we check if we've received any new information, essentially draining the channel
				select {
				case newUpdate := <-receive:
					log.Printf("%v received updated %v", n.GetNodeID(), newUpdate)
					switch newUpdate.Type {
					case ALPHA:
						if newUpdate.Val > alpha && newUpdate.ChildID != childID {
							log.Printf("%v received new alpha, updating to: %v", n.GetNodeID(), newUpdate.Val)
							alpha = newUpdate.Val
						}
					case BETA:
						if newUpdate.Val < beta && newUpdate.ChildID != childID {
							log.Printf("%v received new beta, updating to: %v", n.GetNodeID(), newUpdate.Val)
							beta = newUpdate.Val
						}
					}
				default:
					goOn = false

				}
			}
			if beta <= alpha {
				log.Printf("Trimming all other chillins on node %v. Beta: %v, Alpha %v", n.GetNodeID())
				break
			}
		}
	} else {
		if n.GetMinMax() == 1 {
			val = -100000
		} else {
			val = 100000
		}

		childRecvChannel := make(chan update, 50)
		childSendChannels := make([]chan update, len(chillins))
		finished := make(chan int)
		mywg := sync.WaitGroup{}

		for i := range chillins {
			log.Printf("Start child search processes from %v", n.GetNodeID())
			childSendChannels[i] = make(chan update, 10)
			mywg.Add(1)
			go SearchMulti(chillins[i], depth-1, alpha, beta, childRecvChannel, childSendChannels[i], &mywg, i)
		}
		go func() {
			mywg.Wait()
			finished <- 0
		}()

		for {
			alphaChanged := false
			betaChanged := false
			updatedAlpha := -1
			updatedBeta := -1
			select {
			case challenger := <-childRecvChannel:

				if n.GetMinMax() == 1 {

					if challenger.Val > val || curBestChild == challenger.ChildID {
						log.Printf("New best local chain is %v", challenger.chain)
						val = challenger.Val
						chain = challenger.chain
						curBestChild = challenger.ChildID
					}
					if val > alpha {
						log.Printf("%v new alpha chain %v found with value %v", n.GetNodeID(), chain, val)
						alpha = val
						send <- update{Type: ALPHA, Val: alpha, chain: append(chain, n.GetNodeID()), ChildID: childID}
						alphaChanged = true
						updatedAlpha = challenger.ChildID
					}
				} else {
					if challenger.Val < val || curBestChild == challenger.ChildID {
						log.Printf("New best local chain is %v", challenger.chain)
						val = challenger.Val
						chain = challenger.chain
						curBestChild = challenger.ChildID
					}
					if val < beta {
						log.Printf("%v new beta chain %v found with value %v", n.GetNodeID(), chain, val)
						beta = val
						send <- update{Type: BETA, Val: beta, chain: append(chain, n.GetNodeID()), ChildID: childID}
						betaChanged = true
						updatedBeta = challenger.ChildID
					}

				}
			case newUpdate := <-receive:
				switch newUpdate.Type {
				case ALPHA:
					if newUpdate.Val > alpha && newUpdate.ChildID != childID {
						alpha = newUpdate.Val
						log.Printf("%v received new alpha, updating to: %v", n.GetNodeID(), newUpdate.Val)
						alphaChanged = true
						updatedAlpha = -1
					}
				case BETA:
					if newUpdate.Val < beta && newUpdate.ChildID != childID {
						beta = newUpdate.Val
						log.Printf("%v received new beta, updating to: %v", n.GetNodeID(), newUpdate.Val)
						betaChanged = true
						updatedBeta = -1
					}
				}
			case <-finished:
				//drain our child channel
				goOn := true
				for goOn == true { //now we check if we've received any new information, essentially draining the channel
					select {
					case challenger := <-childRecvChannel:

						if n.GetMinMax() == 1 {

							if challenger.Val > val || curBestChild == challenger.ChildID {
								log.Printf("New best local chain is %v", challenger.chain)
								val = challenger.Val
								chain = challenger.chain
								curBestChild = challenger.ChildID
							}
							if val > alpha {
								log.Printf("%v new alpha chain %v found with value %v", n.GetNodeID(), chain, val)
								alpha = val
								send <- update{Type: ALPHA, Val: alpha, chain: append(chain, n.GetNodeID()), ChildID: childID}
								alphaChanged = true
								updatedAlpha = challenger.ChildID
							}
						} else {
							if challenger.Val < val || curBestChild == challenger.ChildID {
								log.Printf("New best local chain is %v", challenger.chain)
								val = challenger.Val
								chain = challenger.chain
								curBestChild = challenger.ChildID
							}
							if val < beta {
								log.Printf("%v new beta chain %v found with value %v", n.GetNodeID(), chain, val)
								beta = val
								send <- update{Type: BETA, Val: beta, chain: append(chain, n.GetNodeID()), ChildID: childID}
								betaChanged = true
								updatedBeta = challenger.ChildID
							}

						}

					default:
						goOn = false

					}
				}

				log.Printf("%v returning %v", n.GetNodeID(), val)
				send <- update{Type: REAL, Val: val, chain: append(chain, n.GetNodeID()), ChildID: childID}
				wg.Done()
				return
			}

			goOn := true
			for goOn == true { //now we check if we've received any new information, essentially draining the channel
				select {
				case newUpdate := <-receive:
					switch newUpdate.Type {
					case ALPHA:
						if newUpdate.Val > alpha && newUpdate.ChildID != childID {
							alpha = newUpdate.Val
							log.Printf("%v received new alpha, updating to: %v", n.GetNodeID(), newUpdate.Val)
							alphaChanged = true
							updatedAlpha = -1
						}
					case BETA:
						if newUpdate.Val < beta && newUpdate.ChildID != childID {
							beta = newUpdate.Val
							log.Printf("%v received new beta, updating to: %v", n.GetNodeID(), newUpdate.Val)
							betaChanged = true
							updatedBeta = -1
						}
					}

				default:
					goOn = false

				}
			}
			if beta <= alpha {
				log.Printf("Trimming all other chillins on node %v", n.GetNodeID())
				break
			}

			//update our children
			if alphaChanged {
				log.Printf("Node %v updating children alpha to %v", n.GetNodeID(), alpha)
				for i := range childSendChannels {
					childSendChannels[i] <- update{Type: ALPHA, Val: alpha, ChildID: updatedAlpha}
				}
			}
			if betaChanged {
				log.Printf("Node %v updating children beta to %v", n.GetNodeID(), beta)
				for i := range childSendChannels {
					childSendChannels[i] <- update{Type: BETA, Val: beta, ChildID: updatedBeta}
				}
			}
		}
	}
	log.Printf("%v returning %v", n.GetNodeID(), val)
	send <- update{Type: REAL, Val: val, chain: append(chain, n.GetNodeID()), ChildID: childID}
	wg.Done()
	return

}
