package main

import (
	"fmt"
	"time"
	"math/rand"
	"os"
	"runtime"
	"strconv"
)

func createMatrix(s int, c chan [][]int) {
	m := make([][]int, s)
	for i := range m {
		m[i] = make([]int, s)
	}
	c <- m
}

func loadMatrix(x [][]int, flag int, msize int, ch chan int) {
	if flag == 0 {
		for i:=0; i<msize; i++ {
			for j:=0; j<msize; j++ {
				x[i][j] = -1
			}
		}	
	} else {
		for i:=0; i<msize; i++ {
			for j:=0; j<msize; j++ {
				value := rand.Intn(500)
				x[i][j] = value
			}
		}
	}
	ch <- 0
}

func printMatrix ( x [][]int, msize int) {
	for i:=0; i<msize; i++ {
		fmt.Println()
		for j:=0; j<msize; j++ {
			fmt.Print("ary[", i, "][", j, "] = ", x[i][j], " ")
		}
	}	
	fmt.Println("\n")
}

func multMatrices(A [][]int, B [][]int, C [][]int, size int, length int, ctr int, ch chan int) {
	if length == 1 {
		for j:=0; j<size; j++ {
			k := ctr
			C[ctr][j] = 0	
			for h := 0; h < size; h++ {
				C[ctr][j] += A[k][h] * B[h][j]
			}		
		}
		ch <- 0
	} else{
		chans := make(chan int)
		go multMatrices(A, B, C, size, length/2, ctr, chans)
		go multMatrices(A, B, C, size, length/2, ctr + length/2, chans)
		if length % 2 != 0 {
			go multMatrices(A, B, C, size, 1, ctr + length - 1, chans)
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

	if len(os.Args) != 3  {
	  fmt.Println("Usage: fib [<cilk options>] <n>\n")
	  os.Exit(2)
	}
	s := os.Args[2]
	n, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	
	rand.Seed(time.Now().UnixNano())
	
	arr := make(chan [][]int)
	
	go createMatrix(n, arr)
	go createMatrix(n, arr)
	go createMatrix(n, arr)
	A, B, C := <-arr, <-arr, <-arr

	chload := make(chan int)
	
	go loadMatrix(A, 27, n, chload)
	go loadMatrix(B, 53, n, chload)
	go loadMatrix(C, 0, n, chload)

	for i := 0; i < 3; i++ {
   	<-chload
	}

	fmt.Print("\nMatrix A loaded with random numbers:\n")
	printMatrix(A,n)
	fmt.Print("Matrix B loaded with random numbers:\n")
	printMatrix(B,n)
	fmt.Print("Matrix C loaded with -1's:\n")
	printMatrix(C,n)

	chunk := (int)(n/split)
	chans := make(chan int)
	
	var ctr int = 0
		for i:=0; i<split; i++ {
			if ctr == 0	{
				go multMatrices(A, B, C, n, chunk, ctr, chans)
				ctr = ctr + chunk
			} else if i == split - 1 {
				nchunk := n - chunk * i
				ctr = chunk * i
				go multMatrices(A, B, C, n, nchunk, ctr, chans)	
			} else {
				go multMatrices(A, B, C, n, chunk, ctr, chans)	
				ctr =  ctr + chunk
			}
		}
	
	for i := 0; i < split; i++ {
   	<-chans 
	}
	
	fmt.Print("Matrix C holds result after A[]*B[]:\n")
	printMatrix(C,n)
	
	fmt.Print("Go Example: Matrix multiplication using recursive Divide and Conquer\n")
	fmt.Println("running on", runtime.NumCPU(), "processor(s)\n")
	un(startTime)
	fmt.Println()
}
