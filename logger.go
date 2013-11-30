package jsonselect

import (
    "fmt"
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

func getFormattedNodeArray(nodes []*jsonNode) []string {
    var formatted []string
    for _, node := range nodes {
        if node != nil {
            formatted = append(formatted, fmt.Sprint(*node))
        } else {
            formatted = append(formatted, fmt.Sprint(nil))
        }
    }
    return formatted
}

func getFormattedTokens(tokens []*token) []string {
    var output []string
    for _, token := range tokens {
        output = append(output, fmt.Sprint(token.val))
    }
    return output
}

func getFormattedExpression(tokens []*exprElement) []string {
    var output []string
    for _, token := range tokens {
        output = append(output, fmt.Sprint(token.value))
    }
    return output
}
