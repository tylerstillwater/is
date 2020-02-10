# is [![GoDev](https://img.shields.io/static/v1?label=godev&message=documentation&color=informational&style=plastic&&url=https://pkg.go.dev/github.com/tylerb/is@v3.0.0?tab=doc)](https://pkg.go.dev/github.com/tylerb/is/v3@v3.0.0?tab=doc) [![Build Status](https://circleci.com/gh/tylerb/is/tree/v3.svg?style=shield&circle-token=94428439ffc6eda6471dc218471dab20985f444c)](https://circleci.com/gh/tylerb/is/tree/v3)

Is provides a quick, clean and simple framework for writing Go tests.

## Installation

To install, simply execute:

```
go get -u github.com/tylerb/is/v3
```

## Usage

```go
func TestSomething(t *testing.T) {
	assert := is.New(t)

	expected := 10
	actual, _ := awesomeFunction()
	assert.Equal(actual, expected)
}
```

If you'd like a bit more information when a test fails, you may use the `Msg()` method:

```go
func TestSomething(t *testing.T) {
	assert := is.New(t)

	expected := 10
	actual, details := awesomeFunction()
	assert.Msg("actual details: %s", details).Equal(actual, expected)
}
```

By default, any assertion that fails will halt termination of the test. If you would like to run a group of assertions
in a row, you may use the `Lax` method. This is useful for asserting/printing many values at once, so you can correct
all the issues between test runs.

```go
func TestSomething(t *testing.T) {
	assert := is.New(t)

	assert.Lax(func (lax Asserter) {
		lax.Equal(1, 2)
		lax.True(false)
		lax.ShouldPanic(func(){})
	}) 
}
```

If any of the assertions fail inside that function, an additional error will be printed and test execution will halt.
