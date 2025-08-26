package system

import "strings"

type DistroType string

const (
	DistroArch      DistroType = "arch"
	DistroDebian    DistroType = "debian"
	DistroFedora    DistroType = "fedora"
	DistroManjaro   DistroType = "manjaro"
	DistroPopOS     DistroType = "pop-os"
	DistroUbuntu    DistroType = "ubuntu"
	DistroLinuxMint DistroType = "linuxmint"
	DistroCentOS    DistroType = "centos"
	DistroGeneric   DistroType = "generic"
)

type ArchitectureType string

const (
	ArchX86_64  ArchitectureType = "x86_64"
	ArchAMD64   ArchitectureType = "amd64"
	ArchI386    ArchitectureType = "i386"
	ArchAarch64 ArchitectureType = "aarch64"
	ArchUnknown ArchitectureType = "Unknown"
)

func DetectArchitecture(lowerBase string) ArchitectureType {
	if strings.Contains(lowerBase, "x86_64") || strings.Contains(lowerBase, "amd64") {
		return ArchX86_64
	}
	if strings.Contains(lowerBase, "i386") || strings.Contains(lowerBase, "x86") {
		return ArchI386
	}
	if strings.Contains(lowerBase, "aarch") || strings.Contains(lowerBase, "arm64") {
		return ArchAarch64
	}

	return ArchUnknown
}
