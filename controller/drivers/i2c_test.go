package drivers

import (
	"testing"
)

func TestI2c(t *testing.T) {
	var b Bus
	if b.f != nil {
		t.Error()
	}
}
