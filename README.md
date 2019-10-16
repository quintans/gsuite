# gsuite
xUnit tests style with setup/teardown and parameterized tests

This project got inspiration from https://github.com/pavlo/gosuite

## Installation
```
go get -u github.com/quintans/gsuite
```

## Quick Start
Fist create a struct that will hold your tests. It must embed `*assert.Assertions`

```go
type Suite struct {
    *assert.Assertions // Required
    // DB connection
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
func (s *Suite) SetupSuite(t *testing.T) {
}

// TearDownSuite is called once after thevery last test in suite runs
func (s *Suite) TearDownSuite(t *testing.T) {
}

// SetUp is called before each test method
func (s *Suite) SetUp(t *testing.T) {
}

// TearDown is called after each test method
func (s *Suite) TearDown(t *testing.T) {
}

// TestXXX is our test
func (s *Suite) TestXXX(t *testing.T) {
}
```

`SetupSuite()` and `TearDownSuite()` share the same `*Suite` instance. This will be a clone of the instance passed with the `gsuite.Run()`.


`Setup()` and `TearDown()` share the same `*Suite` instance. This will be a clone of the instance passed to `SetupSuite()` and `TearDownSuite()`, meaning that they will have the values set in `SetupSuite()` but tests will be independente instances.

> `*Suite` in `Setup()/TearDown()` will have a sub test `*testing.T` instance, derived from the one passed to `TestSuite()`.

## Parameterized Tests

Parameterized tests are similar as table tests in go. We define a set of test cases that will be passed to a test method.

```go
// testCase defines the structure of each test case
type testCase struct {
	in  string
	out string
}

// TableTestUpper output, will feed into TestUpper
func (s *Suite) TableTestUpper(t *testing.T) []testCase {
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
func (s *Suite) TestUpper(t *testing.T, t testCase) {
	s.Equal(t.out, upper(t.in))
}

```

> Every test will have an independent instance of `*Suite`, meaning `*Suite` in `TestUpper()` will have a sub test `*testing.T` instance, derived from a common one. The cloning of `*Suite` happens like this: `parameterized test -(copy of)-> tests -(copy of)-> root`

## Parallelism

Since you have access to `*testing.T` through the method signature we can set it in any way you wish.
