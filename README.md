# terminal-market-observer
A real-time terminal UI (TUI) that ingests OHLCV candles, maintains a rolling history, and runs deterministic analysis passes (trend, support/resistance, candlestick patterns) while streaming results to a live dashboard.

<img width="1887" height="924" alt="image" src="https://github.com/user-attachments/assets/df429664-7fed-4cfd-bf16-4f3dc97f6ae6" />

This repository is designed as a systems-style example of:

streaming ingestion → normalization → stateful history

deterministic analysis pipeline over rolling windows

terminal UI composition (termdash)

clean shutdown, cancellation, and long-running process behavior

API boundary discipline (typed decode + explicit failure behavior)

# What it does

pulls candle snapshots on an interval

decodes responses into typed OHLCV structs

appends to rolling history with bounded memory

trend classification (uptrend / downtrend / consolidation)

candlestick pattern detection (engulfing, doji variants, hammer, morning star, marubozu, momentum candles)

support/resistance zone estimation via bounce confirmation logic

# Architecture

The system is intentionally separated into layers:

Data Source

remote API clients (HTTP requests, headers, decoding)

time-bounded history pull for bootstrapping

interval-based updates for live tracking

# State

rolling candle store ([]OHLCV with bounded capacity)

derived series for charting and calculations

# Analysis Pipeline

pure functions over history windows

pattern detectors return (bool, label) or typed results

aggregation logic chooses the strongest/most relevant signals

UI + Control Plane

termdash widgets for chart + log + inputs

runtime parameter tuning

keyboard control for pause/exit
