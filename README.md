## DISCLAIMER
Its my first application using go-kit library. Some of the parts of app may look not idiomatic in terms of go-kit.

Another point is: this project was developed in limited time frame. And some of parts not covered well (e.g. not all parts of code has tests).

# LittleBill

This project provides generic Wallet service. It operates by RESTful API.

### Available methods:
- Account creating (with initial balance level)
- List of available accounts in system
- Execute transfers from one account to another (with same currency)
- List executed transfers by account name

Keep in mind that any call may be rejected with some error status. In this scenario caller must repeat same call later.

## How to setup
To launch application on local system you can follow steps:

- `go get github.com/mntor/littlebill`
- `cd "$(go env GOPATH)/src/github.com/mntor/littlebill"`
- `docker build -t littlebill .`
- `docker-compose up`

If no errors have occurred you can make calls to API.

## Examples

(to pretty print output of calls we use [jq](https://stedolan.github.io/jq/) — you can omit '`| jq`' in the end of each call)

Create first account with initial balance of 100.00 USD:
```
➭ curl --silent -X PUT --data '{"name": "acc_1", "balance": 10000, "currency": "USD"}' http://localhost:8080/accounts | jq
```

Create second account with initial balance of 99.99 USD:
```
➭ curl --silent -X PUT --data '{"name": "acc_2", "balance": 9999, "currency": "USD"}' http://localhost:8080/accounts | jq
```

List created accounts:
```
➭ curl --silent http://localhost:8080/accounts | jq
```

Execute transfer of 3.00 USD from acc_1 to acc_2
```
➭ curl --silent -X POST --data '{"account_from": "acc_1", "account_to": "acc_2", "amount": 300}' http://localhost:8080/transfers/execute | jq
```

Check balance changes of accounts:
```
➭ curl --silent http://localhost:8080/accounts | jq
```

View transactions list of acc_1:
```
➭ curl --silent http://localhost:8080/transfers/acc_1 | jq
```

## Thing to improve

- Test coverage
- Prepare benchmarks
- Make rid of retries when storage returns errors
- Use connection pool to storage
- Write OpenAPI documentation
- ...

## Contributing

Before submitting major changes, here are a few guidelines to follow:

1. Check the [open issues][issues] and [pull requests][prs] for existing discussions.
1. Open an [issue][issues] first, to discuss a new feature or enhancement.
1. Write tests, and make sure the test suite passes locally.
1. Open a pull request, and reference the relevant issue(s).
1. After receiving feedback, [squash your commits][squash] and add a [great commit message][message].
1. Have fun!

[issues]: https://github.com/mntor/littlebill/issues
[prs]: https://github.com/mntor/littlebill/pulls
[squash]: http://gitready.com/advanced/2009/02/10/squashing-commits-with-rebase.html
[message]: http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html