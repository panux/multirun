//Package multirun is a quick and easy package to run a batch task on multiple goroutines
package multirun

import (
	"reflect"
	"runtime"
	"sync"
)

//Runnable is an interface for any runnable task type
type Runnable interface {
	Run(int)
}

type simpleRunnable func(int)

func (r simpleRunnable) Run(iter int) {
	r(iter)
}

//SimpleRunnable creates a simple runnable with a function that takes a sequential number
func SimpleRunnable(task func(int)) Runnable {
	return simpleRunnable(task)
}

func runLoop(work <-chan int, locker *sync.RWMutex, r Runnable) {
	locker.RLock()
	defer locker.RUnlock()
	for i := range work {
		r.Run(i)
	}
}

//Run runs a task in parallell
func Run(r Runnable, count int, goroutines int) {
	work := make(chan int)
	locker := new(sync.RWMutex)
	for ; goroutines > 0; goroutines-- { //Start runners
		go runLoop(work, locker, r)
	}
	for count--; count > -1; count-- { //Distribute work
		work <- count
	}
	close(work)
	locker.Lock() //Wait for everything to finish
}

//Array iterates over an array in parallell
func Array(arr interface{}, task func(int)) {
	Run(SimpleRunnable(task), reflect.ValueOf(arr).Len(), runtime.NumCPU())
}
