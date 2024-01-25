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
