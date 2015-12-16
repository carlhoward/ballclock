package main

import "testing"

type testdat struct {
    n int
    t int
    expected string
}

var tests = []testdat{
  { 30, 0, "30 balls cycle after 15 days." },
  { 45, 0, "45 balls cycle after 378 days." },
  { 30, 325, "{\"Min\":[],\"FiveMin\":[22,13,25,3,7],\"Hour\":[6,12,17,4,15],\"Main\":[11,5,26,18,2,30,19,8,24,10,29,20,16,21,28,1,23,14,27,9]}"},
}

func TestBallClock(t *testing.T) {
    for _, dat := range tests {
        bc := NewBallClock(dat.n, dat.t)
        val := bc.Run()
        if val != dat.expected {
            t.Errorf("test failed: \n  expected[%s]\n     found[%s]\n", dat.expected, val);
        }
    }
}

func BenchmarkBallClock(b *testing.B) {
    for i := 0; i < b.N; i++ {
        bc := NewBallClock(30, 0)
        bc.Run()
        bc.Init(45, 0)
        bc.Run()
        bc.Init(30, 325)
        bc.Run()
    }
}
