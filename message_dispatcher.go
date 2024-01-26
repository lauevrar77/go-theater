package theater

import "log"

type MessageHandler func(interface{})
type RequestMessageHandler func(interface{}) Message

type MessageDispatcher struct {
	mailbox                *Mailbox
	system                 *ActorSystem
	messageHandlers        map[string]MessageHandler
	requestMessageHandlers map[string]RequestMessageHandler
	defaultHandler         *MessageHandler
}

func NewMessageDispatcher(mailbox *Mailbox, as *ActorSystem) MessageDispatcher {
	return MessageDispatcher{
		mailbox:                mailbox,
		messageHandlers:        make(map[string]MessageHandler),
		requestMessageHandlers: make(map[string]RequestMessageHandler),
		system:                 as,
	}
}

func (d *MessageDispatcher) RegisterMessageHandler(msgType string, handler MessageHandler) {
	d.messageHandlers[msgType] = handler
}

func (d *MessageDispatcher) RegisterRequestMessageHandler(msgType string, handler RequestMessageHandler) {
	d.requestMessageHandlers[msgType] = handler
}

func (d *MessageDispatcher) RegisterDefaultHandler(handler MessageHandler) {
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

func (d *MessageDispatcher) handleMessage(msg Message) {
	if handler, ok := d.messageHandlers[msg.Type]; ok {
		handler(msg.Content)
	} else if d.defaultHandler != nil {
		(*d.defaultHandler)(msg.Content)
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
		(*d.defaultHandler)(msg.Content)
	} else {
		log.Printf("[MessageDispatcher] No handler for request message type %s (and not default handler provided)", req.Type)
	}
}
