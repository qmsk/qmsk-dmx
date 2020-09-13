package heads

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/qmsk/dmx/api"
	"github.com/qmsk/dmx/logging"
)

func loadToml(obj interface{}, path string) error {
	if meta, err := toml.DecodeFile(path, obj); err != nil {
		return err
	} else if len(meta.Undecoded()) > 0 {
		return fmt.Errorf("Invalid keys in toml: %#v", meta.Undecoded())
	} else {
		return nil
	}
}

func load(obj interface{}, path string) error {
	ext := filepath.Ext(path)

	switch ext {
	case ".toml":
		if err := loadToml(obj, path); err != nil {
			return fmt.Errorf("Load %T file %v: %v", obj, path, err)
		} else {
			logging.Log.Infof("load %v: %T<%v>", path, obj, obj)

			return nil
		}
	default:
		return fmt.Errorf("Unknown %T file ext=%v: %v", obj, ext, path)
	}
}

type walkFunc func(path string, id []string) error

// Load config from given path, using the given stat info
func walkStat(path string, stat os.FileInfo, f walkFunc, prefix []string, top bool) error {
	name := stat.Name()

	if stat.IsDir() {
		logging.Log.Debugf("heads:loadOne path=%v: dir prefix=%v", path, prefix)

		var id []string

		if top {
			id = prefix
		} else {
			id = append(id, prefix...)
			id = append(id, name)
		}

		return walkDir(path, f, id)

	} else {
		// take basename *.ext part
		ext := filepath.Ext(name)
		name = name[:len(name)-len(ext)]

		var id []string

		if top {
			id = prefix
		} else {
			id = append(id, prefix...)
			id = append(id, name)
		}

		logging.Log.Debugf("heads:loadOne path=%v: file id=%v", path, id)

		if err := f(path, id); err != nil {
			return err
		} else {
			return nil
		}
	}
}

// Recursively load multiple config files from a directory
func walkDir(dirPath string, f walkFunc, prefix []string) error {
	if files, err := ioutil.ReadDir(dirPath); err != nil {
		return fmt.Errorf("read dir %v: %v", dirPath, err)
	} else {
		for _, stat := range files {
			path := filepath.Join(dirPath, stat.Name())

			if !(stat.Mode().IsRegular() || stat.Mode().IsDir()) {
				logging.Log.Debugf("heads:loadsDir path=%v: skip irregular", path)

				continue
			}

			if strings.HasPrefix(stat.Name(), ".") {
				logging.Log.Debugf("heads:loadOne path=%v: skip dotfile", path)

				continue
			}

			if err := walkStat(path, stat, f, prefix, false); err != nil {
				return err
			}
		}

		return nil
	}
}

// Load config from path, which may either be a file, or a directory to be loaded recursively
func walk(path string, f walkFunc) error {
	if stat, err := os.Stat(path); err != nil {
		return err
	} else {
		return walkStat(path, stat, f, nil, true)
	}
}

type loadConfig api.Config

func (loadConfig *loadConfig) loadTypes(path string) error {
	return walk(path, func(path string, id []string) error {
		var typeID = api.TypeID(filepath.Join(id...))
		var headType api.HeadType

		if err := load(&headType, path); err != nil {
			return err
		} else {
			loadConfig.HeadTypes[typeID] = headType
		}
	})
}

func (loadConfig *loadConfig) load(path string) error {
	return walk(path, func(path string, id []string) error {
		var typ string
		var name string

		if len(id) > 0 {
			typ = id[0]
		}
		if len(id) > 1 {
			typ = typ + "/*"
			name = filepath.Join(id[1:]...)
		}

		switch typ {
		case "":
			return load(loadConfig, path)

		case "colors":
			return load(&loadConfig.Colors, path)

		case "colors/*":
			var color api.Color

			if err := load(&color, path); err != nil {
				return err
			} else {
				loadConfig.Colors[api.ColorID(name)] = color
			}

		case "heads":
			return load(&loadConfig.Heads, path)

		case "heads/*":
			var headConfig api.HeadConfig

			if err := load(&headConfig, path); err != nil {
				return err
			} else {
				loadConfig.Heads[api.HeadID(name)] = headConfig
			}

		case "groups":
			return load(&loadConfig.Groups, path)

		case "groups/*":
			var groupConfig api.GroupConfig

			if err := load(&groupConfig, path); err != nil {
				return err
			} else {
				loadConfig.Groups[api.GroupID(name)] = groupConfig
			}

		case "presets":
			return load(&loadConfig.Presets, path)

		case "presets/*":
			var presetConfig api.PresetConfig

			if err := load(&presetConfig, path); err != nil {
				return err
			} else {
				loadConfig.Presets[api.PresetID(name)] = presetConfig
			}

		default:
			return fmt.Errorf("Unkonwn config %v: %v", id, path)
		}

		return nil
	})
}

// map relative Head.Type= references
func loadHeadType(config api.Config, headConfig api.HeadConfig) (api.HeadType, error) {
	if headType, exists := config.HeadTypes[headConfig.Type]; !exists {
		return headType, fmt.Errorf("Unknown Type=%v", headConfig.Type)
	} else {
		// merge over global colors
		headType.Colors = config.Colors.Merge(headType.Colors)

		return headType, nil
	}
}

// Config
func loadHeadID(headID api.HeadID, index int) api.HeadID {
	return api.HeadID(fmt.Sprintf("%s.%d", headID, index+1))
}

func loadHeadConfig(headConfig api.HeadConfig, index int, headType api.HeadType) api.HeadConfig {
	var step = len(headType.Channels)

	headConfig.Address = headConfig.Address + api.DMXAddress(index*step)

	return headConfig
}

func loadHeads(config api.Config, f func(headID api.HeadID, headConfig api.HeadConfig, headType api.HeadType) error) error {
	for headID, headConfig := range config.Heads {
		var err error

		if headType, headTypeErr := loadHeadType(config, headConfig); headTypeErr != nil {
			err = headTypeErr
		} else if headConfig.Count == 0 {
			err = f(headID, headConfig, headType)
		} else {
			for index := 0; index < int(headConfig.Count); index++ {
				err = f(loadHeadID(headID, index), loadHeadConfig(headConfig, index, headType), headType)
			}
		}

		if err != nil {
			return fmt.Errorf("Load head %v: %v", headID, err)
		}
	}

	return nil
}
