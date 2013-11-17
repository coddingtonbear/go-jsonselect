
go-jsonselect
=============

A golang implementation of [JSONSelect](http://jsonselect.org/) modeled off of [@mwhooker's Python implementation](https://github.com/mwhooker/jsonselect)


Usage
-----

```golang

import (
    "jsonselect"
)

var json string = `
    {
        "beers": [
            {
                "title": "alpha",
                "rating": 50
            },
            {
                "title": "beta",
                "rating": 90
            }
        ]
    }
`

parser, _ := jsonselect.CreateParserFromString(json)
results, _ := parser.GetValues(".beers object:has(.rating:expr(x>70)))
// Results [map[title: beta rating: 90]]
```
