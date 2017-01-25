package controller

import (
  "math"

	"github.com/davidhinkes/gogeorge/modules"
)

type ActuatorCommand struct {
  Enabled bool
	YokeRoll float64
}

type SensorData struct {
	RollDegrees float64
}

type T struct {
	deltaT   float64
	rollRate modules.Rate
  rollAngleFilter modules.AngleFilter
  rollRatePID *modules.PID
}

func New(deltaT float64) *T {
	return &T{
		deltaT:          deltaT,
		rollAngleFilter: modules.AngleFilter{},
		rollRate:        modules.MakeRate(int(math.Ceil(0.5 / deltaT))),
    rollRatePID: modules.NewPID(-0.1, -0.1, -0.1, int(math.Ceil(0.5/deltaT))),
	}
}

func (t *T) Do(sensorData SensorData) ActuatorCommand {
  rollErrorDegrees := -sensorData.RollDegrees // wing-leveler
  rollRateCommandDegreesPerSecond := absMax(0.25*(rollErrorDegrees), 15)
	rollRateDegreesPerSecond := t.rollRate.Do(t.rollAngleFilter.Do(sensorData.RollDegrees))
  rollRateErrorDegreesPerSecond := rollRateCommandDegreesPerSecond - rollRateDegreesPerSecond
  return ActuatorCommand {
    Enabled: true,
    YokeRoll: t.rollRatePID.Do(rollRateErrorDegreesPerSecond),
  }
}

func absMax(x, max float64) float64 {
  if x > max {
    return max
  }
  if x < -max {
    return -max
  }
  return x
}
