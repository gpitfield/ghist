/*
Package ghist implements a streaming histogram.

ghist generates streaming histograms of float32 values in specified bin counts.
For most applications, a bin count <100 is likely sufficient.
Based loosely on http://jmlr.org/papers/volume11/ben-haim10a/ben-haim10a.pdf.
*/
package ghist

import (
	"math"
	"sort"
)

// Histogram maintains a distribution of Size populated bins
type Histogram struct {
	Count uint64
	bins  []Bin
	Size  int
}

// Bin keeps track of a histogram bin's minimum and maximum values and count
type Bin struct {
	max   float64
	min   float64
	count uint64
}

var zero = Bin{} // for empty comparisons

// Add adds a float64 value to the histogram, modifying it as necessary
func (h *Histogram) Add(value float64) {
	if h.Count < math.MaxUint64 {
		h.Count += 1
	} else {
		panic("Integer overflow: Maximum count exceeded in ghist Histogram")
	}
	// see if it fits in an existing bin
	index := sort.Search(len(h.bins), func(i int) bool { return value >= h.bins[i].min })
	if index < len(h.bins) && h.bins[index].max >= value {
		h.bins[index].count += 1
		return
	}

	// if not, insert it where it belongs in the order
	bin := Bin{min: value, max: value, count: 1}
	h.bins = h.bins[0 : h.Size+1] // grow the bins slice by one
	if index == len(h.bins) {     // we go at the very end since we're too small
		h.bins[h.Size] = bin
	} else { // we go before the index because we're larger than its max
		copy(h.bins[index+1:], h.bins[index:])
		h.bins[index] = bin
	}

	h.merge(h.closest())
}

// Add32 adds a float32 to the histogram, converting it to a float64
func (h *Histogram) Add32(value float32) {
	h.Add(float64(value))
}

// Percentile returns a float64 of the percentile of a given value in the histogram
func (h *Histogram) Percentile(value float64) (percentile float64) {
	var position uint64
	for i := h.Size - 1; i >= 0; i-- { // iterate in reverse to get a percentile
		if value > h.bins[i].max {
			position += h.bins[i].count
		} else { // linear estimate of value's position in its bin
			pct := 1.0
			if h.bins[i].max-h.bins[i].min > 0.0 {
				pct = (value - h.bins[i].min) / (h.bins[i].max - h.bins[i].min)
			}
			position += uint64(float64(h.bins[i].count) * pct)
			break
		}
	}
	return float64(position) / float64(h.Count)
}

// Percentile32 returns a float32 of the percentile of a given value in the histogram
func (h *Histogram) Percentile32(value float32) (percentile float32) {
	return float32(h.Percentile(float64(value)))
}

// sort Interface
// Sort the histogram in descending order for compatibility with sort.Search mechanics
func (h Histogram) Len() int           { return len(h.bins) }
func (h Histogram) Swap(i, j int)      { h.bins[i], h.bins[j] = h.bins[j], h.bins[i] }
func (h Histogram) Less(i, j int) bool { return h.bins[i].min > h.bins[j].max }

// merge merges bin j into bin i
func (h *Histogram) merge(i int, j int) {
	// make sure i < j
	if i > j {
		i, j = j, i
	}
	// merge j into i
	h.bins[i].min = h.bins[j].min
	h.bins[i].count += h.bins[j].count

	// slide everyone above j back one
	copy(h.bins[j:h.Size], h.bins[j+1:h.Size+1])
	h.bins = h.bins[0:h.Size]
}

// closest returns the indexes i, j of the two adjacent bins that span the smallest total distance
func (h *Histogram) closest() (i int, j int) {
	var gap float64
	i = 0
	minGap := h.bins[0].max - h.bins[len(h.bins)-1].min
	for pos, bin := range h.bins[0 : len(h.bins)-1] {
		gap = bin.max - h.bins[pos+1].min
		if gap < minGap {
			minGap = gap
			i = pos
		}
	}
	return i, i + 1
}

// New returns a new Histogram with binCount bins
func New(binCount int) *Histogram {
	return &Histogram{
		Size:  binCount,
		bins:  make([]Bin, binCount+1),
		Count: 0,
	}
}
