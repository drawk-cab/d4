package d4

import (
    "strings"
    "io"
    "math"
)

/* loop length in seconds, will get a click after this, default to 1 day(!) */
const LOOP = 60 * 60 * 24

var SEMITONE = math.Pow(2, 1.0/12)
var SEC = float64(LOOP)
var BPM = float64(LOOP / 60)

const DEBUG = true

func CompileString(in string, sampleRate float64) (Machine, error) {
    return Compile(strings.NewReader(in), sampleRate)
}

func Compile(in io.Reader, sampleRate float64) (Machine, error) {

    // Set iter nonzero to avoid zeros everywhere during dummy run
    s := NewOpcodeMachine(sampleRate)

    s.Init()
    err := s.Program(in)

    return s, err
}
