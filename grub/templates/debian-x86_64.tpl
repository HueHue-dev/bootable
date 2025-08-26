menuentry "{{.MenuTitle}}" {
    set isofile="{{.ISOPath}}"
    loopback loop "$isofile"
    linux (loop)/install.amd/vmlinuz boot=live components quiet splash iso-scan/filename=$isofile ---
    initrd (loop)/install.amd/initrd.gz
}
