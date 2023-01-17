package main

import (
	"fmt"
)

// This is the Go equivalent of a class (or object with methods)
//
// Here we create class called SenderA with a method
// send on it below... the send method is a "value receiver"
// because it receives a copy of the struct (look at the
// parenthesis just after "func" -- that's a copy).
type SenderA struct {
	FirstName    string
	MessageCount int
}

// This is what links the methods to the struct,
// conceptually forming a static class:
//   value receiver (usually 1 letter, first letter of struct)
//     |
//     |
//    \/
func (s SenderA) Send(message string) {
	s.MessageCount++
	fmt.Printf("Send %d from %s: %s\n", s.MessageCount, s.FirstName, message)
}

// Here we create a totally different class "struct"
// but this time we want to be able to update the instance.
// In Golang we do this with "pointer receivers".
// Instead of copying the struct we'll pass a pointer
// to it in its methods (see the * just after "func" again).
type SenderB struct {
	FirstName    string
	MessageCount int
}

// This is what links instance methods to the struct,
// conceptually forming a class:
//   pointer receiver (usually 1 letter, first letter of struct)
//     |
//     |
//    \/
func (s *SenderB) Send(message string) {
	s.MessageCount++
	fmt.Printf("Send %d from %s: %s\n", s.MessageCount, s.FirstName, message)

	// no need to return anything, s was modified OK
}

// Let's create some for an example.
var (
	senderA = SenderA{FirstName: "A"} // creates a new copy of SenderA struct
	senderB = SenderB{FirstName: "B"} // creates a new copy of SenderB struct
)

// Let's test them out.
//
// See the issue in senderA?
// Because it is like a static class, its instance
// values don't change.
//
// Can you think of a way to alter senderA's Send method
// and logic below to make it work?
func runSenders() {
	senderA.Send("message one")   // Message 1 from A: message one
	senderA.Send("message two")   // Message 1 from A: message one
	senderA.Send("message three") // Message 1 from A: message three

	senderB.Send("message one")   // Message 1 from B: message one
	senderB.Send("message two")   // Message 2 from B: message two
	senderB.Send("message three") // Message 3 from B: message three
}

// Lastly, what are interfaces?
//
// Let's say we didn't care *how* a message got sent,
// only that we wanted to send a message.
//
// Go has a "selfish" (compositional) interface style.
// Rather than stating what methods you can offer for
// others to use (e.g. Java Interface), you state
// only what you want to use when you need it.
//
// Here we will make a SenderInterface just for our
// method to use below it. This allows our method to
// say hey I don't care what you give me as along
// as I can call the method "Send".
type SenderInterface interface {
	Send(message string)
}

// Now we can "overload" this method with any
// sender implementation we want.
func SendEmail(sender SenderInterface, message string) {
	sender.Send(message)
}

func runSendersInterface() {
	SendEmail(senderA, "message four")

	// Recall that SenderB is a pointer receiver, so
	// we need to pass it be reference (it's address)
	// to avoid a copy.
	//
	// Unlike C++ Golang makes resolving the pointer easy.
	// I.e. with the interface above we don't care how
	// the underlying memory is implemented either.
	SendEmail(&senderB, "message four")
}

func main() {
	runSenders()

	runSendersInterface()
}
