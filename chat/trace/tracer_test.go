package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer.out == nil {
		t.Error("Return fron New should not be nil")
	} else {
		tracer.Trace("Hello trace package.")
		if buf.String() != "Hello trace package.\n" {
			t.Errorf("Output wrong text: '%s'", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer
	silentTracer.Trace("something")
}
