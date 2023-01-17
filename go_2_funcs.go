package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func main() {
	// Read 1-6 below then come back to 7 🙃
	//
	// ******************************************************************************************************
	// ******************************************************************************************************
	// 7. Deferred functions (for cleanup)
	// ******************************************************************************************************
	// ******************************************************************************************************
	// Usually these are placed at the start of the method or
	// immediately after we start something that needs to stop later.
	// Examples include
	// - database connections
	// - timers
	// - web connections / json reads
	// - running threads

	// Defer schedules this function call to run immediately after its
	// parent method returns (before the caller's next line).
	defer fmt.Println("I print after main returns and before whatever called main continues.")

	// It's common to create a function and immediately schedule the call it with defer.
	defer func() {
		fmt.Println("I print after main returns and before whatever called main continues.")
	}() // <-- notice the immediate function call

	// Real Example: postgres db connection
	//
	db, err := sql.Open( // Doesn't actually open the db, just parses out the connection info below
		"postgres",
		"host=127.0.0.1 port=5432 user=root password=root dbname=users sslmode=disable connect_timeout=3",
	)
	if err != nil {
		log.Fatal("Killing program, check the host/port/user string for syntax errors.")
	}
	conn, err := db.Conn(context.Background())
	if err != nil {
		log.Fatal("Killing program, couldn't reach the db for a connection.")
	}
	defer conn.Close() // ALWAYS HAVE THIS WITH DB

	// Fun fact! If you'd like to get back at the devops team for that prank
	// they pulled, or consistently hogging all the bathroom stalls over lunch,
	// read on :D.
	//
	// Yours truly caused a production outage because I didn't defer a close.
	// DB instances on the server have pools of connections, each with
	// a timeout. If you don't close your connections, the db holds onto
	// them until the timeout. If you are mining data from prod and your
	// program infinitely restarts (because you're lazy and copy-pasted
	// a dockerfile containing on-fail-restart), prod db access will eventually
	// fill up and prod (i.e. your company's web production product) will also halt.
	//
	// The dev ops lead will have face glowing as red as his error logs.
	// Note: a good devops team will protect against lazy SWEs with
	//  reasonable connection caps per client.

	// ******************************************************************************************************
	// ******************************************************************************************************
	// 1. Functions
	// ******************************************************************************************************
	// ******************************************************************************************************
	// Here's a quick ramp up on how to make functions in Go.
	//
	// 1. func main in a file in package main (like above)
	//    func main is a special function that is run as the binary entry point
	//    like java's main, c++'s main, python's main
	//
	//    func main must be in package main to be an entrypoint, but it
	//    does not need to be in the same folder (e.g. the file can be
	//    in a sub dir
	//
	//      Go projects either look like:
	//        myProject/    <- single binary project
	//          main.go     <- has package main and func main
	//          foo/
	//            foo_service.go <- has package foo and func Bar, can be imported elsewhere
	//
	//        myProject/   <- multiple binary project
	//          cmd/
	//            bar/
	//              bar.go     <- has package main and func main
	//            bazz/
	//              bazz.go     <- has package main and func main
	//
	// 2. func Foo at the package level (capitalize first letter)
	//    this is a public method, can be imported elsewhere
	//
	//    Not in package main though, because no one imports from package main
	//      instead, func main will import methods, which import others.
	//
	// 3. foo := func(){}, as a value, like below (similar to python, js, lambdas)
	//   Usually this is used as a helper to wrap some variables in one scope,
	//   an being operated upon to another scope (polymorphism)
	whizzBang := "whizzbang!"
	foo := func(bar string) string {
		// Outer scope variables in here get closured in (try to avoid doing this).
		whizzBangFromOutside := whizzBang
		fmt.Println(whizzBangFromOutside)

		// COMMON MEMORY LEAK BELOW
		//
		// Memory leaks happen if either
		//   the function stays in memory waiting to be reused
		//   or the function returns the address of the closured variable,
		//   extending its life beyond the scope here.
		//     e.g. return &whizzBang
		return bar + "bazz" + whizzBang
	}
	foo("bang") // call it normally

	// What happens to the above if I do this now?
	whizzBang = whizzBang + " more stuff "

	// 4. as an IIFE (immediately invoked function expression)
	//    Create a func and run it right away, not reusable.
	//
	//    Commonly used for goroutines, see further down.
	//    I usually don't see return values on these, but if you
	//    wanted one, just throw a variable on the left, with :=
	func(message string) {
		fmt.Println(message)
	}(whizzBang) // pass variables as input params here or left blank ().

	// 5. A value receiver, see below.
	// Like a class instance method, but immutable
	firstUser := user{}
	firstUser = firstUser.updateNameAndCopy("alice") // copy happened
	fmt.Println(firstUser.firstName)                 // alice

	// 6. A pointer receiver, see below.
	// Like a class instance method
	firstPerson := person{}
	firstPerson.updateMyName("jordan") // no copy happened!
	fmt.Println(firstPerson.firstName) // jordan
}

// 5. Value receivers on a struct (think of like a class)
//    This is Go's equivalent of an immutable class, which includes:
//      a. define the struct and its values (like class instance vars)
//      b. methods called "value receivers" that use the struct's values, change them
//
//    This type of "receiver" is immutable though, the struct
//    is copied any time its members are changed.
//
//    Typically the receiver method starts with func (then the
//    first letter of the struct's name (u for user) then the
//    struct it operates on (user).
//
//    See use case above.
type user struct {
	firstName string
}

func (u user) updateNameAndCopy(newName string) user {
	u.firstName = newName
	return u // a copy of the updated user struct
}

// 6. Pointer receivers on a struct (think of like a typical class)
//    This is Go's equivalent of a class, which includes:
//     a. define the struct and its values (same as above)
//     b. methods called "pointer receivers" the modify structs values (without a copy)
//
//    This type of "receiver" is mutable, the struct
//    is not copied when its members are changed.
//
//    Notice it looks the same as above, but with a pointer * in the receiver.
type person struct {
	firstName string
}

func (p *person) updateMyName(newName string) {
	p.firstName = newName
}
