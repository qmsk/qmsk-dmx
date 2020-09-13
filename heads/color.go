package heads

import (
	"github.com/qmsk/dmx/api"
	"github.com/qmsk/go-web"
)

type ColorHandler interface {
	GetColor() api.Color
	SetColor(api.ColorParams) api.Color
}

// Head.Color
type HeadColor struct {
	red       *Channel
	green     *Channel
	blue      *Channel
	intensity *Channel
}

func (it HeadColor) exists() bool {
	return it.red != nil || it.green != nil || it.blue != nil
}

func (hc HeadColor) GetColor() (color api.Color) {
	if hc.red != nil {
		color.Red = hc.red.GetValue()
	}
	if hc.green != nil {
		color.Green = hc.green.GetValue()
	}
	if hc.blue != nil {
		color.Blue = hc.blue.GetValue()
	}
	return
}

func (hc HeadColor) setColor(color api.Color) api.Color {
	if hc.red != nil {
		color.Red = hc.red.SetValue(color.Red)
	}
	if hc.green != nil {
		color.Green = hc.green.SetValue(color.Green)
	}
	if hc.blue != nil {
		color.Blue = hc.blue.SetValue(color.Blue)
	}
	return color
}

// Set color with intensity, using either head intensity channel or linear RGB scaling
func (hc HeadColor) setIntensity(color api.Color, intensity api.Value) api.Color {
	if hc.intensity != nil {
		hc.setColor(color)
		hc.intensity.SetValue(intensity)
		return color
	} else {
		return hc.setColor(color.Scale(intensity))
	}
}

func (hc HeadColor) SetColor(params api.ColorParams) api.Color {
	if params.ScaleIntensity != nil {
		return hc.setIntensity(params.Color, *params.ScaleIntensity)
	} else {
		return hc.setColor(params.Color)
	}
}

// Group.Color
type GroupColor struct {
	headColors map[api.HeadID]HeadColor
}

func (groupColor GroupColor) exists() bool {
	return len(groupColor.headColors) > 0
}

// Return one color for the group
func (groupColor GroupColor) GetColor() (color api.Color) {
	for _, headColor := range groupColor.headColors {
		// This works fine assuming they are all the same color :)
		return headColor.GetColor()
	}

	return
}

func (groupColor GroupColor) SetColor(params api.ColorParams) api.Color {
	var color api.Color

	for _, headColor := range groupColor.headColors {
		color = headColor.SetColor(params)
	}

	return color
}

// API
type colorView struct {
	handler ColorHandler
	params  api.ColorParams
}

func (view *colorView) IntoREST() interface{} {
	return &view.params
}

func (view *colorView) GetREST() (web.Resource, error) {
	return view.handler.GetColor(), nil
}

func (view *colorView) PostREST() (web.Resource, error) {
	return view.handler.SetColor(view.params), nil
}
