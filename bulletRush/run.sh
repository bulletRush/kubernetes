#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

KUBE_VERBOSE=10
BASE_PATH="$(pwd)"
OUTPUT="$BASE_PATH/_tmp"
OUTPUT_BIN="$OUTPUT/bin"

# Print a status line.  Formatted to show up in a stream of output.
kube::log::status() {
  local V="${V:-0}"
  if [[ $KUBE_VERBOSE < $V ]]; then
    return
  fi

  timestamp=$(date +"[%m%d %H:%M:%S]")
  echo "+++ $timestamp $1"
  shift
  for message; do
    echo "    $message"
  done
}


# Wait for background jobs to finish. Return with
# an error status if any of the jobs failed.
kube::util::wait-for-jobs() {
  local fail=0
  local job
  for job in $(jobs -p); do
    wait "${job}" || fail=$((fail + 1))
  done
  return ${fail}
}


kube::build::build_and_copy_binaries() {
    local targets=(
        kube-apiserver
        kube-controller-manager
        kube-dns
        kube-proxy
        kubeadm
        kubectl
        kubelet
    )
    mkdir -p "$OUTPUT_BIN"
    for component in "${targets[@]}"; do
         component_path="$BASE_PATH/cmd/$component"
         kube::log::status "build: $component_path"
         cd "$component_path"
         go build -i -v
         cp "${component_path}/${component}" "${OUTPUT_BIN}/${component}"
    done
}

kube::build::build_kubernetes_images() {
    local targets=(
        kube-apiserver,busybox
        kube-controller-manager,busybox
        # kube-scheduler,busybox
        kube-proxy,barrettwu/debian-iptables-amd64:v4
    )
    local arch="amd64"
    for component in "${targets[@]}"; do
        local oldifs=$IFS
        IFS=","
        set $component
        IFS=$oldifs
        local binary_name="$1"
        local base_image="$2"
        kube::log::status "Starting Docker build for image: ${binary_name}, based on: ${base_image}"
        (
            local docker_build_path="${OUTPUT}/${binary_name}.dockerbuild"
            local docker_file_path="${docker_build_path}/Dockerfile"
            local binary_file_path="${OUTPUT_BIN}/${binary_name}"
            mkdir -p ${docker_build_path}
            cp ${binary_file_path} ${docker_build_path}/${binary_name}
            printf " FROM ${base_image} \n ADD ${binary_name} /usr/local/bin/${binary_name}\n" > ${docker_file_path}
            local docker_image_tag=barrettwu/${binary_name}-${arch}:dev
            docker build -q -t "${docker_image_tag}" ${docker_build_path}
            echo "Build end: ${docker_image_tag}"
        ) &
    done
    kube::util::wait-for-jobs || { echo "err happened"; return 1; }
    kube::log::status "Docker builds done"
}

# kube::build::build_and_copy_binaries
kube::build::build_kubernetes_images