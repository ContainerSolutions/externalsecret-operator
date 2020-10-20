package controllers

import (
	"context"
	"time"

	storev1alpha1 "github.com/containersolutions/externalsecret-operator/apis/store/v1alpha1"
	"github.com/containersolutions/externalsecret-operator/pkg/backend"
	"github.com/containersolutions/externalsecret-operator/pkg/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

const SecretStoreNamespace = "default"

var _ = Describe("SecretstoreController", func() {
	var (
		SecretStoreName           = "externalsecret-operator-test"
		SecretStoreControllerName = "test-store-ctrl"
		KeyName                   = "test-secret"
		KeyVersion                = "test-version"
		CredentialSecretName      = "credential-secret"

		timeout = time.Second * 30
		// duration = time.Second * 10
		interval = time.Millisecond * 250

		StoreConfig = `
		{
			"type": "dummy",
			"auth": {
				"secretRef": {
					"name": "credential-secret",
					"namespace": "default"
				}
			},
			"parameters": {
				"Suffix": "I am definitely a param"
			}
		}`
	)

	Context("When creating a SecretStore", func() {
		ctx := context.Background()
		It("Should intialize backend with the the given controller name", func() {

			credentialsSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      CredentialSecretName,
					Namespace: SecretStoreNamespace,
				},
				StringData: map[string]string{
					"operator-credentials.json": `{
						"Credential": "-dummyvalue"
					}`,
				},
			}
			Expect(k8sClient.Create(ctx, credentialsSecret)).Should(Succeed())

			credentialsSecretLookupKey := types.NamespacedName{Name: CredentialSecretName, Namespace: SecretStoreNamespace}
			createdCredentialsSecret := &corev1.Secret{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, credentialsSecretLookupKey, createdCredentialsSecret)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			secretStore := &storev1alpha1.SecretStore{}

			secretStore.ObjectMeta = metav1.ObjectMeta{
				Name:      SecretStoreName,
				Namespace: SecretStoreNamespace,
			}

			secretStore.Spec = storev1alpha1.SecretStoreSpec{
				Controller: SecretStoreControllerName,
				Store: runtime.RawExtension{
					Raw: []byte(StoreConfig),
				},
			}

			Expect(k8sClient.Create(ctx, secretStore)).Should(Succeed())

			secretStoreLookupKey := types.NamespacedName{Name: SecretStoreName, Namespace: SecretStoreNamespace}
			createdSecretStore := &storev1alpha1.SecretStore{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, secretStoreLookupKey, createdSecretStore)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Eventually(func() string {
				backend := backend.Instances[SecretStoreControllerName]
				if backend == nil {
					return ""
				}
				secretValue, err := backend.Get(KeyName, KeyVersion)
				if err != nil {
					return ""
				}
				return secretValue
			}, timeout, interval).Should(Equal("test-secrettest-versionI am definitely a param"))

		})
	})

	Context("When creating a SecretStore", func() {
		ctx := context.Background()

		It("Should handle a missing secret store gracefully", func() {
			randomSecretName := "Non existernt Secret Store" + utils.RandomString(20)

			secretStoreLookupKey := types.NamespacedName{Name: randomSecretName, Namespace: SecretStoreNamespace}
			nonExistentSecretStore := &storev1alpha1.SecretStore{}

			// secretStore := &storev1alpha1.SecretStore{}

			err := k8sClient.Get(ctx, secretStoreLookupKey, nonExistentSecretStore)

			Expect(err).ToNot(BeNil())
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})
	})

	Context("When creating a SecretStore", func() {
		ctx := context.Background()

		StoreConfig2 := `
		{
			"type": "dummy",
			"auth": {
				"secretRef": {
					"name": "credential-secret-non-existent",
					"namespace": "default"
				}
			},
			"parameters": {
				"Suffix": "I am definitely a param"
			}
		}`

		It("Should handle a missing credential secret", func() {
			randomSecretStoreName := SecretStoreName + utils.RandomString(20)

			secretStore := &storev1alpha1.SecretStore{}

			secretStore.ObjectMeta = metav1.ObjectMeta{
				Name:      randomSecretStoreName,
				Namespace: SecretStoreNamespace,
			}

			secretStore.Spec = storev1alpha1.SecretStoreSpec{
				Controller: SecretStoreControllerName,
				Store: runtime.RawExtension{
					Raw: []byte(StoreConfig2),
				},
			}

			Expect(k8sClient.Create(ctx, secretStore)).Should(Succeed())

			secretStoreLookupKey := types.NamespacedName{Name: randomSecretStoreName, Namespace: SecretStoreNamespace}
			createdSecretStore := &storev1alpha1.SecretStore{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, secretStoreLookupKey, createdSecretStore)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})
	})

	Context("When creating a SecretStore", func() {
		ctx := context.Background()
		// blank params trigger error during dummy Init()
		StoreConfig3 := `
		{
			"type": "dummy",
			"auth": {
				"secretRef": {
					"name": "credential-secret",
					"namespace": "default"
				}
			},
			"parameters": {}
		}`

		It("Should handle Init() failure", func() {

			randomSecretStoreName := SecretStoreName + utils.RandomString(20)

			secretStore := &storev1alpha1.SecretStore{}

			secretStore.ObjectMeta = metav1.ObjectMeta{
				Name:      randomSecretStoreName,
				Namespace: SecretStoreNamespace,
			}

			secretStore.Spec = storev1alpha1.SecretStoreSpec{
				Controller: SecretStoreControllerName,
				Store: runtime.RawExtension{
					Raw: []byte(StoreConfig3),
				},
			}

			Expect(k8sClient.Create(ctx, secretStore)).Should(Succeed())
		})
	})
})
