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

    // Not actually a type, obviously
    J_OPER jsonType = "oper"
)

type JsonNode struct {
    value interface{}
    typ jsonType
    json *simplejson.Json
    parent *JsonNode
    parent_key string
    idx int
    siblings int
}

func (p *parser) getFlooredDocumentMap(node *JsonNode) []*JsonNode {
    var newMap []*JsonNode
    return p.findSubordinateJsonNodes(node.json, newMap, nil, "", -1, -1)
}

func (p *parser) findSubordinateJsonNodes(jdoc *simplejson.Json, nodes []*JsonNode, parent *JsonNode, parent_key string, idx int, siblings int) []*JsonNode {
    node := JsonNode{}
    node.parent = parent
    node.json = jdoc
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
            nodes = p.findSubordinateJsonNodes(element, nodes, &node, "", i + 1, length)
        }
    }
    data, err := jdoc.Map()
    if err == nil {
        node.value, _ = jdoc.Map()
        node.typ = J_OBJECT
        for key := range data {
            element := jdoc.Get(key)
            nodes = p.findSubordinateJsonNodes(element, nodes, &node, key, -1, -1)
        }
    }

    nodes = append(nodes, &node)
    return nodes
}

func (p *parser) mapDocument() {
    var nodes []*JsonNode
    p.nodes = p.findSubordinateJsonNodes(p.data, nodes, nil, "", -1, -1)
}
