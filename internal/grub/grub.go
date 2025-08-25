package grub

import (
	"bootable/internal/helper"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed all:templates/*
var grubTemplatesFS embed.FS

type EntryData struct {
	MenuTitle string
	ISOPath   string
}

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

type TemplateKey struct {
	Distro DistroType
	Arch   ArchitectureType
}

func Install(devicePath, mountPoint string) error {
	bootDir := filepath.Join(mountPoint, "boot")
	if err := os.MkdirAll(filepath.Join(bootDir, "grub"), 0o755); err != nil {
		return err
	}

	if err := helper.Run(
		"grub-install",
		"--target=i386-pc",
		"--boot-directory", bootDir,
		"--recheck",
		devicePath,
	); err != nil {
		return fmt.Errorf("BIOS grub-install failed: %w", err)
	}

	efiDir := mountPoint
	if err := os.MkdirAll(filepath.Join(efiDir, "EFI"), 0o755); err != nil {
		return err
	}
	if err := helper.Run(
		"grub-install",
		"--target=x86_64-efi",
		"--efi-directory", efiDir,
		"--boot-directory", bootDir,
		"--removable",
	); err != nil {
		return fmt.Errorf("UEFI grub-install failed: %w", err)
	}

	return nil
}
func WriteConfig(mountPoint string, isoPaths []string) error {
	builder := NewCFGBuilder().
		SetMountPoint(mountPoint).
		SetISOPaths(isoPaths).
		SetTemplatesFS(grubTemplatesFS)

	configurator := NewGrubConfigurator(builder)
	return configurator.Construct()
}

func detectArchitecture(lowerBase string) ArchitectureType {
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
