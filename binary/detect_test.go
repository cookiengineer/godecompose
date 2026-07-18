package binary

import "testing"

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name   string
		magic  []byte
		format Format
	}{
		{"ELF", []byte{0x7F, 'E', 'L', 'F'}, FormatELF},
		{"PE", []byte{'M', 'Z', 0, 0}, FormatPE},
		{"MachO_LE_32", []byte{0xCE, 0xFA, 0xED, 0xFE}, FormatMachO},
		{"MachO_BE_32", []byte{0xFE, 0xED, 0xFA, 0xCE}, FormatMachO},
		{"MachO_LE_64", []byte{0xCF, 0xFA, 0xED, 0xFE}, FormatMachO},
		{"MachO_BE_64", []byte{0xFE, 0xED, 0xFA, 0xCF}, FormatMachO},
		{"FatBE", []byte{0xCA, 0xFE, 0xBA, 0xBE}, FormatMachO},
		{"FatLE", []byte{0xBE, 0xBA, 0xFE, 0xCA}, FormatMachO},
		{"Unknown", []byte{0x00, 0x00, 0x00, 0x00}, FormatUnknown},
		{"TooShort", []byte{0x7F}, FormatUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectFormat(tt.magic)
			if got != tt.format {
				t.Errorf("detectFormat(%x) = %v, want %v", tt.magic, got, tt.format)
			}
		})
	}
}

func TestFormatString(t *testing.T) {
	tests := []struct {
		format Format
		want   string
	}{
		{FormatUnknown, "unknown"},
		{FormatELF, "ELF"},
		{FormatPE, "PE"},
		{FormatMachO, "Mach-O"},
	}

	for _, tt := range tests {
		got := tt.format.String()
		if got != tt.want {
			t.Errorf("Format(%d).String() = %q, want %q", tt.format, got, tt.want)
		}
	}
}
