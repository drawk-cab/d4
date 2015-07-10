package d4

type Word struct {
    opcode float64 // so we can stick everything in a big array
    t_dependent bool
}

const W_NOOP = 0xff
const W_NUMBER = 0xfe // used to signal the next item is a number not an opcode
const W_EOF = 0x00 // stop

const W_BEGIN_COMMENT = 0xf0
const W_END_COMMENT = 0xf1
const W_BEGIN_DEF = 0xf2
const W_END_DEF = 0xf3
const W_IF = 0xf4
const W_THEN = 0xf5
const W_ELSE = 0xf6
const W_CONSTANT = 0xf7
const W_VARIABLE = 0xf8
const W_LOOP = 0xf9

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

    // . as noop is consistent with Forth because in d4,
    // for a value to be output it merely needs to be on the stack.

    ".":        Word{ W_NOOP, false },
    "NOOP":     Word{ W_NOOP, false },

    "(":        Word{ W_BEGIN_COMMENT,  false },
    ")":        Word{ W_END_COMMENT,  false },
    ":":        Word{ W_BEGIN_DEF,  false },
    ";":        Word{ W_END_DEF,  false },
    "IF":       Word{ W_IF,  false },
    "THEN":     Word{ W_THEN,  false },
    "ELSE":     Word{ W_ELSE,  false },
    "CONSTANT": Word{ W_CONSTANT,  false },
    "VARIABLE": Word{ W_VARIABLE,  false },
    "LOOP":     Word{ W_LOOP,  false },

    "FALSE":    Word{ W_FALSE,    false },
    "TRUE":     Word{ W_TRUE,    false },
    "+":        Word{ W_PLUS,    false },
    "-":        Word{ W_MINUS,    false },
    "*":        Word{ W_TIMES,    false },
    "/":        Word{ W_DIVIDE,    false },
    "MOD":      Word{ W_MOD,    false },
    "DMOD":     Word{ W_DMOD,    false },

    "=":        Word{ W_EQUALS,    false },
    ">":        Word{ W_GREATER,    false },
    "<":        Word{ W_LESS,    false },
    "NOT":      Word{ W_NOT,    false },
    "AND":      Word{ W_AND,    false },
    "OR":       Word{ W_OR,    false },
    "MAX":      Word{ W_MAX,    false },
    "MIN":      Word{ W_MIN,    false },


    "DUP":      Word{ W_DUP,    false },
    "DDUP":     Word{ W_DDUP,    false },
    "OVER":     Word{ W_OVER,    false },
    "DROP":     Word{ W_DROP,    false },
    "NIP":      Word{ W_NIP,    false },
    "TUCK":     Word{ W_TUCK,    false },
    "SWAP":     Word{ W_SWAP,    false },
    "ROT":      Word{ W_ROT,    false },

    "HZ":       Word{ W_HZ,    false },
    "BPM":      Word{ W_BPM,    false },
    "S":        Word{ W_S,    false },

    "FLAT":     Word{ W_FLAT,    false },
    "SHARP":    Word{ W_SHARP,    false },
    "#":        Word{ W_SHARP,    false },
    "HIGH":     Word{ W_HIGH,    false },
    "'":        Word{ W_HIGH,    false },
    "LOW":      Word{ W_LOW,    false },
    ",":        Word{ W_LOW,    false },

    "ON":       Word{ W_ON,    false },

    "T":        Word{ W_T,    true },
    "SIN":      Word{ W_SIN,    true },
    "SAW":      Word{ W_SAW,    true },
    "TR":       Word{ W_TR,    true },
    "PULSE":    Word{ W_PULSE,    true },
    "SQ":       Word{ W_SQ,    true },
}
