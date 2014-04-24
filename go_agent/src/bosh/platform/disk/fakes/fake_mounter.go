package fakes

type FakeMounter struct {
	MountPartitionPaths []string
	MountMountPoints    []string
	MountMountOptions   [][]string
	MountErr            error

	RemountAsReadonlyPath string

	RemountFromMountPoint string
	RemountToMountPoint   string
	RemountMountOptions   []string

	SwapOnPartitionPaths []string

	UnmountPartitionPath string
	UnmountDidUnmount    bool
	UnmountErr           error

	IsMountPointResult bool
	IsMountPointErr    error

	IsMountedDevicePathOrMountPoint string
	IsMountedResult                 bool
	IsMountedErr                    error
}

func (m *FakeMounter) Mount(partitionPath, mountPoint string, mountOptions ...string) error {
	m.MountPartitionPaths = append(m.MountPartitionPaths, partitionPath)
	m.MountMountPoints = append(m.MountMountPoints, mountPoint)
	m.MountMountOptions = append(m.MountMountOptions, mountOptions)
	return m.MountErr
}

func (m *FakeMounter) RemountAsReadonly(mountPoint string) (err error) {
	m.RemountAsReadonlyPath = mountPoint
	return
}

func (m *FakeMounter) Remount(fromMountPoint, toMountPoint string, mountOptions ...string) (err error) {
	m.RemountFromMountPoint = fromMountPoint
	m.RemountToMountPoint = toMountPoint
	m.RemountMountOptions = mountOptions
	return
}

func (m *FakeMounter) SwapOn(partitionPath string) (err error) {
	m.SwapOnPartitionPaths = append(m.SwapOnPartitionPaths, partitionPath)
	return
}

func (m *FakeMounter) Unmount(partitionPath string) (didUnmount bool, err error) {
	m.UnmountPartitionPath = partitionPath
	return m.UnmountDidUnmount, m.UnmountErr
}

func (m *FakeMounter) IsMountPoint(path string) (result bool, err error) {
	return m.IsMountPointResult, m.IsMountPointErr
}

func (m *FakeMounter) IsMounted(devicePathOrMountPoint string) (bool, error) {
	m.IsMountedDevicePathOrMountPoint = devicePathOrMountPoint
	return m.IsMountedResult, m.IsMountedErr
}
