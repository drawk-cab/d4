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

const DEBUG = false

func NewMachineString(in string, sampleRate float64, clip float64) (Machine, error) {
    return NewMachine(strings.NewReader(in), sampleRate, clip)
}

func NewMachine(in io.Reader, sampleRate float64, clip float64) (Machine, error) {
    s := NewOpcodeMachine(sampleRate, clip)
    s.Init(nil)
    err := s.Program(in)
    return s, err
}

func CloneMachine(in io.Reader, m Machine) (Machine, error) {
    s := NewOpcodeMachine(0, 0)
    s.Init(m)
    err := s.Program(in)
    return s, err
}
