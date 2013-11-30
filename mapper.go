package jsonselect

import (
    "sort"
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

type jsonNode struct {
    value interface{}
    typ jsonType
    json *simplejson.Json
    parent *jsonNode
    parent_key string
    idx int
    siblings int
    position int
}

func buildAncestorNodeMap(node *jsonNode, ancestors map[*simplejson.Json]bool) {
    if node.parent != nil {
        ancestors[node.parent.json] = true
        buildAncestorNodeMap(node.parent, ancestors)
    }
}

func (p *Parser) getFlooredDocumentMap(node *jsonNode) map[*simplejson.Json]*jsonNode {
    ancestorNodeMap := make(map[*simplejson.Json]bool)
    buildAncestorNodeMap(node, ancestorNodeMap)

    flooredMap := make(map[*simplejson.Json]*jsonNode)
    for key, value := range p.nodes {
        _, ok := ancestorNodeMap[key]
        if !ok {
            flooredMap[key] = value
        }
    }

    return flooredMap
}

func (p *Parser) buildJsonNodeMap(jdoc *simplejson.Json, nodes map[*simplejson.Json]*jsonNode, parent *jsonNode, parent_key string, idx int, siblings int, position *int) map[*simplejson.Json]*jsonNode {
    node := jsonNode{}
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
            p.buildJsonNodeMap(element, nodes, &node, "", i + 1, length, position)
        }
    }
    data, err := jdoc.Map()
    if err == nil {
        node.value, _ = jdoc.Map()
        node.typ = J_OBJECT
        for key := range data {
            element := jdoc.Get(key)
            p.buildJsonNodeMap(element, nodes, &node, key, -1, -1, position)
        }
    }

    *position++
    node.position = *position
    if node.json != nil {
        nodes[node.json] = &node
    }
    return nodes
}

func (p *Parser) mapDocument() {
    p.nodes = make(map[*simplejson.Json]*jsonNode)
    var position int = 0
    p.buildJsonNodeMap(p.Data, p.nodes, nil, "", -1, -1, &position)
}

type nodeSortFunction func(n1, n2 *jsonNode) bool

type nodeSorter struct {
    nodes []*jsonNode
    function func(n1, n2 *jsonNode) bool
}

func (s nodeSortFunction) Sort(nodes []*jsonNode) {
    sorter := &nodeSorter{
        nodes: nodes,
        function: s,
    }
    sort.Sort(sorter)
}

func (s *nodeSorter) Len() int {
    return len(s.nodes)
}

func (s *nodeSorter) Swap(i, j int) {
    s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i]
}

func (s *nodeSorter) Less(i, j int) bool {
    return s.function(s.nodes[i], s.nodes[j])
}
