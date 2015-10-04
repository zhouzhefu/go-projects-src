package gotest

import (
	"testing"
)

func Test_Division_1(t *testing.T) {
	if result, err := Division(6, 2); result != 3 || err != nil {
		t.Error("Division func test failed with unexpected result")
	} else {
		t.Log("First test [Test_Division_1] passed")
	}
}

func TestDivision_2(t *testing.T) {
	// t.Error("Second test [TestDivision_2] just cannot pass")
	if result, err := Division(6, 0); result == 0 || err != nil {
		t.Log("Second test [TestDivision_2] passed")
	}
}

func Benchmark_Division(b *testing.B) {
	b.StopTimer()

	// do something prepration

	b.StartTimer()
	for i:=0; i<b.N; i++ {
		Division(4, 6)
	}
}