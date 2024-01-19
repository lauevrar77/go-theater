# Go Theater
This is an implementation of an actor system in Golang

## Example usage
```go
package main 

import (
	"fmt"
	"time"
)

type TimeTicker struct {
	me      *ActorRef
	mailbox *Mailbox
	system  *ActorSystem
	target  *ActorRef
}

func NewTimtTicker(target *ActorRef) *TimeTicker {
	return &TimeTicker{
		target: target,
	}
}

func (tt *TimeTicker) Initialize(me *ActorRef, mailbox *Mailbox, system *ActorSystem) {
	tt.me = me
	tt.mailbox = mailbox
	tt.system = system
}

func (tt *TimeTicker) Run() {
	target, err := tt.system.ByRef(*tt.target)
	if err != nil {
		return
	}
	run := true
	for run {
		select {
		case <-*tt.mailbox:
			run = false
		default:
			time.Sleep(1 * time.Second)
			*target <- Message{}
		}
	}
	fmt.Println("ticker quit")
}

type TimePrinter struct {
	me      *ActorRef
	mailbox *Mailbox
	system  *ActorSystem
}

func (tp *TimePrinter) Initialize(me *ActorRef, mailbox *Mailbox, system *ActorSystem) {
	tp.me = me
	tp.mailbox = mailbox
	tp.system = system
}

func (tp *TimePrinter) Run() {
	ticker := NewTimtTicker(tp.me)
	tickerMailbox, err := tp.system.Spawn("time-ticker", ticker, 10)
	if err != nil {
		return
	}
	cpt := 0
	for {
		msg := <-*tp.mailbox
		fmt.Println(msg)
		cpt += 1
		if cpt == 10 {
			*tickerMailbox <- Message{}
			break
		}
	}
	fmt.Println("printer quit")
}

func main() {
	printer := TimePrinter{}
	system := NewActorSystem()
	_, err := system.Spawn("time-printer", &printer, 10)
	if err != nil {
		panic(err)
	}
	system.Run()
}
```
