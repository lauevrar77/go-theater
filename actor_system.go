package theater

import (
	"fmt"
	"log"
	"sync"
)

type ActorSystem struct {
	actors           map[ActorRef]Mailbox
	actorsWaitGroup  sync.WaitGroup
	cleanerWaitGroup sync.WaitGroup
	deadActorsQueue  chan ActorRef
	lock             sync.Mutex
}

func NewActorSystem() ActorSystem {
	return ActorSystem{
		actors:           make(map[ActorRef]Mailbox),
		actorsWaitGroup:  sync.WaitGroup{},
		cleanerWaitGroup: sync.WaitGroup{},
		deadActorsQueue:  make(chan ActorRef, 1000),
	}
}

func (as *ActorSystem) Spawn(ref ActorRef, behavior ActorBehavior, mailboxSize int) (*Mailbox, error) {
	as.lock.Lock()
	if _, ok := as.actors[ref]; ok {
		return nil, fmt.Errorf("Actor already exists")
	}
	mailbox := make(Mailbox, mailboxSize)
	as.actors[ref] = mailbox
	as.lock.Unlock()

	behavior.Initialize(&ref, &mailbox, as)
	actor := Actor{Me: ref, Mailbox: &mailbox, Behavior: behavior, DeadQueue: as.deadActorsQueue}
	go actor.Run(&as.actorsWaitGroup)
	as.actorsWaitGroup.Add(1)
	return &mailbox, nil
}

func (as *ActorSystem) removeDeadActor(actorRef ActorRef) {
	as.lock.Lock()
	if mailbox, ok := as.actors[actorRef]; ok {
		close(mailbox)
		delete(as.actors, actorRef)
		log.Printf("[ActorSystem] Removed dead actor %v", actorRef)
	}
	as.lock.Unlock()
}

func (as *ActorSystem) ByRef(ref ActorRef) (*Mailbox, error) {
	mailbox, ok := as.actors[ref]
	if !ok {
		return nil, fmt.Errorf("Actor not found")
	}
	return &mailbox, nil
}

func (as *ActorSystem) Run() {
	go cleanDeadActors(as.deadActorsQueue, as)
	as.cleanerWaitGroup.Add(1)
	as.actorsWaitGroup.Wait()
	close(as.deadActorsQueue)
	as.cleanerWaitGroup.Wait()
}

func cleanDeadActors(deadActorsQueue chan ActorRef, actorSystem *ActorSystem) {
	for actorRef := range deadActorsQueue {
		actorSystem.removeDeadActor(actorRef)
	}
	actorSystem.cleanerWaitGroup.Done()
}
