package jsonselect

import (
    "regexp"
)

type itemType string

type lexerItem struct {
    typ itemType
    val string
}

type scannerItem struct {
    regex string
    typ itemType
}

const (
    S_TYPE itemType = "s_type"
    S_IDENTIFIER itemType = "s_identifier"
    S_QUOTED_IDENTIFIER itemType = "s_quoted_identifier"
    S_PCLASS itemType = "s_pclass"
    S_PCLASS_FUNC itemType = "s_pclass_func"
    S_NTH_FUNC itemType = "s_nth_func"
    S_OPER itemType = "s_oper"
    S_EMPTY itemType = "s_empty"
    S_UNK itemType = "s_unk"
    S_FLOAT itemType = "s_float"
    S_WORD itemType = "s_word"
    S_BINOP itemType = "s_binop"
    S_VALS itemType = "s_vals"
    S_KEYWORD itemType = "s_keyword"
    S_PVAR itemType = "s_pvar"
    S_EXPR itemType = "s_expr"
    S_NUMBER itemType = "s_number"
    S_STRING itemType = "s_string"
    S_PAREN itemType = "s_paren"
)

var scanner = [17]scannerItem{
    scannerItem{
        `\([^\)]+\)`,
        S_EXPR,
    },
    scannerItem{
        `[~*,>]`,
        S_OPER,
    },
    scannerItem{
        `\s`,
        S_EMPTY,
    },
    scannerItem{
        `(-?\d+(\.\d*)([eE][+\-]?\d+)?)`,
        S_FLOAT,
    },
    scannerItem{
        `string|boolean|null|array|object|number`,
        S_TYPE,
    },
    scannerItem{
        `\"([_a-zA-Z]|[^\0-\0177]|\\[^\s0-9a-fA-F])([_a-zA-Z0-9\-]|[^\u0000-\u0177]|(\\[^\s0-9a-fA-F]))*\"`,
        S_WORD,
    },
    scannerItem{
        `\.?\"([^"\\]|\\[^"])*\"`,
        S_QUOTED_IDENTIFIER,
    },
    scannerItem{
        `\.([_a-zA-Z]|[^\0-\0177]|\\[^\s0-9a-fA-F])([_a-zA-Z0-9\-]|[^\u0000-\u0177]|(\\[^\s0-9a-fA-F]))*`,
        S_IDENTIFIER,
    },
    scannerItem{
        `:(root|empty|first-child|last-child|only-child)`,
        S_PCLASS,
    },
    scannerItem{
        `:(has|expr|val|contains)`,
        S_PCLASS_FUNC,
    },
    scannerItem{
        `:(nth-child|nth-last-child)`,
        S_NTH_FUNC,
    },
    scannerItem{
        `(&&|\|\||[\$\^<>!\*]=|[=+\-*/%<>])`,
        S_BINOP,
    },
    scannerItem{
        `true|false|null`,
        S_VALS,
    },
    scannerItem{
        `n`,
        S_PVAR,
    },
    scannerItem{
        `odd|even`,
        S_KEYWORD,
    },
}

var expressionScanner = [7]scannerItem{
    scannerItem{
        `\s`,
        S_KEYWORD,
    },
    scannerItem{
        `true|false|null`,
        S_VALS,
    },
    scannerItem{
        `-?\d+(\.\d*)?([eE][+\-]?\d+)?`,
        S_NUMBER,
    },
    scannerItem{
        `\"([^\]|\[^\"])*\"`,
        S_STRING,
    },
    scannerItem{
        `x`,
        S_PVAR,
    },
    scannerItem{
        `(&&|\|\||[\$\^<>!\*]=|[=+\-*/%<>])`,
        S_BINOP,
    },
    scannerItem{
        `\(|\)`,
        S_PAREN,
    },
}
