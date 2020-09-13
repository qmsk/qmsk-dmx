package heads

import (
	"fmt"

	"github.com/qmsk/dmx/api"
	"github.com/qmsk/dmx/logging"
)

type Options struct {
	Log logging.Option `long:"log.heads"`

	LibraryPath []string `long:"heads-library" value-name:"PATH"`
}

func (options Options) LoadConfig(path string) (api.Config, error) {
	var config loadConfig

	for _, libraryPath := range options.LibraryPath {
		if err := config.loadTypes(libraryPath); err != nil {
			return api.Config(config), fmt.Errorf("loadTypes %v: %v", libraryPath, err)
		}
	}

	if err := config.load(path); err != nil {
		return api.Config(config), err
	}

	return api.Config(config), nil
}

func (options Options) NewController(config api.Config) (*Controller, error) {
	options.Log.Package = "heads"

	var controller = MakeController()

	controller.log = options.Log.Logger("package", "heads")
	controller.events.log = options.Log.Logger("events", nil)

	if err := controller.Load(config); err != nil {
		return nil, err
	}

	return &controller, nil
}
