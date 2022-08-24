# gsuite
xUnit tests style with setup/teardown and parameterized tests

This project got inspiration from https://github.com/pavlo/gosuite

Also, this is pretty much what testify suite (https://github.com/stretchr/testify) offers, but with the option of having parallel tests and table driven tests.

The only advantage that this brings to traditional go tests is organization and the call to teardown event if a panic occurs

## Installation
```
go get -u github.com/quintans/gsuite
```

## Quick Start
Fist create a struct that will hold your tests

```go
type MySuite struct {
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
	gsuite.Run(t, &MySuite{})
}
```

To setup the environment before the tests we have the option to write setup and teardown methods.

The available methods are.

```go
// SetUpSuite is called once before the very first test in suite runs
func (s *MySuite) SetupSuite(t *gsuite.T) {
}

// TearDownSuite is called once after thevery last test in suite runs
func (s *MySuite) TearDownSuite(t *gsuite.T) {
}

// SetUp is called before each test method
func (s *MySuite) SetUp(t *gsuite.T) {
}

// TearDown is called after each test method
func (s *MySuite) TearDown(t *gsuite.T) {
}

// TestXXX is our test
func (s *MySuite) TestXXX(t *gsuite.T) {
}
```

The execution calls is depicted in the following tree

```
<run>
 ├─ SetupSuite
 │  ├─ SetUp
 │  ├─ TestXXX
 │  └─ TearDown
 └─ TearDownSuite
```

Each level will have it own copy of `MySuite`, derived from its parent level.

This means that `SetupSuite()` and `TearDownSuite()` share the same `*MySuite` instance and this instance will be a copy of the instance passed to `gsuite.Run()`.

`Setup()`, `TearDown()` and `TestXXX()` share the same `*MySuite` instance and this will be a copy of the instance passed to `SetupSuite()` and `TearDownSuite()`, meaning that they will have the values set in `SetupSuite()` but tests will be independent instances.

Any change made in the `Setup()` will only affect the test being executed, allowing them to run in isolation.

Every test will have its own independent sub test `*testing.T`, passed in the argument `t *gsuite.T` accessible by `t.T()`, meaning `*TestSuite` in `TestUpper()` will have a sub test `*testing.T` instance, derived from a common one.

The **shallow copy** of `*MySuite` happens like this:


## Table tests

```go
func (s *TestSuite) TestUpper(t *gsuite.T) {
	testCases := map[string]struct {
		in  string
		out string
	}{
		"one": {
			in:  "hello",
			out: "HELLO",
		},
		"two": {
			in:  "world",
			out: "WORLD",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *gsuite.T) {
			tableTestCalledTimes++
			t.Equal(tc.out, upper(tc.in))
		})
	}
}

func upper(s string) string {
	return strings.ToUpper(s)
}
```

## Parallelism

Since you have access to `*testing.T`, through `t.T()`, we can set it in any way you wish.
