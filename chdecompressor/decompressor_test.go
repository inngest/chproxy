package chdecompressor

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestRead(t *testing.T) {
	testCases := []struct {
		lz4      string
		zstd     string
		expected string
	}{
		{
			"\xfe\xf7\xd3-\xd9%\b\xb8\xeaz\xef\xe4Zt\xe8E\x823\x00\x00\x00+\x00\x00\x00\xf2\vSELECT number FROM system.\x13\x00\xb0s LIMIT 10\n",
			"\x93:6;\x93aw%\xa8$!\xb0Z\x8fG*\x90=\x00\x00\x00+\x00\x00\x00(\xb5/\xfd +Y\x01\x00SELECT number FROM system.numbers LIMIT 10\n",
			"SELECT number FROM system.numbers LIMIT 10\n",
		},
		{
			"\xbc\x83\xfb\x856BM\x82\xcdç\xa35\xd7\x18\xa7\x82(\x00\x00\x00\x1d\x00\x00\x00\xf0\x0eSELECT * FROM system.metrics\n",
			"\xd4Hq6t<d\x18\xa7\xf44$\xc7h\x89\x80\x90/\x00\x00\x00\x1d\x00\x00\x00(\xb5/\xfd \x1d\xe9\x00\x00SELECT * FROM system.metrics\n",
			"SELECT * FROM system.metrics\n",
		},
		{
			"`\xf7d\xb6t\x9bWM\x15|{\x02\xc5s\xd0 \x82e\x00\x00\x00z\x00\x00\x00\xc0SELECT col1,\x06\x00\x112\x06\x00\x113\x06\x00\x114\x06\x00\x115\x06\x00\x116\x06\x00\x117\x06\x00\x118\x06\x00\x119\x06\x00\"10\a\x00\x02>\x00\x121?\x00\x121@\x00\x121A\x00\xf0\b15 FROM system.metrics\n",
			"\xdb\x03\xea\x18\xbae\xcb٤\x92\xf0\\<\x0e\xae\xac\x90q\x00\x00\x00z\x00\x00\x00(\xb5/\xfd z\xfd\x02\x00\xe2\x04\x12\x1a\x90\xc5\t\xbbR\x98BX[\xce\x19U\x8a\xd3\"\x9b\xec\x16\xd9l\x1bj\u007f(p\x01\x83\xfb\x80\xf0X#\x1a{\xae\xe0\xec<\x00\xc6*\xa1\xccԽ\xb1H%A\x0e\xdd\x1b\x95\xee\x8dI\xf7F\xa4{\xa3\xa1{c\xeaނe5+\x10\b\x00`\b\x9c\x81\x02\x0e\x16\xe0\x89\xe0\x04\xb4\x03²\xd7\xd4\x03",
			"SELECT col1, col2, col3, col4, col5, col6, col7, col8, col9, col10, col11, col12, col13, col14, col15 FROM system.metrics\n",
		},
		{
			"v\xad_\xfcX\x156C\xf6\xfb\xc9'\xb7\x8b\x11H\x82\v\x00\x00\x00\x01\x00\x00\x00\x10\n",
			"\x85C-a\xcb\xde\xd0i\xbb\v\xfa\xd9_\n\u007f\b\x90\x13\x00\x00\x00\x01\x00\x00\x00(\xb5/\xfd \x01\t\x00\x00\n",
			"\n",
		},
	}

	for _, tc := range testCases {
		testDecompress(t, tc.lz4, tc.expected)
		// testDecompress(t, tc.zstd, tc.expected)
	}
}

func testDecompress(t *testing.T, str, expected string) {
	bb := bytes.NewBufferString(str)
	r := NewReader(bb)
	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("cannot decompress %q: %s", expected, err)
	}
	if string(b) != expected {
		t.Fatalf("got %q; expected: %q", string(b), expected)
	}
}

func TestReadNegative(t *testing.T) {
	testCases := []string{
		"",
		"\xfe",
		"\x93:6;",
		"\xfe\xf7\xd3-\xd9%\b\xb8\xeaz\xef\xe4Zt\xe8E\x823\x00\x00\x00+\x00",
		"\x93:6;\x93aw%\xa8$!\xb0Z\x8fG*\x90=\x00\x00\x00+\x00\x00\x00(\xb5/\xfd +Y\x01",
		"\xfe\xf7\xd3-\xd9%\b\xb8\xeaz\xef\xe4Zt\xe8E\x823\x00\x00\x00+\x00\x00\x00\xf2\vSELECT number FROM system.\x13\x00\xb0s",
		"\x93:6;\x93aw%\xa8$!\xb0Z\x8fG*\x90=\x00\x00\x00+\x00\x00\x00(\xb5/\xfd +Y\x01\x00SELECT number FROM system.",
	}
	for _, v := range testCases {
		bb := bytes.NewBufferString(v)
		r := NewReader(bb)
		b := make([]byte, 128)
		_, err := r.Read(b)
		if err == nil {
			t.Fatalf("expected to get error for %q; got nil", v)
		}
	}
}
