package ghist

// Percentile returns a float64 of the percentile of a given value in the histogram
func (h *Histogram) Percentile(value float64) (percentile float64) {
	var position uint64
	for i := h.Size - 1; i >= 0; i-- { // iterate in reverse to get a percentile
		if h.Bins[i].Count == 0 {
			continue
		}
		if value > h.Bins[i].Max {
			position += h.Bins[i].Count
		} else if value >= h.Bins[i].Min { // linear estimate of value's position in its bin
			pct := 0.5
			if h.Bins[i].Max-h.Bins[i].Min > 0.0 {
				pct = (value - h.Bins[i].Min) / (h.Bins[i].Max - h.Bins[i].Min)
			}
			position += uint64(float64(h.Bins[i].Count) * pct)
			break
		}
	}
	return float64(position) / float64(h.Count)
}

// Percentile32 returns Percentile() as a float32
func (h *Histogram) Percentile32(value float32) (percentile float32) {
	return float32(h.Percentile(float64(value)))
}

// Median returns an estimate of the distribution's median value by interpolating within
// the bucket that contains the median value
func (h *Histogram) Median() (median float64) {
	var (
		midPoint  uint64 = h.Count / 2
		seenCount uint64 = 0
	)

	for i := 0; i < h.Size; i++ {
		seenCount += h.Bins[i].Count
		if seenCount >= midPoint {
			if h.Bins[i].Count > 1 {
				offset := 1 - (float64(seenCount-midPoint) / float64(h.Bins[i].Count-1))
				return h.Bins[i].Max - offset*(h.Bins[i].Max-h.Bins[i].Min)
			}
			return h.Bins[i].Max
		}
	}
	return
}

// Median32 returns Median() as a float32
func (h *Histogram) Median32() (median float32) {
	return float32(h.Median())
}

// Mean returns the mean value of the histogram
func (h *Histogram) Mean() (mean float64) {
	if h.Count > 0 {
		mean = h.Sum / float64(h.Count)
	}
	return
}

// Mean32 returns Mean() as a float32
func (h *Histogram) Mean32() (mean float32) {
	return float32(h.Mean())
}

// Mode returns the most populated Bin in the histogram
func (h *Histogram) Mode() (mode Bin) {
	var (
		modeIndex        = 0
		maxCount  uint64 = 0
	)
	if h.Size > 0 {
		for i := 0; i < h.Size; i++ {
			if h.Bins[i].Count > maxCount {
				modeIndex = i
				maxCount = h.Bins[i].Count
			}
		}
		mode = h.Bins[modeIndex]
	}
	return
}
