package codec

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/farshidtz/senml/v2"
)

// CSVHeader is the fixed header to support records with different value types
const CSVHeader = "Time,Update Time,Name,Unit,Value,String Value,Boolean Value,Data Value,Sum"

func WriteCSV(p senml.Pack, w io.Writer, options ...Option) error {
	o := &Options{
		PrettyPrint: false,
	}
	for _, opt := range options {
		opt(o)
	}

	csvWriter := csv.NewWriter(w)

	if o.WithHeader {
		err := csvWriter.Write(strings.Split(CSVHeader, ","))
		if err != nil {
			return err
		}
	}

	// normalize first to add base values to row values
	p.Normalize()

	for i := range p {
		row := make([]string, 9)
		row[0] = strconv.FormatFloat(p[i].Time, 'f', -1, 64)
		row[1] = strconv.FormatFloat(p[i].UpdateTime, 'f', -1, 64)
		row[2] = p[i].Name
		row[3] = p[i].Unit
		if p[i].Value != nil {
			row[4] = strconv.FormatFloat(*p[i].Value, 'f', -1, 64)
		}
		row[5] = p[i].StringValue
		if p[i].BoolValue != nil {
			row[6] = fmt.Sprintf("%t", *p[i].BoolValue)
		}
		row[7] = p[i].DataValue
		if p[i].Sum != nil {
			row[8] = strconv.FormatFloat(*p[i].Sum, 'f', -1, 64)
		}

		err := csvWriter.Write(row)
		if err != nil {
			return err
		}
	}
	csvWriter.Flush() // TODO flush during the iterations?
	if err := csvWriter.Error(); err != nil {
		return err
	}
	return nil
}

// EncodeCSV serializes the SenML pack into CSV bytes
func EncodeCSV(p senml.Pack, options ...Option) ([]byte, error) {

	var buf bytes.Buffer
	err := WriteCSV(p, &buf, options...)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ReadCSV(r io.Reader, options ...Option) (senml.Pack, error) {
	o := &Options{
		PrettyPrint: false,
	}
	for _, opt := range options {
		opt(o)
	}

	csvReader := csv.NewReader(r)

	if o.WithHeader {
		row, err := csvReader.Read()
		if err == io.EOF {
			return nil, fmt.Errorf("missing header or no input")
		}
		if err != nil {
			return nil, err
		}
		if joined := strings.Join(row, ","); joined != CSVHeader {
			return nil, fmt.Errorf("unexpected header: %s. Expected: %s", joined, CSVHeader)
		}
	}

	var p senml.Pack
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		var record senml.Record
		// Time
		record.Time, err = strconv.ParseFloat(row[0], 10)
		if err != nil {
			return nil, err
		}
		// Update Time
		record.UpdateTime, err = strconv.ParseFloat(row[8], 10)
		if err != nil {
			return nil, err
		}
		// Name
		record.Name = row[1]
		// Unit
		record.Unit = row[2]
		// Value
		if row[3] != "" {
			value, err := strconv.ParseFloat(row[3], 10)
			if err != nil {
				return nil, err
			}
			record.Value = &value
		}
		// String Value
		record.StringValue = row[4]
		// Boolean Value
		if row[5] != "" {
			boolValue, err := strconv.ParseBool(row[5])
			if err != nil {
				return nil, err
			}
			record.BoolValue = &boolValue
		}
		// Data Value
		record.DataValue = row[6]
		// Sum
		if row[7] != "" {
			sum, err := strconv.ParseFloat(row[7], 10)
			if err != nil {
				return nil, err
			}
			record.Sum = &sum
		}

		p = append(p, record)
	}

	return p, nil
}

// DecodeCSV takes a SenML pack in CSV bytes and decodes it into a Pack
func DecodeCSV(b []byte, options ...Option) (senml.Pack, error) {

	p, err := ReadCSV(bytes.NewReader(b), options...)
	if err != nil {
		return nil, err
	}

	return p, nil
}
