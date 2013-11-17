package main

import (
    "fmt"
    "jsonselect"
)

var json string = `
    {
        "stories": [
            {
                "title": "alpha",
                "good": false,
                "rating": 45
            },
            {
                "title": "beta",
                "good": true,
                "rating": 90
            }
        ]
    }
`

func main() {
    parser, err := jsonselect.CreateParserFromString(json)
    if err != nil {
        fmt.Println(err)
        return
    }
    results, err := parser.Parse(".stories object:has(.rating:expr(x>50))")
    if err != nil {
        fmt.Println(err)
    } else {
        if len(results) > 0 {
            for idx, result := range results {
                fmt.Println(idx, *result)
            }
        } else {
            fmt.Println("No matches")
        }
    }
}
