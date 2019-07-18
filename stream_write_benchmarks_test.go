package xlsx_benchmarks

import (
	"bytes"
	"encoding/xml"
	"github.com/plandem/ooxml"
	"math/rand"
	"testing"
	"time"
)

const streamWriteRows = 2000000

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func writeData(cb func(row *Row)) {
	for iRow := 0; iRow < streamWriteRows; iRow++ {
		row := Row{
			Ref:    iRow,
			Height: rand.Float32(),
		}

		cb(&row)
	}
}

func BenchmarkStreamWrite_Memory(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		stream, err := ooxml.NewStreamFileWriter("test.xlsx", true, nil)
		if err != nil {
			panic(err)
		}

		writeData(func(row *Row) {
			if err := stream.Encode(row); err != nil {
				panic(err)
			}
		})

		_ = stream.Close()
	}
}

func BenchmarkStreamWrite_File(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		stream, err := ooxml.NewStreamFileWriter("test.xlsx", false, nil)
		if err != nil {
			panic(err)
		}

		writeData(func(row *Row) {
			if err := stream.Encode(row); err != nil {
				panic(err)
			}
		})

		_ = stream.Close()
	}
}

func BenchmarkStreamWrite_Marshal(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		items := make([]*Row, 0)
		writeData(func(row *Row) {
			items = append(items, row)
		})

		buf := &bytes.Buffer{}
		enc := xml.NewEncoder(buf)
		if err := enc.Encode(items); err != nil {
			panic(err)
		}
	}
}
