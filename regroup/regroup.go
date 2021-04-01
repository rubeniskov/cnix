package regroup

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ReGroup struct {
	matcher *regexp.Regexp
}

type Parseable interface {
	UnmarshalText(data []byte) error 
}

// Compile compiles given expression as regex and return new ReGroup with this expression as matching engine.
// If the expression can't be compiled as regex, a CompileError will be returned
func Compile(expr string) (*ReGroup, error) {
	matcher, err := regexp.Compile(expr)
	if err != nil {
		return nil, fmt.Errorf("compilation error: %v", err)
	}

	return &ReGroup{matcher: matcher}, nil
}

// MustCompile calls Compile and panicked if it retuned an error
func MustCompile(expr string) *ReGroup {
	reGroup, err := Compile(expr)
	if err != nil {
		panic(`regroup: Compile(` + expr + `): ` + err.Error())
	}
	return reGroup
}

func (r *ReGroup) MatchToTarget(s string, target interface{}) error {
	match := r.matcher.FindStringSubmatch(s)
	
	if match == nil {
		return fmt.Errorf("regroup: no matches found")
	}

	matchNames := r.matchGroupMap(match)

	targetRef, err := r.validateTarget(target)
	if err != nil {
		return err
	}

	/********/

	targetType := targetRef.Type()

	for i := 0; i < targetType.NumField(); i++ {
		fieldRef := targetRef.Field(i)
		if !fieldRef.CanSet() {
			continue
		}

		pep := targetType.Field(i)
		
		regroupKey, regroupOption := r.groupAndOption(pep)
		
		if regroupKey == "" {
			return nil
		}

		matchValue := matchNames[regroupKey]
		
		if regroupOption == "required" && matchValue == "" {
			return fmt.Errorf("empty: %s\n", regroupKey)
		}

		var val reflect.Value

		if matchValue != "" {
			fieldRefType := fieldRef.Type()

			if fieldRefType.Kind() == reflect.Ptr {
				if fieldRef.IsNil() {
					return fmt.Errorf("can't set value to nil pointer in field: %s", fieldRefType.Name())
				}
				fieldRef = fieldRef.Elem()
			}
			
			switch fieldRef.Type().Kind() {
				case reflect.Bool:
					b, err := strconv.ParseBool(matchValue)
					if err != nil {
						return err
					}
					val = reflect.ValueOf(b).Convert(fieldRef.Type())
					fieldRef.Set(val);
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					b, err := strconv.ParseUint(matchValue, 10, 64)
					if err != nil {
						return err
					}
					val = reflect.ValueOf(b).Convert(fieldRef.Type())
					fieldRef.Set(val);
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					b, err := strconv.ParseInt(matchValue, 10, 64)
					if err != nil {
						return err
					}
					val = reflect.ValueOf(b).Convert(fieldRef.Type())
					fieldRef.Set(val);
				case reflect.Float32, reflect.Float64:
					b, err := strconv.ParseFloat(matchValue, 64)
					if err != nil {
						return err
					}
					val = reflect.ValueOf(b).Convert(fieldRef.Type())
					fieldRef.Set(val);
				case reflect.Struct:
					niface := reflect.New(fieldRef.Type()).Interface()
					if piface, ok := niface.(Parseable); ok {
						piface.UnmarshalText([]byte(matchValue))
						fieldRef.Set(reflect.ValueOf(niface).Elem())
					} else {
						return fmt.Errorf("invalid struct, must implement UnmarshalText method to serialize see: https://golang.org/pkg/encoding/#TextUnmarshaler")
					}
				case reflect.String:
					val = reflect.ValueOf(matchValue).Convert(fieldRef.Type())
					fieldRef.Set(val);
				default:
					return fmt.Errorf("unsupported type")
			}
		}

	}

	return nil
}

func (r *ReGroup) matchGroupMap(match []string) map[string]string {
	ret := make(map[string]string)
	for i, name := range r.matcher.SubexpNames() {
		if i != 0 && name != "" {
			ret[name] = match[i]
		}
	}
	return ret
}

func (r *ReGroup) validateTarget(target interface{}) (reflect.Value, error) {
	targetPtr := reflect.ValueOf(target)
	if targetPtr.Kind() != reflect.Ptr {
		return reflect.Value{}, fmt.Errorf("regroup: target not valid")
	}
	return targetPtr.Elem(), nil
}

func (r *ReGroup) groupAndOption(fieldType reflect.StructField) (group, option string) {
	regroupKey := fieldType.Tag.Get("regroup")
	if regroupKey == "" {
		return "", ""
	}
	splitted := strings.Split(regroupKey, ",")
	if len(splitted) == 1 {
		return strings.TrimSpace(splitted[0]), ""
	}
	return strings.TrimSpace(splitted[0]), strings.TrimSpace(strings.ToLower(splitted[1]))
}