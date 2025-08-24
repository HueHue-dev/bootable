menuentry "{{.MenuTitle}}" {
  set isofile="{{.ISOPath}}"
  loopback loop (hd0,msdos1)$isofile
  linux (loop)/manjaro/boot/vmlinuz-x86_64 tz=UTC
  initrd (loop)/manjaro/boot/initramfs-x86_64.img
}
