package converter

import (
	"fmt"
	"io"
)

// WriteReport writes a summary of the conversion to the given writer.
func WriteReport(w io.Writer, inputFile string, format Format, lineCount int, err error) {
	if err != nil {
		fmt.Fprintf(w, "✗ Conversion failed [%s → %s]: %v\n", inputFile, format, err)
		return
	}
	fmt.Fprintf(w, "✓ Converted %s → %s (%d variables)\n", inputFile, format, lineCount)
}

// SupportedFormats returns a human-readable list of supported output formats.
func SupportedFormats() []Format {
	return []Format{
		FormatJSON,
		FormatYAML,
		FormatTOML,
		FormatShell,
		FormatDotEnv,
	}
}

// FormatNames returns the string names of all supported formats.
func FormatNames() []string {
	formats := SupportedFormats()
	names := make([]string, len(formats))
	for i, f := range formats {
		names[i] = string(f)
	}
	return names
}
