package avatar

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"regexp"
	"strings"
	"unicode"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type InitialsOptions struct {
	BgColor       color.Color
	Size          int
	FontPath      string
	FontSize      float64
	Font          *truetype.Font
	TextColor     color.Color
	NInitials     int
	GradientTable GradientTable
}

type Initials struct {
	source        string
	options       *InitialsOptions
	originalImage image.Image
	squareImage   image.Image
	circleImage   image.Image
}

func (i Initials) Source() []byte {
	return []byte(i.source)
}

func (i Initials) loadOriginalImage() (image.Image, error) {
	text := i.source

	nInitials := i.options.nInitials()

	if nInitials > 0 {
		text = getInitials(text, nInitials)
	}

	size := i.options.size() * 3 // 3 times bigger for better quality

	// Draw background img
	var dst *image.RGBA
	if len(i.options.GradientTable) != 0 {
		dst = newGradientImage(size, size, i.options.GradientTable)
	} else {
		imgRect := image.Rect(0, 0, size, size)
		dst = image.NewRGBA(imgRect)
		draw.Draw(
			dst,
			dst.Bounds(),
			image.NewUniform(i.options.bgColor()),
			image.ZP,
			draw.Src)
	}

	ftFont, err := i.options.font()
	if err != nil {
		return nil, err
	}

	fontSize := i.options.FontSize
	if fontSize == 0 {
		fontSize = getFontSizeThatFits([]byte(text), float64(size), ftFont)
	}

	fontFace := truetype.NewFace(ftFont, &truetype.Options{
		Size: fontSize,
	})

	fd := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(i.options.textColor()),
		Face: fontFace,
	}

	// Figure out baseline and adv for string in img
	txtWidth := fd.MeasureBytes([]byte(text))
	txtWidthInt := int(txtWidth >> 6)

	bounds, _ := fd.BoundBytes([]byte(text))
	txtHeight := bounds.Max.Y - bounds.Min.Y
	txtHeightInt := int(txtHeight >> 6)

	advLine := (size / 2) - (txtWidthInt / 2)
	baseline := (size / 2) + (txtHeightInt / 2)

	fd.Dot = fixed.Point26_6{X: fixed.Int26_6(advLine << 6), Y: fixed.Int26_6(baseline << 6)}

	fd.DrawBytes([]byte(text))

	return dst, nil
}

func (i Initials) originalImg() image.Image {
	return i.originalImage
}

// Generates the square avatar
// It returns the avatar image in []byte or an error something went wrong
func (i Initials) Square() ([]byte, error) {
	return square(i, i.options)
}

func (i Initials) Circle() ([]byte, error) {
	return circle(i, i.options)
}

func (o InitialsOptions) bgColor() color.Color {
	return bgColor(o.BgColor)
}

func (o InitialsOptions) size() int {
	return size(o.Size)
}

func (o InitialsOptions) nInitials() int {
	if o.NInitials == 0 {
		return defaultNInitials
	}
	return o.NInitials
}

func (o InitialsOptions) font() (*truetype.Font, error) {
	if o.Font != nil {
		return o.Font, nil
	}

	if strings.TrimSpace(o.FontPath) == "" {
		return nil, errors.New("No font path given")
	}

	fontBytes, err := ioutil.ReadFile(o.FontPath)
	if err != nil {
		return nil, err
	}

	return freetype.ParseFont(fontBytes)
}

func (o InitialsOptions) textColor() color.Color {
	if o.TextColor == nil {
		return defaultTxtColor
	}

	return o.TextColor
}

func getFontSizeThatFits(text []byte, imgWidth float64, ftFont *truetype.Font) float64 {
	fontSize := float64(100)

	drawer := font.Drawer{
		Face: truetype.NewFace(ftFont, &truetype.Options{
			Size: fontSize,
		}),
	}

	tw := float64(drawer.MeasureBytes(text) >> 6)

	ratio := fontSize / tw

	return ratio * (imgWidth - (40./100)*imgWidth)
}

func getInitials(text string, nChars int) string {
	if len(text) == 0 {
		return ""
	}

	var initials []rune
	var previous = ' '

	regEmail := regexp.MustCompile("^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")
	skipFromAt := regEmail.MatchString(text)

	for _, ch := range []rune(text) {
		if skipFromAt == true && ch == '@' {
			break
		}

		if isSymbol(ch) {
			previous = ch
			continue
		}

		if ((unicode.IsUpper(ch) && unicode.IsLower(previous)) || (unicode.IsLower(ch) && len(initials) == 0)) || isSymbol(previous) {
			initials = append(initials, ch)
			previous = ch
		}
	}

	for i := len(initials); i < nChars && len([]rune(text)) > i; i++ {
		if isSymbol(rune(text[i])) {
			continue
		}

		initials = append(initials, []rune(text)[i])
	}

	return string(initials)
}

func isSymbol(ch rune) bool {
	return unicode.IsSymbol(ch) || unicode.IsSpace(ch) || unicode.IsPunct(ch)
}
