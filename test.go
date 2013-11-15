package main

import (
    "fmt"
    "jsonselect"
)

var json string = `
    {
        "user": false,
        "name": {
            "first": "Adam",
            "last": "Coddington",
            "details": {
                "address": {
                    "city": null
                }
            },
            "age": 30
        },
        "stories": [
            {
                "title": "alpha",
                "good": false
            },
            {
                "title": "beta",
                "good": true
            }
        ],
        "things": [
            {
                "type": "bike",
                "title": "LHT"
            }
        ]
    }
`

func main() {
    parser, err := jsonselect.CreateParser(json)
    if err != nil {
        fmt.Println(err)
        return
    }
    results, err := parser.Parse(".stories :nth-child(1)")
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
