package api

type GroupConfig struct {
	Heads []HeadID
	Name  string
}

type GroupID string

type Groups map[GroupID]Group

type Group struct {
	GroupConfig // TODO: separate config field
	ID          GroupID
	Heads       []HeadID
	Colors      Colors

	Intensity *Intensity
	Color     *Color
}

type GroupParams struct {
	Intensity *IntensityParams `json:",omitempty"`
	Color     *ColorParams     `json:",omitempty"`
}
