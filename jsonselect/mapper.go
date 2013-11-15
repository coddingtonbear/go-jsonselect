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
    node.parent = &parent
    if len(parent_key) > 0 {
        node.parent_key = parent_key
    }
    if idx > -1 {
        node.idx = idx
    }
    if siblings > -1 {
        node.siblings = siblings
    }

    length, err := jdoc.ArrayLength()
    if err == nil {
        node.value, _ = jdoc.Array()
        node.typ = J_ARRAY
        for i := 0; i < length; i++ {
            element := jdoc.GetIndex(i)
            value := strconv.FormatInt(int64(i), 10)
            new_nodes := findSubordianteNodes(element, nodes, jdoc, "", i + 1, len(node.value))
            for _, node := range new_nodes {
                nodes = append(nodes, node)
            }
        }
    }
    data, err := jdoc.Map()
    if err == nil {
        node.value, _ = jdoc.Map()
        node.typ = J_MAP
        for key := range data {
            element := jdoc.Get(key)
            new_nodes := findSubordianteNodes(element, nodes, jdoc, key, -1, -1)
            for _, node := range new_nodes {
                nodes = append(nodes, node)
            }
        }
    }

    nodes = append(nodes, node)
    return nodes
}

func (p *Parser) mapDocument() ([]*Node, error) {
    var nodes []*Node
    nodes = p.findSubordinateNodes(p.data, nodes, nil, "", -1, -1)
    return nodes, nil
}
