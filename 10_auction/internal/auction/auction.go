//go:generate scribblec-param.sh ../../Auction.scr -d ../../ -param Protocol github.com/nickng/scribble-go-examples/10_auction/Auction -param-api Auctioneer -param-api Bidder

package auction

import (
	"encoding/gob"
	"fmt"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/10_auction/Auction/Protocol"
	"github.com/nickng/scribble-go-examples/10_auction/Auction/Protocol/Auctioneer_1to1"
	"github.com/nickng/scribble-go-examples/10_auction/Auction/Protocol/Bidder_1toK"
	"github.com/nickng/scribble-go-examples/10_auction/message"
	"github.com/nickng/scribble-go-examples/scributil"
)

func init() {
	gob.Register(new(message.Initial))
	gob.Register(new(message.Highest))
	gob.Register(new(message.BidOrSkip))
	gob.Register(new(message.Winner))
}

// Auctioneer is the implementation of Auctioneer[1].
func Auctioneer(p *Protocol.Protocol, K, self int, cc scributil.ClientConn, host string, baseport int, wg *sync.WaitGroup) {
	Auctioneer := p.New_Auctioneer_1to1(K, self)

	wgCli := new(sync.WaitGroup)
	wgCli.Add(K)
	mu := new(sync.Mutex) // for concurrent access of connection map
	for k := 1; k <= K; k++ {
		go func(k int) {
			mu.Lock()
			if err := Auctioneer.Bidder_1toK_Dial(k, host, baseport+k, cc.Dial, cc.Formatter()); err != nil {
				log.Fatalf("cannot dial: %v", err)
			}
			mu.Unlock()
			wgCli.Done()
		}(k)
	}
	wgCli.Wait()

	Auctioneer.Run(func(s *Auctioneer_1to1.Init) Auctioneer_1to1.End {
		initBids := make([]message.Initial, K)
		s0 := s.Bidder_1toK_Gather_InitialBid(initBids)
		scributil.Debugf("[info] Auctioneer: initial bids: %v.\n", initBids)
		highest, _, _ := topInitialBid(initBids)
		s1 := s0.Bidder_1toK_Scatter_HighestBid(highest)
		for {
			bids := make([]message.BidOrSkip, K)
			k := 0
			s2 := s1.Foreach(func(s *Auctioneer_1to1.Init_19) Auctioneer_1to1.End {
				sEnd := s.Bidder_I_Gather_BidOrSkip(bids[k:])
				k++
				return *sEnd
			})
			scributil.Debugf("[info] Auctioneer: received bids: %v.\n", bids)

			if winner := findWinner(bids); winner == nil { // no winner
				highest, topbid, topbidder := topBid(bids)
				scributil.Debugf("[info] Auctioneer: Highest bid %d by Bidder[%d].\n", topbid, topbidder)
				s1 = s2.Bidder_1toK_Scatter_HighestBid(highest)
			} else {
				winner := findWinner(bids)
				scributil.Debugf("[info] Auctioneer: Bidder[%d] wins.\n", winner[0].BidderID)
				sEnd := s2.Bidder_1toK_Scatter_Winner(winner)
				return *sEnd
			}
		}
	})
	wg.Done()
}

// Bidder is the implementations of Bidder[1] ... Bidder[K].
func Bidder(p *Protocol.Protocol, K, self int, sc scributil.ServerConn, port int, wg *sync.WaitGroup) {
	Bidder := p.New_Bidder_1toK(K, self)

	ln, err := sc.Listen(port)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	if err := Bidder.Auctioneer_1to1_Accept(1, ln, sc.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}

	Bidder.Run(func(s *Bidder_1toK.Init) Bidder_1toK.End {
		// As an over-simplification the initial bid is the bidder ID (i.e. 1..K) and
		// the maximum bid is always initial bid + 5 (so the last bidder always wins)
		// This allows for a deterministic output to see what is happening
		// in the duration of the protocol.
		initBid := self
		maxBid := initBid + 5
		fmt.Printf("Bidder[%d]: Initial bid %d (max bid %d)\n", self, initBid, maxBid)
		initBidMsg := []message.Initial{message.Initial{Bid: initBid}}
		s0 := s.Auctioneer_1_Scatter_InitialBid(initBidMsg)
		scributil.Debugf("[info] Bidder[%d]: initial bids: %d.\n", self, initBid)
		highestBid := make([]message.Highest, 1)
		s1 := s0.Auctioneer_1_Gather_HighestBid(highestBid)
		for {
			bid := []message.BidOrSkip{message.BidOrSkip{}}
			if highestBid[0].Bid+1 > maxBid {
				bid[0].MakeBid = false
				fmt.Printf("Bidder[%d]: I give up.\n", self)
			} else {
				bid[0].Bid = highestBid[0].Bid + 1
				bid[0].MakeBid = true
				fmt.Printf("Bidder[%d]: new bid at %d.\n", self, bid[0].Bid)
			}
			s2 := s1.Auctioneer_1_Scatter_BidOrSkip(bid)
			scributil.Debugf("[info] Bidder[%d] sent bid: %v.\n", self, bid)
			switch s3 := s2.Auctioneer_1_Branch().(type) {
			case *Bidder_1toK.HighestBid:
				s1 = s3.Recv_HighestBid(&highestBid[0])
			case *Bidder_1toK.Winner:
				winner := new(message.Winner)
				sEnd := s3.Recv_Winner(winner)
				if winner.BidderID != self {
					fmt.Printf("Bidder[%d]: I lost to Bidder[%d]!\n", self, winner.BidderID)
				} else {
					fmt.Printf("Bidder[%d]: I won!\n", self)
				}
				return *sEnd
			}
		}
	})
	wg.Done()
}

func topInitialBid(bids []message.Initial) (highest []message.Highest, bid int, bidderID int) {
	maxBid, maxIdx := bids[0].Bid, 0
	for i, bid := range bids {
		if bid.Bid > maxBid {
			maxBid, maxIdx = bid.Bid, i
		}
	}
	return highestBids(len(bids), maxBid, maxIdx)
}

func topBid(bids []message.BidOrSkip) (highest []message.Highest, bid int, bidderID int) {
	maxBid, maxIdx := bids[0].Bid, 0
	for i, bid := range bids {
		if bid.Bid > maxBid {
			maxBid, maxIdx = bid.Bid, i
		}
	}
	return highestBids(len(bids), maxBid, maxIdx)
}

func highestBids(nBids, maxBid, maxIdx int) (highest []message.Highest, bid, bidderID int) {
	highestBids := make([]message.Highest, nBids)
	for i := range highestBids {
		highestBids[i].Bid = maxBid
	}
	return highestBids, maxBid, maxIdx + 1
}

func findWinner(bids []message.BidOrSkip) []message.Winner {
	bidders, index := 0, 0
	for i := range bids {
		if bids[i].MakeBid {
			bidders, index = bidders+1, i
		}
	}
	if bidders == 1 {
		winner := make([]message.Winner, len(bids))
		for i := range winner {
			winner[i].BidderID = index + 1
		}
		return winner
	}
	return nil
}
