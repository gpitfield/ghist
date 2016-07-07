package ghist

import (
	"testing"
)

var (
	testValues      = []float64{0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0}
	testValuesNeg   = []float64{0.0, -1.0, -2.0, -3.0, -4.0, -5.0, -6.0, -7.0, -8.0, -9.0}
	testValuesMixed = []float64{0.0, -1.0, 2.0, -3.0, 4.0, -5.0, 6.0, -7.0, 8.0, -9.0}
)

// TODO: move percentiles to TestStatistics, and add checks for max/min/count in TestHistogram
func TestHistogram(t *testing.T) {
	hist := New(5)
	var (
		count uint64  = 0
		sum   float64 = 0.0
	)
	// create a test histogram with n+1 instances of testValues[n]
	for j, val := range testValues {
		for i := 0; i <= j; i++ {
			hist.Add(val)
			sum += val
			count++
		}
	}
	if hist.Count != count {
		t.Fatalf("Wrong number of entries in histogram; got %d expected %d\n", hist.Count, count)
	}
	if hist.Sum != sum {
		t.Fatalf("Wrong sum of values in histogram; got %d expected %d\n", hist.Sum, sum)
	}
	if pct := hist.Percentile(0.0); pct != 0.0 {
		t.Fatalf("Wrong percentile for value 0.0; got %f expected %f\n", pct, 0.0)
	}
	if pct := hist.Percentile(1.0); pct != 3.0/55.0 {
		t.Fatalf("Wrong percentile for value 1.0; got %f expected %f\n", pct, 3.0/55.0)
	}
	if pct := hist.Percentile(9.0); pct != 1.0 {
		t.Fatalf("Wrong percentile for value 9.0; got %f expected %f\n", pct, 1.0)
	}

	hist = New(5)
	count = 0
	sum = 0
	// create a test histogram with n+1 instances of testValuesNeg[n]
	for j, val := range testValuesNeg {
		for i := 0; i <= j; i++ {
			hist.Add(val)
			sum += val
			count++
		}
	}
	if hist.Count != count {
		t.Fatalf("Wrong number of entries in histogram; got %d expected %d\n", hist.Count, count)
	}
	if hist.Sum != sum {
		t.Fatalf("Wrong sum of values in histogram; got %d expected %d\n", hist.Sum, sum)
	}
	if pct := hist.Percentile(0.0); pct != 1.0 {
		t.Fatalf("Wrong percentile for value 0.0; got %f expected %f\n", pct, 1.0)
	}
	if pct := hist.Percentile(-9.0); pct != 0.0 {
		t.Fatalf("Wrong percentile for value -9.0; got %f expected %f\n", pct, 0.0)
	}

	hist = New(5)
	count = 0
	sum = 0
	// create a test histogram with n+1 instances of testValuesMixed[n]
	for j, val := range testValuesMixed {
		for i := 0; i <= j; i++ {
			hist.Add(val)
			sum += val
			count++
		}
	}
	if hist.Count != count {
		t.Fatalf("Wrong number of entries in histogram; got %d expected %d\n", hist.Count, count)
	}
	if hist.Sum != sum {
		t.Fatalf("Wrong sum of values in histogram; got %d expected %d\n", hist.Sum, sum)
	}
	if pct := hist.Percentile(-9.0); pct != 0.0 {
		t.Fatalf("Wrong percentile for value -9.0; got %f expected %f\n", pct, 0.0)
	}
	if pct := hist.Percentile(8.0); pct != 1.0 {
		t.Fatalf("Wrong percentile for value 8.0; got %f expected %f\n", pct, 1.0)
	}
	if pct := hist.Percentile(0.0); pct != 30.0/55.0 {
		t.Fatalf("Wrong percentile for value 0.0; got %f expected %f\n", pct, 30.0/55.0)
	}
}

func TestStatistics(t *testing.T) {
	hist := New(5)
	var count uint64 = 0
	// create a test histogram with n+1 instances of testValues[n]
	for j, val := range testValues {
		for i := 0; i <= j; i++ {
			hist.Add(val)
			count++
		}
	}
	if hist.Mean() != 6.0 {
		t.Fatalf("Wrong mean for histogram; got %f expected %f\n", hist.Mean(), 6.0)
	}
	if hist.Median() != 6.5 {
		t.Fatalf("Wrong median for histogram; got %f expected %f\n", hist.Median(), 6.5)
	}
	if hist.Mode().Count != 19 {
		t.Fatalf("Wrong mode for histogram; got %d expected %d\n", hist.Mode().Count, 19)
	}
}
