package jsonselect

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type tokenType string

type token struct {
	typ tokenType
	val interface{}
}

type scannerItem struct {
	regex *regexp.Regexp
	typ   tokenType
}

const (
	S_TYPE              tokenType = "type"
	S_IDENTIFIER        tokenType = "identifier"
	S_QUOTED_IDENTIFIER tokenType = "quoted_identifier"
	S_PCLASS            tokenType = "pclass"
	S_PCLASS_FUNC       tokenType = "pclass_func"
	S_NTH_FUNC          tokenType = "nth_func"
	S_OPER              tokenType = "operator"
	S_EMPTY             tokenType = "empty"
	S_UNK               tokenType = "unknown"
	S_FLOAT             tokenType = "float"
	S_WORD              tokenType = "word"
	S_BINOP             tokenType = "binop"
	S_BOOL              tokenType = "bool"
	S_NIL               tokenType = "null"
	S_KEYWORD           tokenType = "keyword"
	S_PVAR              tokenType = "pvar"
	S_EXPR              tokenType = "expr"
	S_NUMBER            tokenType = "number"
	S_STRING            tokenType = "string"
	S_PAREN             tokenType = "paren"
)

var selectorScanner = []scannerItem{
	scannerItem{
		regexp.MustCompile(`^\([^\)]+\)`),
		S_EXPR,
	},
	scannerItem{
		// we match any of the operators and all surrounding whitespace
		// to ensure we don't get extra space operators
		regexp.MustCompile(`^\s*[~*,> ]\s*`),
		S_OPER,
	},
	scannerItem{
		regexp.MustCompile(`^\s`),
		S_EMPTY,
	},
	scannerItem{
		regexp.MustCompile(`^(-?\d+(\.\d*)([eE][+\-]?\d+)?)`),
		S_FLOAT,
	},
	scannerItem{
		regexp.MustCompile(`^string|boolean|null|array|object|number`),
		S_TYPE,
	},
	scannerItem{
		regexp.MustCompile(`^\"([_a-zA-Z]|\\[^\s0-9a-fA-F])([_a-zA-Z0-9\-]|(\\[^\s0-9a-fA-F]))*\"`),
		S_WORD,
	},
	scannerItem{
		regexp.MustCompile(`^\.?\"([^"\\]|\\[^"])*\"`),
		S_QUOTED_IDENTIFIER,
	},
	scannerItem{
		regexp.MustCompile(`^\.([_a-zA-Z]|\\[^\s0-9a-fA-F])([_a-zA-Z0-9\-]|(\\[^\s0-9a-fA-F]))*`),
		S_IDENTIFIER,
	},
	scannerItem{
		regexp.MustCompile(`^:(root|empty|first-child|last-child|only-child)`),
		S_PCLASS,
	},
	scannerItem{
		// we match any trailing whitespace to ensure that we don't get a space operator
		// if whitespace exists before the expression. We must support whitespace before
		// the expression in order to pass the basic_has-whitespace test.
		regexp.MustCompile(`^:(has|expr|val|contains)\s*`),
		S_PCLASS_FUNC,
	},
	scannerItem{
		// we match any trailing whitespace to ensure that we don't get a space operator
		// if whitespace exists before the expression. We must support whitespace before
		// the expression in order to pass the basic_has-whitespace test.
		regexp.MustCompile(`^:(nth-child|nth-last-child)\s*`),
		S_NTH_FUNC,
	},
	scannerItem{
		regexp.MustCompile(`^(&&|\|\||[\$\^<>!\*]=|[=+\-*/%<>])`),
		S_BINOP,
	},
	scannerItem{
		regexp.MustCompile(`^true|false`),
		S_BOOL,
	},
	scannerItem{
		regexp.MustCompile(`^null`),
		S_NIL,
	},
	scannerItem{
		regexp.MustCompile(`^n`),
		S_PVAR,
	},
	scannerItem{
		regexp.MustCompile(`^odd|even`),
		S_KEYWORD,
	},
}

var expressionScanner = []scannerItem{
	scannerItem{
		regexp.MustCompile(`^\s`),
		S_EMPTY,
	},
	scannerItem{
		regexp.MustCompile(`^true|false`),
		S_BOOL,
	},
	scannerItem{
		regexp.MustCompile(`^null`),
		S_NIL,
	},
	scannerItem{
		regexp.MustCompile(`^-?\d+(\.\d*)?([eE][+\-]?\d+)?`),
		S_NUMBER,
	},
	scannerItem{
		regexp.MustCompile(`^\"([^\]|\[^\"])*\"`),
		S_STRING,
	},
	scannerItem{
		regexp.MustCompile(`^x`),
		S_PVAR,
	},
	scannerItem{
		regexp.MustCompile(`^(&&|\|\||[\$\^<>!\*]=|[=+\-*/%<>])`),
		S_BINOP,
	},
	scannerItem{
		regexp.MustCompile(`^\(|\)`),
		S_PAREN,
	},
}

func lexNextToken(input string, scanners []scannerItem) (*token, int, error) {
	for _, scanner := range scanners {
		if scanner.regex.MatchString(input) {
			idx := scanner.regex.FindStringIndex(input)
			if idx[0] == 0 {
				if scanner.typ == S_EXPR && strings.Count(input[idx[0]:idx[1]], "(") != strings.Count(input[idx[0]:idx[1]], ")") {
					terminated := false
					var curr int
					for curr = idx[1]; curr <= len(input[idx[0]:]); curr++ {
						if strings.Count(input[idx[0]:curr], "(") == strings.Count(input[idx[0]:curr], ")") {
							terminated = true
							break
						}
					}
					if terminated == false {
						return nil, len(input), errors.New(fmt.Sprintf("Unterminated expression: %s", input[idx[0]:curr-1]))
					} else {
						idx[1] = curr
					}
				}
				token := getToken(
					scanner.typ,
					input[idx[0]:idx[1]],
				)
				return &token, idx[1], nil
			}
		}
	}
	return nil, len(input), errors.New(fmt.Sprintf("Selector parsing error at %s", input))
}

func lex(input string, scanners []scannerItem) ([]*token, error) {

	// trim whitespace to ensure we don't get hanging space operators
	input = strings.TrimSpace(input)

	var tokens []*token
	var start = 0
	for start < len(input) {
		token, new_value, err := lexNextToken(input[start:], scanners)
		if err != nil {
			return nil, err
		}
		start = start + new_value
		if token.typ != S_EMPTY {
			tokens = append(
				tokens,
				token,
			)
		}
	}
	logger.Print("Tokenization results: ", input)
	if logger.Enabled {
		for i, token := range tokens {
			logger.Print("[", i, "] ", token)
		}
	}
	return tokens, nil
}

func getToken(typ tokenType, val string) token {
	switch typ {
	case S_IDENTIFIER, S_PCLASS:
		return token{typ, val[1:]}
	case S_PCLASS_FUNC, S_NTH_FUNC:
		// we match trailing whitespace in S_PCLASS_FUNC and S_NTH_FUNC to ensure
		// we don't get a space operator before the expression. So we must trim the
		// matched whitespace here
		return token{typ, strings.TrimSpace(val[1:])}
	case S_QUOTED_IDENTIFIER:
		return token{S_IDENTIFIER, val[2 : len(val)-1]}
	case S_NIL:
		return token{typ, nil}
	case S_BOOL:
		result, _ := strconv.ParseBool(val)
		return token{typ, result}
	case S_NUMBER:
		result, _ := strconv.ParseInt(val, 10, 64)
		return token{typ, result}
	case S_EMPTY:
		return token{typ, " "}
	case S_FLOAT:
		result, _ := strconv.ParseFloat(val, 32)
		return token{typ, result}
	case S_WORD:
		return token{typ, val[1 : len(val)-1]}
	case S_STRING:
		return token{S_STRING, val[1 : len(val)-1]}
	case S_OPER:
		// If the operator is padded with whitespace, we match the whole string so we must
		// trim leading and trailing whitespace.
		inner := strings.TrimSpace(val)
		if inner == "" {
			// If we're left with an empty string, we want a space operator.
			inner = " "
		}
		return token{S_OPER, inner}
	default:
		return token{typ, val}
	}
}
