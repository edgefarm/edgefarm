package net

import (
	"time"
)

const (
	nodeLabelKey     = "e2e-net"                 // must match tag in manifest
	app1Name         = "e2e-net1"                // must match name in manifest
	comp1Name        = "publisher1"              // must match name in manifest
	net1Name         = "e2e-network"             // must match name in manifest
	streamName       = "export-stream-aggregate" // must match name in manifest
	pubId1           = 123                       // must match manifest
	testingNameSpace = "e2e-net"
	dsPollTimeout    = time.Minute * 5
	kubeConfig       = ""
)
