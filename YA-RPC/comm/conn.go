package comm

import (
	"net"
)

func Write(conn net.Conn, pkg *RPCPackage) error {
	b, err := pkg.Encode()
	if err != nil {
		return err
	}
	_, err = conn.Write(b)
	return err
}

func Read(conn net.Conn) (*RPCPackage, error) {
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	b := buf[:n]
	pkg, err := Decode(b)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}
