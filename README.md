# is [![GoDoc](https://godoc.org/github.com/tylerb/is?status.png)](http://godoc.org/github.com/tylerb/is) [![wercker status](https://app.wercker.com/status/ca19d08c7449039f40f9aec9418e4155/s "wercker status")](https://app.wercker.com/project/bykey/ca19d08c7449039f40f9aec9418e4155)

Is provides a quick, clean and simple framework for writing Go tests.

## Installation

To install, simply execute:

```
go get gopkg.in/tylerb/is.v1
```

I am using [gopkg.in](http://http://labix.org/gopkg.in) to control releases.

## Usage

Using `Is` is simple:

```go
func TestSomething(t *testing.T) {
	is := is.New(t)

	expected := 10
	result, _ := awesomeFunction()
	is.Equal(expected,result)
}
```

If you'd like a bit more information when a test fails, you may use the `Msg()` method:

```go
func TestSomething(t *testing.T) {
	is := is.New(t)

	expected := 10
	result, details := awesomeFunction()
	is.Msg("result details: %s", details).Equal(expected,result)
}
```

## Contributing

If you would like to contribute, please:

1. Create a GitHub issue regarding the contribution. Features and bugs should be discussed beforehand.
2. Fork the repository.
3. Create a pull request with your solution.

All pull requests should:

1. Pass [gometalinter -t .](https://github.com/alecthomas/gometalinter) with no warnings.
2. Be `go fmt` formatted.
