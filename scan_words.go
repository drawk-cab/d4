package d4

import (
    "unicode"
    "unicode/utf8"
)

const R_LETTER = 0
const R_DIGIT = 1
const R_OTHER = 2

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
            seq = R_LETTER
        case unicode.IsDigit(r), r == '.':
            seq = R_DIGIT
        default:
            seq = R_OTHER
    }

    // Scan until rune not matching current set.
    for width, i := 0, start; i < len(data); i += width {
        r, width = utf8.DecodeRune(data[i:])
        if (seq == R_OTHER && i != start) ||
           ((unicode.IsLetter(r) || r == '_') != (seq == R_LETTER)) ||
           ((unicode.IsDigit(r) || r == '.') != (seq == R_DIGIT)) ||
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
