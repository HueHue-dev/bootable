package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"bootable/grub"
	"bootable/helper"
	"bootable/system"
)

var (
	devicePath string
	isoPaths   []string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a bootable USB drive with multiple ISOs",
	Long: `This command will partition, format, and install a multiboot
			bootloader on the specified USB drive, then copy the provided
			ISO files to it.`,
	Run: func(cmd *cobra.Command, args []string) {
		devicePath, _ := cmd.Flags().GetString("device")
		if devicePath == "" {
			fmt.Println("Error: --device flag is required")
			os.Exit(1)
		}

		isoPaths, _ := cmd.Flags().GetStringSlice("isos")
		if len(isoPaths) == 0 {
			fmt.Println("Error: at least one --iso flag is required")
			os.Exit(1)
		}

		fmt.Printf("WARNING: This will erase all data on %s.\n", devicePath)
		fmt.Print("Press 'y' to continue: ")
		var consent string
		fmt.Scanln(&consent)
		if consent != "y" && consent != "Y" {
			fmt.Println("Operation canceled.")
			os.Exit(1)
		}

		fmt.Println("Partitioning and formatting device...")
		if err := system.Format(devicePath, "MULTIBOOT"); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to format device: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Partitioning and formatting completed.")

		part := system.PartitionPath(devicePath, 1)
		mountPoint, err := os.MkdirTemp("", "bootable-mnt-*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create temp mount dir: %v\n", err)
			os.Exit(1)
		}
		defer os.RemoveAll(mountPoint)

		fmt.Printf("Mounting %s to %s...\n", part, mountPoint)
		if err := helper.Run("mount", part, mountPoint); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to mount partition: %v\n", err)
			os.Exit(1)
		}
		defer func() {
			fmt.Println("Unmounting...")
			_ = helper.Run("sync")
			_ = helper.Run("umount", mountPoint)
		}()

		fmt.Println("Installing GRUB bootloader...")
		if err := grub.Install(devicePath, mountPoint); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to install GRUB: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("GRUB installation completed.")

		targetDir := filepath.Join(mountPoint, "ISOs")
		if err := os.MkdirAll(targetDir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create target directory: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Copying ISO files to the USB drive...")
		var copied []string
		for _, iso := range isoPaths {
			dst := filepath.Join(targetDir, filepath.Base(iso))
			fmt.Printf("Copying %s -> %s\n", iso, dst)
			if err := copyFile(iso, dst); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to copy %s: %v\n", iso, err)
				os.Exit(1)
			}
			copied = append(copied, filepath.Join("/ISOs", filepath.Base(iso)))
		}

		if err := grub.WriteConfig(mountPoint, copied); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write GRUB config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Multiboot USB creation completed!")
		os.Exit(0)
	},
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Sync()
		_ = out.Close()
	}()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

func init() {
	var rootCmd = &cobra.Command{Use: "bootable"}

	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&devicePath, "device", "d", "", "Path to the USB device (e.g., /dev/sdb)")
	createCmd.Flags().StringSliceVarP(&isoPaths, "isos", "i", []string{}, "Paths to the ISO files (can be repeated)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
