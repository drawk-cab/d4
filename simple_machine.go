package d4

import (
    "math"
    "fmt"
    "strconv"
    "bufio"
    "io"
)

const SM_NORMAL = 0
const SM_COLON = 1
const SM_DEF = 2
const SM_CONSTANT = 3
const SM_COMMENT = 4
const SM_IF_FALSE = 5

type SimpleMachine struct {
    iter int
    sampleRate float64
    step float64
    clip float64
    code []string
    words map[string][]string
    variables map[string][]float64
    constants map[string][]float64
    old_variables []map[string][]float64
}

func NewSimpleMachine( sampleRate float64 ) *SimpleMachine {
    return &SimpleMachine{0, sampleRate, 1/(LOOP*sampleRate), 1.0, nil, nil, nil, nil, nil}
}

func (m *SimpleMachine) Init() error {
    return nil
}

func (m *SimpleMachine) Program( in io.Reader ) error {
    scanner := bufio.NewScanner(in)

    scanner.Split(ScanForthWords)

    var code []string
    var words map[string][]string
    var new_word string

    words = make(map[string][]string)

    mode := []int{SM_NORMAL}

    for scanner.Scan() {
        w := scanner.Text()
        switch mode[len(mode)-1] {
            case SM_COLON:
                new_word = w

                _, exists := words[new_word]
                if exists {
                    panic("Word "+new_word+" has already been defined")
                } else {                
                    words[new_word] = nil
                    mode[len(mode)-1] = SM_DEF
                }

            case SM_DEF:
                switch w {
                    case ";":
                        mode = mode[:len(mode)-1]
                    case "(":
                        mode = append(mode, SM_COMMENT)
                    default:
                        words[new_word] = append(words[new_word], w)
                }

            case SM_COMMENT:
                switch w {
                    case "(":
                        mode = append(mode, SM_COMMENT)
                    case ")":
                        mode = mode[:len(mode)-1]
                }

            case SM_NORMAL:
                switch w {
                    case ":":
                        mode = append(mode, SM_COLON)
                    case "(":
                        mode = append(mode, SM_COMMENT)
                    default:
                        code = append(code, w)
                }
        }
    }

    err := scanner.Err()

    if err == nil {
        m.code = code
        fmt.Println("Code:  ",code)
        m.words = words
        fmt.Println("Words: ",words)
    }

    return err
}

func (m *SimpleMachine) Fill32( buf []float32 ) error {
    for i := range buf {

        stack, err := m.Run()

        if (err != nil) {
            return err
        }

        /* Add up whatever is left on the stack */
        r := float64(0)
        for _, s := range stack {
            r += s
        }

        /* Tweak the scale if it's clipping */
        if r > m.clip {
            m.clip = r
        }
        buf[i] = float32(r / m.clip)
    }

    return nil
}

func (m *SimpleMachine) Run() ([]float64, error) {
    var stack []float64
    _, phase := math.Modf( float64(m.iter) * m.step )
    m.iter += 1
    return m.run_code(stack, phase, m.code, false)
}

func (m *SimpleMachine) run_code( stack []float64, phase float64, code []string, debug bool ) ([]float64, error) {
    var err error
    var pop float64
    var mode = SM_NORMAL

    if debug == true {
        fmt.Println("==",code)
    }

    for _, w := range code {
        l := len(stack)-1
        switch mode {
            case SM_IF_FALSE:
                if w == "THEN" || w == "ELSE" {
                    mode = SM_NORMAL
                }

            case SM_NORMAL:
                switch w {

                    /* Forth words */

                    case "TRUE":
                        stack = append(stack, 1)
                    case "FALSE":
                        stack = append(stack, 0)

                    case "IF":
                        pop, stack = stack[l], stack[:l]
                        if pop != 0 {
                            // Test succeeded, carry on
                        } else {
                            mode = SM_IF_FALSE
                        }
                    case "THEN":
                        // Must have been executing an ELSE clause, do nothing
                    case "ELSE":
                        // Must have been executing an IF clause, skip to THEN
                        mode = SM_IF_FALSE

                    case "+":
                        pop, stack = stack[l], stack[:l]
                        stack[l-1] += pop
                    case "-":
                        pop, stack = stack[l], stack[:l]
                        stack[l-1] -= pop
                    case "*":
                        pop, stack = stack[l], stack[:l]
                        stack[l-1] *= pop
                    case "/":
                        pop, stack = stack[l], stack[:l]
                        stack[l-1] /= pop
                    case "MOD":
                        pop, stack = stack[l], stack[:l]
                        stack[l-1] = math.Mod( stack[l-1], pop )
        
                    case "=":
                        if stack[l] == stack[l-1] {
                            stack = append(stack, 1)
                        } else {
                            stack = append(stack, 0)
                        }

                    case ">":
                        pop, stack = stack[l], stack[:l]
                        if pop > stack[l-1] {
                            stack[l-1] = 1
                        } else {
                            stack[l-1] = 0
                        }

                    case "<":
                        pop, stack = stack[l], stack[:l]
                        if pop < stack[l-1] {
                            stack[l-1] = 1
                        } else {
                            stack[l-1] = 0
                        }

                    case "NOT":
                        if stack[l] == 0 {
                            stack[l] = 1
                        } else {
                            stack[l] = 0
                        }

                    case "AND":
                        pop, stack = stack[l], stack[:l]
                        if pop != 0 && stack[l-1] != 0 {
                            stack[l-1] = 1
                        } else {
                            stack[l-1] = 0
                        }

                    case "OR":
                        pop, stack = stack[l], stack[:l]
                        if pop != 0 || stack[l-1] != 0 {
                            stack[l-1] = 1
                        } else {
                            stack[l-1] = 0
                        }

                    case "DUP":
                        stack = append(stack, stack[l])

                    case "DDUP":
                        stack = append(stack, stack[l-1], stack[l])

                    case "OVER":
                        stack = append(stack, stack[l-1])

                    case "DROP":
                        stack = stack[:l]

                    case "NIP":
                        pop, stack = stack[l], stack[:l]
                        stack[l-1] = pop

                    case "TUCK":
                        stack = append(stack, stack[l])
                        stack[l], stack[l-1] = stack[l-1], stack[l]

                    case "SWAP":
                        stack[l], stack[l-1] = stack[l-1], stack[l]

                    case "ROT":
                        stack[l], stack[l-1], stack[l-2] = stack[l-2], stack[l], stack[l-1]

                    case "CONSTANT":
                        mode = SM_CONSTANT

                    case "LOOP":
                        // TODO 

                    /* Useful words */
        
                    case "MAX":
                        pop, stack = stack[l], stack[:l]
                        stack[l-1] = math.Max(pop,stack[l-1])

                    case "MIN":
                        pop, stack = stack[l], stack[:l]
                        stack[l-1] = math.Min(pop,stack[l-1])

                    case ".", "NOOP":
                        // . as noop is consistent with Forth because in d4,
                        // for a value to be output it merely needs to be on the stack.

                    /* musical words */

                    case "HZ":
                        stack[l] *= SEC

                    case "BPM":
                        stack[l] *= BPM

                    case "S":
                        stack[l] /= SEC

                    case "T":
                        stack = append(stack, phase)

                    case "ON":
                        /* (time, length, base -- age, on (if on) OR off (if off) */
                        var sched, dur, now float64
                        sched, dur, now, stack = stack[l-2], stack[l-1], stack[l], stack[:l-1]
                        age := now - sched
                        if age > 0 && age < dur {
                            stack[l-2] = age
                            stack = append(stack, 1)
                        } else {
                            stack[l-2] = 0
                        }

                    /* intervals */

                    case "#","SHARP":
                        stack[l] *= SEMITONE
                    case "FLAT":
                        stack[l] /= SEMITONE
                    case "'","HIGH":
                        stack[l] *= 2
                    case ",","LOW":
                        stack[l] /= 2

                    /* oscillators */

                    case "SIN":
                        stack[l] = math.Sin(stack[l] * phase * 2 * math.Pi)

                    case "SAW":
                        stack[l] = math.Mod(stack[l] * phase * 2, 2) - 1

                    case "DIA":
                        _, frac := math.Modf(stack[l] * phase)
                        if frac < 0.5 {
                            stack[l] = frac * 4 - 1
                        } else {
                            stack[l] = 3 - frac * 4
                        }

                    case "SQ":
                        _, frac := math.Modf(stack[l] * phase)
                        if frac < 0.5 {
                            stack[l] = 1
                        } else {
                            stack[l] = -1
                        }

                    default:
                        word_def, ok := m.words[w]
                        if ok {
                            if debug == true {
                                fmt.Println(">> ",w)
                            }
                            stack, err = m.run_code(stack, phase, word_def, debug)
                        } else {
                            num, err := strconv.ParseFloat(w, 64)
                            if err != nil {
                                return stack, fmt.Errorf("unknown word: %s", w)
                            }
                            stack = append(stack, num)
                        }
                }
        }
        if debug == true {
            fmt.Println(w,"--",stack)
        }
    }
    if debug == true {
        fmt.Println("<<",stack)
    }

    return stack, err
}
