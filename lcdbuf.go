
package lcdbuf

import (
	"fmt"
	"os"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"gitcafe.com/nuomi-studio/pcf.git"
)

/*
	(0,0)                (127,0)
  ..0..................
  ..1..................
  ..0..................
  ..1..................
  ..0..................
  ..1..................
  ..0..................
  ..1..................
	(0,8)                (127,8)
*/

type Buf struct {
	Pix []byte
	W, H int
}

func PCFText(font *pcf.File, str string) (buf *Buf) {
	// hack
	H := 16
	Base := 13

	buf = &Buf{
		H: H,
	}

	metrics := []pcf.MetricEntry{}
	bits := [][]byte{}
	strides := []int{}

	for _, r := range str {
		b, me, stride, err := font.Lookup(r)
		if err != nil {
			continue
		}
		metrics = append(metrics, me)
		bits = append(bits, b)
		strides = append(strides, stride)
		buf.W += me.CharWidth
	}

	buf.Pix = make([]byte, buf.W*(buf.H+7)/8)

	x := 0
	for j, me := range metrics {
		b := bits[j]
		stride := strides[j]
		for i := 0; i < len(b); i += stride {
			x1 := x + me.LeftSidedBearing
			y1 := Base - me.CharAscent + i/stride
			buf.DrawHoriBits(x1, y1, b[i:i+stride], me.CharWidth)
		}
		x += me.CharWidth
	}

	return
}

func DrawOffset(orig, cur *Buf, x, y, offset int) {
	for i := 0; i < cur.H/8; i++ {
		for j := 0; j < cur.W && j+x<orig.W-1; j++ {
			b := cur.Pix[i*cur.W+j+offset]
			if oi := (y/8+i)*orig.W+j+x; oi < len(orig.Pix) {
				orig.Pix[oi] = b
			}
		}
	}
	return
}

func Draw(orig, cur *Buf, x, y int) {
	DrawOffset(orig, cur, x, y, 0)
}

func (l *Buf) DrawHoriBits(x, y int, b []byte, w int) {
	for bi := 0; bi < w; bi++ {
		bm := byte(1)<<byte(7-(bi%8))
		bv := b[bi/8]
		pi := y/8*l.W+x+bi
		pm := byte(1)<<byte(y%8)
		pv := l.Pix[pi]
		pv &= ^pm
		if bv & bm != 0 {
			pv |= pm
		}
		l.Pix[pi] = pv
	}
}

func (l *Buf) Inverse() {
	for i := range l.Pix {
		l.Pix[i] = ^l.Pix[i]
	}
}

func (l *Buf) Clear() {
	for i := range l.Pix {
		l.Pix[i] = 0
	}
}

func (l *Buf) DumpGoCode(fname string) {
	f, err := os.Create(fname)
	if err != nil {
		return
	}
	fmt.Fprintf(f, "pic := lcdbuf.Buf{\n")
	fmt.Fprintf(f, "  W: %d,\n", l.W)
	fmt.Fprintf(f, "  H: %d,\n", l.H)
	fmt.Fprintf(f, "  Pix: []byte{\n")
	for i := 0; i < l.H/8; i++ {
		fmt.Fprint(f, "    ")
		for j := 0; j < l.W; j++ {
			fmt.Fprintf(f, "0x%x,", l.Pix[i*l.W+j])
		}
		fmt.Fprint(f, "\n")
	}
	fmt.Fprintln(f, "  },")
	fmt.Fprintln(f, "}")
	f.Close()
}

func (l *Buf) DumpCCode(fname string) {
	f, err := os.Create(fname)
	if err != nil {
		return
	}
	fmt.Fprintf(f, "/* W:%d H:%d */\n", l.W, l.H)
	fmt.Fprintf(f, "u8 arr[] = {\n")
	for i := 0; i < l.H/8; i++ {
		for j := 0; j < l.W; j++ {
			fmt.Fprintf(f, "0x%x,", l.Pix[i*l.W+j])
		}
		fmt.Fprint(f, "\n")
	}
	fmt.Fprintln(f, "}")
	f.Close()
}

func (l *Buf) DumpAscii(fname string) {
	f, err := os.Create(fname)
	if err != nil {
		return
	}
	for i := 0; i < l.H; i++ {
		for j := 0; j < l.W; j++ {
			if l.Pix[(i/8)*l.W+j] & (1<<byte(i%8)) != 0 {
				fmt.Fprint(f, "@")
			} else {
				fmt.Fprint(f, ".")
			}
		}
		fmt.Fprintln(f)
	}
	f.Close()
}

func New(w, h int) *Buf {
	rh := (h+7)/8
	return &Buf{
		Pix: make([]byte, rh*w),
		W: w,
		H: h,
	}
}

func FromImageFile(fname string) (buf *Buf, err error) {
	var f *os.File
	if f, err = os.Open(fname); err != nil {
		return
	}

	var img image.Image
	img, _, err = image.Decode(f)
	if err != nil {
		return
	}
	rect := img.Bounds().Max

	h := (rect.Y + 7) / 8

	buf = &Buf{
		W: rect.X,
		H: rect.Y,
		Pix: make([]byte, rect.X*h),
	}

	for i := 0; i < h; i++ {
		for j := 0; j < rect.X; j++ {
			by := uint32(0)
			for k := 0; k < 8; k++ {
				c := img.At(j, i*8+k)
				r,g,b,_ := c.RGBA()
				y := (r+g+b)/3
				if y < 65536/2 {
					by |= 1<<uint32(k)
				}
			}
			buf.Pix[i*buf.W+j] = byte(by)
		}
	}

	return
}

