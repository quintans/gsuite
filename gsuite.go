package gsuite

import (
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

const (
	testMethodPrefix    = "Test"
	setupMethod         = "Setup"
	setupSuiteMethod    = "SetupSuite"
	tearDownMethod      = "TearDown"
	tearDownSuiteMethod = "TearDownSuite"
	embededType         = "Suite"
)

type SetupSuiter interface {
	SetupSuite(t *T)
}

type TearDownSuiter interface {
	TearDownSuite(t *T)
}

type Setuper interface {
	Setup(t *T)
}

type TearDowner interface {
	TearDown(t *T)
}

type (
	T struct {
		*require.Assertions
		t *testing.T
	}
)

func (t *T) T() *testing.T {
	return t.t
}

func (t *T) Run(name string, test func(t *T)) {
	t.t.Run(name, func(tt *testing.T) {
		test(newT(tt))
	})
}

func Run(tt *testing.T, suite interface{}) {
	verifySuite(tt, suite)
	suiteType := reflect.TypeOf(suite)

	myT := newT(tt)

	clone := shallowCopy(suite)

	// call setup suite
	if s, ok := clone.(SetupSuiter); ok {
		s.SetupSuite(myT)
	}
	if s, ok := clone.(TearDownSuiter); ok {
		defer s.TearDownSuite(myT)
	}

	for i := 0; i < suiteType.NumMethod(); i++ {
		m := suiteType.Method(i)
		if strings.HasPrefix(m.Name, testMethodPrefix) {
			tt.Run(m.Name, func(tt *testing.T) {
				myT := newT(tt)

				subClone := shallowCopy(clone)

				if s, ok := subClone.(Setuper); ok {
					s.Setup(myT)
				}
				if s, ok := subClone.(TearDowner); ok {
					defer s.TearDown(myT)
				}

				m.Func.Call([]reflect.Value{
					reflect.ValueOf(subClone),
					reflect.ValueOf(myT),
				})
			})
		}
	}
}

func shallowCopy(inter interface{}) interface{} {
	nInter := reflect.New(reflect.TypeOf(inter).Elem())

	val := reflect.ValueOf(inter).Elem()
	nVal := nInter.Elem()
	for i := 0; i < val.NumField(); i++ {
		ori := val.Field(i)
		dest := nVal.Field(i)

		if !dest.CanSet() {
			ori = reflect.NewAt(ori.Type(), unsafe.Pointer(ori.UnsafeAddr())).Elem()
			dest = reflect.NewAt(dest.Type(), unsafe.Pointer(dest.UnsafeAddr())).Elem()
		}
		dest.Set(ori)

	}

	return nInter.Interface()
}

func verifySuite(t *testing.T, suite interface{}) {
	st := reflect.TypeOf(suite)

	for i := 0; i < st.NumMethod(); i++ {
		m := st.Method(i)

		if strings.HasPrefix(m.Name, testMethodPrefix) {
			if m.Type.NumIn() != 2 {
				t.Fatalf("Test %s in %s should have the following signature: func(*gsuite.T)", m.Name, st)
			}
		}
	}
}

func newT(t *testing.T) *T {
	return &T{
		Assertions: require.New(t),
		t:          t,
	}
}
