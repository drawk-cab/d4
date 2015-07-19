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

var DEBUG = false

func NewMachineString(in string, sample_rate float64, save_s float64,
                      clip float64, imports map[string]string, workers int) (Machine, error) {
    return NewMachine(strings.NewReader(in), sample_rate, save_s, clip, imports, workers)
}

func NewMachine(in io.Reader, sample_rate float64, save_s float64,
                clip float64, imports map[string]string, workers int) (Machine, error) {
    s := NewOpcodeMachine(sample_rate, save_s, clip, imports, workers)
    s.Init(nil)
    err := s.Program(in)
    return s, err
}

func CloneMachine(in io.Reader, m Machine) (Machine, error) {
    s := NewOpcodeMachine(0, 0, 0, nil, 0)
    s.Init(m)

    err := s.Program(in)
    return s, err
}
