package go_ntskem

import (
	"testing"
)

func TestGenerateKey(t *testing.T) {
	nts := NTSKEM{}
	nts.New(12)
	nts.GenerateKey()
}

func TestEncapsulate(t *testing.T) {

}

func TestDecapsulate(t *testing.T) {

}
