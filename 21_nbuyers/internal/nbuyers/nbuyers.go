//go:generate scribblec-param.sh ../../NBuyers.scr -d ../../ -param NBuyers github.com/nickng/scribble-go-examples/21_nbuyers/NBuyers -param-api Buyer -param-api Seller

package nbuyer

import (
	"encoding/gob"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nickng/scribble-go-examples/21_nbuyers/NBuyers/NBuyers"
	Buyer_1 "github.com/nickng/scribble-go-examples/21_nbuyers/NBuyers/NBuyers/family_1/Buyer_1to1and1toKand1toKsub1_not_2toKandKtoK"
	Buyer_i "github.com/nickng/scribble-go-examples/21_nbuyers/NBuyers/NBuyers/family_1/Buyer_1toKand1toKsub1and2toK_not_1to1andKtoK"
	Buyer_K "github.com/nickng/scribble-go-examples/21_nbuyers/NBuyers/NBuyers/family_1/Buyer_1toKand2toKandKtoK_not_1to1and1toKsub1"
	"github.com/nickng/scribble-go-examples/21_nbuyers/NBuyers/NBuyers/family_1/Seller_1to1"
	Buyer_1_family2 "github.com/nickng/scribble-go-examples/21_nbuyers/NBuyers/NBuyers/family_2/Buyer_1to1and1toKand1toKsub1_not_2toKandKtoK"
	Buyer_K_family2 "github.com/nickng/scribble-go-examples/21_nbuyers/NBuyers/NBuyers/family_2/Buyer_1toKand2toKandKtoK_not_1to1and1toKsub1"
	"github.com/nickng/scribble-go-examples/21_nbuyers/message"
	"github.com/nickng/scribble-go-examples/scributil"
)

// WalletAmount is the amount in Buyer's wallet,
// and controls if a book is affordable.
const WalletAmount = 30

func init() {
	gob.Register(new(message.Address))
	gob.Register(new(message.Date))
}

// Seller implements Seller[1].
func Seller(p *NBuyers.NBuyers, K, self int, sc scributil.ServerConn, baseport int, wg *sync.WaitGroup) {
	Seller := p.New_family_1_Seller_1to1(K, self)

	wgSvr := new(sync.WaitGroup)
	wgSvr.Add(K)
	for i := 1; i <= K; i++ {
		go func(i int) {
			ln, err := sc.Listen(baseport + i)
			if err != nil {
				log.Fatalf("cannot listen: %v", err)
			}
			scributil.Debugf("[connection] Seller: listening for Buyer[%d] at :%d.\n", i, baseport+i)
			if 1 == i {
				// Buyer[1..K-1]\[2..K]
				if err := Seller.Buyer_1to1and1toKand1toKsub1_not_2toKandKtoK_Accept(i, ln, sc.Formatter()); err != nil {
					log.Fatalf("cannot accept: %v", err)
				}
			}
			if 2 <= i && i <= K-1 {
				// Buyer[1..K]\[1,K]
				if err := Seller.Buyer_1toKand1toKsub1and2toK_not_1to1andKtoK_Accept(i, ln, sc.Formatter()); err != nil {
					log.Fatalf("cannot accept: %v", err)
				}
			}
			if i == K {
				// Buyer[1..K]\[1,1..K-1]
				if err := Seller.Buyer_1toKand2toKandKtoK_not_1to1and1toKsub1_Accept(i, ln, sc.Formatter()); err != nil {
					log.Fatalf("cannot accept: %v", err)
				}
			}
			wgSvr.Done()
		}(i)
	}
	wgSvr.Wait()
	scributil.Debugf("[connection] Seller: Ready.\n")

	Seller.Run(func(s *Seller_1to1.Init) Seller_1to1.End {
		title := make([]string, 1)
		s0 := s.Buyer_1_Gather_Title(title)
		quote := make([]int, K)
		for i := range quote {
			quote[i] = 100
		}
		var sEnd Seller_1to1.End
		s1 := s0.Buyer_1toK_Scatter_Quote(quote)
		switch s2 := s1.Buyer_K_Branch().(type) {
		case *Seller_1to1.OK:
			scributil.Debugf("Seller: OK.\n")
			s3 := s2.Recv_OK()
			addr := make([]message.Address, K)
			s4 := s3.Buyer_K_Gather_Address(addr)
			scributil.Debugf("Seller: received %v.\n", addr)
			deliverDate := tomorrow()
			s5 := s4.Buyer_K_Scatter_Date(deliverDate)
			scributil.Debugf("Seller: sent %v.\n", deliverDate)
			sEnd = *s5
		case *Seller_1to1.Quit:
			scributil.Debugf("Seller: Quit.\n")
			s3 := s2.Recv_Quit()
			sEnd = *s3
		}
		return sEnd
	})
	wg.Done()
}

// Buyer1 implements Buyer[1]
func Buyer1(p *NBuyers.NBuyers, K, self int, seller scributil.ClientConn, host string, port int, buyer scributil.ClientConn, buyer2Host string, buyer2Port int, wg *sync.WaitGroup) {
	Buyer1 := p.New_family_1_Buyer_1to1and1toKand1toKsub1_not_2toKandKtoK(K, self)

	scributil.Debugf("[connection] Buyer[%d]: dialling to Seller at %s:%d.\n", self, host, port)
	if err := Buyer1.Seller_1to1_Dial(1, host, port, seller.Dial, seller.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	scributil.Debugf("[connection] Buyer[%d]: dialling to Buyer[%d] at %s:%d.\n", self, self+1, buyer2Host, buyer2Port)
	if err := Buyer1.Buyer_1toKand1toKsub1and2toK_not_1to1andKtoK_Dial(self+1, buyer2Host, buyer2Port, buyer.Dial, buyer.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	scributil.Debugf("[connection] Buyer[%d]: Ready.\n", self)

	Buyer1.Run(func(s *Buyer_1.Init) Buyer_1.End {
		title := []string{"Harry Potter and the Philosopher's Stone"}
		s0 := s.Seller_1_Scatter_Title(title)
		scributil.Debugf("Buyer[%d]: sent %v.\n", self, title)
		quoteTotal := make([]int, 1)
		s1 := s0.Seller_1_Gather_Quote(quoteTotal)
		scributil.Debugf("Buyer[%d]: received %v.\n", self, quoteTotal)
		quoteRemain := []int{quoteTotal[0] - quoteTotal[0]/K}
		sEnd := s1.Buyer_2_Scatter_Quote(quoteRemain)
		scributil.Debugf("Buyer[%d]: sent %v.\n", self, quoteRemain)
		return *sEnd
	})
	wg.Done()
}

// Buyer implements Buyer[2..K-1]
func Buyer(p *NBuyers.NBuyers, K, self int, seller scributil.ClientConn, host string, port int, buyer1 scributil.ServerConn, buyer1Port int, buyerN scributil.ClientConn, buyerNHost string, buyerNPort int, wg *sync.WaitGroup) {
	Buyer := p.New_family_1_Buyer_1toKand1toKsub1and2toK_not_1to1andKtoK(K, self)

	scributil.Debugf("[connection] Buyer[%d]: dialling to Seller at %s:%d.\n", self, host, port)
	if err := Buyer.Seller_1to1_Dial(1, host, port, seller.Dial, seller.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	ln, err := buyer1.Listen(buyer1Port)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	scributil.Debugf("[connection] Buyer[%d]: listening for Buyer[%d] at :%d.\n", self, self-1, buyer1Port)
        if self == 2 {
		if err := Buyer.Buyer_1to1and1toKand1toKsub1_not_2toKandKtoK_Accept(self-1, ln, buyer1.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	} else {
		if err := Buyer.Buyer_1toKand1toKsub1and2toK_not_1to1andKtoK_Accept(self-1, ln, buyer1.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
        }
	scributil.Debugf("[connection] Buyer[%d]: dialling to Buyer[%d] at %s:%d.\n", self, self+1, buyerNHost, buyerNPort)
	if self+1 == K {
		if err := Buyer.Buyer_1toKand2toKandKtoK_not_1to1and1toKsub1_Dial(self+1, buyerNHost, buyerNPort, buyerN.Dial, buyerN.Formatter()); err != nil {
			log.Fatalf("cannot dial: %v", err)
		}
	} else {
		if err := Buyer.Buyer_1toKand1toKsub1and2toK_not_1to1andKtoK_Dial(self+1, buyerNHost, buyerNPort, buyerN.Dial, buyerN.Formatter()); err != nil {
			log.Fatalf("cannot dial: %v", err)
		}
	}
	scributil.Debugf("[connection] Buyer[%d]: Ready.\n", self)

	Buyer.Run(func(s *Buyer_i.Init) Buyer_i.End {
		quoteTotal := make([]int, 1)
		s0 := s.Seller_1_Gather_Quote(quoteTotal)
		scributil.Debugf("Buyer[%d]: received %v.\n", self, quoteTotal)
		amtToPay := make([]int, 1)
		s1 := s0.Buyer_selfsub1_Gather_Quote(amtToPay)
		scributil.Debugf("Buyer[%d]: received %v.\n", self, amtToPay)
		quoteRemain := []int{amtToPay[0] - quoteTotal[0]/K}
		sEnd := s1.Buyer_selfplus1_Scatter_Quote(quoteRemain)
		scributil.Debugf("Buyer[%d]: sent %v.\n", self, quoteRemain)
		return *sEnd
	})
	wg.Done()
}

// BuyerK impements Buyer[K].
func BuyerK(p *NBuyers.NBuyers, K, self int, seller scributil.ClientConn, host string, port int, buyer scributil.ServerConn, buyerPort int, wg *sync.WaitGroup) {
	BuyerK := p.New_family_1_Buyer_1toKand2toKandKtoK_not_1to1and1toKsub1(K, self)

	if err := BuyerK.Seller_1to1_Dial(1, host, port, seller.Dial, seller.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	ln, err := buyer.Listen(buyerPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	scributil.Debugf("[connection] Buyer[%d]: listening for Buyer[%d] at :%d.\n", self, self-1, buyerPort)
	if err := BuyerK.Buyer_1toKand1toKsub1and2toK_not_1to1andKtoK_Accept(self-1, ln, buyer.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	scributil.Debugf("[connection] Buyer[%d]: Ready.\n", self)

	BuyerK.Run(func(s *Buyer_K.Init) Buyer_K.End {
		quoteTotal := make([]int, 1)
		s0 := s.Seller_1_Gather_Quote(quoteTotal)
		scributil.Debugf("Buyer[%d]: received %v.\n", self, quoteTotal)
		amtToPay := make([]int, 1)
		s1 := s0.Buyer_selfsub1_Gather_Quote(amtToPay)
		scributil.Debugf("Buyer[%d]: received %v.\n", self, amtToPay)
		if isAffordable(amtToPay) {
			fmt.Printf("Book is affordable (Cost: %d <= %d)\n", amtToPay[0], WalletAmount)
			s2 := s1.Seller_1_Scatter_OK()
			scributil.Debugf("Buyer[%d]: sent OK.\n", self)
			addr := getAddress()
			s3 := s2.Seller_1_Scatter_Address(addr)
			scributil.Debugf("Buyer[%d]: sent %v.\n", self, addr)
			deliveryDate := make([]message.Date, 1)
			sEnd := s3.Seller_1_Gather_Date(deliveryDate)
			scributil.Debugf("Buyer[%d]: received %v.\n", self, deliveryDate)
			fmt.Printf("Receiving book on %s\n", deliveryDate[0].D.String())
			return *sEnd
		}
		sEnd := s1.Seller_1_Scatter_Quit()
		fmt.Printf("Book is not affordable (Cost: %d > %d)\n", amtToPay[0], WalletAmount)
		scributil.Debugf("Buyer[%d]: sent Quit.\n", self)
		return *sEnd
	})
	wg.Done()
}

// Buyer1_family2 implements Buyer[1]
func Buyer1_family2(p *NBuyers.NBuyers, K, self int, seller scributil.ClientConn, host string, port int, buyer scributil.ClientConn, buyer2Host string, buyer2Port int, wg *sync.WaitGroup) {
	Buyer1 := p.New_family_2_Buyer_1to1and1toKand1toKsub1_not_2toKandKtoK(K, self)

	scributil.Debugf("[connection] Buyer[%d]: dialling to Seller at %s:%d.\n", self, host, port)
	if err := Buyer1.Seller_1to1_Dial(1, host, port, seller.Dial, seller.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	scributil.Debugf("[connection] Buyer[%d]: dialling to Buyer[%d] at %s:%d.\n", self, self+1, buyer2Host, buyer2Port)
	if err := Buyer1.Buyer_1toKand2toKandKtoK_not_1to1and1toKsub1_Dial(self+1, buyer2Host, buyer2Port, buyer.Dial, buyer.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	scributil.Debugf("[connection] Buyer[%d]: Ready.\n", self)

	Buyer1.Run(func(s *Buyer_1_family2.Init) Buyer_1_family2.End {
		title := []string{"Harry Potter and the Philosopher's Stone"}
		s0 := s.Seller_1_Scatter_Title(title)
		scributil.Debugf("Buyer[%d]: sent %v.\n", self, title)
		quoteTotal := make([]int, 1)
		s1 := s0.Seller_1_Gather_Quote(quoteTotal)
		scributil.Debugf("Buyer[%d]: received %v.\n", self, quoteTotal)
		quoteRemain := []int{quoteTotal[0] - quoteTotal[0]/K}
		sEnd := s1.Buyer_2_Scatter_Quote(quoteRemain)
		scributil.Debugf("Buyer[%d]: sent %v.\n", self, quoteRemain)
		return *sEnd
	})
	wg.Done()
}

// BuyerK_family2 implements Buyer[K].
func BuyerK_family2(p *NBuyers.NBuyers, K, self int, seller scributil.ClientConn, host string, port int, buyer scributil.ServerConn, buyerPort int, wg *sync.WaitGroup) {
	BuyerK := p.New_family_2_Buyer_1toKand2toKandKtoK_not_1to1and1toKsub1(K, self)

	if err := BuyerK.Seller_1to1_Dial(1, host, port, seller.Dial, seller.Formatter()); err != nil {
		log.Fatalf("cannot dial: %v", err)
	}
	ln, err := buyer.Listen(buyerPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer ln.Close()
	scributil.Debugf("[connection] Buyer[%d]: listening for Buyer[%d] at :%d.\n", self, self-1, buyerPort)
	if err := BuyerK.Buyer_1to1and1toKand1toKsub1_not_2toKandKtoK_Accept(self-1, ln, buyer.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	scributil.Debugf("[connection] Buyer[%d]: Ready.\n", self)

	BuyerK.Run(func(s *Buyer_K_family2.Init) Buyer_K_family2.End {
		quoteTotal := make([]int, 1)
		s0 := s.Seller_1_Gather_Quote(quoteTotal)
		scributil.Debugf("Buyer[%d]: received %v.\n", self, quoteTotal)
		amtToPay := make([]int, 1)
		s1 := s0.Buyer_selfsub1_Gather_Quote(amtToPay)
		scributil.Debugf("Buyer[%d]: received %v.\n", self, amtToPay)
		if isAffordable(amtToPay) {
			fmt.Printf("Book is affordable (Cost: %d <= %d)\n", amtToPay[0], WalletAmount)
			s2 := s1.Seller_1_Scatter_OK()
			scributil.Debugf("Buyer[%d]: sent OK.\n", self)
			addr := getAddress()
			s3 := s2.Seller_1_Scatter_Address(addr)
			scributil.Debugf("Buyer[%d]: sent %v.\n", self, addr)
			deliveryDate := make([]message.Date, 1)
			sEnd := s3.Seller_1_Gather_Date(deliveryDate)
			scributil.Debugf("Buyer[%d]: received %v.\n", self, deliveryDate)
			fmt.Printf("Receiving book on %s\n", deliveryDate[0].D.String())
			return *sEnd
		}
		sEnd := s1.Seller_1_Scatter_Quit()
		fmt.Printf("Book is not affordable (Cost: %d > %d)\n", amtToPay[0], WalletAmount)
		scributil.Debugf("Buyer[%d]: sent Quit.\n", self)
		return *sEnd
	})
	wg.Done()
}

func getAddress() []message.Address {
	return []message.Address{message.Address{
		Line1:    "Imperial College London",
		Line2:    "South Kensington",
		Country:  "United Kingdom",
		PostCode: "SW7 2AZ",
	}}
}

func tomorrow() []message.Date {
	return []message.Date{message.Date{
		D: time.Now().Add(24 * time.Hour),
	}}
}

func isAffordable(num []int) bool {
	return num[0] <= WalletAmount
}
