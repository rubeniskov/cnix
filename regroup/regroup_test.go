package regroup

import (
	"fmt"
	"strings"
	"testing"
)

type Custom struct {
	Foo string
	Bar string
}

func (c *Custom) UnmarshalText(data []byte) error {
	str := string(data)
	chnks := strings.Split(str, ",")
	c.Foo = chnks[0]
	c.Bar = chnks[1]
	return nil
}

func TestRecorder(t *testing.T) {

	type SuiteTypes struct {
		StringVal 		string 	`regroup:"string"`
		UintVal 		uint 	`regroup:"uint"`
		Uint8Val 		uint8 	`regroup:"uint"`
		Uint16Val 		uint16 	`regroup:"uint"`
		Uint32Val 		uint32 	`regroup:"uint"`
		Uint64Val 		uint64 	`regroup:"uint"`
		IntVal 			int 	`regroup:"int"`
		Int8Val 		int8 	`regroup:"int"`
		Int16Val 		int16 	`regroup:"int"`
		Int32Val 		int32 	`regroup:"int"`
		Int64Val 		int64 	`regroup:"int"`
		Float32Val 		float32 `regroup:"float"`
		Float64Val 		float64 `regroup:"float"`
		StructVal 		Custom 	`regroup:"custom"`
		StructPtrVal 	*Custom `regroup:"custom"`
	}

	tests := [...]struct {
		name   string
		pattern   string
		subject   string
		target 	  interface{}
		expected  interface{}
	}{
		{
			"Default parse groups",
			`(?P<string>\w+)\s(?P<uint>\d+)\s(?P<int>-?\d+)\s(?P<float>\d+.\d+)\s(?P<custom>\w+,\w+)`,
			"foobar 100 -100 100.10 foo,bar",
			&SuiteTypes{
				StructPtrVal: &Custom{},
			},
			&SuiteTypes{ 
				"foobar", 
				100, 100, 100, 100, 100, 
				-100, -100, -100, -100, -100, 100.10, 100.10, 
				Custom{"foo", "bar"},
				&Custom{"foo", "bar"}, 
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			re := MustCompile(tt.pattern)

			err := re.MatchToTarget(tt.subject, tt.target)

			if err != nil {
				t.Error(err)
			}

			fmt.Printf("%#v\n", tt.target)

			if p, ok := tt.target.(*SuiteTypes); ok {
				fmt.Printf("%d\n", p.UintVal)
			}
		})
	}
}