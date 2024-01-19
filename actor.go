package theater

import "sync"

type Actor struct {
	Mailbox  *Mailbox
	Behavior ActorBehavior
}

func (a *Actor) Run(wg *sync.WaitGroup) {
	a.Behavior.Run()
	wg.Done()
}

type ActorBehavior interface {
	Initialize(*ActorRef, *Mailbox, *ActorSystem)
	Run()
}
