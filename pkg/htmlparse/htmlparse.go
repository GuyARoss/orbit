package htmlparse

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// innerHTML is a utility function that assists with the parsing the content of html tags
// it does this by returning the subset of the two provided strings "subStart" & "subEnd"
func innerHTML(str string, subStart string, subEnd string) string {
	return strings.Split(strings.Join(strings.Split(str, subStart)[1:], ""), subEnd)[0]
}

// HTMLDoc represents a basic document model that will be rendered upon build request
type HTMLDoc struct {
	Head []string
	Body []string
}

// render renders the document out to a single string
func (s *HTMLDoc) Render() string {
	return fmt.Sprintf(`
	<!doctype html>
	<html lang="en">
	<head>%s</head>
	<body>%s</body>
	</html>`, strings.Join(s.Head, ""), strings.Join(s.Body, ""))
}

func (s *HTMLDoc) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(s.Render())
	if err != nil {
		return err
	}

	err = file.Sync()
	return err
}

func (s *HTMLDoc) Merge(doc *HTMLDoc) *HTMLDoc {
	s.Body = append(doc.Body, s.Body...)
	s.Head = append(doc.Head, s.Head...)

	return s
}

func DocFromFile(path string) *HTMLDoc {
	data, err := ioutil.ReadFile(path)

	if len(data) == 0 || err != nil {
		return &HTMLDoc{}
	}

	return &HTMLDoc{
		Head: []string{innerHTML(string(data), "<head>", "</head>")},
		Body: []string{innerHTML(string(data), "<body>", "</body>")},
	}
}
