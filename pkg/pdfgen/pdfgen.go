package pdfgen

import (
	"context"
	"io"
	"os"
	"os/exec"
)

func Generate(ctx context.Context, w io.Writer, html []byte) error {
	tempHTMLFile, err := os.CreateTemp("", "*.html")
	if err != nil {
		return err
	}

	defer os.Remove(tempHTMLFile.Name())

	if _, err = tempHTMLFile.Write(html); err != nil {
		return err
	}

	if err = tempHTMLFile.Close(); err != nil {
		return err
	}

	tempPDFFile, err := os.CreateTemp("", "*.pdf")
	if err != nil {
		return err
	}

	tempPDFFilePath := tempPDFFile.Name()
	defer os.Remove(tempPDFFilePath)

	if err = tempPDFFile.Close(); err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "weasyprint", tempHTMLFile.Name(), tempPDFFilePath)

	err = cmd.Run()

	if err != nil {
		return err
	}

	pdfFile, err := os.Open(tempPDFFilePath)
	if err != nil {
		return err
	}

	defer pdfFile.Close()

	if _, err = io.Copy(w, pdfFile); err != nil {
		return err
	}

	return nil
}

