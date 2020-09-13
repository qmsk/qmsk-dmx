package heads

import (
	"github.com/qmsk/dmx"
	"github.com/qmsk/dmx/api"
	"github.com/qmsk/dmx/logging"
	"github.com/qmsk/go-web"
)

type heads map[api.HeadID]*Head

func (heads heads) makeAPI() api.Heads {
	var apiHeads = make(api.Heads)

	for headID, head := range heads {
		apiHeads[headID] = head.makeAPI()
	}
	return apiHeads
}

func (heads heads) makeAPIList() []api.Head {
	var apiHeads = make([]api.Head, 0, len(heads))

	for _, head := range heads {
		apiHeads = append(apiHeads, head.makeAPI())
	}

	return apiHeads
}

type headsView struct {
	heads heads
}

func (view headsView) Index(name string) (web.Resource, error) {
	if name == "" {
		return view, nil
	} else if head := view.heads[api.HeadID(name)]; head != nil {
		return headView{head: head}, nil
	} else {
		return nil, nil
	}
}

func (view headsView) GetREST() (web.Resource, error) {
	return view.heads.makeAPIList(), nil
}

// A single DMX receiver using multiple consecutive DMX channels from a base address within a single universe
type Head struct {
	log logging.Logger

	id       api.HeadID
	config   api.HeadConfig
	headType api.HeadType
	output   *Output
	events   Events
	groups   groups

	channels  channels
	intensity *HeadIntensity
	color     *HeadColor
}

func (head *Head) String() string {
	return string(head.id)
}

func (head *Head) Name() string {
	if head.config.Name != "" {
		return head.config.Name
	} else {
		return string(head.id)
	}
}

func (head *Head) init() {
	head.channels = make(channels)

	for index, channelConfig := range head.headType.Channels {
		var channel = &Channel{
			id:      channelConfig.ID(),
			config:  channelConfig,
			index:   uint(index),
			output:  head.output,
			address: dmx.Address(head.config.Address) + dmx.Address(index),
		}

		channel.init()

		head.channels[channel.id] = channel
	}

	// setup parameters
	if headIntensity := head.Intensity(); headIntensity.exists() {
		head.intensity = &headIntensity
	}
	if headColor := head.Color(); headColor.exists() {
		head.color = &headColor
	}
}

// Head is member of Group
func (head *Head) initGroup(group *Group) {
	head.groups[group.id] = group
}

func (head *Head) Channel(config api.ChannelConfig) *Channel {
	return head.channels[config.ID()]
}

func (head *Head) Intensity() HeadIntensity {
	return HeadIntensity{
		channel: head.Channel(api.ChannelConfig{Intensity: true}),
	}
}

func (head *Head) Color() HeadColor {
	return HeadColor{
		red:       head.Channel(api.ChannelConfig{Color: api.ChannelColorRed}),
		green:     head.Channel(api.ChannelConfig{Color: api.ChannelColorGreen}),
		blue:      head.Channel(api.ChannelConfig{Color: api.ChannelColorBlue}),
		intensity: head.Channel(api.ChannelConfig{Intensity: true}),
	}
}

func (head *Head) makeAPI() api.Head {
	var apiHead = api.Head{
		ID:     head.id,
		Config: head.config,
		Type:   head.headType,

		Channels: head.channels.makeAPI(),
	}

	if head.intensity != nil {
		var intensity = head.intensity.GetIntensity()

		apiHead.Intensity = &intensity
	}

	if head.color != nil {
		var color = head.color.GetColor()

		apiHead.Color = &color
	}

	return apiHead
}

func (head *Head) applyAPI(params api.HeadParams) error {
	head.log.Info("Apply: %#v", params)

	for channelID, channelParams := range params.Channels {
		if channel := head.channels[channelID]; channel == nil {
			return web.Errorf(422, "No channel for head %v: %v", head.id, channelID)
		} else {
			channel.SetChannel(channelParams)
		}
	}

	if params.Intensity != nil {
		if head.intensity == nil {
			return web.Errorf(422, "No intensity for head %v", head.id)
		} else {
			head.intensity.SetIntensity(*params.Intensity)
		}
	}

	if params.Color != nil {
		if head.color == nil {
			return web.Errorf(422, "No color for head %v", head.id)
		} else {
			head.color.SetColor(*params.Color)
		}
	}

	return nil
}

func (head *Head) update() error {
	head.log.Info("Apply")

	head.output.Refresh()

	head.events.update(api.Event{
		Heads:  api.Heads{head.id: head.makeAPI()},
		Groups: head.groups.makeAPI(),
	})

	return nil
}

// GET /heads/:id => api.Head
type headView struct {
	head   *Head
	params api.HeadParams
}

func (view *headView) Index(name string) (web.Resource, error) {
	switch name {
	case "":
		return view, nil
	case "channels":
		return &channelsView{view.head.channels}, nil
	case "intensity":
		return &intensityView{handler: view.head.intensity}, nil
	case "color":
		return &colorView{handler: view.head.color}, nil
	default:
		return nil, nil
	}
}

func (view *headView) IntoREST() interface{} {
	return &view.params
}

func (view *headView) GetREST() (web.Resource, error) {
	return view.head.makeAPI(), nil
}
func (view *headView) PostREST() (web.Resource, error) {
	if err := view.head.applyAPI(view.params); err != nil {
		return nil, err
	} else if err := view.head.update(); err != nil {
		return nil, err
	} else {
		return view.head.makeAPI(), nil
	}
}
