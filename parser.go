package argsparser

import (
	"errors"
	"flag"
	"os"
	"reflect"
)

const (
	paramKey    = "param"
	usageKey    = "usage"
	emptyString = ""
)

var (
	UnsupportedType = errors.New("unsupported type")
)

type Arguments interface {
	Validate() error
}

type ArgsParser struct {
	args    Arguments
	values  []interface{}
	flagSet *flag.FlagSet
}

// New Create new ArgsParser object
func New(args Arguments) *ArgsParser {
	return newParser(os.Args[0], args)
}

func newParser(name string, args Arguments) *ArgsParser {
	return &ArgsParser{
		args:    args,
		flagSet: flag.NewFlagSet(name, flag.ExitOnError),
	}
}

// Init Init the ArgsParser
func (ap *ArgsParser) Init() error {
	return ap.register()
}

// register register all flags
func (ap *ArgsParser) register() error {
	argv := reflect.ValueOf(ap.args)
	if argv.Type().Kind() != reflect.Ptr {
		return UnsupportedType
	}
	argv = argv.Elem()
	at := argv.Type()
	fieldsNum := at.NumField()
	ap.values = make([]interface{}, fieldsNum)
	for i := 0; i < fieldsNum; i++ {
		fv := at.Field(i)
		ap.regFlag(i, &fv)
	}
	return nil
}
func (ap *ArgsParser) regFlag(i int, st *reflect.StructField) {
	param := st.Tag.Get(paramKey)
	usage := st.Tag.Get(usageKey)
	if param == emptyString {
		return
	}
	switch st.Type.Kind() {
	case reflect.String:
		ap.values[i] = ap.flagSet.String(param, emptyString, usage)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ap.values[i] = ap.flagSet.Int(param, 0, usage)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ap.values[i] = ap.flagSet.Uint(param, 0, usage)
	case reflect.Bool:
		ap.values[i] = ap.flagSet.Bool(param, false, usage)

	}
}

func (ap *ArgsParser) injectValues() {
	value := reflect.ValueOf(ap.args).Elem()
	for i, v := range ap.values {
		switch v.(type) {
		case nil:
			continue
		case *string:
			s := v.(*string)
			value.Field(i).SetString(*s)
		case *int, *int8, *int16, *int32, *int64:
			iv := v.(*int)
			value.Field(i).SetInt(int64(*iv))
		case *uint, *uint8, *uint16, *uint32, *uint64:
			uiv := v.(*uint)
			value.Field(i).SetUint(uint64(*uiv))
		case *bool:
			bv := v.(*bool)
			value.Field(i).SetBool(*bv)
		}
	}
}

// Parse Parse args and inject values into Arguments object
func (ap *ArgsParser) Parse() error {
	return ap.ParseValues(os.Args[1:])
}

func (ap *ArgsParser) ParseValues(values []string) error {
	ap.flagSet.Parse(values)
	ap.injectValues()
	err := ap.args.Validate()
	if err != nil {
		flag.PrintDefaults()
	}
	return err
}
func (ap *ArgsParser) PrintHelp() {
	ap.flagSet.PrintDefaults()
}
