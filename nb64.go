package nb64

import (
	"strconv"
)

const URLEncode = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

const (
	Dot   = '.'
	Tilde = '~'
)

type CorruptInputError int64

func (e CorruptInputError) Error() string {
	return "illegal input data at input byte " + strconv.FormatInt(int64(e), 10)
}

type Encoding struct {
	magic  byte
	encode [256]byte
	decode [256]byte
}

func NewEncoding(encode string, magic byte) *Encoding {
	if len(encode) != 64 {
		panic("encoding alphabet is not 64-bytes long")
	}
	for i := 0; i < len(encode); i++ {
		if encode[i] == magic {
			panic("encoding alphabet contains magic character")
		}
		// if encode[i] == '\n' || encode[i] == '\r' {
		// 	panic("encoding alphabet contains newline character")
		// }
	}

	e := &Encoding{magic: magic}
	for i := 0; i < len(encode); i++ {
		e.encode[i] = encode[i]
		e.encode[i|0x40] = encode[i]
		e.encode[i|0x80] = encode[i]
		e.encode[i|0xC0] = encode[i]
	}

	for i := 0; i < 256; i++ {
		e.decode[i] = 0xFF
	}
	for i := 0; i < len(encode); i++ {
		e.decode[encode[i]] = byte(i)
	}
	return e
}

var URLEncoding = NewEncoding(URLEncode, Dot)

func (e *Encoding) Encode(plain []byte) ([]byte, error) {
	l := len(plain)
	out := make([]byte, (l*7+1)/3)
	var from, to, n int
	for to < l {
		// <=127 base64
		for from = to; to < l && plain[to] <= 127; to++ {
		}
		for ; from < to; from += 6 {
			m := to - from
			switch m {
			default:
				m = 6
				out[n+6] = e.encode[plain[from+5]] // 6   0
				out[n+5] = plain[from+5] >> 6      // 5+1 6
				fallthrough
			case 5:
				out[n+5] = e.encode[plain[from+4]<<1|out[n+5]] // 5+1 6
				out[n+4] = plain[from+4] >> 5                  // 4+2 5
				fallthrough
			case 4:
				out[n+4] = e.encode[plain[from+3]<<2|out[n+4]] // 4+2 5
				out[n+3] = plain[from+3] >> 4                  // 3+3 4
				fallthrough
			case 3:
				out[n+3] = e.encode[plain[from+2]<<3|out[n+3]] // 3+3 4
				out[n+2] = plain[from+2] >> 3                  // 2+4 3
				fallthrough
			case 2:
				out[n+2] = e.encode[plain[from+1]<<4|out[n+2]] // 2+4 3
				out[n+1] = plain[from+1] >> 2                  // 1+5 2
				fallthrough
			case 1:
				out[n+1] = e.encode[plain[from]<<5|out[n+1]] // 1+5 2
				out[n] = e.encode[plain[from]>>1]            //   6 1
			}
			n += m + 1
		}

		// >127 64base
		for from = to; to < l && plain[to] > 127; to++ {
		}
		if from < to {
			out[n] = e.magic
			n++
			for ; from < to; from++ {
				b := plain[from]
				switch {
				case b < 0xC0: // utf8 字节尾
					out[n] = e.encode[b&0x3F]
					n++
				case b < 0xE0: // utf8 双字节头
					out[n] = e.encode[0]
					out[n+1] = e.encode[b&0x1F]
					n += 2
				case b < 0xF0: // utf8 三字节头
					out[n] = e.encode[b&0x0F]
					n++
				case b == 0xF0: // utf8 四字节头，只支持到18bits
				default:
					return nil, CorruptInputError(from)
				}
			}
			out[n] = e.magic
			n++
		}
	}
	return out[:n], nil
}

// Decode convert bn64 bytes to origin bytes.
func (e *Encoding) Decode(enc []byte) ([]byte, error) {
	l := len(enc)
	ll := l
	if l > 6 {
		ll = (l - 2) * 4 / 3
	}
	out := make([]byte, ll)
	var from, to, n int
	for to < l {
		// <=127 base64
		for from = to; to < l && enc[to] != e.magic; to++ {
			if e.decode[enc[to]] == 0xFF {
				return nil, CorruptInputError(to)
			}
		}
		for ; from < to; from += 7 {
			m := to - from
			switch m {
			default:
				m = 7
				out[n+5] = e.decode[enc[from+6]] // 5 = 5:1<<6 + 6:6>>0
				fallthrough
			case 6:
				out[n+5] = ((e.decode[enc[from+5]] << 6) & 0x7F) | out[n+5] // 5 = 5:1<<6 + 6:6>>0
				out[n+4] = e.decode[enc[from+5]] >> 1                       // 4 = 4:2<<5 + 5:5>>1
				fallthrough
			case 5:
				out[n+4] = ((e.decode[enc[from+4]] << 5) & 0x7F) | out[n+4] // 4 = 4:2<<5 + 5:5>>1
				out[n+3] = e.decode[enc[from+4]] >> 2                       // 3 = 3:3<<4 + 4:4>>2
				fallthrough
			case 4:
				out[n+3] = ((e.decode[enc[from+3]] << 4) & 0x7F) | out[n+3]
				out[n+2] = e.decode[enc[from+3]] >> 3
				fallthrough
			case 3:
				out[n+2] = ((e.decode[enc[from+2]] << 3) & 0x7F) | out[n+2]
				out[n+1] = e.decode[enc[from+2]] >> 4
				fallthrough
			case 2:
				out[n+1] = ((e.decode[enc[from+1]] << 2) & 0x7F) | out[n+1]
				out[n] = e.decode[enc[from+1]] >> 5
				fallthrough
			case 1:
				out[n] = ((e.decode[enc[from]] << 1) & 0x7F) | out[n]
			}
			n += m - 1
		}
		to++

		// >127 64base
		for from = to; to < l && enc[to] != e.magic; to++ {
			if e.decode[enc[to]] == 0xFF {
				return nil, CorruptInputError(to)
			}
		}
		if (to-from)%3 != 0 {
			return nil, CorruptInputError(to)
		}
		for ; from < to; from += 3 {
			switch {
			case e.decode[enc[from]] == 0 && e.decode[enc[from+1]] < 0x20:
				// utf8 双字节8-11bits， < 800(1000 0000 0000, --10 0000 --00 0000)
				//  --000000 --0xxxxx --xxxxxx
				// =>        110xxxxx 10xxxxxx '·'
				out[n] = 0xC0 | e.decode[enc[from+1]]
				out[n+1] = 0x80 | e.decode[enc[from+2]]
				n += 2

			case e.decode[enc[from]] >= 0x10:
				//              --1xxxxx --xxxxxx --xxxxxx
				//              --x1xxxx --xxxxxx --xxxxxx
				// =>  11110000 10xxxxxx 10xxxxxx 10xxxxxx '𐄁'
				// utf8 四字节头17-18bits(只支持到18bits)，
				// >= 10000(1 0000 0000 0000 0000, 6位分隔 1 0000 --00 0000 --00 0000)
				out[n] = 0xF0
				out[n+1] = 0x80 | e.decode[enc[from]]
				out[n+2] = 0x80 | e.decode[enc[from+1]]
				out[n+3] = 0x80 | e.decode[enc[from+2]]
				n += 4

			default:
				//    --00xxxx --xxxxxx --xxxxxx
				// => 1110xxxx 10xxxxxx 10xxxxxx '中'
				// utf8 三字节 // 12-16bits
				out[n] = 0xE0 | e.decode[enc[from]]
				out[n+1] = 0x80 | e.decode[enc[from+1]]
				out[n+2] = 0x80 | e.decode[enc[from+2]]
				n += 3
			}
		}
		to++
	}
	return out[:n], nil
}
