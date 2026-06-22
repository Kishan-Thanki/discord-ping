package bot

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

// initFonts loads and caches the TTF font faces onto the Bot struct.
func (b *Bot) initFonts() {
	fontBytes, err := os.ReadFile("assets/font.ttf")
	if err != nil {
		slog.Warn("Could not read assets/font.ttf, welcome images will not have text", "err", err)
		return
	}

	f, err := truetype.Parse(fontBytes)
	if err != nil {
		slog.Error("Failed to parse font", "err", err)
		return
	}

	b.fontFace48 = truetype.NewFace(f, &truetype.Options{Size: 48})
	b.fontFace32 = truetype.NewFace(f, &truetype.Options{Size: 32})

	slog.Info("Fonts loaded and cached into memory successfully")
}

func (b *Bot) generateWelcomeImage(avatarURL, username string) (*bytes.Buffer, error) {
	const w = 800
	const h = 400

	dc := gg.NewContext(w, h)

	// Draw Background
	dc.SetColor(color.RGBA{R: 30, G: 33, B: 36, A: 255})
	dc.Clear()

	// Aesthetic circles
	dc.SetColor(color.RGBA{R: 114, G: 137, B: 218, A: 50})
	dc.DrawCircle(800, 0, 300)
	dc.Fill()
	dc.DrawCircle(0, 400, 200)
	dc.Fill()

	// Fetch Avatar
	avatarImg, err := fetchImage(avatarURL)
	if err == nil {
		// Draw Avatar with Circular Mask
		dc.DrawCircle(w/2, h/2-30, 80)
		dc.Clip()
		scaledAvatar := resizeImage(avatarImg, 160, 160)
		dc.DrawImageAnchored(scaledAvatar, w/2, h/2-30, 0.5, 0.5)
		dc.ResetClip()

		// Draw stroke
		dc.SetColor(color.RGBA{R: 114, G: 137, B: 218, A: 255})
		dc.SetLineWidth(6)
		dc.DrawCircle(w/2, h/2-30, 80)
		dc.Stroke()
	}

	// Draw "Welcome" if font is loaded
	if b.fontFace48 != nil {
		dc.SetFontFace(b.fontFace48)
		dc.SetColor(color.White)
		dc.DrawStringAnchored("WELCOME TO THE SERVER", w/2, h-80, 0.5, 0.5)
	}

	// Draw Username if font is loaded
	if b.fontFace32 != nil {
		dc.SetFontFace(b.fontFace32)
		dc.SetColor(color.RGBA{R: 153, G: 170, B: 181, A: 255})
		dc.DrawStringAnchored(username, w/2, h-40, 0.5, 0.5)
	}

	buf := new(bytes.Buffer)
	if err := dc.EncodePNG(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func fetchImage(url string) (image.Image, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	return img, err
}

func resizeImage(src image.Image, targetW, targetH int) image.Image {
	bounds := src.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	for y := 0; y < targetH; y++ {
		for x := 0; x < targetW; x++ {
			srcX := x * srcW / targetW
			srcY := y * srcH / targetH
			dst.Set(x, y, src.At(bounds.Min.X+srcX, bounds.Min.Y+srcY))
		}
	}
	return dst
}
