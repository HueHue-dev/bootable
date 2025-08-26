menuentry "{{.MenuTitle}}" {
    set isofile="{{.ISOPath}}"
    loopback loop "$isofile"
    linux (loop)/casper/vmlinuz.efi boot=live components quiet splash iso-scan/filename=$isofile ---
    initrd (loop)/casper/initrd.gz
}
