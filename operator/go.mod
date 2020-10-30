module example-solution-operator

go 1.13

require (
	github.com/aws/aws-sdk-go v1.25.48 // indirect
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.17+incompatible // indirect
	github.com/edsrzf/mmap-go v1.0.0 // indirect
	github.com/go-openapi/spec v0.19.4
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/google/pprof v0.0.0-20190723021845-34ac40c74b70 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/operator-framework/operator-sdk v0.17.0
	github.com/spf13/pflag v1.0.5
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20191107075043-30be4d16710a
	sigs.k8s.io/controller-runtime v0.5.2
	sigs.k8s.io/structured-merge-diff v1.0.2 // indirect
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.17.4 // Required by prometheus-operator
)
