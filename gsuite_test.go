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

func (s *TestSuite) SetupSuite(t *gsuite.T) {
	setUpSuiteCalledTimes++
	s.setUpSuiteCalledTimes++
}

func (s *TestSuite) TearDownSuite(t *gsuite.T) {
	tearDownSuiteCalledTimes++
}

func (s *TestSuite) Setup(t *gsuite.T) {
	setUpCalledTimes++
	s.setUpCalledTimes++
}

func (s *TestSuite) TearDown(t *gsuite.T) {
	tearDownCalledTimes++
	s.tearDownCalledTimes++
}

func (s *TestSuite) TestFirstTestMethod(t *gsuite.T) {
	testFirstCalledTimes++
	t.Equal(1, setUpSuiteCalledTimes)
	t.Equal(1, s.setUpSuiteCalledTimes)
	t.Equal(0, tearDownSuiteCalledTimes)
	t.Equal(1, setUpCalledTimes)
	t.Equal(1, s.setUpCalledTimes)
	t.Equal(0, tearDownCalledTimes)
}

func (s *TestSuite) TestSecondTestMethod(t *gsuite.T) {
	testSecondCalledTimes++
	t.Equal(1, setUpSuiteCalledTimes)
	t.Equal(1, s.setUpSuiteCalledTimes)
	t.Equal(0, tearDownSuiteCalledTimes)
	t.Equal(2, setUpCalledTimes)
	t.Equal(1, s.setUpCalledTimes)
	t.Equal(1, tearDownCalledTimes)
}

// TestThirdTestMethod will be called with each element from the output slice of TableTestThirdTestMethod
func (s *TestSuite) TestThirdTestMethod(t *gsuite.T) {
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
