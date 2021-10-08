package mapper

type Config struct {
	From, To, TypeName      string
	CsvSeparator            string
	WordCaseType            CaseType
	Verbose, AddPackageName bool
}

type GeneratedType struct {
	AddPackageName bool
	PkgName        string
	Name           string
	Fields         []Field
}
type Field struct {
	CsvTag string
	Name   string
	Type   string
}

type parsedCsv struct {
	headers, values []string
}
