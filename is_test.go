package is

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

var numberTypes = []reflect.Type{
	reflect.TypeOf(int(0)),
	reflect.TypeOf(int8(0)),
	reflect.TypeOf(int16(0)),
	reflect.TypeOf(int32(0)),
	reflect.TypeOf(int64(0)),
	reflect.TypeOf(uint(0)),
	reflect.TypeOf(uint8(0)),
	reflect.TypeOf(uint16(0)),
	reflect.TypeOf(uint32(0)),
	reflect.TypeOf(uint64(0)),
	reflect.TypeOf(float32(0)),
	reflect.TypeOf(float64(0)),
}

type testStruct struct {
	v int
}

var tests = []struct {
	a      interface{}
	b      interface{}
	c      interface{}
	d      interface{}
	e      interface{}
	cTypes []reflect.Type
}{
	{
		a:      0,
		b:      0,
		c:      1,
		d:      0,
		e:      1,
		cTypes: numberTypes,
	},
	{
		a: "test",
		b: "test",
		c: "testing",
		d: "",
		e: "testing",
	},
	{
		a: struct{}{},
		b: struct{}{},
		c: struct{ v int }{v: 1},
		d: testStruct{},
		e: testStruct{v: 1},
	},
	{
		a: &struct{}{},
		b: &struct{}{},
		c: &struct{ v int }{v: 1},
		d: &testStruct{},
		e: &testStruct{v: 1},
	},
	{
		a: []int64{0, 1},
		b: []int64{0, 1},
		c: []int64{0, 2},
		d: []int64{},
		e: []int64{0, 2},
	},
	{
		a: map[string]int64{"answer": 42},
		b: map[string]int64{"answer": 42},
		c: map[string]int64{"answer": 43},
		d: map[string]int64{},
		e: map[string]int64{"answer": 42},
	},
	{
		a: true,
		b: true,
		c: false,
		d: false,
		e: true,
	},
}

func Test(t *testing.T) {
	assert := New(t)

	for i, test := range tests {
		for _, cType := range test.cTypes {
			fail = func(is *asserter, format string, args ...interface{}) {
				fmt.Print(fmt.Sprintf(fmt.Sprintf("(test #%d) - ", i)+format, args...))
				t.FailNow()
			}
			assert.Equal(test.a, reflect.ValueOf(test.b).Convert(cType).Interface())
		}
		assert.Equal(test.a, test.b)
	}

	for i, test := range tests {
		for _, cType := range test.cTypes {
			fail = func(is *asserter, format string, args ...interface{}) {
				fmt.Print(fmt.Sprintf(fmt.Sprintf("(test #%d) - ", i)+format, args...))
				t.FailNow()
			}
			assert.NotEqual(test.a, reflect.ValueOf(test.c).Convert(cType).Interface())
		}
		assert.NotEqual(test.a, test.c)
	}

	for i, test := range tests {
		for _, cType := range test.cTypes {
			fail = func(is *asserter, format string, args ...interface{}) {
				fmt.Print(fmt.Sprintf(fmt.Sprintf("(test #%d) - ", i)+format, args...))
				t.FailNow()
			}
			assert.Zero(reflect.ValueOf(test.d).Convert(cType).Interface())
		}
		assert.Zero(test.d)
	}

	for i, test := range tests {
		for _, cType := range test.cTypes {
			fail = func(is *asserter, format string, args ...interface{}) {
				fmt.Print(fmt.Sprintf(fmt.Sprintf("(test #%d) - ", i)+format, args...))
				t.FailNow()
			}
			assert.NotZero(reflect.ValueOf(test.e).Convert(cType).Interface())
		}
		assert.NotZero(test.e)
	}

	fail = func(is *asserter, format string, args ...interface{}) {
		fmt.Print(fmt.Sprintf(format, args...))
		t.FailNow()
	}
	assert.Nil(nil)
	assert.NotNil(&testStruct{v: 1})
	assert.Err(errors.New("error"))
	assert.NotErr(nil)
	assert.True(true)
	assert.False(false)
	assert.Zero(nil)
	assert.Nil((*testStruct)(nil))
	assert.OneOf(1, 2, 3, 1)
	assert.NotOneOf(1, 2, 3)

	lens := []interface{}{
		[]int{1, 2, 3},
		[3]int{1, 2, 3},
		map[int]int{1: 1, 2: 2, 3: 3},
	}
	for _, l := range lens {
		assert.Len(l, 3)
	}

	fail = func(is *asserter, format string, args ...interface{}) {}
	assert.Equal((*testStruct)(nil), &testStruct{})
	assert.Equal(&testStruct{}, (*testStruct)(nil))
	assert.Equal((*testStruct)(nil), (*testStruct)(nil))

	fail = func(is *asserter, format string, args ...interface{}) {
		fmt.Print(fmt.Sprintf(format, args...))
		t.FailNow()
	}
	assert.ShouldPanic(func() {
		panic("The sky is falling!")
	})
}

func TestMsg(t *testing.T) {
	assert := New(t)

	assert = assert.Msg("something %s", "else")
	if assert.(*asserter).failFormat != "something %s" {
		t.Fatal("failFormat not set")
	}
	if assert.(*asserter).failArgs[0].(string) != "else" {
		t.Fatal("failArgs not set")
	}

	assert = assert.AddMsg("another %s %s", "couple", "things")
	if assert.(*asserter).failFormat != "something %s - another %s %s" {
		t.Fatal("failFormat not set")
	}
	if assert.(*asserter).failArgs[0].(string) != "else" {
		t.Fatal("failArgs not set")
	}
	if assert.(*asserter).failArgs[1].(string) != "couple" {
		t.Fatal("failArgs not set")
	}
	if assert.(*asserter).failArgs[2].(string) != "things" {
		t.Fatal("failArgs not set")
	}
}

func TestLaxNoFailure(t *testing.T) {
	assert := New(t)

	hit := 0

	fail = func(is *asserter, format string, args ...interface{}) {
		hit++
	}

	assert.Lax(func(lax Asserter) {
		lax.Equal(1, 1)
	})

	fail = failDefault

	assert.Equal(hit, 0)
}

func TestLaxFailure(t *testing.T) {
	assert := New(t)

	hitLax := 0
	hitStrict := 0

	fail = func(is *asserter, format string, args ...interface{}) {
		if is.strict {
			hitStrict++
			return
		}
		is.failed = true
		hitLax++
	}

	assert.Lax(func(lax Asserter) {
		lax.Equal(1, 2)
	})

	fail = failDefault

	assert.Equal(hitLax, 1)
	assert.Equal(hitStrict, 1)
}

func TestOneOf(t *testing.T) {
	assert := New(t)

	hit := 0
	fail = func(is *asserter, format string, args ...interface{}) {
		hit++
	}
	assert.OneOf(2, 1, 2, 3)
	assert.OneOf(4, 1, 2, 3)
	assert.NotOneOf(2, 1, 2, 3)
	assert.NotOneOf(4, 1, 2, 3)

	fail = failDefault
	assert.Equal(hit, 2)
}

func TestFailures(t *testing.T) {
	assert := New(t)

	hit := 0
	fail = func(is *asserter, format string, args ...interface{}) {
		hit++
	}

	assert.NotEqual(1, 1)
	assert.Err(nil)
	assert.NotErr(errors.New("error"))
	assert.Nil(&hit)
	assert.NotNil(nil)
	assert.True(false)
	assert.False(true)
	assert.Zero(1)
	assert.NotZero(0)
	assert.Len([]int{}, 1)
	assert.Len(nil, 1)
	assert.ShouldPanic(func() {})

	fail = failDefault
	assert.Equal(hit, 12)
}

func TestWaitForTrue(t *testing.T) {
	assert := New(t)

	hit := 0
	fail = func(is *asserter, format string, args ...interface{}) {
		hit++
	}

	assert.WaitForTrue(200*time.Millisecond, func() bool {
		return false
	})
	assert.Equal(hit, 1)

	assert.WaitForTrue(200*time.Millisecond, func() bool {
		return true
	})
	assert.Equal(hit, 1)
}

type equaler struct {
	equal  bool
	called bool
}

func (e *equaler) Equal(in interface{}) bool {
	e.called = true
	v, ok := in.(*equaler)
	if !ok {
		return false
	}
	return e.equal == v.equal
}

func TestEqualer(t *testing.T) {
	assert := New(t)

	hit := 0
	fail = func(is *asserter, format string, args ...interface{}) {
		hit++
	}

	a := &equaler{equal: true}
	b := &equaler{}

	assert.Equal(a, b)
	if !a.called {
		t.Fatalf("a.Equal should have been called")
	}

	assert.Equal(b, a)
	if !b.called {
		t.Fatalf("b.Equal should have been called")
	}

	if hit != 2 {
		t.Fatalf("fail func should have been called 2 times, but was called %d times", hit)
	}

	a.called = false
	b.called = false
	b.equal = true
	hit = 0

	assert.NotEqual(a, b)
	if !a.called {
		t.Fatalf("a.Equal should have been called")
	}

	assert.NotEqual(b, a)
	if !b.called {
		t.Fatalf("b.Equal should have been called")
	}

	if hit != 2 {
		t.Fatalf("fail func should have been called 2 times, but was called %d times", hit)
	}
}
