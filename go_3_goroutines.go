package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	// Below is a goroutine, any time you use the word "go".
	// The function call after go is run in a concurrent goroutine,
	//   aka lightweight thread.
	//
	// This is not a good goroutine though, because our main thread
	// won't wait for it to finish. The program might exit before it
	// has a chance to run.
	go printFoo()

	// Below is another common way to see goroutines launched, with IIFE functions.
	// Also not a good goroutine, because our main thread won't wait for it.
	go func() {
		fmt.Println("I'm concurrent! but might not happen, because nothing waiting on me :(")
	}()

	// Below are channels, for use with goroutines. These help us send and receive
	// data to goroutines, and are how we wait on goroutines.
	//
	// Setup a channel of bools that can hold at most 1 value
	// until someone pulls the value out of the channel.
	// If the channel fills up (e.g. has 1 bool in it), anyone
	// who wants to put another bool in the channel will block
	// and wait until the channel is freed up.
	//
	// The channel starts off empty.
	boolChannel := make(chan bool, 1)
	// Kick off a goroutine and skip over it.
	go func() {
		// Running parallel to main thread now in here
		fmt.Println("I'm concurrent, and someone will wait on me :D")

		// Put a bool in the channel
		// blocks here until the channel has free space to put the message on
		boolChannel <- true
	}()
	// Read the value from the channel into result
	result := <-boolChannel // blocks here until something is in the channel it can take out
	if result {
		fmt.Println("Finished waiting for that goroutine")
	}

	// A more common way to get messages from a channel is use a for loop on it.
	// The for loop will block in its thread until receiving a message on the channel
	// or a channel close event.
	//
	// I typically see this when we only want a goroutine to run while consuming
	// a finite list of items, which involves heavier processing we don't want slowing down
	// the main thread.
	//
	// Performance: in a setup like below, finding a good channel length
	//   takes a bit playing around but is important. Why?
	//      We've got a goroutine pulling message out of the channel.
	//      We've got a main thread putting messages in the channel.
	//
	//      If the channel fills up, what happens? Main thread is blocked until
	//        the channel frees up -- this stops everything. In a normal application
	//        that means no other server requests can make it through, users will see
	//        it freeze, business will call you up at 3AM on full volume asking why the overnight
	//        dark pool trading job is hung and your firm is losing millions, asking
	//        you to get dressed and come in the office even though there's 2 feet
	//        of snow on the ground and you've got take your kid to school in the morning
	//        despite you not being the one that forgot to find a good buffer length on
	//        your channel, it was actually the front end javascript dev business pulled in to do a
	//        hack in Go because they wanted to ship a perf improvement before christmas so that the
	//        new just-hired bigshot vp that severely overpromised on our timeline could
	//        dodge owning a failure and stay on track for that 1% bonus to keep up payments
	//        on the ridiculous matte-black bmw m8 they keep bringing up at every status meeting
	//
	//      ... eh hem ...
	//
	//      Playing around involves guess-and-check with how fast we load data into the channel
	//      compared to how fast we can pull it out. Here's how to game it:
	//
	//        If the consumer always slower than producer:
	//          put both in goroutines, consumer will eventually block trying to load
	//            Still include buffer though! Unbuffered channels
	//              add waiting overhead on both sides of the channel (vs one side)
	//                e.g. very rarely see unbuffered channels e.g. make(chan string)
	//
	//          start with 10, see how fast it is
	//              e.g. make(chan string, 10)
	//            double the buffer, see how fast
	//            double again, see how fast, when it stops getting faster, there's your buffer
	//
	//        If the consumer is 20% slower sometimes, 20% faster sometimes:
	//          Buffer is great, buffer 2x or 40% of the variable amount.
	//
	//          If you don't know the amount (e.g. a stream of data), do above,
	//            start with 10, see how fast, double, check, double until no faster.
	//
	//        If the consumer is faster than the producer, still include a buffer
	//           but keep it small (just big enough to avoid overhead blocking on wait)
	//
	// Below a message channel is setup where the main thread sends messages
	//   and goroutine receives them. In this example, the main thread (producer)
	//   will be faster because array iteration is way faster than printing.
	//   We move the printing out to a goroutine, and add some buffer of 3
	messageChan := make(chan string, 3)
	go func() {
		for {
			select {
			case message, ok := <-messageChan:
				if !ok {
					log.Print(errors.New("messageChan unexpectedly closed"))
					return // stops the loop and goroutine
				}

				fmt.Println(message)
				// default: <-- NEVER USE FOR CHANNEL READS unless you really need to.
				//   Without default the thread will stop CPU usage until the channel has a message.
				//   But if you include default here, the CPU thread will be 100% busy
				//     running the outer loop while waiting for a message.
			}
		}
	}()
	//
	// Send some messages (goroutine will receive and print them)
	// This will block main if the channel fills up.
	for _, message := range []string{"foo", "bar", "bazz", "wham", "whack", "bang", "pop", "zow"} {
		select {
		case messageChan <- message: // tries to add the message to the channel if it is not full
			// message sent ok
		default:
			// drop the message, because the channel was full
			// ALWAYS INCLUDE default FOR CHANNEL WRITES
			fmt.Println("messageChan full, dropping message: " + message)
		}
	}
	//
	// This will notify the infinite loop in the goroutine above to stop,
	// and that goroutine will finish.
	//
	// Most production go apps do not close channels
	// DO NOT CLOSE CHANNELS UNLESS YOU REALLY (REALLY) KNOW WHAT YOU'RE DOING
	// close(messageChan) <-- will likely cause a panic

	// So how are channels typically used? Here's an example.
	// Let's spam stock market data, and convert it to emojis.
	//
	// We'll dedicate
	//   two goroutines as spammers
	//   four goroutines as emoji converters
	//   a stockTicker channel between spammers and converters
	//
	// Feel free to make more and max out your processor.
	// Prof accepts no liability for 🔥🔥🔥🔥, do at own risk.
	//
	// What's the max values in the channel? Hard to say, but if
	// you have more consumers than producers lower, otherwise
	// if more producers higher (more buffer)
	stockTickerChan := make(chan string, 100)
	//
	// The workers will send a message when they're done, one
	// per worker. If all four finish at the same time I need
	// to have enough channel buffer to support that (4 workers)
	workers := 4
	workerDoneChan := make(chan bool, workers)
	//
	// If 5 seconds goes by this helpful go utility
	// will send a time message after 5 seconds.
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop() // ALWAYS schedule it to stop later, otherwise mem leak
	tickerChan := ticker.C

	// Fire up 4 workers listening for data on the ticker channel.
	// When they get a symbol, they'll convert it.
	maxEmojis := 100000
	go stockEmojiWorker(stockTickerChan, workerDoneChan, "AAPL", "🍎", maxEmojis)
	go stockEmojiWorker(stockTickerChan, workerDoneChan, "GOOG", "🤓", maxEmojis)
	go stockEmojiWorker(stockTickerChan, workerDoneChan, "FB", "🤢", maxEmojis)
	go stockEmojiWorker(stockTickerChan, workerDoneChan, "AMZN", "📦", maxEmojis)
	//
	// Fire up 2 spammers. As soon as these start running, data will
	// flow through the channel to the workers. Each worker arbitrarily
	// grabs a value off the channel.
	go stockSymbolSpammer(stockTickerChan)
	go stockSymbolSpammer(stockTickerChan)
	//
	// I want to know if the program timed out or we hit the max emojis.
	// Defer a message, and update this variable later.
	exitMessage := ""
	defer fmt.Println(exitMessage)
	//
	// The main thread fired off the goroutines, set the defer, and skipped here.
	// We need to wait for the goroutines now otherwise the program will exit.
	//
	// But this time I'm waiting on one of two things to happen... either
	// the workers all report done, or the timeout happens.
	//
	// Store the amount of workers we've heard from.
	finishedWorkers := 0
	for {
		// Select Case -- NOT TO BE CONFUSED WITH SWITCH
		//    Select case was designed to wait on messages from multiple
		//    goroutines, handling them one at a time.
		//
		// These are almost always in an infinite loop, until some condition reached.
		// This select case waits on either the timeout or 4 worker done messages.
		select {
		case <-tickerChan:
			// If the timeout utility sends a time, we'll get in here and done.
			exitMessage = "\n\ngot a timeout message"

			return // exits main, program done.
		case <-workerDoneChan:
			finishedWorkers++
			if finishedWorkers == workers {
				exitMessage = "\n\nheard from all the workers"

				return // exits main, program done.
			}
		}
	}
}

func printFoo() {
	fmt.Println("foo")
}

// I spam whatever channel you give me.
func stockSymbolSpammer(stockChan chan string) {
	stockSymbols := []string{"AAPL", "GOOG", "FB", "AMZN"}

	for {
		// Randomly pick a symbol
		randomStock := stockSymbols[rand.Intn(len(stockSymbols))]

		// Put it in the channel
		stockChan <- randomStock
	}
}

// I convert whatever stocks you give me to emoji, up to max emojis.
func stockEmojiWorker(stockChan chan string, doneChan chan bool, ticker, icon string, maxEmojis int) {
	emojiCount := 0

	// Continuously read from the channel with range, so helpful!
	for stockSymbol := range stockChan {
		// If it matches the ticker, convert it to emoji.
		if stockSymbol == ticker {
			emojiCount++
			fmt.Print(icon)
		}

		if emojiCount > maxEmojis {
			break
		}
	}

	doneChan <- true
}
