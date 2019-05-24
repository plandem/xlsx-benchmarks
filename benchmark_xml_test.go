package xlsx_benchmarks_test

import (
	"encoding/xml"
	"github.com/plandem/xlsx/internal/ml"
	"log"
	"os"
	"testing"
)

const benchmarkXMLFile = "./test_files/saved_huge_xlsx/xl/worksheets/sheet1.xml"

func xmlReadUnmarshal() {
	f, err := os.Open(benchmarkXMLFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ws := ml.Worksheet{}
	decoder := xml.NewDecoder(f)
	decoder.Decode(&ws)
}

func xmlReadSax() {
	f, err := os.Open(benchmarkXMLFile)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	decoder := xml.NewDecoder(f)

	var row = ml.Row{}

	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "row":
				row = ml.Row{}
				decoder.DecodeElement(&row, &se)
			}
		}
	}

	for _, c := range row.Cells {
		log.Printf("%+v,", c.Value)
	}
}

func BenchmarkXML(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		xmlReadUnmarshal()
	}
}

//
//func BenchmarkXML(b *testing.B) {
//	b.ReportAllocs()
//
//	for i := 0; i < b.N; i++ {
//		xmlReadSax()
//	}
//}

//func BenchmarkXML_Unmarshal(b *testing.B) {
//	b.ReportAllocs()
//
//	for i := 0; i < b.N; i++ {
//		xmlReadUnmarshal()
//	}
//}
//
//func BenchmarkXML_SAX(b *testing.B) {
//	b.ReportAllocs()
//
//	for i := 0; i < b.N; i++ {
//		xmlReadSax()
//	}
//}
