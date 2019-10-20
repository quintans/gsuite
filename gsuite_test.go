package gsuite_test

import (
	"strings"
	"testing"

	"github.com/quintans/gsuite"
	"github.com/stretchr/testify/assert"
)

var (
	setUpSuiteCalledTimes    int
	tearDownSuiteCalledTimes int
	setUpCalledTimes         int
	tearDownCalledTimes      int
	testFirstCalledTimes     int
	testSecondCalledTimes    int
	tableTestCalledTimes     int
)

type TestSuite struct {
	*gsuite.Suite

	setUpSuiteCalledTimes    int
	tearDownSuiteCalledTimes int
	setUpCalledTimes         int
	tearDownCalledTimes      int
}

func TestCalls(t *testing.T) {
	s := &TestSuite{}
	gsuite.Run(t, s)

	assert.Equal(t, 1, setUpSuiteCalledTimes)
	assert.Equal(t, 1, tearDownSuiteCalledTimes)
	assert.Equal(t, 3, setUpCalledTimes)
	assert.Equal(t, 3, tearDownCalledTimes)
	assert.Equal(t, 1, testFirstCalledTimes)
	assert.Equal(t, 1, testSecondCalledTimes)
	assert.Equal(t, 2, tableTestCalledTimes)
}

func (s *TestSuite) SetupSuite() {
	setUpSuiteCalledTimes++
	s.setUpSuiteCalledTimes++
}

func (s *TestSuite) TearDownSuite() {
	tearDownSuiteCalledTimes++
}

func (s *TestSuite) Setup() {
	setUpCalledTimes++
	s.setUpCalledTimes++
}

func (s *TestSuite) TearDown() {
	tearDownCalledTimes++
	s.tearDownCalledTimes++
}

func (s *TestSuite) TestFirstTestMethod() {
	testFirstCalledTimes++
	s.Equal(1, setUpSuiteCalledTimes)
	s.Equal(1, s.setUpSuiteCalledTimes)
	s.Equal(0, tearDownSuiteCalledTimes)
	s.Equal(1, setUpCalledTimes)
	s.Equal(1, s.setUpCalledTimes)
	s.Equal(0, tearDownCalledTimes)
}

func (s *TestSuite) TestSecondTestMethod() {
	testSecondCalledTimes++
	s.Equal(1, setUpSuiteCalledTimes)
	s.Equal(1, s.setUpSuiteCalledTimes)
	s.Equal(0, tearDownSuiteCalledTimes)
	s.Equal(2, setUpCalledTimes)
	s.Equal(1, s.setUpCalledTimes)
	s.Equal(1, tearDownCalledTimes)
}

type testCase struct {
	in  string
	out string
}

// TableTestThirdTestMethod output, will feed into TestThirdTestMethod
func (s *TestSuite) TableTestThirdTestMethod() []testCase {
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

// TestThirdTestMethod will be called with each element from the output slice of TableTestThirdTestMethod
func (s *TestSuite) TestThirdTestMethod(tc testCase) {
	tableTestCalledTimes++
	s.Equal(tc.out, upper(tc.in))
}

func upper(s string) string {
	return strings.ToUpper(s)
}
