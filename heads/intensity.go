package heads

import (
	"github.com/qmsk/dmx/api"
	"github.com/qmsk/go-web"
)

type IntensityHandler interface {
	GetIntensity() api.Intensity
	SetIntensity(api.IntensityParams) api.Intensity
}

// Head.Intensity
type HeadIntensity struct {
	channel *Channel
}

func (it HeadIntensity) exists() bool {
	return it.channel != nil
}

func (it HeadIntensity) GetIntensity() api.Intensity {
	if it.channel != nil {
		return api.Intensity{it.channel.GetValue()}
	} else {
		return api.Intensity{}
	}
}

func (it HeadIntensity) SetIntensity(params api.IntensityParams) api.Intensity {
	var intensity = params.Intensity

	if params.ScaleIntensity != nil {
		intensity = intensity.Scale(*params.ScaleIntensity)
	}

	return api.Intensity{it.channel.SetValue(intensity.Intensity)}
}

// Group.Intensity
type GroupIntensity struct {
	heads map[api.HeadID]HeadIntensity
}

func (groupIntensity GroupIntensity) exists() bool {
	return len(groupIntensity.heads) > 0
}

func (groupIntensity GroupIntensity) GetIntensity() (intensity api.Intensity) {
	for _, headIntensity := range groupIntensity.heads {
		return headIntensity.GetIntensity()
	}
	return
}

func (groupIntensity GroupIntensity) SetIntensity(params api.IntensityParams) api.Intensity {
	var intensity api.Intensity

	for _, headIntensity := range groupIntensity.heads {
		intensity = headIntensity.SetIntensity(params)
	}
	return intensity
}

// Web API
type intensityView struct {
	handler IntensityHandler
	params  api.IntensityParams
}

func (view *intensityView) IntoREST() interface{} {
	return &view.params
}

func (view *intensityView) GetREST() (web.Resource, error) {
	return view.handler.GetIntensity(), nil
}

func (view *intensityView) PostREST() (web.Resource, error) {
	return view.handler.SetIntensity(view.params), nil
}
