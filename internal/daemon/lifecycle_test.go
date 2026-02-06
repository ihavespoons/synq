package daemon

import (
	"os"
	"testing"
)

func TestWriteAndReadPID(t *testing.T) {
	tmp := t.TempDir()

	if err := WritePID(tmp); err != nil {
		t.Fatal(err)
	}

	pid, err := ReadPID(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if pid != os.Getpid() {
		t.Errorf("PID = %d, want %d", pid, os.Getpid())
	}
}

func TestReadPID_NotExist(t *testing.T) {
	tmp := t.TempDir()
	pid, err := ReadPID(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if pid != 0 {
		t.Errorf("PID = %d, want 0", pid)
	}
}

func TestIsRunning(t *testing.T) {
	tmp := t.TempDir()

	// No PID file -> not running.
	running, _ := IsRunning(tmp)
	if running {
		t.Error("expected not running")
	}

	// Write current PID -> running.
	if err := WritePID(tmp); err != nil {
		t.Fatal(err)
	}
	running, pid := IsRunning(tmp)
	if !running {
		t.Error("expected running")
	}
	if pid != os.Getpid() {
		t.Errorf("PID = %d, want %d", pid, os.Getpid())
	}
}

func TestRemovePID(t *testing.T) {
	tmp := t.TempDir()
	if err := WritePID(tmp); err != nil {
		t.Fatal(err)
	}

	RemovePID(tmp)

	pid, err := ReadPID(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if pid != 0 {
		t.Errorf("PID = %d after removal, want 0", pid)
	}
}
