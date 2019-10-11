#!/usr/bin/env bash

# true if tuned is installed
tuned () {
  command -v tuned-adm &> /dev/null
  [ $? == 0 ];
}

# true if there is an active profile
tuned_active () {
  tuned-adm active | grep ':' &> /dev/null
  [ $? == 0 ];
}

tuned_list () {
  profiles=( $(tuned-adm list | sed -n '1!p' | sed '$d' | awk -F'- '  '{print $2}' | tr -d ' ') )
}

array_contains () {
    local seeking=$1; shift
    local in=1
    for element; do
        if [[ $element == $seeking ]]; then
            in=0
            break
        fi
    done
    return $in
}

tuned_preferred_profile () {
  preferences=( "openshift-node" "throughput-performance" "network-latency" )
  for preference in "${preferences[@]}"
  do
    if array_contains $preference "${profiles[@]}" ; then
      profile=$preference
      break
    fi
  done
}

enable_tuned_nuodb () {
  tuned-adm profile nuodb &> /dev/null
  systemctl restart tuned &> /dev/null
}

if tuned and tuned_active; then

  tuned_list
  tuned_preferred_profile

  mkdir /usr/lib/tuned/nuodb
  cat > /usr/lib/tuned/nuodb/tuned.conf <<EOF
[main]
summary=Optimize for NuoDB Distributed RDBMS
include=${profile}

[sysctl]
vm.swappiness = 1
vm.dirty_background_ratio = 3
vm.dirty_ratio = 80
vm.dirty_expire_centisecs = 500
vm.dirty_writeback_centisecs = 100
kernel.shmmax = 4398046511104
kernel.shmall = 1073741824
kernel.shmmni = 4096
kernel.sem = 250 32000 100 128
fs.file-max = 6815744
fs.aio-max-nr = 1048576
net.ipv4.ip_local_port_range = 9000 65500
net.core.rmem_default = 262144
net.core.rmem_max = 4194304
net.core.wmem_default = 262144
net.core.wmem_max = 1048576
kernel.panic_on_oops = 1

[vm]
transparent_hugepages=never
EOF

  chmod 0644 /usr/lib/tuned/nuodb/tuned.conf
  enable_tuned_nuodb

else

  if [ -d /sys/kernel/mm/transparent_hugepage ]; then
    thp_path=/sys/kernel/mm/transparent_hugepage
  elif [ -d /sys/kernel/mm/redhat_transparent_hugepage ]; then
    thp_path=/sys/kernel/mm/redhat_transparent_hugepage
  else
    return 0
  fi

  echo 'never' > ${thp_path}/enabled
  echo 'never' > ${thp_path}/defrag

  unset thp_path

fi
