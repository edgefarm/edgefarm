package apps

import (
	"time"
)

const (
	nodeLabelKey     = "e2e"             // must match tag in manifest
	appName          = "e2e-app"         // must match name in manifest
	appName2         = "e2e-app2"        // must match name in manifest
	podNamePrefix1   = "test-component1" // must match name in manifest
	podNamePrefix2   = "test-component2" // must match name in manifest
	podNamePrefix3   = "test-component3" // must match name in manifest
	testingNameSpace = "default"         // must match name in manifest
	dsPollTimeout    = time.Minute * 5
	kubeConfig       = ""
)
