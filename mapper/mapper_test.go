package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapper_parseCsv(t *testing.T) {
	t.Run("works with comma", func(t *testing.T) {
		csvEx := `name,nickname,blah
ivan,koshqua,blah`
		m := new(Mapper)
		m.Config = Config{
			CsvSeparator: ",",
		}
		ps, err := m.parseCsv([]byte(csvEx))
		assert.NoError(t, err)
		expectedHeaders := []string{"name", "nickname", "blah"}
		expectedValues := []string{"ivan", "koshqua", "blah"}
		assert.Equal(t, expectedHeaders, ps.headers)
		assert.Equal(t, expectedValues, ps.values)
	})
	t.Run("works with pipe", func(t *testing.T) {
		csvEx := `name|nickname|blah
ivan|koshqua|blah`
		m := new(Mapper)
		m.Config = Config{
			CsvSeparator: "|",
		}
		ps, err := m.parseCsv([]byte(csvEx))
		assert.NoError(t, err)
		expectedHeaders := []string{"name", "nickname", "blah"}
		expectedValues := []string{"ivan", "koshqua", "blah"}
		assert.Equal(t, expectedHeaders, ps.headers)
		assert.Equal(t, expectedValues, ps.values)
	})

}

func TestMapper_normalizeHeaders(t *testing.T) {
	type fields struct {
		config Config
	}
	type args struct {
		p *parsedCsv
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "pascal",
			fields: fields{
				config: Config{
					WordCaseType: Pascal,
				},
			},
			args: args{
				p: &parsedCsv{headers: []string{"1Numeric*Blah", "Abra^Kadabra"}},
			},
			want:    []string{"NumericBlah", "AbraKadabra"},
			wantErr: false,
		},
		{
			name: "camel",
			fields: fields{
				config: Config{
					WordCaseType: Camel,
				},
			},
			args: args{
				p: &parsedCsv{
					headers: []string{"1numeric*Bl ah", "abra Kadabra&***"},
				},
			},
			want: []string{"NumericBlah", "AbraKadabra"},
		},
		{
			name: "kebab",
			fields: fields{
				config: Config{
					WordCaseType: Kebab,
				},
			},
			args: args{
				p: &parsedCsv{
					headers: []string{"1numeric-*bl ah", "abra-kadabra&***"},
				},
			},
			want: []string{"NumericBlah", "AbraKadabra"},
		},
		{
			name: "snake",
			fields: fields{
				config: Config{
					WordCaseType: Snake,
				},
			},
			args: args{
				p: &parsedCsv{
					headers: []string{"1numeric_*bl ah", "abra_kadabra&***"},
				},
			},
			want: []string{"NumericBlah", "AbraKadabra"},
		},
		{
			name: "space",
			fields: fields{
				config: Config{
					WordCaseType: Space,
				},
			},
			args: args{
				p: &parsedCsv{
					headers: []string{"1numeric_* blah", "abra kadabra&***"},
				},
			},
			want: []string{"NumericBlah", "AbraKadabra"},
		},
		{
			name: "non-existent",
			fields: fields{
				config: Config{
					WordCaseType: CaseType(7),
				},
			},
			args: args{
				p: &parsedCsv{
					headers: []string{"1numeric_* blah", "abra kadabra&***"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mapper{
				Config: tt.fields.config,
			}
			got, err := m.normalizeHeaders(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("normalizeHeaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapper_getFieldTypes(t *testing.T) {
	tests := []struct {
		name string
		vals []string
		want []string
	}{
		{
			name: "some_types",
			vals: []string{"12.3", "12", "true", "false", "blah"},
			want: []string{typeFloat, typeInt, typeBool, typeBool, typeString},
		},
		{
			name: "some_other",
			vals: []string{"112381723^12312&", "2021-02-06", "999111aaa", "????", "false1"},
			want: []string{typeString, typeString, typeString, typeString, typeString},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mapper{}
			got := m.getFieldTypes(tt.vals)
			assert.Equal(t, tt.want, got)
		})
	}
}
