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

package main

import (
    "fmt"
    "os"
    "strconv"
    "encoding/json"
	"github.com/carlhoward/track"
)

// global tracks
var minTrack track.Track
var fiveminTrack track.Track
var hourTrack track.Track
var mainTrack track.Track
var tracks []*track.Track

// clock operational parameters
var numberBalls int
var targetMinutes int
var elapsedMinutes int

func advance() {
    minTrack.AddBall(mainTrack.RemFirstBall())
}

func getTime() string {
    mins := minTrack.Count() + 5 * fiveminTrack.Count()
    hours := hourTrack.Count() + 1
    s := "time: " + strconv.Itoa(hours) + ":"
    if (mins < 10) {
        s += "0"
    }
    return s + strconv.Itoa(mins)
}

func cacheState() {
    for _, track := range tracks {
        track.CacheState()
    }
}

// test whether all tracks' current states matches their cached states
func sameState() bool {
    for _, track := range tracks {
        if !track.SameState() {
            return false
        }
    }
    return true
}

func getState() map[string][]int {
    m := make(map[string][]int)
    for _, track := range tracks {
        m[track.Name] = track.Balls
    }
    return m
}

// clears all tracks and adds balls 1 through n to the main track
func initTracks(n int) {
    clear()
    // add balls to the overflow track
    // this gives starting time of 1:00
    for b := 1; b <= n; b++ {
        mainTrack.AddBall(b)
    }
}

func clear() {
    elapsedMinutes = 0
    for _, track := range tracks {
        track.Clear()
    }
}

func exitUsage() {
        fmt.Printf("Usage: %s ballCount [timeToRun]\n", os.Args[0])
        os.Exit(0)
}

// utility function to parse the given command line argument as an integer.
func getIntegerArg(argIndex int) int {
    intVal, err := strconv.Atoi(os.Args[argIndex])
    if err != nil {
        exitUsage()
    }
    return intVal
}

func main() {
	// parse command line arguments to determine operational
	// parameters and mode of operation.
    if len(os.Args) == 2 {
        // MODE_1: run until state repeats, then display days elapsed
        numberBalls = getIntegerArg(1)
        targetMinutes = 0
    } else if len(os.Args) == 3 {
        // MODE_2: run for given time in minutes, then display state
        numberBalls = getIntegerArg(1)
        targetMinutes = getIntegerArg(2)
    } else {
        exitUsage()
    }
    if numberBalls < 27 || numberBalls > 127 {
        fmt.Printf("ballCount must  be at least 27 and at most 127\n")
        exitUsage()
    }

    // create all tracks
    minTrack = track.Track{Name: "Min", Limit: 4}
    fiveminTrack = track.Track{Name: "FiveMin", Limit: 11}
    hourTrack = track.Track{Name: "Hour", Limit: 11}
    mainTrack = track.Track{Name: "Main", Limit: 0}
    tracks = []*track.Track{&minTrack, &fiveminTrack, &hourTrack, &mainTrack}

    // link tracks
    minTrack.NextTrack = &fiveminTrack
    fiveminTrack.NextTrack = &hourTrack
    hourTrack.NextTrack = &mainTrack

    // assign the overflow track
    for _, track := range tracks {
        track.MainTrack = &mainTrack
    }

    initTracks(numberBalls)
    cacheState()

    // main loop
    for {
        advance()
        elapsedMinutes++
        if targetMinutes != 0 {
            if elapsedMinutes == targetMinutes {
                strJ, _ := json.Marshal(getState())
                fmt.Printf("State for %d balls after %d minutes:\n%s\n", numberBalls, elapsedMinutes, strJ)
                break
            }
        } else if(sameState()) {
            days := elapsedMinutes / 60 / 24
            fmt.Printf("%d balls cycle after %d days.\n", numberBalls, days)
            break
        }
    }

}
