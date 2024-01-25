package theater

import "sync"

type Actor struct {
	Me        ActorRef
	Mailbox   *Mailbox
	Behavior  ActorBehavior
	DeadQueue chan ActorRef
}

func (a *Actor) Run(wg *sync.WaitGroup) {
	a.Behavior.Run()
	a.DeadQueue <- a.Me
	wg.Done()
}

type ActorBehavior interface {
	Initialize(*ActorRef, *Mailbox, *ActorSystem)
	Run()
}
