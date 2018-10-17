//go:generate scribblec-param.sh ../../QuoteRequest.scr -d ../../ -param WebService github.com/nickng/scribble-go-examples/11_quote-request/QuoteRequest -param-api Buyer -param-api Supplier -param-api Manufacturer

package quotereq

import (
	"encoding/gob"
	"log"
	"sync"

	"github.com/nickng/scribble-go-examples/11_quote-request/QuoteRequest/WebService"
	"github.com/nickng/scribble-go-examples/11_quote-request/QuoteRequest/WebService/Buyer_1to1"
	"github.com/nickng/scribble-go-examples/11_quote-request/QuoteRequest/WebService/Manufacturer_1toM"
	"github.com/nickng/scribble-go-examples/11_quote-request/QuoteRequest/WebService/Supplier_1to1and1toS_not_2toS"
	"github.com/nickng/scribble-go-examples/11_quote-request/QuoteRequest/WebService/Supplier_1toSand2toS_not_1to1"
	"github.com/nickng/scribble-go-examples/11_quote-request/message"
	"github.com/nickng/scribble-go-examples/scributil"
)

func init() {
	gob.Register(new(message.Quote))
	gob.Register(new(message.QuoteReq))
}

// Buyer implements Buyer[1].
func Buyer(p *WebService.WebService, S, self int, supp scributil.ClientConn, host string, basePort int, wg *sync.WaitGroup) {
	Buyer := p.New_Buyer_1to1(S, self)

	for s := 1; s <= S; s++ {
		scributil.Debugf("[connection] Buyer: dialling to Supplier[%d] at %s:%d.\n", s, host, basePort+s)
		if s == 1 {
			// Supplier[1]
			if err := Buyer.Supplier_1to1and1toS_not_2toS_Dial(s, host, basePort+s, supp.Dial, supp.Formatter()); err != nil {
				log.Fatalf("cannot dial: %v", err)
			}
		} else {
			// Supplier[2..S]
			if err := Buyer.Supplier_1toSand2toS_not_1to1_Dial(s, host, basePort+s, supp.Dial, supp.Formatter()); err != nil {
				log.Fatalf("cannot dial: %v", err)
			}
		}
	}
	scributil.Debugf("Buyer: Ready.\n")

	Buyer.Run(func(s *Buyer_1to1.Init) Buyer_1to1.End {
		var qr []message.QuoteReq
		for i := 1; i <= S; i++ {
			qr = append(qr, message.QuoteReq{ItemID: 10})
		}
		s0 := s.Supplier_1toS_Scatter_QuoteReq(qr)
		scributil.Debugf("Buyer[%d]: sent %s\n", self, qr)
		for {
			quote := make([]message.Quote, S)
			s1 := s0.Supplier_1toS_Gather_Reply(quote)
			scributil.Debugf("Buyer[%d]: received %s\n", self, quote)
			if accept(quote) {
				sEnd := s1.Supplier_1toS_Scatter_Order(quote)
				scributil.Debugf("Buyer[%d]: sent order %s\n", self, quote)
				return *sEnd
			}
			s2 := s1.Supplier_1toS_Scatter_Modify(quote)
			scributil.Debugf("Buyer[%d]: sent modify %s\n", self, quote)
			switch s3 := s2.Supplier_1_Branch().(type) {
			case *Buyer_1to1.Confirm:
				scributil.Debugf("Buyer[%d]: received confifm", self)
				sEnd := s3.Recv_Confirm()
				return *sEnd
			case *Buyer_1to1.Modify:
				var q message.Quote
				s1 = s3.Recv_Modify(&q)
				scributil.Debugf("Buyer[%d]: received modify", self, q)
			case *Buyer_1to1.Reject:
				sEnd := s3.Recv_Reject()
				scributil.Debugf("Buyer[%d]: received reject\n", self)
				return *sEnd
			case *Buyer_1to1.Renegotiate:
				s0 = s3.Recv_Renegotiate()
				scributil.Debugf("Buyer[%d]: received renegotiate\n", self)
			}
		}
	})
	wg.Done()
}

// Supplier1 implements Supplier[1].
func Supplier1(p *WebService.WebService, M, S, self int, buy scributil.ServerConn, buyPort int, supp scributil.ClientConn, suppHost string, suppBasePort int, manu scributil.ClientConn, manuHost string, manuBasePort int, wg *sync.WaitGroup) {
	Supplier1 := p.New_Supplier_1to1and1toS_not_2toS(M, S, self)

	// First connect to Manufacturers.
	for m := 1; m <= M; m++ {
		scributil.Debugf("[connection] Supplier[%d]: dialling to Manufacturer[%d] at %s:%d.\n", self, m, manuHost, manuBasePort+S*(m-1))
		// Manufacturer[1..M]
		if err := Supplier1.Manufacturer_1toM_Dial(m, manuHost, manuBasePort+S*(m-1), manu.Dial, manu.Formatter()); err != nil {
			log.Fatalf("cannot dial: %v", err)
		}
	}
	// Dial Supplier broadcast.
	for s := 2; s <= S; s++ {
		scributil.Debugf("[connection] Supplier[%d]: dialling to Supplier[%d] at :%d.\n", self, s, suppBasePort+s)
		if err := Supplier1.Supplier_1toSand2toS_not_1to1_Dial(s, suppHost, suppBasePort+s, supp.Dial, supp.Formatter()); err != nil {
			log.Fatalf("cannot accept: %v", err)
		}
	}
	// Listen to Buyer.
	lnBuy, err := buy.Listen(buyPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer lnBuy.Close()
	scributil.Debugf("[connection] Supplier[%d]: listening for Buyer at :%d.\n", self, buyPort)
	if err := Supplier1.Buyer_1to1_Accept(1, lnBuy, buy.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	scributil.Debugf("Supplier[%d]: Ready.\n", self)

	Supplier1.Run(func(s *Supplier_1to1and1toS_not_2toS.Init) Supplier_1to1and1toS_not_2toS.End {
		qr := make([]message.QuoteReq, 1)
		s0 := s.Buyer_1_Gather_QuoteReq(qr)
		scributil.Debugf("Supplier[%d]: received %s\n", self, qr)
		var qrs []message.QuoteReq
		for {
			for i := 1; i <= M; i++ {
				qrs = append(qrs, message.QuoteReq{ItemID: qr[0].ItemID})
			}
			s1 := s0.Manufacturer_1toM_Scatter_QuoteReq(qrs)
			scributil.Debugf("Supplier[%d]: sent %s\n", self, qrs)
			quote := make([]message.Quote, M)
			s2 := s1.Manufacturer_1toM_Gather_Reply(quote)
			scributil.Debugf("Supplier[%d]: received %s\n", self, quote)
			s3 := s2.Buyer_1_Scatter_Reply(quote)
			scributil.Debugf("Supplier[%d]: sent %s\n", self, quote)
			switch s4 := s3.Buyer_1_Branch().(type) {
			case *Supplier_1to1and1toS_not_2toS.Modify_Supplier_State7:
				var q message.Quote
				s5 := s4.Recv_Modify(&q)
				scributil.Debugf("Supplier[%d]: received modify %s\n", self, q)
				if accept([]message.Quote{q}) {
					s6 := s5.Buyer_1_Scatter_Confirm()
					scributil.Debugf("Supplier[%d]: sent confirm\n", self)
					s7 := s6.Manufacturer_1toM_Scatter_Finish()
					scributil.Debugf("Supplier[%d]: sent finish\n", self)
					sEnd := s7.Supplier_2toS_Scatter_Finish()
					scributil.Debugf("Supplier[%d]: sent finish\n", self)
					return *sEnd
				}
				if withinRange() {
					s6 := s5.Buyer_1_Scatter_Modify([]message.Quote{message.Quote{ItemID: q.ItemID, Quote: q.Quote}})
					scributil.Debugf("Supplier[%d]: sent modify", self)
					s7 := s6.Manufacturer_1toM_Scatter_Finish()
					scributil.Debugf("Supplier[%d]: sent finish\n", self)
					var quotes []message.Quote
					for i := 2; i <= S; i++ {
						quotes = append(quotes, q)
					}
					s3 = s7.Supplier_2toS_Scatter_Modify(quotes)
					scributil.Debugf("Supplier[%d]: sent modify %s\n", self, quotes)
				} else {
					if negotiatePossible() {
						s6 := s5.Buyer_1_Scatter_Renegotiate()
						scributil.Debugf("Supplier[%d]: sent renegotiate\n", self)
						s7 := s6.Manufacturer_1toM_Scatter_Renegotiate()
						scributil.Debugf("Supplier[%d]: sent renegotiate\n", self)
						s0 = s7.Supplier_2toS_Scatter_Renegotiate()
						scributil.Debugf("Supplier[%d]: sent renegotiate\n", self)
					} else {
						s6 := s5.Buyer_1_Scatter_Reject()
						scributil.Debugf("Supplier[%d]: sent reject\n", self)
						s7 := s6.Manufacturer_1toM_Scatter_Finish()
						scributil.Debugf("Supplier[%d]: sent finish\n", self)
						sEnd := s7.Supplier_2toS_Scatter_Finish()
						scributil.Debugf("Supplier[%d]: sent finish\n", self)
						return *sEnd
					}
				}
			case *Supplier_1to1and1toS_not_2toS.Order:
				var o message.Quote
				s5 := s4.Recv_Order(&o)
				scributil.Debugf("Supplier[%d]: received order %s\n", self, o)
				sEnd := s5.Manufacturer_1toM_Scatter_Finish()
				scributil.Debugf("Supplier[%d]: sent finish\n", self)
				return *sEnd
			}
		}
	})
	wg.Done()
}

// Supplier2toS implements Supplier[2..S].
func Supplier2toS(p *WebService.WebService, M, S, self int, buy scributil.ServerConn, buyPort int, supp scributil.ServerConn, suppPort int, manu scributil.ClientConn, manuHost string, manuBasePort int, wg *sync.WaitGroup) {
	Supplier2toS := p.New_Supplier_1toSand2toS_not_1to1(M, S, self)

	// First connect to Manufacturers.
	for m := 1; m <= M; m++ {
		scributil.Debugf("[connection] Supplier[%d]: dialling to Manufacturer[%d] at %s:%d.\n", self, m, manuHost, manuBasePort+S*(m-1))
		// Manufacturer[1..M]
		if err := Supplier2toS.Manufacturer_1toM_Dial(m, manuHost, manuBasePort+S*(m-1), manu.Dial, manu.Formatter()); err != nil {
			log.Fatalf("cannot dial: %v", err)
		}
	}
	// Listen to Supplier broadcast.
	lnSupp, err := supp.Listen(suppPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer lnSupp.Close()
	scributil.Debugf("[connection] Supplier[%d]: listening for Supplier[%d] at :%d.\n", self, 1, suppPort)
	if err := Supplier2toS.Supplier_1to1and1toS_not_2toS_Accept(1, lnSupp, supp.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	// Listen to Buyer.
	lnBuy, err := buy.Listen(buyPort)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	defer lnBuy.Close()
	scributil.Debugf("[connection] Supplier[%d]: listening for Buyer at :%d.\n", self, buyPort)
	if err := Supplier2toS.Buyer_1to1_Accept(1, lnBuy, buy.Formatter()); err != nil {
		log.Fatalf("cannot accept: %v", err)
	}
	scributil.Debugf("Supplier[%d]: Ready.\n", self)

	Supplier2toS.Run(func(s *Supplier_1toSand2toS_not_1to1.Init) Supplier_1toSand2toS_not_1to1.End {
		qr := make([]message.QuoteReq, 1)
		s0 := s.Buyer_1_Gather_QuoteReq(qr)
		scributil.Debugf("Supplier[%d]: received %s\n", self, qr)
		var qrs []message.QuoteReq
		for {
			for i := 1; i <= M; i++ {
				qrs = append(qrs, message.QuoteReq{ItemID: qr[0].ItemID})
			}
			s1 := s0.Manufacturer_1toM_Scatter_QuoteReq(qrs)
			scributil.Debugf("Supplier[%d]: sent %s\n", self, qrs)
			quote := make([]message.Quote, M)
			s2 := s1.Manufacturer_1toM_Gather_Reply(quote)
			s3 := s2.Buyer_1_Scatter_Reply(quote)
			scributil.Debugf("Supplier[%d]: sent %s\n", self, quote)
			switch s4 := s3.Buyer_1_Branch().(type) {
			case *Supplier_1toSand2toS_not_1to1.Order_Supplier_State5:
				var o message.Quote
				sEnd := s4.Recv_Order(&o)
				scributil.Debugf("Supplier[%d]: received order %s\n", self, o)
				return *sEnd
			case *Supplier_1toSand2toS_not_1to1.Modify_Supplier_State5:
				var q message.Quote
				s5 := s4.Recv_Modify(&q)
				scributil.Debugf("Supplier[%d]: received modify %s\n", self, q)
				switch s6 := s5.Supplier_1_Branch().(type) {
				case *Supplier_1toSand2toS_not_1to1.Modify_Supplier_State6:
					var q message.Quote
					s3 = s6.Recv_Modify(&q)
					scributil.Debugf("Supplier[%d]: received modify %s\n", self, q)
				case *Supplier_1toSand2toS_not_1to1.Finish_Supplier_State6:
					sEnd := s6.Recv_Finish()
					scributil.Debugf("Supplier[%d]: received finish\n", self)
					return *sEnd
				case *Supplier_1toSand2toS_not_1to1.Renegotiate_Supplier_State6:
					s0 = s6.Recv_Renegotiate()
					scributil.Debugf("Supplier[%d]: received renegotiate\n", self)
				}
			}
		}
	})
	wg.Done()
}

// Manufacturer implements Manufacturer[1..M].
func Manufacturer(p *WebService.WebService, M, S, self int, supp scributil.ServerConn, suppBasePort int, wg *sync.WaitGroup) {
	Manufacturer := p.New_Manufacturer_1toM(M, S, self)

	wgSvr := new(sync.WaitGroup)
	wgSvr.Add(S)
	for s := 1; s <= S; s++ {
		go func(s int) {
			ln, err := supp.Listen(suppBasePort + s)
			if err != nil {
				log.Fatalf("cannot listen: %v", err)
			}
			defer ln.Close()
			scributil.Debugf("[connection] Manufacturer[%d]: listening for Supplier[%d] at :%d.\n", self, s, suppBasePort+s)
			if s == 1 {
				// Supplier[1]
				if err := Manufacturer.Supplier_1to1and1toS_not_2toS_Accept(s, ln, supp.Formatter()); err != nil {
					log.Fatalf("cannot accept: %v", err)
				}
			} else {
				// Supplier[2..S]
				if err := Manufacturer.Supplier_1toSand2toS_not_1to1_Accept(s, ln, supp.Formatter()); err != nil {
					log.Fatalf("cannot accept: %v", err)
				}
			}
			wgSvr.Done()
		}(s)
	}
	wgSvr.Wait()
	scributil.Debugf("Manufacturer[%d]: Ready.\n", self)

	Manufacturer.Run(func(s *Manufacturer_1toM.Init) Manufacturer_1toM.End {
		for {
			qr := make([]message.QuoteReq, S)
			s0 := s.Supplier_1toS_Gather_QuoteReq(qr)
			scributil.Debugf("Manufacturer[%d]: received %s\n", self, qr)
			var quote []message.Quote
			for i := 1; i <= S; i++ {
				quote = append(quote, message.Quote{ItemID: qr[i-1].ItemID, Quote: self})
			}
			s1 := s0.Supplier_1toS_Scatter_Reply(quote)
			scributil.Debugf("Manufacturer[%d]: sent %s\n", self, quote)
			switch s2 := s1.Supplier_1_Branch().(type) {
			case *Manufacturer_1toM.Finish:
				sEnd := s2.Recv_Finish()
				scributil.Debugf("Manufacturer[%d]: received finish\n", self)
				return *sEnd
			case *Manufacturer_1toM.Renegotiate_Manufacturer_State3:
				s = s2.Recv_Renegotiate()
				scributil.Debugf("Manufacturer[%d]: received renegotiate\n", self)
			}
		}
	})
	wg.Done()
}

func negotiatePossible() bool {
	return false
}

func accept(quote []message.Quote) bool {
	return true
}

func withinRange() bool {
	return false
}
