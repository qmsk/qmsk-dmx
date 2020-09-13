package api

type Intensity struct {
	Intensity Value
}

func (intensity Intensity) Scale(scale Value) Intensity {
	return Intensity{Intensity: intensity.Intensity * scale}
}

type IntensityParams struct {
	ScaleIntensity *Value
	Intensity
}
