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
	Bins  []Bin
	Size  int
}

// Bin keeps track of a histogram bin's minimum and maximum values and count
type Bin struct {
	Max   float64
	Min   float64
	Count uint64
}

var zero = Bin{} // for empty comparisons

// Add adds a float64 value to the histogram, modifying it as necessary
func (h *Histogram) Add(value float64) {
	if h.Count == math.MaxUint64 {
		panic("Integer overflow: Attempt to exceed maximum Count in ghist Histogram")
	}
	h.Count += 1

	// see if it fits in an existing bin
	index := sort.Search(len(h.Bins), func(i int) bool { return value >= h.Bins[i].Min })
	if index < len(h.Bins) && h.Bins[index].Max >= value {
		h.Bins[index].Count += 1
		return
	}

	// if not, insert it where it belongs in the order
	bin := Bin{Min: value, Max: value, Count: 1}
	h.Bins = h.Bins[0 : h.Size+1] // grow the bins slice by one
	if index == len(h.Bins) {     // we go at the very end since we're too small
		h.Bins[h.Size] = bin
	} else { // we go before the index because we're larger than its max
		copy(h.Bins[index+1:], h.Bins[index:])
		h.Bins[index] = bin
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
		if value > h.Bins[i].Max {
			position += h.Bins[i].Count
		} else { // linear estimate of value's position in its bin
			pct := 1.0
			if h.Bins[i].Max-h.Bins[i].Min > 0.0 {
				pct = (value - h.Bins[i].Min) / (h.Bins[i].Max - h.Bins[i].Min)
			}
			position += uint64(float64(h.Bins[i].Count) * pct)
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
func (h Histogram) Len() int           { return len(h.Bins) }
func (h Histogram) Swap(i, j int)      { h.Bins[i], h.Bins[j] = h.Bins[j], h.Bins[i] }
func (h Histogram) Less(i, j int) bool { return h.Bins[i].Min > h.Bins[j].Max }

// merge merges bin j into bin i
func (h *Histogram) merge(i int, j int) {
	// make sure i < j
	if i > j {
		i, j = j, i
	}
	// merge j into i
	h.Bins[i].Min = h.Bins[j].Min
	h.Bins[i].Count += h.Bins[j].Count

	// slide everyone above j back one
	copy(h.Bins[j:h.Size], h.Bins[j+1:h.Size+1])
	h.Bins = h.Bins[0:h.Size]
}

// closest returns the indexes i, j of the two adjacent bins that span the smallest total distance
func (h *Histogram) closest() (i int, j int) {
	var gap float64
	i = 0
	minGap := h.Bins[0].Max - h.Bins[len(h.Bins)-1].Min
	for pos, bin := range h.Bins[0 : len(h.Bins)-1] {
		gap = bin.Max - h.Bins[pos+1].Min
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
		Bins:  make([]Bin, binCount+1),
		Count: 0,
	}
}
