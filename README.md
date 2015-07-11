# d4
Diminished Forth definition and interpreter

This is the tiny language floatbeat uses to define sounds, designed for 
fun and live coding of music. It's like Forth but simpler.

## Principles

### Words

A word is a contiguous sequence of (unicode.IsDigit or .) or (unicode.IsLetter or _) or a single other non-whitespace rune.

Whitespace is only needed in order to separate adjacent sequences of digits or letters.

    2DUP --> [2] [DUP]
    2.4* --> [2.4] [*]
    a__b!=5 --> [a__b] [!] [=] [5]

Words are case sensitive.

There are no built-in words beginning with a lowercase letter.

### There is only one data type, number (float64)

### Programs are run iteratively.

A typical run will consist of many interations of the supplied code.

The iteration number is available in the built-in word `T`.

Built-in words `S`, `BPM`, `HZ` convert `T` to time units based
on a supplied sample rate (default 22050)

**TODO** All defined `CONSTANT`s are saved at the end of the first iteration (`T` = 0).
They are not evaluated again.

**TODO** All defined `VARIABLE`s are saved at the end of each iteration.

**TODO** These values can be accessed in later iterations using `AGO`.

The word `.` pops the value off the top of the stack and adds it to an output stack ready to be returned.

## Standard Forth words

* `TRUE` === `1`
* `FALSE` === `0`
* `+` `-` `*` `/`
* `=` `>` `<` `NOT` `OR` `AND`
* `:` `;` define a word
* `IF` `THEN` `ELSE` : can't be nested, define words if you want to do this
* `DROP` ( x -- )
* `DUP` ( x -- x x )
* `DDUP` ( x y -- x y x y ) : standard Forth `2DUP`
* `OVER` ( x y -- x y x )
* `NIP` ( x y -- y )
* `TUCK` ( x y -- y x y )
* `SWAP` ( x y -- y x )
* `ROT` ( x y z -- y z x )
* `DMOD` ( number, modulus -- remainder, floor ) : standard Forth `/MOD`
* **TODO** `CONSTANT name` : Evaluated on first run only, must not depend on `T`
* **TODO** `VARIABLE`
* **TODO** `!` `@`

## Built-in units

* `HZ`
* `S`
* `BPM`

## Extra words

* **TODO** `.`, `NOOP` : noop

* **TODO** `AGO` ( variable_name time -- value ) : return the value the named `VARIABLE` ended up with after the iteration the specified time ago.

* **TODO** `LAST` ( variable_name -- value ) : return the value the named `VARIABLE` ended up with in the previous iteration.

* **TODO** `T` : iteration number (use `S`, `BPM`, `HZ` to convert to seconds)

* `ON` ( schedule_t, duration, base_t -- 0 or age, 1 ) : Is the note of length `duration` scheduled for `schedule_t` currently in progress at time `base_t`, and if so, how old is it?

    _example_ `1S 0.5S T ON IF 440HZ SIN THEN` plays an A for 0.5 sec, 1 sec after the start.

* **TODO** `FROM`...`CHOOSE` : treat TOS as a pointer and execute the word with that index between `FROM` and `CHOOSE` . So `2 FROM a b c d CHOOSE` would execute only `c` because `c` has index 2 in the set.

    _example_ `T 1S 4 DMOD FROM play_c play_d play_e play_f CHOOSE` will execute `play_c` in the first second, `play_d` in the second, and so on. Because of the `DMOD`, the age of each note will be on TOS for the play routines.

## Musical words

* `SIN` ( freq -- pcm ) : Sine wave oscillator

    _example_ `440 HZ SIN`

* `SAW` ( freq -- pcm ) : Sawtooth oscillator

* `PULSE` ( freq width -- pcm ) : Pulse wave oscillator

* `SQ` ( freq -- pcm ) === `0.5 PULSE` : Square wave oscillator

* `TR` ( freq -- pcm ) : Triangle wave oscillator* 

* `CLIP` ( n ) : When outputting to a buffer (sound hardware...), scale the results down by this factor.

    _example_ `A SIN C SIN 2 CLIP` will scale down the 2 notes so they fit into -1...+1

