# Bootable USB Creator

`bootable` is a command-line tool for creating multiboot USB drives with multiple operating systems. It handles partitioning, formatting, GRUB installation, and copies ISO files to the USB, automatically configuring GRUB entries for supported distributions.

## How to Use

**WARNING:** This tool will erase all data on the specified USB device. Ensure you have backed up any important data before proceeding.

1.  **Identify your USB device path:**
    On Linux, this is typically something like `/dev/sdb`, `/dev/sdc`, etc. **Be cautious** to select the correct device, as choosing the wrong one can lead to data loss on your main system. You can use commands like `lsblk` or `fdisk -l` to identify your USB drive.

2.  **Run the `create` command:**

    ```bash
    bootable create --device /dev/sdX --isos /path/to/your/iso1.iso --isos /path/to/your/iso2.iso
    ```

    *   Replace `/dev/sdX` with the actual path to your USB device (e.g., `/dev/sdb`).
    *   Replace `/path/to/your/iso1.iso` and `/path/to/your/iso2.iso` with the full paths to your ISO files. You can specify multiple `--isos` flags.

    The tool will prompt you for confirmation before proceeding.

## Supported ISOs and Architectures

The tool uses a template-based system to automatically configure GRUB for different distributions. Support is determined by internal templates and a mapping logic that tries to detect the distribution and architecture from the ISO filename.

Currently, the following distributions and architectures are explicitly supported:

*   **Arch Linux:**
    *   `x86_64` (e.g., `arch-x86_64.iso`)

*   **Debian / Ubuntu:**
    *   `AMD64` (commonly found as `amd64` or `x86_64` in the filename, e.g., `debian-live-12.5.0-amd64.iso`, `linuxmint-21.3-cinnamon-64bit.iso`)

*   **Fedora / CentOS:**
    *   `x86_64` (e.g., `Fedora-Workstation-Live-x86_64-39.iso`)
    *   `Aarch64` (if `fedora-aarch.tpl` is available and mapped, e.g., `Fedora-Workstation-Live-aarch64-39.iso`)

*   **Manjaro:**
    *   `x86_64` (e.g., `manjaro-kde-23.1.4-240316-linux66.iso`)

*   **Pop!_OS / Ubuntu:**
    *   `x86_64` (e.g., `pop-os_22.04_amd64_nvidia_14.iso`, `ubuntu-22.04.4-desktop-amd64.iso`)

**Note on Detection:** The tool infers the distribution and architecture from the ISO's filename. While common naming conventions are used, variations in ISO filenames might affect detection. If an exact match isn't found, the tool will attempt to use a more generic template or a general fallback template if available.
