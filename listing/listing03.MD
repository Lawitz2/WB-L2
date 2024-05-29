Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

package main

import (
"fmt"
"os"
)

func Foo() error {
var err *os.PathError = nil
return err
}

func main() {
err := Foo()
fmt.Println(err)
fmt.Println(err == nil)
}
Ответ:

err содержит в себе интерфейс с нил указателем, 
и во втором случае программа скажет "неверно" потому что интерфейс с нил указателем это не нил.
https://research.swtch.com/interfaces