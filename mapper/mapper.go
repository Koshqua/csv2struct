package mapper

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

var fieldRe = regexp.MustCompile(`(?m)[^a-zA-Z0-9]|^\d`)

type CaseType int

//go:embed struct.temp
var structTemp string

const (
	Pascal CaseType = iota + 1
	Camel
	Kebab
	Snake
	Space
)

const (
	typeInt    = "int64"
	typeBool   = "bool"
	typeString = "string"
	typeFloat  = "float64"
)

type Mapper struct {
	Config Config
}

func (m *Mapper) parseCsv(file io.Reader) (*parsedCsv, error) {
	r := csv.NewReader(file)
	r.Comma = []rune(m.Config.CsvSeparator)[0]
	r.TrimLeadingSpace = true
	p := new(parsedCsv)
	headers, err := r.Read()
	if err != nil && err != io.EOF {
		return p, fmt.Errorf("couldn't read headers w error %v", err)
	}
	p.headers = headers
	values, err := r.Read()
	if err != nil && err != io.EOF {
		return p, fmt.Errorf("couldn't read values w error %v", err)
	}
	p.values = values
	return p, nil
}

func (m *Mapper) normalizeHeaders(p *parsedCsv) ([]string, error) {
	switch m.Config.WordCaseType {
	case Pascal:
		normalized := make([]string, 0, len(p.headers))
		for _, h := range p.headers {
			h = string(fieldRe.ReplaceAll([]byte(h), []byte("")))
			normalized = append(normalized, h)
		}
		return normalized, nil
	case Camel:
		normalized := make([]string, 0, len(p.headers))
		for _, h := range p.headers {
			h = string(fieldRe.ReplaceAll([]byte(h), []byte("")))
			normalized = append(normalized, strings.Title(h))
		}
		return normalized, nil
	case Snake:
		normalized := make([]string, 0, len(p.headers))
		for _, h := range p.headers {
			words := strings.Split(h, "_")
			for i := 0; i < len(words); i++ {
				words[i] = string(fieldRe.ReplaceAll([]byte(words[i]), []byte("")))
				words[i] = strings.Title(strings.ToLower(words[i]))
			}
			h = strings.Join(words, "")
			h = string(fieldRe.ReplaceAll([]byte(h), []byte("")))
			normalized = append(normalized, h)
		}
		return normalized, nil
	case Kebab:
		normalized := make([]string, 0, len(p.headers))
		for _, h := range p.headers {
			words := strings.Split(h, "-")
			for i := 0; i < len(words); i++ {
				words[i] = string(fieldRe.ReplaceAll([]byte(words[i]), []byte("")))
				words[i] = strings.Title(strings.ToLower(words[i]))
			}
			h = strings.Join(words, "")
			h = string(fieldRe.ReplaceAll([]byte(h), []byte("")))
			normalized = append(normalized, h)
		}
		return normalized, nil
	case Space:
		normalized := make([]string, 0, len(p.headers))
		for _, h := range p.headers {
			words := strings.Split(h, " ")
			for i := 0; i < len(words); i++ {
				words[i] = string(fieldRe.ReplaceAll([]byte(words[i]), []byte("")))
				words[i] = strings.Title(strings.ToLower(words[i]))
			}
			h = strings.Join(words, "")
			h = string(fieldRe.ReplaceAll([]byte(h), []byte("")))
			normalized = append(normalized, h)
		}
		return normalized, nil
	}
	return nil, fmt.Errorf("non existent case")
}

func (m *Mapper) getFieldTypes(vals []string) []string {
	types := make([]string, 0, len(vals))
	for _, val := range vals {
		if _, err := strconv.ParseBool(val); err == nil {
			types = append(types, typeBool)
			continue
		}
		if _, err := strconv.ParseInt(val, 10, 64); err == nil {
			types = append(types, typeInt)
			continue
		}
		if _, err := strconv.ParseFloat(val, 64); err == nil {
			types = append(types, typeFloat)
			continue
		}
		types = append(types, typeString)
	}
	return types
}

func (m *Mapper) CreateStructFromCsv() (string, error) {
	csvFile, err := os.OpenFile(m.Config.From, os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}
	defer csvFile.Close()
	parsed, err := m.parseCsv(csvFile)
	if err != nil {
		return "", err
	}
	normalizedHeader, err := m.normalizeHeaders(parsed)
	if err != nil {
		return "", err
	}
	valueTypes := m.getFieldTypes(parsed.values)
	if len(normalizedHeader) != len(valueTypes) && len(normalizedHeader) != len(parsed.headers) {
		return "", fmt.Errorf("got %v headers and %v values, not matching", len(normalizedHeader), len(valueTypes))
	}
	pkgName := getPackageName()
	if pkgName == "" {
		pkgName = "placeholder_package_name"
	}
	typeToGen := GeneratedType{
		PkgName: pkgName,
		Name:    m.Config.TypeName,
	}
	fields := make([]Field, 0, len(normalizedHeader))
	for i := 0; i < len(normalizedHeader); i++ {
		field := Field{
			Name:   normalizedHeader[i],
			Type:   valueTypes[i],
			CsvTag: parsed.headers[i],
		}
		fields = append(fields, field)
	}
	log.Printf("%#+v", fields)
	typeToGen.Fields = fields
	templ := template.Must(template.New("type_temp").Parse(structTemp))
	w := bytes.NewBuffer([]byte{})
	err = templ.Execute(w, typeToGen)
	return w.String(), err
}

func getPackageName() string {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedName,
	}
	pkgs, err := packages.Load(cfg)
	if err != nil {
		panic(err)
	}
	if len(pkgs) != 1 {
		panic(fmt.Errorf("expected to get only only 1 package, got %v", len(pkgs)))
	}
	return pkgs[0].Name
}
