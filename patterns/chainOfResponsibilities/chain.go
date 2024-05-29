package main

import (
	"fmt"
)

/*
Цепочка обязанностей применяется при обработке запросов, когда либо может поступать нессколько видов запросов, либо конкретика запроса неизвестна.
*/

type request struct {
	id   int
	data string
}

type Ihandler interface {
	handle(request)
	setNext(Ihandler)
}

type handlerA struct {
	nextInChain Ihandler
}

func (h *handlerA) setNext(ihandler Ihandler) {
	h.nextInChain = ihandler
}

func (h *handlerA) handle(r request) {
	if len(r.data) < 10 {
		fmt.Println("handler A is handling the request...")
		return
	} else {
		h.nextInChain.handle(r)
	}
}

type handlerB struct {
	nextInChain Ihandler
}

func (h *handlerB) handle(r request) {
	if len(r.data) < 30 {
		fmt.Println("handler B is handling the request...")
		return
	} else {
		fmt.Println("we're unable to handle your request")
		return
	}
}

func (h *handlerB) setNext(ihandler Ihandler) {
	h.nextInChain = ihandler
}

func main() {
	handA := &handlerA{}
	handB := &handlerB{}
	handA.nextInChain = handB

	req1 := request{
		id:   1,
		data: "i am a pretty long request",
	}

	req2 := request{
		id:   1,
		data: "short",
	}

	req3 := request{
		id:   1,
		data: "i am a very very very very very very very very long request",
	}

	handA.handle(req1)
	handA.handle(req2)
	handA.handle(req3)
}
