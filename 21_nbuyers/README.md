## nbuyers

The n-buyers protocol describes how N Buyers collaborate to buy from a Seller.
This implementation uses a simple strategy, where each Buyer pays for a equal
proportion of the cost, i.e. Cost/N each, and the final Buyer decides if the
product is *affordable*.

In the example code (shared memory version), the Cost of the product is 100 and
the product is affordable if the final Buyer pays less than 30. Hence with 3
buyers, i.e.

    $ go run main.go -K 3
    Book is not affordable (Cost: 34 > 30)

The cost for each Buyer is 33 (34 for last Buyer) and is not affordable; but
if the number of Buyers is more than or equal to 4, then the project is
affordable, i.e.

    $ go run main.go -K 4
    Book is affordable (Cost: 25 <= 30)
    Receiving book on [date]
