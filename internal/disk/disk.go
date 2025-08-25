package disk

import (
	"fmt"
	"os/exec"
	"time"

	"bootable/internal/helper"
)

func Format(devicePath, label string) error {
	_ = exec.Command("umount", "-f", PartitionPath(devicePath, 1)).Run()

	if err := helper.Run("wipefs", "-a", devicePath); err != nil {
		return fmt.Errorf("wipefs failed: %w", err)
	}

	sfdiskInput := ",,c,*\n"
	if err := helper.RunWithStdin(sfdiskInput, "sfdisk", devicePath); err != nil {
		return fmt.Errorf("sfdisk failed: %w", err)
	}

	if err := helper.Run("partprobe", devicePath); err != nil {
		_ = helper.Run("blockdev", "--rereadpt", devicePath)
	}

	time.Sleep(1 * time.Second)

	part := PartitionPath(devicePath, 1)

	if err := helper.Run("mkfs.vfat", "-F32", "-n", label, part); err != nil {
		return fmt.Errorf("mkfs.vfat failed on %s: %w", part, err)
	}

	return nil
}

func PartitionPath(dev string, n int) string {
	if dev == "" {
		return ""
	}
	if isLastCharDigit(dev) {
		return fmt.Sprintf("%sp%d", dev, n)
	}

	return fmt.Sprintf("%s%d", dev, n)
}

func isLastCharDigit(s string) bool {
	if s == "" {
		return false
	}
	last := s[len(s)-1]

	return last >= '0' && last <= '9'
}
