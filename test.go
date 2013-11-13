package main

import (
    "fmt"
    "jsonselect"
)

func main() {
    result, err := jsonselect.Lex(
        "test",
        jsonselect.EXPRESSION_SCANNER,
    )
    fmt.Println(result)
    fmt.Println(err)
}
