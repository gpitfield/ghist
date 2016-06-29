package ghist

import (
	"testing"
)

var (
	testValues      = []float64{0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0}
	testValuesNeg   = []float64{0.0, -1.0, -2.0, -3.0, -4.0, -5.0, -6.0, -7.0, -8.0, -9.0}
	testValuesMixed = []float64{0.0, -1.0, 2.0, -3.0, 4.0, -5.0, 6.0, -7.0, 8.0, -9.0}
)

func TestHistogram(t *testing.T) {
	hist := New(5)
	var count uint64 = 0
	// create a test histogram with n+1 instances of testValues[n]
	for j, val := range testValues {
		for i := 0; i <= j; i++ {
			hist.Add(val)
			count++
		}
	}
	if hist.Count != count {
		t.Fatalf("Wrong number of entries in histogram; got %d expected %d\n", hist.Count, count)
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
	// create a test histogram with n+1 instances of testValuesNeg[n]
	for j, val := range testValuesNeg {
		for i := 0; i <= j; i++ {
			hist.Add(val)
			count++
		}
	}
	if hist.Count != count {
		t.Fatalf("Wrong number of entries in histogram; got %d expected %d\n", hist.Count, count)
	}
	if pct := hist.Percentile(0.0); pct != 1.0 {
		t.Fatalf("Wrong percentile for value 0.0; got %f expected %f\n", pct, 1.0)
	}
	if pct := hist.Percentile(-9.0); pct != 0.0 {
		t.Fatalf("Wrong percentile for value -9.0; got %f expected %f\n", pct, 0.0)
	}

	hist = New(5)
	count = 0
	// create a test histogram with n+1 instances of testValuesMixed[n]
	for j, val := range testValuesMixed {
		for i := 0; i <= j; i++ {
			hist.Add(val)
			count++
		}
	}
	if hist.Count != count {
		t.Fatalf("Wrong number of entries in histogram; got %d expected %d\n", hist.Count, count)
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
