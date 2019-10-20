# gsuite
xUnit tests style with setup/teardown and parameterized tests

This project got inspiration from https://github.com/pavlo/gosuite

Also, this is pretty much what testify suite (https://github.com/stretchr/testify) offers, but with the option of having parallel tests and table driven tests.

## Installation
```
go get -u github.com/quintans/gsuite
```

## Quick Start
Fist create a struct that will hold your tests. It must embed `*gsuite.Suite`

```go
type TestSuite struct {
    *gsuite.Suite // Required
    // DB connection pool
    // rest/grpc clients
    // temporary directories
    // fakes
    // etc
}
```


and to run it, we use the classic go testing entry point and call our suite

```go
func TestSuite(t *testing.T) {
	gsuite.Run(t, &Suite{})
}
```

To setup the environment before the tests we have the option to write setup and teardown methods.

The available methods are.

```go
// SetUpSuite is called once before the very first test in suite runs
func (s *TestSuite) SetupSuite() {
}

// TearDownSuite is called once after thevery last test in suite runs
func (s *TestSuite) TearDownSuite() {
}

// SetUp is called before each test method
func (s *TestSuite) SetUp() {
}

// TearDown is called after each test method
func (s *TestSuite) TearDown() {
}

// TestXXX is our test
func (s *TestSuite) TestXXX() {
}
```

`SetupSuite()` and `TearDownSuite()` share the same `*TestSuite` instance. This will be a clone of the instance passed with the `gsuite.Run()`.


`Setup()` and `TearDown()` share the same `*TestSuite` instance. This will be a clone of the instance passed to `SetupSuite()` and `TearDownSuite()`, meaning that they will have the values set in `SetupSuite()` but tests will be independente instances.

> `*TestSuite` in `Setup()/TearDown()` will have a sub test `*testing.T` instance, derived from the one passed to `TestSuite()`.

## Parameterized Tests

Parameterized tests are similar as table tests in go. We define a set of test cases that will be passed to a test method.

```go
// testCase defines the structure of each test case
type testCase struct {
	in  string
	out string
}

// TableTestUpper output, will feed into TestUpper
func (s *TestSuite) TableTestUpper() []testCase {
	return []testCase{
		{
			in:  "hello",
			out: "HELLO",
		},
		{
			in:  "world",
			out: "WORLD",
		},
	}
}

// TestUpper will be called with each element from the output slice of TableTestUpper
func (s *TestSuite) TestUpper(t testCase) {
	s.Equal(t.out, upper(t.in))
}

```

Every test will use a **shallow copy** of the initial `*TestSuite` passed as argument to `gsuite.Run()`, so any change made in the `Setup()` will only affect the test being executed, allowing them to run in isolation.

Every test will have its own independent sub test `*testing.T`, accessible by `s.T()`, meaning `*TestSuite` in `TestUpper()` will have a sub test `*testing.T` instance, derived from a common one.

The **shallow copy** of `*TestSuite` happens like this: `run -> test -> table test`


## Parallelism

Since you have access to `*testing.T`, through `s.T()`, we can set it in any way you wish.
