package jsonselect

import (
    "github.com/latestrevision/go-simplejson"
)

type Parser struct {
    data *simplejson.Json
}

type Node struct {
    value interface{}
    parent *Node
    parent_key string
    idx int
    siblings int
}

func CreateParser(body []byte) (*Parser, error) {
    json, err := simplejson.NewJson(body)
    if err != nil {
        return nil, err
    }
    return &Parser{json}, nil
}
