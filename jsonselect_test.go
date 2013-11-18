package jsonselect

import (
    "io/ioutil"
    "strings"
    "testing"
    "github.com/latestrevision/go-simplejson"
)

func getTestParser(testDocuments map[string]*simplejson.Json, testName string) (*Parser, error) {
    jsonDocument := testDocuments[testName[0:strings.Index(testName, "_")]]
    return CreateParser(jsonDocument)
}

func runTestsInDirectory(t *testing.T, baseDirectory string) {
    var testDocuments = make(map[string]*simplejson.Json)
    var testSelectors = make(map[string]string)
    var testOutput = make(map[string][]string)

    files, err := ioutil.ReadDir(baseDirectory)
    if err != nil {
        t.Error("Error encountered while loading conformance tests ", err)
    }

    for _, fileInfo := range(files) {
        name := fileInfo.Name()
        if strings.HasSuffix(name, ".json") {
            json_document, err := ioutil.ReadFile(baseDirectory + name)
            if err != nil {
                t.Error("Error encountered while reading ", name, ": ", err)
                continue
            }
            parsed_document, err := simplejson.NewJson(json_document)
            if err != nil {
                t.Error("Error encountered while deserializing ", name, ": ", err)
                continue
            }
            testDocuments[name[0:len(name)-len(".json")]] = parsed_document
        } else if (strings.HasSuffix(name, ".output")) {
            output_document, err := ioutil.ReadFile(baseDirectory + name)
            if err != nil {
                t.Error("Error encountered while reading ", name, ": ", err)
                continue
            }
            var actualOutput []string
            for _, str := range strings.Split(string(output_document), "\n") {
                if len(str) > 0 {
                    actualOutput = append(actualOutput, str)
                }
            }
            testOutput[name[0:len(name)-len(".output")]] = actualOutput
        } else if (strings.HasSuffix(name, ".selector")) {
            selector_document, err := ioutil.ReadFile(baseDirectory + name)
            if err != nil {
                t.Error("Error encountered while reading ", name, ": ", err)
                continue
            }
            testSelectors[name[0:len(name)-len(".selector")]] = string(selector_document)
        }
    }

    for testName := range testSelectors {
        parser, err := getTestParser(testDocuments, testName)
        if err != nil {
            t.Error("Unable to find testing document for ", testName)
        }
        selectorString := testSelectors[testName]
        expectedOutput := testOutput[testName]

        results, err := parser.GetJsonElements(selectorString)
        if err != nil {
            t.Error("Error encountered while finding matches ", err)
        }
        var stringResults []string
        for _, result := range results {
            encoded, err := result.Encode()
            if err != nil {
                t.Error("Error encoding result '", result, "' in to JSON")
            }
            stringResults = append(stringResults, string(encoded))
        }

        if len(stringResults) != len(expectedOutput) {
            t.Error("Test ", testName, " failed due to number of results being mismatched: [Actual] ", stringResults, " (", len(stringResults), ") != [Expected] ", expectedOutput, " (", len(expectedOutput), ")")
        }
        for idx, result := range stringResults {
            expectedEncoded := expectedOutput[idx]
            if expectedEncoded != result {
                t.Error(
                    "Test ", testName, " failed on item #", idx, ": [Actual] ", result, " != [Expected] ", expectedEncoded,
                )
            }
        }
    }
}

func TestLevel1(t *testing.T) {
    runTestsInDirectory(t, "./conformance_tests/level_1/")
}

func TestLevel2(t *testing.T) {
    runTestsInDirectory(t, "./conformance_tests/level_2/")
}

//func xTestLevel3(t *testing.T) {
//    runTestsInDirectory(t, "./conformance_tests/level_3/")
//}
