package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/kociumba/gabagool-file/gabagool"
)

func displayData(data []byte, dataType gabagool.DataTypes) error {

	ext := ".txt"
	def := ""

	switch dataType {
	case gabagool.Text:
		ext = ".txt"
	case gabagool.Image:
		ext = ".png"
	case gabagool.Bytes:
		ext = ".bin"
	}

	log.Info("", "Data content:", string(data))

	// Create a temporary file with the text data
	tempFile, err := os.CreateTemp("", "gabagool-*"+ext)
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	log.Info(tempFile.Name())

	// TODO figure out some forced way of specyfying dimensions of the image
	// if dataType == gabagool.Image {
	// 	// Convert the canvas pixels to a png
	// 	// we need to loop over the pixels and use the image package
	// 	img := image.NewRGBA(image.Rect(0, 0, f.Width, f.Height))
	// 	for x := 0; x < f.Width; x++ {
	// 		for y := 0; y < f.Height; y++ {
	// 			img.Set(x, y, color.RGBA{R: f.Pixels[y*f.Width+x].R, G: f.Pixels[y*f.Width+x].G, B: f.Pixels[y*f.Width+x].B, A: f.Pixels[y*f.Width+x].A})
	// 		}
	// 	}
	// 	// create a temporary file to write the image to
	// 	tempFile, err := os.CreateTemp("", "gabagool-*"+ext)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer os.Remove(tempFile.Name())
	// 	// write the image to the file
	// 	err = png.Encode(tempFile, img)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer tempFile.Close()
	// }

	// Write the text data to the temporary file
	_, err = tempFile.Write(data)
	if err != nil {
		return err
	}

	if err := tempFile.Close(); err != nil {
		return err
	}

	// Close the temporary file
	defer tempFile.Close()

	switch dataType {
	case gabagool.Text:
		log.Info("text")
		def = getDefaultTextEditor()
	case gabagool.Image:
		log.Info("image")
		def = getDefaultImageViewer()
	case gabagool.Bytes:
		log.Info("bytes")
		def = getDefaultBinaryEditor()
	}

	def = strings.ReplaceAll(def, "\"%1\"", "\"")
	def = strings.ReplaceAll(def, "\"", "")

	log.Info(def)

	// Launch the default text editor with the temporary file
	cmd := exec.Command(def, `"`+tempFile.Name()+`"`)

	log.Info(cmd)

	return cmd.Run()
}
