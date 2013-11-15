package jsonselect

import (
    "errors"
    "regexp"
    "strconv"
    "github.com/latestrevision/go-simplejson"
)

type Parser struct {
    data *simplejson.Json
    nodes []*Node
}

func CreateParser(body string) (*Parser, error) {
    json, err := simplejson.NewJson([]byte(body))
    if err != nil {
        return nil, err
    }
    parser := Parser{json, nil}
    parser.mapDocument()
    return &parser, err
}

func (p *Parser) Parse(selector string) ([]*Node, error) {
    tokens, err := Lex(selector, SCANNER)
    if err != nil {
        return nil, err
    }

    results, err := p.selectorProduction(tokens)
    if err != nil {
        return nil, err
    }

    return results, nil
}

func (p *Parser) selectorProduction(tokens []*token) ([]*Node, error) {
    var results []*Node
    var validators []func(*Node)bool
    var matched bool
    var value interface{}
    var result func()

    _, matched, _ = p.peek(tokens, S_TYPE)
    if matched {
        value, _ = p.match(tokens, S_TYPE)
        validators = append(
            validators,
            p.typeProduction(value),
        )
    }
    _, matched, _ = p.peek(tokens, S_IDENTIFIER)
    if matched {
        value, _ = p.match(tokens, S_IDENTIFIER)
        validators = append(
            validators,
            p.keyProduction(value),
        )
    }
    _, matched, _ = p.peek(tokens, S_PCLASS)
    if matched {
        value, _ = p.match(tokens, S_PCLASS)
        validators = append(
            validators,
            p.pclassProduction(value),
        )
    }
    _, matched, _ = p.peek(tokens, S_NTH_FUNC)
    if matched {
        value, _ = p.match(tokens, S_NTH_FUNC)
        validators = append(
            validators,
            p.nthChildProduction(value, tokens),
        )
    }
    _, matched, _ = p.peek(tokens, S_PCLASS_FUNC)
    if matched {
        value, _ = p.match(tokens, S_PCLASS_FUNC)
        validators = append(
            validators,
            p.pclassFuncProduction(value, tokens),
        )
    }

    if len(validators) < 1 {
        return nil, errors.New("No selector recognized")
    }

    results, err := p.matchNodes(validators)
    if err != nil {
        return nil, err
    }

    _, matched, _ = p.peek(tokens, S_OPER)
    if matched {
        value, _ = p.match(tokens, S_OPER)
        rvals, err := p.selectorProduction(tokens)
        if err != nil {
            return nil, err
        }
        switch value {
            case ",":
                for _, val := range rvals {
                    // TODO: This is quite slow
                    // it seems like it's probably quite easy to expand
                    // the list just once
                    results = append(results, val)
                }
            case ">":
                results = parents(results, rvals)
            case "~":
                results = siblings(results, rvals)
            case " ":
                results = ancestors(results, rvals)
            default:
                return nil, errors.New("Unrecognized operator")
        }
    } else if len(tokens) > 0 {
        rvals, err := p.selectorProduction(tokens)
        if err != nil {
            return nil, err
        }
        results = ancestors(results, rvals)
    }

    return results, nil
}

func (p *Parser) peek(tokens []*token, typ tokenType) (interface{}, bool, error) {
    if len(tokens) < 1 {
        return nil, false, errors.New("No more tokens")
    }
    if tokens[0].typ == typ {
        return tokens[0].val, true, nil
    }
    return nil, false, nil
}

func (p *Parser) match(tokens []*token, typ tokenType) (interface{}, error) {
    value, matched, error := p.peek(tokens, typ)
    if !matched {
        return nil, errors.New("Match not successful")
    }
    element, tokens := tokens[0], tokens[1:]
    return value, nil
}

func (p *Parser) matchNodes(validators []func(*Node)bool) ([]*Node, error) {
    // TODO:
    var matches []*Node
    return matches, nil
}

func (p *Parser) typeProduction(value interface{}) func(*Node)bool {
    return func(node *Node) bool {
        return node.typ == value.(jsonType)
    }
}
func (p *Parser) keyProduction(value interface{}) func(*Node)bool {
    // TODO: Verify this -- I'm not sure what this is supposed to match
    return func(node *Node) bool {
        if node.parent_key == ""{
            return false
        }
        return node.parent_key == value.(string)
    }
}

func (p *Parser) pclassProduction(value interface{}) func(*Node)bool {
    pclass := value.(string)
    if pclass == "first-child" {
        return func(node *Node) bool {
            return node.idx == 1
        }
    } else if pclass == "last-child" {
        return func(node *Node) bool {
            return node.siblings > 0 && node.idx == node.siblings
        }
    } else if pclass == "only-child" {
        return func(node *Node) bool {
            return node.siblings == 1
        }
    } else if pclass == "root" {
        return func(node *Node) bool {
            return node.parent == nil
        }
    } else if pclass == "empty" {
        return func(node *Node) bool {
            return node.typ == J_ARRAY && len(node.value.(string)) < 1
        }
    }
    return func(node *Node) bool {
        return false
    }
}

func (p *Parser) nthChildProduction(value interface{}, tokens []*token) func(*Node)bool {
    nthChildRegexp := regexp.MustCompile(`^\s*\(\s*(?:([+\-]?)([0-9]*)n\s*(?:([+\-])\s*([0-9]))?|(odd|even)|([+\-]?[0-9]+))\s*\)`)
    args, _ := p.match(tokens, S_EXPR)
    var a int
    var b int
    var reverse bool = false

    pattern := nthChildRegexp.FindStringSubmatch(args.(string))

    if pattern[5] != "" {
        a = 2
        if pattern[5] == "odd" {
            b = 1
        } else {
            b = 0
        }
    } else if pattern[6] != ""{
        a = 0
        b, _ := strconv.ParseInt(pattern[6], 10, 64)
    } else {
        sign := "+"
        if pattern[1] != "" {
            sign = pattern[1]
        }
        coeff := "1"
        if pattern[2] != "" {
            coeff = pattern[2]
        }
        a, _ := strconv.ParseInt(coeff, 10, 64)
        if sign == "-" {
            a = -1 * a
        }
        g3, _ := strconv.ParseInt(pattern[3], 10, 64)
        g4, _ := strconv.ParseInt(pattern[4], 10, 64)
        if pattern[3] != "" {
            b = int(g3 + g4)
        } else {
            b = 0
        }
    }

    if value.(string) == "nth-last-child" {
        reverse = true
    }

    return func(node *Node)bool {
        if node.siblings == 0 {
            return false
        }

        idx := node.idx - 1
        if reverse {
            idx = node.siblings - idx
        } else {
            idx++
        }

        if a == 0 {
            return b == idx
        } else {
            return ((idx - b) % a) == 0 && (idx * a + b) >= 0
        }
    }
}

func (p *Parser) pclassFuncProduction(value interface{}, tokens []*token) func(*Node)bool {
    args, _ := p.match(tokens, S_EXPR)
    pclass = value.(string)

    if pclass == "expr" {
        tokens, _ := Lex(args.(string), EXPRESSION_SCANNER)
        return func(node *Node)bool {
            return p.parseExpression(tokens, node).(bool)
        }
    }

    args, _ = Lex(args.(string)[1:len(args.(string))-1], SCANNER)

    if pclass == "has" {
        for _, token := range args {
            if token.value == ">" {
                token.typ = S_EMPTY
            }
        }
        rvals := p.selectorProduction(args)

        var ancestors []*Node
        for _, node := range rvals {
            ancestors = append(ancestors, node.parent)
        }
        return func(node *Node)bool {
            return ndoeIsMemberOfList(node, ancestors)
        }
    } else if pclass == "contains" {
        return func(node *Node)bool {
            return node.typ == J_STRING && strings.Count(node.value.(string), args[0].val.(string)) > 0
        }
    } else if pclass == "val" {
        return func(node *Node)bool {
            return node.typ == J_STRING && node.value.(string) == args[0].val.(string)
        }
    }

    // If we didn't find a known pclass, do not match anything.
    return func(node *Node)bool {
        return false
    }
}

func (p *Parser) evaluateParsedExpression(tokens []*token, node *Node, cmap map[string]func(*token, *token)int) interface{} {
    var matched bool
    if len(tokens) < 1 {
        return false
    }

    value, matched, _ = p.peek(tokens, S_PAREN)
    if matched && value.(string) == "(" {
        p.match(tokens, S_PAREN)
        lhs := p.evaluateParsedExpression(tokens, node, cmap)
        return lhs
    }

    _, matched, _ = p.peek(tokens, S_PVAR)
    if matched {
        p.match(tokens, S_PVAR)
        lhs = node.value
    } else {
        relevantTokens := []tokenType{S_STRING, S_BOOL, S_NIL, S_NUMBER}
        for _, ttype := range relevantTokens {
            _, matched, _ = p.peek(tokens, ttype)
            if matched {
                lhs := p.match(tokens, ttype)
                break
            }
        }
    }

    value, matched, _ = p.peek(tokens, S_PAREN)
    if matched && value.(string) == ")" {
        p.match(tokens, S_PAREN)
        return lhs
    }

    binop, _ := p.match(tokens, S_BINOP)
    comparatorFunction = cmap[binop.(string)]
    rhs := p.evaluateParsedExpression(tokens, node, cmap)

    return comparatorFunction(lhs, rhs)
}

func (p *Parser) parseExpression(tokens []*token, node *Node) interface{} {
    comparatorMap = map[string]func(lhs *token, rhs *token)int {
        "*": func(lhs *token, rhs *token)int {
            return lhs.val.(int) * rhs.val.(int)
        },
        "/": func(lhs *token, rhs *token)int {
            return lhs.val.(int) / rhs.val.(int)
        },
        "%": func(lhs *token, rhs *token)int {
            return lhs.val.(int) % rhs.val.(int)
        },
        "+": func(lhs *token, rhs *token)int {
            return lhs.val.(int) + rhs.val.(int)
        },
        "-": func(lhs *token, rhs *token)int {
            return lhs.val.(int) - rhs.val.(int)
        },
        "<=": func(lhs *token, rhs *token)int {
            return lhs.val.(int) <= rhs.val.(int)
        },
        "<": func(lhs *token, rhs *token)int {
            return lhs.val.(int) < rhs.val.(int)
        },
        ">=": func(lhs *token, rhs *token)int {
            return lhs.val.(int) > rhs.val.(int)
        },
        ">": func(lhs *token, rhs *token)int {
            return lhs.val.(int) > rhs.val.(int)
        },
        "$=": func(lhs *token, rhs *token)int {
            lhs_str = lhs.(string)
            rhs_str = rhs.(string)
            return strings.LastIndex(lhs_str, rhs_str) == len(lhs_str) - len(rhs_str)
        },
        "^=": func(lhs *token, rhs *token)int {
            lhs_str = lhs.(string)
            rhs_str = rhs.(string)
            return strings.Index(lhs_str, rhs_str) == 0
        },
        "*=": func(lhs *token, rhs *token)int {
            lhs_str = lhs.(string)
            rhs_str = rhs.(string)
            return strings.Index(lhs_str, rhs_str) != 0
        },
        "=": func(lhs *token, rhs *token)int {
            return lhs.(string) == rhs.(string)
        },
        "!=": func(lhs *token, rhs *token)int {
            return lhs.(string) != rhs.(string)
        },
        "&&": func(lhs *token, rhs *token)int {
            return lhs.(string) && rhs.(string)
        },
        "||": func(lhs *token, rhs *token)int {
            return lhs.(string) || rhs.(string)
        },
    }
    return p.evaluateParsedExpression(tokens, node, comparatorMap)
}
