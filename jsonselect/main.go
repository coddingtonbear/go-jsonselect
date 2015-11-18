package main

// jsonselect is a command-line tool to apply JSONSelect filters to
// JSON through stdin, line by line.
//
// You can use as (eg. filtering Mixpanel-like event data):
//
// Extract the `event` prop, one by line:
//
//     cat jsonfile | jsonselect .event
//
// Extract two lines for each incoming line, one is the `event` property, the other the JSONPath equivalent to `.properties.os_name`
//
//     cat jsonfile | jsonselect .event ".properties .os_name"
//
// Same thing, on a single line, separated by \t characters:
//
//     cat jsonfile | jsonselect -s .event ".properties .os_name"
//
// Nicely indented properties dictionary, prefixed with the `event` as quoted (-q) text:
//
//     cat jsonfile | jsonselect -q -i .event .properties
//



import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	jsonselect "github.com/coddingtonbear/go-jsonselect"
)

func main() {
	var singleLine bool
	var quotedStrings bool
	var indent bool

	flag.BoolVar(&singleLine, "s", false, "Put things on a single line")
	flag.BoolVar(&quotedStrings, "q", false, "Keep strings quoted instead of unquoting them")
	flag.BoolVar(&indent, "i", false, "Nicely indent any JSON output")
	flag.Parse()

	errored := false

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		elements, err := elementsForAllPatterns(scanner.Text(), flag.Args())
		if err != nil {
			log.Println("Error:", err)
			return
		}

		for i, el := range elements {
			if i > 0 {
				if singleLine {
					fmt.Print("\t")
				} else {
					fmt.Print("\n")
				}
			}

			stringElement, isString := el.(string)
			if !quotedStrings && isString {
				fmt.Print(stringElement)
			} else {
				var out []byte
				if indent {
					out, err = json.MarshalIndent(el, "", "  ")
				} else {
					out, err = json.Marshal(el)
				}
				if err != nil {
					log.Printf("Error marshalling %v: %s\n", el, err)
					errored = true
				}

				os.Stdout.Write(out)
			}

		}
		fmt.Print("\n")

	}
	if scanner.Err() != nil {
		log.Println("Error scanning file:", scanner.Err())
		errored = true
	}

	if errored {
		os.Exit(1)
	}
}

func elementsForAllPatterns(body string, patterns []string) ([]interface{}, error) {
	var out []interface{}
	for _, pattern := range patterns {
		parser, err := jsonselect.CreateParserFromString(body)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling JSON, killing feed: %s", err)
		}

		elements, err := parser.GetValues(pattern)
		if err != nil {
			return nil, fmt.Errorf("Error parsing document: %s", err)
		}
		for _, el := range elements {
			out = append(out, el)
		}
	}

	return out, nil
}
