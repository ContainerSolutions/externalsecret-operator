module github.com/containersolutions/externalsecret-operator

go 1.15

require (
	cloud.google.com/go v0.66.0
	github.com/ContainerSolutions/onepassword v0.1.0
	github.com/aws/aws-sdk-go v1.34.29
	github.com/go-logr/logr v0.2.1
	github.com/go-logr/zapr v0.2.0 // indirect
	github.com/googleapis/gax-go v1.0.3
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.13.0
	github.com/smartystreets/goconvey v1.6.4
	github.com/xanzy/go-gitlab v0.39.0
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	google.golang.org/api v0.32.0
	google.golang.org/genproto v0.0.0-20200921165018-b9da36f5f452
	google.golang.org/grpc v1.31.1
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.6.3
)
