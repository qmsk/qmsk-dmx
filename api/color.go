package api

type ColorID string

type Colors map[ColorID]Color

// Merge in colors from given map, overriding colors in this map
func (colors Colors) Merge(merge Colors) Colors {
	var merged = make(Colors)

	for colorID, color := range colors {
		merged[colorID] = color
	}

	for colorID, color := range merge {
		merged[colorID] = color
	}

	return merged
}

type Color struct {
	Red   Value
	Green Value
	Blue  Value
}

func (color Color) IsZero() bool {
	return color.Red == 0.0 && color.Green == 0.0 && color.Blue == 0.0
}

// Linear RGB intensity scaling
func (color Color) Scale(intensity Value) Color {
	return Color{
		Red:   color.Red * intensity,
		Green: color.Green * intensity,
		Blue:  color.Blue * intensity,
	}
}

type ColorParams struct {
	ScaleIntensity *Value
	Color
}
