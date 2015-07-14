package d4

import (
    "strings"
    "math"
    "fmt"
    "strconv"
    "bufio"
    "io"
)

type Job struct {
    id int
    iter int64
}

type JobResult struct {
    id int
    value []float64
    err error
}

const M_NORMAL = 0
const M_COLON = 1
const M_DEF = 2
const M_CONSTANT = 3
const M_COMMENT = 4
const M_IF_FALSE = 5
const M_CHOOSE_FALSE = 6
const M_LITERAL = 7
const M_IMPORT = 8

type OpcodeMachine struct {
    MachineData
    step float64
    code []float64
    words map[string][]string
    save_addr int
    saves []map[float64]float64
    save_ptr int
    opcode_info map[float64]Word
}

func NewOpcodeMachine( sample_rate float64, save_s float64, clip float64, imports map[string]string ) *OpcodeMachine {

    // import names are case insensitive and stored as capitals
    upper_imports := map[string]string{}
    for name, code := range imports {
        upper_imports[strings.ToUpper(name)] = code
    }

    save_len := int(save_s * sample_rate)

    return &OpcodeMachine{MachineData{0, sample_rate, save_len, clip, upper_imports},
                          1/(LOOP*sample_rate), nil, nil, 1000, nil, save_len-1, nil}
}

func (m *OpcodeMachine) GetData() MachineData {
    return m.MachineData
}

func (m *OpcodeMachine) Init(clone_from Machine) error {

    m.opcode_info = map[float64]Word{
        W_NUMBER: Word{ "n", W_NUMBER, false, 0 }, // this opcode is created, not supplied
    }

    if clone_from != nil {
        m.MachineData = clone_from.GetData()
        m.step = 1/(LOOP*m.sample_rate)
    }

    m.saves = make([]map[float64]float64, m.save_len)

    return nil
}

func (m *OpcodeMachine) Program( in io.Reader ) error {

    words, need_imports, err := m.read(in)
    m.words = words

    if err != nil {
        return err
    }

    for _, name := range need_imports {

        in = strings.NewReader( m.imports[name] )
        new_words, new_imports, err := m.read( in )

        if err != nil {
            return err
        }

        if len(new_imports)>0 {
            return fmt.Errorf("Program error: import %s tried to import %s", name, new_imports)
            // TODO: allow nested imports
        }

        for w, defn := range new_words {
            _, ok := m.words[w]
            if !ok {
                m.words[w] = defn
            } else {
                // don't overwrite existing word
            }
        }
    }

    // We now have a set of word definitions (counting '' for everything outside a word definition)
    // which we can translate into opcodes

    var breadcrumb []string = []string{}
    var code = []float64{}

    code, err = m.compile(code, "", breadcrumb)

    if err != nil {
        return err
    }
    
    code = append(code, W_EOF)

    m.code, err = m.optimize(code)

    return err
}

func (m *OpcodeMachine) read( in io.Reader ) (map[string][]string, []string, error) {

    imports := []string{}

    words := map[string][]string{ "": []string{} }

    scanner := bufio.NewScanner(in)

    scanner.Split(ScanForthWords)

    cur_word := ""
    mode := []int{M_NORMAL}

    for scanner.Scan() {
        w := strings.ToUpper(scanner.Text())
        switch mode[len(mode)-1] {

            case M_COLON:
                if w == ":" {
                    mode[len(mode)-1] = M_IMPORT
                } else {
                    cur_word = w

                    _, exists := words[cur_word]
                    if exists {
                        return words, imports, fmt.Errorf("Scan error: %s has already been defined", cur_word)
                    } else {                
                        _, exists := WORDS[cur_word]
                        if exists {
                            return words, imports, fmt.Errorf("Scan error: %s is a built-in word and cannot be redefined", cur_word)
                        } else {
                            words[cur_word] = nil
                            mode[len(mode)-1] = M_DEF
                        }
                    }
                }

            case M_CONSTANT:

                words[w] = []string{strconv.Itoa(m.save_addr)} // everything is a string at this point
                m.save_addr += 1

                words[cur_word] = append(words[cur_word], w)

                mode = mode[:len(mode)-1]

            case M_DEF:
                switch w {
                    case ":":
                        return words, imports, fmt.Errorf("Scan error: : found inside definition")
                    case ")":
                        return words, imports, fmt.Errorf("Scan error: ) found outside comment")
                    case ";":
                        cur_word = ""
                        mode = mode[:len(mode)-1]
                    case "(":
                        mode = append(mode, M_COMMENT)
                    case "SAVE", "CONSTANT":
                        mode = append(mode, M_CONSTANT)
                    default:
                        words[cur_word] = append(words[cur_word], w)
                }

            case M_IMPORT:
                switch w {
                    case ":":
                        return words, imports, fmt.Errorf("Scan error: : found inside import statement")
                    case "(":
                        mode = append(mode, M_COMMENT)
                    case ")":
                        return words, imports, fmt.Errorf("Scan error: ) found outside comment")
                    case ";":
                        mode = mode[:len(mode)-1]
                    default:
                        imports = append(imports, w)
                }

            case M_COMMENT:
                switch w {
                    case "(":
                        mode = append(mode, M_COMMENT)
                    case ")":
                        mode = mode[:len(mode)-1]
                }

            case M_NORMAL:
                switch w {
                    case ":":
                        mode = append(mode, M_COLON)
                    case "(":
                        mode = append(mode, M_COMMENT)
                    case ";":
                        return words, imports, fmt.Errorf("Scan error: ; found outside definition")
                    case ")":
                        return words, imports, fmt.Errorf("Scan error: ) found outside comment")
                    case "SAVE", "CONSTANT":
                        mode = append(mode, M_CONSTANT)
                    default:
                        words[cur_word] = append(words[cur_word], w)
                }
        }
        if DEBUG {
            fmt.Println("scan",w," -- ",mode)
        }
    }

    err := scanner.Err()

    return words, imports, err
}

func (m *OpcodeMachine) compile( code []float64, word string, breadcrumb []string ) ([]float64, error) {
    var err error

    //word = strings.ToUpper(word)

    defn, ok := m.words[word]
    if ok {

        if DEBUG {
            fmt.Println(word,"=",defn,"--",breadcrumb)
        }

        // word is a defined word
        for _, outer_word := range breadcrumb {
            if outer_word == word {
                return code, fmt.Errorf("Compile error: recursive definition %s", breadcrumb)
            }
        }

        for _, w := range defn {
            w = strings.ToUpper(w)

            word_info, ok := WORDS[w]
            if ok {
                code = append(code, word_info.opcode)
                m.opcode_info[word_info.opcode] = word_info
            } else {
                new_breadcrumb := append(breadcrumb, word)
                code, err = m.compile( code, w, new_breadcrumb )
                if err != nil {
                    return code, err
                }
            }
        }
    } else {

        num, err := strconv.ParseFloat(word, 64)
        if err != nil {
            return code, fmt.Errorf("Compile error: unknown word %s", word)
        }
        code = append(code, W_NUMBER, num)
    }
    return code, err
}

func (m *OpcodeMachine) optimize( code []float64 ) ([]float64, error) {
    var output []float64

    literal := []float64{}

    mode := []int{M_NORMAL}

    for _, w := range code {
        switch mode[len(mode)-1] {
            case M_LITERAL:
                switch w {
                    case W_BEGIN_LITERAL:
                        mode = append(mode, M_LITERAL)
                    case W_END_LITERAL:
                        mode = mode[:len(mode)-1]
                        if mode[len(mode)-1] == M_NORMAL {
                            // outside [], time to evaluate
                            literal = append(literal, W_EOF)
                            if DEBUG == true {
                                fmt.Println("Evaluating",literal)
                            }
                            literal_output, literal_stack, err := m.RunCode(literal,0)
                            if DEBUG == true {
                                fmt.Println("Replacing with",literal_stack)
                            }
                            if err != nil {
                                return output, err
                            }
                            if len(literal_output) > 0 {
                                return output, fmt.Errorf("Optimize error: attempted output from within [ ]")
                            }
                            for _, value := range literal_stack {
                                output = append(output, W_NUMBER, value)
                            }
                            literal = []float64{}
                        }
                    default:
                        literal = append(literal, w)
                }
            case M_NORMAL:
                switch w {
                    case W_BEGIN_LITERAL:
                        mode = append(mode, M_LITERAL)
                    case W_END_LITERAL:
                        return output, fmt.Errorf("Optimize error: ] found outside literal")
                    default:
                        output = append(output, w)                
                }
        }
    }

    // if EOF during a literal, tack it on the end
    for _, w := range literal {
        output = append(output, w)
    }

    return output, nil //TODO
}

func (m *OpcodeMachine) Fill32( buf []float32, workers int ) error {
    if workers == 1 {
        return m.fill32_single(buf)
    } else {
        return m.fill32_parallel(buf, workers)
    }
}

func (m *OpcodeMachine) fill32_parallel( buf []float32, workers int ) error {

    jobs := make(chan *Job, len(buf))
    results := make(chan *JobResult, len(buf))

    for w := 0; w <= workers; w++ {
        go m.work()
    }

    for i := range buf {
        jobs <- &Job{i, m.iter}
        m.iter++
    }
    close(jobs)

    for result := range results {
        if result.err != nil {
            return result.err
        }

        r := float64(0)
        for _, s := range result.value {
            r += s
        }
        buf[result.id] = float32( r / m.clip )
    }
    //fmt.Println(output, i)

    return nil
}

func (m *OpcodeMachine) fill32_single( buf []float32 ) error {
    var output []float64
    var err error

    for i := range buf {

        output, err = m.Run()

        if (err != nil) {
            return err
        }

        r := float64(0)
        for _, s := range output {
            r += s
        }

        buf[i] = float32(r / m.clip)
    }

    //fmt.Println(output, i)

    return err
}

func (m *OpcodeMachine) Run() ([]float64, error) {
    m.iter += 1
    m.save_ptr -= 1
    if m.save_ptr < 0 {
        m.save_ptr += m.save_len
    }
    m.saves[m.save_ptr] = nil

    output, _, err := m.RunCode(m.code, m.iter)
    return output, err
}

func (m *OpcodeMachine) work() (jobs chan *Job) {
    for j := range jobs {
        fmt.Println("doing job",j)
        output, _, err := m.RunCode(m.code, j.iter)
        fmt.Println("output",output,err)
        //results <- JobResult{j.id, output, err}
    }
    return nil
}

func (m *OpcodeMachine) RunCode(code []float64, iter int64) ([]float64, []float64, error) {

    output := []float64{}
    stack := []float64{}
    
    _, phase := math.Modf( float64(m.iter) * m.step )

    var err error
    var pop float64
    var w_info Word
    var choose_value int

    code_ptr := 0
    top := -1

    mode_breadcrumb := []int{}
    mode := M_NORMAL

    w := code[code_ptr]

    for w != W_EOF {
        w_info = m.opcode_info[w]
        switch mode {
            case M_IF_FALSE:
                switch w {
                    case W_THEN:
                        var old_mode int
                        old_mode, mode_breadcrumb = mode_breadcrumb[len(mode_breadcrumb)-1], mode_breadcrumb[:len(mode_breadcrumb)-1]
                        mode = old_mode
                    case W_ELSE:
                        mode = M_NORMAL
                }

            case M_CHOOSE_FALSE:
                switch w {
                    case W_CHOOSE:
                        var old_mode int
                        old_mode, mode_breadcrumb = mode_breadcrumb[len(mode_breadcrumb)-1], mode_breadcrumb[:len(mode_breadcrumb)-1]
                        mode = old_mode
                    case W_CHOOSE_SEP:
                        choose_value -= 1
                        if choose_value == 0 {
                            mode = M_NORMAL
                        }
                }

            case M_NORMAL:
                if w_info.needs > top+1 {
                    return output, stack, fmt.Errorf("Runtime error: %s needs %d items on stack, got %s", w_info.name, w_info.needs, stack)
                }
                switch w {

                    case W_NOOP, W_BEGIN_LITERAL, W_END_LITERAL:
                        // noop
                    case W_NUMBER:
                        code_ptr += 1
                        stack = append(stack, code[code_ptr])
                        top += 1
                    case W_OUTPUT:
                        pop, stack = stack[top], stack[:top]
                        top -= 1
                        output = append(output, pop)
                    case W_CLIP:
                        pop, stack = stack[top], stack[:top]
                        m.clip = pop

                    /* Runtime control */

                    case W_IF:
                        pop, stack = stack[top], stack[:top]
                        top -= 1
                        mode_breadcrumb = append(mode_breadcrumb, mode)
                        if pop == 0 {
                            mode = M_IF_FALSE
                        }

                    case W_THEN:
                        var old_mode int
                        old_mode, mode_breadcrumb = mode_breadcrumb[len(mode_breadcrumb)-1], mode_breadcrumb[:len(mode_breadcrumb)-1]
                        mode = old_mode

                    case W_ELSE:
                        // Must have been executing an IF clause, skip to THEN
                        mode = M_IF_FALSE

                    case W_FROM:
                        choose_value, stack = int(stack[top]), stack[:top]
                        top -= 1
                        mode_breadcrumb = append(mode_breadcrumb, mode)
                        if choose_value != 0 {
                            mode = M_CHOOSE_FALSE
                        }

                    case W_CHOOSE_SEP:
                        choose_value -= 1
                        if choose_value != 0 {
                            mode = M_CHOOSE_FALSE
                        }

                    case W_CHOOSE:
                        var old_mode int
                        old_mode, mode_breadcrumb = mode_breadcrumb[len(mode_breadcrumb)-1], mode_breadcrumb[:len(mode_breadcrumb)-1]
                        mode = old_mode

                    /* Memory */

                    case W_PEEK:
                        if stack[top] < 1000 || stack[top] != math.Floor(stack[top]) {
                            return output, stack, fmt.Errorf("Runtime error: word before ? was not a save name")
                        }

                        val, ok := m.saves[m.save_ptr][stack[top]]
                        if ok {
                            stack[top] = val
                        } else {
                            return output, stack, fmt.Errorf("Runtime error: nothing at address %d at ptr %d", stack[top], m.save_ptr)
                        }

                    case W_OLD:
                        pop, stack = stack[top], stack[:top]
                        top -= 1

                        old_ptr := (m.save_ptr + int(pop)) % m.save_len

                        if old_ptr >= len(m.saves) || old_ptr < 0 {
                            return output, stack, fmt.Errorf("Runtime error: tried to read ptr %s but save_len is %s", old_ptr, m.save_len)
                        }

                        if stack[top] < 1000 || stack[top] != math.Floor(stack[top]) {
                            return output, stack, fmt.Errorf("Runtime error: bad save name passed to OLD")
                        }

                        val, ok := m.saves[old_ptr][stack[top]]
                        //fmt.Printf("Looked at address %d at ptr %d (now %d): %s", stack[top], old_ptr, m.save_ptr, val)
                        if ok {
                            stack[top] = val
                        } else {
                            stack[top] = 0
                        }

                    case W_POKE:
                        if m.save_ptr >= len(m.saves) || m.save_ptr < 0 {
                            return output, stack, fmt.Errorf("Runtime error: tried to save ptr %s but save_len is %s", m.save_ptr, len(m.saves),m.save_len)
                        }
                        if m.saves[m.save_ptr] == nil {
                            m.saves[m.save_ptr] = map[float64]float64{}
                        }

                        _, ok := m.saves[m.save_ptr][stack[top]]
                        if ok {
                            return output, stack, fmt.Errorf("Runtime error: address %d already in use", stack[top])                            
                        } else {
                            pop, stack = stack[top], stack[:top]
                            top -= 1
                            m.saves[m.save_ptr][pop] = stack[top]
                        }

                    /* Forth words */

                    case W_TRUE:
                        stack = append(stack, 1)
                        top += 1
                    case W_FALSE:
                        stack = append(stack, 0)
                        top += 1

                    case W_PLUS:
                        pop, stack = stack[top], stack[:top]
                        top -= 1
                        stack[top] += pop
                    case W_MINUS:
                        pop, stack = stack[top], stack[:top]
                        top -= 1
                        stack[top] -= pop
                    case W_TIMES:
                        pop, stack = stack[top], stack[:top]
                        top -= 1
                        stack[top] *= pop
                    case W_DIVIDE:
                        pop, stack = stack[top], stack[:top]
                        if pop == 0 {
                            return output, stack, fmt.Errorf("Runtime error: divide by zero")
                        }
                        top -= 1
                        stack[top] /= pop
                    case W_MOD:
                        pop, stack = stack[top], stack[:top]
                        top -= 1
                        stack[top] = math.Mod( stack[top], pop )

                    case W_DMOD:
                        if stack[top] == 0 {
                            return output, stack, fmt.Errorf("Runtime error: divide by zero")
                        }
                        result, remainder := math.Modf( stack[top-1] / stack[top] )
                        stack[top] = result
                        stack[top-1] = remainder * pop
        
                    case W_EQUALS:
                        pop, stack = stack[top], stack[:top]
                        top -= 1
                        if stack[top] == stack[top-1] {
                            stack = append(stack, 1)
                        } else {
                            stack = append(stack, 0)
                        }

                    case W_GREATER:
                        pop, stack = stack[top], stack[:top]
                        top -= 1
                        if stack[top] > pop {
                            stack[top] = 1
                        } else {
                            stack[top] = 0
                        }

                    case W_LESS:
                        pop, stack = stack[top], stack[:top]
                        top -= 1
                        if stack[top] < pop {
                            stack[top] = 1
                        } else {
                            stack[top] = 0
                        }

                    case W_NOT:
                        if stack[top] == 0 {
                            stack[top] = 1
                        } else {
                            stack[top] = 0
                        }

                    case W_AND:
                        pop, stack = stack[top], stack[:top]
                        top -= 1
                        if pop != 0 && stack[top] != 0 {
                            stack[top] = 1
                        } else {
                            stack[top] = 0
                        }

                    case W_OR:
                        pop, stack = stack[top], stack[:top]
                        top -= 1
                        if pop != 0 || stack[top] != 0 {
                            stack[top] = 1
                        } else {
                            stack[top] = 0
                        }

                    case W_DUP:
                        stack = append(stack, stack[top])
                        top += 1

                    case W_DDUP:
                        stack = append(stack, stack[top-1], stack[top])
                        top += 2

                    case W_OVER:
                        stack = append(stack, stack[top-1])
                        top += 1

                    case W_DROP:
                        stack = stack[:top]
                        top -= 1

                    case W_NIP:
                        pop, stack = stack[top], stack[:top]
                        top -= 1
                        stack[top] = pop

                    case W_TUCK:
                        stack = append(stack, stack[top])
                        top += 1
                        stack[top], stack[top-1] = stack[top-1], stack[top]

                    case W_SWAP:
                        stack[top], stack[top-1] = stack[top-1], stack[top]

                    case W_ROT:
                        stack[top], stack[top-1], stack[top-2] = stack[top-2], stack[top], stack[top-1]

                    case W_CONSTANT:
                        mode = M_CONSTANT

                    case W_LOOP:
                        // TODO 

                    /* Useful words */
        
                    case W_MAX:
                        pop, stack = stack[top], stack[:top]
                        top -= 1
                        stack[top] = math.Max(pop,stack[top])

                    case W_MIN:
                        pop, stack = stack[top], stack[:top]
                        top -= 1
                        stack[top] = math.Min(pop,stack[top])

                    /* musical words */

                    case W_HZ:
                        stack[top] *= SEC * phase

                    case W_BPM:
                        stack[top] *= BPM * phase

                    case W_S:
                        stack[top] *= m.sample_rate

                    case W_T:
                        stack = append(stack, float64(m.iter))
                        top += 1

                    case W_ON:
                        /* (time, length, base -- age, on (if on) OR off (if off) */
                        var sched, dur, now float64
                        sched, dur, now, stack = stack[top-2], stack[top-1], stack[top], stack[:top-1]
                        age := now - sched
                        if age > 0 && age < dur {
                            stack[top-2] = age
                            stack = append(stack, 1)
                            top -= 1
                        } else {
                            stack[top-2] = 0
                            top -= 2
                        }

                    /* intervals */

                    case W_SHARP:
                        stack[top] *= SEMITONE
                    case W_FLAT:
                        stack[top] /= SEMITONE
                    case W_HIGH:
                        stack[top] *= 2
                    case W_LOW:
                        stack[top] /= 2

                    /* oscillators */

                    case W_SIN:
                        _, frac := math.Modf(stack[top])
                        stack[top] = math.Sin(frac * 2 * math.Pi)

                    case W_SAW:
                        stack[top] = math.Mod(stack[top] * 2, 2) - 1

                    case W_TR:
                        _, frac := math.Modf(stack[top])
                        if frac < 0.5 {
                            stack[top] = frac * 4 - 1
                        } else {
                            stack[top] = 3 - frac * 4
                        }

                    case W_PULSE:
                        pop, stack = stack[top], stack[:top]
                        top -= 1
                        _, frac := math.Modf(stack[top])
                        if frac < pop {
                            stack[top] = 1
                        } else {
                            stack[top] = -1
                        }

                    case W_SQ:
                        _, frac := math.Modf(stack[top])
                        if frac < 0.5 {
                            stack[top] = 1
                        } else {
                            stack[top] = -1
                        }

                    default:
                        return output, stack, fmt.Errorf("Runtime error: unknown opcode %d", w)
                }
        }

        if DEBUG == true {
            fmt.Println(w_info.name,": stack=",stack,top,"mode=",mode_breadcrumb,mode,"choose=",choose_value,"out=",output)
        }

        code_ptr += 1
        w = code[code_ptr]
    }
    if DEBUG == true {
        fmt.Println("<<",stack,"out",output)
    }

    return output, stack, err
}
