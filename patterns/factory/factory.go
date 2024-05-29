package main

import "fmt"

/*
Фабричный метод применяется, когда стракты для выполнения программы всегда нужны одни и те же. Не требуется знать
все параметры стракта, вызов фабрики всё сделает за вас.
*/

type burgerI interface {
	setName(string)
	setMeatType(string)
	setCheese(bool)
	getName() string
	getMeatType() string
	getCheese() bool
}

type burgerFactory interface {
	createHamburger() burgerI
	createCheeseburger() burgerI
}

type burgerfac struct{}

func (b *burgerfac) createHamburger() burgerI {
	return &burger{
		name:     "hamburger",
		meatType: "beef",
		cheese:   false,
	}
}

func (b *burgerfac) createCheeseburger() burgerI {
	return &burger{
		name:     "cheeseburger",
		meatType: "chicken",
		cheese:   true,
	}
}

type burger struct {
	name     string
	meatType string
	cheese   bool
}

func (b *burger) setName(s string) {
	b.name = s
}

func (b *burger) setMeatType(s string) {
	b.meatType = s
}

func (b *burger) setCheese(b2 bool) {
	b.cheese = true
}

func (b *burger) getName() string {
	return b.name
}

func (b *burger) getMeatType() string {
	return b.meatType
}

func (b *burger) getCheese() bool {
	return b.cheese
}

func main() {
	burgerFac := &burgerfac{}
	ham := burgerFac.createHamburger()
	cheese := burgerFac.createCheeseburger()
	fmt.Println(ham, cheese)
}
