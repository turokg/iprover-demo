package internal

import "time"

const (
	Addr = ":8080"

	// RunTimeout task is killed after it
	RunTimeout = time.Second * 10 * 60
	// KillTimeout wait after SIGINT, then process is killed
	KillTimeout = time.Second * 700
	// StopWord stops execution, after received in the websocked
	StopWord = "stop"
	// BinPath application which is launched
	//BinPath = "/Users/eyukorovin/iprover/iprover-demo/iprover-stub/iprover"
	BinPath     = "/Users/eyukorovin/iprover/iprover/iproveropt"
	ProblemsDir = "/Users/eyukorovin/iprover/iprover-demo/backend/problems/"

	LaunchBuffer = 1000 //TODO посмотри, наверно надо увеличить
)
