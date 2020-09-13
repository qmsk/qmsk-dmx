package api

import (
	"fmt"
)

type TypeID string
type HeadType struct {
	Vendor string
	Model  string
	Mode   string
	URL    string

	Channels []ChannelConfig
	Colors   Colors
}

func (headType HeadType) String() string {
	return fmt.Sprintf("%v/%v=%v", headType.Vendor, headType.Model, headType.Mode)
}

type HeadConfig struct {
	Type     TypeID
	Universe Universe
	Address  DMXAddress
	Name     string
	Count    uint // Clone multiple copies of the head at id.N
	Groups   []GroupID
}

type HeadID string

type Heads map[HeadID]Head

type Head struct {
	ID     HeadID
	Config HeadConfig
	Type   HeadType

	Channels  Channels   `json:",omitempty"`
	Intensity *Intensity `json:",omitempty"`
	Color     *Color     `json:",omitempty"`
}

type HeadParams struct {
	Channels  map[ChannelID]ChannelParams `json:",omitempty"`
	Intensity *IntensityParams            `json:",omitempty"`
	Color     *ColorParams                `json:",omitempty"`
}
