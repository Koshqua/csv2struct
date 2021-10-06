package mapper

type Config struct {
	From, To, TypeName string
	CsvSeparator       string
	WordCaseType       CaseType
	Verbose            bool
}

type GeneratedType struct {
	PkgName string
	Name    string
	Fields  []Field
}
type Field struct {
	CsvTag string
	Name   string
	Type   string
}

type parsedCsv struct {
	headers, values []string
}
