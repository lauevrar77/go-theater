package theater

import (
	"fmt"
	"sync"
)

type ActorSystem struct {
	actors map[ActorRef]Mailbox
	wg     sync.WaitGroup
}

func NewActorSystem() ActorSystem {
	return ActorSystem{
		actors: make(map[ActorRef]Mailbox),
		wg:     sync.WaitGroup{},
	}
}

func (as *ActorSystem) Spawn(ref ActorRef, behavior ActorBehavior, mailboxSize int) (*Mailbox, error) {
	if _, ok := as.actors[ref]; ok {
		return nil, fmt.Errorf("Actor already exists")
	}
	mailbox := make(Mailbox, mailboxSize)
	as.actors[ref] = mailbox
	behavior.Initialize(&ref, &mailbox, as)
	as.wg.Add(1)
	actor := Actor{Mailbox: &mailbox, Behavior: behavior}
	go actor.Run(&as.wg)
	return &mailbox, nil
}

func (as *ActorSystem) ByRef(ref ActorRef) (*Mailbox, error) {
	mailbox, ok := as.actors[ref]
	if !ok {
		return nil, fmt.Errorf("Actor not found")
	}
	return &mailbox, nil
}

func (as *ActorSystem) Run() {
	as.wg.Wait()
}
