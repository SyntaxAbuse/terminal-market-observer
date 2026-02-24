package engine // Package owns state + orchestrates analysis passes.

import ( // Imports required by the engine.
	"context" // Context controls cancellation and deadlines.
	"time"    // Time handles polling intervals and timestamps.
) // End imports.

type Candle struct { // Candle represents a single OHLCV point in time.
	Timestamp time.Time // Timestamp of the candle.
	Open      float64   // Open price.
	High      float64   // High price.
	Low       float64   // Low price.
	Close     float64   // Close price.
	Volume    float64   // Volume.
	Trades    int       // Trade count.
} // End Candle.

type Settings struct { // Settings are runtime-adjustable knobs.
	PercentThreshold  float64 // Percent tolerance for S/R clustering.
	MinConfirmations  int     // Minimum confirmations for S/R.
	MaxHistory        int     // Max candles to retain.
	PollInterval      time.Duration // Poll interval for live updates.
} // End Settings.

type Snapshot struct { // Snapshot is what the UI renders.
	Now          time.Time // Current time.
	Iteration    int       // Tick count.
	Trend        string    // Current trend classification.
	Pattern      string    // Pattern classification.
	ReversalHint string    // Reversal hint text.
	LastCandle   Candle    // Most recent candle.
	Series       []float64 // Series used for plotting.
	Support      []float64 // Support levels.
	Resistance   []float64 // Resistance levels.
	StartPrice   float64   // Starting price for percent diff.
} // End Snapshot.

type Source interface { // Source fetches candles from an external provider.
	FetchLatest(ctx context.Context) ([]Candle, error) // Fetch latest candles (usually 1).
	FetchBootstrap(ctx context.Context) ([]Candle, error) // Fetch historical candles for initial state.
} // End Source.

type Engine struct { // Engine owns state and runs updates.
	src      Source   // Data source implementation.
	cfg      Settings // Runtime settings.
	history  []Candle // Rolling candle history.
	series   []float64 // Derived series (close prices).
	start    float64  // Starting price for diff.
	iter     int      // Iteration counter.
} // End Engine.

func New(src Source, cfg Settings) *Engine { // New constructs an engine instance.
	return &Engine{ // Return initialized engine.
		src: src, // Assign data source.
		cfg: cfg, // Assign settings.
	} // End return.
} // End New.

func (e *Engine) Bootstrap(ctx context.Context) error { // Bootstrap loads initial history.
	candles, err := e.src.FetchBootstrap(ctx) // Fetch historical candles.
	if err != nil { // If fetching failed.
		return err // Return error.
	} // End error check.

	for _, c := range candles { // Append each candle into history.
		e.appendCandle(c) // Append with bounds checks.
	} // End loop.

	if len(e.series) > 0 { // If we have price data.
		e.start = e.series[len(e.series)-1] // Use last close as starting price.
	} // End start assignment.

	return nil // Successful bootstrap.
} // End Bootstrap.

func (e *Engine) Tick(ctx context.Context) (Snapshot, error) { // Tick fetches latest candles, updates state, returns snapshot.
	e.iter++ // Increment iteration counter.

	candles, err := e.src.FetchLatest(ctx) // Fetch latest candle(s).
	if err != nil { // If fetch failed.
		return Snapshot{}, err // Return empty snapshot + error.
	} // End error check.

	for _, c := range candles { // Append each new candle.
		e.appendCandle(c) // Append into rolling history + series.
	} // End loop.

	return e.analyze(), nil // Run analysis and return snapshot.
} // End Tick.

func (e *Engine) appendCandle(c Candle) { // appendCandle adds a candle with bounded history.
	e.history = append(e.history, c) // Append candle.
	e.series = append(e.series, c.Close) // Append close to series.

	if e.cfg.MaxHistory > 0 && len(e.history) > e.cfg.MaxHistory { // Enforce history cap.
		e.history = e.history[1:] // Drop oldest candle.
	} // End cap enforcement.

	if e.cfg.MaxHistory > 0 && len(e.series) > e.cfg.MaxHistory { // Enforce series cap too.
		e.series = e.series[1:] // Drop oldest close.
	} // End cap enforcement.
} // End appendCandle.
