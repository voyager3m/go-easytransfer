// EasyTransfer project go-easytransfer.go
/******************************************************
This program is free software: you can redistribute it and/or
modify it under the terms of the GNU General Public License as
published by the Free Software Foundation, either version 3 of
the License, or(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
<http://www.gnu.org/licenses/>

This work is licensed under the Creative Commons Attribution-ShareAlike 3.0 Unported License.
To view a copy of this license, visit http://creativecommons.org/licenses/by-sa/3.0/ or
send a letter to Creative Commons, 444 Castro Street, Suite 900, Mountain View, California, 94041, USA.
********************************************************/
package easytransfer

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
)

var inbuf bytes.Buffer

func SendData(val interface{}, stream io.Writer) {
	header := []byte{0x06, 0x85}

	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, val)
	CS := byte(buf.Len())

	var outbuf bytes.Buffer
	outbuf.Write(header)
	outbuf.WriteByte(CS)
	outbytes := buf.Bytes()
	for _, j := range outbytes {
		CS ^= j
		outbuf.WriteByte(j)
	}
	outbuf.WriteByte(CS)
	stream.Write(outbuf.Bytes())
}

func readToBuffer(stream io.Reader, count int) {
	buf := make([]byte, 1)
	var readed = 0
	for {
		r, err := stream.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error read from stream: %v", err)
		}
		if r > 0 {
			readed++
			inbuf.WriteByte(buf[0])
		} else {
			break
		}
		if readed == count {
			break
		}
	}
}

func ReceiveData(val interface{}, stream io.Reader) bool {
	readToBuffer(stream, 4)
	if inbuf.Len() < 4 {
		log.Printf("inbuf len < 4: %v\n", inbuf.Len())
		return false
	}
	for {
		b, err := inbuf.ReadByte()
		if err != nil {
			log.Printf("read inbuf error: %v", err)
			return false
		}
		if b == 0x06 {
			break
		}
		readToBuffer(stream, 1)
		if inbuf.Len() == 0 {
			log.Printf("start header not found")
			return false
		}
	}

	b, err := inbuf.ReadByte()
	if err != nil {
		log.Printf("read inbuf error: %v\n", err)
		return false
	}
	if b != 0x85 {
		log.Printf("b != 0x85 (%d)\n", b)
		return false
	}

	if b, err = inbuf.ReadByte(); err != nil {
		log.Printf("error read buffer: %v\n", err)
		return false
	}

	// get size of struct
	var size = byte(0)
	{
		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, val)
		size = byte(buf.Len())
	}

	if size != b {
		log.Printf("struct sizes are diff %d != %d\n", size, b)
		return false
	}

	readToBuffer(stream, int(size))

	var buf bytes.Buffer
	var cs = size
	for i := 0; i < int(size); i++ {
		b, _ := inbuf.ReadByte()
		cs ^= b
		buf.WriteByte(b)
	}

	if b, err = inbuf.ReadByte(); err != nil {
		log.Printf("error read buffer: %v\n", err)
		return false
	}

	if b != cs {
		log.Printf("bad CRC: %d != %d", b, cs)
		return false
	}

	if err = binary.Read(&buf, binary.LittleEndian, val); err != nil {
		log.Printf("Error write to interface: %v\n", err)
		return false

	}

	return true
}
