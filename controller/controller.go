package controller

import (
  "fmt"
  "math"

	"github.com/davidhinkes/gogeorge/modules"
	"github.com/davidhinkes/gogeorge/internal/io"
)

type ActuatorCommand struct {
  Enabled bool
	YokeRoll float64
}

type SensorData struct {
	RollDegrees float64
  HeadingDegrees float64
}

type T struct {
	deltaT   float64
  headingDesiredDegrees *io.Float64
	rollRate modules.Rate
  rollAngleFilter modules.AngleFilter
  rollRatePID *modules.PID
  master *io.Bool
}

func New(deltaT float64) *T {
  name := func(s string) string {
    return fmt.Sprintf("controller/rollratepid/%v", s)
  }
	return &T{
		deltaT:          deltaT,
    master:         io.NewBool("controller/master", false),
    headingDesiredDegrees: io.NewFloat64("controller/heading_desired_degrees", 0),
		rollAngleFilter: modules.AngleFilter{},
		rollRate:        modules.MakeRate(int(math.Ceil(0.5 / deltaT))),
    rollRatePID: modules.NewPID(
      io.NewFloat64(name("p"), 0.01),
      io.NewFloat64(name("i"), 0.0005),
      io.NewFloat64(name("d"), 0.005),
      int(math.Ceil(0.5/deltaT))),
	}
}

func normalizeDegrees(x float64) float64 {
  for x > 180 {
    x -= 360
  }
  for x < -180 {
    x += 360
  }
  return x
}

func (t *T) Do(sensorData SensorData) ActuatorCommand {
  if !t.master.Get() {
    return ActuatorCommand{}
  }
  headingErrorDegrees := normalizeDegrees(t.headingDesiredDegrees.Get() - sensorData.HeadingDegrees)
  rollDesiredDegrees := absMax(headingErrorDegrees/2, 30)
  rollErrorDegrees := rollDesiredDegrees-sensorData.RollDegrees // wing-leveler
  rollRateCommandDegreesPerSecond := absMax(rollErrorDegrees, 5)
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
