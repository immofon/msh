package msh

import (
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	Read(strings.NewReader("$MeshFormat\n4.1 0 8\n$EndMeshFormat\n"))
}
