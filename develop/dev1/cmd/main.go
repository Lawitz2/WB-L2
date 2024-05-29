package main

import (
	timePrec "WB-L2/gettime"
	"fmt"
	"log/slog"
	"os"
)

/*
Создать программу печатающую точное время с использованием NTP -библиотеки. Инициализировать как go module.
Использовать библиотеку github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Требования:
Программа должна быть оформлена как go module
Программа должна корректно обрабатывать ошибки библиотеки: выводить их в STDERR и возвращать ненулевой код выхода в OS

*/

func main() {
	t, err := timePrec.GetTimePrec()
	if err != nil {
		slog.Error("err getting time: ", err.Error())
		os.Exit(1)
	}
	fmt.Printf("time: %s", t.String())
}
