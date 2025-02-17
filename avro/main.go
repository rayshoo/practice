package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/linkedin/goavro/v2"
)

func main() {
	avro1()
	avro2()
	avro3()
}

func avro1() {
	/*
		type LongList struct {
			next *LongList
		}
	*/
	codec, err := goavro.NewCodec(`
        {
          "type": "record",
          "name": "LongList",
          "fields" : [
            {"name": "next", "type": ["null", "LongList"], "default": null}
          ]
        }`)
	if err != nil {
		fmt.Println(err)
	}

	// NOTE: May omit fields when using default value
	textual := []byte(`{"next":{"LongList":{}}}`)

	// TEST
	// textual := []byte(`{"next":{"LongList":{"next": null}}}`)
	// textual := []byte(`{"next":{"LongList":{"next": {"LongList":{}}}}}`)
	// textual := []byte(`{"next":null}`)

	// Convert textual Avro data (in Avro JSON format) to native Go form
	native, _, err := codec.NativeFromTextual(textual)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(native)

	// Convert native Go form to binary Avro data
	binary, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		fmt.Println(err)
	}

	// Convert binary Avro data back to native Go form
	native, _, err = codec.NativeFromBinary(binary)
	if err != nil {
		fmt.Println(err)
	}

	// Convert native Go form to textual Avro data
	textual, err = codec.TextualFromNative(nil, native)
	if err != nil {
		fmt.Println(err)
	}

	// NOTE: Textual encoding will show all fields, even those with values that
	// match their default values
	fmt.Println(string(textual))
	// Output: {"next":{"LongList":{"next":null}}}
}

func avro2() {
	/*
		type test_schema struct {
			time int64
			customer string
		}
	*/
	avroSchema := `
	{
	  "type": "record",
	  "name": "test_schema",
	  "fields": [
		{
		  "name": "time",
		  "type": "long"
		},
		{
		  "name": "customer",
		  "type": "string"
		}
	  ]
	}`

	// Writing OCF data
	var ocfFileContents bytes.Buffer
	writer, err := goavro.NewOCFWriter(goavro.OCFConfig{
		W:      &ocfFileContents,
		Schema: avroSchema,
	})
	if err != nil {
		fmt.Println(err)
	}
	err = writer.Append([]map[string]interface{}{
		{
			"time":     1617104831727,
			"customer": "customer1",
		},
		{
			"time":     1717104831727,
			"customer": "customer2",
		},
	})
	// fmt.Println(ocfFileContents.String())

	// Reading OCF data
	ocfReader, err := goavro.NewOCFReader(strings.NewReader(ocfFileContents.String()))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Records in OCF File")
	for ocfReader.Scan() {
		record, err := ocfReader.Read()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("record", record)
	}
	err = os.WriteFile("data.avro", ocfFileContents.Bytes(), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
}

func avro3() {
	const Max = 10
	/*
		type Example struct {
			username string
			age int
		}
	*/
	// schemaJSON := `{
	//    "type": "record",
	//    "name": "Example",
	//    "fields": [
	//        {"name": "username", "type": "string"},
	//        {"name": "age", "type": "int"}
	//    ]
	// }`
	schemaJSON, err := os.ReadFile("schema.avsc")
	if err != nil {
		panic(err)
	}

	codec, err := goavro.NewCodec(string(schemaJSON))
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	data := map[string]interface{}{
		"username": "JohnDoe",
		"age":      30,
	}

	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	native, _, err := codec.NativeFromTextual(b)
	if err != nil {
		panic(err)
	}

	binaryData, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(binaryData))

	ch := make(chan []byte)
	done := make(chan bool)
	defer close(ch)
	defer close(done)

	// Produce
	go func() {
		count := 0
		for range ticker.C {
			ch <- binaryData
			if count++; count == Max {
				done <- true
				return
			}
		}
	}()

	// Consumes
	for {
		select {
		case data := <-ch:
			native, _, err := codec.NativeFromBinary(data)
			if err != nil {
				panic(err)
			}
			textual, err := codec.TextualFromNative(nil, native)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(textual))

		case result := <-done:
			if result {
				return
			}
		}
	}
}
