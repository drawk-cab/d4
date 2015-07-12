package d4

import (
    "testing"
    "time"
    "fmt"
    "os"
    "bufio"
)

const ITERS = 10000

func chk(err error) {
    if (err != nil) {
        panic(err)
    }
}

func test(t *testing.T, name string, code string, expect_error bool, expect []float64, iterations int, debug bool) {
    DEBUG = debug
    machine, err := NewMachineString(code, 22050, 1)
    if err == nil {
        test_machine(t, name, machine, expect_error, expect, iterations)
    } else {
        if !expect_error {
            t.Errorf("unexpected compile error", err)
        }
    }
}

func test_file(t *testing.T, name string, filename string, expect_error bool, expect []float64, iterations int, debug bool) {
    DEBUG = debug
    opened_file, err := os.OpenFile(filename, os.O_RDONLY, 0755)
    if err != nil {
        panic(err)
    }
    in := bufio.NewReader( opened_file )
    machine, err := NewMachine(in, 22050, 1)
    if err == nil {
        test_machine(t, name, machine, expect_error, expect, iterations)
    } else {
        if !expect_error {
            t.Errorf("unexpected compile error", err)
        }
    }
}

func test_machine(t *testing.T, name string, machine Machine, expect_error bool, expect []float64, iterations int) {
    result, err := machine.Run()

    if err != nil && !expect_error {
        t.Errorf("%s : unexpected runtime error", name, err)
    } else {
        if err == nil && expect_error {
            t.Errorf("%s : expected an error but didn't get one", name)
        }
    }

    if len(result) != len(expect) {
        t.Errorf("%s : result %f, want %f", name, result, expect)
        return
    }

    for i, result_i := range result {    
        if result_i != expect[i] {
            t.Errorf("%s : result %f, want %f", name, result, expect)
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
        fmt.Printf("%s : %d iterations took %s (%d kHz)\n", name, iterations, elapsed,
            (int64(iterations) * 1000000 / elapsed.Nanoseconds()))
    }

    return
}

func TestEmpty(t *testing.T) {
    test( t,  "empty",
              "",
              false, []float64{},
              ITERS, false,
    )
}

func TestPush(t *testing.T) {
    test( t,  "push",
              "47.3 .",
              false, []float64{47.3},
              ITERS, false,
    )
}

func TestAdd(t *testing.T) {
    test( t,  "add",
              "47 21 + .",
              false, []float64{68},
              ITERS, false,
    )
}

func TestSub(t *testing.T) {
    test( t,  "sub",
              "47 21 - .",
              false, []float64{26},
              ITERS, false,
    )
}

func TestMul(t *testing.T) {
    test( t,  "mul",
              "47 2 * .",
              false, []float64{94},
              ITERS, false,
    )
}

func TestDiv(t *testing.T) {
    test( t,  "div",
              "94 2 / .",
              false, []float64{47},
              ITERS, false,
    )
}

func TestSwap(t *testing.T) {
    test( t,  "swap",
              "47 21 SWAP . .",
              false, []float64{47, 21},
              ITERS, false,
    )
}

func TestComment(t *testing.T) {
    test( t,  "comment",
              "47 (Hi; I am: a comment (quite a hard one)) 21 SWAP . .",
              false, []float64{47, 21},
              ITERS, false,
    )
}

func TestDefinition(t *testing.T) {
    test( t,  "definition",
              "3 :five 5; :two five 3 -; five + two + .",
              false, []float64{10},
              ITERS, false,
    )
}

func TestRecursiveDefinition(t *testing.T) {
    test( t,  "recursive definition (error)",
              ":here there; :there yonder; :yonder here; here .",
              true, nil,
              1, false,
    )
}

func TestIfElse(t *testing.T) {
    test( t,  "if-else",
              "4 11 10 > IF 1 . ELSE 2 . THEN 11 10 < IF 1 . ELSE 2 . THEN .",
              false, []float64{1, 2, 4},
              ITERS, false,
    )
}

func TestNestedIfElse(t *testing.T) {
    test( t,  "nested if-else",
              "4 11 10 > IF 11 10 > IF 1 . ELSE 2 . THEN ELSE 3 . THEN 11 10 < IF 1 . ELSE 11 10 < IF 2 . ELSE 3 . THEN THEN .",
              false, []float64{1, 3, 4},
              ITERS, false,
    )
}

func TestChoose(t *testing.T) {
    test( t,  "choose",
              "3 FROM 7, 8, 9, 10, 11, 12 CHOOSE .",
              false, []float64{10},
              ITERS, false,
    )
}

func TestOscillators(t *testing.T) {
    test( t,  "oscillators",
              "0.25 SIN . 0.25 SAW . 0.25 SQ . 0.25 TR .",
              false, []float64{1, -0.5, 1, 0},
              ITERS, false,
    )
}

func TestOptimize(t *testing.T) {
    test( t,  "optimize",
              "[400 2 * 800 /] [96 30 - [12 21 +] /] +.",
              false, []float64{3},
              ITERS, false,
    )
}

func TestOptimizeBaseline(t *testing.T) {
    test( t,  "optimize baseline",
              "400 2 * 800 / 96 30 - 12 21 + / +.",
              false, []float64{3},
              ITERS, false,
    )
}

func TestOptimizeUnderrun(t *testing.T) {
    test( t,  "optimize underrun (error)",
              "4 [2 *] 3 +.",
              true, nil,
              1, false,
    )
}

func TestOptimizeUnfinished(t *testing.T) {
    test( t,  "optimize unfinished",
              "[4 2 *.",
              false, []float64{8},
              ITERS, false,
    )
}

func TestOptimizeOutput(t *testing.T) {
    test( t,  "optimize with output (error)",
              "[4 2 * .] 3 +.",
              true, nil,
              1, false,
    )
}

func TestImport(t *testing.T) {
    test( t,  "import",
              "::scale (comment); ut.",
              false, []float64{220},
              ITERS, false,
    )
}

func TestFile(t *testing.T) {
    test_file( t,  "gloucester",
              "tests/gloucester.d4",
              false, []float64{0},
              1, true,
    )
}
