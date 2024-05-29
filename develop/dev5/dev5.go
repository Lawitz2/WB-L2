package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

/*
Реализовать утилиту фильтрации по аналогии с консольной утилитой (man grep — смотрим описание и основные параметры).

Реализовать поддержку утилитой следующих ключей:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", напечатать номер строки
*/
func main() {
	A := flag.Int("A", 0, "after found")
	B := flag.Int("B", 0, "before found")
	C := flag.Int("C", 0, "before + after (context)")
	c := flag.Bool("c", false, "amount of lines")
	i := flag.Bool("i", false, "case insensitivity")
	v := flag.Bool("v", false, "invert")
	F := flag.Bool("F", false, "full match")
	n := flag.Bool("n", false, "line num")

	flag.Parse()

	if *A < 0 || *B < 0 || *C < 0 {
		fmt.Println("incorrect flag input, can't put negative A/B/C values")
		os.Exit(1)
	}

	if (*A != 0 || *B != 0) && *C != 0 {
		fmt.Println("incorrect flag input, can't put A and/or B with C")
		os.Exit(1)
	}

	if *C > 0 {
		*A, *B = *C, *C
	}

	var searchString string
	var lineNum, linesFound, printCounter, lenCh int
	var found bool
	var s, text string

	searchString = flag.Args()[0]

	bufferForA := make(chan string, *B)
	defer close(bufferForA)

	if *i {
		s = strings.ToLower(searchString)
	} else {
		s = searchString
	}

	for _, f := range flag.Args()[1:] {

		file, err := os.Open(f)
		if err != nil {
			slog.Error(err.Error())
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineNum++

			if *i {
				text = strings.ToLower(scanner.Text())
			} else {
				text = scanner.Text()
			}

			if *F {
				if s == text {
					found = true
				}
			} else {
				if strings.Contains(text, s) {
					found = true
				}
			}

			if *v {
				found = !found
			}

			if !found && *B > 0 {
				if len(bufferForA) == *B {
					<-bufferForA
				}
				bufferForA <- scanner.Text()
			}

			if found {
				linesFound++
				if *B > 0 {
					lenCh = len(bufferForA)
					for k := 0; k < lenCh; k++ {
						if *n {
							fmt.Printf("%d: ", lineNum+k-lenCh)
						}
						fmt.Printf("%s\n", <-bufferForA)
					}
				}

				printCounter = *A + 1
			}

			if printCounter > 0 {
				if *n {
					fmt.Printf("%d: ", lineNum)
				}
				fmt.Printf("%s\n", scanner.Text())
				printCounter--
			}

			found = false
		}

		if *c {
			fmt.Printf("lines matching search parameters: %d", linesFound)
		}
	}
}
