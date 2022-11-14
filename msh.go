package msh

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Mesh struct {
	Format        MeshFormat
	PhysicalNames []PhysicalName
}

type MeshFileType int

const (
	MeshFileType_ASCII  MeshFileType = 0
	MeshFileType_Binary MeshFileType = 1
)

type MeshFormat struct {
	Version  string       // only support "4.1"
	FileType MeshFileType // only support ASCII mode, 0
	DataSize int          // sizeof(size_t)
}

type PhysicalName struct {
	Dimension   int
	PhysicalTag int
	Name        string
}

func ReadLine(r io.Reader) (line <-chan string, errch <-chan error) {
	line_ := make(chan string, 100)
	errch_ := make(chan error, 1)

	go func(line chan<- string, errch chan<- error) {
		defer close(line)
		defer close(errch)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line <- scanner.Text()
		}

		errch <- scanner.Err()
	}(line_, errch_)
	return line_, errch_
}

func ExpectFieldName(linech <-chan string) string {
	for line := range linech {
		if strings.HasPrefix(line, "$") && !strings.HasPrefix(line, "$End") {
			return line[1:]
		}
	}
	return ""
}

func ExpectEndFieldName(linech <-chan string, field_name string) {
	v := "$End" + field_name
	for line := range linech {
		if line == v {
			return
		}
	}
}

func ExpectMeshFormat(linech <-chan string) MeshFormat {
	format := MeshFormat{}
	args := strings.Split(<-linech, " ")
	if len(args) != 3 {
		panic("ExpectMeshFormat: len(args) MUST be 3!")
	}
	format.Version = args[0]
	file_type, err := strconv.Atoi(args[1])
	if err != nil {
		panic(err)
	}
	data_size, err := strconv.Atoi(args[2])
	if err != nil {
		panic(err)
	}

	format.FileType = MeshFileType(file_type)
	format.DataSize = (data_size)
	ExpectEndFieldName(linech, "MeshFormat")
	return format
}

func Read(r io.Reader) (*Mesh, error) {
	linech, errch := ReadLine(r)
	mesh := &Mesh{}

LOOP:
	for {
		field_name := ExpectFieldName(linech)
		fmt.Println("Field Name:", field_name)
		switch field_name {
		case "": // End
			break LOOP
		case "MeshFormat":
			mesh.Format = ExpectMeshFormat(linech)
			fmt.Println(mesh.Format)
		default: // Unsupported field name
			fmt.Println("!Unsupported field name:", field_name)
		}
	}

	return &Mesh{}, <-errch
}
