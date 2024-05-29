package main

import "fmt"

/*
Команда (command) используется когда нужно выполнить несколько операций последовательно, запланировать их выполнение или выполнить их удаленно.
*/

type Icommand interface {
	execute()
}

type button struct {
	command Icommand
}

func (but *button) press() {
	but.command.execute()
}

type onCommand struct {
	device Idevice
}

func (rec *onCommand) execute() {
	rec.device.on()
}

type offCommand struct {
	device Idevice
}

func (rec *offCommand) execute() {
	rec.device.off()
}

type Idevice interface {
	on()
	off()
}

type tv struct {
	isRunning bool
}

func (t *tv) on() {
	if t.isRunning {
		fmt.Println("tv is already running")
	} else {
		t.isRunning = true
		fmt.Println("tv is running")
	}
}

func (t *tv) off() {
	if !t.isRunning {
		fmt.Println("tv is already off")
	} else {
		t.isRunning = false
		fmt.Println("tv is off")
	}
}

func main() {
	samsung := &tv{}

	onCom := &onCommand{device: samsung}
	offCom := &offCommand{device: samsung}

	onBut := &button{onCom}
	offBut := &button{offCom}

	onBut.press()
	onBut.press()
	offBut.press()
}
