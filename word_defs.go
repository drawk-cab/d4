package d4

type Word struct {
    name string
    opcode float64    // so we can stick everything in a big array
    t_dependent bool
    needs int         // how many values must be on the stack
}

const W_NOOP = 0xff
const W_NUMBER = 0xf0 // used to signal the next item is a number not an opcode
const W_OUTPUT = 0xf1
const W_CLIP = 0xf2

const W_EOF = 0x00 // stop

const W_BEGIN_COMMENT = 0xe0
const W_END_COMMENT = 0xe1
const W_BEGIN_DEF = 0xe2
const W_END_DEF = 0xe3
const W_IF = 0xe4
const W_THEN = 0xe5
const W_ELSE = 0xe6
const W_CONSTANT = 0xe7
const W_VARIABLE = 0xe8
const W_LOOP = 0xe9
const W_CHOOSE = 0xea
const W_FROM = 0xeb
const W_CHOOSE_SEP = 0xec

const W_FALSE = 0x01
const W_TRUE = 0x02
const W_PLUS = 0x03
const W_MINUS = 0x04
const W_TIMES = 0x05
const W_DIVIDE = 0x06
const W_MOD = 0x07
const W_DMOD = 0x08

const W_EQUALS = 0x10
const W_GREATER = 0x11
const W_LESS = 0x12
const W_NOT = 0x13
const W_AND = 0x14
const W_OR = 0x15
const W_MAX = 0x16
const W_MIN = 0x17

const W_DUP = 0x20
const W_DDUP = 0x21
const W_OVER = 0x22
const W_DROP = 0x23
const W_NIP = 0x24
const W_TUCK = 0x25
const W_SWAP = 0x26
const W_ROT = 0x27

const W_HZ = 0x30
const W_BPM = 0x31
const W_S = 0x32

const W_FLAT = 0x40
const W_SHARP = 0x41
const W_HIGH = 0x42
const W_LOW = 0x43

const W_ON = 0x50

const W_T = 0x80
const W_SIN = 0x81
const W_SAW = 0x82
const W_TR = 0x83
const W_PULSE = 0x84
const W_SQ = 0x85

var WORDS = map[string]Word{

    "NOOP":     Word{ "NOOP", W_NOOP, false, 0 },
    ".":        Word{ ".", W_OUTPUT, false, 1 },
    "CLIP":     Word{ "CLIP", W_CLIP,   false, 1 },

    "(":        Word{ "(", W_BEGIN_COMMENT,  false, 0 },
    ")":        Word{ ")", W_END_COMMENT,  false, 0 },
    ":":        Word{ ":", W_BEGIN_DEF,  false, 0 },
    ";":        Word{ ";", W_END_DEF,  false, 0 },
    "IF":       Word{ "IF", W_IF,  false, 1 },
    "THEN":     Word{ "THEN", W_THEN,  false, 0 },
    "ELSE":     Word{ "ELSE", W_ELSE,  false, 0 },
    "CONSTANT": Word{ "CONSTANT", W_CONSTANT,  false, 1 },
    "VARIABLE": Word{ "VARIABLE", W_VARIABLE,  false, 1 },
    "LOOP":     Word{ "LOOP", W_LOOP,  false, 0 },
    "CHOOSE":   Word{ "CHOOSE", W_CHOOSE,  false, 0 },
    "FROM":     Word{ "FROM", W_FROM,  false, 1 },
    ",":        Word{ ",", W_CHOOSE_SEP, false, 0 },

    "FALSE":    Word{ "FALSE", W_FALSE,    false, 0 },
    "TRUE":     Word{ "TRUE", W_TRUE,    false, 0 },
    "+":        Word{ "+", W_PLUS,    false, 2 },
    "-":        Word{ "-", W_MINUS,    false, 2 },
    "*":        Word{ "*", W_TIMES,    false, 2 },
    "/":        Word{ "/", W_DIVIDE,    false, 2 },
    "MOD":      Word{ "MOD", W_MOD,    false, 2 },
    "DMOD":     Word{ "DMOD", W_DMOD,    false, 2 },

    "=":        Word{ "=", W_EQUALS,    false, 2 },
    ">":        Word{ ">", W_GREATER,    false, 2 },
    "<":        Word{ "<", W_LESS,    false, 2 },
    "NOT":      Word{ "NOT", W_NOT,    false, 1 },
    "AND":      Word{ "AND", W_AND,    false, 2 },
    "OR":       Word{ "OR", W_OR,    false, 2 },
    "MAX":      Word{ "MAX", W_MAX,    false, 2 },
    "MIN":      Word{ "MIN", W_MIN,    false, 2 },


    "DUP":      Word{ "DUP", W_DUP,    false, 1 },
    "DDUP":     Word{ "DDUP", W_DDUP,    false, 2 },
    "OVER":     Word{ "OVER", W_OVER,    false, 2 },
    "DROP":     Word{ "DROP", W_DROP,    false, 1 },
    "NIP":      Word{ "NIP", W_NIP,    false, 2 },
    "TUCK":     Word{ "TUCK", W_TUCK,    false, 1 },
    "SWAP":     Word{ "SWAP", W_SWAP,    false, 2 },
    "ROT":      Word{ "ROT", W_ROT,    false, 3 },

    "HZ":       Word{ "HZ", W_HZ,    false, 1 },
    "BPM":      Word{ "BPM", W_BPM,    false, 1 },
    "S":        Word{ "S", W_S,    false, 1 },

    "FLAT":     Word{ "FLAT", W_FLAT,    false, 1 },
    "SHARP":    Word{ "SHARP", W_SHARP,    false, 1 },
    "#":        Word{ "#", W_SHARP,    false, 1 },
    "HIGH":     Word{ "HIGH", W_HIGH,    false, 1 },
    "'":        Word{ "'", W_HIGH,    false, 1 },
    "LOW":      Word{ "LOW", W_LOW,    false, 1 },
    "_":        Word{ "_", W_LOW,    false, 1 },

    "ON":       Word{ "ON", W_ON,    false, 3 },

    "T":        Word{ "T", W_T,    true, 0 },
    "SIN":      Word{ "SIN", W_SIN,    true, 1 },
    "SAW":      Word{ "SAW", W_SAW,    true, 1 },
    "TR":       Word{ "TR", W_TR,    true, 1 },
    "PULSE":    Word{ "PULSE", W_PULSE,    true, 2 },
    "SQ":       Word{ "SQ", W_SQ,    true, 1 },
}
