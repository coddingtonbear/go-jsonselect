package jsonselect

import (
    "log"
    "strconv"
)

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
        return float64(as_int)
    }
    as_string, ok := in.(string)
    if ok {
        parsed_float_string, err := strconv.ParseFloat(as_string, 64)
        if err == nil {
            return parsed_float_string
        }
        parsed_int_string, err := strconv.ParseInt(as_string, 10, 32)
        if err == nil {
            return float64(parsed_int_string)
        }
    }
    log.Print("Error transforming ", in, " into Float64")
    return -1
}

func getInt32(in interface{}) int32 {
    value := int32(getFloat64(in))
    if value == -1 {
        log.Print("Errot ransforming ", in, " into Int32")
    }
    return value
}
