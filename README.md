# csv2struct 
Easy to use library and CLI utility to generate Go struct from CSV files.
As a benefit, it's fully compatible with [csvutil](https://github.com/jszwec/csvutil).
So, structs generated by this utility can be used with that library. 

## Install 
`go install github.com/Koshqua/csv2struct@latest`

## Usage 

```bash
NAME:
   csv2struct - Converts csv files to golang structs compatible with https://github.com/jszwec/csvutil

USAGE:
   csv2struct [global options] command [command options] [arguments...]

AUTHOR:
   Ivan Malovanyi (https://github.com/Koshqua)

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --from value, -f value        specify which csv file to use
   --to value, -t value          specify the output .go file
   --typename value, --tn value  specify how to name output type
   --csvsep value, --cs value    specify the csv separator (default: ",")
   --casetype value, --ct value  specify the headers case type, possible values are: pascal, camel, kebab, snake, space (default: "pascal")
   --verbose, -v                 verbose logging (with debug) (default: false)
   --help, -h                    show help (default: false)
```
## Example
```bash
csv2struct -f ./test.csv -t ./blah.go -tn Blah --casetype space  
```
Also, it's available as library. 
Will provide usage examples a bit later...


## Contribution
Your contribution to the project is welcomed and appreciated.
