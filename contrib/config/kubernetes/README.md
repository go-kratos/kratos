# Kubernetes Config

### Usage in the Kubernetes Cluster
It is required to 
> serviceaccount should be set to the actual account of your environment, the default account will be `namespace::default` if the `spec.serviceAccount` is unset. 
execute this command:
```
kubectl create clusterrolebinding go-kratos:kube --clusterrole=view --serviceaccount=mesh:default
```
or use `kubectl apply -f bind-role.yaml`
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

### Usage outside the Kubernetes Cluster
Set the path `~/.kube/config` to KubeConfig
```go
    config.NewSource(SourceOption{
		Namespace:     "mesh",
		LabelSelector: "",
		KubeConfig:    filepath.Join(homedir.HomeDir(), ".kube", "config"),
	})
```