package codec

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/farshidtz/senml/v2"
)

const (
	// CSV should be in normalized form
	csvString = `946684799,10,dev123temp,degC,22.1,,,,0
946684799,0,dev123room,degC,,kitchen,,,
946684800,0,dev123data,degC,,,,abc,
946684800,0,dev123ok,degC,,,true,,
`

	csvStringWithHeader = `Time,Update Time,Name,Unit,Value,String Value,Boolean Value,Data Value,Sum
` + csvString
)

func TestEncodeCSV(t *testing.T) {

	t.Run("without header", func(t *testing.T) {
		dataOut, err := EncodeCSV(referencePack(true))
		if err != nil {
			t.Fatalf("Encoding error: %s", err)
		}

		if string(dataOut) != csvString {
			t.Logf("Expected:\n'%s'", csvString)
			t.Fatalf("Got:\n'%s'", dataOut)
		}
	})

	t.Run("with header", func(t *testing.T) {
		dataOut, err := EncodeCSV(referencePack(true), WithHeader)
		if err != nil {
			t.Fatalf("Encoding error: %s", err)
		}

		if string(dataOut) != csvStringWithHeader {
			t.Logf("Expected:\n'%s'", csvStringWithHeader)
			t.Fatalf("Got:\n'%s'", dataOut)
		}
	})
}

func TestDecodeCSV(t *testing.T) {

	t.Run("compare fields", func(t *testing.T) {

		pack, err := DecodeCSV([]byte(csvString))
		if err != nil {
			t.Fatalf("Error decoding: %s", err)
		}

		ref := referencePack(true)
		ref.Normalize()

		if err := compareFields(pack, ref); err != nil {
			t.Fatalf("Error matching records: %s", err)
		}
	})

	t.Run("wrong header", func(t *testing.T) {
		data := []byte("Bad,Time,Name,Unit,Value,String Value,Boolean Value,Data Value,Sum,Update Time")
		_, err := DecodeCSV(data, WithHeader)
		if err == nil {
			t.Fatalf("No error for wrong header")
		}
	})

}

// EXAMPLES

func ExampleEncodeCSV() {
	var pack senml.Pack = []senml.Record{
		{Time: 1276020000, Name: "room1/temp_label", StringValue: "hot"},
		{Time: 1276020100, Name: "room1/temp_label", StringValue: "cool"},
	}

	// encode to CSV (format: name,excel-time,value,unit)
	csvBytes, err := EncodeCSV(pack, WithHeader)
	if err != nil {
		panic(err) // handle the error
	}
	fmt.Printf("%s\n", csvBytes)
	// Output:
	// Time,Update Time,Name,Unit,Value,String Value,Boolean Value,Data Value,Sum
	// 1276020000,0,room1/temp_label,,,hot,,,
	// 1276020100,0,room1/temp_label,,,cool,,,
}

func ExampleDecodeCSV() {
	input := `Time,Update Time,Name,Unit,Value,String Value,Boolean Value,Data Value,Sum
946684799,10,dev123temp,degC,22.1,,,,0
946684799,0,dev123room,degC,,kitchen,,,`

	// decode JSON
	pack, err := DecodeCSV([]byte(input), WithHeader)
	if err != nil {
		panic(err) // handle the error
	}

	// validate the SenML Pack
	err = pack.Validate()
	if err != nil {
		panic(err) // handle the error
	}
}

func ExampleWriteCSV() {
	var pack senml.Pack = []senml.Record{
		{Time: 1276020000, Name: "room1/temp_label", StringValue: "hot"},
		{Time: 1276020100, Name: "room1/temp_label", StringValue: "cool"},
	}

	var writer io.Writer = os.Stdout // write to stdout
	err := WriteCSV(pack, writer, WithHeader)
	if err != nil {
		panic(err) // handle the error
	}
	// Output:
	// xTime,Update Time,Name,Unit,Value,String Value,Boolean Value,Data Value,Sum
	// 1276020000,0,room1/temp_label,,,hot,,,
	// 1276020100,0,room1/temp_label,,,cool,,,
}
