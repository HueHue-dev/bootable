package grub

import (
	"bootable/grub/configurator"
	"bootable/helper"
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed all:templates/*
var grubTemplatesFS embed.FS

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
	builder := configurator.NewCFGBuilder().
		SetMountPoint(mountPoint).
		SetISOPaths(isoPaths).
		SetTemplatesFS(grubTemplatesFS)

	grubConfigurator := configurator.NewGrubConfigurator(builder)

	return grubConfigurator.Construct()
}
