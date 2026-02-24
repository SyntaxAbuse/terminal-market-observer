package engine // Same package; this file focuses on analysis.

import ( // Imports for analysis.
	"math" // Used for percent diff.
	"time" // Used for timestamps.
) // End imports.

func (e *Engine) analyze() Snapshot { // analyze derives the renderable snapshot.
	now := time.Now() // Capture current time.

	last := Candle{} // Default last candle.
	if len(e.history) > 0 { // If candles exist.
		last = e.history[len(e.history)-1] // Take most recent candle.
	} // End last candle selection.

	// These calls are placeholders for your existing functions refactored into internal/analysis. // Comment boundary.
	trend := classifyTrend(e.history) // Compute trend (pure).
	pattern := classifyPattern(e.history) // Compute pattern (pure).
	reversal := detectReversal(e.history) // Compute reversal (pure).
	support := findSupport(e.series, e.cfg.PercentThreshold, e.cfg.MinConfirmations) // Compute support (pure).
	resist := findResistance(e.series, e.cfg.PercentThreshold, e.cfg.MinConfirmations) // Compute resistance (pure).

	diff := 0.0 // Default diff.
	pct := 0.0  // Default percent diff.

	if e.start > 0 { // If we have a valid starting price.
		diff = math.Abs(e.start - last.Close) // Compute absolute delta.
		pct = (diff / e.start) * 100.0        // Compute percent delta.
	} // End start guard.

	_ = pct // pct can be included in Snapshot if you want; keeping this minimal here. // Comment for clarity.

	return Snapshot{ // Build snapshot for UI.
		Now:          now,      // Set current time.
		Iteration:    e.iter,    // Set iteration.
		Trend:        trend,     // Set trend string.
		Pattern:      pattern,   // Set pattern string.
		ReversalHint: reversal,  // Set reversal hint.
		LastCandle:   last,      // Set last candle.
		Series:       append([]float64(nil), e.series...), // Copy series to avoid external mutation.
		Support:      support,   // Set support levels.
		Resistance:   resist,    // Set resistance levels.
		StartPrice:   e.start,   // Set starting price.
	} // End snapshot.
} // End analyze.
