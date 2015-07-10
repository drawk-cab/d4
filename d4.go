package d4

import (
    "strings"
    "math"
    "fmt"
    "bufio"
    "io"
    "strconv"
    "unicode"
    "unicode/utf8"
)

const sampleRate = 22050

const LETTER = 0
const DIGIT = 1
const OTHER = 2

const NORMAL = 0
const COLON = 1
const DEF = 2
const CONSTANT = 3
const DELAY = 4
const COMMENT = 5
const IF_TRUE = 6
const IF_FALSE = 7

/* loop length in seconds, will get a click after this, default to 1 day(!) */
const loop = 60 * 60 * 24

var semitone = math.Pow(2, 1.0/12)

var SEC = float64(loop)
var BPM = float64(loop / 60)

// ScanForthWords is a split function for a Scanner that 
// looks for a d4 word and returns it.
// It will never return an empty string.

func ScanForthWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
    // Skip leading spaces.
    start := 0
    var r rune

    for width := 0; start < len(data); start += width {
        r, width = utf8.DecodeRune(data[start:])
        if !unicode.IsSpace(r) {
            break
        }
    }

    if atEOF && len(data) == 0 {
        return 0, nil, nil
    }

    var seq int
    switch {
        case unicode.IsLetter(r), r == '_':
            seq = LETTER
        case unicode.IsDigit(r), r == '.':
            seq = DIGIT
        default:
            seq = OTHER
    }

    // Scan until rune not matching current set.
    for width, i := 0, start; i < len(data); i += width {
        r, width = utf8.DecodeRune(data[i:])
        if (seq == OTHER && i != start) ||
           ((unicode.IsLetter(r) || r == '_') != (seq == LETTER)) ||
           ((unicode.IsDigit(r) || r == '.') != (seq == DIGIT)) ||
           (unicode.IsSpace(r)) {
                return i, data[start:i], nil
        }
    }

    // If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
    if atEOF && len(data) > start {
        return len(data), data[start:], nil
    }

    // Request more data.
    return 0, nil, nil
}

type Machine struct {
    iter int
    sampleRate float64
    step float64
    code []string
    words map[string][]string
    variables map[string][]float64
    constants map[string][]float64
    old_variables []map[string][]float64
}

func CompileString(in string, sampleRate float64) *Machine {
    return Compile(strings.NewReader(in), sampleRate)
}

func Compile(in io.Reader, sampleRate float64) *Machine {

    // Set iter nonzero to avoid zeros everywhere during dummy run
    s := &Machine{0, sampleRate, 1/(loop*sampleRate), nil, nil, nil, nil, nil}

    s.code, s.words = s.read_code(in)

    // TODO: optimize

    return s
}


func (m Machine) read_code(in io.Reader) ([]string, map[string][]string) {
    scanner := bufio.NewScanner(in)

    scanner.Split(ScanForthWords)

    var code []string
    var words map[string][]string
    var new_word string

    words = make(map[string][]string)

    mode := []int{NORMAL}

    for scanner.Scan() {
        w := scanner.Text()
        switch mode[len(mode)-1] {
            case COLON:
                new_word = w

                _, exists := words[new_word]
                if exists {
                    panic("Word "+new_word+" has already been defined")
                } else {                
                    words[new_word] = nil
                    mode[len(mode)-1] = DEF
                }

            case DEF:
                switch w {
                    case ";":
                        mode = mode[:len(mode)-1]
                    case "(":
                        mode = append(mode, COMMENT)
                    default:
                        words[new_word] = append(words[new_word], w)
                }

            case COMMENT:
                switch w {
                    case "(":
                        mode = append(mode, COMMENT)
                    case ")":
                        mode = mode[:len(mode)-1]
                }

            case NORMAL:
                switch w {
                    case ":":
                        mode = append(mode, COLON)
                    case "(":
                        mode = append(mode, COMMENT)
                    default:
                        code = append(code, w)
                }
        }
    }

    fmt.Println("Code:  ",code)
    fmt.Println("Words: ",words)
    chk(scanner.Err())

    return code, words
}

func (m *Machine) Run(debug bool) []float64 {
    var stack []float64
    _, phase := math.Modf( float64(m.iter) * m.step )
    return m.run_code(stack, phase, m.code, debug)
}

func (m *Machine) run_code(stack []float64, phase float64, code []string, debug bool) []float64 {
    var pop float64
    var mode = NORMAL

    if debug == true {
        fmt.Println("==",code)
    }

    for _, w := range code {
        l := len(stack)-1
        switch mode {
            case IF_FALSE:
                if w == "THEN" || w == "ELSE" {
                    mode = NORMAL
                }

            case NORMAL:
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
                            mode = IF_FALSE
                        }
                    case "THEN":
                        // Must have been executing an ELSE clause, do nothing
                    case "ELSE":
                        // Must have been executing an IF clause, skip to THEN
                        mode = IF_FALSE

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
                        mode = CONSTANT

                    /* Useful words */
        
                    case "MAX":
                        pop, stack = stack[l], stack[:l]
                        stack[l-1] = math.Max(pop,stack[l-1])
                    case "MIN":
                        pop, stack = stack[l], stack[:l]
                        stack[l-1] = math.Min(pop,stack[l-1])

                    /* musical words */

                    case "DELAY":
                        mode = DELAY

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
                        stack[l] *= semitone
                    case "FLAT":
                        stack[l] /= semitone
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
                            stack = m.run_code(stack, phase, word_def, debug)
                        } else {
                            num, err := strconv.ParseFloat(w, 64)
                            if err != nil {
                                panic( "Unknown word: "+w )
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
    return stack
}

func chk(err error) {
    if err != nil {
        panic(err)
    }
}
