package d4

import (
    "testing"
    "time"
    "fmt"
)

func test(t *testing.T, code string, expect []float64, iterations int) *Machine {
    machine := CompileString(code, 22050)
    result := machine.Run(false)

    for i, result_i := range result {    
        if result_i != expect[i] {
            t.Errorf("%s failed: result %f, want %f", code, result, expect)
            break
        }   
    }

    if iterations > 0 {
        then := time.Now()
        for i := 1; i <= iterations; i++ {
            _ = machine.Run(false)
        }
        fmt.Printf("%d iterations took %s\n", iterations, time.Since(then))
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
