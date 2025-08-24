package grub

import (
	"bootable/internal/helper"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type EntryData struct {
	MenuTitle string
	ISOPath   string
}

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
	builder := NewCFGBuilder().
		SetMountPoint(mountPoint).
		SetISOPaths(isoPaths).
		SetTemplatesFS(grubTemplatesFS)

	configurator := NewGrubConfigurator(builder)
	return configurator.Construct()
}

func pickTemplateForISO(lowerBase string, templateMap map[string]*template.Template) *template.Template {
	switch {
	case strings.Contains(lowerBase, "ubuntu"),
		strings.Contains(lowerBase, "debian"),
		strings.Contains(lowerBase, "linuxmint"),
		strings.Contains(lowerBase, "elementary"):
		return templateMap["debian"]

	case strings.Contains(lowerBase, "pop-os"):
		return templateMap["pop-os"]

	case strings.Contains(lowerBase, "arch"),
		strings.Contains(lowerBase, "archlinux"),
		strings.Contains(lowerBase, "endeavouros"),
		strings.Contains(lowerBase, "manjaro"):
		return templateMap["arch"]

	case strings.Contains(lowerBase, "fedora"),
		strings.Contains(lowerBase, "centos"),
		strings.Contains(lowerBase, "rhel"),
		strings.Contains(lowerBase, "rocky"),
		strings.Contains(lowerBase, "almalinux"):
		return templateMap["fedora"]

	case strings.Contains(lowerBase, "opensuse"),
		strings.Contains(lowerBase, "tumbleweed"),
		strings.Contains(lowerBase, "leap"),
		strings.Contains(lowerBase, "suse"):
		return templateMap["opensuse"]

	case strings.Contains(lowerBase, "alpine"):
		return templateMap["alpine"]

	default:
		return templateMap["generic"]
	}
}

func getRawTemplate(templateName string, templatesFS embed.FS) (string, error) {
	filePath := "templates/" + templateName

	content, err := templatesFS.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	return string(content), nil
}
