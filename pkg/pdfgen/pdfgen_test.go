package pdfgen_test

import (
	"bytes"
	"context"
	"testing"

	"invoicehub/pkg/pdfgen"

	"github.com/stretchr/testify/require"
)

func Test_Generate(t *testing.T) {
	ctx := context.Background()

	buff := bytes.NewBuffer(nil)

	html := `<html><body><h1>test</h1></body></html>`
	err := pdfgen.Generate(ctx, buff, []byte(html))
	require.NoError(t, err)

	require.NotEmpty(t, buff.Bytes())
}
