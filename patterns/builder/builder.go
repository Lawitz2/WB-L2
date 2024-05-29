package main

import "fmt"

/*
Строитель применяется, если требуется создать стракт с множеством параметров, некоторые из которых могут быть опциональны
*/

type Burger struct { //product\complex object we need to make
	meatType string
	cheese   bool
}

type IburgerBuilder interface { //interface for constructing said product
	setMeatType(string) IburgerBuilder
	setCheese(bool) IburgerBuilder
	makeBurger() *Burger
}

type burgerBuilder struct { //implements IburgerBuilder interface
	burger *Burger
}

func (b *burgerBuilder) setMeatType(s string) IburgerBuilder {
	b.burger.meatType = s
	return b
}

func (b *burgerBuilder) setCheese(b2 bool) IburgerBuilder {
	b.burger.cheese = b2
	return b
}

func (b *burgerBuilder) makeBurger() *Burger {
	return b.burger
}

func newBurgerBuilder() IburgerBuilder { //returns new burgerBuilder
	return &burgerBuilder{
		burger: &Burger{},
	}
}

type Cook struct { //provides interface for creating products
	builder IburgerBuilder
}

func (c *Cook) makeBurger(meatType string, cheese bool) *Burger {
	c.builder.setMeatType(meatType)
	c.builder.setCheese(cheese)

	return c.builder.makeBurger()
}

func main() {
	stove := newBurgerBuilder()

	chef := &Cook{stove}

	myBurger := chef.makeBurger("beef", true)
	fmt.Println(myBurger.meatType, myBurger.cheese)
}
