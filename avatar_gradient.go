package avatar

import (
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/draw"
)

type Color struct {
	colorful.Color
}

// This table contains the "keypoints" of the colorgradient you want to generate.
// The position of each keypoint has to live in the range [0,1]
type GradientTable []struct {
	Col Color
	Pos float64
}

// This is the meat of the gradient computation. It returns a HCL-blend between
// the two colors around `t`.
// Note: It relies heavily on the fact that the gradient keypoints are sorted.
func (gt GradientTable) getInterpolatedColorFor(t float64) Color {
	for i := 0; i < len(gt)-1; i++ {
		c1 := gt[i]
		c2 := gt[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			// We are in between c1 and c2. Go blend them!
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			return Color{Color: c1.Col.BlendHcl(c2.Col.Color, t).Clamped()}
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return gt[len(gt)-1].Col
}

// This is a very nice thing Golang forces you to do!
// It is necessary so that we can write out the literal of the colortable below.
func MustParseHex(s string) Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("MustParseHex: " + err.Error())
	}
	return Color{Color: c}
}

func ParseHex(s string) (Color, error) {
	c, err := colorful.Hex(s)
	if err != nil {
		return Color{}, err
	}
	return Color{Color: c}, nil
}

func newGradientImage(height, width int, gradientTable GradientTable) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := height - 1; y >= 0; y-- {
		c := gradientTable.getInterpolatedColorFor(float64(y) / float64(height))
		draw.Draw(img, image.Rect(0, y, width, y+1), &image.Uniform{C: c}, image.ZP, draw.Src)
	}

	return img
}
