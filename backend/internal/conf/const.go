package conf

import "time"

// logic constants
const (
	Addr = ":8080"

	// RunTimeout task is killed after it
	RunTimeout = time.Second * 60 * 10
	// KillTimeout wait after SIGINT, then process is killed
	KillTimeout = time.Second * 7
	// StopWord stops execution, after received in the websocked
	StopWord = "stop"
	// BinPath application which is launched
	BinPath = "/Users/eyukorovin/iprover-demo/iprover-stub/iprover"
)
