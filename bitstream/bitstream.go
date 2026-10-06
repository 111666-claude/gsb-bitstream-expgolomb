// Package bitstream 按位读取缓冲。
package bitstream

import (
	"errors"
	"fmt"
)

// Counter 记录比较次数。
type Counter struct {
	Scanned int
}

// ErrEOF 在请求的位超出缓冲时返回。
var ErrEOF = errors.New("eof")

// Reader 是位读取器。
type Reader struct {
	data []byte
	pos  int
}

// NewReader 建一个位读取器。
func NewReader(data []byte) *Reader {
	return &Reader{data: data}
}

// bit 读游标处的一位并前进游标。
func (r *Reader) bit(c *Counter) (uint32, error) {
	c.Scanned++
	if r.pos >= len(r.data)*8 {
		return 0, ErrEOF
	}
	b := uint32(r.data[r.pos/8]>>uint(7-r.pos%8)) & 1
	r.pos++
	return b, nil
}

// Read 读 n 位，从每个字节的最高位开始。
func (r *Reader) Read(n int, c *Counter) (uint32, error) {
	var v uint32
	for i := 0; i < n; i++ {
		b, err := r.bit(c)
		if err != nil {
			return 0, err
		}
		v = v<<1 | b
	}
	return v, nil
}

// UE 读一个无符号指数哥伦布码。
func (r *Reader) UE(c *Counter) (uint32, error) {
	z := 0
	for {
		b, err := r.bit(c)
		if err != nil {
			return 0, err
		}
		if b == 1 {
			break
		}
		z++
		if z >= 32 {
			return 0, fmt.Errorf("ue: 前导零过多：%w", ErrEOF)
		}
	}
	b, err := r.Read(z, c)
	if err != nil {
		return 0, err
	}
	return 1<<uint(z) - 1 + b, nil
}

// SE 读一个有符号指数哥伦布码。
func (r *Reader) SE(c *Counter) (int32, error) {
	k, err := r.UE(c)
	if err != nil {
		return 0, err
	}
	if k%2 == 0 {
		return -int32(k / 2), nil
	}
	return int32((k + 1) / 2), nil
}
