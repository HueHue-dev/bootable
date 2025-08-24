menuentry "{{.MenuTitle}}" {
  set isofile="{{.ISOPath}}"
  loopback loop $isofile
  linux (loop)/live/vmlinuz boot=live config findiso=$isofile
  initrd (loop)/live/initrd.img
}
