```
kubebuilder init --domain x893675.github.io --repo github.com/x893675/demo-crd

kubebuilder edit --multigroup=true

kubebuilder create api --group demo --version v1 --kind Guestbook

kubebuilder create api --group iam --version v1 --kind Token
```