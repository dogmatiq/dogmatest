package location_test

// This file contains definitions used within tests that check for specific line
// numbers. To minimize test disruption edit this file as infrequently as
// possible.
//
// New additions should always be made at the end so that the line numbers of
// existing definitions do not change. The padding below can be removed as
// imports statements added.

import (
	. "github.com/dogmatiq/testkit/location"
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
	// import padding
)

func doNothing()             {}
func doPanic()               { panic("<panic>") }
func ofCallLayer1() Location { return OfCall() }
func ofCallLayer2() Location { return ofCallLayer1() }

type ofMethodT struct{}

func (ofMethodT) Method() {}
