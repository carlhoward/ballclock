package main

import (
    "fmt"
    "os"
    "strconv"
    "encoding/json"
)

type track struct {
    name string
    balls []int
    limit int
    nextTrack *track
    mainTrack *track
    cachedState []int
}

func (t *track) count() int {
    return len(t.balls)
}

func (t *track) val() string {
    return fmt.Sprintf("%s %v\n", t.name, t.balls)
}

func (t *track) cacheState() {
    t.cachedState = make([]int, len(t.balls), len(t.balls))
    copy(t.cachedState, t.balls)
}

func (t *track) sameState() bool {
    if len(t.balls) != len(t.cachedState) {
        return false
    }
    for i := 0; i < len(t.balls); i++ {
        if t.balls[i] != t.cachedState[i] {
            return false
        }
    }
    return true
}

func (t *track) remFirstBall() int {
    l := len(t.balls)
    if (l == 0) {
        fmt.Println("Error: empty overflow track")
        os.Exit(1)
    }
    b := t.balls[0]
    t.balls = t.balls[1:]
    return b
}

func (t *track) remLastBall() int {
    l := len(t.balls)
    b := t.balls[l - 1]
    t.balls = t.balls[:l-1]
    return b
}

func (t *track) clear() {
    t.balls = t.balls[:0]
}

func (t *track) addBall(b int) {
//    fmt.Printf("Adding %d to %v\n", b, t)
    if (t.limit != 0 && len(t.balls) == t.limit) {
        // cascade
        for (t.count() > 0) {
            t.mainTrack.addBall(t.remLastBall())
        }
        t.nextTrack.addBall(b)
    } else {
        t.balls = append(t.balls, b)
    }
}

func advance() {
    minTrack.addBall(mainTrack.remFirstBall())
}

func getTime() string {
    mins := minTrack.count() + 5 * fiveminTrack.count()
    hours := hourTrack.count() + 1
    s := "time: " + strconv.Itoa(hours) + ":"
    if (mins < 10) {
        s += "0"
    }
    return s + strconv.Itoa(mins)
}

func cacheState() {
    for _, track := range tracks {
        track.cacheState()
    }
}

func sameState() bool {
    for _, track := range tracks {
        if !track.sameState() {
            return false
        }
    }
    return true
}

func getState() map[string][]int {
    m := make(map[string][]int)
    for _, track := range tracks {
        m[track.name] = track.balls
    }
    return m
}

func initTracks(n int) {
    clear()
    // add balls to the overflow track
    // this gives starting time of 1:00
    for b := 1; b <= n; b++ {
        mainTrack.addBall(b)
    }
}

func clear() {
    elapsedMinutes = 0
    for _, track := range tracks {
        track.clear()
    }
}

func exitUsage() {
        fmt.Printf("Usage: %s ballCount [timeToRun]\n", os.Args[0])
        os.Exit(0)
}

var minTrack track
var fiveminTrack track
var hourTrack track
var mainTrack track
var tracks []*track

var numberBalls int
var targetMinutes int
var elapsedMinutes int

func getIntegerArg(argIndex int) int {
    intVal, err := strconv.Atoi(os.Args[argIndex])
    if err != nil {
        exitUsage()
    }
    return intVal
}

func main() {

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
    minTrack = track{name: "Min", limit: 4}
    fiveminTrack = track{name: "FiveMin", limit: 11}
    hourTrack = track{name: "Hour", limit: 11}
    mainTrack = track{name: "Main", limit: 0}
    tracks = []*track{&minTrack, &fiveminTrack, &hourTrack, &mainTrack}

    // link tracks
    minTrack.nextTrack = &fiveminTrack
    fiveminTrack.nextTrack = &hourTrack
    hourTrack.nextTrack = &mainTrack

    // assign the overflow track
    for _, track := range tracks {
        track.mainTrack = &mainTrack
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
