package jsonselect

import (
    "github.com/latestrevision/go-simplejson"
)

type jsonType string

const (
    J_STRING jsonType = "string"
    J_NUMBER jsonType = "number"
    J_OBJECT jsonType = "object"
    J_ARRAY jsonType = "array"
    J_BOOLEAN jsonType = "boolean"
    J_NULL jsonType = "null"
)

type Node struct {
    value interface{}
    typ jsonType
    parent *Node
    parent_key string
    idx int
    siblings int
}

func (p *Parser) findSubordinateNodes(jdoc *simplejson.Json, nodes []*Node, parent *Node, parent_key string, idx int, siblings int) []*Node {
    node := Node{}
    node.parent = parent
    if len(parent_key) > 0 {
        node.parent_key = parent_key
    }
    if idx > -1 {
        node.idx = idx
    }
    if siblings > -1 {
        node.siblings = siblings
    }

    string_value, err := jdoc.String()
    if err == nil {
        node.value = string_value
        node.typ = J_STRING
    }

    int_value, err := jdoc.Int()
    if err == nil {
        node.value = int_value
        node.typ = J_NUMBER
    }

    float_value, err := jdoc.Float64()
    if err == nil {
        node.value = float_value
        node.typ = J_NUMBER
    }

    bool_value, err := jdoc.Bool()
    if err == nil {
        node.value = bool_value
        node.typ = J_BOOLEAN
    }

    if jdoc.IsNil() {
        node.value = nil
        node.typ = J_NULL
    }

    length, err := jdoc.ArrayLength()
    if err == nil {
        node.value, _ = jdoc.Array()
        node.typ = J_ARRAY
        for i := 0; i < length; i++ {
            element := jdoc.GetIndex(i)
            nodes = p.findSubordinateNodes(element, nodes, &node, "", i + 1, length)
        }
    }
    data, err := jdoc.Map()
    if err == nil {
        node.value, _ = jdoc.Map()
        node.typ = J_OBJECT
        for key := range data {
            element := jdoc.Get(key)
            nodes = p.findSubordinateNodes(element, nodes, &node, key, -1, -1)
        }
    }

    nodes = append(nodes, &node)
    return nodes
}

func (p *Parser) mapDocument() []*Node {
    var nodes []*Node
    nodes = p.findSubordinateNodes(p.data, nodes, nil, "", -1, -1)
    return nodes
}
