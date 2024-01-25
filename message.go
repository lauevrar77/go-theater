package theater

type Message struct {
	Type    string
	Content interface{}
}

type RequestMessage struct {
	Type      string
	Content   interface{}
	RespondTo ActorRef
}
