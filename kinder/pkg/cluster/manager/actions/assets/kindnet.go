/*
Copyright 2021 The Kubernetes Authors.

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

package assets

// KindnetImage054 is the image for kindnet 0.5.4
const KindnetImage054 = "kindest/kindnetd:0.5.4"

// KindnetManifest054 holds the kindnet manifest for 0.5.4
const KindnetManifest054 = `
# kindnetd networking manifest
# would you kindly template this file
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kindnet
rules:
  - apiGroups:
      - policy
    resources:
      - podsecuritypolicies
    verbs:
      - use
    resourceNames:
      - kindnet
  - apiGroups:
      - ""
    resources:
      - nodes
    verbs:
      - list
      - watch
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kindnet
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kindnet
subjects:
  - kind: ServiceAccount
    name: kindnet
    namespace: kube-system
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kindnet
  namespace: kube-system
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: kindnet
  namespace: kube-system
  labels:
    tier: node
    app: kindnet
    k8s-app: kindnet
spec:
  selector:
    matchLabels:
      app: kindnet
  template:
    metadata:
      labels:
        tier: node
        app: kindnet
        k8s-app: kindnet
    spec:
      hostNetwork: true
      tolerations:
        - operator: Exists
          effect: NoSchedule
      serviceAccountName: kindnet
      containers:
        - name: kindnet-cni
          image: kindest/kindnetd:0.5.4
          env:
            - name: HOST_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.hostIP
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: POD_SUBNET
              value: "192.168.0.0/16"
          volumeMounts:
            - name: cni-cfg
              mountPath: /etc/cni/net.d
            - name: xtables-lock
              mountPath: /run/xtables.lock
              readOnly: false
            - name: lib-modules
              mountPath: /lib/modules
              readOnly: true
          resources:
            requests:
              cpu: "100m"
              memory: "50Mi"
            limits:
              cpu: "100m"
              memory: "50Mi"
          securityContext:
            privileged: false
            capabilities:
              add: ["NET_RAW", "NET_ADMIN"]
      volumes:
        - name: cni-cfg
          hostPath:
            path: /etc/cni/net.d
        - name: xtables-lock
          hostPath:
            path: /run/xtables.lock
            type: FileOrCreate
        - name: lib-modules
          hostPath:
            path: /lib/modules
`
