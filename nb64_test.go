package nb64

import (
	"bytes"
	"encoding/base64"
	"math"
	"math/rand"
	"testing"
	"time"
	"unicode/utf8"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var max = int(math.Pow(2, 18))
var debugEncoding *Encoding

func init() {
	var debugEncode [64]byte
	for i := 0; i < 64; i++ {
		debugEncode[i] = byte(i)
	}
	debugEncoding = NewEncoding(string(debugEncode[:]), '~')
}

func TestNB64(t *testing.T) {
	src := []byte("on`t carry")
	// src := []byte("1𐄁1中1·1")
	enc, err := debugEncoding.Encode(src)
	if err != nil {
		t.Log(err)
		t.Logf("plain:  %s", src)
		t.Logf("encode: %s", enc)
		t.FailNow()
	}
	// fmt.Printf("%s %07b, %06b: %v\n", src, src, dst, err)
	dst, err := debugEncoding.Decode(enc)
	if err != nil || !bytes.Equal(dst, src) {
		t.Log(err)
		t.Logf("plain:  %07b %s", src, src)
		t.Logf("encode: %06b %s", enc, enc)
		t.Logf("decode: %07b %s", dst, dst)
		t.FailNow()
	}
}

func TestBoundary(t *testing.T) {
	//   Unicode符号范围    |        UTF-8编码方式
	//   (十六进制)         |          （二进制）
	// --------------------+---------------------------------------------
	// 0000 0000-0000 007F | 0xxxxxxx
	// 0000 0080-0000 07FF | 110xxxxx 10xxxxxx
	// 0000 0800-0000 FFFF | 1110xxxx 10xxxxxx 10xxxxxx
	// 0001 0000-0010 FFFF | 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
	src := []byte("\u0000\u007F\u0080\u07FF\u0800\uFFFF\U00010000\uFFFD")
	enc, err := debugEncoding.Encode(src)
	if err != nil {
		for i, v := range string(src) {
			t.Logf("%d: %d %x", i, v, v)
		}
		t.Fatalf("%07b: %v", src, err)
	}
	dst, err := debugEncoding.Decode(enc)
	if err != nil || !bytes.Equal(plain, src) {
		for i, v := range string(src) {
			t.Logf("%d: %d %x", i, v, v)
		}
		t.Fatalf("%07b %06b %07b", src, enc, dst)
	}

	src = []byte("\U0004FFFF")
	_, err = debugEncoding.Encode(src)
	if err == nil {
		for i, v := range string(src) {
			t.Logf("%d: %d %x", i, v, v)
		}
		t.Fatalf("%07b: %v", src, err)
	}
}

// var chars = [][]byte{nil, []byte("\u0001"), []byte("\u0080"), []byte("\u0800"), []byte("\U00010000")}
var chars = [][]byte{nil, []byte("$"), []byte("¢"), []byte("€"), []byte("𐄁")}

// var ln = make([]int, 48)
// var ln2 = make([]int, 96)

func recurse(t *testing.T, src []byte) {
	for i := 4; i > 0; i-- {
		src = append(src, chars[i]...)
		// fmt.Printf("%s\n", src)
		enc, err := debugEncoding.Encode(src)
		if err != nil {
			t.Fatalf("%s %08b, %06b: %v", src, src, enc, err)
		}
		dst, err := debugEncoding.Decode(enc)
		if err != nil {
			t.Fatalf("%s:%08b, %06b, %08b: %v", src, src, dst, dst, err)
		}
		if !bytes.Equal(dst, src) {
			t.Fatalf("%s %08b, %06b, %08b", src, src, dst, dst)
		}

		// if ln[len(src)] < len(dst) {
		// 	ln[len(src)] = len(dst)
		// }
		// if ln2[len(dst)] < len(src) {
		// 	ln2[len(dst)] = len(src)
		// }

		if len(src) < 20 || len(dst) < 20 {
			recurse(t, src)
		}

		src = src[:len(src)-i]
	}
}

func TestOrder(t *testing.T) {
	src := make([]byte, 0, 48)
	recurse(t, src)

	// for i, v := range ln {
	// 	fmt.Printf("%3d\t%3d\n", i, v)
	// }
	// for i, v := range ln2 {
	// 	fmt.Printf("%3d\t%3d\n", i, v)
	// }
}

func TestSelect(t *testing.T) {
	var chars = [][]byte{nil, []byte("$"), []byte("¢"), []byte("€"), []byte("𐄁")}
	for k := 0; k < 1000; k++ {
		var src []byte
		for i := 0; i < 256; i++ {
			c := chars[r.Intn(5)]
			src = append(src, c...)
		}
		dst, err := debugEncoding.Encode(src)
		if err != nil {
			t.Fatalf("%02x", src)
		}
		_, err = debugEncoding.Decode(dst)
		if err != nil {
			t.Fatalf("%02x", dst)
		}
	}
}

func TestRandRunes(t *testing.T) {
	for i := 0; i < 1000; i++ {
		testRandRunes(t)
	}
}

func testRandRunes(t *testing.T) {
	t.Helper()
	src, n := make([]byte, 4096), 0
	for i := 0; i < 128; i++ {
		cp := r.Intn(max)
		n += utf8.EncodeRune(src[n:], rune(cp))
	}
	enc, _ := URLEncoding.Encode(src[:n])
	dst, err := URLEncoding.Decode(enc)
	if err != nil || !bytes.Equal(plain, src[:n]) {
		t.Fatalf("src: %x, dst: %x, plain: %x", src[:n], enc, dst)
	}
	// t.Logf("src: %5d, dst: %x, plain: %x", src[:n], dst, plain)
}

func TestDecode(t *testing.T) {
	enc := []byte("r5UA.GIRQEBE7s.QYeWVA.ACi.QA.GYv.QYdmxA.JD9.QA.ACi.QZ9-_Ig.Fl9.Qa9OTmg.FtpQEBFtQ.XA.DAC.")
	dst, err := URLEncoding.Decode(enc)
	if err != nil {
		t.Fatalf("%s, %v", dst, err)
	}
}

func TestNB64EncodeDecode(t *testing.T) {
	src := plain[50:60]
	enc, err := URLEncoding.Encode(src)
	if err != nil {
		t.Fatalf("%s, %v", src, err)
	}
	dst, err := URLEncoding.Decode(enc)
	if err != nil {
		t.Fatalf("%s, %v", enc, err)
	}
	if !bytes.Equal(dst, src) {
		t.Logf("plain:  %s", src)
		t.Logf("encode: %s", enc)
		t.Logf("decode: %s", dst)
		t.FailNow()
	}
}
func TestBase64EncodeDecode(t *testing.T) {
	enc := make([]byte, base64.RawURLEncoding.EncodedLen(len(plain)))
	base64.RawURLEncoding.Encode(enc, plain)

	dst := make([]byte, base64.RawURLEncoding.DecodedLen(len(enc)))
	_, err := base64.RawURLEncoding.Decode(dst, enc)
	if err != nil {
		t.Fatalf("%s, %v", enc, err)
	}
	if !bytes.Equal(dst, plain) {
		t.Logf("decode:\n%s", dst)
		t.Logf("plain:\n%s", plain)
		t.FailNow()
	}
}

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		URLEncoding.Encode(plain)
	}
}

func BenchmarkBase64Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dst := make([]byte, (len(plain)*4+3)/3)
		base64.RawURLEncoding.Encode(dst, plain)
	}
}

func BenchmarkDecode(b *testing.B) {
	enc, err := URLEncoding.Encode(plain)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		URLEncoding.Decode(enc)
	}
}

func BenchmarkBase64Decode(b *testing.B) {
	enc := make([]byte, (len(plain)*4+3)/3)
	base64.RawURLEncoding.Encode(enc, plain)
	for i := 0; i < b.N; i++ {
		dst := make([]byte, (len(enc)*3+4)/4)
		base64.RawURLEncoding.Decode(dst, enc)
	}
}

var plain = []byte(`You don` + "`" + `t care for me,
（你不关心我）
you don` + "`" + `t carry where I have been,
（不在乎我去何处）
I` + "`" + `ve done all I could,
（我已竭尽所能）
so that I could be with you.
（只为和你一起）
Anyway you want,
（无论你想什么）
I do everything you need,
（我愿意为你做任何事）
Maybe now you can see,
（或许你现在才明白）
that our love was went to be.
（我们的爱是什么）
But I was so wrong,
（但我错了）
always thought I could be strong.
（总以为自己能坚强）
When you left me here,
（当你离开我）
you took my heart away dear.
（也带走了我的心）
I feel so alone,
（我感到如此孤单）
I’ve miss you so long.
（已错过你太久）
I just can’t carry on,
（我无法坚持）
feeling lost at all alone.
（孤独而不知所措）
You love me with a whole broken heart,
（你爱着我但你的心已破碎）
left me here thinking why we fall apart.
（留下我面对这无言的结局）
But I was so wrong,
（但我错了）
always thought I could be strong.
（总以为自己能坚强）
When you left me here,
（当你离开我）
you took my heart away dear.
（也带走了我的心）
I feel so alone,
（我感到如此孤单）
I've been missing you so long . that our love was meant to be .
（已错过你太久）
I just can’t carry on,
（我无法坚持）
feeling lost at all alone.
（孤独而不知所措）
You love me with a whole broken heart,
（你爱着我但你的心已破碎）
left me here thinking why we fall apart.
（留下我面对这无言的结局）
But I was so wrong,
（但我错了）
always thought I could be strong.
（总以为自己能坚强）
When you left me here,
（当你离开我）
you took my heart away dear.
（也带走了我的心）
I feel so alone,
（我感到如此孤单）`)
