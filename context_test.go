package learn_go_context

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

// Context
// is a package that tell another function to do something, like cancellation, timeout, deadline, pass value
// and context is Immutable
// Parent Child (inheritance)
// if you add something new to the context, it will make a new child of that context
// https://pkg.go.dev/context
func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

// Context with Value
// context.WithValue(parent, key, value)
// If you're trying to get value from context, and the context doesnt have, it will check its parent value
// https://pkg.go.dev/context#WithValue
func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextC, "e", "E")

	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)

	fmt.Println(contextE.Value("e")) // get from it self
	fmt.Println(contextE.Value("c")) // get from parent
	fmt.Println(contextE.Value("b")) // nil, cause it self and parent doesnt have the value of "b"
	fmt.Println(contextD.Value("c")) // nil, cause it self and parent doesnt have the value of "c"
}

// Context with Cancel
// can be used to cancel another function(goroutine)
// context.WithCancel(parent)
// https://pkg.go.dev/context#WithCancel

// Go routine with leak
// func CreateCounter() chan int {
// 	destination := make(chan int)

// 	go func() {
// 		defer close(destination)
// 		counter := 1
// 		for {
// 			destination <- counter
// 			counter++
// 		}
// 	}()

// 	return destination
// }

// func TestContextWithCancel(t *testing.T) {
// 	fmt.Println("Total Goroutine :", runtime.NumGoroutine())

// 	destination := CreateCounter()
// 	for i := range destination {
// 		fmt.Println("Counter", i)
// 		if i == 10 {
// 			break
// 		}
// 	}

// 	fmt.Println("Total Goroutine :", runtime.NumGoroutine())
// }

// Handle leak with cancel
func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done(): // wait for signal
				return
			default: // else, keep doing the default
				destination <- counter
				counter++
			}
		}
	}()

	return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println("Total Goroutine :", runtime.NumGoroutine())
	parent := context.Background()            // create the parent context
	ctx, cancel := context.WithCancel(parent) // create context with cancel

	destination := CreateCounter(ctx)
	for i := range destination {
		fmt.Println("Counter", i)
		if i == 10 {
			break
		}
	}
	cancel() // pass cancel signal to context

	fmt.Println("Total Goroutine :", runtime.NumGoroutine())
}

// Context with timeout
// context.WithTimeout(parent, duration)
// https://pkg.go.dev/context#WithTimeout
func CreateCounterTimeout(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done(): // wait for signal
				return
			default: // else, keep doing the default
				destination <- counter
				counter++
				time.Sleep(1 * time.Second) // make it slow for timeout simulation
			}
		}
	}()

	return destination
}

func TestContextWithTimeout(t *testing.T) {
	fmt.Println("Total Goroutine :", runtime.NumGoroutine())
	parent := context.Background()                            // create the parent context
	ctx, cancel := context.WithTimeout(parent, 5*time.Second) // create context with timeout, after 5 sec the goroutine will cancel
	defer cancel()

	destination := CreateCounterTimeout(ctx)
	for i := range destination {
		fmt.Println("Counter", i)
	}

	fmt.Println("Total Goroutine :", runtime.NumGoroutine())
}

// Context with deadline
// context.WithDeadline(parent, time)
// https://pkg.go.dev/context#WithDeadline
func TestContextWithDeadline(t *testing.T) {
	fmt.Println("Total Goroutine :", runtime.NumGoroutine())
	parent := context.Background()                                              // create the parent context
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(10*time.Second)) // create context with deadline(real time)
	defer cancel()

	destination := CreateCounterTimeout(ctx)
	for i := range destination {
		fmt.Println("Counter", i)
	}

	fmt.Println("Total Goroutine :", runtime.NumGoroutine())
}
