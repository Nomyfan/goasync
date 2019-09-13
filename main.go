package main

import (
	"fmt"
	"goasync/async"
	"time"
)

func goAsync(done chan bool) {
	fmt.Println("enter goAsync")
	time.Sleep(3 * time.Second)
	fmt.Println("leave goAsync")
	done <- true
}

func goAsync1() {
	fmt.Println("enter goAsync")
	time.Sleep(3 * time.Second)
	fmt.Println("leave goAsync")
}

func goGetXY() interface{} {
	time.Sleep(3 * time.Second)
	return time.Now().Nanosecond() % 1000
}

func doSth() {
	fmt.Println("doing sth")
	time.Sleep(1 * time.Second)
	fmt.Println("sth done")
}

func main() {

	//fmt.Println("enter main")
	//done := make(chan bool)
	//go goAsync(done)
	//doSth()
	//<-done
	//fmt.Println("leave main")

	//fmt.Println("enter main")
	//task := async.NewTask(goAsync1)
	//task.InvokeAsync()
	//doSth()
	//task.Await()
	//fmt.Println("leave main")

	taskX := async.NewResultTask(goGetXY)
	taskY := async.NewResultTask(goGetXY)
	start := time.Now()
	taskX.InvokeAsync()
	taskY.InvokeAsync()
	taskX.Await()
	taskY.Await()

	x, _ := taskX.GetResult().(int)
	y, _ := taskY.GetResult().(int)
	elapsed := time.Now().Sub(start)
	sum := x + y
	fmt.Printf("sum of x and y is %d \n", sum)
	fmt.Printf("Elapsed: %d ms\n", elapsed.Milliseconds())
}
