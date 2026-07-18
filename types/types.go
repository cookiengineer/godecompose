// Package types defines shared enumerations and data structures used across
// the godecompose framework.
package types

// Arch represents a CPU architecture.
type Arch int

const (
	ArchUnknown Arch = iota
	ArchX86_64
	ArchAArch64
	ArchRISCv64
)

func (a Arch) String() string {
	switch a {
	case ArchX86_64:
		return "x86_64"
	case ArchAArch64:
		return "aarch64"
	case ArchRISCv64:
		return "riscv64"
	default:
		return "unknown"
	}
}

func ArchFromString(s string) Arch {
	switch s {
	case "x86_64", "amd64", "x64":
		return ArchX86_64
	case "aarch64", "arm64":
		return ArchAArch64
	case "riscv64":
		return ArchRISCv64
	default:
		return ArchUnknown
	}
}

// Platform represents an operating system or kernel.
type Platform int

const (
	PlatformUnknown Platform = iota
	PlatformLinux
	PlatformWindows
	PlatformDarwin
	PlatformFreeBSD
	PlatformOpenBSD
	PlatformNetBSD
)

func (p Platform) String() string {
	switch p {
	case PlatformLinux:
		return "linux"
	case PlatformWindows:
		return "windows"
	case PlatformDarwin:
		return "darwin"
	case PlatformFreeBSD:
		return "freebsd"
	case PlatformOpenBSD:
		return "openbsd"
	case PlatformNetBSD:
		return "netbsd"
	default:
		return "unknown"
	}
}

func PlatformFromString(s string) Platform {
	switch s {
	case "linux":
		return PlatformLinux
	case "windows":
		return PlatformWindows
	case "darwin", "macos", "mac":
		return PlatformDarwin
	case "freebsd":
		return PlatformFreeBSD
	case "openbsd":
		return PlatformOpenBSD
	case "netbsd":
		return PlatformNetBSD
	default:
		return PlatformUnknown
	}
}
