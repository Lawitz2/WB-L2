package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

/*
Реализовать утилиту аналог консольной команды cut (man cut).
Утилита должна принимать строки через STDIN, разбивать по разделителю (TAB) на колонки и выводить запрошенные.

Реализовать поддержку утилитой следующих ключей:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

*/

func main() {
	f := flag.String("f", "1", "fields")
	d := flag.String("d", "\t", "delimiter")
	s := flag.Bool("s", false, "only strings with delimiter")
	flag.Parse()

	var flagArray []int

	box := strings.Split(*f, ",")
	for _, i := range box {
		num, err := strconv.Atoi(i)
		if err == nil {
			flagArray = append(flagArray, num)
		} else {
			if strings.Contains(i, "-") {
				box2 := strings.Split(i, "-")
				lower, _ := strconv.Atoi(box2[0])
				upper, _ := strconv.Atoi(box2[1])
				for n := lower; n <= upper; n++ {
					flagArray = append(flagArray, n)
				}
			}
		}
	}

	file, err := os.Open("input2.txt")
	if err != nil {
		slog.Error(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var text []string
	for scanner.Scan() {
		text = strings.Split(scanner.Text(), *d)
		fmt.Printf("encoded text: %v, slice len: %d\n", text, len(text))

		if *s && len(text) < 2 {
			continue
		}

		for _, fl := range flagArray {
			if len(text) < fl {
				continue
			} else {
				fmt.Printf("%s%s", text[fl-1], *d)
			}
		}
		fmt.Println("")
	}

}
