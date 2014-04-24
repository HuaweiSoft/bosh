package disk

import ()

type linuxBindMounter struct {
	delegateMounter Mounter
}

func NewLinuxBindMounter(delegateMounter Mounter) linuxBindMounter {
	return linuxBindMounter{delegateMounter}
}

func (m linuxBindMounter) Mount(partitionPath, mountPoint string, mountOptions ...string) error {
	mountOptions = append(mountOptions, "--bind")
	return m.delegateMounter.Mount(partitionPath, mountPoint, mountOptions...)
}

func (m linuxBindMounter) RemountAsReadonly(mountPoint string) error {
	// Remounting mount points mounted originally by warden with '-r ro --bind' flags does not work
	return nil
}

func (m linuxBindMounter) Remount(fromMountPoint, toMountPoint string, mountOptions ...string) error {
	return m.delegateMounter.Remount(fromMountPoint, toMountPoint, mountOptions...)
}

func (m linuxBindMounter) SwapOn(partitionPath string) (err error) {
	return m.delegateMounter.SwapOn(partitionPath)
}

func (m linuxBindMounter) Unmount(partitionOrMountPoint string) (bool, error) {
	return m.delegateMounter.Unmount(partitionOrMountPoint)
}

func (m linuxBindMounter) IsMountPoint(path string) (bool, error) {
	return m.delegateMounter.IsMountPoint(path)
}

func (m linuxBindMounter) IsMounted(partitionOrMountPoint string) (bool, error) {
	return m.delegateMounter.IsMounted(partitionOrMountPoint)
}
