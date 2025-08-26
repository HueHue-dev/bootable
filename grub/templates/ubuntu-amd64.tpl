menuentry "{{.MenuTitle}}" {
    set isofile="{{.ISOPath}}"
    loopback loop "$isofile"
    linux (loop)/casper/vmlinuz noprompt noeject iso-scan/filename=$isofile ---
    initrd (loop)/casper/initrd
}
