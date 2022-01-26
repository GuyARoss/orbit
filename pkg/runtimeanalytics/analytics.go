package runtimeanalytics

import "time"

type RuntimeAnalytics struct {
	startTime time.Time
}

func (r *RuntimeAnalytics) StartCapture() {
	r.startTime = time.Now()
}

func (r *RuntimeAnalytics) StopCapture() float64 {
	return time.Since(r.startTime).Seconds()
}
