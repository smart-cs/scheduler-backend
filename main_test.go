package main_test

import (
	"os"
	"testing"
	"time"

	"github.com/nickwu241/schedulecreator-backend"
)

func runMainAndAssertExit(t *testing.T, shouldExit bool) {
	finished := make(chan bool)
	go func() {
		main.RunMain()
		finished <- true
	}()
	// 200 ms should be long enough to bootstrap the server
	time.Sleep(time.Millisecond * 200)

	select {
	case _ = <-finished:
		// Channel received a message, main exited.
		if !shouldExit {
			t.FailNow()
		}
	default:
		// Channel hasn't received a message, main hasn't exited.
		if shouldExit {
			t.FailNow()
		}
	}
}

func TestMain(t *testing.T) {
	t.Log("main should block forever because it starts a server")
	runMainAndAssertExit(t, false)
}

func TestMain_EnvPort(t *testing.T) {
	t.Log("main with $PORT set should block forever")
	os.Setenv("PORT", "4321")
	runMainAndAssertExit(t, false)
}

func TestMain_BadEnvPort(t *testing.T) {
	t.Log("main should exit if $PORT is set to a non-integer value")
	os.Setenv("PORT", "this value is not an integer")
	runMainAndAssertExit(t, true)
}
