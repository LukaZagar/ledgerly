package parse

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// Table is a decoded CSV file: a header row plus the data rows, with the
// delimiter that was used so callers can report it.
type Table struct {
	Header    []string
	Rows      [][]string
	Delimiter rune
}

// candidateDelimiters are tried, most-specific first, when auto-detecting.
var candidateDelimiters = []rune{';', '\t', ',', '|'}

// ReadOptions tunes how a CSV is read. The zero value auto-detects the
// delimiter and skips nothing, which is what ReadCSV uses.
type ReadOptions struct {
	// Delimiter is the field separator; 0 means auto-detect.
	Delimiter rune
	// SkipRows drops this many decoded lines before the header is read, for
	// banks that print a preamble above the actual table.
	SkipRows int
}

// ReadCSV reads a CSV from r, auto-detecting both the byte encoding and the
// field delimiter. German bank exports are frequently semicolon-separated and
// encoded as Windows-1252 with a UTF-8 BOM, so we cannot assume the Go
// defaults.
func ReadCSV(r io.Reader) (*Table, error) {
	return ReadCSVOpts(r, ReadOptions{})
}

// ReadCSVWithDelimiter is like ReadCSV but uses the supplied delimiter instead
// of guessing it.
func ReadCSVWithDelimiter(r io.Reader, delim rune) (*Table, error) {
	return ReadCSVOpts(r, ReadOptions{Delimiter: delim})
}

// ReadCSVOpts reads a CSV honouring the given options.
func ReadCSVOpts(r io.Reader, opts ReadOptions) (*Table, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	text := decodeToUTF8(raw)
	if opts.SkipRows > 0 {
		text = dropLines(text, opts.SkipRows)
	}
	delim := opts.Delimiter
	if delim == 0 {
		delim = detectDelimiter(text)
	}
	return parseCSV(text, delim)
}

// dropLines removes the first n lines from text.
func dropLines(text string, n int) string {
	for ; n > 0; n-- {
		if i := strings.IndexByte(text, '\n'); i >= 0 {
			text = text[i+1:]
		} else {
			return ""
		}
	}
	return text
}

func parseCSV(text string, delim rune) (*Table, error) {
	cr := csv.NewReader(strings.NewReader(text))
	cr.Comma = delim
	cr.FieldsPerRecord = -1 // tolerate ragged rows; banks add trailing fields
	cr.LazyQuotes = true
	cr.TrimLeadingSpace = false

	records, err := cr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv: %w", err)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("read csv: no rows")
	}
	return &Table{
		Header:    records[0],
		Rows:      records[1:],
		Delimiter: delim,
	}, nil
}

// detectDelimiter picks the delimiter that appears most often in the first
// non-empty line of the file.
func detectDelimiter(text string) rune {
	line := firstLine(text)
	best := ';'
	bestCount := -1
	for _, d := range candidateDelimiters {
		if c := strings.Count(line, string(d)); c > bestCount {
			best, bestCount = d, c
		}
	}
	return best
}

func firstLine(text string) string {
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) != "" {
			return line
		}
	}
	return text
}

// decodeToUTF8 strips a UTF-8 BOM and transcodes the bytes to valid UTF-8. If
// the input is already valid UTF-8 it is returned unchanged; otherwise it is
// interpreted as Windows-1252, the de-facto encoding of older bank exports.
func decodeToUTF8(b []byte) string {
	b = bytes.TrimPrefix(b, []byte{0xEF, 0xBB, 0xBF}) // UTF-8 BOM
	if utf8.Valid(b) {
		return string(b)
	}
	var sb strings.Builder
	sb.Grow(len(b))
	for _, c := range b {
		sb.WriteRune(windows1252ToRune(c))
	}
	return sb.String()
}

// windows1252ToRune maps a single Windows-1252 byte to its Unicode rune. Bytes
// below 0x80 are ASCII; 0xA0–0xFF match Latin-1 (and thus Unicode) directly;
// 0x80–0x9F are the printable extras that differ, the most relevant being the
// Euro sign at 0x80.
func windows1252ToRune(c byte) rune {
	if c < 0x80 || c >= 0xA0 {
		return rune(c)
	}
	if r, ok := win1252High[c]; ok {
		return r
	}
	return rune(c)
}

var win1252High = map[byte]rune{
	0x80: '€', 0x82: '‚', 0x83: 'ƒ', 0x84: '„', 0x85: '…',
	0x86: '†', 0x87: '‡', 0x88: 'ˆ', 0x89: '‰', 0x8A: 'Š',
	0x8B: '‹', 0x8C: 'Œ', 0x8E: 'Ž', 0x91: '‘', 0x92: '’',
	0x93: '“', 0x94: '”', 0x95: '•', 0x96: '–', 0x97: '—',
	0x98: '˜', 0x99: '™', 0x9A: 'š', 0x9B: '›', 0x9C: 'œ',
	0x9E: 'ž', 0x9F: 'Ÿ',
}
