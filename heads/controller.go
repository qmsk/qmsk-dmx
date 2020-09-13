package heads

import (
	"fmt"
	"github.com/qmsk/dmx"
	"github.com/qmsk/dmx/api"
	"github.com/qmsk/dmx/logging"
)

func MakeController() Controller {
	return Controller{
		outputs: make(outputs),
		heads:   make(heads),
		groups:  make(groups),
		presets: make(presets),
		events:  &events{},
	}
}

type Controller struct {
	log     logging.Logger
	outputs outputs
	heads   heads
	groups  groups
	presets presets
	events  *events // optional
}

func (controller *Controller) output(universe api.Universe) *Output {
	output := controller.outputs[universe]
	if output == nil {
		output = &Output{
			events: controller.events,
		}

		output.init(controller.log, universe)

		controller.outputs[universe] = output
	}

	return output
}

// Patch output
// XXX: not goroutine-safe...
func (controller *Controller) Output(universe api.Universe, config OutputConfig, writer dmx.Writer) {
	controller.output(universe).connect(config, writer)
}

func (controller *Controller) Load(config api.Config) error {
	// preload groups
	for groupID, groupConfig := range config.Groups {
		controller.addGroup(groupID, groupConfig)
	}

	if err := loadHeads(config, func(headID api.HeadID, headConfig api.HeadConfig, headType api.HeadType) error {
		controller.addHead(headID, headConfig, headType)

		return nil
	}); err != nil {
		return err
	}

	// once all heads are patched, init groups
	for _, group := range controller.groups {
		group.init()
	}

	//  load presets
	for presetID, presetConfig := range config.Presets {
		if err := controller.addPreset(presetID, presetConfig); err != nil {
			return fmt.Errorf("Load preset %v: %v", presetID, err)
		}
	}

	return nil
}

func (controller *Controller) addGroup(id api.GroupID, config api.GroupConfig) *Group {
	group := controller.groups[id]

	if group == nil {
		group = &Group{
			//XXX: log:    controller.options.Log.Logger("group", id),
			id:      id,
			config:  config,
			heads:   make(heads),
			colors:  make(api.Colors),
			events:  controller.events,
			outputs: controller.outputs,
		}

		controller.groups[id] = group
	}

	return group
}

func (controller *Controller) group(id api.GroupID) *Group {
	if group := controller.groups[id]; group == nil {
		return controller.addGroup(id, api.GroupConfig{})
	} else {
		return group
	}
}

// Patch head
func (controller *Controller) addHead(id api.HeadID, config api.HeadConfig, headType api.HeadType) *Head {
	var output = controller.output(config.Universe)
	var head = Head{
		// XXX: log:      controller.options.Log.Logger("head", id),
		id:       id,
		config:   config,
		headType: headType,
		output:   output,
		events:   controller.events,
		groups:   make(groups),
	}

	// load head parameters
	head.init()

	// map controller
	controller.heads[id] = &head

	// map groups
	for _, groupID := range config.Groups {
		controller.group(groupID).addHead(&head)
	}

	return &head
}

func (controller *Controller) addPreset(id api.PresetID, config api.PresetConfig) error {
	var preset = Preset{
		// XXX: log:    controller.options.Log.Logger("preset", id),
		events: controller.events,
		id:     id,
		config: config,
		groups: make(map[api.GroupID]PresetParameters),
		heads:  make(map[api.HeadID]PresetParameters),
	}

	if config.All != nil {
		preset.initAll(controller.heads, controller.groups)
	}

	for groupID, params := range config.Groups {
		if group := controller.groups[api.GroupID(groupID)]; group == nil {
			return fmt.Errorf("Unknown group for preset %v: %v", id, groupID)
		} else {
			preset.initGroup(group, params)
		}
	}

	for headID, params := range config.Heads {
		if head := controller.heads[api.HeadID(headID)]; head == nil {
			return fmt.Errorf("Unknown head for preset %v: %v", id, headID)
		} else {
			preset.initHead(head, params)
		}
	}

	controller.presets[id] = &preset

	return nil
}

func (controller *Controller) Each(fn func(head *Head)) {
	for _, head := range controller.heads {
		fn(head)
	}
}

// refresh outputs
func (controller *Controller) RefreshOutputs() error {
	var refreshErr error

	for _, output := range controller.outputs {
		if err := output.Refresh(); err != nil {
			refreshErr = err
		}
	}

	return refreshErr
}
