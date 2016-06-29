# ghist
Ghist implements a streaming histogram in Go, loosely along the lines of http://jmlr.org/papers/volume11/ben-haim10a/ben-haim10a.pdf.

It's spelled *ghist* because it's a Go Histogram, and it's pronounced *gist*, because that's what a streaming histogram provides.

To use ghist, instantiate a new histogram with ghist.New(binCount), where binCount is the number of bins you want in the histogram.

As values stream in, add them to the histogram with AddValue(value float64).

To find the percentile of a value, use ValuePercentile(value float64).

Streaming histograms dynamically resize bins, trying to maintain as much information resolution in the histogram as possible. 
For most uses, a binCount < 100 is fine.
The bins start out very narrow, but as more values are added the bins grow to cover almost the entirety of the range of values.
