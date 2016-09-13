package ghist

import "fmt"

func (h Histogram) String() (str string) {
	str = fmt.Sprintf("%d bin ghist summarizing %d items\n", len(h.Bins), h.Count)
	str += fmt.Sprintf("Mean: %2f    Median: %2f\n", h.Mean(), h.Median())
	str += fmt.Sprintf("Mode: %v\n\nBins:\n", h.Mode().String())
	if h.MaxBinRatio > 0 {
		str += fmt.Sprintf("MaxBinRatio: %d\n", h.MaxBinRatio)
	}
	for i, bin := range h.Bins {
		str += fmt.Sprintf("%d: %v\n", i, bin.String())
	}
	return
}

func (b Bin) String() string {
	return fmt.Sprintf("%d in [%2f:%2f] totaling %2f", b.Count, b.Max, b.Min, b.Sum)
}
