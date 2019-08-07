package templates

const cloudcontrollerTpl = `
sudo bash -c 'cat << EOF | kubectl create -f -
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cloud-controller-manager
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: system:cloud-controller-manager
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: cloud-controller-manager
  namespace: kube-system
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  labels:
    k8s-app: cloud-controller-manager
  name: cloud-controller-manager
  namespace: kube-system
spec:
  selector:
    matchLabels:
      k8s-app: cloud-controller-manager
  template:
    metadata:
      labels:
        k8s-app: cloud-controller-manager
    spec:
      serviceAccountName: cloud-controller-manager
      containers:
      - name: cloud-controller-manager
        image: k8s.gcr.io/cloud-controller-manager:v{{ .K8SVersion }}
        command:
        - /usr/local/bin/cloud-controller-manager
        - --cloud-provider=aws
        - --leader-elect=true
        - --use-service-account-credentials
        - --allocate-node-cidrs=true
        - --configure-cloud-routes=true
        - --cluster-cidr={{ .CIDR }}
        volumeMounts:
        - mountPath: /etc/ssl/certs
          name: ssl-certs-host
          readOnly: true
      tolerations:
      - key: node.cloudprovider.kubernetes.io/uninitialized
        value: "true"
        effect: NoSchedule
      - key: "node-role.kubernetes.io/master"
        effect: NoSchedule
      nodeSelector:
        node-role.kubernetes.io/master: ""
      volumes:
      - hostPath:
          path: /etc/ssl/certs
        name: ssl-certs-host
EOF'
`
