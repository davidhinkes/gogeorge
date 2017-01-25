package controller

import (
  "testing"
)

func TestController(t *testing.T) {
  New(0.1).Do(SensorData{})
}
