package jsonselect

import (
    "encoding/json"
    "log"
    "strconv"
)

type exprElement struct {
    value interface{}
    typ jsonType
}

func nodeIsMemberOfList(needle *Node, haystack []*Node) bool {
    for _, element := range haystack {
        if element == needle {
            return true
        }
    }
    return false
}

func nodeIsAncestorOfHaystackMember(needle *Node, haystack []*Node) bool {
    if nodeIsMemberOfList(needle, haystack) {
        return true
    }
    if needle.parent == nil {
        return false
    }
    return nodeIsAncestorOfHaystackMember(needle.parent, haystack)
}

func parents(lhs []*Node, rhs []*Node) []*Node {
    var results []*Node

    for _, element := range rhs {
        if nodeIsMemberOfList(element.parent, lhs) {
            results = append(results, element)
        }
    }

    return results
}

func ancestors(lhs []*Node, rhs []*Node) []*Node {
    var results []*Node

    for _, element := range rhs {
        if nodeIsAncestorOfHaystackMember(element, lhs) {
            results = append(results, element)
        }
    }

    return results
}

func siblings(lhs []*Node, rhs []*Node) []*Node {
    var parents []*Node
    var results []*Node

    for _, element := range lhs {
        parents = append(parents, element.parent)
    }

    for _, element := range rhs {
        if nodeIsMemberOfList(element.parent, parents){
            results = append(results, element)
        }
    }

    return results
}

func getFloat64(in interface{}) float64 {
    as_float, ok := in.(float64)
    if ok {
        return as_float
    }
    as_int, ok := in.(int64)
    if ok {
        value := float64(as_int)
        return value
    }
    as_string, ok := in.(string)
    if ok {
        parsed_float_string, err := strconv.ParseFloat(as_string, 64)
        if err == nil {
            value := parsed_float_string
            return value
        }
        parsed_int_string, err := strconv.ParseInt(as_string, 10, 32)
        if err == nil {
            value := float64(parsed_int_string)
            return value
        }
    }
    result := float64(-1)
    log.Print("Error transforming ", in, " into Float64")
    return result
}

func getInt32(in interface{}) int32 {
    value := int32(getFloat64(in))
    if value == -1 {
        log.Print("Error transforming ", in, " into Int32")
    }
    return value
}

func getJsonString(in interface{}) string {
    marshaled_result, err := json.Marshal(in)
    if err != nil {
        log.Print("Error transforming ", in, " into JSON string")
    }
    result := string(marshaled_result)
    return result
}

func exprElementIsTruthy(e exprElement) bool {
    switch e.typ {
        case J_STRING:
            return len(e.value.(string)) > 0
        case J_NUMBER:
            return e.value.(float64) > 0
        case J_OBJECT:
            return true
        case J_ARRAY:
            return true
        case J_BOOLEAN:
            return e.value.(bool)
        case J_NULL:
            return false
        default:
            return false
    }
}

func exprElementsMatch(lhs exprElement, rhs exprElement) bool {
    return lhs.typ == rhs.typ
}
