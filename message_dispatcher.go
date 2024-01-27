package theater

import (
	"log"
)

type RegisteredMessageHandler func(interface{})
type UnregisteredMessageHandler func(Message)
type RegisteredRequestHandler func(interface{}) Message

type MessageDispatcher struct {
	mailbox                *Mailbox
	system                 *ActorSystem
	messageHandlers        map[string]RegisteredMessageHandler
	requestMessageHandlers map[string]RegisteredRequestHandler
	defaultHandler         *UnregisteredMessageHandler
}

func NewMessageDispatcher(mailbox *Mailbox, as *ActorSystem) MessageDispatcher {
	return MessageDispatcher{
		mailbox:                mailbox,
		messageHandlers:        make(map[string]RegisteredMessageHandler),
		requestMessageHandlers: make(map[string]RegisteredRequestHandler),
		system:                 as,
	}
}

func (d *MessageDispatcher) RegisterMessageHandler(msgType string, handler RegisteredMessageHandler) {
	d.messageHandlers[msgType] = handler
}

func (d *MessageDispatcher) RegisterRequestMessageHandler(msgType string, handler RegisteredRequestHandler) {
	d.requestMessageHandlers[msgType] = handler
}

func (d *MessageDispatcher) RegisterDefaultHandler(handler UnregisteredMessageHandler) {
	d.defaultHandler = &handler
}

func (d *MessageDispatcher) Receive() {
	msg := <-*d.mailbox
	if msg.Type == "RequestMessage" {
		d.handleRequestMessage(msg)
	} else {
		d.handleMessage(msg)
	}
}

func (d *MessageDispatcher) TryReceive() {
	select {
	case msg := <-*d.mailbox:
		if msg.Type == "RequestMessage" {
			d.handleRequestMessage(msg)
		} else {
			d.handleMessage(msg)
		}
		break
	default:
	}
}

func (d *MessageDispatcher) Send(target ActorRef, msg Message) error {
	mailbox, err := d.system.ByRef(target)
	if err != nil {
		return err
	}
	*mailbox <- msg
	return nil
}

func (d *MessageDispatcher) handleMessage(msg Message) {
	if handler, ok := d.messageHandlers[msg.Type]; ok {
		handler(msg.Content)
	} else if d.defaultHandler != nil {
		(*d.defaultHandler)(msg)
	} else {
		log.Printf("[MessageDispatcher] No handler for message type %s (and not default handler provided)", msg.Type)
	}
}

func (d *MessageDispatcher) handleRequestMessage(msg Message) {
	req := msg.Content.(RequestMessage)
	if handler, ok := d.requestMessageHandlers[req.Type]; ok {
		response := handler(req.Content)
		mailbox, err := d.system.ByRef(req.RespondTo)
		if err != nil {
			log.Printf("[MessageDispatcher] No mailbox found for actor %s", req.RespondTo)
		}
		*mailbox <- response
	} else if d.defaultHandler != nil {
		(*d.defaultHandler)(msg)
	} else {
		log.Printf("[MessageDispatcher] No handler for request message type %s (and not default handler provided)", req.Type)
	}
}
