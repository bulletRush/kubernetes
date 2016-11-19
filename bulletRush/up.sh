#!/bin/bash
set -o errexit
k8s_path=/data/kubernetes
kubeadm=${k8s_path}/cmd/kubeadm/kubeadm
${kubeadm} reset
systemctl restart kubelet
${kubeadm} init --image-prefix=barrettwu --use-kubernetes-version=dev
kubectl taint nodes --all dedicated-

manifest_path=/data/manifests
kubectl create -f ${manifest_path}/weave-kube
kubectl create -f ${manifest_path}/dashboard/src/deploy/kubernetes-dashboard.yaml

kubectl get pod --all-namespaces
kubectl describe svc/kubernetes-dashboard --namespace=kube-system

