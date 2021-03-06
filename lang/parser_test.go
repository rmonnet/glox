package lang

import (
	"strings"
	"testing"
)

func TestExpr(t *testing.T) {

	t.Run("primary nodes", func(t *testing.T) {
		script := `
			123;
			"hello";
			true;
			false;
			nil;
			(123);
			aVariableName_01;`
		expect := []string{
			"123",
			"\"hello\"",
			"true",
			"false",
			"nil",
			"(group 123)",
			"(aVariableName_01)"}
		matchAST(t, expect, script)
	})

	t.Run("unary operators", func(t *testing.T) {
		script := `
			- 123.45;
			! false;`
		expect := []string{
			"(- 123.45)",
			"(! false)"}
		matchAST(t, expect, script)
	})

	t.Run("binary operators", func(t *testing.T) {
		script := `
			123 + 456;
			123.9 - 456.9;
			123 / 456;
			123 * 456;`
		expect := []string{
			"(+ 123 456)",
			"(- 123.9 456.9)",
			"(/ 123 456)",
			"(* 123 456)"}
		matchAST(t, expect, script)
	})

	t.Run("logical operators", func(t *testing.T) {
		script := `
			-1 < 2;
			1 <= -2;
			1 > 2;
			-1 >= -2;
			1 == 2;
			"a" != "b";
			true and false;
			true or false;
			false or (1 <= 2);`
		expect := []string{
			"(< (- 1) 2)",
			"(<= 1 (- 2))",
			"(> 1 2)",
			"(>= (- 1) (- 2))",
			"(== 1 2)",
			"(!= \"a\" \"b\")",
			"(and true false)",
			"(or true false)",
			"(or false (group (<= 1 2)))"}
		matchAST(t, expect, script)
	})

	t.Run("assigment", func(t *testing.T) {
		script := `
			myVar = 123.456;
			myVar = "blue" + " " + "violet";`
		expect := []string{
			"(assign myVar 123.456)",
			"(assign myVar (+ (+ \"blue\" \" \") \"violet\"))"}
		matchAST(t, expect, script)
	})

	t.Run("Call", func(t *testing.T) {
		script := `
			clock();
			add(12, 34);
			echo("hello");`
		expect := []string{
			"(call (clock) (args))",
			"(call (add) (args 12 34))",
			"(call (echo) (args \"hello\"))"}
		matchAST(t, expect, script)
	})

	t.Run("Get", func(t *testing.T) {
		script := `
			cake.flavor;
			Cake("french").flavor;`
		expect := []string{
			"(get (cake) flavor)",
			"(get (call (Cake) (args \"french\")) flavor)"}
		matchAST(t, expect, script)
	})

	t.Run("Set", func(t *testing.T) {
		script := `
			cake.flavor = "vanilla";
			Cake().flavor = "vanilla";`
		expect := []string{
			"(set (cake) flavor \"vanilla\")",
			"(set (call (Cake) (args)) flavor \"vanilla\")"}
		matchAST(t, expect, script)
	})

	t.Run("block", func(t *testing.T) {
		script := `
			{
				print 123;
				{
					a = 3;
				}
			}`
		expect := []string{
			"(block (print 123) (block (assign a 3)))"}
		matchAST(t, expect, script)
	})

	t.Run("fun", func(t *testing.T) {
		script := `
			fun square(x) { return x * x; }
			fun echo(text) { print text; }
		 	fun triple(x) {
		 		var dbl = double(x);
		 		return x * dbl;
		 	}`
		expect := []string{
			"(fun square (params x) (return (* (x) (x))))",
			"(fun echo (params text) (print (text)))",
			"(fun triple (params x) " +
				"(var dbl (call (double) (args (x)))) " +
				"(return (* (x) (dbl))))"}
		matchAST(t, expect, script)

	})

	t.Run("if", func(t *testing.T) {
		script := `
			if (x > 34) {
				print "x greater than 34";
			} else {
				print "x less than 34";
			}
			if (! morning) print "hi";`
		expect := []string{
			"(if (> (x) 34) (block (print \"x greater than 34\")) " +
				"(block (print \"x less than 34\")))",
			"(if (! (morning)) (print \"hi\"))"}
		matchAST(t, expect, script)
	})

	t.Run("return", func(t *testing.T) {
		script := `
			fun yesterday() { return clock() - 24*3600; }
			fun doNothing() { return; }`
		expect := []string{
			"(fun yesterday (params) " +
				"(return (- (call (clock) (args)) (* 24 3600))))",
			"(fun doNothing (params) (return))"}
		matchAST(t, expect, script)
	})

	t.Run("var declaration", func(t *testing.T) {
		script := `
			var a = 123;
			var a_b = true or 3 < 4;
			var c;`
		expect := []string{
			"(var a 123)",
			"(var a_b (or true (< 3 4)))",
			"(var c)"}
		matchAST(t, expect, script)
	})

	t.Run("while", func(t *testing.T) {
		script := `
			while (i < 10) {
				i = i + 2;
			}
			while (i < 10) i = i + 2;`
		expect := []string{
			"(while (< (i) 10) (block (assign i (+ (i) 2))))",
			"(while (< (i) 10) (assign i (+ (i) 2)))"}
		matchAST(t, expect, script)
	})

	t.Run("for", func(t *testing.T) {
		script := `
			for (i = 0; i < 5; i = i + 1) {
				print i;
			}
			for (i = 0; i < 5; i = i + 1) print i;
			for (var i = 0; i < 5; i = i + 1) print i;
			for (; i < 5; i = i + 1) print i;
			for (; i < 5;) print i;
			for (;;) print i;`
		expect := []string{
			"(block (assign i 0) (while (< (i) 5) (block " +
				"(block (print (i))) (assign i (+ (i) 1)))))",
			"(block (assign i 0) (while (< (i) 5) (block " +
				"(print (i)) (assign i (+ (i) 1)))))",
			"(block (var i 0) (while (< (i) 5) (block " +
				"(print (i)) (assign i (+ (i) 1)))))",
			"(while (< (i) 5) (block " +
				"(print (i)) (assign i (+ (i) 1))))",
			"(while (< (i) 5) (print (i)))",
			"(while true (print (i)))"}
		matchAST(t, expect, script)
	})

	t.Run("class", func(t *testing.T) {
		script := `
			class Cake {
				hello() {
					print "hello";
				}
				getName() {
					return this.name;
				}
			}
			class ChocolateCake < Cake {
				getName() {
					return super.getName() + " au chocolat";
				}
			}`
		expect := []string{
			"(class Cake nil (fun hello (params) (print \"hello\")) " +
				"(fun getName (params) (return (get (this) name))))",
			"(class ChocolateCake Cake (fun getName (params) " +
				"(return (+ (call (super getName) (args)) \" au chocolat\"))))"}
		matchAST(t, expect, script)
	})
}

func TestCompilerErrors(t *testing.T) {

	t.Run("missing ;", func(t *testing.T) {
		script := `print i`
		errMsg := "[line 1] Error at end: Expect ';' after value.\n"
		expectError(t, errMsg, script)
	})

	t.Run("invalid assignment target", func(t *testing.T) {
		script := `"name" = "Bob";`
		errMsg := "[line 1] Error at '=': Invalid assignment target.\n"
		expectError(t, errMsg, script)
	})

	t.Run("expect expression (synch advance)", func(t *testing.T) {
		script := `
			var a;
			fun echo(n) { print n;}}
			a = 1;
			fun add(a, b) { return a + b;}`
		errMsg := "[line 3] Error at '}': Expect expression.\n"
		expectError(t, errMsg, script)
	})

	t.Run("expect expression (sync return)", func(t *testing.T) {
		script := `
			var a;
			fun echo(n) { print n;}}
			fun add(a, b) { return a + b;}`
		errMsg := "[line 3] Error at '}': Expect expression.\n"
		expectError(t, errMsg, script)
	})

}

func TestAstPrettyPrint(t *testing.T) {

	script := `
		fun isPositive(n) {
			var res;
			if (n > 0) {
				res = true;
			} else {
				res = false;
			}
			return res;
		}
	    fun countTo(n) {
	    	var i = 0;
	    	while (true) {
				print i;
				i = i + 1;
				if (i > n) return;
			}
		}
		class Boat {
			init(name) {
				this.name = name;
			}
			sailTo(port) {
				print this.name + " is sailing to " + port;
			}
		}
		class SpeedBoat < Boat {}
		print isPositive(20);
		var myBoat = Boat("Never Mind");
		myBoat.sailTo("Lisboa");`

	expect := "\n(block\n" +
		"  (fun isPositive (params n)\n" +
		"    (var res)\n" +
		"    (if (> (n) 0)\n" +
		"      (block\n" +
		"        (assign res true))\n" +
		"      (block\n" +
		"        (assign res false)))\n" +
		"    (return (res)))\n" +
		"  (fun countTo (params n)\n" +
		"    (var i 0)\n" +
		"    (while true\n" +
		"      (block\n" +
		"        (print (i))\n" +
		"        (assign i (+ (i) 1))\n" +
		"        (if (> (i) (n))\n" +
		"          (return)))))\n" +
		"  (class Boat nil\n" +
		"    (fun init (params name)\n" +
		"      (set (this) name (name)))\n" +
		"    (fun sailTo (params port)\n" +
		"      (print (+ (+ (get (this) name) \" is sailing to \") (port)))))\n" +
		"  (class SpeedBoat Boat)\n" +
		"  (print (call (isPositive) (args 20)))\n" +
		"  (var myBoat (call (Boat) (args \"Never Mind\")))\n" +
		"  (call (get (myBoat) sailTo) (args \"Lisboa\")))"

	scanner := &Scanner{}
	tokens := scanner.ScanTokens(script)
	parser := &Parser{}
	program := &BlockStmt{parser.Parse(tokens)}
	got := program.PrettyPrint("\n", "  ")
	if expect != got {
		t.Errorf("Expected '%s' but got '%s'", expect, got)
	}

}

// ------------------
// Helper functions
// ------------------

func matchAST(t *testing.T, expect []string, script string) {

	t.Helper()

	scanner := &Scanner{}
	tokens := scanner.ScanTokens(script)
	parser := &Parser{}
	got := parser.Parse(tokens)

	if scanner.HadError() {
		t.Fatal("Error encountered while scanning")
	}

	if parser.HadError() {
		t.Fatal("Error encountered while parsing")
	}

	length := len(expect)
	if len(got) > length {
		length = len(got)
	}

	for i := 0; i < length; i++ {

		if i >= len(got) {
			t.Errorf("Expected statement\n'%s'\nwas missing in %dth position",
				expect[i], i+1)
		} else if i >= len(expect) {
			t.Errorf("Unexpected statement\n'%s'\nin %dth position",
				got[i], i+1)
		} else if got[i].String() != expect[i] {
			t.Errorf("Expected statement\n'%s'\nbut got\n'%s'\nin %dth position",
				expect[i], got[i], i+1)
		}
	}
}

func expectError(t *testing.T, errMsg string, script string) {

	t.Helper()

	b := &strings.Builder{}
	scanner := &Scanner{}
	scanner.RedirectErrors(b)
	tokens := scanner.ScanTokens(script)
	parser := &Parser{}
	parser.RedirectErrors(b)
	parser.Parse(tokens)

	if !parser.HadError() {
		t.Errorf("Expected Error '%s' but got none", errMsg)
	}

	got := b.String()
	if got != errMsg {
		t.Errorf("Expected Error '%s' but got '%s'", errMsg, got)
	}
}
