package main

import (
	"fmt"
	"sync"
)

func fanOut(done <-chan interface{}, source <-chan int, count int) []chan int {
	var result []chan int
	for i := 0; i < count; i++ {
		result = append(result, make(chan int))
	}
	/*
	  for _, each := range result {
	     defer close(each)

	*/
	go func() {
		for value := range source {
			var wg sync.WaitGroup
			wg.Add(len(result))
			for _, out := range result {
				go func(out chan int, value int) {
					defer wg.Done()
					out <- value
				}(out, value)
			}
			wg.Wait()
		}
	}()

	return result
}

func main() {
	s := "gopher"
	fmt.Printf("Hello and welcome, %s!\n", s)

	source := make(chan int)
	done := make(chan interface{})
	channels := fanOut(done, source, 3)
	go printChannel(channels[0], "rød")
	go printChannel(channels[1], "grønn")
	go printChannel(channels[2], "blå")

	for i := 1; i <= 5; i++ {
		//TIP <p>To start your debugging session, right-click your code in the editor and select the Debug option.</p> <p>We have set one <icon src="AllIcons.Debugger.Db_set_breakpoint"/> breakpoint
		// for you, but you can always add more by pressing <shortcut actionId="ToggleLineBreakpoint"/>.</p>
		// fmt.Println("i =", 100/i)
		// time.Sleep(1000 * time.Millisecond)
		source <- i
	}

	close(source)
}

func printChannel(outChan chan int, prefix string) {
	for source := range outChan {
		fmt.Printf("%s: %d \n", prefix, source)
	}
}
