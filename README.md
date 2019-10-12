# GoAsync

Task-based async calling like C#

## Example1
```go
package main

import (
	"fmt"
	"goasync"
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

	fmt.Println("enter main")
	done := make(chan bool)
	go goAsync(done)
	doSth()
	<-done
	fmt.Println("leave main")

	fmt.Println("enter main")
	task := goasync.NewTask(goAsync1)
	task.InvokeAsync()
	doSth()
	task.Await()
	fmt.Println("leave main")

	taskX := goasync.NewResultTask(goGetXY)
	taskY := goasync.NewResultTask(goGetXY)
	start := time.Now()
	taskX.InvokeAsync()
	taskY.InvokeAsync()
	//taskX.Await() // Useless call here
	//taskY.Await() // Useless call here

	x, _ := taskX.GetResult().(int)
	y, _ := taskY.GetResult().(int)
	elapsed := time.Now().Sub(start)
	sum := x + y
	fmt.Printf("sum of x and y is %d \n", sum)
	fmt.Printf("Elapsed: %d ms\n", elapsed.Milliseconds())
}
```

## Example2
```go
package main

import (
	"fmt"
	"goasync"
	"time"
)

func main() {
	task := goasync.StartNewResult(func() interface{} {
		time.Sleep(3 * time.Second)
		return 2
	}).ContinueWithResult(func(t *goasync.Task) interface{} {
		v, _ := t.GetResult().(int)
		return v * v
	}).ContinueWith(func(t *goasync.Task) {
		fmt.Println("result is",t.GetResult())
	})
	fmt.Println("run first...")
	task.Await()
}
```

Output
```
run first...
result is 4
```