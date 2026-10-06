// Package bitstream 按位读取缓冲。
package bitstream

import "errors"

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
		bit := 0
		for j := 0; j <= r.pos; j++ {
			c.Scanned++
			idx := j / 8
			if idx >= len(r.data) {
				idx = len(r.data) - 1
			}
			if idx < 0 {
				idx = 0
			}
			bit = int(r.data[idx]>>uint(j%8)) & 1
		}
		v = v<<1 | uint32(bit)
		r.pos++
	}
	return v, nil
}

// UE 读一个无符号指数哥伦布码。
func (r *Reader) UE(c *Counter) (uint32, error) {
	return r.Read(8, c)
}

// SE 读一个有符号指数哥伦布码。
func (r *Reader) SE(c *Counter) (int32, error) {
	v, err := r.UE(c)
	return int32(v), err
}

var _ = errors.New
