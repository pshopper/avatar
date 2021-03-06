package avatar

import (
	"image/color"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitials_Circle(t *testing.T) {
	size := 200

	newAvatar, err := NewAvatarFromInitials("John Smith", &InitialsOptions{
		Size:      size,
		NInitials: 2,
		FontPath:  getTestResource("test_data", "Arial.ttf"),
	})
	assert.NoError(t, err)

	round, err := newAvatar.Circle()
	assert.NoError(t, err)

	roundOutputPath := getTestResource("output", "round_john_smith_initials.png")
	roundFile, err := os.Create(roundOutputPath)
	assert.NoError(t, err)
	roundFile.Write(round)
}

func TestInitials_Square(t *testing.T) {
	size := 200

	newAvatar, err := NewAvatarFromInitials("John Smith", &InitialsOptions{
		Size:      size,
		NInitials: 2,
		FontPath:  getTestResource("test_data", "Arial.ttf"),
		TextColor: color.White,
		BgColor:   color.RGBA{0, 0, 255, 255},
	})
	assert.NoError(t, err)

	square, err := newAvatar.Square()
	assert.NoError(t, err)

	squareOutputPath := getTestResource("output", "square_john_smith_initials.png")
	squareFile, err := os.Create(squareOutputPath)
	assert.NoError(t, err)
	squareFile.Write(square)
}

func TestInitialsGradient_Circle(t *testing.T) {
	size := 800

	newAvatar, err := NewAvatarFromInitials("РЩ", &InitialsOptions{
		Size:      size,
		FontSize:  float64(size) / 2,
		NInitials: 2,
		FontPath:  getTestResource("test_data", "SFProText-Medium.ttf"),
		GradientTable: GradientTable{
			{MustParseHex("#82a7e8"), 0.0},
			{MustParseHex("#4b6ecd"), 1.0},
		},
	})
	assert.NoError(t, err)

	round, err := newAvatar.Circle()
	assert.NoError(t, err)

	roundOutputPath := getTestResource("output", "round_ivanov_ivan_initials.png")
	roundFile, err := os.Create(roundOutputPath)
	assert.NoError(t, err)
	roundFile.Write(round)
}
