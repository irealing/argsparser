package main

import (
	"errors"
	"reflect"
	"flag"
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
	args   Arguments
	values []interface{}
}

// New Create new ArgsParser object
func New(args Arguments) *ArgsParser {
	return &ArgsParser{args: args}
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
		ap.values[i] = flag.String(param, emptyString, usage)
	case reflect.Int:
		ap.values[i] = flag.Int(param, 0, usage)
	case reflect.Uint:
		ap.values[i] = flag.Uint(param, 0, usage)
	case reflect.Bool:
		ap.values[i] = flag.Bool(param, false, usage)

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
		case *int:
			iv := v.(*int)
			value.Field(i).SetInt(int64(*iv))
		case *uint:
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
	flag.Parse()
	ap.injectValues()
	err := ap.args.Validate()
	if err != nil {
		flag.PrintDefaults()
	}
	return err
}
