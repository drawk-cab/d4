package d4

import "io"

type Machine interface {
    Init() error
    Program(io.Reader) error
    Run() ([]float64, error)
    Fill32([]float32) error
}
