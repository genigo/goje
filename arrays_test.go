package goje

import (
	"database/sql/driver"
	"testing"
)

func TestStringArrayScan(t *testing.T) {
	cases := []struct {
		literal string
		want    StringArray
	}{
		{"{}", StringArray{}},
		{"{go,sql}", StringArray{"go", "sql"}},
		{`{"has space","has,comma"}`, StringArray{"has space", "has,comma"}},
		{`{"quote\"inside","back\\slash"}`, StringArray{`quote"inside`, `back\slash`}},
		{"{NULL,x}", StringArray{"", "x"}},
		{"{a}", StringArray{"a"}},
	}

	for _, tc := range cases {
		var got StringArray
		if err := got.Scan(tc.literal); err != nil {
			t.Errorf("Scan(%q): %+v", tc.literal, err)
			continue
		}
		if len(got) != len(tc.want) {
			t.Errorf("Scan(%q) = %v, want %v", tc.literal, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("Scan(%q) = %v, want %v", tc.literal, got, tc.want)
				break
			}
		}
	}

	var nilArr StringArray
	if err := nilArr.Scan(nil); err != nil || nilArr != nil {
		t.Errorf("Scan(nil) = %v err=%+v, want nil slice", nilArr, err)
	}

	if err := (&StringArray{}).Scan("not an array"); err == nil {
		t.Error("invalid literal accepted")
	}
}

func TestArrayValue(t *testing.T) {
	val, err := StringArray{"go", "has,comma", `q"uote`}.Value()
	if err != nil {
		t.Fatal(err)
	}
	want := `{go,"has,comma","q\"uote"}`
	if val != want {
		t.Errorf("Value() = %v, want %v", val, want)
	}

	// value/scan round trip through the database text format
	roundTrips := []StringArray{
		{},
		{"a"},
		{"", "b"},
		{"with space", "with,comma", `with"quote`, `with\slash`, "{}"},
	}
	for _, arr := range roundTrips {
		val, err := arr.Value()
		if err != nil {
			t.Fatal(err)
		}
		var got StringArray
		if err := got.Scan(val); err != nil {
			t.Fatal(err)
		}
		if len(got) != len(arr) {
			t.Fatalf("round trip %v -> %v -> %v", arr, val, got)
		}
		for i := range got {
			if got[i] != arr[i] {
				t.Fatalf("round trip %v -> %v -> %v", arr, val, got)
			}
		}
	}
}

func TestNumericArrayScan(t *testing.T) {
	var ints Int32Array
	if err := ints.Scan("{1,2,3}"); err != nil {
		t.Fatal(err)
	}
	if len(ints) != 3 || ints[0] != 1 || ints[2] != 3 {
		t.Errorf("Int32Array = %v", ints)
	}
	if v, _ := ints.Value(); v != "{1,2,3}" {
		t.Errorf("Int32Array.Value() = %v", v)
	}

	var i64 Int64Array
	if err := i64.Scan("{9223372036854775807}"); err != nil {
		t.Fatal(err)
	}
	if i64[0] != 9223372036854775807 {
		t.Errorf("Int64Array = %v", i64)
	}

	var fl Float64Array
	if err := fl.Scan("{1.5,2.25}"); err != nil {
		t.Fatal(err)
	}
	if fl[0] != 1.5 || fl[1] != 2.25 {
		t.Errorf("Float64Array = %v", fl)
	}

	var b BoolArray
	if err := b.Scan("{t,f,true,false}"); err != nil {
		t.Fatal(err)
	}
	if !b[0] || b[1] || !b[2] || b[3] {
		t.Errorf("BoolArray = %v", b)
	}
	if v, _ := b.Value(); v != "{true,false,true,false}" {
		t.Errorf("BoolArray.Value() = %v", v)
	}
}

func TestArrayScanBytes(t *testing.T) {
	// database/sql may deliver the literal as []byte
	var got StringArray
	if err := got.Scan([]byte("{a,b}")); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[1] != "b" {
		t.Errorf("bytes scan = %v", got)
	}
}

var _ driver.Valuer = StringArray{}
