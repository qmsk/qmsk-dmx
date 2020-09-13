package api

type ChannelColor string

const (
	ChannelColorRed   = "red"
	ChannelColorGreen = "green"
	ChannelColorBlue  = "blue"
)

type ChannelConfig struct {
	Control   string       `json:",omitempty"`
	Intensity bool         `json:",omitempty"`
	Color     ChannelColor `json:",omitempty"`
}

func (config ChannelConfig) String() string {
	if config.Control != "" {
		return "control:" + config.Control
	}
	if config.Intensity {
		return "intensity"
	}
	if config.Color != "" {
		return "color:" + string(config.Color)
	}

	return ""
}

func (config ChannelConfig) ID() ChannelID {
	return ChannelID(config.String())
}

type ChannelID string
type Channels map[ChannelID]Channel

type Channel struct {
	ID      ChannelID
	Config  ChannelConfig
	Index   uint
	Address DMXAddress

	DMX   DMXValue
	Value Value
}

type ChannelParams struct {
	DMX   *DMXValue `json:",omitempty"`
	Value *Value    `json:",omitempty"`
}
