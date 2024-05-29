package main

import (
	"fmt"
	"math"
)

/*
Посетитель позволяет расширить количество действий, которые объекты могут выполнять, без изменения кода самого объекта
*/

type shape interface {
	getType() string
	accept(visitor)
}

type square struct {
	side float64
}

func (s *square) accept(v visitor) {
	v.visitForSquare(*s)
}

func (s *square) getType() string {
	return "i am a square"
}

type circle struct {
	radius float64
}

func (c *circle) accept(v visitor) {
	v.visitForCircle(*c)
}

func (c *circle) getType() string {
	return "i am a circle"
}

type perfecttriangle struct {
	side float64
}

func (t *perfecttriangle) accept(v visitor) {
	v.visitForTriangle(*t)
}

func (t *perfecttriangle) getType() string {
	return "i am a perfecttriangle"
}

// Добавление функционала без изменения оригинальной структуры
type visitor interface {
	visitForSquare(square)
	visitForCircle(circle)
	visitForTriangle(perfecttriangle)
}

type areaCalc struct {
	area float64
}

func (a *areaCalc) visitForSquare(s square) {
	a.area = s.side * s.side
	fmt.Printf("square area: %f\n", a.area)
}

func (a *areaCalc) visitForCircle(c circle) {
	a.area = c.radius * c.radius * math.Pi
	fmt.Printf("circle area: %f\n", a.area)
}

func (a *areaCalc) visitForTriangle(t perfecttriangle) {
	a.area = t.side * t.side * math.Sqrt(3) / 4
	fmt.Printf("triangle area: %f\n", a.area)
}

func main() {
	sq := square{side: 2}
	tr := perfecttriangle{side: 3.5}
	cir := circle{radius: 2.5}

	area := areaCalc{}
	area.visitForSquare(sq)
	area.visitForTriangle(tr)
	area.visitForCircle(cir)
}
