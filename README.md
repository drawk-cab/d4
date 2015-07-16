# d4
Diminished Forth definition and interpreter

This is the tiny language [floatbeat](http://github.com/colourcountry/floatbeat) uses to define sounds, designed for 
fun and fast live coding of music.

It's a stack based language similar to Forth.

## Principles

### Words

A word is a contiguous sequence of (unicode.IsDigit or .) or (unicode.IsLetter or _) or a single other non-whitespace rune.

Whitespace is only needed in order to separate adjacent sequences of digits or letters.

    2DUP --> [2] [DUP]
    2.4* --> [2.4] [*]
    a__b!=5 --> [a__b] [!] [=] [5]

Words are case insensitive.

### There is only one data type, number (float64)

### Programs are run iteratively.

A typical run will consist of many interations of the supplied code.

The iteration number is available in the built-in word `T`.

Built-in words `S`, `BPM`, `HZ` convert `T` to time units based
on a supplied sample rate (default 22050)

Use `KEEP` to define a value to be used later: `75 KEEP my_var` Retrieve with `my_var?`

`KEEP`d values are evaluated at runtime but are fixed once set.
(`KEEP` is less confusing than `CONSTANT`, as values may vary from one iteration to the next)

`KEEP`d values can be accessed in later iterations using `OLD`.

You can build filters using `DELTA` and `PREWARP`.

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
* `CONSTANT`
* `!` `@` `?`

## Built-in units

* `HZ` (freq -- counter) : Create a counter which increases by `freq` every second

* `S` (time -- counter) : Create a counter which increases by 1 every `time` seconds

* `BPM` (freq -- counter) : Create a counter which increases by `freq` every minute

## Extra words

* `.` : remove the item on top of stack (TOS) and add it to the output stack

* `&` === `DUP .` *i.e.* add the item on TOS to the output stack without removing it from the stack

* `NOOP` : noop

* `T` : iteration number

* `ON` ( schedule_t, duration, base_t -- 0 or age, 1 ) : Is the note of length `duration` scheduled for `schedule_t` currently in progress at time `base_t`, and if so, how old is it?

    _example_ `1S 0.5S T ON IF 440HZ SIN THEN` plays an A for 0.5 sec, 1 sec after the start.

* `FROM`...`CHOOSE` : treat TOS as a pointer and execute the word with that index between `FROM` and `CHOOSE` . So `2 FROM a b c d CHOOSE` would execute only `c` because `c` has index 2 in the set.

    _example_ `T 1S 4 DMOD FROM play_c play_d play_e play_f CHOOSE` will execute `play_c` in the first second, `play_d` in the second, and so on. Because of the `DMOD`, the age of each note will be on TOS for the play routines.

* `KEEP x` === `CONSTANT x !`

* `OLD` ( definition_name time -- value ) : get the value saved under `definition_name` in the iteration the specified time ago

    _example_ `my_var 0.3s OLD`

* `DELTA` ( definition_name -- delta ) : get the difference between the current value of `definition_name` and its value in a previous iteration. Useful for implementing filters.

* `PREWARP` ( freq -- prewarp_freq ) : Evaluate **tan**( π `freq` / *samplerate* ) , useful for implementing filters.


* `::`...`;` : import pre-set definition packages

    _example_ `:: scale divisions;` will import the `scale` package (defines the equal tempered scale)
    and the `divisions` package

## Musical words

* `#`, `SHARP`, `♯` ( freq -- freq ) : Sharpen a frequency by 1 semitone (equal tempered)

* `FLAT`, `♭` ( freq -- freq ) : Flatten a frequency by 1 semitone (equal tempered)

* `SIN` ( counter -- pcm ) : Sine wave oscillator

    _example_ `440 HZ SIN`

* `SAW` ( freq -- pcm ) : Sawtooth oscillator

* `PULSE` ( freq width -- pcm ) : Pulse wave oscillator

* `SQ` ( freq -- pcm ) === `0.5 PULSE` : Square wave oscillator

* `TR` ( freq -- pcm ) : Triangle wave oscillator* 

* `CLIP` ( n ) : When outputting to a buffer (sound hardware...), scale the results down by this factor.

    _example_ `A SIN C SIN 2 CLIP` will scale down the 2 notes so they fit into -1...+1

