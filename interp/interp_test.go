package interp

import (
	"fmt"
	"os"
)

func ExampleLit() {

	runScript(`
		print 123;
		print true;
		print false;
		print "hello world!";
		print nil;
		`)
	// Output:
	// 123
	// true
	// false
	// hello world!
	// nil
}

func ExampleUnaryExpr() {

	runScript(`
		print - 12;
		print -13;
		print ! true;
	`)
	// Output:
	// -12
	// -13
	// false
}

func ExampleBinaryExpr() {

	runScript(`
		print 1 + 2;
		print "a" + "b";
		print "a" + 1;
		print "a" + true;
		print "a" + nil;
		print 1 - 2;
		print 1 / 2;
		print 2 * 3;
		print 1 + 2 * 3; /// checking operator priorities
		print (1 + 2) * 3;
	`)
	// Output:
	// 3
	// ab
	// a1
	// atrue
	// anil
	// -1
	// 0.5
	// 6
	// 7
	// 9
}

func ExampleLogicalExpr() {

	runScript(`
		print true and false;
		print true or false;
		print 1 < 2;
		print 2 <= 2;
		print 2 <= 1;
		print 1 > 2;
		print 2 >= 2;
		print 1 >= 2;
		print 3 == 3;
		print 3 != 3;
	`)
	// Output:
	// false
	// true
	// true
	// true
	// false
	// false
	// true
	// false
	// true
	// false
}

func ExampleVarDeclStmt() {

	runScript(`
		var a = 1;
		var b;
		print a + 1;
		print b;
	`)
	// Output:
	// 2
	// nil
}

func ExampleAssignExpr() {

	runScript(`
		var a;
		a = "hello";
		print a;
	`)
	// Output:
	// hello
}

func ExampleFunDeclStmt() {

	runScript(`
		fun hello(name) {
			print "Hello, " + name + "!";
		}
		print hello;
		hello("Bob");
	`)
	// Output:
	// <fun hello>
	// Hello, Bob!
}

func ExampleReturnStmt() {

	runScript(`
		fun add(a, b) {
			return a + b;
		}
		print add(10, 3);
		fun doNothingWithNegative(n) {
			if (n < 0) {
				print "I am not doing anything with " + n;
				return;
			}
			print "I am working with " + n;
		}
		doNothingWithNegative(-10);
	`)
	// Output:
	// 13
	// I am not doing anything with -10
}

func ExampleClassDeclStmt() {

	runScript(`
		class Cake {
			bake() {
				print "baking the cake!";
			}
		}
		print Cake;
		var myCake = Cake();
		print myCake;
		myCake.bake();
	`)
	// Output:
	// <class Cake>
	// <instance Cake>
	// baking the cake!
}

func ExampleThisExpr() {

	runScript(`
	class Cake {
		init(name) {
			this.name = name;
		}
		bake() {
			print "baking the " + this.name + "!";
		}
	}
	Cake("pie").bake();
	`)
	// Output:
	// baking the pie!
}

func ExampleSuperExpr() {

	runScript(`
		class Cake {
			bake(time) {
				print "cook for " + time + " minutes.";
			}
		}
		class ChocolateCake < Cake {
			bake(time) {
				print "cover in chocolate.";
				super.bake(time);
			}
		}
		ChocolateCake().bake(30);
	`)
	// Output:
	// cover in chocolate.
	// cook for 30 minutes.
}

func ExampleIfStmt() {

	runScript(`
		fun isPositive(n) {
			if (n > 0) {
				print n + " is positive";
			} else if (n < 0) {
				print n + " is negative";
			} else {
				print n + " is null";
			}
		}
		isPositive(10);
		isPositive(-10);
		isPositive(0);
	`)
	// Output:
	// 10 is positive
	// -10 is negative
	// 0 is null
}

func ExampleIfStmt_noBlock() {

	runScript(`
		fun isPositive(n) {
			if (n > 0) print n + " is positive";
			else if (n < 0)	print n + " is negative";
			else print n + " is null";
		}
		isPositive(10);
		isPositive(-10);
		isPositive(0);
	`)
	// Output:
	// 10 is positive
	// -10 is negative
	// 0 is null
}

func ExampleIfStmt_noElseBranch() {

	runScript(`
		fun isPositive(n) {
			if (n > 0) {
				print n + " is positive";
			}
		}
		isPositive(10);
		isPositive(-10);
	`)
	// Output:
	// 10 is positive
}

func ExampleWhileStmt() {

	runScript(`
		var i = 0;
		while (i < 3) {
			print i;
			i = i + 1;
		}
	`)
	// Output:
	// 0
	// 1
	// 2
}

func ExampleWhileStmt_forLoop() {

	runScript(`
		for (var i = 0; i < 3; i = i + 1) print i;
	`)
	// Output:
	// 0
	// 1
	// 2
}

// ------------------
// Error Conditions
// ------------------

func Example_compileErrorMissingSemicolon() {

	i := runScript(`print a`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 1] Error at end: Expect ';' after value.
	// true
	// false
}

// ------------------
// Helper Functions
// ------------------

func runScript(script string) *Interp {

	// we redirect both regular and error output to stdout
	// so we can use the golang testable example pattern
	// to check script execution.
	interp := New(os.Stdout, os.Stdout)
	interp.Run(script)
	return interp
}
