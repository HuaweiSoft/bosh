source $base_dir/lib/prelude_common.bash
source $base_dir/lib/helpers.sh

work=$1
chroot=${chroot:=$work/chroot}
mkdir -p $work $chroot
mkdir -p $chroot/tmp

CUSTOM_YUM_CONF=$base_dir/etc/custom_yum.conf
cp -v ${CUSTOM_YUM_CONF} $chroot/tmp

# Source settings if present
if [ -f $settings_file ]
then
  source $settings_file
fi

# Source /etc/lsb-release if present
if [ -f $chroot/etc/lsb-release ]
then
  source $chroot/etc/lsb-release
fi

function pkg_mgr {
  centos_file=$chroot/etc/centos-release
  ubuntu_file=$chroot/etc/debian_version

  if [ -f $ubuntu_file ]
  then
    echo "Found $ubuntu_file - Assuming Ubuntu"
    run_in_chroot $chroot "apt-get update"
    run_in_chroot $chroot "apt-get -f -y --force-yes --no-install-recommends $*"
    run_in_chroot $chroot "apt-get clean"
  elif [ -f $centos_file ]
  then
    echo "Found $centos_file - Assuming CentOS"
    run_in_chroot $chroot "yum -c /tmp/$(basename ${CUSTOM_YUM_CONF}) update --assumeyes"
    run_in_chroot $chroot "yum -c /tmp/$(basename ${CUSTOM_YUM_CONF}) --verbose --assumeyes $*"
    run_in_chroot $chroot "yum clean all"
  else
    echo "Unknown OS, exiting"
    exit 2
  fi
}

function cleanup_build_artifacts {
 [[ -f $chroot/tmp/$(basename ${CUSTOM_YUM_CONF}) ]] && rm -v $chroot/tmp/$(basename ${CUSTOM_YUM_CONF})
}
