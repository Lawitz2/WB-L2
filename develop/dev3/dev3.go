package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

/*
Отсортировать строки в файле по аналогии с консольной утилитой sort (man sort — смотрим описание и основные параметры):
на входе подается файл из несортированными строками, на выходе — файл с отсортированными.

Реализовать поддержку утилитой следующих ключей:

-k — указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел)
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки
*/
func output(s [][]string) {
	out, err := os.Create("output.txt")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	writer := bufio.NewWriter(out)

	for ind, i := range s {
		for ind2, d := range i {
			writer.WriteString(d)
			if ind2 < len(i) {
				writer.WriteString(" ")
			}
		}
		if ind < len(s)-1 {
			writer.WriteString("\n")
		}
	}
	writer.Flush()
}

func outputnum(s [][]int) {
	out, err := os.Create("output.txt")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	writer := bufio.NewWriter(out)

	for ind, i := range s {
		for ind2, d := range i {
			box := strconv.Itoa(d)
			writer.WriteString(box)
			if ind2 < len(i) {
				writer.WriteString(" ")
			}
		}
		if ind < len(s)-1 {
			writer.WriteString("\n")
		}
	}
	writer.Flush()
}

func main() {
	var text []string
	r := flag.Bool("r", false, "reverse sort")
	u := flag.Bool("u", false, "do not print same strings")
	n := flag.Bool("n", false, "sort by numeric value")
	k := flag.Int("k", 0, "sort by K-th column, starting 0")

	flag.Parse()

	input, err := os.Open("input.txt")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	sameStringsMap := make(map[string]struct{})

	scanner := bufio.NewScanner(input)

	if !*u {
		for scanner.Scan() {
			text = append(text, scanner.Text())
		}
	} else {
		for scanner.Scan() {
			sameStringsMap[scanner.Text()] = struct{}{}
		}
	}

	if *u {
		for key := range sameStringsMap {
			text = append(text, key)
		}
	}

	var text2 [][]string

	for _, s := range text {
		text2 = append(text2, strings.Split(s, " "))
	}

	if *n {
		text2numbers := make([][]int, len(text2))
		for ind, i := range text2 {
			for _, d := range i {
				box, _ := strconv.Atoi(d)
				text2numbers[ind] = append(text2numbers[ind], box)
			}
		}

		sort.Slice(text2numbers, func(i, j int) bool {
			if len(text2numbers[i]) > *k && len(text2numbers[j]) > *k {
				return text2numbers[i][*k] < text2numbers[j][*k]
			} else {
				return len(text2numbers[i]) > *k
			}
		})

		if *r {
			slices.Reverse(text2numbers)
		}

		outputnum(text2numbers)
		return
	}

	sort.Slice(text2, func(i, j int) bool {
		if len(text2[i]) > *k && len(text2[j]) > *k {
			return text2[i][*k] < text2[j][*k]
		} else {
			return len(text2[i]) > *k
		}
	})

	if *r {
		slices.Reverse(text2)
	}

	output(text2)
}
