package jsonselect

import (
    "io/ioutil"
    "reflect"
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
            var stringTemporary string
            for _, str := range strings.Split(string(output_document), "\n") {
                stringTemporary = stringTemporary + str
                // Try to parse -- if it works, we have the whole object
                if strings.Index(stringTemporary, "{") == 0 {
                    if strings.Count(stringTemporary, "{") != strings.Count(stringTemporary, "}") {
                        continue
                    }
                    actualOutput = append(actualOutput, stringTemporary)
                    stringTemporary = ""
                } else if strings.Index(stringTemporary, "[") == 0 {
                    if strings.Count(stringTemporary, "[") != strings.Count(stringTemporary, "]") {
                        continue
                    }
                    actualOutput = append(actualOutput, stringTemporary)
                    stringTemporary = ""
                } else if len(stringTemporary) > 0 {
                    actualOutput = append(actualOutput, stringTemporary)
                    stringTemporary = ""
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
            t.Error("Test ", testName, "failed: ", err)
        }
        selectorString := testSelectors[testName]
        expectedOutput := testOutput[testName]

        results, err := parser.GetJsonElements(selectorString)
        if err != nil {
            t.Error("Test ", testName, "failed: ", err)
        }
        var stringResults []string
        for _, result := range results {
            encoded, err := result.Encode()
            if err != nil {
                t.Error("Test ", testName, "failed: ", err)
            }
            stringResults = append(stringResults, string(encoded))
        }

        if len(stringResults) != len(expectedOutput) {
            t.Error("Test ", testName, " failed due to number of results being mismatched; ", len(stringResults), " != ", len(expectedOutput), ": [Actual] ", stringResults, " != [Expected] ", expectedOutput)
        } else {
            var matched bool = true
            for idx, result := range stringResults {
                expectedEncoded := expectedOutput[idx]
                // If the string begins with {, let's load it using simplejson,
                // convert it to a map, and use reflection to see if they do match.
                if strings.Index(strings.TrimSpace(expectedEncoded), "{") == 0 {
                    expectedJson, err := simplejson.NewJson([]byte(expectedEncoded))
                    if err != nil {
                        t.Error(
                            "Test ", testName, " failed due to a JSON decoding error while decoding expectation: ", err,
                        )
                    }
                    resultJson, err := simplejson.NewJson([]byte(result))
                    if err != nil {
                        t.Error(
                            "Test ", testName, " failed due to a JSON decoding error while decoding result: ", err,
                        )
                    }
                    result := reflect.DeepEqual(
                        expectedJson.MustMap(),
                        resultJson.MustMap(),
                    )
                    if !result {
                        matched = false
                    }
                } else if strings.Index(strings.TrimSpace(expectedEncoded), "[") == 0 {
                    expectedJson, err := simplejson.NewJson([]byte(expectedEncoded))
                    if err != nil {
                        t.Error(
                            "Test ", testName, " failed due to a JSON decoding error while decoding expectation: ", err,
                        )
                    }
                    resultJson, err := simplejson.NewJson([]byte(result))
                    if err != nil {
                        t.Error(
                            "Test ", testName, " failed due to a JSON decoding error while decoding result: ", err,
                        )
                    }
                    result := reflect.DeepEqual(
                        expectedJson.MustArray(),
                        resultJson.MustArray(),
                    )
                    if !result {
                        matched = false
                    }
                } else if expectedEncoded != result {
                    matched = false
                }
                if !matched {
                    t.Error(
                        "Test ", testName, " failed on item #", idx, ": [Actual] ", result, " != [Expected] ", expectedEncoded,
                    )
                }
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
