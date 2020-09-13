package api

type PresetConfigParams struct {
	Intensity *IntensityParams
	Color     *ColorParams
}

func (params PresetConfigParams) Scale(scale Value) PresetConfigParams {
	var ret PresetConfigParams

	if params.Intensity != nil {
		var intensity = *params.Intensity

		ret.Intensity = &intensity
		ret.Intensity.ScaleIntensity = &scale
	}

	if params.Color != nil {
		var color = *params.Color

		ret.Color = &color
		ret.Color.ScaleIntensity = &scale
	}

	return ret
}

type PresetGroups map[GroupID]PresetConfigParams
type PresetHeads map[HeadID]PresetConfigParams

type PresetConfig struct {
	Name   string
	All    *PresetConfigParams
	Groups PresetGroups
	Heads  PresetHeads
}

type PresetID string
type Presets map[PresetID]Preset

type Preset struct {
	ID     PresetID
	Config PresetConfig

	// TODO: superfluous re config?
	Groups PresetGroups
	Heads  PresetHeads
}

type PresetParams struct {
	Intensity *Value // TODO: rename to Scale
}
