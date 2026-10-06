package bitstream

import "testing"

func TestBitsAllOnes(t *testing.T) {
	v, err := NewReader([]byte{0xff}).Read(8, &Counter{})
	if err != nil || v != 255 {
		t.Fatalf("全 1 应读成 255：%d %v", v, err)
	}
}

func TestBitsAllZeroes(t *testing.T) {
	v, err := NewReader([]byte{0x00}).Read(8, &Counter{})
	if err != nil || v != 0 {
		t.Fatalf("全 0 应读成 0：%d %v", v, err)
	}
}

func TestZeroWidth(t *testing.T) {
	v, err := NewReader([]byte{0xff}).Read(0, &Counter{})
	if err != nil || v != 0 {
		t.Fatalf("读 0 位应是 0：%d %v", v, err)
	}
}

func TestZeroWidthEmpty(t *testing.T) {
	v, err := NewReader(nil).Read(0, &Counter{})
	if err != nil || v != 0 {
		t.Fatalf("空缓冲读 0 位应是 0：%d %v", v, err)
	}
}

func TestNewReaderNotNil(t *testing.T) {
	if NewReader(nil) == nil {
		t.Fatal("应返回读取器")
	}
}

func TestCounterCounts(t *testing.T) {
	c := &Counter{}
	NewReader([]byte{0x00}).Read(1, c)
	if c.Scanned == 0 {
		t.Fatal("应统计比较次数")
	}
}
