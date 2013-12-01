package jsonselect

import (
    "encoding/json"
    "log"
    "strconv"
    "github.com/latestrevision/go-simplejson"
)

func nodeIsMemberOfList(needle *jsonNode, haystack map[*simplejson.Json]*jsonNode) bool {
    _, ok := haystack[needle.json]
    return ok
}

func nodeIsAncestorOfHaystackMember(needle *jsonNode, haystack map[*simplejson.Json]*jsonNode) bool {
    if nodeIsMemberOfList(needle, haystack) {
        return true
    }
    if needle.parent == nil {
        return false
    }
    return nodeIsAncestorOfHaystackMember(needle.parent, haystack)
}

func parents(lhs map[*simplejson.Json]*jsonNode, rhs map[*simplejson.Json]*jsonNode) map[*simplejson.Json]*jsonNode {
    results := make(map[*simplejson.Json]*jsonNode)

    for _, element := range rhs {
        if nodeIsMemberOfList(element.parent, lhs) {
            results[element.json] = element
        }
    }

    logger.Print(len(results), " of [", len(rhs), "]RHS elements are parents of [", len(lhs), "]LHS")
    return results
}

func ancestors(lhs map[*simplejson.Json]*jsonNode, rhs map[*simplejson.Json]*jsonNode) map[*simplejson.Json]*jsonNode {
    results := make(map[*simplejson.Json]*jsonNode)

    for _, element := range rhs {
        if nodeIsAncestorOfHaystackMember(element, lhs) {
            results[element.json] = element
        }
    }

    logger.Print(len(results), " of [", len(rhs), "]RHS elements are ancestors of [", len(lhs), "]LHS")
    return results
}

func siblings(lhs map[*simplejson.Json]*jsonNode, rhs map[*simplejson.Json]*jsonNode) map[*simplejson.Json]*jsonNode {
    parents := make(map[*simplejson.Json]*jsonNode)
    results := make(map[*simplejson.Json]*jsonNode)

    for _, element := range lhs {
        parents[element.parent.json] = element.parent
    }

    for _, element := range rhs {
        if nodeIsMemberOfList(element.parent, parents){
            results[element.json] = element
        }
    }

    logger.Print(len(results), " of [", len(rhs), "]RHS elements are siblings of [", len(lhs), "]LHS")
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
    panic("ERR")
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
    as_string, ok := in.(string)
    if ok {
        return as_string
    }
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

func nodeMapToArray(nodes map[*simplejson.Json]*jsonNode) []*jsonNode {
    output := make([]*jsonNode, 0, len(nodes))
    for _, value := range nodes{
        output = append(
            output,
            value,
        )
    }
    return output
}
