package main

import "fmt"

/*
Состояние используется, когда объекту требуется выполнять разные действия, в зависимости от котекста ситуации
*/

type person struct {
	base  string
	state Istate
}

func (p *person) setState(st Istate) {
	p.state = st
}

type Istate interface {
	getAction()
}

type calmState struct{}

func (c *calmState) getAction() {
	fmt.Println("i am calm so i put it down")
}

type madState struct{}

func (m *madState) getAction() {
	fmt.Println("i am mad so i throw it at the wall")
}

func main() {
	me := &person{}
	calm := &calmState{}
	mad := &madState{}

	me.base = "i have a phone in my hands, "

	me.setState(calm)
	fmt.Printf("%s", me.base)
	me.state.getAction() // i am calm so i put it down

	me.setState(mad)
	fmt.Printf("%s", me.base)
	me.state.getAction() // i am mad so i throw it at the wall
}
