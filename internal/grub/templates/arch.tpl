menuentry "{{.MenuTitle}}" {
  set isofile="{{.ISOPath}}"
  loopback loop (hd0,msdos1)$isofile
  linux (loop)/arch/boot/vmlinuz-linux tz=UTC
  initrd (loop)/arch/boot/x86_64/initramfs-linux.img
}
