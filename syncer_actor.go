package theater

type SyncerActor struct {
	me      *ActorRef
	mailbox *Mailbox
	system  *ActorSystem

	target     *ActorRef
	resultChan chan Message
	msg        Message
}

func NewSyncerActor(target *ActorRef, msg Message, resultChan chan Message) *SyncerActor {
	return &SyncerActor{
		target:     target,
		resultChan: resultChan,
		msg:        msg,
	}
}

func (sa *SyncerActor) Initialize(me *ActorRef, mailbox *Mailbox, system *ActorSystem) {
	sa.me = me
	sa.mailbox = mailbox
	sa.system = system
}

func (sa *SyncerActor) Run() {
	targetMailbox, err := sa.system.ByRef(*sa.target)
	if err != nil {
		return
	}
	*targetMailbox <- Message{
		Type: "RequestMessage",
		Content: RequestMessage{
			Type:      sa.msg.Type,
			Content:   sa.msg.Content,
			RespondTo: *sa.me,
		},
	}
	result := <-*sa.mailbox
	sa.resultChan <- result
	close(sa.resultChan)
}
