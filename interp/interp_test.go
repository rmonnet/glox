package interp

import (
	"fmt"
	"os"
)

// -------------
// Expressions
// -------------

func ExampleAssignExpr() {

	runScript(`
		var a;
		a = "hello";
		print a;
		print a = 2;
	`)
	// Output:
	// hello
	// 2
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
		fun echo(n) { print n; }
		class Bar {}
		print "" + echo;
		print "" + Bar;
		print "" + Bar();
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
	// <fun echo>
	// <class Bar>
	// <instance Bar>
}

func ExampleCallExpr() {

	// with recursion
	runScript(`
		fun count(n) {
			if (n > 1) count(n-1);
			print n;
		}
		count(3);
	`)
	// Output:
	// 1
	// 2
	// 3
}

func ExampleCallExpr_implicitReturnNil() {

	runScript(`
		fun proc() {
			print "don't return anything";
		}
		print proc();
	`)
	// Output:
	// don't return anything
	// nil
}

func ExampleCallExpr_closure() {

	runScript(`
		fun makeCounter() {
			var i = 0;
			fun count() {
				i = i + 1;
				return i;
			}
			return count;
		}
		var counter1 = makeCounter();
		var counter2 = makeCounter();
		print counter1();
		print counter1();
		print counter2();
	`)
	// Output:
	// 1
	// 2
	// 1
}

func ExampleCallExpr_closure2() {

	runScript(`
		var a = "global";
		{
			fun showA() {
				print a;
			}
			showA();
			var a = "block";
			showA();
		}
	`)
	// Output:
	// global
	// global
}

func ExampleCallExpr_firstOrderFun() {

	runScript(`
		fun printIt(n) { print n; }
		fun thrice(fn) {
			for (var i = 1; i <= 3; i = i + 1) {
				fn(i);
			}
		}
		thrice(printIt);
	`)
	// Output:
	// 1
	// 2
	// 3
}

func ExampleGetExpr() {

	runScript(`
		class Saxophone {
			play() {
				print "Careless Whisper";
			}
		}
		class golfClub {
			play() {
				print "Fore!";
			}
		}
		fun playIt(thing) {
			thing.play();
		}
		var sax = Saxophone();
		var cart = golfClub();
		sax.play();
		cart.play();
		playIt(sax);
		playIt(cart);
	`)
	// Output:
	// Careless Whisper
	// Fore!
	// Careless Whisper
	// Fore!
}

func ExampleGetExpr_boundMethodPathological() {

	// when you store a bound method in a variable,
	// it stays linked to the original instance.
	runScript(`
		class Person {
			sayName() {
				print this.name;
			}
		}
		var jane = Person();
		jane.name = "Jane";
		var bill = Person();
		bill.name = "Bill";
		bill.oldSayName = bill.sayName;
		bill.sayName = jane.sayName;
		bill.sayName();
		bill.oldSayName();
	`)
	// Output:
	// Jane
	// Bill
}

func ExampleGetExpr_methodInheritance() {

	runScript(`
		class Level1 {
			doIt() {
				print "do it from level 1";
			}
			fakeIt() {
				print "fake it from level 1";
			}
			dreamIt() {
				print "dream it from level 1";
			}
		}
		class Level2 < Level1 {
			fakeIt() {
				print "fake it from level 2";
			}
			dreamIt() {
				print "dream it from level 2";
			}
		}
		class Level3 < Level2{
			dreamIt() {
				print "dream it from level 3";
			}
		}
		var l = Level3();
		l.dreamIt();
		l.fakeIt();
		l.doIt();
	`)
	// Output:
	// dream it from level 3
	// fake it from level 2
	// do it from level 1
}
func ExampleGetExpr_invokeInitDirectly() {

	runScript(`
		class Foo {
			init() {
				print this;
			}
		}
		var foo = Foo();
		print foo.init();
	`)
	// Output:
	// <instance Foo>
	// <instance Foo>
	// <instance Foo>
}

func ExampleGetExpr_methodVsVariable() {

	// We can store a method in a variable (bound method)
	// we can also store a plain function in a field
	runScript(`
		class Box {
			store(thing) {
				print "stored " + thing + " in the box";
			}
		}
		fun notMethod(arg) {
			print "called function with " + arg;
		}
		var box = Box();
		var store = box.store;
		box.function = notMethod;
		store("cookies");
		box.function(111);
	`)
	// Output:
	// stored cookies in the box
	// called function with 111
}

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

func ExampleLogicalExpr() {

	runScript(`
		print true and false;
		print false and true;
		print true or false;
		print false or true;
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
	// false
	// true
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

func ExampleLogicalExpr_truthy() {

	runScript(`
		fun isTruthy(a) { 
			if (a) {
				return true;
			} else {
				return false;
			}
		}
		class Boat {}
		fun echo(s) { print s; }
		print isTruthy(nil);
		print isTruthy(1);
		print isTruthy("");
		print isTruthy("hello");
		print isTruthy(echo);
		print isTruthy(Boat);
		print isTruthy(Boat());
	`)
	// Output:
	// false
	// true
	// true
	// true
	// true
	// true
	// true
}

func ExampleSetExpr() {

	runScript(`
		class Person {
			getName() {
				return this.name;
			}
		}
		var p = Person();
		p.name = "Bob";
		print p.name;
		print p.getName();
	`)
	// Output:
	// Bob
	// Bob

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

func ExampleVarExpr() {

	// redefinition of a variable is allowed
	runScript(`
		var a = "before";
		print a;
		var a = "after";
		print a;
	`)
	// Output:
	// before
	// after
}

func ExampleVarExpr_shadowing() {

	runScript(`
		var volume = 11;
		volume = 0;
		{
			var volume = 3 * 4 * 5;
			print volume;
		}
		print volume;
	`)
	// Output:
	// 60
	// 0
}

func ExampleVarExpr_enclosingVars() {

	runScript(`
		var global = "outside";
		{
			var local = "inside";
			print global + "/" + local;
		}
	`)
	// Output:
	// outside/inside

}

func ExampleVarExpr_enclosingVars2() {

	runScript(`
		var a = "global a";
		var b = "global b";
		var c = "global c";
		{
			var a = "outer a";
			var b = "outer b";
			{
				var a = "inner a";
				print a;
				print b;
				print c;
			}
			print a;
			print b;
			print c;
		}
		print a;
		print b;
		print c;
	`)
	// Output:
	// inner a
	// outer b
	// global c
	// outer a
	// outer b
	// global c
	// global a
	// global b
	// global c
}

// ------------
// Statements
// ------------

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

func ExampleReturnStmt_initAlwaysReturnThis() {

	runScript(`
		class Boat {
			init() {
				return;
			}
		}
		print Boat();
	`)
	// Output:
	// <instance Boat>
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

func ExampleWhileStmt_infiniteForLoop() {

	// if we use a for loop with no "condition", it loops forever
	// since we don't have a "break" statement, testing within
	// a function to use "return" as "break".
	runScript(`
		fun printTo(n) {
			var i = 0;
			for (;;) {
				print i;
				i = i + 1;
				if (i >n) return;
			}
		}
		printTo(3);
	`)
	// Output:
	// 0
	// 1
	// 2
	// 3
}

// ------------------
// Standard Library
// ------------------

func Example_libClock() {

	runScript(`
		var then = clock();
		var now = clock();
		print clock;
		print (now - then) <= 1;
	`)
	// Output:
	// <native fun>
	// true
}

// -----------------
// Compiler Errors
// -----------------

func Example_compileErrorMissingSemicolon() {

	i := runScript(`print a`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 1] Error at end: Expect ';' after value.
	// true
	// false
}

func Example_compileErrorMultipleLocalVarDecl() {

	// lox doesn't allow the same local variable to be
	// defined in the same *local* scope (global redeclarations
	// are allowed).

	i := runScript(`
		fun bad() {
			var a = "first";
			var a = "second";
		}
	`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 4] Error at 'a': Variable already declared in this scope.
	// true
	// false
}

func Example_compileErrorReturnValueFromInit() {

	i := runScript(`
		class BadBoat {
			init() {
				return "boat";
			}
		}
	`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 4] Error at 'return': Can't return a value from an initializer.
	// true
	// false
}

func Example_compilerErrorSelfReferencingClass() {

	i := runScript(`class Bar < Bar {}`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 1] Error at 'Bar': A class can't inherit from itself.
	// true
	// false
}

func Example_compilerErrorSelfReferencingVar() {

	i := runScript(`
		var a = "outer";
		{
			var a = a;
		}
	`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 4] Error at 'a': Can't read local variable in its own initializer.
	// true
	// false
}

func Example_compileErrorSuperWithNoSuperClass() {

	i := runScript(`
		class Eclair {
			cook() {
				super.cook();
			}
		}
	`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 4] Error at 'super': Can't use 'super' in a class with no superclass.
	// true
	// false
}

func Example_compileErrorThisInMethod() {

	i := runScript(`
		fun notAMethod() {
			print this;
		}
	`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 3] Error at 'this': Can't use 'this' outside of a class.
	// true
	// false
}

func Example_compilerErrorTopLevelReturn() {

	i := runScript(`return "at top level";`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 1] Error at 'return': Can't return from top-level code.
	// true
	// false
}

func Example_compilerErrorTopLevelSuper() {

	i := runScript(`super.greet();`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 1] Error at 'super': Can't use 'super' outside a class.
	// true
	// false
}

func Example_compileErrorTopLevelThis() {

	i := runScript(`print this;`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 1] Error at 'this': Can't use 'this' outside of a class.
	// true
	// false
}

// ----------------
// Runtime Errors
// ----------------

func Example_runtimeErrorArityMismatch() {

	i := runScript(`
		fun add(a, b, c) {
			return a + b +c;
		}
		print add(1,2,3);
		print add(1,2);
	`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// 6
	// [line 6] Expected 3 arguments but got 2.
	// false
	// true
}

func Example_runtimeErrorBadCall() {

	i := runScript(`
		var m = "hello";
		print m(2);
	`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 3] Can only call functions and classes.
	// false
	// true
}

func Example_runtimeErrorBadFieldLookup() {

	i := runScript(`
		var m = "hello";
		print m.name;
	`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 3] Only class instances have fields.
	// false
	// true
}

func Example_runtimeErrorBadFieldSet() {

	i := runScript(`
		var m = "hello";
		m.name = "Bob";
	`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 3] Only class instances have fields.
	// false
	// true
}

func Example_runtimeErrorBadOperandNumber() {

	i := runScript(`print (1 < "a");`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 1] Operand must be a number.
	// false
	// true

}

func Example_runtimeError_BadPlusOperands() {

	i := runScript(`true + 1;`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 1] Operands must be two numbers or at least one string.
	// false
	// true
}
func Example_runtimeErrorSuperclassNotAClass() {

	i := runScript(`
	var Cake = "Not a Class";
	class Eclair < Cake {}
	`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 3] Superclass must be a class.
	// false
	// true

}

func Example_runtimeErrorUndefinedAssignment() {

	i := runScript(`a = 123;`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 1] Undefined variable 'a'.
	// false
	// true
}

func Example_runtimeErrorUndefinedField() {

	i := runScript(`
		class Bar {}
		print Bar().name;	
	`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 3] Undefined field or method 'name'.
	// false
	// true
}

func Example_runtimeErrorUndefinedMethod() {

	i := runScript(`
		class Bar {}
		print Bar().name();	
	`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 3] Undefined field or method 'name'.
	// false
	// true
}

func Example_runtimeErrorUndefinedVariable() {

	i := runScript(`print a;`)
	fmt.Println(i.HadCompileError())
	fmt.Println(i.HadRuntimeError())
	// Output:
	// [line 1] Undefined variable 'a'.
	// false
	// true
}

// ------------------
// Helper Functions
// ------------------

func runScript(script string) *Interp {

	// we redirect both regular and error output to stdout
	// so we can use the golang testable example pattern
	// to check script execution.
	interp := New(os.Stdout, os.Stdout)
	interp.Run(script, false)
	return interp
}
