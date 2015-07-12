package d4

import "io"

type MachineData struct {
    iter int
    sampleRate float64
    clip float64
}

type Machine interface {
    Init(Machine) error
    Program(io.Reader) (error)
    Run() ([]float64, error)
    RunCode([]float64) ([]float64, []float64, error)
    Fill32([]float32) error
    GetData() MachineData
}
