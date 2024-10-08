package main

import (
	"fmt"
)

// SenderA shows an example of how to make a "class" in Go.
// The struct is like an object or dict, but by adding
// methods (below) we can do things with the data in the struct.
//
// Here we create class called SenderA with a method
// send on it below... the send method is a "value receiver"
// because it receives a copy of the struct (look at the
// parenthesis just after "func" -- that's a copy).
type SenderA struct {
	FirstName    string
	MessageCount int
}

// Send is a "value receiver", like a "static" class method.
//
// READ ONLY --
// These methods DO NOT alter the original values, they'll
//
//	  use a copy of the fields in the struct.
//
//		value receiver (usually 1 letter, first letter of struct)
//		  |
//		  |
//		 \/
func (s SenderA) Send(message string) {

	// Common Bug
	//   Value receivers use a COPY of fields, not the original.
	//   s.MessageCount is initially 0
	//   Each call to Send()
	//   	s.MessageCount++ happens on copy of 0, makes it 1
	//   	s.MessageCount++ happens on copy of 0, makes it 1
	//   	s.MessageCount++ happens on copy of 0, makes it 1
	//    ...
	s.MessageCount++ // ❌ Common Bug

	fmt.Printf("Send from %s: %s\n", s.FirstName, message)

	// ❌ Common Bug
	s.FirstName = "what will happen????" // nothing, not saved
}

// SenderB shows a similar setup, but with a "pointer receiver".
type SenderB struct {
	FirstName    string
	MessageCount int
}

// Send this time is a "pointer receiver", like an "instance" class method.
// conceptually forming a class:
//
//		pointer receiver (usually 1 letter, first letter of struct)
//	      notice the "*", it lets us update to the original fields
//		   | |
//		   | |
//		  \/\/
func (s *SenderB) Send(message string) {
	// As a pointer receiver, we can update the internal fields
	//   on the struct with each call, like MessageCount.
	s.MessageCount++ // ✅ update is saved
	fmt.Printf("Send %d from %s: %s\n", s.MessageCount, s.FirstName, message)

	// no need to return anything, s was modified
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

// SenderInterface - wait, what are interfaces?
//
// Let's say we didn't care *how* a message got sent,
// only that we wanted to send a message.
//
// Go has a "selfish" (compositional) interface style,
// which annoyingly is the opposite style of most
// other languages like Java, C++, C.
//
// If you ask anyone what annoys them the most about Go,
//
//	they'll bring this up. There are good reasons why
//	Go authors built it this way, (compile safety, perf),
//	but it's still annoying.
//
// In Java, C++ C, you make a method signature, then
//
//	 have a few objects/classes meet that method signature
//	 by having a method on the object.
//	    Example: 3 Classes (Stream, Database, Game)
//
//	    Begin by declaring the method contract
//
//	    Interface Conn
//	       Connect() void
//	       Disconnect() void
//
//		   Each class includes those methods
//		   Compiler enforces Interface checks
//	       class Stream<implements Conn> { public Connect(){...}, public Disconnect(){...} }
//	       class Database<implements Conn> { public Connect(){...}, public Disconnect(){...} }
//	       class Game<implements Conn> { public Connect(){...}, public Disconnect(){...} }
//
//	    End with POLYMORPHISM for all
//	       connections = [new Stream(), new Database(), new Game()]
//	       connections.forEach(c => c.Connect)
//	       connections.forEach(c => c.Disconnect)
//
// In Go, it's the opposite. Using the same above (stream, database, game):
//
//			Begin with POLYMORPHIC interface for just `func connectAll` below
//
//	     Interface Conn
//	       Connect() void
//	       Disconnect() void
//
//			func connectAll(connStructs []Conn){
//		       for i := range connStructs {
//	           connStructs[i].Connect()
//	           connStructs[i].DisConnect()
//		       }
//			}
//
//	     Compiler does NOT enforce interface checks below
//
//	     struct Stream{}
//	     func (s *Stream) Connect{}
//	     func (s *Stream) Disconnect{}
//
//	     struct Database{}... struct Game....
//
//	     Compiler DOES enforce interface checks if you do this:
//
//	     struct Stream{}
//	     func NewStream() Stream{
//	       s := Stream{}
//	       var _ Conn = &s   // Interface checked here
//	     }
//	     func (s *Stream) Disconnect{}
//
// Here we will make a SenderInterface just for our
// method to use below it. This allows our method to
// say hey I don't care what you give me as along
// as I can call the method "Send".
type SenderInterface interface {
	Send(message string)
}

// SendEmail allows us to "overload" it with any
// sender implementation we want. This is
// "polymorphism" in Go.
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
