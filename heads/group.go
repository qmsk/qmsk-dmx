package heads

import (
	"github.com/qmsk/dmx/api"
	"github.com/qmsk/dmx/logging"
	"github.com/qmsk/go-web"
)

type groups map[api.GroupID]*Group

func (groups groups) makeAPI() api.Groups {
	apiGroups := make(api.Groups)

	for groupID, group := range groups {
		apiGroups[groupID] = group.makeAPI()
	}

	return apiGroups
}

func (groups groups) makeAPIList() (apiGroups []api.Group) {
	for _, group := range groups {
		apiGroups = append(apiGroups, group.makeAPI())
	}

	return
}

type groupsView struct {
	groups groups
}

func (view groupsView) GetREST() (web.Resource, error) {
	return view.groups.makeAPI(), nil
}

func (view groupsView) Index(name string) (web.Resource, error) {
	if name == "" {
		return view, nil
	} else if group := view.groups[api.GroupID(name)]; group != nil {
		return groupView{group: group}, nil
	} else {
		return nil, nil
	}
}

// Group
type Group struct {
	log     logging.Logger
	id      api.GroupID
	config  api.GroupConfig
	heads   heads
	events  Events
	outputs outputs

	intensity *GroupIntensity
	color     *GroupColor
	colors    api.Colors
}

func (group *Group) addHead(head *Head) {
	group.heads[head.id] = head
}

// initialize group parameters from heads
func (group *Group) init() {
	// reverse-mappings for apply updates
	for _, head := range group.heads {
		head.initGroup(group)
	}

	if groupIntensity := group.makeIntensity(); groupIntensity.exists() {
		group.intensity = &groupIntensity
	}

	if groupColor := group.makeColor(); groupColor.exists() {
		group.color = &groupColor
	}

	// merge head ColorMaps
	for _, head := range group.heads {
		if colorMap := head.headType.Colors; colorMap != nil {
			group.colors.Merge(colorMap)
		}
	}
}

func (group *Group) makeIntensity() GroupIntensity {
	var groupIntensity = GroupIntensity{
		heads: make(map[api.HeadID]HeadIntensity),
	}

	for headID, head := range group.heads {
		if headIntensity := head.intensity; headIntensity != nil {
			groupIntensity.heads[headID] = *headIntensity
		}
	}

	return groupIntensity
}

func (group *Group) makeColor() GroupColor {
	var groupColor = GroupColor{
		headColors: make(map[api.HeadID]HeadColor),
	}

	for headID, head := range group.heads {
		if headColor := head.color; headColor != nil {
			groupColor.headColors[headID] = *headColor
		}
	}

	return groupColor
}

func (group *Group) makeAPIHeads() []api.HeadID {
	var heads = make([]api.HeadID, 0)

	for headID, _ := range group.heads {
		heads = append(heads, headID)
	}
	return heads
}

func (group *Group) makeAPI() api.Group {
	var apiGroup = api.Group{
		GroupConfig: group.config,
		ID:          group.id,
		Heads:       group.makeAPIHeads(),
		Colors:      group.colors,
	}

	if group.intensity != nil {
		var intensity = group.intensity.GetIntensity()

		apiGroup.Intensity = &intensity
	}

	if group.color != nil {
		var color = group.color.GetColor()

		apiGroup.Color = &color
	}

	return apiGroup
}

func (group *Group) applyAPI(params api.GroupParams) error {
	if params.Intensity != nil {
		group.intensity.SetIntensity(*params.Intensity)
	}

	if params.Color != nil {
		group.color.SetColor(*params.Color)
	}

	return nil
}

func (group *Group) update() error {
	group.log.Info("Apply")

	group.outputs.Refresh()

	group.events.update(api.Event{
		Heads:  group.heads.makeAPI(),
		Groups: api.Groups{group.id: group.makeAPI()},
	})

	return nil
}

type groupView struct {
	group  *Group
	params api.GroupParams
}

func (view *groupView) IntoREST() interface{} {
	return &view.params
}

func (view *groupView) GetREST() (web.Resource, error) {
	return view.group.makeAPI(), nil
}

func (view *groupView) PostREST() (web.Resource, error) {
	if err := view.group.applyAPI(view.params); err != nil {
		return nil, err
	} else if err := view.group.update(); err != nil {
		return nil, err
	} else {
		return view.group.makeAPI(), nil
	}
}
