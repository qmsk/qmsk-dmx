package api

type Config struct {
	HeadTypes map[TypeID]HeadType
	Colors    Colors

	Heads  map[HeadID]HeadConfig
	Groups map[GroupID]GroupConfig

	Presets map[PresetID]PresetConfig
}
