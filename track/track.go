//
// package implementing functionality of a track in a ball clock
//
package track

import (
    "fmt"
    "os"
)

type Track struct {
    Name string
    Balls []int
    Limit int
    NextTrack *Track
    MainTrack *Track
    CachedState []int
}

func (t *Track) Count() int {
    return len(t.Balls)
}

func (t *Track) Val() string {
    return fmt.Sprintf("%s %v\n", t.Name, t.Balls);
}

func (t *Track) CacheState() {
    t.CachedState = make([]int, len(t.Balls), len(t.Balls))
    copy(t.CachedState, t.Balls);
}

func (t *Track) SameState() bool {
    if len(t.Balls) != len(t.CachedState) {
        return false;
    }
    for i := 0; i < len(t.Balls); i++ {
        if t.Balls[i] != t.CachedState[i] {
            return false;
        }
    }
    return true;
}

func (t *Track) RemFirstBall() int {
    l := len(t.Balls)
    if (l == 0) {
        fmt.Println("Error: empty overflow track");
        os.Exit(1);
    }
    b := t.Balls[0]
    t.Balls = t.Balls[1:]
    return b
}

func (t *Track) RemLastBall() int {
    l := len(t.Balls)
    b := t.Balls[l - 1]
    t.Balls = t.Balls[:l-1]
    return b
}

func (t *Track) Clear() {
    t.Balls = t.Balls[:0]
}

func (t *Track) AddBall(b int) {
//    fmt.Printf("Adding %d to %v\n", b, t);
    if (t.Limit != 0 && len(t.Balls) == t.Limit) {
        // cascade
        for (t.Count() > 0) {
            t.MainTrack.AddBall(t.RemLastBall())
        }
        t.NextTrack.AddBall(b);
    } else {
        t.Balls = append(t.Balls, b);
    }
}

