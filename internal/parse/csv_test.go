package parse

import (
	"strings"
	"testing"
)

func TestReadCSVSemicolon(t *testing.T) {
	const in = "Buchungstag;Auftraggeber;Verwendungszweck;Betrag\n" +
		"02.01.2026;REWE Markt GmbH;Einkauf;-23,45\n" +
		"03.01.2026;Arbeitgeber AG;Gehalt;2.500,00\n"

	tab, err := ReadCSV(strings.NewReader(in))
	if err != nil {
		t.Fatalf("ReadCSV: %v", err)
	}
	if tab.Delimiter != ';' {
		t.Errorf("delimiter = %q, want ';'", tab.Delimiter)
	}
	if len(tab.Header) != 4 || tab.Header[0] != "Buchungstag" {
		t.Errorf("unexpected header: %v", tab.Header)
	}
	if len(tab.Rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(tab.Rows))
	}
	if tab.Rows[0][1] != "REWE Markt GmbH" {
		t.Errorf("row[0][1] = %q", tab.Rows[0][1])
	}
}

func TestReadCSVStripsBOMAndDecodesWindows1252(t *testing.T) {
	// UTF-8 BOM, then a Windows-1252 encoded body with ü (0xFC) and € (0x80).
	body := []byte{0xEF, 0xBB, 0xBF}
	body = append(body, []byte("Name;Betrag\n")...)
	body = append(body, []byte("M")...)
	body = append(body, 0xFC) // ü in Windows-1252 / Latin-1
	body = append(body, []byte("ller;100")...)
	body = append(body, 0x80) // € in Windows-1252
	body = append(body, '\n')

	tab, err := ReadCSV(strings.NewReader(string(body)))
	if err != nil {
		t.Fatalf("ReadCSV: %v", err)
	}
	if tab.Header[0] != "Name" {
		t.Errorf("BOM not stripped, header[0] = %q", tab.Header[0])
	}
	if tab.Rows[0][0] != "Müller" {
		t.Errorf("decoded name = %q, want Müller", tab.Rows[0][0])
	}
	if tab.Rows[0][1] != "100€" {
		t.Errorf("decoded amount = %q, want 100€", tab.Rows[0][1])
	}
}

func TestDetectDelimiter(t *testing.T) {
	tests := []struct {
		in   string
		want rune
	}{
		{"a;b;c\n1;2;3", ';'},
		{"a,b,c\n1,2,3", ','},
		{"a\tb\tc\n1\t2\t3", '\t'},
	}
	for _, tt := range tests {
		if got := detectDelimiter(tt.in); got != tt.want {
			t.Errorf("detectDelimiter(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
