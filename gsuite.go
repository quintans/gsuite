package gsuite

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	tableMethodPrefix   = "Table"
	testMethodPrefix    = "Test"
	setupMethod         = "Setup"
	setupSuiteMethod    = "SetupSuite"
	tearDownMethod      = "TearDown"
	tearDownSuiteMethod = "TearDownSuite"
	embededType         = "Assertions"
)

type SetupSuiter interface {
	SetupSuite(t *testing.T)
}
type TearDownSuiter interface {
	TearDownSuite(t *testing.T)
}

type Setuper interface {
	Setup(t *testing.T)
}
type TearDowner interface {
	TearDown(t *testing.T)
}

type TestSuite struct {
}

func Run(t *testing.T, suite interface{}) {
	verifySuite(t, suite)
	suiteType := reflect.TypeOf(suite)

	clone := cloner(suite)
	setEmbededAssertions(t, clone)

	// call setup suite
	if s, ok := clone.(SetupSuiter); ok {
		s.SetupSuite(t)
	}
	if s, ok := clone.(TearDownSuiter); ok {
		defer s.TearDownSuite(t)
	}

	for i := 0; i < suiteType.NumMethod(); i++ {
		m := suiteType.Method(i)
		if strings.HasPrefix(m.Name, testMethodPrefix) {
			t.Run(m.Name, func(t *testing.T) {
				subClone := cloner(clone)
				setEmbededAssertions(t, subClone)

				if s, ok := subClone.(Setuper); ok {
					s.Setup(t)
				}
				if s, ok := subClone.(TearDowner); ok {
					defer s.TearDown(t)
				}

				if m.Type.NumIn() == 3 {
					// call table tests
					if tm, ok := suiteType.MethodByName(tableMethodPrefix + m.Name); ok {
						out := tm.Func.Call([]reflect.Value{reflect.ValueOf(subClone), reflect.ValueOf(t)})
						s := out[0]
						for i := 0; i < s.Len(); i++ {
							t.Run(strconv.Itoa(i), func(tt *testing.T) {
								setEmbededAssertions(tt, subClone)
								in := []reflect.Value{reflect.ValueOf(subClone), reflect.ValueOf(tt), s.Index(i)}
								m.Func.Call(in)
								setEmbededAssertions(t, subClone)
							})
						}
					}
				} else {
					in := []reflect.Value{reflect.ValueOf(subClone), reflect.ValueOf(t)}
					m.Func.Call(in)
				}
			})
		}
	}
}

func cloner(inter interface{}) interface{} {
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

type template struct {
	*assert.Assertions
}

var assertField, _ = reflect.TypeOf(&template{}).Elem().FieldByName(embededType)

func verifySuite(t *testing.T, suite interface{}) {
	st := reflect.TypeOf(suite)
	embedded, _ := st.Elem().FieldByName(embededType)
	if embedded.Type != assertField.Type {
		t.Fatalf("Struct %v needs to have *assert.Assertions.", st)
	}

	for i := 0; i < st.NumMethod(); i++ {
		m := st.Method(i)
		// verify that every table test method has a corresponding test method with one argument
		if strings.HasPrefix(m.Name, tableMethodPrefix) {
			if m.Type.NumOut() != 1 && m.Type.Kind() != reflect.Slice {
				t.Fatalf("Table test method %s should have only one output of type slice", m.Name)
			}
			test := m.Name[len(tableMethodPrefix):]
			tm, ok := st.MethodByName(test)
			if !ok {
				t.Fatalf("Table test method %s does not have the corresponding test method %s in %v", m.Name, test, st)
			}
			if tm.Type.NumIn() != 3 {
				t.Fatalf("Test method %s should have 2 arguments (testing.T and test data) to be used in conjuction with %s", test, m.Name)
			}
		}

		if strings.HasPrefix(m.Name, testMethodPrefix) {
			if m.Type.NumIn() < 2 || m.Type.NumIn() > 3 {
				t.Fatalf("Test %s in %s should have the following signature: func(*testing.T [, Test data])", m.Name, st)
			} else if m.Type.NumIn() == 3 {
				tableTest := tableMethodPrefix + m.Name
				_, ok := st.MethodByName(tableTest)
				if !ok {
					t.Fatalf("There is no table test method %s for the test method %s in %v", tableTest, m.Name, st)
				}
			}
		}
	}
}

func setEmbededAssertions(t *testing.T, suite interface{}) {
	v := reflect.ValueOf(suite)
	f := v.Elem().FieldByName(embededType)
	f.Set(reflect.ValueOf(assert.New(t)))
}
