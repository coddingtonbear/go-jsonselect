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
    }
`

func main() {
    parser, err := jsonselect.CreateParser(json)
    if err != nil {
        fmt.Println(err)
    }
    results, err := parser.Parse(".bin")
    if err != nil {
        fmt.Println(err)
    } else {
        for _, result := range results {
            fmt.Println(*result)
        }
    }
}
