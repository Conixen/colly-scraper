package main

import (
	"fmt"
	"strings"
	"time"
)

// Use colly.Limit to avoid IP bans and respect server load.
// Target only public stats, prices, and news to stay compliant.
// Rotate User-Agents to mimic real browser behavior.

// Target: Yahoo Finance (Finance) - Best for stock tickers and historical data.