package main

import (
    "fmt"
    "jsonselect"
)

func main() {
    results, err := jsonselect.Lex(
        ".biscuits:has(magic) string:nth-last-child(even)",
        jsonselect.SCANNER,
    )
    if err != nil {
        fmt.Println(err)
    } else {
        for _, result := range results {
            fmt.Println(*result)
        }
    }
}
