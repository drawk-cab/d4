package d4

import (
    "testing"
    "time"
    "fmt"
    "os"
    "bufio"
)

func test(t *testing.T, code string, expect []float64, iterations int) *Machine {
    machine := CompileString(code, 22050)
    return test_machine(t, machine, expect, iterations)
}

func test_file(t *testing.T, filename string, expect []float64, iterations int) *Machine {
    opened_file, err := os.OpenFile(filename, os.O_RDONLY, 0755)
    if err != nil {
        panic(err)
    }
    in := bufio.NewReader( opened_file )
    machine := Compile(in, 22050)
    return test_machine(t, machine, expect, iterations)
}

func test_machine(t *testing.T, machine *Machine, expect []float64, iterations int) *Machine {
    result := machine.Run(false)

    for i, result_i := range result {    
        if result_i != expect[i] {
            t.Errorf("result %f, want %f", result, expect)
            break
        }   
    }

    if iterations > 0 {
        then := time.Now()
        for i := 1; i <= iterations; i++ {
            _ = machine.Run(false)
        }
        elapsed := time.Since(then)
        fmt.Printf("%d iterations took %s (%d kHz)\n", iterations, elapsed,
            (int64(iterations) * 1000000 / elapsed.Nanoseconds()))
    }

    return machine
}

func TestEmpty(t *testing.T) {
    _ = test( t, 
              "",
              []float64{},
              1000,
    )
}

func TestPush(t *testing.T) {
    _ = test( t, 
              "47.3",
              []float64{47.3},
              1000,
    )
}

func TestAdd(t *testing.T) {
    _ = test( t, 
              "47 21 +",
              []float64{68},
              1000,
    )
}

func TestSub(t *testing.T) {
    _ = test( t, 
              "47 21 -",
              []float64{26},
              1000,
    )
}

func TestSwap(t *testing.T) {
    _ = test( t, 
              "47 21 SWAP",
              []float64{21, 47},
              1000,
    )
}

func TestLoopTune(t *testing.T) {
    _ = test_file( t, 
              "tests/loop-tune.d4",
              []float64{0},
              1000,
    )
}
