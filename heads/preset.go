package heads

import (
	"net/http"

	"github.com/qmsk/dmx/api"
	"github.com/qmsk/dmx/logging"
	"github.com/qmsk/go-web"

	"github.com/BurntSushi/toml"
)

// Config
type PresetParameters struct {
	intensityHandler IntensityHandler
	colorHandler     ColorHandler

	Intensity *api.IntensityParams
	Color     *api.ColorParams
}

func (presetParameters PresetParameters) IsZero() bool {
	if presetParameters.Intensity != nil && !presetParameters.Intensity.IsZero() {
		return false
	}
	if presetParameters.Color != nil && !presetParameters.Color.IsZero() {
		return false
	}
	return true
}

// Do the parameters set in this preset override any of those parameters set in the other preset?
func (presetParameters PresetParameters) Overrides(other PresetParameters) bool {
	if presetParameters.Intensity == nil || other.Intensity == nil {

	} else if !presetParameters.Intensity.Equals(*other.Intensity) {
		return true
	}

	if presetParameters.Color == nil || other.Color == nil {

	} else if !presetParameters.Color.Equals(*other.Color) {
		return true
	}

	return false
}

func (presetParameters PresetParameters) scale(scale api.Value) PresetParameters {
	if presetParameters.Intensity != nil {
		intensityParams := *presetParameters.Intensity
		intensityParams.ScaleIntensity = &scale

		presetParameters.Intensity = &intensityParams
	}

	if presetParameters.Color != nil {
		colorParams := *presetParameters.Color
		colorParams.ScaleIntensity = &scale

		presetParameters.Color = &colorParams
	}

	return presetParameters
}

func (params PresetParameters) apply() {
	if params.intensityHandler != nil && params.Intensity != nil {
		params.intensityHandler.SetIntensity(*params.Intensity)
	}

	if params.colorHandler != nil && params.Color != nil {
		params.colorHandler.SetColor(*params.Color)
	}
}

func (presetParameters PresetParameters) Set(params api.PresetParams) {

}

type presets map[api.PresetID]*Preset

func (presets presets) Get() api.Presets {
	var apiPresests = make(api.Presets)

	for presetID, preset := range presets {
		apiPresests[presetID] = preset.makeAPI()
	}

	return apiPresests
}

type presetsView struct {
	presets presets
}

func (view presetsView) GetREST() (web.Resource, error) {
	return view.presets.makeAPI(), nil
}

func (view presetsView) Index(name string) (web.Resource, error) {
	if name == "" {
		return view, nil
	} else if preset := view.presets[api.PresetID(name)]; preset != nil {
		return &presetView{preset: preset}, nil
	} else {
		return nil, nil
	}
}

type Preset struct {
	log    logging.Logger
	events Events

	id     api.PresetID
	config api.PresetConfig

	allHeads  heads
	allGroups groups
	groups    map[api.GroupID]PresetParameters
	heads     map[api.HeadID]PresetParameters
}

func (preset *Preset) initAll(heads heads, groups groups) {
	preset.allHeads = heads
	preset.allGroups = groups
}

func (preset *Preset) initGroup(group *Group, params api.PresetConfigParams) error {
	var groupParameters = PresetParameters{
		intensityHandler: group.intensity,
		colorHandler:     group.color,

		Intensity: params.Intensity,
		Color:     params.Color,
	}

	preset.groups[group.id] = groupParameters

	return nil
}

func (preset *Preset) initHead(head *Head, params api.PresetConfigParams) error {
	var headParameters = PresetParameters{
		intensityHandler: head.intensity,
		colorHandler:     head.color,

		Intensity: params.Intensity,
		Color:     params.Color,
	}

	preset.heads[head.id] = headParameters

	return nil
}

func (preset *Preset) Get() api.Preset {
	return api.Preset{
		ID:     preset.id,
		Config: preset.config,

		Groups: preset.config.Groups,
		Heads:  preset.config.Heads,
	}
}

func (preset *Preset) Set(params api.PresetParams) error {
	var event eventBuilder

	preset.log.Info("Apply")

	if allParams := preset.Config.All; allParams != nil {
		if params.Intensity != nil {
			allParams = allParams.Scale(*params.Intensity)
		}

		for _, head := range preset.allHeads {
			// all params are optional
			if allParams.Intensity != nil && head.intensity != nil {
				head.intensity.SetIntensity(*params.Intensity)
			}

			if allParams.Color != nil && head.color != nil {
				head.color.SetColor(*params.Color)
			}
		}

		// update everything
		event.addHeads(preset.allHeads)
		event.addGroups(preset.allGroups)
	}

	for _, groupParams := range preset.groups {
		if apiPresetParams.Intensity != nil {
			apiGroupParams = apiGroupParams.scaleIntensity(*apiPresetParams.Intensity)
		}

		groupParams.apply()

		event.addGroup(apiGroupParams.group)
		event.addHeads(apiGroupParams.group.heads)
	}

	for _, headParams := range preset.heads {
		if apiPresetParams.Intensity != nil {
			apiHeadParams = apiHeadParams.scaleIntensity(*apiPresetParams.Intensity)
		}

		if err := apiHeadParams.Apply(); err != nil {
			return err
		}

		event.addHead(apiHeadParams.head)
		event.addGroups(apiHeadParams.head.groups)
	}

	preset.events.update(event)

	return nil
}

// GET /config/preset.toml
type httpConfigPreset struct {
	controller *Controller
}

// Export a preset configuration from the current state
func (controller *Controller) ConfigPreset() api.PresetConfig {
	var allParameters = PresetParameters{
		Intensity: &APIIntensity{},
		Color:     &APIColor{},
	}

	var presetConfig = PresetConfig{
		All:    &allParameters,
		Groups: make(map[string]PresetParameters),
		Heads:  make(map[string]PresetParameters),
	}

	for groupID, group := range controller.groups {
		var presetParameters = PresetParameters{
			Intensity: group.intensity.makeAPI(),
			Color:     group.color.makeAPI(),
		}

		if presetParameters.Overrides(allParameters) {

		} else {
			continue
		}

		presetConfig.Groups[string(groupID)] = presetParameters
	}

	for headID, head := range controller.heads {
		var presetParameters = PresetParameters{
			Intensity: head.parameters.Intensity.makeAPI(),
			Color:     head.parameters.Color.makeAPI(),
		}

		var baseParameters = allParameters

		for groupID, _ := range head.groups {
			if groupParameters, exists := presetConfig.Groups[string(groupID)]; exists {
				baseParameters = groupParameters
			}
		}

		if presetParameters.Overrides(baseParameters) {

		} else {
			continue
		}

		presetConfig.Heads[string(headID)] = presetParameters
	}

	return presetConfig
}

func (httpConfigPreset httpConfigPreset) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/toml")

	var presetConfig = httpConfigPreset.heads.ConfigPreset()

	if err := toml.NewEncoder(w).Encode(presetConfig); err != nil {
		panic(err)
	}
}
