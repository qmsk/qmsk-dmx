package heads

import (
	"github.com/qmsk/dmx"
	"github.com/qmsk/dmx/api"
	"github.com/qmsk/go-web"
)

// Channels
type channels map[api.ChannelID]*Channel

func (channels channels) makeAPI() api.Channels {
	var apiChannels = make(api.Channels)

	for id, channel := range channels {
		apiChannels[id] = channel.GetChannel()
	}

	return apiChannels
}

type channelsView struct {
	channels channels
}

func (view channelsView) Index(name string) (web.Resource, error) {
	if channel := view.channels[api.ChannelID(name)]; channel == nil {
		return nil, nil
	} else {
		return channel, nil
	}
}
func (view channelsView) GetREST() (web.Resource, error) {
	return view.channels.makeAPI(), nil
}

type Channel struct {
	id     api.ChannelID
	config api.ChannelConfig
	output *Output
	index  uint

	address dmx.Address
}

func (channel *Channel) init() {
	channel.output.SetDMX(channel.address, 0)
}

func (channel *Channel) GetDMX() api.DMXValue {
	return api.DMXValue(channel.output.GetDMX(channel.address))
}
func (channel *Channel) GetValue() api.Value {
	return api.Value(channel.output.GetValue(channel.address))
}

func (channel *Channel) SetDMX(value api.DMXValue) {
	channel.output.SetDMX(channel.address, dmx.Channel(value))
}
func (channel *Channel) SetValue(value api.Value) api.Value {
	return api.Value(channel.output.SetValue(channel.address, Value(value)))
}

func (channel *Channel) SetChannel(params api.ChannelParams) {
	if params.DMX != nil {
		channel.SetDMX(*params.DMX)
	}
	if params.Value != nil {
		channel.SetValue(*params.Value)
	}
}

func (channel *Channel) GetChannel() api.Channel {
	return api.Channel{
		ID:      channel.id,
		Index:   channel.index,
		Config:  channel.config,
		Address: api.DMXAddress(channel.address),
		DMX:     channel.GetDMX(),
		Value:   channel.GetValue(),
	}
}

// API
type channelView struct {
	channel *Channel
	params  api.ChannelParams
}

func (view *channelView) GetREST() (web.Resource, error) {
	return view.channel.GetChannel(), nil
}

func (view *channelView) IntoREST() interface{} {
	return &view.params
}

func (view *channelView) PostREST() (web.Resource, error) {
	view.channel.SetChannel(view.params)

	return view.channel.GetChannel(), nil
}
