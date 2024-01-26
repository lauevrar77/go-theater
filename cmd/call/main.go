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
	me            *theater.ActorRef
	mailbox       *theater.Mailbox
	system        *theater.ActorSystem
	dispatcher    theater.MessageDispatcher
	shoudContinue bool
}

func (tp *TimeGiver) Initialize(me *theater.ActorRef, mailbox *theater.Mailbox, system *theater.ActorSystem) {
	tp.me = me
	tp.mailbox = mailbox
	tp.system = system
	tp.dispatcher = theater.NewMessageDispatcher(mailbox, system)

	tp.dispatcher.RegisterMessageHandler("ComputeTime", tp.PrintTime)
	tp.dispatcher.RegisterRequestMessageHandler("ComputeTime", tp.ReturnTime)
	tp.dispatcher.RegisterDefaultHandler(tp.Quit)
}

func (tp *TimeGiver) PrintTime(payload interface{}) {
	msg := payload.(ComputeTime)
	fmt.Println(tp.computeTime(msg.dur))
}

func (tp *TimeGiver) ReturnTime(payload interface{}) theater.Message {
	msg := payload.(ComputeTime)
	return theater.Message{
		Type: "ComputedTime",
		Content: ComputedTime{
			time: tp.computeTime(msg.dur),
		},
	}
}

func (tp *TimeGiver) Quit(payload interface{}) {
	tp.shoudContinue = false
}

func (tp *TimeGiver) Run() {
	tp.shoudContinue = true
	for tp.shoudContinue {
		tp.dispatcher.Receive()
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

	*mailbox <- theater.Message{
		Type: "ComputeTime",
		Content: ComputeTime{
			dur: 1 * time.Second,
		},
	}

	responseMsg, err := system.Call(&giverRef, theater.Message{
		Type: "ComputeTime",
		Content: ComputeTime{
			dur: 1 * time.Second,
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(responseMsg.Content.(ComputedTime).time)

	*mailbox <- theater.Message{
		Type:    "Quit",
		Content: nil,
	}
	wg.Wait()

}
