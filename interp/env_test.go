package interp

import (
	"testing"

	"gitlab.com/rcmonnet/glox/lang"
)

func TestDefineAndGet(t *testing.T) {

	t.Run("get existing value in current environment", func(t *testing.T) {

		defer func() {
			err := recover()
			if err != nil {
				rte := err.(runtimeError)
				t.Fatalf("get returned an error: %s", rte.Error())
			}
		}()

		pi := 3.14
		env := newEnv(nil)
		env.define("pi", pi)
		lookupVal := env.get(newToken("pi"))
		if pi != lookupVal {
			t.Errorf("Expected %f but got %f", pi, lookupVal)
		}
	})

	t.Run("get existing value in ancestor environment", func(t *testing.T) {

		defer func() {
			err := recover()
			if err != nil {
				rte := err.(runtimeError)
				t.Fatalf("get returned an error: %s", rte.Error())
			}
		}()

		pi := 3.14
		env := newEnv(nil)
		env.define("pi", pi)
		env = newEnv(env)
		lookupVal := env.get(newToken("pi"))
		if pi != lookupVal {
			t.Errorf("Expected %f but got %f", pi, lookupVal)
		}
	})

	t.Run("get non existing value", func(t *testing.T) {

		defer func() {
			err := recover()
			if err == nil {
				t.Fatal("Expected get to raise a runtimeError")
			} else {
				_, ok := err.(runtimeError)
				if !ok {
					t.Fatal("Expected get to raise a runtimeError")
				}
			}
		}()

		env := newEnv(nil)
		env.get(newToken("pi"))
	})

}

func TestAssign(t *testing.T) {

	t.Run("assign existing value in current environment", func(t *testing.T) {

		defer func() {
			err := recover()
			if err != nil {
				rte := err.(runtimeError)
				t.Fatalf("assign returned an error: %s", rte.Error())
			}
		}()

		pi := 3.14
		betterPi := 3.14159
		env := newEnv(nil)
		env.define("pi", pi)
		env.assign(newToken("pi"), betterPi)
		lookupVal := env.get(newToken("pi"))
		if betterPi != lookupVal {
			t.Errorf("Expected %f but got %f", pi, lookupVal)
		}
	})

	t.Run("assign existing value in ancestor environment", func(t *testing.T) {

		defer func() {
			err := recover()
			if err != nil {
				rte := err.(runtimeError)
				t.Fatalf("assign returned an error: %s", rte.Error())
			}
		}()

		pi := 3.14
		betterPi := 3.14159
		env := newEnv(nil)
		env.define("pi", pi)
		env = newEnv(env)
		env.assign(newToken("pi"), betterPi)
		lookupVal := env.get(newToken("pi"))
		if betterPi != lookupVal {
			t.Errorf("Expected %f but got %f", pi, lookupVal)
		}
	})

	t.Run("assign non existing value", func(t *testing.T) {

		defer func() {
			err := recover()
			if err == nil {
				t.Fatal("Expected assign to raise a runtimeError")
			} else {
				_, ok := err.(runtimeError)
				if !ok {
					t.Fatal("Expected assign to raise a runtimeError")
				}
			}
		}()

		pi := 3.14
		env := newEnv(nil)
		env.assign(newToken("pi"), pi)
	})

}

func TestDepth(t *testing.T) {

	env := newEnv(nil)

	if env.depth() != 0 {
		t.Fatalf("Expected depth of first environment to be 0")
	}

	env = newEnv(env)

	if env.depth() != 1 {
		t.Fatalf("Expected depth of second environment to be 1")
	}
}

// ------------------
// Helper functions
// ------------------

func newToken(name string) *lang.Token {

	tk := &lang.Token{}
	tk.Type = lang.IdentifierToken
	tk.Lexeme = name
	return tk
}
