package configurator

import (
	"bootable/system"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type ICFGBuilder interface {
	SetMountPoint(mountPoint string) ICFGBuilder
	SetISOPaths(isoPaths []string) ICFGBuilder
	SetTemplatesFS(grubTemplatesFS embed.FS) ICFGBuilder
	createGrubCfgFile() error
	insertHeaderTemplate() error
	insertIsoSpecificTemplates() error
	GetResult() error
}

type CFGBuilder struct {
	mountPoint      string
	isoPaths        []string
	grubTemplatesFS embed.FS
	cfgFile         *os.File
	grubCfgPath     string
}

type TemplateKey struct {
	Distro system.DistroType
	Arch   system.ArchitectureType
}

type EntryData struct {
	MenuTitle string
	ISOPath   string
}

func NewCFGBuilder() *CFGBuilder {
	return &CFGBuilder{}
}

func (b *CFGBuilder) SetMountPoint(mountPoint string) ICFGBuilder {
	b.mountPoint = mountPoint
	b.grubCfgPath = filepath.Join(mountPoint, "boot", "grub", "grub.cfg")
	return b
}

func (b *CFGBuilder) SetISOPaths(isoPaths []string) ICFGBuilder {
	b.isoPaths = isoPaths
	return b
}

func (b *CFGBuilder) SetTemplatesFS(grubTemplatesFS embed.FS) ICFGBuilder {
	b.grubTemplatesFS = grubTemplatesFS
	return b
}

func (b *CFGBuilder) createGrubCfgFile() error {
	if err := os.MkdirAll(filepath.Dir(b.grubCfgPath), 0o755); err != nil {
		return err
	}

	f, err := os.Create(b.grubCfgPath)
	if err != nil {
		return err
	}
	b.cfgFile = f
	return nil
}

func (b *CFGBuilder) insertHeaderTemplate() error {
	headerTemplate, err := b.getRawTemplate("header.tpl")
	if err != nil {
		return fmt.Errorf("error getting header template: %w", err)
	}

	_, err = b.cfgFile.WriteString(headerTemplate + "\n")
	return err
}

func (b *CFGBuilder) loadTemplate(tmplKey TemplateKey) *template.Template {
	tmplName := string(tmplKey.Distro) + "-" + string(tmplKey.Arch)

	tmplContent, err := b.getRawTemplate(tmplName + ".tpl")
	if err != nil {
		panic(fmt.Sprintf("Failed to load any GRUB template: %v. This should not happen.", tmplName))
	}
	tmpl, err := template.New(tmplName).Parse(tmplContent)

	return tmpl
}

func (b *CFGBuilder) insertIsoSpecificTemplates() error {
	for _, iso := range b.isoPaths {
		base := filepath.Base(iso)
		title := fmt.Sprintf("Boot ISO: %s", base)

		templateKey := b.getTemplateKeyFromIso(strings.ToLower(base))

		tmpl := b.loadTemplate(templateKey)

		var entry strings.Builder
		data := EntryData{
			MenuTitle: title,
			ISOPath:   iso,
		}

		if err := tmpl.Execute(&entry, data); err != nil {
			return fmt.Errorf("failed to execute template for %s: %w", iso, err)
		}

		_, err := b.cfgFile.WriteString(entry.String() + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *CFGBuilder) GetResult() error {
	if b.cfgFile != nil {
		defer b.cfgFile.Close()
		return b.cfgFile.Sync()
	}
	return fmt.Errorf("grub.cfg file was not created or processed")
}

func (b *CFGBuilder) getRawTemplate(templateName string) (string, error) {
	filePath := "templates/" + templateName

	content, err := b.grubTemplatesFS.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	return string(content), nil
}

func (b *CFGBuilder) getTemplateKeyFromIso(lowerBase string) TemplateKey {
	switch {
	case strings.Contains(lowerBase, string(system.DistroDebian)):
		archKey := system.DetectArchitecture(lowerBase)
		templateKey := TemplateKey{Distro: system.DistroDebian, Arch: archKey}
		return templateKey

	case strings.Contains(lowerBase, string(system.DistroPopOS)):
		archKey := system.DetectArchitecture(lowerBase)
		templateKey := TemplateKey{Distro: system.DistroPopOS, Arch: archKey}
		return templateKey

	case strings.Contains(lowerBase, string(system.DistroUbuntu)):
		archKey := system.DetectArchitecture(lowerBase)
		templateKey := TemplateKey{Distro: system.DistroDebian, Arch: archKey}
		return templateKey

	case strings.Contains(lowerBase, string(system.DistroManjaro)):
		archKey := system.DetectArchitecture(lowerBase)
		templateKey := TemplateKey{Distro: system.DistroManjaro, Arch: archKey}
		return templateKey

	case strings.Contains(lowerBase, string(system.DistroArch)):
		archKey := system.DetectArchitecture(lowerBase)
		templateKey := TemplateKey{Distro: system.DistroArch, Arch: archKey}
		return templateKey

	case strings.Contains(lowerBase, string(system.DistroFedora)):
		archKey := system.DetectArchitecture(lowerBase)
		templateKey := TemplateKey{Distro: system.DistroFedora, Arch: archKey}
		return templateKey

	default:
		return TemplateKey{Distro: system.DistroGeneric, Arch: system.ArchUnknown}
	}
}
