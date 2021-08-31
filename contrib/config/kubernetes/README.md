# kube
Kubernetes is a service discovery.

### kube集群内部署
集群内部署需要权限
kubectl执行
> serviceaccount 请调整为实际环境account。在未指定spec.serviceAccount情况下默认为namespace::default
```
kubectl create clusterrolebinding go-kratos:kube --clusterrole=view --serviceaccount=mesh:default
```
或者 kubect apply -f bind-role.yaml
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: go-kratos:kube
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: view
subjects:
- kind: ServiceAccount
  name: default
  namespace: mesh
```

### 集群外运行
> 指定 .kube 文件访问
```go
    config.NewSource(SourceOption{
		Namespace:     "mesh",
		LabelSelector: "",
		KubeConfig:    filepath.Join(homedir.HomeDir(), ".kube", "config"),
	})
```