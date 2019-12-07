package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	var input string

	fmt.Scanln(&input)

	code := parse(input)

	in := make(chan int, 1)
	in <- 1
	close(in)

	out := make(chan int)

	sem := make(chan struct{})

	var outs []int

	go func() {
		for n := range out {
			outs = append(outs, n)
		}
		sem <- struct{}{}
	}()

	if err := run(code, in, out); err != nil {
		log.Fatal(err)
	}

	close(out)

	<-sem

	fmt.Printf("Part 1: %d\n", outs[len(outs)-1])
}

func parse(code string) []int {
	ss := strings.Split(code, ",")
	ii := make([]int, 0, len(ss))

	for n := range ss {
		if i, err := strconv.Atoi(ss[n]); err == nil {
			ii = append(ii, i)
		}
	}

	return ii
}

func opcode(n int) int {
	return n % 100
}

func paramMode(n, param int) int {
	n /= 100
	for i := 0; i < param; i++ {
		n /= 10
	}
	return n % 10
}

func readVal(code []int, idx, pMode int) int {
	switch pMode {
	case 0:
		return code[code[idx]]
	case 1:
		return code[idx]
	}
	return 0
}

func writeVal(code []int, idx, pMode, val int) {
	switch pMode {
	case 0:
		code[code[idx]] = val
	case 1:
		code[idx] = val
	}
}

func run(code []int, in <-chan int, out chan<- int) error {
	var cur int

	for {
		switch opcode(code[cur]) {
		case 1:
			p1 := readVal(code, cur+1, paramMode(code[cur], 0))
			p2 := readVal(code, cur+2, paramMode(code[cur], 1))
			writeVal(code, cur+3, paramMode(code[cur], 2), p1+p2)
			cur += 4
		case 2:
			p1 := readVal(code, cur+1, paramMode(code[cur], 0))
			p2 := readVal(code, cur+2, paramMode(code[cur], 1))
			writeVal(code, cur+3, paramMode(code[cur], 2), p1*p2)
			cur += 4
		case 3:
			writeVal(code, cur+1, paramMode(code[cur], 0), <-in)
			cur += 2
		case 4:
			out <- readVal(code, cur+1, paramMode(code[cur], 0))
			cur += 2
		case 99:
			return nil
		default:
			return fmt.Errorf("unknown instruction %d at position %d", code[cur], cur)
		}
		if cur >= len(code) {
			return errors.New("out of bounds")
		}
	}
}
