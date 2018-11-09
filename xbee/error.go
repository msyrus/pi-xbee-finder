package xbee

import "errors"

var (
	ErrInvalidData       = errors.New("data is not a valid xbee data")
	ErrChkSumMismatched  = errors.New("checksum does not match")
	ErrUnsupportedLength = errors.New("length of packet is not supported")
)
