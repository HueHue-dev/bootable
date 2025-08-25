menuentry "{{.MenuTitle}}" {
    set isofile="{{.ISOPath}}"
    loopback loop $isofile
    linux (loop)/boot/aarch/loader/vmlinuz root=live:CDLABEL=$isofile quiet rhgb rd.live.image
    initrd (loop)/boot/aarch/loader/initrd.img
}