package theater

type SyncerActor struct {
	me         ActorRef
	dispatcher MessageDispatcher
	system     *ActorSystem

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

func (sa *SyncerActor) Initialize(me ActorRef, dispatcher MessageDispatcher, system *ActorSystem) {
	dispatcher.RegisterDefaultHandler(sa.OnResponse)

	sa.me = me
	sa.dispatcher = dispatcher
	sa.system = system
}

func (sa *SyncerActor) OnResponse(message Message) {
	sa.resultChan <- message
	close(sa.resultChan)
}

func (sa *SyncerActor) Run() {
	err := sa.dispatcher.Send(*sa.target, Message{
		Type: "RequestMessage",
		Content: RequestMessage{
			Type:      sa.msg.Type,
			Content:   sa.msg.Content,
			RespondTo: sa.me,
		},
	})
	if err != nil {
		return
	}

	sa.dispatcher.Receive()
}
