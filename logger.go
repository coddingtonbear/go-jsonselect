package jsonselect

import (
    "log"
    "os"
)

type Logger struct {
    Enabled bool
}


var logger = Logger{false}
var handler = log.New(os.Stderr, "jsonselect: ", 0)

func (*Logger) Print(a ...interface{}) {
    if logger.Enabled {
        handler.Print(a...)
    }
}

func (*Logger) Println(a ...interface{}) {
    if logger.Enabled {
        handler.Println(a...)
    }
}

func EnableLogger() {
    logger.Enabled = true
}
