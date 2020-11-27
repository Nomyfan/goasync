# GoAsync

Task-based async calling like C#

## Example
```go
package main

import (
	"fmt"
	"goasync"
	"time"
)

func main() {
	task1 := goasync.StartNewResultTask(func() goasync.Any {
		<-time.After(5 * time.Second)
		return 10
	}).ContinueWithAnyThenAny(func(any goasync.Any) goasync.Any {
		return 10 * any.(int)
	}).ContinueWithAnyThenVoid(func(m10 goasync.Any) {
		fmt.Printf("The result is %v\n", m10)
	}).ContinueWithVoidThenVoid(func() {
		fmt.Println("done")
	})

	task2 := goasync.StartNewVoidTask(func() {
		<-time.After(3 * time.Second)
	}).ContinueWithVoidThenAny(func() goasync.Any {
		return 10
	}).ContinueWithAnyThenVoid(func(any goasync.Any) {
		fmt.Printf("After 3 seconds, got result %v\n", any)
	})

	task1.Await()
	task2.Await()
}
```
Output
```shell script
After 3 seconds, got result 10
The result is 100
done
```