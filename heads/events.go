package heads

import (
	"github.com/qmsk/dmx/api"
	"github.com/qmsk/dmx/logging"
	"github.com/qmsk/go-web"
)

// WebSocket handler
func (controller *Controller) WebEvents() web.Events {
	controller.events.eventChan = make(chan web.Event)

	return web.MakeEvents(web.EventConfig{
		EventPush: controller.events.eventChan,
		StateFunc: func() web.State {
			controller.log.Info("WebEvents State request")

			// must be goroutine-safe
			return controller.makeAPI()
		},
	})
}

type eventBuilder api.Event

func (eventBuilder *eventBuilder) addHead(head *Head) {
	if eventBuilder.Heads == nil {
		eventBuilder.Heads = api.Heads{head.id: head.makeAPI()}
	} else {
		eventBuilder.Heads[head.id] = head.makeAPI()
	}
}

func (eventBuilder *eventBuilder) addHeads(heads heads) {
	if eventBuilder.Heads == nil {
		eventBuilder.Heads = heads.makeAPI()
	} else {
		for headID, head := range heads {
			eventBuilder.Heads[headID] = head.makeAPI()
		}
	}
}

func (eventBuilder *eventBuilder) addGroup(group *Group) {
	if eventBuilder.Groups == nil {
		eventBuilder.Groups = api.Groups{group.id: group.makeAPI()}
	} else {
		eventBuilder.Groups[group.id] = group.makeAPI()
	}
}

func (eventBuilder *eventBuilder) addGroups(groups groups) {
	if eventBuilder.Groups == nil {
		eventBuilder.Groups = groups.makeAPI()
	} else {
		for groupID, group := range groups {
			eventBuilder.Groups[groupID] = group.makeAPI()
		}
	}
}

type events struct {
	log       logging.Logger
	eventChan chan web.Event
}

// push update events via websocket
func (events *events) update(event api.Event) {
	if events.eventChan != nil {
		events.log.Infof("update")
		events.eventChan <- event
	}
}

type Events interface {
	update(event api.Event)
}
