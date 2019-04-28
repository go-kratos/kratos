package drawimg

import (
	"errors"
	"fmt"
	"image"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype/raster"
	"github.com/golang/freetype/truetype"
)

const (
	nGlyphs     = 256
	nXFractions = 4
	nYFractions = 1
)

type cacheEntry struct {
	valid        bool
	glyph        truetype.Index
	advanceWidth fixed.Int26_6
	mask         *image.Alpha
	offset       image.Point
}

// ParseFont just calls the Parse function from the freetype/truetype package.
func ParseFont(b []byte) (*truetype.Font, error) {
	return truetype.Parse(b)
}

// Pt converts from a co-ordinate pair measured in pixels to a fixed.Point26_6 co-ordinate pair measured in fixed.Int26_6 units.
func Pt(x, y int) fixed.Point26_6 {
	return fixed.Point26_6{
		X: fixed.Int26_6(x << 6),
		Y: fixed.Int26_6(y << 6),
	}
}

// A Context holds the state for drawing text in a given font and size.
type Context struct {
	r        *raster.Rasterizer
	f        *truetype.Font
	glyphBuf truetype.GlyphBuf
	clip     image.Rectangle

	dst           draw.Image
	src           image.Image
	fontSize, dpi float64
	scale         fixed.Int26_6
	hinting       font.Hinting
	cache         [nGlyphs * nXFractions * nYFractions]cacheEntry
}

// PointToFixed converts the given number of points (as in "a 12 point font") into a 26.6 fixed point number of pixels.
func (c *Context) PointToFixed(x float64) fixed.Int26_6 {
	return fixed.Int26_6(x * float64(c.dpi) * (64.0 / 72.0))
}

// drawContour draws the given closed contour with the given offset.
func (c *Context) drawContour(ps []truetype.Point, dx, dy fixed.Int26_6) {
	if len(ps) == 0 {
		return
	}
	start := fixed.Point26_6{
		X: dx + ps[0].X,
		Y: dy - ps[0].Y,
	}
	others := []truetype.Point(nil)
	if ps[0].Flags&0x01 != 0 {
		others = ps[1:]
	} else {
		last := fixed.Point26_6{
			X: dx + ps[len(ps)-1].X,
			Y: dy - ps[len(ps)-1].Y,
		}
		if ps[len(ps)-1].Flags&0x01 != 0 {
			start = last
			others = ps[:len(ps)-1]
		} else {
			start = fixed.Point26_6{
				X: (start.X + last.X) / 2,
				Y: (start.Y + last.Y) / 2,
			}
			others = ps
		}
	}
	c.r.Start(start)
	q0, on0 := start, true
	for _, p := range others {
		q := fixed.Point26_6{
			X: dx + p.X,
			Y: dy - p.Y,
		}
		on := p.Flags&0x01 != 0
		if on {
			if on0 {
				c.r.Add1(q)
			} else {
				c.r.Add2(q0, q)
			}
		} else {
			if on0 {
				// No-op.
			} else {
				mid := fixed.Point26_6{
					X: (q0.X + q.X) / 2,
					Y: (q0.Y + q.Y) / 2,
				}
				c.r.Add2(q0, mid)
			}
		}
		q0, on0 = q, on
	}
	if on0 {
		c.r.Add1(start)
	} else {
		c.r.Add2(q0, start)
	}
}

// rasterize returns the advance width, glyph mask and integer-pixel offset to render the given glyph at the given sub-pixel offsets.
// The 26.6 fixed point arguments fx and fy must be in the range [0, 1).
func (c *Context) rasterize(glyph truetype.Index, fx, fy fixed.Int26_6) (
	fixed.Int26_6, *image.Alpha, image.Point, error) {

	if err := c.glyphBuf.Load(c.f, c.scale, glyph, c.hinting); err != nil {
		return 0, nil, image.Point{}, err
	}
	xmin := int(fx+c.glyphBuf.Bounds.Min.X) >> 6
	ymin := int(fy-c.glyphBuf.Bounds.Max.Y) >> 6
	xmax := int(fx+c.glyphBuf.Bounds.Max.X+0x3f) >> 6
	ymax := int(fy-c.glyphBuf.Bounds.Min.Y+0x3f) >> 6
	if xmin > xmax || ymin > ymax {
		return 0, nil, image.Point{}, errors.New("freetype: negative sized glyph")
	}
	fx -= fixed.Int26_6(xmin << 6)
	fy -= fixed.Int26_6(ymin << 6)
	c.r.Clear()
	e0 := 0
	for _, e1 := range c.glyphBuf.Ends {
		c.drawContour(c.glyphBuf.Points[e0:e1], fx, fy)
		e0 = e1
	}
	a := image.NewAlpha(image.Rect(0, 0, xmax-xmin, ymax-ymin))
	c.r.Rasterize(raster.NewAlphaSrcPainter(a))
	return c.glyphBuf.AdvanceWidth, a, image.Point{xmin, ymin}, nil
}

// glyph returns the advance width, glyph mask and integer-pixel offset to  render the given glyph at the given sub-pixel point.
func (c *Context) glyph(glyph truetype.Index, p fixed.Point26_6) (fixed.Int26_6, *image.Alpha, image.Point, error) {
	ix, fx := int(p.X>>6), p.X&0x3f
	iy, fy := int(p.Y>>6), p.Y&0x3f
	tg := int(glyph) % nGlyphs
	tx := int(fx) / (64 / nXFractions)
	ty := int(fy) / (64 / nYFractions)
	t := ((tg*nXFractions)+tx)*nYFractions + ty
	if e := c.cache[t]; e.valid && e.glyph == glyph {
		return e.advanceWidth, e.mask, e.offset.Add(image.Point{ix, iy}), nil
	}
	advanceWidth, mask, offset, err := c.rasterize(glyph, fx, fy)
	if err != nil {
		return 0, nil, image.Point{}, err
	}
	c.cache[t] = cacheEntry{true, glyph, advanceWidth, mask, offset}
	return advanceWidth, mask, offset.Add(image.Point{ix, iy}), nil
}

// DrawString draws s at p and returns p advanced by the text extent.
func (c *Context) DrawString(s string, p fixed.Point26_6) (fixed.Point26_6, error) {
	if c.f == nil {
		return fixed.Point26_6{}, errors.New("freetype: DrawText called with a nil font")
	}
	prev, hasPrev := truetype.Index(0), false
	for _, rune := range s {
		index := c.f.Index(rune)
		if hasPrev {
			kern := c.f.Kern(c.scale, prev, index)
			if c.hinting != font.HintingNone {
				kern = (kern + 32) &^ 63
			}
			p.X += kern
		}
		advanceWidth, mask, offset, err := c.glyph(index, p)
		if err != nil {
			return fixed.Point26_6{}, err
		}
		p.X += advanceWidth
		glyphRect := mask.Bounds().Add(offset)
		dr := c.clip.Intersect(glyphRect)
		if !dr.Empty() {
			mp := image.Point{0, dr.Min.Y - glyphRect.Min.Y}
			draw.DrawMask(c.dst, dr, c.src, image.ZP, mask, mp, draw.Over)
		}
		prev, hasPrev = index, true
	}
	return p, nil
}

// recalc recalculates scale and bounds values from the font size, screen resolution and font metrics, and invalidates the glyph cache.
func (c *Context) recalc() {
	c.scale = fixed.Int26_6(c.fontSize * c.dpi * (64.0 / 72.0))
	if c.f == nil {
		c.r.SetBounds(0, 0)
	} else {
		b := c.f.Bounds(c.scale)
		xmin := +int(b.Min.X) >> 6
		ymin := -int(b.Max.Y) >> 6
		xmax := +int(b.Max.X+63) >> 6
		ymax := -int(b.Min.Y-63) >> 6
		c.r.SetBounds(xmax-xmin, ymax-ymin)
	}
	for i := range c.cache {
		c.cache[i] = cacheEntry{}
	}
}

// SetDPI sets the screen resolution in dots per inch.
func (c *Context) SetDPI(dpi float64) {
	if c.dpi == dpi {
		return
	}
	c.dpi = dpi
	c.recalc()
}

// SetFont sets the font used to draw text.
func (c *Context) SetFont(f *truetype.Font) {
	if c.f == f {
		return
	}
	c.f = f
	c.recalc()
}

// SetFontSize sets the font size in points (as in "a 12 point font").
func (c *Context) SetFontSize(fontSize float64) {
	if c.fontSize == fontSize {
		return
	}
	c.fontSize = fontSize
	c.recalc()
}

// SetHinting sets the hinting policy.
func (c *Context) SetHinting(hinting font.Hinting) {
	c.hinting = hinting
	for i := range c.cache {
		c.cache[i] = cacheEntry{}
	}
}

// SetDst sets the destination image for draw operations.
func (c *Context) SetDst(dst draw.Image) {
	c.dst = dst
}

// SetSrc sets the source image for draw operations. This is typically an image.Uniform.
func (c *Context) SetSrc(src image.Image) {
	c.src = src
}

// SetClip sets the clip rectangle for drawing.
func (c *Context) SetClip(clip image.Rectangle) {
	c.clip = clip
}

// NewContext creates a new Context.
func NewContext() *Context {
	return &Context{
		r:        raster.NewRasterizer(0, 0),
		fontSize: 12,
		dpi:      72,
		scale:    12 << 6,
	}
}

// FontBox get  textbox.
func (c *Context) FontBox(s string, p fixed.Point26_6) (fixed.Point26_6, error) {
	if c.f == nil {
		return fixed.Point26_6{}, errors.New("freetype: DrawText called with a nil font")
	}
	prev, hasPrev := truetype.Index(0), false
	for _, rune := range s {
		index := c.f.Index(rune)
		if hasPrev {
			kern := c.f.Kern(c.scale, prev, index)
			if c.hinting != font.HintingNone {
				kern = (kern + 32) &^ 63
			}
			p.X += kern
		}
		advanceWidth, _, _, err := c.glyph(index, p)
		if err != nil {
			return fixed.Point26_6{}, err
		}
		p.X += advanceWidth
	}
	return p, nil
}

// Int26_6ToString convert Int26_6 to string.
func Int26_6ToString(x fixed.Int26_6) string {
	const shift, mask = 6, 1<<6 - 1
	if x >= 0 {
		return fmt.Sprintf("%d:%02d", int32(x>>shift), int32(x&mask))
	}
	x = -x
	if x >= 0 {
		return fmt.Sprintf("-%d:%02d", int32(x>>shift), int32(x&mask))
	}
	return "-33554432:00" // The minimum value is -(1<<25).
}
