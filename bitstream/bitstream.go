// Package bitstream 按位读取缓冲。
package bitstream

import "errors"

// ErrEOF 表示读取越过缓冲末尾。
var ErrEOF = errors.New("eof")

// Counter 记录比较次数。
type Counter struct {
	Scanned int
}

// Reader 是位读取器。
type Reader struct {
	data []byte
	pos  int
}

// NewReader 建一个位读取器。
func NewReader(data []byte) *Reader {
	return &Reader{data: data}
}

// Read 读 n 位。
func (r *Reader) Read(n int, c *Counter) (uint32, error) {
	var v uint32
	for i := 0; i < n; i++ {
		idx := r.pos / 8
		if idx >= len(r.data) {
			return v, ErrEOF
		}
		if c != nil {
			c.Scanned++
		}
		bit := int(r.data[idx]>>(7-uint(r.pos%8))) & 1
		v = v<<1 | uint32(bit)
		r.pos++
	}
	return v, nil
}

// UE 读一个无符号指数哥伦布码。
func (r *Reader) UE(c *Counter) (uint32, error) {
	z := 0
	for {
		bit, err := r.Read(1, c)
		if err != nil {
			return 0, err
		}
		if bit == 1 {
			break
		}
		z++
	}
	b, err := r.Read(z, c)
	if err != nil {
		return 0, err
	}
	return (uint32(1) << uint(z)) - 1 + b, nil
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
