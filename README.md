# Go Theater
This is an implementation of an actor system in Golang

## Example usage
```go
package main

import (
	"fmt"
	"github.com/lauevrar77/go-theater"
	"time"
)

type TimeTicker struct {
	me         theater.ActorRef
	dispatcher theater.MessageDispatcher
	system     *theater.ActorSystem

	target theater.ActorRef
	run    bool
}

func NewTimtTicker(target theater.ActorRef) *TimeTicker {
	return &TimeTicker{
		target: target,
		run:    true,
	}
}

func (tt *TimeTicker) Initialize(me theater.ActorRef, dispatcher theater.MessageDispatcher, system *theater.ActorSystem) {
	dispatcher.RegisterDefaultHandler(tt.Quit)

	tt.me = me
	tt.dispatcher = dispatcher
	tt.system = system
}

func (tt *TimeTicker) Quit(msg theater.Message) {
	tt.run = false
}

func (tt *TimeTicker) Run() {
	for tt.run {
		tt.dispatcher.Send(tt.target, theater.Message{})
		time.Sleep(1 * time.Second)
		tt.dispatcher.TryReceive()
	}
	fmt.Println("ticker quit")
}

type TimePrinter struct {
	me         theater.ActorRef
	dispatcher theater.MessageDispatcher
	system     *theater.ActorSystem
}

func (tp *TimePrinter) Initialize(me theater.ActorRef, dispatcher theater.MessageDispatcher, system *theater.ActorSystem) {
	dispatcher.RegisterDefaultHandler(tp.OnMessage)

	tp.me = me
	tp.dispatcher = dispatcher
	tp.system = system
}

func (tp *TimePrinter) OnMessage(msg theater.Message) {
	fmt.Println(msg)
}

func (tp *TimePrinter) Run() {
	ticker := NewTimtTicker(tp.me)
	_, err := tp.system.Spawn("time-ticker", ticker, 10)
	if err != nil {
		return
	}

	cpt := 0
	for {
		tp.dispatcher.Receive()
		cpt += 1
		if cpt == 10 {
			tp.dispatcher.Send(theater.ActorRef("time-ticker"), theater.Message{})
			break
		}
	}
	fmt.Println("printer quit")
}

func main() {
	printer := TimePrinter{}
	system := theater.NewActorSystem()
	_, err := system.Spawn("time-printer", &printer, 10)
	if err != nil {
		panic(err)
	}
	system.Run()
}
```

# How to get a synchronous answer ?
```go
package main

import (
	"fmt"
	"github.com/lauevrar77/go-theater"
	"time"
)

type TimeTicker struct {
	me         theater.ActorRef
	dispatcher theater.MessageDispatcher
	system     *theater.ActorSystem

	target theater.ActorRef
	run    bool
}

func NewTimtTicker(target theater.ActorRef) *TimeTicker {
	return &TimeTicker{
		target: target,
		run:    true,
	}
}

func (tt *TimeTicker) Initialize(me theater.ActorRef, dispatcher theater.MessageDispatcher, system *theater.ActorSystem) {
	dispatcher.RegisterDefaultHandler(tt.Quit)

	tt.me = me
	tt.dispatcher = dispatcher
	tt.system = system
}

func (tt *TimeTicker) Quit(msg theater.Message) {
	tt.run = false
}

func (tt *TimeTicker) Run() {
	for tt.run {
		tt.dispatcher.Send(tt.target, theater.Message{})
		time.Sleep(1 * time.Second)
		tt.dispatcher.TryReceive()
	}
	fmt.Println("ticker quit")
}

type TimePrinter struct {
	me         theater.ActorRef
	dispatcher theater.MessageDispatcher
	system     *theater.ActorSystem
}

func (tp *TimePrinter) Initialize(me theater.ActorRef, dispatcher theater.MessageDispatcher, system *theater.ActorSystem) {
	dispatcher.RegisterDefaultHandler(tp.OnMessage)

	tp.me = me
	tp.dispatcher = dispatcher
	tp.system = system
}

func (tp *TimePrinter) OnMessage(msg theater.Message) {
	fmt.Println(msg)
}

func (tp *TimePrinter) Run() {
	ticker := NewTimtTicker(tp.me)
	_, err := tp.system.Spawn("time-ticker", ticker, 10)
	if err != nil {
		return
	}

	cpt := 0
	for {
		tp.dispatcher.Receive()
		cpt += 1
		if cpt == 10 {
			tp.dispatcher.Send(theater.ActorRef("time-ticker"), theater.Message{})
			break
		}
	}
	fmt.Println("printer quit")
}

func main() {
	printer := TimePrinter{}
	system := theater.NewActorSystem()
	_, err := system.Spawn("time-printer", &printer, 10)
	if err != nil {
		panic(err)
	}
	system.Run()
}
```

# How to call an actor in an HTTP Handler
If you need a response, just call the actor synchronously (see above) otherwise, get the mailbox from the `ActorSystem` and send a message
