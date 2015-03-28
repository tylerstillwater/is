package is

import (
	"fmt"
	"log"
	"testing"
	"time"
)

// Is provides methods that leverage the existing testing capabilities found
// in the Go test framework. The methods provided allow for a more natural,
// efficient and expressive approach to writing tests. The goal is to write
// fewer lines of code while improving communication of intent.
type Is struct {
	fail func(format string, args ...interface{})
}

// New creates a new instance of the Is object and stores a reference to the
// provided testing object.
func New(tb testing.TB) *Is {
	if tb == nil {
		log.Fatalln("You must provide a testing object.")
	}
	v := &Is{
		fail: func(format string, args ...interface{}) {
			fmt.Print(decorate(fmt.Sprintf(format, args...)))

			// If there is an After.Msg, we need to give it time to print before we
			// fail. Is this a terrible approach? Yeah. But I couldn't think of
			// another way to do it and maintain the interface. Am I a horrible
			// person for doing this? Probably.
			go func() {
				time.Sleep(5 * time.Millisecond)
				tb.FailNow()
			}()
		},
	}
	return v
}

// Equal performs a deep compare of the provided objects and fails if they are
// not equal.
//
// Equal does not respect type differences. If the types are different and
// comparable (eg int32 and int64), but the values are the same, the objects
// are considered equal.
func (is *Is) Equal(a interface{}, b interface{}) After {
	result := isEqual(a, b)
	if !result {
		is.fail("expected objects '%s' and '%s' to be equal, but got: %v and %v",
			objectTypeName(a),
			objectTypeName(b), a, b)
	}
	return newAfter(result == true)
}

// NotEqual performs a deep compare of the provided objects and fails if they are
// equal.
//
// NotEqual does not respect type differences. If the types are different and
// comparable (eg int32 and int64), but the values are different, the objects
// are considered not equal.
func (is *Is) NotEqual(a interface{}, b interface{}) After {
	result := isEqual(a, b)
	if result {
		is.fail("expected objects '%s' and '%s' not to be equal",
			objectTypeName(a),
			objectTypeName(b))
	}
	return newAfter(result == false)
}

// Err checks the provided error object to determine if an error is present.
func (is *Is) Err(e error) After {
	result := isNil(e)
	if result {
		is.fail("expected error")
	}
	return newAfter(result == false)
}

// NotErr checks the provided error object to determine if an error is not
// present.
func (is *Is) NotErr(e error) After {
	result := isNil(e)
	if !result {
		is.fail("expected no error, but got: %v", e)
	}
	return newAfter(result == true)
}

// Nil checks the provided object to determine if it is nil.
func (is *Is) Nil(o interface{}) After {
	result := isNil(o)
	if !result {
		is.fail("expected object '%s' to be nil, but got: %v", objectTypeName(o), o)
	}
	return newAfter(result == true)
}

// NotNil checks the provided object to determine if it is not nil.
func (is *Is) NotNil(o interface{}) After {
	result := isNil(o)
	if result {
		is.fail("expected object '%s' not to be nil", objectTypeName(o))
	}
	return newAfter(result == false)
}

// True checks the provided boolean to determine if it is true.
func (is *Is) True(b bool) After {
	result := b == true
	if !result {
		is.fail("expected boolean to be true")
	}
	return newAfter(result == true)
}

// False checks the provided boolean to determine if is false.
func (is *Is) False(b bool) After {
	result := b == false
	if !result {
		is.fail("expected boolean to be false")
	}
	return newAfter(result == true)
}

// Zero checks the provided object to determine if it is the zero value
// for the type of that object. The zero value is the same as what the object
// would contain when initialized but not assigned.
//
// This method, for example, would be used to determine if a string is empty,
// an array is empty or a map is empty. It could also be used to determine if
// a number is 0.
//
// In cases such as slice, map, array and chan, a nil value is treated the
// same as an object with len == 0
func (is *Is) Zero(o interface{}) After {
	result := isZero(o)
	if !result {
		is.fail("expected object '%s' to be zero value, but it was: %v", objectTypeName(o), o)
	}
	return newAfter(result == true)
}

// NotZero checks the provided object to determine if it is not the zero
// value for the type of that object. The zero value is the same as what the
// object would contain when initialized but not assigned.
//
// This method, for example, would be used to determine if a string is not
// empty, an array is not empty or a map is not empty. It could also be used
// to determine if a number is not 0.
//
// In cases such as slice, map, array and chan, a nil value is treated the
// same as an object with len == 0
func (is *Is) NotZero(o interface{}) After {
	result := isZero(o)
	if result {
		is.fail("expected object '%s' not to be zero value", objectTypeName(o))
	}
	return newAfter(result == false)
}

// After is an interface describing methods that may be called as part of a
// fluent interface after the Is methods. For example,
//
// 				is.Equal(res.Code, 200).Msg("Bad response. Body is: %v",res.Body())
//
// checks to see if the code is 200. If that check fails, the message is
// printed to give the programmer more information about why this check
// failed. If the check succeeds, nothing is printed.
type After interface {
	Msg(format string, args ...interface{})
}

type after struct {
	success bool
}

// ensure after implements After interface
var _ After = (*after)(nil)

func (a *after) Msg(format string, args ...interface{}) {
	if !a.success {
		fmt.Printf("\t"+format+"\n", args...)
	}
}

// newAfter is a function variable that points to defaultAfter by default.
// This is overridden in tests to test the after functionality.
var newAfter = defaultAfter

// defaultAfter is used to construct an after object
func defaultAfter(success bool) After {
	return &after{success: success}
}
