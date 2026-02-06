package logger

import "testing"

func TestInit(t *testing.T) {
	Init(false)
	log := Get()
	if log == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestInitVerbose(t *testing.T) {
	Init(true)
	log := Get()
	if log == nil {
		t.Fatal("expected non-nil logger")
	}
}
