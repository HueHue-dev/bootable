package grub

import (
	"embed"
	"fmt"
	"io/fs"
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
	file            *os.File
	grubCfgPath     string
	templateMap     map[string]*template.Template
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
	b.file = f
	return nil
}

func (b *CFGBuilder) insertHeaderTemplate() error {
	headerTemplate, err := getRawTemplate("header.tpl", b.grubTemplatesFS)
	if err != nil {
		return fmt.Errorf("error getting header template: %w", err)
	}

	_, err = b.file.WriteString(headerTemplate + "\n")
	return err
}

func (b *CFGBuilder) loadTemplates() error {
	b.templateMap = make(map[string]*template.Template)
	err := fs.WalkDir(b.grubTemplatesFS, "templates", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".tpl") {
			content, err := b.grubTemplatesFS.ReadFile(path)
			if err != nil {
				return err
			}
			tmplName := strings.TrimSuffix(d.Name(), ".tpl")
			tmpl, err := template.New(tmplName).Parse(string(content))
			if err != nil {
				return err
			}
			b.templateMap[tmplName] = tmpl
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}
	return nil
}

func (b *CFGBuilder) insertIsoSpecificTemplates() error {
	if b.templateMap == nil {
		if err := b.loadTemplates(); err != nil {
			return err
		}
	}

	for _, iso := range b.isoPaths {
		base := filepath.Base(iso)
		title := fmt.Sprintf("Boot ISO: %s", base)

		tmpl := pickTemplateForISO(strings.ToLower(base), b.templateMap)
		if tmpl == nil {
			return fmt.Errorf("no template found for ISO: %s", iso)
		}

		var entry strings.Builder
		data := EntryData{
			MenuTitle: title,
			ISOPath:   iso,
		}

		if err := tmpl.Execute(&entry, data); err != nil {
			return fmt.Errorf("failed to execute template for %s: %w", iso, err)
		}

		_, err := b.file.WriteString(entry.String() + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *CFGBuilder) GetResult() error {
	if b.file != nil {
		defer b.file.Close()
		return b.file.Sync()
	}
	return fmt.Errorf("grub.cfg file was not created or processed")
}
