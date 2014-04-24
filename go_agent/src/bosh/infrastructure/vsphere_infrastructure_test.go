package infrastructure_test

import (
	"errors"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh/infrastructure"
	fakedpresolv "bosh/infrastructure/devicepathresolver/fakes"
	boshlog "bosh/logger"
	boshdisk "bosh/platform/disk"
	fakeplatform "bosh/platform/fakes"
	boshsettings "bosh/settings"
	fakesys "bosh/system/fakes"
)

var _ = Describe("vSphere Infrastructure", func() {
	var (
		logger             boshlog.Logger
		vsphere            Infrastructure
		platform           *fakeplatform.FakePlatform
		devicePathResolver *fakedpresolv.FakeDevicePathResolver
	)

	BeforeEach(func() {
		platform = fakeplatform.NewFakePlatform()
		devicePathResolver = fakedpresolv.NewFakeDevicePathResolver()
		logger = boshlog.NewLogger(boshlog.LevelNone)
	})

	JustBeforeEach(func() {
		vsphere = NewVsphereInfrastructure(platform, devicePathResolver, logger)
	})

	Describe("GetSettings", func() {
		It("vsphere get settings", func() {
			platform.GetFileContentsFromCDROMContents = []byte(`{"agent_id": "123"}`)

			settings, err := vsphere.GetSettings()
			Expect(err).NotTo(HaveOccurred())

			Expect(platform.GetFileContentsFromCDROMPath).To(Equal("env"))
			Expect(settings.AgentID).To(Equal("123"))
		})
	})

	Describe("SetupNetworking", func() {
		It("vsphere setup networking", func() {
			networks := boshsettings.Networks{"bosh": boshsettings.Network{}}

			vsphere.SetupNetworking(networks)

			Expect(platform.SetupManualNetworkingNetworks).To(Equal(networks))
		})
	})

	Describe("GetEphemeralDiskPath", func() {
		It("vsphere get ephemeral disk path", func() {
			realPath, found := vsphere.GetEphemeralDiskPath("does not matter")
			Expect(found).To(Equal(true))

			Expect(realPath).To(Equal("/dev/sdb"))
		})
	})

	Describe("MountPersistentDisk", func() {
		BeforeEach(func() {
			devicePathResolver.RegisterRealDevicePath("fake-volume-id", "fake-real-device-path")
		})

		It("creates the mount directory with the correct permissions", func() {
			err := vsphere.MountPersistentDisk("fake-volume-id", "/mnt/point")
			Expect(err).ToNot(HaveOccurred())

			mountPoint := platform.Fs.GetFileTestStat("/mnt/point")
			Expect(mountPoint.FileType).To(Equal(fakesys.FakeFileTypeDir))
			Expect(mountPoint.FileMode).To(Equal(os.FileMode(0700)))
		})

		It("returns error when creating mount directory fails", func() {
			platform.Fs.MkdirAllError = errors.New("fake-mkdir-all-err")

			err := vsphere.MountPersistentDisk("fake-volume-id", "/mnt/point")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-mkdir-all-err"))
		})

		It("partitions the disk", func() {
			err := vsphere.MountPersistentDisk("fake-volume-id", "/mnt/point")
			Expect(err).ToNot(HaveOccurred())

			Expect(platform.FakeDiskManager.FakePartitioner.PartitionDevicePath).To(Equal("fake-real-device-path"))
			partitions := []boshdisk.Partition{
				{Type: boshdisk.PartitionTypeLinux},
			}
			Expect(platform.FakeDiskManager.FakePartitioner.PartitionPartitions).To(Equal(partitions))
		})

		It("formats the disk", func() {
			err := vsphere.MountPersistentDisk("fake-volume-id", "/mnt/point")
			Expect(err).ToNot(HaveOccurred())

			Expect(platform.FakeDiskManager.FakeFormatter.FormatPartitionPaths).To(Equal([]string{"fake-real-device-path1"}))
			Expect(platform.FakeDiskManager.FakeFormatter.FormatFsTypes).To(Equal([]boshdisk.FileSystemType{boshdisk.FileSystemExt4}))
		})

		It("mounts the disk", func() {
			err := vsphere.MountPersistentDisk("fake-volume-id", "/mnt/point")
			Expect(err).ToNot(HaveOccurred())

			Expect(platform.FakeDiskManager.FakeMounter.MountPartitionPaths).To(Equal([]string{"fake-real-device-path1"}))
			Expect(platform.FakeDiskManager.FakeMounter.MountMountPoints).To(Equal([]string{"/mnt/point"}))
		})
	})
})
