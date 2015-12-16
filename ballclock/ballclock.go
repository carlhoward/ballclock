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
    "github.com/carlhoward/ballclock/track"
)

type BallClock struct {
    // clock operational parameters 
    NumberBalls int
    TargetMinutes int
    // tracks
    Min *track.Track
    FiveMin *track.Track
    Hour *track.Track
    Main *track.Track
    Tracks []*track.Track
}

type State struct {
    Min []int
    FiveMin []int
    Hour []int
    Main []int
}

func (bc *BallClock) advance() {
    bc.Min.AddBall(bc.Main.RemFirstBall())
}

func (bc *BallClock) GetTime() string {
    mins := bc.Min.Count() + 5 * bc.FiveMin.Count()
    hours := bc.Hour.Count() + 1
    s := "time: " + strconv.Itoa(hours) + ":"
    if (mins < 10) {
        s += "0"
    }
    return s + strconv.Itoa(mins)
}

func (bc *BallClock) cacheState() {
    for _, track := range bc.Tracks {
        track.CacheState()
    }
}

// test whether all tracks' current states matches their cached states
func (bc *BallClock) sameState() bool {
    for _, track := range bc.Tracks {
        if !track.SameState() {
            return false
        }
    }
    return true
}

func (bc *BallClock) getState() State {
    return State{Min: bc.Min.Balls, FiveMin: bc.FiveMin.Balls, Hour: bc.Hour.Balls, Main: bc.Main.Balls}
}

// clears all tracks and adds balls 1 through n to the main track
func (bc *BallClock) Init(numberBalls int, targetMinutes int) {
    bc.NumberBalls = numberBalls
    bc.TargetMinutes = targetMinutes
    bc.clear()
    // add balls to the overflow track
    // this gives starting time of 1:00
    for b := 1; b <= numberBalls; b++ {
        bc.Main.AddBall(b)
    }
    bc.cacheState()
}

func (bc *BallClock) clear() {
    for _, track := range bc.Tracks {
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

func NewBallClock(numberBalls int, targetMinutes int) *BallClock {
    bc := new(BallClock)
    
    // create all tracks
    bc.Min = &track.Track{Name: "Min", Limit: 4}
    bc.FiveMin = &track.Track{Name: "FiveMin", Limit: 11}
    bc.Hour = &track.Track{Name: "Hour", Limit: 11}
    bc.Main = &track.Track{Name: "Main", Limit: 0}
    bc.Tracks = []*track.Track{bc.Min, bc.FiveMin, bc.Hour, bc.Main}

    // link tracks
    bc.Min.NextTrack = bc.FiveMin
    bc.FiveMin.NextTrack = bc.Hour
    bc.Hour.NextTrack = bc.Main

    // assign the overflow track
    for _, track := range bc.Tracks {
        track.MainTrack = bc.Main
    }

    bc.Init(numberBalls, targetMinutes)
    return bc
}

func (bc *BallClock) Run() string {
    elapsedMinutes := 0
    for {
        bc.advance()
        elapsedMinutes++
        if bc.TargetMinutes != 0 {
            if elapsedMinutes == bc.TargetMinutes {
                strJ, _ := json.Marshal(bc.getState())
                return string(strJ)
            }
        } else if(bc.sameState()) {
            days := elapsedMinutes / 60 / 24
            return fmt.Sprintf("%d balls cycle after %d days.", bc.NumberBalls, days)
        }
    }
}

func main() {
    var numberBalls int
    var targetMinutes int

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

    bc := NewBallClock(numberBalls, targetMinutes)
    fmt.Sprintf("%s\n", bc.Run())
}

