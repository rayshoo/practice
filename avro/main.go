package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/linkedin/goavro"
)

func main() {
	// Avro 스키마 정의
	schemaJSON := `{
	   "type": "record",
	   "name": "Example",
	   "fields": [
	       {"name": "username", "type": "string"},
	       {"name": "age", "type": "int"}
	   ]
	}`
	// schemaJSON, err := os.ReadFile("schema.avsc")
	// if err != nil {
	// 	panic(err)
	// }

	codec, err := goavro.NewCodec(string(schemaJSON))
	if err != nil {
		panic(err)
	}

	// 파일 생성
	file, err := os.Create("data.avro")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Avro 데이터 생성 및 파일에 적재
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 예시 데이터 생성
		data := map[string]interface{}{
			"username": "JohnDoe",
			"age":      30,
		}

		b, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))

		// Avro 데이터 생성
		native, _, err := codec.NativeFromTextual(b)
		if err != nil {
			panic(err)
		}

		binaryData, err := codec.BinaryFromNative(nil, native)
		if err != nil {
			panic(err)
		}

		// 파일에 Avro 데이터 쓰기
		if _, err := file.Write(binaryData); err != nil {
			panic(err)
		}

		// 줄 바꿈
		if _, err := file.WriteString("\n"); err != nil {
			panic(err)
		}

		fmt.Println("Avro 데이터를 파일에 적재했습니다.")
	}
}
