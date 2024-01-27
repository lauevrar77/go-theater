package theater

import (
	"fmt"
	"github.com/google/uuid"
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

	dispatcher := NewMessageDispatcher(&mailbox, as)
	behavior.Initialize(ref, dispatcher, as)
	actor := Actor{Me: ref, Behavior: behavior, DeadQueue: as.deadActorsQueue}
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
		return nil, fmt.Errorf("Actor %v not found", ref)
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

func (as *ActorSystem) Call(target *ActorRef, msg Message) (*Message, error) {
	resultChan := make(chan Message, 1)
	syncerName := fmt.Sprintf("syncer-%s-%s", string(*target), uuid.New().String())
	syncer := NewSyncerActor(target, msg, resultChan)
	_, err := as.Spawn(
		ActorRef(syncerName),
		syncer,
		1,
	)
	if err != nil {
		return nil, err
	}
	result := <-resultChan
	return &result, nil
}

func cleanDeadActors(deadActorsQueue chan ActorRef, actorSystem *ActorSystem) {
	for actorRef := range deadActorsQueue {
		actorSystem.removeDeadActor(actorRef)
	}
	actorSystem.cleanerWaitGroup.Done()
}
