package log

import (
	"os"
	"sync"
	"time"
)

type flusher struct {
	interval time.Duration
	ticker   *time.Ticker
	ch       chan time.Duration
	sync.Mutex
}

func (f *flusher) start(lfi time.Duration) {
	f.interval = f.envOverride(lfi)
	f.ticker = time.NewTicker(f.interval)
	f.ch = make(chan time.Duration)

	go f.run()
}

func (f *flusher) run() {
	for {
		select {
		case d := <-f.ch:
			f.replaceTicker(d)
		case <-f.ticker.C:
			f.flush()
		}
	}
}

func (f *flusher) updateFlushInterval(d time.Duration) {
	f.ch <- d
}

func (f *flusher) replaceTicker(d time.Duration) {
	f.Lock()
	defer f.Unlock()

	if d == f.interval {
		return
	}

	Infof("Adjusting flush interval: %v -> %v", f.interval, d)

	f.interval = d

	f.ticker.Stop()

	f.ticker = time.NewTicker(d)
}

func (f *flusher) flush() {
	f.Lock()
	defer f.Unlock()

	logging.lockAndFlushAll()
}

func (f *flusher) envOverride(dflt time.Duration) time.Duration {
	v := os.Getenv("LOG_FLUSH_INTERVAL")
	if v == "" {
		return dflt
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		os.Stderr.Write([]byte("WARNING: Cannot parse $LOG_FLUSH_INTERVAL value \"" + v + "\": " + err.Error()))
		os.Stderr.Sync()
		os.Setenv("LOG_FLUSH_INTERVAL", "")
		return dflt
	}

	return d
}
