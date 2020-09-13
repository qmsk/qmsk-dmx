package heads

import (
	"net/http"

	"github.com/qmsk/dmx/api"
	"github.com/qmsk/go-web"
)

func (controller *Controller) WebAPI() web.API {
	return web.MakeAPI(controller)
}

func (controller *Controller) Index(name string) (web.Resource, error) {
	switch name {
	case "":
		return controller, nil
	case "groups":
		return &groupsView{controller.groups}, nil
	case "outputs":
		return &outputsView{controller.outputs}, nil
	case "heads":
		return &headsView{controller.heads}, nil
	case "presets":
		return controller.presets, nil
	default:
		return nil, nil
	}
}

func (controller *Controller) makeAPI() api.Index {
	return api.Index{
		Outputs: controller.outputs.makeAPI(),
		Heads:   controller.heads.makeAPI(),
		Groups:  controller.groups.makeAPI(),
		Presets: controller.presets.makeAPI(),
	}
}

func (controller *Controller) GetREST() (web.Resource, error) {
	return controller.makeAPI(), nil
}

func (controller *Controller) Apply() error {
	controller.log.Info("Apply")

	// Refresh DMX output
	if err := controller.Refresh(); err != nil {
		controller.log.Warn("Refresh: ", err)

		return err
	}

	return nil
}

func (controller *Controller) WebConfigPreset() http.Handler {
	return httpConfigPreset{controller}
}
