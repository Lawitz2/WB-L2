Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

package main

import (
"fmt"
)


func test() (x int) {
defer func() {
x++
}()
x = 1
return
}


func anotherTest() int {
var x int
defer func() {
x++
}()
x = 1
return x
}


func main() {
fmt.Println(test())
fmt.Println(anotherTest())
}
Ответ:

2, 1
В первом случае дефер увеличит Х на 1 перед возвращением, во втором значение запоминается при вызове ретурн