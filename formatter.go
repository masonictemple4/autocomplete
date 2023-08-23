package autocomplete

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

// The formatter is used to define formatters to assign to the data providers.
// This allow us to both provide a stable API with default options, but also offers
// the flexibility to define custom formatters if you see fit.
//
// TIP: Verify you have implemented the Formatter interface correctly by
//
//	var _ formatter.Formatter = (*YourTypeHere)(nil)
//
// NOTE: Though it is not required to satisfy the interface, it is the standard
// to create a type alias if you're not using a user defined struct. For example:
// `type DefaultFormat []string`
//
// Implementing the Formatter interface only requires one method.
// Format. It takes the file data and returns a slice of strings
// (keywords).
type Formatter interface {
	FormatRead(data []byte, fileName string) ([]string, error)
	FormatWrite(keywords []string, fileName string) ([]byte, error)
}

// DefaultFormat requires that your file decode into a slice of strings.
// Basically a non-nested JSON array of strings.
//
// TYPE: type DefaultFormat []string
//
// Example: keywords.json
//
//	[
//	  "keyword1",
//	  "keyword2",
//	  "keyword3"
//	]
//
// Example: keywords.txt
//
//	keyword1
//	keyword2
//	keyword3
//
// Example: keywords.csv
//
//	keyword1,keyword2,keyword3
//
// Example: keywords.yaml
//
//   - keyword1
//   - keyword2
//   - keyword3
type DefaultFormat []string

func (f DefaultFormat) FormatRead(data []byte, fileName string) ([]string, error) {
	fType := detectFileType(fileName)
	switch fType {
	case "json":
		var obj DefaultFormat
		if err := json.Unmarshal(data, &obj); err != nil {
			return nil, err
		}
		return obj, nil
	case "txt":
		return strings.Split(string(data), "\n"), nil
	case "csv":
		// Use your preferred CSV parsing library here
		// For instance, you can use the 'encoding/csv' package provided by the standard library
		srcRdr := bytes.NewReader(data)
		reader := csv.NewReader(srcRdr)

		full, err := reader.ReadAll()
		if err != nil {
			return nil, err
		}

		var results []string
		// Skips headers
		for _, innerObj := range full[1:] {
			results = append(results, innerObj...)
		}

		return results, nil
	case "yaml":
		var obj DefaultFormat
		if err := yaml.Unmarshal(data, &obj); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Invalid file type")
	}
	// Technically this should be covered by the default block.. But the compiler
	// isn't picking up on that.
	return nil, errors.New("Unhandled error")

}
func (f DefaultFormat) FormatWrite(keywords []string, fileName string) ([]byte, error) {
	fType := detectFileType(fileName)
	switch fType {
	case "json":
		return json.Marshal(keywords)
	case "txt":
		return []byte(strings.Join(keywords, "\n")), nil
	case "csv":
		// Use your preferred CSV parsing library here
		// For instance, you can use the 'encoding/csv' package provided by the standard library
		var buf bytes.Buffer
		writer := csv.NewWriter(&buf)
		writer.Write(keywords)
		writer.Flush()
		return buf.Bytes(), nil
	case "yaml":
		return yaml.Marshal(keywords)
	default:
		return nil, errors.New("Invalid file type")
	}

}

// KeywordObjectList requires a top level object named "keywords"
// with a value of a slice of strings.
//
//	TYPE: type KeywordObjectList struct {
//		Keywords []string `json:"keywords" yaml:"keywords"`
//	}
//
// Example: keywords.json
//
//	{
//	  "keywords": [
//	    "keyword1",
//	    "keyword2",
//	    "keyword3"
//	  ]
//	}
//
// Example: keywords.yaml
//
//	keywords:
//	  - keyword1
//	  - keyword2
//	  - keyword3
//
// Example: keywords.csv
//
//	keywords
//	keyword1,keyword2,keyword3
//
// Example: keywords.txt
//
//	keywords
//	keyword1
//	keyword2
//	keyword3
type KeywordObjectListFormat struct {
	Keywords []string `json:"keywords" yaml:"keywords"`
}

func (k KeywordObjectListFormat) FormatRead(data []byte, fileName string) ([]string, error) {
	fType := detectFileType(fileName)

	switch fType {
	case "json":
		var obj KeywordObjectListFormat
		err := json.Unmarshal(data, &obj)
		if err != nil {
			return nil, err
		}
		return obj.Keywords, nil
	case "txt":
		results := strings.Split(string(data), "\n")
		if results[0] == "keywords" {
			return results[1:], nil
		}
		return results, nil
	case "csv":
		// Use your preferred CSV parsing library here
		reader := csv.NewReader(bytes.NewReader(data))
		return reader.Read()
	case "yaml":
		var obj KeywordObjectListFormat
		err := yaml.Unmarshal(data, &obj)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Invalid file type")
	}
	// Technically this should be covered by the default block.. But the compiler
	// isn't picking up on that.
	return nil, errors.New("Unhandled error")
}

func (k KeywordObjectListFormat) FormatWrite(keywords []string, fileName string) ([]byte, error) {
	fType := detectFileType(fileName)

	switch fType {
	case "json":
		obj := KeywordObjectListFormat{Keywords: keywords}
		return json.Marshal(obj)
	case "txt":
		var buffer bytes.Buffer
		buffer.WriteString("keywords\n")
		for _, keyword := range keywords {
			buffer.WriteString(keyword)
			buffer.WriteString("\n")
		}
		return buffer.Bytes(), nil
	case "csv":
		var buf bytes.Buffer
		writer := csv.NewWriter(&buf)
		writer.Write(keywords)
		writer.Flush()
		return buf.Bytes(), nil
	case "yaml":
		obj := KeywordObjectListFormat{Keywords: keywords}
		return yaml.Marshal(obj)
	default:
		return nil, errors.New("Invalid file type")
	}
}

// There might be a better way of doing this in the future. I have tried with the bytes
// using http.DetectContentType(data) and not as much help as it should be. Will have to
// research later to see if there is another way of detecting file type.
func detectFileType(fileName string) string {
	parts := strings.Split(fileName, ".")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-1]
}
