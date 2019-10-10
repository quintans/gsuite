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

type Suite struct {
	*assert.Assertions

	setUpSuiteCalledTimes    int
	tearDownSuiteCalledTimes int
	setUpCalledTimes         int
	tearDownCalledTimes      int
}

func TestCalls(t *testing.T) {
	s := &Suite{}
	gsuite.Run(t, s)

	assert.Equal(t, 1, setUpSuiteCalledTimes)
	assert.Equal(t, 1, tearDownSuiteCalledTimes)
	assert.Equal(t, 3, setUpCalledTimes)
	assert.Equal(t, 3, tearDownCalledTimes)
	assert.Equal(t, 1, testFirstCalledTimes)
	assert.Equal(t, 1, testSecondCalledTimes)
	assert.Equal(t, 2, tableTestCalledTimes)
}

func (s *Suite) SetupSuite(t *testing.T) {
	setUpSuiteCalledTimes++
	s.setUpSuiteCalledTimes++
}

func (s *Suite) TearDownSuite(t *testing.T) {
	tearDownSuiteCalledTimes++
}

func (s *Suite) Setup(t *testing.T) {
	setUpCalledTimes++
	s.setUpCalledTimes++
}

func (s *Suite) TearDown(t *testing.T) {
	tearDownCalledTimes++
	s.tearDownCalledTimes++
}

func (s *Suite) TestFirstTestMethod(t *testing.T) {
	testFirstCalledTimes++
	s.Equal(1, setUpSuiteCalledTimes)
	s.Equal(1, s.setUpSuiteCalledTimes)
	s.Equal(0, tearDownSuiteCalledTimes)
	s.Equal(1, setUpCalledTimes)
	s.Equal(1, s.setUpCalledTimes)
	s.Equal(0, tearDownCalledTimes)
}

func (s *Suite) TestSecondTestMethod(t *testing.T) {
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
func (s *Suite) TableTestThirdTestMethod(t *testing.T) []testCase {
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
func (s *Suite) TestThirdTestMethod(t *testing.T, tc testCase) {
	tableTestCalledTimes++
	s.Equal(tc.out, upper(tc.in))
}

func upper(s string) string {
	return strings.ToUpper(s)
}
