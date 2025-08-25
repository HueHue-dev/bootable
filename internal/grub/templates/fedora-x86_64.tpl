menuentry "{{.MenuTitle}}" {
    set isofile="{{.ISOPath}}"
    loopback loop $isofile
    linux (loop)/boot/x86_64/loader/vmlinuz root=live:CDLABEL=$isofile quiet rhgb rd.live.image
    initrd (loop)/boot/x86_64/loader/initrd.img
}