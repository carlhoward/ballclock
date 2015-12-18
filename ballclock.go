
//
// Program to simulate operation of a mechanical ball clock
//
// This program is an exercise in learning the go programming language.
// 
// The clock simulation executes in 2 modes of operation determined by
// the number of arguments given:
//
// MODE_1 usage: ballclock numberBalls
//     -> run until state repeats, then display days elapsed
//
// MODE_2 usage: ballclock numberBalls targetMinutes
//     -> run for given time in minutes, then display state
//
// representational tracks are implemented as lifo queues while the overflow
// track is implemented as a ring buffer (fifo queue).
//

package main

import (
    "fmt"
    "os"
    "strconv"
    "encoding/json"
)

var numberBalls int

var overflow ring // ring buffer as overflow track
var min track // minute track
var fmin track // five-minute track
var hour track // hour track 	

type buffer interface {
    push(int)
    pop() int
    val() []int
}

//  ring buffer definitions (fifo queue)

type ring struct {
    balls [128]int
    start int
    end int
}

func (r *ring) push(b int) {
    r.balls[r.end] = b
    r.end = (r.end + 1) & 0x7F
}

func (r *ring) pop() int {
    b := r.balls[r.start]
    r.start = (r.start + 1) & 0x7F
    return b
}

func (r *ring) val() []int {
    if (r.start < r.end) {
        return overflow.balls[r.start : r.end]
    } else {
        return append(r.balls[r.start : 128], r.balls[0 : r.end]...)
    }
}

// return true if ring buffer is in original state
func (r *ring) check() bool {
    val := 1;
    if r.start < r.end {
        for i := r.start; i < r.end; i++ {
            if r.balls[i] != val {
                return false
            }
            val++
        }
    } else {
        for i := r.start; i < 128; i++ {
            if r.balls[i] != val {
                return false
            }
            val++
        }
        for i := 0; i < r.end; i++ {
            if r.balls[i] != val {
                return false
            }
            val++
        }
    }
    return true
}

//  track definitions (lifo queue)

type track struct {
    balls [12]int
    count int
    limit int
    next buffer
}

func (t *track) push(b int) {
    if (t.count == t.limit) {
        // empty track into overflow and cascade to next track
        for i := 0; i < t.limit; i++ {
            overflow.push(t.pop())
        }
        t.next.push(b)
    } else {
        t.balls[t.count] = b
        t.count++
    }
}

func (t *track) pop() int {
    t.count--
    return t.balls[t.count]
}

func (t *track) val() []int {
    return t.balls[0:t.count]
}

// main functions

type state struct {
    Min []int
    FiveMin []int
    Hour []int
    Main []int
}

func getState() state {
    return state{Min: min.val(), FiveMin: fmin.val(), Hour: hour.val(), Main: overflow.val()}
}

func exitUsage() {
        fmt.Printf("Usage: %s ballCount [timeToRun]\n", os.Args[0])
        os.Exit(0)
}

func getIntegerArg(argIndex int) int {
    intVal, err := strconv.Atoi(os.Args[argIndex])
    if err != nil {
        exitUsage()
    }
    return intVal
}

func main() {
    targetMins := 0
    halfDays := 0
    if len(os.Args) == 2 {
        // MODE_1: run until state repeats, then display days elapsed
        numberBalls = getIntegerArg(1)
    } else if len(os.Args) == 3 {
        // MODE_2: run for given time in minutes, then display state
        numberBalls = getIntegerArg(1)
        targetMins = getIntegerArg(2)
    } else {
        exitUsage()
    }
    if numberBalls < 27 || numberBalls > 127 {
        fmt.Printf("ballCount must  be at least 27 and at most 127\n")
        exitUsage()
    }

    // create tracks
    min = track{count: 0, limit: 4}
    fmin = track{count: 0, limit: 11}
    hour = track{count: 0, limit: 11}
    overflow = ring{start: 0, end: 0}

    // link tracks
    min.next = &fmin
    fmin.next = &hour
    hour.next = &overflow

    // initialize balls [1 .. n]
    for b := 1; b <= numberBalls; b++ {
        overflow.push(b)
    }
    
    // pick run loop based on mode of operation
    if targetMins == 0 {
        for {
            // loop for 12 hours (= 720 minutes)
            for i := 0;  i < 720; i++ {
                b := overflow.pop()
                min.push(b)
            }
            halfDays++
            if overflow.check() {
                days := halfDays / 2
                fmt.Printf("%d balls cycle after %d days.\n", numberBalls, days)
                break
            }
        }
    } else {
        // loop for allotted time
        for i := 0;  i < targetMins; i++ {
            min.push(overflow.pop())
        }
        strJ, _ := json.Marshal(getState())
        fmt.Printf("%s\n", string(strJ))
    }

}
