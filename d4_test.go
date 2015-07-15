package d4

import (
    "testing"
    "time"
    "fmt"
    "os"
    "bufio"
)

var TEST_IMPORTS map[string]string = map[string]string{ "IMPORT": ":imported 57;" }

func chk(err error) {
    if (err != nil) {
        panic(err)
    }
}

func test(t *testing.T, name string, code string, expect_error bool, expect []float64, debug bool) {
    DEBUG = debug

    machine, err := NewMachineString(code, 22050, 1.0, 1, TEST_IMPORTS, 1)
    if err == nil {
        test_machine(t, name, machine, expect_error, expect)
    } else {
        if !expect_error {
            t.Errorf("%s: unexpected compile error: ", name, err)
        }
    }
}

func test_file(t *testing.T, name string, filename string, expect_error bool, expect []float64, debug bool) {
    DEBUG = debug

    opened_file, err := os.OpenFile(filename, os.O_RDONLY, 0755)
    if err != nil {
        panic(err)
    }
    in := bufio.NewReader( opened_file )
    machine, err := NewMachine(in, 22050, 1.0, 1, TEST_IMPORTS, 1)

    if err == nil {
        test_machine(t, name, machine, expect_error, expect)
    } else {
        if !expect_error {
            t.Errorf("unexpected compile error %s", err)
        }
    }
}

func test_machine(t *testing.T, name string, machine Machine, expect_error bool, expect []float64) {
    result, err := machine.Run()

    if err != nil && !expect_error {
        t.Errorf("%s : unexpected runtime error: %s", name, err)
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

    return
}

func test_fill32(t *testing.T, name string, filename string, expect_error bool, expect []float32, buf_size int, workers int) {
    DEBUG = false

    opened_file, err := os.OpenFile(filename, os.O_RDONLY, 0755)
    if err != nil {
        panic(err)
    }

    in := bufio.NewReader( opened_file )
    machine, err := NewMachine(in, 22050, 1.0, 1, TEST_IMPORTS, workers)
    if err != nil {
        panic(err)
    }


    buf := make([]float32, buf_size)

    then := time.Now()
    err = machine.Fill32(buf)
    elapsed := time.Since(then)

    if err != nil && !expect_error {
        t.Errorf("%s (fill32) : unexpected runtime error", name, err)
    } else {
        if err == nil && expect_error {
            t.Errorf("%s (fill32) : expected an error but didn't get one", name)
        }
    }

    if expect != nil {
        for i, buf_i := range buf {    
            if buf_i != expect[i] {
                t.Errorf("%s : result %f, want %f", name, buf, expect)
                return
            }
        }
    }

    fmt.Printf("%s : filled %d in %s (%d kHz)\n", name, buf_size, elapsed,
            (int64(buf_size) * 1000000 / elapsed.Nanoseconds()))
}


func TestEmpty(t *testing.T) {
    test( t,  "empty",
              "",
              false, []float64{},
              false,
    )
}

func TestPush(t *testing.T) {
    test( t,  "push",
              "47.3 .",
              false, []float64{47.3},
              false,
    )
}

func TestLeftovers(t *testing.T) {
    test( t,  "leftovers",
              "47.3",
              true, nil,
              false,
    )
}

func TestAdd(t *testing.T) {
    test( t,  "add",
              "47 21 + .",
              false, []float64{68},
              false,
    )
}

func TestSub(t *testing.T) {
    test( t,  "sub",
              "47 21 - .",
              false, []float64{26},
              false,
    )
}

func TestMul(t *testing.T) {
    test( t,  "mul",
              "47 2 * .",
              false, []float64{94},
              false,
    )
}

func TestDiv(t *testing.T) {
    test( t,  "div",
              "94 2 / .",
              false, []float64{47},
              false,
    )
}

func TestSwap(t *testing.T) {
    test( t,  "swap",
              "47 21 SWAP . .",
              false, []float64{47, 21},
              false,
    )
}

func TestComment(t *testing.T) {
    test( t,  "comment",
              "47 (Hi; I am: a comment (quite a hard one)) 21 SWAP . .",
              false, []float64{47, 21},
              false,
    )
}

func TestDefinition(t *testing.T) {
    test( t,  "definition",
              "3 :five 5; :two five 3 -; five + two + .",
              false, []float64{10},
              false,
    )
}

func TestRecursiveDefinition(t *testing.T) {
    test( t,  "recursive definition (error)",
              ":here there; :there yonder; :yonder here; here .",
              true, nil,
              false,
    )
}

func TestIfElse(t *testing.T) {
    test( t,  "if-else",
              "4 11 10 > IF 1 . ELSE 2 . THEN 11 10 < IF 1 . ELSE 2 . THEN .",
              false, []float64{1, 2, 4},
              false,
    )
}

func TestNestedIfElse(t *testing.T) {
    test( t,  "nested if-else",
              "4 11 10 > IF 11 10 > IF 1 . ELSE 2 . THEN ELSE 3 . THEN 11 10 < IF 1 . ELSE 11 10 < IF 2 . ELSE 3 . THEN THEN .",
              false, []float64{1, 3, 4},
              false,
    )
}

func TestChoose(t *testing.T) {
    test( t,  "choose",
              "3 FROM 7, 8, 9, 10, 11, 12 CHOOSE .",
              false, []float64{10},
              false,
    )
}

func TestNestedChoose(t *testing.T) {
    test( t,  "choose",
              "3 FROM 7, 8, 9, 3 FROM 10, 1 FROM 99, 100 CHOOSE, 12, 1 FROM 20, 0 FROM 31 CHOOSE, 21 CHOOSE CHOOSE, 13, 14 CHOOSE .",
              false, []float64{31},
              false,
    )
}

func TestOscillators(t *testing.T) {
    test( t,  "oscillators",
              "0.25 SIN . 0.25 SAW . 0.25 SQ . 0.25 TR .",
              false, []float64{1, -0.5, 1, 0},
              false,
    )
}

func TestOptimize(t *testing.T) {
    test( t,  "optimize",
              "[400 2 * 800 /] [96 30 - [12 21 +] /] +.",
              false, []float64{3},
              false,
    )
}

func TestOptimizeBaseline(t *testing.T) {
    test( t,  "optimize baseline",
              "400 2 * 800 / 96 30 - 12 21 + / +.",
              false, []float64{3},
              false,
    )
}

func TestOptimizeUnderrun(t *testing.T) {
    test( t,  "optimize underrun (error)",
              "4 [2 *] 3 +.",
              true, nil,
              false,
    )
}

func TestOptimizeUnfinished(t *testing.T) {
    test( t,  "optimize unfinished",
              "[4 2 *.",
              false, []float64{8},
              false,
    )
}

func TestOptimizeOutput(t *testing.T) {
    test( t,  "optimize with output (error)",
              "[4 2 * .] 3 +.",
              true, nil,
              false,
    )
}

func TestImport(t *testing.T) {
    test( t,  "import",
              "::import (comment); imported.",
              false, []float64{57},
              false,
    )
}

func TestConstant(t *testing.T) {
    test( t,  "constant",
              ":a 47 dup constant b! 11 +; a. b? b.",
              false, []float64{58, 47, 1000},
              false,
    )
}

func TestKeep(t *testing.T) {
    test( t,  "keep",
              ":a 47 dup keep b 11 +; a. b? b.",
              false, []float64{58, 47, 1000},
              false,
    )
}

func TestConstantUndefined(t *testing.T) {
    test( t,  "constant-undefined",
              "50?",
              true, nil,
              false,
    )
}

func TestConstantAlreadyDefined(t *testing.T) {
    test( t,  "constant-already-defined",
              "40 constant a! 80 a!",
              true, nil,
              false,
    )
}

func TestOld(t *testing.T) {
    test( t,  "old",
              "1 t+ dup. keep a a 1 old .",
              false, []float64{2, 0},
              false,
    )
    // old value is before the beginning of time hence 0 (t starts at 1)
    // TODO actually test the old value after running a bit
}

func TestDelta(t *testing.T) {
    test( t,  "old",
              "1 t+ dup. keep a a delta .",
              false, []float64{2, 0},
              false,
    )
}

/*
const BENCHMARK_FILE = "tests/gloucester.d4"

func TestFile(t *testing.T) {
    test_file( t,  "benchmark-validate",
              BENCHMARK_FILE,
              false, []float64{},
              true,
    )
}

func TestBenchmark(t *testing.T) {
    test_fill32( t,  "benchmark-single",
              BENCHMARK_FILE,
              false, nil,
              100000,
              1,
    )
}

func TestBenchmarkParallel(t *testing.T) {
    test_fill32( t,  "benchmark-parallel-4",
              BENCHMARK_FILE,
              false, nil,
              100000,
              4,
    )
}
*/
