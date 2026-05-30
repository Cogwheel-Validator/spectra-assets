package transformer

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"os"

	xdraw "golang.org/x/image/draw"
)

const logoSize = 256

// decodeImageFile opens path and verifies it is a decodable image, returning the decoded result.
func decodeImageFile(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("not a valid image: %w", err)
	}
	return img, nil
}

// fetchImageFromURL downloads and decodes an image from a URL.
func fetchImageFromURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("invalid image from %s: %w", url, err)
	}
	return img, nil
}

// resizeCover scales src to a size×size square using cover mode (scale to fill, then center-crop),
// matching the default behaviour of Node.js sharp's resize with fit: 'cover'.
func resizeCover(src image.Image, size int) image.Image {
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	scaleX := float64(size) / float64(srcW)
	scaleY := float64(size) / float64(srcH)
	scale := scaleX
	if scaleY > scaleX {
		scale = scaleY
	}

	scaledW := int(float64(srcW)*scale + 0.5)
	scaledH := int(float64(srcH)*scale + 0.5)

	scaled := image.NewRGBA(image.Rect(0, 0, scaledW, scaledH))
	xdraw.CatmullRom.Scale(scaled, scaled.Bounds(), src, srcBounds, xdraw.Src, nil)

	offsetX := (scaledW - size) / 2
	offsetY := (scaledH - size) / 2

	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	xdraw.Draw(dst, dst.Bounds(), scaled, image.Pt(offsetX, offsetY), xdraw.Src)

	return dst
}

// encodeJPEG encodes img as a JPEG and returns the bytes.
func encodeJPEG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	options := &jpeg.Options{
		Quality: 80,
	}
	if err := jpeg.Encode(&buf, img, options); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
