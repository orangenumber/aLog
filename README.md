# alog

(c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>  
Last update: 10/18/2020


## Intro

aLog is the simple logger with a goal of zero allocation.
Please note that Alog's `Printf` does not support same formatting of `fmt.Printf`. _(see usage)_


---

## Usage

- `Printf(format string, s ...interface{})`
    - To make zero allocation, a simpler version of Printf has been created.
    - This only supports following: `%d`, `%f`, `%s`, and `%t`.
        - For a float (`%f`), currently two decimal points are supported.
- `Printj(optionalPrefix string, a interface{})`: This is used print JSON for a struct. This takes optional string prefix.
- `Print(s ...interface{})`
- `SetOutput(output io.Writer)`
- `SetPrefix(prefix string)`
- `SetFlag(flag uint16)`
    - Available flags:
        - `F_TIME`: Print time (`14:01:02`)
        - `F_MMDD`: Print date (`01/02`)
        - `F_MICROSEC`: Print microsecond (`10:01:02.000`)
        - `F_PREFIX`: Print prefix if any
        - `F_UTC`: Use UTC time
        - `F_DATE`: Print date (`2020/01/02` format, including year)
        - `F_USE_BUF_1K`: Use buffer of 1K.
        - `F_USE_BUF_2K`: Use buffer of 2K.
        - `F_STD`: Use `F_MMDD`, `F_TIME`, and `F_PREFIX`

__Note:__ for a higher performance, use `if` condition in front of the log call,
    by doing so, potentially many string operation won't need to run and save
    potential memory allocation.

```go
var debug = false
var info = true
var l = alog.New(os.Stdout, "", alog.STD)
```

```go
if debug {l.Print("[Debug] " + someValue + " / " + someValue2 )}
```

### Creating an alog instance

Alog can be created by `New(output io.Writer, prefix string, flag uint16) *ALogger`.

```go
l := alog.New(os.Stdout, "[test] ", alog.STD)
l.Printf("Log says: %s", "Hello")
```


### Without creating an alog instance

Alog also can be used without creating an object.
`os.Stdout` will be used for a default output. However, this can be customized as well as flag.

```go
package main

import (
    "github.com/gonyyi/alog"
    "os"
)

func main() {
    alog.Print("Hello") // standard output
    alog.SetOutput(alog.Discard) // change output to Discard.
    alog.Print("Discarded Hello")
    alog.SetOutput(os.Stderr) // use standard error
    alog.Print("Stderr Hello")

    alog.SetFlag(alog.F_STD) // set standard (MMDD, TIME, PREFIX)
    alog.SetPrefix("[test] ") // set prefix
}
```


---

## Example

__With new object__

```go
package main

import (
    "github.com/gonyyi/alog"
    "os"
)

func main() {
	out, _ := os.Create("test.log")
	x := alog.New(out, "test ", alog.O_PREFIX|alog.O_TIME|alog.O_DATE)
    x.Printf("Name: %s\nAge: %d", "Gon", 17)
    a := struct {
        Name  string `json:"name"`
        City  string `json:"city"`
        Count int    `json:"cnt"`
    }{
        Name: "Gon",
        City: "Conway",
    }
    x.Printj("log|", &a)
    x.Close()
}
```
