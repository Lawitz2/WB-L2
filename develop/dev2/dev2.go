package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"unicode"
)

/*
Создать Go-функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы/руны, например:
"a4bc2d5e" => "aaaabccddddde"
"abcd" => "abcd"
"45" => "" (некорректная строка)
"" => ""

Дополнительно
Реализовать поддержку escape-последовательностей.
Например:
qwe\4\5 => qwe45 (*)
qwe\45 => qwe44444 (*)
qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка, функция должна возвращать ошибку. Написать unit-тесты.
*/
func unzip(s string) (string, error) {
	var isletter, esc bool
	var builder strings.Builder
	var err error
	var box rune

	if s == "" {
		return "", nil // пустая строка - корректный ввод
	}

	for _, i := range s {
		switch {
		case esc || unicode.IsLetter(i):
			isletter = true
			builder.Write([]byte(string(i)))
			box = i
			esc = false

		case i == '\\':
			esc = true
			continue

		case unicode.IsNumber(i):
			if !isletter {
				continue
			}

			count, err := strconv.Atoi(string(i))
			if err != nil {
				//slog.Error("atoi error: ", err.Error())
				os.Exit(1)
			}

			for d := 0; d < count-1; d++ {
				builder.Write([]byte(string(box)))
			}
			isletter = false
		}
	}
	if builder.String() == "" {
		err = errors.New("incorrect string input")
	}
	return builder.String(), err
}

func main() {
	test := "a45"
	fmt.Println("input string: ", test)

	s, err := unzip(test)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	fmt.Println("unpacked string: ", s)
}
