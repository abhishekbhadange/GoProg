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

func addMatrices(A [][]int, B [][]int, C [][]int, length int, size int, ctr int, ch chan int) {
	if length == 1 {
		for j:=0; j<size; j++ {
			C[ctr][j] = A[ctr][j] + B[ctr][j]
		}
		ch <- 0
	} else {
		chans := make(chan int)
		go addMatrices(A, B, C, length/2, size, ctr, chans)
		go addMatrices(A, B, C, length/2, size, ctr + length/2, chans)
		if length % 2 != 0 {
			go addMatrices(A, B, C, 1, size, ctr + length - 1, chans)
		}
		if length % 2 != 0 {
    		<-chans				
			<-chans
			<-chans

		} else {
			<-chans
			<-chans
		}
	}
	ch <- 0
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
	var nchunk int
	
	chans := make(chan int)
	
	for i:=0; i<split; i++ {
		if ctr == 0 {
			go addMatrices(A, B, C, chunk, size, ctr, chans)
			ctr = ctr + chunk
		} else if i == split-1 {
			nchunk = size - chunk * i
			ctr = chunk * i
			go addMatrices(A, B, C, nchunk, size, ctr, chans)
		} else {
			go addMatrices(A, B, C, chunk, size, ctr, chans)
			ctr = ctr + chunk
		}
	}
	
	for i := 0; i < split; i++ {
		<-chans 
	}
	
	fmt.Println("\nMatrix C holds result after A[]+B[]:\n")
	printMatrix(C,size)
	
	fmt.Print("Go Example: Matrix addition using recursive Divide and Conquer\n");
	fmt.Println("running on", runtime.NumCPU(), "processor(s)\n")
	un(startTime)
	fmt.Println()
}
	