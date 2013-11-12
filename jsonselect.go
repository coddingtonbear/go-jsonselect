package jsonselect

import (
    "github.com/latestrevision/go-simplejson"
)

type Parser struct {
    data *simplejson.Json
}

func CreateParser(body []byte) (*Parser, error) {
    json, err := simplejson.NewJson(body)
    if err != nil {
        return nil, err
    }
    return &Parser{json}, nil
}
