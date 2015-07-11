package d4

import (
    "testing"
    "time"
    "fmt"
    "os"
    "bufio"
)

const ITERS = 1

func chk(err error) {
    if (err != nil) {
        panic(err)
    }
}

func test(t *testing.T, code string, expect_error bool, expect []float64, iterations int) {
    machine, err := NewMachineString(code, 22050)
    if err == nil {
        test_machine(t, machine, expect_error, expect, iterations)
    } else {
        if !expect_error {
            t.Errorf("unexpected compile error", err)
        }
    }
}

func test_file(t *testing.T, filename string, expect_error bool, expect []float64, iterations int) {
    opened_file, err := os.OpenFile(filename, os.O_RDONLY, 0755)
    if err != nil {
        panic(err)
    }
    in := bufio.NewReader( opened_file )
    machine, err := NewMachine(in, 22050)
    if err == nil {
        test_machine(t, machine, expect_error, expect, iterations)
    } else {
        if !expect_error {
            t.Errorf("unexpected compile error", err)
        }
    }
}

func test_machine(t *testing.T, machine Machine, expect_error bool, expect []float64, iterations int) {
    result, err := machine.Run()

    if err != nil && !expect_error {
        t.Errorf("unexpected runtime error", err)
    } else {
        if err == nil && expect_error {
            t.Errorf("expected an error but didn't get one")
        }
    }

    if len(result) != len(expect) {
        t.Errorf("result %f, want %f", result, expect)
        return
    }

    for i, result_i := range result {    
        if result_i != expect[i] {
            t.Errorf("result %f, want %f", result, expect)
            return
        }   
    }

    if iterations > 0 {
        then := time.Now()
        for i := 1; i <= iterations; i++ {
            _, err := machine.Run()
            _ = err
        }
        elapsed := time.Since(then)
        fmt.Printf("%d iterations took %s (%d kHz)\n", iterations, elapsed,
            (int64(iterations) * 1000000 / elapsed.Nanoseconds()))
    }

    return
}

func TestEmpty(t *testing.T) {
    test( t, 
              "",
              false, []float64{},
              ITERS,
    )
}

func TestPush(t *testing.T) {
    test( t, 
              "47.3",
              false, []float64{47.3},
              ITERS,
    )
}

func TestAdd(t *testing.T) {
    test( t, 
              "47 21 +",
              false, []float64{68},
              ITERS,
    )
}

func TestSub(t *testing.T) {
    test( t, 
              "47 21 -",
              false, []float64{26},
              ITERS,
    )
}

func TestSwap(t *testing.T) {
    test( t, 
              "47 21 SWAP",
              false, []float64{21, 47},
              ITERS,
    )
}

func TestComment(t *testing.T) {
    test( t, 
              "47 (Hi; I am: a comment (quite a hard one)) 21 SWAP",
              false, []float64{21, 47},
              ITERS,
    )
}

func TestDefinition(t *testing.T) {
    test( t, 
              "3 :five 5; :two five 3 -; five + two +",
              false, []float64{10},
              ITERS,
    )
}

func TestRecursiveDefinition(t *testing.T) {
    test( t, 
              ":here there; :there yonder; :yonder here; here",
              true, nil,
              ITERS,
    )
}

func TestIfElse(t *testing.T) {
    test( t, 
              "11 10 > IF 1 ELSE 2 THEN 11 10 < IF 1 ELSE 2 THEN",
              false, []float64{1, 2},
              ITERS,
    )
}

func TestNestedIfElse(t *testing.T) {
    test( t, 
              "11 10 > IF 11 10 > IF 1 ELSE 2 THEN ELSE 3 THEN 11 10 < IF 1 ELSE 11 10 < IF 2 ELSE 3 THEN THEN",
              false, []float64{1, 3},
              ITERS,
    )
}

func TestOscillators(t *testing.T) {
    test( t, 
              "0.25 SIN 0.25 SAW 0.25 SQ 0.25 TR",
              false, []float64{1, -0.5, 1, 1},
              ITERS,
    )
}

func TestLoopTune(t *testing.T) {
    test_file( t, 
              "tests/loop-tune.d4",
              false, []float64{0},
              ITERS,
    )
}
