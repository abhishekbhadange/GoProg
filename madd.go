package main

import (
	"fmt"
	"math/rand"
	"time"
	"os"
	"strconv"
	"runtime"
	)
	
func createMatrix(s int, c chan [][]int) {
	m := make([][]int, s)
	for i := range m {
		m[i] = make([]int, s)
	}
	c <- m
}

func loadMatrix(mat [][]int, size int, flag int, ch chan int) {
	if flag==0 {
		for i:=0; i<size; i++ {
			for j:=0; j<size; j++ {
				mat[i][j] = -1
			}
		}
	} else {
		for i:=0; i<size; i++ {
			for j:=0; j<size; j++ {
				value := rand.Intn(500)
				mat[i][j] = value
			}
		}
	}
	ch <- 1
}

func printMatrix(x [][]int, msize int) {
	for i:=0; i<msize; i++ {
		fmt.Println()
		for j:=0; j<msize; j++ {
			fmt.Print("ary[", i, "][", j, "] = ", x[i][j], " ")
		}
	}	
	fmt.Println("\n")
}

func addMatrices(A [][]int, B [][]int, C [][]int, chunk int, size int, ctr int, c chan int) {
	for i:=ctr; i<chunk; i++ {
		ctr++
		for j:=0; j<size; j++ {
			C[i][j] = A[i][j] + B[i][j]
		}
	}
	c <- 1 
}

func trace() time.Time {
	return time.Now()
}

func un(startTime time.Time) {
	endTime := time.Now()
	fmt.Println("Time:", (endTime.Sub(startTime)).Seconds(), "s")
}

func main() {

	startTime := trace()
	
	split, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	runtime.GOMAXPROCS(split)
	
	rand.Seed(time.Now().UTC().UnixNano())
	
	if len(os.Args) != 3  {
	  fmt.Println("Usage: fib [<cilk options>] <n>\n")
	  os.Exit(2)
	}
	
	size, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	rand.Seed(time.Now().UnixNano())
	
	arr := make(chan [][]int)
	
	go createMatrix(size, arr)
	go createMatrix(size, arr)
	go createMatrix(size, arr)
	A, B, C := <-arr, <-arr, <-arr
	
	ch := make(chan int)
	
	go loadMatrix(A, size, 12, ch)
	go loadMatrix(B, size, 10, ch)
	go loadMatrix(C, size, 0, ch)
	
	for i:=0; i<3; i++ {
		<-ch
	}
 	
	fmt.Println("\nMatrix A loaded with random numbers:\n")
	printMatrix(A,size)
	fmt.Println("\nMatrix B loaded with random numbers:\n")
	printMatrix(B,size)
	fmt.Println("\nMatrix C loaded with -1's:\n")
	printMatrix(C,size)
	
	chunk := (int)(size/split)
	
	ctr := 0
	
	for i:=0; i<split; i++ {
		if ctr == 0 {
			go addMatrices(A, B, C, chunk, size, ctr, ch)
			ctr = ctr + chunk
		} else if i == split-1 {
			ctr = chunk * i
			go addMatrices(A, B, C, size, size, ctr, ch)
		} else {
			go addMatrices(A, B, C, chunk + ctr, size, ctr, ch)
			ctr =  ctr + chunk
		}
	}
	
	fmt.Println("\nMatrix C holds result after A[]+B[]:\n")
	printMatrix(C,size)
	
	fmt.Print("Go Example: Matrix addition using Static Work Division\n");
	fmt.Println("running on", runtime.NumCPU(), "processor(s)\n")
	un(startTime)
	fmt.Println()
}
	