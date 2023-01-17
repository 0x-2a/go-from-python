// Prof Hibschman's quick ramp-up for Python/JS devs
// to pick up writing Golang. Part 1, more to follow.

// We always need a package name, like python package naming.
// For 435 will usually be main.
package main

// Imports pull in Go from files other directories, like python.
// GoLand will auto-import these for you.
//
// These imports are all from Golang (the standard library),
// no need to install additional libraries.
import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Like Java, Python... main kicks things off.
func main() {
	// ******************************************************************************************************
	// ******************************************************************************************************
	// Printing out
	// ******************************************************************************************************
	// ******************************************************************************************************

	// Print to console
	fmt.Println("Hello World!")                              // Like python print or console.log
	fmt.Printf("Hello %s\n", "from Printf")                  // Interpolated printing
	fmt.Printf("Print anything with %v \n", []string{"foo"}) // Interpolated printing

	// For pro projects, typically log is used instead of fmt.Print because log has log levels.
	// In pro runtimes it is common to shut off debug and trace log levels.
	log.Println("A message with return")
	log.Print("a message with no return")
	log.Printf("a message with no return")
	// log.Fatal kills the program

	// ******************************************************************************************************
	// ******************************************************************************************************
	// Variables
	// ******************************************************************************************************
	// ******************************************************************************************************
	str := "a string" // declare a variable + set the value, "walrus operator"
	str = "updated!"  // update a variable

	// The var syntax is used mostly with maps and arrays, waits until later to allocate.
	// Unlike JS or python, reading this will give a default value for primitives
	//
	// string     ""
	// int/float  0
	// struct     gives a struct instance with default values
	//
	// array map function interface will be "nil" though, meaning there is no pointer yet to an array,map, etc
	var myStr string // declares a variable (without a value, like None in Python, undefined in JS)

	num := 42        // declare int (auto sized to 32 or 64 bit)
	numFloat := 42.4 // declare float64

	i, j, k := 1, 2, 3 // declare multiple

	// The most common equivalent to python/js arrays in Go are called "slices"
	// A slice is an array with variable length.
	wordsSlice := []string{"foo", "bar", "bazz"} // a slice of strings
	fmt.Println(wordsSlice[0])

	numsSlice := []int{0, 1, 3} // all items have to be the same type
	if len(numsSlice) > 0 {     // how to check length
		fmt.Println(wordsSlice[0])
	}

	// When Go developers say "array" they mean "fixed-sized" array.
	// These are rarely used.
	myFixedArr := [3]int{} // a fixed array with 3 integers

	// Maps are like Python Dict {} or JS plain object {}
	keyValMap := map[string]string{
		"foo": "bar",
		"bim": "bazz",
	}

	// Structs are like classes.
	// They start as just plain storage of variables.
	// In other languages, this is like a model or data transfer object (DTO)
	type User struct {
		Name     string
		Password string
	}
	aliceUser := User{Name: "Alice", Password: "Gopher123"}
	bobUser := User{Name: "Bob", Password: "Gopher456"}

	// You can use structs in slices / arrays
	usersSlice := []User{
		aliceUser,
		bobUser,
		{Name: "Cindy", Password: "Gopher789"}, // declare inline without the type
	}

	// You can use structs on BOTH sides of maps too!
	userBuddyMap := map[User]User{
		aliceUser: bobUser,
		bobUser:   aliceUser,
	}

	// Read from the map
	alicesBuddy, keyExists := userBuddyMap[aliceUser]
	if keyExists {
		fmt.Println(alicesBuddy)
	}

	// Empty Variables
	var emptyInt int       //  declares an empty variable, sets default value of the type, 0
	var emptyString string //  declares an empty variable, sets default value of the type, ""

	// Declares an empty "slice", dynamic array -- like array in Python, JS
	var slice []string      // very common to do this when adding items to a temporary slice
	var goArr [4]string     // declares an empty array, fixed size, rarely used in my experience
	var aMap map[string]int // declares a nil map, rarely do it this way, usually initialize the map (see maps later)

	_ = str // lets you compile with an unused variable, dumping it to _
	_ = myStr
	_ = myFixedArr
	_ = aliceUser
	_ = usersSlice

	// ******************************************************************************************************
	// ******************************************************************************************************
	// If/for/switch/Comparison
	// ******************************************************************************************************
	// ******************************************************************************************************

	// If else (parenthesis discouraged unless it clarifies a boolean)
	if emptyString == "" {
		fmt.Println("true")
	} else if emptyString == "foo" {
		fmt.Println("else if true")
	} else {
		fmt.Println("else here")
	}

	// For range loops, like python range over list
	for i, value := range slice {
		// For the C++ devs, the value is a copy
		// You can speed it up by omitting value above
		fmt.Printf("%d %s", i, value)
	}
	// Ignore i
	for _, value := range slice {
		fmt.Printf("just wanted the value %s", value)
	}
	// Ignore value
	for i := range slice {
		// Fastest version, does not copy the array value each time.
		fmt.Printf("just wanted the i %d", i)
	}

	// Compare bool, string, nums all with ==
	aString := "foo"
	bString := "foo"
	fmt.Println(aString == bString) // true

	aNum := 1
	bNum := 1
	fmt.Println(aNum == bNum) // true

	// Range over map (like python dict range)
	// The order is random.
	for key, value := range keyValMap {
		fmt.Printf("%s %s", key, value)
	}

	// For loops
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}

	// There is no "while" in golang
	// Because you do it like this.
	for num < 10 {
		fmt.Print("how to make a while loop")
		break // how to break out
	}

	// Infinite loop is commonly used in a goroutine (separate thread)
	// to constantly pull items off a list where another thread is adding them.
	// See channels later.
	for {
		fmt.Print("how to make an infinite loop")
		break
	}

	// Switch
	// - The cases do not fall through
	// - Don't confuse with select case, which is used for goroutines
	switch str {
	case "foo":
		fmt.Println("got foo")
	case "bar":
		fmt.Println("got bar, did not fall through")
	default:
		fmt.Println("did not match above, guess i'll go then")
	}

	// Select - a switch for channels
	// The closest we get to "async" in python3.
	//
	// This example below has one thread writing messages
	//   and one thread reading them.
	//
	// Make a "channel" to pass messages, aka a pipe, aka shared queue
	//   chan string - we're passing strings through it from one thread to another
	//   1024 - Set how many messages will fit in the pipe until the sender is blocked
	messageChannel := make(chan string, 1024)

	// Fire up a thread that pulls messages out of the channel and prints.
	go func() {
		for {
			select { // ALWAYS USE SELECT WITH CHANNELS

			// Wait/Read from the channel
			// - Blocks this thread until there is a message
			case message, ok := <-messageChannel:
				if !ok {
					fmt.Println("channel closed")
				} else {
					fmt.Println(message)
				}

				//default:
				// default means "when the cases are all blocked do work here"
				//
				// Typically when reading from channels you won't want a default case.
			}
		}
	}()

	// Fire up another thread that puts messages onto the channel.
	go func() {
		for {
			select { // ALWAYS USE SELECT WITH CHANNELS

			// Constantly add strings containing "foo" to the channel.
			case messageChannel <- "foo":
				fmt.Println("sent a message")
			default:
				// Typically include when sending, we don't want to slow down or block a sender
				fmt.Println("channel is full!")
			}
		}
	}()

	// ******************************************************************************************************
	// ******************************************************************************************************
	// String Helpers
	// ******************************************************************************************************
	// ******************************************************************************************************

	// Contains includes
	hasWord := strings.Contains("some words", "word") // true

	// Split
	someString := "one,two,three,four "
	words := strings.Split(someString, ",") // []string{"one", "two", ...}

	// Join
	backTogether := strings.Join(words, " ") // "one two three four "

	// Get each letter
	letters := "abcd"
	firstLetter := string(letters[0]) // "a"
	// without casting to string, it's a "rune", not a char
	firstRune := letters[0] // golang uses "runes", not chars, which are like integer versions of the symbol

	// Iterate over letters
	for _, eachRune := range letters {
		letter := string(eachRune)
		fmt.Println(letter)
	}

	// Interpolation To String
	easySentence := "Add " + " words " + " together " + " without " + " perf " + " hits "
	easySentenceNum := "Digits " + strconv.Itoa(42) + " yay "
	sentence := fmt.Sprintf("A word here: %s, an int here: %d, a float here: %.2f", "hello", 42, 42.42)

	// Conversion From String
	idInt, _ := strconv.Atoi("234")               // String to int
	idInt64, _ := strconv.ParseInt("234", 10, 64) // String to int64
	boolStr := strconv.FormatBool(true)

	// ******************************************************************************************************
	// ******************************************************************************************************
	// Arrays, Slices, Lists
	// ******************************************************************************************************
	// ******************************************************************************************************

	// Create some slices (like python lists), far more common than go arrays
	emptySlice := []int{}  // declares an empty slice, this style typically avoided for slices
	var myEmptySlice []int // delcare an empty "nil" slice, can append to it right away

	strSlice := []string{"a", "b", "c"} // Variable size array (like python list)
	moreLetters := []string{"e", "f", "g"}
	numbersSlice := []int{2, 3, 5, 7, 11, 13} // Variable size array (like python list)

	// Create fixed size Array
	numFixedSizeArr := [6]int{2, 3, 5, 7, 11, 13} // Fixed size array

	// Check array/slice length
	if len(strSlice) > 1 {
		fmt.Println(strSlice[1])
	}

	// Read/Write to array/slice
	numbersSlice[0] = 1              // update array
	numFixedSizeArr[0] = 1           // update array
	numFromArr := numFixedSizeArr[4] // read array value (does a copy)
	numFromSlice := numbersSlice[4]  // same for slices

	// Get part of array/slice
	partOfArr := numFixedSizeArr[1:4] // 0-based, inclusive, exclusive
	partOfSlice := numbersSlice[1:4]  // 0-based, inclusive, exclusive
	everyThingBefore4 := numbersSlice[:4]
	everyThingStartingAt2 := numbersSlice[2:]

	// Add to the slice
	strSlice = append(strSlice, "d")            // Add one
	strSlice = append(strSlice, moreLetters...) // Add many

	// Sort the slice
	sort.Slice(moreLetters, func(i, j int) bool {
		return moreLetters[i] < moreLetters[j] // Ascending
	})

	// Slice of structs (like typed python dicts, or Typescript objects)
	// typically type declarations go at top of file, not inside functions
	type student struct {
		year int
		name string
	}

	nameYearSlice := []student{
		{2, "bob"},
		{3, "alice"},
		{5, "cindy"},
	}

	// ******************************************************************************************************
	// ******************************************************************************************************
	// Maps (like a python default dict)
	// ******************************************************************************************************
	// ******************************************************************************************************

	// Maps
	emptyMap := map[string]int{}      // creates a blank map ready for string keys that point to int values
	emptyMapB := make(map[string]int) // same as above, but a more formal style of writing it
	emptyMapC := map[User]User{}      // you can have non-primitive keys too! awesome and rare language feature
	nameToAge := map[string]int{      // create map with stuff in it
		"Bob":   42,
		"Alice": 33,
	}

	// Read from the map
	bobAge, keyExists := nameToAge["Bob"]
	if keyExists {
		fmt.Println(bobAge) // 42
	} else {
		fmt.Println(bobAge) // 0, the default value for int
	}

	// Write to the map, Update the map
	nameToAge["Bob"] = 34

	// Delete from the map
	delete(nameToAge, "Bob") // Remove key val, ignores if none there

	// You can use structs on BOTH sides of maps, awesome and rare language feature!
	dorisUser := User{Name: "Doris", Password: "Gopher123"}
	evanUser := User{Name: "Evan", Password: "Gopher456"}
	buddyMap := map[User]User{
		dorisUser: evanUser,
		evanUser:  dorisUser,
	}

	// Read from the map, checking if we have it
	dorisBuddyUser, keyExists := buddyMap[dorisUser]
	if keyExists {
		fmt.Println(dorisBuddyUser) // {Name: "Evan", Password: "Gopher456"}
	}

	// Iterate through the map
	for k, v := range nameToAge {
		delete(nameToAge, k) // this is safe!
		fmt.Printf("key[%s] value[%d]\n", k, v)
	}

	// Check how many keys the map has
	if len(nameToAge) == 0 {
		// empty!
	}

	// ******************************************************************************************************
	// ******************************************************************************************************
	// Time
	// ******************************************************************************************************
	// ******************************************************************************************************

	// Get the current time.
	now := time.Now()

	// Time in Zone (see https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)
	nyLocation, _ := time.LoadLocation("America/New_York")
	now.In(nyLocation)

	// Time from string (format, time string).
	tm, _ := time.Parse("2006-01-02 03:04:05", "2021-01-03 00:00:00")

	// Time to string
	timeStr := now.Format("Mon 2006-01-02 03:04:05 MST")

	// Time from timestamp
	time.Unix(1630357720, 0).In(nyLocation)
	time.Unix(1630357720, 0).In(time.UTC)

	// Time to unix timestamp
	time.Now().Unix()

	// Millis
	millis := time.Now().UnixNano() / int64(time.Millisecond)
	// Go 1.17 has
	// millisNew := time.Now().UnixMilli()

	_ = firstRune
	_ = emptyMap
	_ = emptyMapB
	_ = emptyMapC
	_, _, _ = tm, timeStr, millis
	_, _, _, _, _, _, _ = num, numFloat, emptyInt, emptyString, slice, goArr, wordsSlice
	_, _, _ = numsSlice, keyValMap, aMap
	_, _, _, _ = hasWord, backTogether, firstLetter, easySentence
	_, _, _, _, _ = easySentenceNum, sentence, idInt, idInt64, boolStr
	_, _, _ = i, j, k
	_, _, _, _, _, _, _, _, _ = emptySlice, myEmptySlice, numFromArr, numFromSlice, partOfArr, partOfSlice, everyThingBefore4, everyThingStartingAt2, nameYearSlice
}
