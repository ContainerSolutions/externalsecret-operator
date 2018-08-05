package main

import (
	"context"
	"runtime"

	"github.com/ContainerSolutions/externalconfig-operator/pkg/secrets"
	stub "github.com/ContainerSolutions/externalconfig-operator/pkg/stub"
	sdk "github.com/operator-framework/operator-sdk/pkg/sdk"
	k8sutil "github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	sdkVersion "github.com/operator-framework/operator-sdk/version"

	"github.com/sirupsen/logrus"
)

func printVersion() {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func main() {
	var backendKey secrets.ContextKey = "backend"
	printVersion()

	sdk.ExposeMetricsPort()

	resource := "externalconfig-operator.container-solutions.com/v1alpha1"
	kind := "ExternalConfig"
	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		logrus.Fatalf("Failed to get watch namespace: %v", err)
	}
	backend := secrets.NewAWSSecretsManagerBackend()
	if err != nil {
		logrus.Fatalf("Failed to initialize the secrets backend: %v", err)
	}
	ctx := context.WithValue(context.Background(), backendKey, backend)
	resyncPeriod := 5
	logrus.Infof("Watching %s, %s, %s, %d", resource, kind, namespace, resyncPeriod)
	sdk.Watch(resource, kind, namespace, resyncPeriod)
	sdk.Handle(stub.NewHandler())
	sdk.Run(ctx)
}
