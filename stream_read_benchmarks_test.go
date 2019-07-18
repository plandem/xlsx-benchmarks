package xlsx_benchmarks

import (
	"archive/zip"
	"encoding/xml"
	"github.com/plandem/ooxml"
	"github.com/plandem/ooxml/ml"
	"testing"
)

type Cell struct {
	Formula   *ml.Reserved     `xml:"f,omitempty"`
	Value     string           `xml:"v,omitempty"`
	InlineStr *ml.Reserved     `xml:"is,omitempty"`
	ExtLst    *ml.Reserved     `xml:"extLst,omitempty"`
	Ref       string           `xml:"r,attr"`
	Style     int              `xml:"s,attr,omitempty"`
	Type      string           `xml:"t,attr,omitempty"`
	Cm        ml.OptionalIndex `xml:"cm,attr,omitempty"`
	Vm        ml.OptionalIndex `xml:"vm,attr,omitempty"`
	Ph        bool             `xml:"ph,attr,omitempty"`
}

type Row struct {
	Cells        []*Cell      `xml:"c"`
	ExtLst       *ml.Reserved `xml:"extLst,omitempty"`
	Ref          int          `xml:"r,attr,omitempty"`
	Spans        string       `xml:"spans,attr,omitempty"`
	Style        int          `xml:"s,attr,omitempty"`
	CustomFormat bool         `xml:"customFormat,attr,omitempty"`
	Height       float32      `xml:"ht,attr,omitempty"`
	Hidden       bool         `xml:"hidden,attr,omitempty"`
	CustomHeight bool         `xml:"customHeight,attr,omitempty"`
	OutlineLevel uint8        `xml:"outlineLevel,attr,omitempty"`
	Collapsed    bool         `xml:"collapsed,attr,omitempty"`
	ThickTop     bool         `xml:"thickTop,attr,omitempty"`
	ThickBot     bool         `xml:"thickBot,attr,omitempty"`
	Phonetic     bool         `xml:"ph,attr,omitempty"`
}

func openSheet(cb func(f *zip.File)) {
	z, err := zip.OpenReader(bigFile)
	if err != nil {
		panic(err)
	}

	defer z.Close()

	for _, f := range z.File {
		if f.Name == `xl/worksheets/sheet1.xml` {
			cb(f)
			break
		}
	}
}

func BenchmarkStreamRead_Stream(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		openSheet(func(f *zip.File) {
			stream, _ := ooxml.NewStreamFileReader(f)

			for {
				t, _ := stream.Token()
				if t == nil {
					break
				}

				switch start := t.(type) {
				case xml.StartElement:
					if start.Name.Local == "row" {
						var row Row
						if stream.DecodeElement(&row, &start) == nil {
							_ = row
						}
					}
				}
			}

			_ = stream.Close()
		})
	}
}

func BenchmarkStreamRead_Unmarshal(b *testing.B) {
	b.ReportAllocs()

	type Sheet struct {
		Rows []*Row `xml:"sheetData>row"`
	}

	for i := 0; i < b.N; i++ {
		openSheet(func(f *zip.File) {
			if reader, err := f.Open(); err != nil {
				panic(err)
			} else {
				var sheet Sheet

				decoder := xml.NewDecoder(reader)
				if err = decoder.Decode(&sheet); err != nil {
					panic(err)
				}

				_ = sheet
				_ = reader.Close()
			}
		})
	}
}
