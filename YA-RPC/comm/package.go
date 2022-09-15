package comm

import (
	"bytes"
	"encoding/gob"
)

type RPCPackage struct {
	FuncName string
	Arg      interface{}
	Err      error
}

func NewRPCPackage(funcName string, arg interface{}, err error) *RPCPackage {
	pkg := &RPCPackage{
		FuncName: funcName,
		Arg:      arg,
		Err:      err,
	}
	return pkg
}

func (p *RPCPackage) Encode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *RPCPackage) Decode(serial []byte) error {
	buf := bytes.NewBuffer(serial)
	decoder := gob.NewDecoder(buf)
	var pkg RPCPackage
	if err := decoder.Decode(&pkg); err != nil {
		return err
	}
	p.FuncName = pkg.FuncName
	p.Arg = pkg.Arg
	p.Err = pkg.Err
	return nil
}

func Encode(p *RPCPackage) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(serial []byte) (*RPCPackage, error) {
	buf := bytes.NewBuffer(serial)
	decoder := gob.NewDecoder(buf)
	var pkg RPCPackage
	if err := decoder.Decode(&pkg); err != nil {
		return nil, err
	}
	return &pkg, nil
}
