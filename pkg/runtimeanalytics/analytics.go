// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// LICENSE file in the root directory of this source tree.
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
