package xbee

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
)

type FrameType int

const (
	UnknownFrameType FrameType = iota
	DataSampleFrameType
)

type PacketType int

const (
	UnknownPacketType PacketType = iota
	UnicastPacketType
	BroadcastPacketType
)

// 7E 0012 92 0013A20040A9C7E5 B271 01 01 0002 00 0000 FC

// 7E 0012 92 0013A20040A9C7D4 BC24 02 01 0002 00 0002 4D
//
// start bit - 7E
// byte len before chksm - 0012 -> 18
// frame type - 92 (data sample frame)
// device id - 0013A20040A9C7D4
// net adrs - BC24
// packet type - 02 (broadcast packet)
// number of samples - 01 (fixed)
// digital chan mask - 0002 -> 0000000000000010
// analog chan mask - 00
// value - 0002
// chksm - 4D

type Frame struct {
	Len        int
	Type       FrameType
	DeviceID   string
	NetAddr    string
	PacketType PacketType
	SampleN    int
	DataD      map[int]bool
	DataA      map[int]int
}

func ParseFrom(r io.Reader, dataCh chan<- *Frame, errCh chan<- error) error {
	buf := &bytes.Buffer{}
	for {
		byts := make([]byte, 1000)
		n, err := r.Read(byts)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		buf.Write(byts[:n])
		if buf.Len() < 20 {
			continue
		}

		tmpBuf := &bytes.Buffer{}
		for {
			b, err := buf.ReadByte()
			if err != nil {
				break
			}
			if b != 0x7e {
				continue
			}
			tmpBuf.WriteByte(b)

			if buf.Len() < 2 {
				break
			}

			byts := buf.Next(2)
			tmpBuf.Write(byts)

			l := int(binary.BigEndian.Uint16(byts))
			if buf.Len() < (l + 1) {
				break
			}

			tmpBuf.Write(buf.Next(l + 1))

			f, err := ParseFrame(tmpBuf.Next(tmpBuf.Len()))
			if err != nil {
				errCh <- err
				continue
			}
			dataCh <- f
		}
		buf.WriteTo(tmpBuf)
		tmpBuf.WriteTo(buf)
	}
}

func ParseFrame(data []byte) (*Frame, error) {
	// fmt.Println(hex.EncodeToString(data))
	if len(data) < 4 || data[0] != 0x7e {
		return nil, ErrInvalidData
	}

	if !checksum(data) {
		return nil, ErrChkSumMismatched
	}

	l := int(binary.BigEndian.Uint16(data[1:3]))
	if l != 18 {
		return nil, ErrUnsupportedLength
	}

	f := Frame{
		Len:      l,
		DeviceID: hex.EncodeToString(data[4:12]),
		NetAddr:  hex.EncodeToString(data[12:14]),
		SampleN:  int(data[15]),
	}

	switch data[3] {
	case 0x92:
		f.Type = DataSampleFrameType
	}

	switch data[14] {
	case 0x01:
		f.PacketType = UnicastPacketType
	case 0x02:
		f.PacketType = BroadcastPacketType
	}

	dmask := binary.BigEndian.Uint16(data[16:18])
	ddata := binary.BigEndian.Uint16(data[19:21])
	amask := data[18]

	f.DataD = map[int]bool{}
	f.DataA = map[int]int{}

	for i := 0; dmask != 0; i++ {
		if dmask&1 == 1 {
			f.DataD[i] = (ddata & 1) == 1
		}
		dmask = dmask >> 1
		ddata = ddata >> 1
	}

	for i := 0; amask != 0; i++ {
		if amask&1 == 1 {
			f.DataA[i] = int(binary.BigEndian.Uint16(data[19:21]))
		}
		amask = amask >> 1
	}

	return &f, nil
}

func checksum(data []byte) bool {
	sum := byte(0)
	for _, d := range data[3:] {
		sum = sum + d
	}
	return sum == 0xff
}
