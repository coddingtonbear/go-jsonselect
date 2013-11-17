package jsonselect

import (
    "errors"
    "io/ioutil"
    "log"
    "os"
    "regexp"
    "strconv"
    "strings"
    "github.com/latestrevision/go-simplejson"
)

var logger = log.New(ioutil.Discard, "jsonselect: ", 0,)

type Parser struct {
    data *simplejson.Json
    nodes []*Node
}

func EnableLogger() {
    logger = log.New(
        os.Stderr,
        "jsonselect: ",
        0,
    )
}

func CreateParserFromString(body string) (*Parser, error) {
    json, err := simplejson.NewJson([]byte(body))
    if err != nil {
        return nil, err
    }
    parser := Parser{json, nil}
    parser.mapDocument()
    return &parser, err
}

func CreateParser(json *simplejson.Json) (*Parser, error) {
    log.SetOutput(ioutil.Discard)
    parser := Parser{json, nil}
    parser.mapDocument()
    return &parser, nil
}

func (p *Parser) evaluateSelector(selector string) ([]*Node, error) {
    tokens, err := Lex(selector, SCANNER)
    if err != nil {
        return nil, err
    }

    nodes, err := p.selectorProduction(tokens, p.nodes)
    if err != nil {
        return nil, err
    }

    logger.Print(len(nodes), " matches found")

    return nodes, nil
}

func (p *Parser) GetJsonElements(selector string) ([]*simplejson.Json, error) {
    nodes, err := p.evaluateSelector(selector)
    if err != nil {
        return nil, err
    }

    var results []*simplejson.Json
    for _, node := range nodes {
        results = append(
            results,
            node.json,
        )
    }
    return results, nil
}

func (p *Parser) GetNodes(selector string) ([]*Node, error) {
    return p.evaluateSelector(selector)
}

func (p *Parser) GetValues(selector string) ([]interface{}, error) {
    nodes, err := p.evaluateSelector(selector)
    if err != nil {
        return nil, err
    }

    var results []interface{}
    for _, node := range nodes {
        results = append(
            results,
            node.value,
        )
    }

    return results, nil
}

func (p *Parser) selectorProduction(tokens []*token, documentMap []*Node) ([]*Node, error) {
    var results []*Node
    var validators []func(*Node)bool
    var matched bool
    var value interface{}
    var validator func(*Node)bool

    _, matched, _ = p.peek(tokens, S_TYPE)
    if matched {
        value, tokens, _ = p.match(tokens, S_TYPE)
        validators = append(
            validators,
            p.typeProduction(value),
        )
    }
    _, matched, _ = p.peek(tokens, S_IDENTIFIER)
    if matched {
        value, tokens, _ = p.match(tokens, S_IDENTIFIER)
        validators = append(
            validators,
            p.keyProduction(value),
        )
    }
    _, matched, _ = p.peek(tokens, S_PCLASS)
    if matched {
        value, tokens, _ = p.match(tokens, S_PCLASS)
        validators = append(
            validators,
            p.pclassProduction(value),
        )
    }
    _, matched, _ = p.peek(tokens, S_NTH_FUNC)
    if matched {
        value, tokens, _ = p.match(tokens, S_NTH_FUNC)
        validator, tokens = p.nthChildProduction(value, tokens)
        validators = append(validators, validator)
    }
    _, matched, _ = p.peek(tokens, S_PCLASS_FUNC)
    if matched {
        value, tokens, _ = p.match(tokens, S_PCLASS_FUNC)
        validator, tokens = p.pclassFuncProduction(value, tokens, documentMap)
        validators = append(validators, validator)
    }

    if len(validators) < 1 {
        return nil, errors.New("No selector recognized")
    }

    results, err := p.matchNodes(validators, documentMap)
    if err != nil {
        return nil, err
    }

    _, matched, _ = p.peek(tokens, S_OPER)
    if matched {
        value, tokens, _ = p.match(tokens, S_OPER)
        rvals, err := p.selectorProduction(tokens, documentMap)
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
        rvals, err := p.selectorProduction(tokens, documentMap)
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

func (p *Parser) match(tokens []*token, typ tokenType) (interface{}, []*token, error) {
    value, matched, _ := p.peek(tokens, typ)
    if !matched {
        return nil, tokens, errors.New("Match not successful")
    }
    _, tokens = tokens[0], tokens[1:]
    return value, tokens, nil
}

func (p *Parser) matchNodes(validators []func(*Node)bool, documentMap []*Node) ([]*Node, error) {
    var matches []*Node
    for _, node := range documentMap {
        var passed = true
        for _, validator := range validators {
            if !validator(node) {
                passed = false
                break
            }
        }
        if passed {
            matches = append(matches, node)
        }
    }
    return matches, nil
}

func (p *Parser) typeProduction(value interface{}) func(*Node)bool {
    logger.Print("Creating typeProduction validator ", value)
    return func(node *Node) bool {
        logger.Print("typeProduction ? ", node.typ, " == ", value)
        return string(node.typ) == value
    }
}
func (p *Parser) keyProduction(value interface{}) func(*Node)bool {
    // TODO: Verify this -- I'm not sure what this is supposed to match
    logger.Print("Creating keyProduction validator ", value)
    return func(node *Node) bool {
        logger.Print("keyProduction ? ", node.parent_key, " == ", value)
        if node.parent_key == ""{
            return false
        }
        return string(node.parent_key) == value
    }
}

func (p *Parser) pclassProduction(value interface{}) func(*Node)bool {
    pclass := value.(string)
    logger.Print("Creating pclassProduction validator ", pclass)
    if pclass == "first-child" {
        return func(node *Node) bool {
            logger.Print("pclassProduction first-child ? ", node.idx, " == 1")
            return node.idx == 1
        }
    } else if pclass == "last-child" {
        return func(node *Node) bool {
            logger.Print("pclassProduction last-child ? ", node.siblings, " > 0 AND ", node.idx, " == ", node.siblings)
            return node.siblings > 0 && node.idx == node.siblings
        }
    } else if pclass == "only-child" {
        return func(node *Node) bool {
            logger.Print("pclassProduction ony-child ? ", node.siblings, " == 1")
            return node.siblings == 1
        }
    } else if pclass == "root" {
        return func(node *Node) bool {
            logger.Print("pclassProduction root ? ", node.parent, " == nil")
            return node.parent == nil
        }
    } else if pclass == "empty" {
        return func(node *Node) bool {
            logger.Print("pclassProduction empty ? ", node.typ, " == ", J_ARRAY, " AND ", len(node.value.(string)), " < 1")
            return node.typ == J_ARRAY && len(node.value.(string)) < 1
        }
    }
    logger.Print("Error: Unknown pclass: ", pclass)
    return func(node *Node) bool {
        logger.Print("Asserting false due to failed pclassProduction")
        return false
    }
}

func (p *Parser) nthChildProduction(value interface{}, tokens []*token) (func(*Node)bool, []*token) {
    nthChildRegexp := regexp.MustCompile(`^\s*\(\s*(?:([+\-]?)([0-9]*)n\s*(?:([+\-])\s*([0-9]))?|(odd|even)|([+\-]?[0-9]+))\s*\)`)
    args, tokens, _ := p.match(tokens, S_EXPR)
    var a int
    var b int
    var reverse bool = false

    pattern := nthChildRegexp.FindStringSubmatch(args.(string))

    logger.Print("Creating nthChildProduction validator ", pattern)

    if pattern[5] != "" {
        a = 2
        if pattern[5] == "odd" {
            b = 1
        } else {
            b = 0
        }
    } else if pattern[6] != ""{
        a = 0
        b_temp, _ := strconv.ParseInt(pattern[6], 10, 64)
        b = int(b_temp)
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
        logger.Print("nthChildProduction ? ", node.siblings, " == 0")
        if node.siblings == 0 {
            return false
        }

        idx := node.idx - 1
        if reverse {
            idx = node.siblings - idx
        } else {
            idx++
        }

        logger.Print("nthChildProduction (continued) ? ", a, " == 0")
        if a == 0 {
            return b == idx
        } else {
            logger.Print("nthChildProduction (continued) ? ", idx-b % a, " == 0 AND ", idx*a+b, " >= 0")
            return ((idx - b) % a) == 0 && (idx * a + b) >= 0
        }
    }, tokens
}

func (p *Parser) pclassFuncProduction(value interface{}, tokens []*token, documentMap []*Node) (func(*Node)bool, []*token) {
    sargs, tokens, _ := p.match(tokens, S_EXPR)
    pclass := value.(string)

    logger.Print("Creating pclassFuncProduction validator ", pclass)

    if pclass == "expr" {
        tokens, _ := Lex(sargs.(string), EXPRESSION_SCANNER)
        var tokens_to_return []*token
        return func(node *Node)bool {
            result := p.parseExpression(tokens, node)
            logger.Print("pclassFuncProduction expr ? ", result, " > 0")
            return  result > 0
        }, tokens_to_return
    }

    lexString := sargs.(string)[1:len(sargs.(string))-1]
    args, _ := Lex(lexString, SCANNER)

    logger.Print("pclassFuncProduction lex results for [", lexString, "]: (follow)")
    for i, arg := range args {
        logger.Print("[", i, "]: ", arg)
    }

    if pclass == "has" {
        for _, token := range args {
            if token.val == ">" {
                token.typ = S_EMPTY
            }
        }
        rvals, _ := p.selectorProduction(args, documentMap)

        var ancestors []*Node
        for _, node := range rvals {
            ancestors = append(ancestors, node.parent)
        }
        return func(node *Node)bool {
            logger.Print("pclassFuncProduction expr ? ", &node, " ∈ ", ancestors)
            return nodeIsMemberOfList(node, ancestors)
        }, tokens
    } else if pclass == "contains" {
        return func(node *Node)bool {
            logger.Print("pclassFuncProduction contains ? ", node.typ, " == ", J_STRING, " AND ", strings.Count(node.value.(string), args[0].val.(string)), " > 0")
            return node.typ == J_STRING && strings.Count(node.value.(string), args[0].val.(string)) > 0
        }, tokens
    } else if pclass == "val" {
        return func(node *Node)bool {
            logger.Print("pclassFuncProduction val ? ", node.typ, " == ", J_STRING, " AND ", node.value.(string), " == ", args[0].val.(string))
            return node.typ == J_STRING && node.value.(string) == args[0].val.(string)
        }, tokens
    }

    // If we didn't find a known pclass, do not match anything.
    logger.Print("Error: Unknown pclass: ", pclass)
    return func(node *Node)bool {
        logger.Print("Asserting false due to failed pclassFuncProduction")
        return false
    }, tokens
}

func (p *Parser) evaluateParsedExpression(tokens []*token, node *Node, cmap map[string]func(interface{}, interface{})float64) float64 {
    var matched bool
    var lhs float64
    var rhs float64

    if len(tokens) < 1 {
        return -1
    }

    value, matched, _ := p.peek(tokens, S_PAREN)
    if matched && value.(string) == "(" {
        _, tokens, _ = p.match(tokens, S_PAREN)
        lhs = p.evaluateParsedExpression(tokens, node, cmap)
        return lhs
    }

    _, matched, _ = p.peek(tokens, S_PVAR)
    if matched {
        _, tokens, _ = p.match(tokens, S_PVAR)
        lhs = getFloat64(node.value)
    } else {
        relevantTokens := []tokenType{S_STRING, S_BOOL, S_NIL, S_NUMBER}
        for _, ttype := range relevantTokens {
            _, matched, _ = p.peek(tokens, ttype)
            if matched {
                var matchedValue interface{}
                matchedValue, tokens, _ = p.match(tokens, ttype)
                lhs = getFloat64(matchedValue)
                break
            }
        }
    }

    value, matched, _ = p.peek(tokens, S_PAREN)
    if matched && value.(string) == ")" {
        _, tokens, _ = p.match(tokens, S_PAREN)
        // Short-circuit?
        return lhs
    }

    binop, tokens, _ := p.match(tokens, S_BINOP)
    comparatorFunction := cmap[binop.(string)]
    rhs = p.evaluateParsedExpression(tokens, node, cmap)

    return comparatorFunction(lhs, rhs)
}

func (p *Parser) parseExpression(tokens []*token, node *Node) float64 {
    comparatorMap := map[string]func(lhs interface{}, rhs interface{})float64{
        "*": func(lhs interface{}, rhs interface{})float64 {
            return getFloat64(lhs) * getFloat64(rhs)
        },
        "/": func(lhs interface{}, rhs interface{})float64 {
            return getFloat64(lhs) / getFloat64(rhs)
        },
        "%": func(lhs interface{}, rhs interface{})float64 {
            return float64(getInt32(lhs) % getInt32(rhs))
        },
        "+": func(lhs interface{}, rhs interface{})float64 {
            return getFloat64(lhs) + getFloat64(rhs)
        },
        "-": func(lhs interface{}, rhs interface{})float64 {
            return getFloat64(lhs) - getFloat64(rhs)
        },
        "<=": func(lhs interface{}, rhs interface{})float64 {
            if getFloat64(lhs) <= getFloat64(rhs) {
                return 1
            }
            return 0
        },
        "<": func(lhs interface{}, rhs interface{})float64 {
            if getFloat64(lhs) < getFloat64(rhs) {
                return 1
            }
            return 0
        },
        ">=": func(lhs interface{}, rhs interface{})float64 {
            if getFloat64(lhs) > getFloat64(rhs) {
                return 1
            }
            return 0
        },
        ">": func(lhs interface{}, rhs interface{})float64 {
            if getFloat64(lhs) > getFloat64(rhs) {
                return 1
            }
            return 0
        },
        "$=": func(lhs interface{}, rhs interface{})float64 {
            lhs_str := lhs.(string)
            rhs_str := rhs.(string)
            if strings.LastIndex(lhs_str, rhs_str) == len(lhs_str) - len(rhs_str) {
                return 1
            }
            return 0
        },
        "^=": func(lhs interface{}, rhs interface{})float64 {
            lhs_str := lhs.(string)
            rhs_str := rhs.(string)
            if strings.Index(lhs_str, rhs_str) == 0 {
                return 1
            }
            return 0
        },
        "*=": func(lhs interface{}, rhs interface{})float64 {
            lhs_str := lhs.(string)
            rhs_str := rhs.(string)
            if strings.Index(lhs_str, rhs_str) != 0 {
                return 1
            }
            return 0
        },
        "=": func(lhs interface{}, rhs interface{})float64 {
            if lhs.(string) == rhs.(string) {
                return 1
            }
            return 0
        },
        "!=": func(lhs interface{}, rhs interface{})float64 {
            if lhs.(string) != rhs.(string) {
                return 1
            }
            return 0
        },
        "&&": func(lhs interface{}, rhs interface{})float64 {
            if lhs.(bool) && rhs.(bool) {
                return 1
            }
            return 0
        },
        "||": func(lhs interface{}, rhs interface{})float64 {
            if lhs.(bool) || rhs.(bool) {
                return 1
            }
            return 0
        },
    }
    return p.evaluateParsedExpression(tokens, node, comparatorMap)
}
