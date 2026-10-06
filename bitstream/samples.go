package bitstream

import "strings"

// Scenario 是一个固定场景。
type Scenario struct {
	Name string
	Run  func() map[string]any
}

// EncUE 编码一个无符号指数哥伦布码的位串。
func EncUE(v uint32) string {
	x := v + 1
	bits := ""
	for x > 0 {
		bits = string(rune('0'+x%2)) + bits
		x /= 2
	}
	return strings.Repeat("0", len(bits)-1) + bits
}

// EncSE 编码一个有符号指数哥伦布码的位串。
func EncSE(v int32) string {
	if v <= 0 {
		return EncUE(uint32(-v) * 2)
	}
	return EncUE(uint32(v)*2 - 1)
}

// Pack 把位串补齐到字节。
func Pack(bits string) []byte {
	for len(bits)%8 != 0 {
		bits += "0"
	}
	out := []byte{}
	for i := 0; i < len(bits); i += 8 {
		var b byte
		for j := 0; j < 8; j++ {
			b <<= 1
			if bits[i+j] == '1' {
				b |= 1
			}
		}
		out = append(out, b)
	}
	return out
}

// Samples 返回全部场景。
func Samples() []Scenario {
	return []Scenario{
		{Name: "bitorder", Run: func() map[string]any {
			v, err := NewReader(Pack("1011")).Read(4, &Counter{})
			if err != nil {
				return map[string]any{"error": err.Error()}
			}
			return map[string]any{"value": v}
		}},
		{Name: "ue", Run: func() map[string]any {
			v, err := NewReader(Pack(EncUE(5))).UE(&Counter{})
			if err != nil {
				return map[string]any{"error": err.Error()}
			}
			return map[string]any{"value": v}
		}},
		{Name: "se", Run: func() map[string]any {
			v, err := NewReader(Pack(EncSE(-2))).SE(&Counter{})
			if err != nil {
				return map[string]any{"error": err.Error()}
			}
			return map[string]any{"value": v}
		}},
		{Name: "bounds", Run: func() map[string]any {
			v, err := NewReader([]byte{0x01}).Read(100, &Counter{})
			if err != nil {
				return map[string]any{"error": err.Error()}
			}
			return map[string]any{"value": v}
		}},
	}
}

// Work 跑规模线场景。
func Work(n int) map[string]any {
	data := make([]byte, n)
	for i := range data {
		data[i] = byte(i)
	}
	r := NewReader(data)
	c := &Counter{}
	for i := 0; i < n; i++ {
		r.Read(8, c)
	}
	return map[string]any{"reads": n, "scanned": c.Scanned}
}
