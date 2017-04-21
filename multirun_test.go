package multirun

import (
	"math/rand"
	"testing"
)

//Test multithreaded adding
func TestAdd(t *testing.T) {
	ran := rand.New(rand.NewSource(10))
	dat := make([]int, 1000000)
	expectsum := 0
	for i := range dat { //Populate array and find expected sum
		dat[i] = ran.Intn(100)
		expectsum += dat[i]
	}
	out := make([]int, len(dat)/100)
	Run(SimpleRunnable(func(iter int) {
		slice := dat[iter*100 : iter*100+100]
		sum := 0
		for _, v := range slice {
			sum += v
		}
		out[iter] = sum
	}), len(out), 7)
	sum := 0
	for _, v := range out {
		sum += v
	}
	if expectsum != sum {
		t.Fatalf("Expected sum %d but got %d\n", expectsum, sum)
	}
}
