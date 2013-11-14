package jsonselect

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
