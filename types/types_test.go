package types

import "testing"

func TestArchString(t *testing.T) {
	tests := []struct {
		arch Arch
		want string
	}{
		{ArchUnknown, "unknown"},
		{ArchX86_64, "x86_64"},
		{ArchAArch64, "aarch64"},
		{ArchRISCv64, "riscv64"},
	}

	for _, tt := range tests {
		got := tt.arch.String()
		if got != tt.want {
			t.Errorf("Arch(%d).String() = %q, want %q", tt.arch, got, tt.want)
		}
	}
}

func TestArchFromString(t *testing.T) {
	tests := []struct {
		input string
		want  Arch
	}{
		{"x86_64", ArchX86_64},
		{"amd64", ArchX86_64},
		{"x64", ArchX86_64},
		{"aarch64", ArchAArch64},
		{"arm64", ArchAArch64},
		{"riscv64", ArchRISCv64},
		{"mips", ArchUnknown},
		{"", ArchUnknown},
	}

	for _, tt := range tests {
		got := ArchFromString(tt.input)
		if got != tt.want {
			t.Errorf("ArchFromString(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestPlatformString(t *testing.T) {
	tests := []struct {
		platform Platform
		want     string
	}{
		{PlatformUnknown, "unknown"},
		{PlatformLinux, "linux"},
		{PlatformWindows, "windows"},
		{PlatformDarwin, "darwin"},
		{PlatformFreeBSD, "freebsd"},
		{PlatformOpenBSD, "openbsd"},
		{PlatformNetBSD, "netbsd"},
	}

	for _, tt := range tests {
		got := tt.platform.String()
		if got != tt.want {
			t.Errorf("Platform(%d).String() = %q, want %q", tt.platform, got, tt.want)
		}
	}
}

func TestPlatformFromString(t *testing.T) {
	tests := []struct {
		input string
		want  Platform
	}{
		{"linux", PlatformLinux},
		{"windows", PlatformWindows},
		{"darwin", PlatformDarwin},
		{"macos", PlatformDarwin},
		{"mac", PlatformDarwin},
		{"freebsd", PlatformFreeBSD},
		{"openbsd", PlatformOpenBSD},
		{"netbsd", PlatformNetBSD},
		{"solaris", PlatformUnknown},
		{"", PlatformUnknown},
	}

	for _, tt := range tests {
		got := PlatformFromString(tt.input)
		if got != tt.want {
			t.Errorf("PlatformFromString(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
