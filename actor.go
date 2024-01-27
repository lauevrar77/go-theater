package theater

import "sync"

type Actor struct {
	Me          ActorRef
	Behavior    ActorBehavior
	Dispatcher  MessageDispatcher
	DeadQueue   chan ActorRef
	ActorSystem *ActorSystem
}

func NewActor(
	ref ActorRef,
	mailbox *Mailbox,
	behavior ActorBehavior,
	deadQueue chan ActorRef,
	as *ActorSystem,
) Actor {
	dispatcher := NewMessageDispatcher(mailbox, as)
	behavior.Initialize(ref, dispatcher, as)
	return Actor{
		Me:          ref,
		Behavior:    behavior,
		Dispatcher:  dispatcher,
		DeadQueue:   deadQueue,
		ActorSystem: as,
	}
}

func (a *Actor) Run(wg *sync.WaitGroup) {
	a.Behavior.Run()
	a.DeadQueue <- a.Me
	wg.Done()
}

type ActorBehavior interface {
	Initialize(ActorRef, MessageDispatcher, *ActorSystem)
	Run()
}
