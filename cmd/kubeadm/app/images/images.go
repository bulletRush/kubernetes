/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package images

import (
	"fmt"
	"runtime"

	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
)

const (
	KubeEtcdImage = "etcd"

	KubeAPIServerImage         = "apiserver"
	KubeControllerManagerImage = "controller-manager"
	KubeSchedulerImage         = "scheduler"
	KubeProxyImage             = "proxy"

	KubeDNSImage         = "kube-dns"
	KubeDNSmasqImage     = "dnsmasq"
	KubeExechealthzImage = "exechealthz"
	Pause                = "pause"
	KubeDiscoveryImage   = "kube-discovery"
	gcrPrefix            = "gcr.io/google_containers"
	etcdVersion          = "2.2.5"

	kubeDNSVersion     = "1.7"
	dnsmasqVersion     = "1.3"
	exechealthzVersion = "1.1"
	pauseVersion       = "3.0"
)

func GetCoreImageList(cfg *kubeadmapi.MasterConfiguration) map[string]string {
	return map[string]string{
		KubeEtcdImage:              fmt.Sprintf("%s/%s-%s:%s", cfg.ImagePrefix, "etcd", runtime.GOARCH, etcdVersion),
		KubeAPIServerImage:         fmt.Sprintf("%s/%s-%s:%s", cfg.ImagePrefix, "kube-apiserver", runtime.GOARCH, cfg.KubernetesVersion),
		KubeControllerManagerImage: fmt.Sprintf("%s/%s-%s:%s", cfg.ImagePrefix, "kube-controller-manager", runtime.GOARCH, cfg.KubernetesVersion),
		KubeSchedulerImage:         fmt.Sprintf("%s/%s-%s:%s", cfg.ImagePrefix, "kube-scheduler", runtime.GOARCH, cfg.KubernetesVersion),
		KubeProxyImage:             fmt.Sprintf("%s/%s-%s:%s", cfg.ImagePrefix, "kube-proxy", runtime.GOARCH, cfg.KubernetesVersion),
	}
}

func GetCoreImage(image string, cfg *kubeadmapi.MasterConfiguration, overrideImage string) string {
	if overrideImage != "" {
		return overrideImage
	}

	return GetCoreImageList(cfg)[image]
}

func GetAddonImageList(cfg *kubeadmapi.MasterConfiguration) map[string]string {
	return map[string]string{
		KubeDNSImage:         fmt.Sprintf("%s/%s-%s:%s", cfg.ImagePrefix, "kubedns", runtime.GOARCH, kubeDNSVersion),
		KubeDNSmasqImage:     fmt.Sprintf("%s/%s-%s:%s", cfg.ImagePrefix, "kube-dnsmasq", runtime.GOARCH, dnsmasqVersion),
		KubeExechealthzImage: fmt.Sprintf("%s/%s-%s:%s", cfg.ImagePrefix, "exechealthz", runtime.GOARCH, exechealthzVersion),
		KubeDiscoveryImage:   fmt.Sprintf("%s/%s-%s:%s", cfg.ImagePrefix, "kube-discovery", runtime.GOARCH, "1.0"),
		Pause:                fmt.Sprintf("%s/%s-%s:%s", cfg.ImagePrefix, "pause", runtime.GOARCH, pauseVersion),
	}
}

func GetAddonImage(cfg *kubeadmapi.MasterConfiguration, image string) string {
	return GetAddonImageList(cfg)[image]
}
