package main

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func encodingExercise() {
	_ = base64.StdEncoding.EncodeToString([]byte("hello"))
	_, _ = base64.StdEncoding.DecodeString("aGVsbG8=")
	_ = base64.URLEncoding.EncodeToString([]byte("hello"))

	_ = hex.EncodeToString([]byte("hello"))
	_, _ = hex.DecodeString("68656c6c6f")
	_ = hex.NewEncoder(nil)
	_ = hex.NewDecoder(nil)

	type Person struct {
		Name string `xml:"name"`
		Age  int    `xml:"age"`
	}
	_, _ = xml.Marshal(&Person{Name: "test", Age: 30})
	_ = xml.Unmarshal([]byte("<Person><name>test</name></Person>"), &Person{})
	_ = xml.NewEncoder(nil)
	_ = xml.NewDecoder(nil)

	var buf [8]byte
	_ = binary.Read(nil, binary.LittleEndian, &buf)
	_ = binary.Write(nil, binary.LittleEndian, buf[:])
	_ = binary.Size(buf)
}

func regexpExercise() {
	re := regexp.MustCompile(`\w+`)
	_ = re.MatchString("hello")
	_ = re.FindString("hello world")
	_ = re.FindAllString("hello world", -1)
	_ = re.ReplaceAllString("hello world", "x")
	_ = re.Split("hello,world", -1)
}

func filepathExercise() {
	_ = filepath.Join("a", "b", "c")
	_ = filepath.Base("/etc/hosts")
	_ = filepath.Dir("/etc/hosts")
	_ = filepath.Ext("file.txt")
	_, _ = filepath.Abs(".")
	_ = filepath.Clean("/etc/../etc/hosts")
	_, _ = filepath.Rel("/etc", "/etc/hosts")
	_, _ = filepath.Split("/etc/hosts")
	_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error { return nil })
	_, _ = filepath.Glob("*.go")
	_, _ = filepath.Match("*.go", "main.go")
	_ = filepath.IsAbs("/etc/hosts")
}

func main() {
	encodingExercise()
	regexpExercise()
	filepathExercise()
	fmt.Println("encoding e2e test done")
}
