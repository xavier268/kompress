package kompress

//=========== schedulers =============

// Common scheduling functions to decide when to update.
// Used for reading and writing.
// Should be deterministic (no time-based decisions !)

// A Scheduler defines how often the tree will be updated.
type Scheduler func(e *engine) bool

// SCAlways always update.
func SCAlways(e *engine) bool {
	return true
}

// SCNever never update
func SCNever(e *engine) bool {
	return false
}

// SCDelta updates if one of the frequency differs a lot ...
func SCDelta(e *engine) bool {
	const DELTA = 5
	for i, f := range e.freq {
		if e.actfreq[i] > f+DELTA {
			return true
		}
	}
	return false
}
