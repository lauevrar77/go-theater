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
	me      *theater.ActorRef
	mailbox *theater.Mailbox
	system  *theater.ActorSystem
	target  *theater.ActorRef
}

func NewTimtTicker(target *theater.ActorRef) *TimeTicker {
	return &TimeTicker{
		target: target,
	}
}

func (tt *TimeTicker) Initialize(me *theater.ActorRef, mailbox *theater.Mailbox, system *theater.ActorSystem) {
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
			*target <- theater.Message{}
		}
	}
	fmt.Println("ticker quit")
}

type TimePrinter struct {
	me      *theater.ActorRef
	mailbox *theater.Mailbox
	system  *theater.ActorSystem
}

func (tp *TimePrinter) Initialize(me *theater.ActorRef, mailbox *theater.Mailbox, system *theater.ActorSystem) {
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
			*tickerMailbox <- theater.Message{}
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
	"sync"
	"time"
)

type ComputeTime struct {
	dur time.Duration
}

type ComputedTime struct {
	time time.Time
}

type TimeGiver struct {
	me      *theater.ActorRef
	mailbox *theater.Mailbox
	system  *theater.ActorSystem
}

func (tp *TimeGiver) Initialize(me *theater.ActorRef, mailbox *theater.Mailbox, system *theater.ActorSystem) {
	tp.me = me
	tp.mailbox = mailbox
	tp.system = system
}

func (tp *TimeGiver) Run() {
	shouldContinue := true
	for shouldContinue {
		msg := <-*tp.mailbox
		switch msg.Type {
		case "RequestMessage":
			reqContent := msg.Content.(theater.RequestMessage)
			switch reqContent.Type {
			case "ComputeTime":
				content := reqContent.Content.(ComputeTime)
				computed := tp.computeTime(content.dur)
				mailbox, err := tp.system.ByRef(reqContent.RespondTo)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				}
				*mailbox <- theater.Message{
					Type: "ComputedTime",
					Content: ComputedTime{
						time: computed,
					},
				}
				break
			default:
				fmt.Println("Unknown message type")
				break
			}
			break
		case "ComputeTime":
			content := msg.Content.(ComputeTime)
			fmt.Println(tp.computeTime(content.dur))
			break
		case "Quit":
			fmt.Println("Quitting")
			shouldContinue = false
		}
	}
}

func (tp *TimeGiver) computeTime(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}

func main() {
	giver := TimeGiver{}
	giverRef := theater.ActorRef("time-giver")
	system := theater.NewActorSystem()

	_, err := system.Spawn(giverRef, &giver, 10)
	if err != nil {
		panic(err)
	}

	mailbox, _ := system.ByRef(giverRef)

	wg := sync.WaitGroup{}
	go func() {
		system.Run()
		wg.Done()
	}()
	wg.Add(1)

	responseMsg, err := system.Call(&giverRef, theater.Message{
		Type: "ComputeTime",
		Content: ComputeTime{
			dur: 1 * time.Second,
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(responseMsg)

	*mailbox <- theater.Message{
		Type:    "Quit",
		Content: nil,
	}
	wg.Wait()

}
```

# How to call an actor in an HTTP Handler
If you need a response, just call the actor synchronously (see above) otherwise, get the mailbox from the `ActorSystem` and send a message
