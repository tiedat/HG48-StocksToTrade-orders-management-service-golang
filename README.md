# StocksToTrade order management

A service responsible for the taking and processing of orders.

## Development

[Install](https://golang.org/doc/install) Go 1.16+.

## Testing

Install [Docker](https://www.docker.com) before running tests.

```text
$ go test -v
=== RUN   TestOrdersHandler
=== RUN   TestOrdersHandler/ok
    order_test.go:32: [{"id":1,"email":"foo@example.com","create_at":"2021-07-28T03:52:28.327702Z","updated_at":"2021-07-28T03:52:28.327702Z"},{"id":2,"email":"bar@example.com","create_at":"2021-07-28T03:52:28.327702Z","updated_at":"2021-07-28T03:52:28.327702Z"}]

--- PASS: TestOrdersHandler (0.00s)
    --- PASS: TestOrdersHandler/ok (0.00s)
=== RUN   TestPing
=== RUN   TestPing/ok
--- PASS: TestPing (0.00s)
    --- PASS: TestPing/ok (0.00s)
PASS
ok  	gitlab.com/HG48-StocksToTrade/orders-management-service-golang	4.521s
```
>>>>>>> 02bd9ea (Update README.md information)
