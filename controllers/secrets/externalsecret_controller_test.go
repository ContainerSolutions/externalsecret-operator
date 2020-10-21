package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	secretsv1alpha1 "github.com/containersolutions/externalsecret-operator/apis/secrets/v1alpha1"
	storev1alpha1 "github.com/containersolutions/externalsecret-operator/apis/store/v1alpha1"
	"github.com/containersolutions/externalsecret-operator/pkg/utils"
)

const ExternalSecretNamespace = "default"

var _ = Describe("ExternalsecretController", func() {
	var (
		ExternalSecretName     = "externalsecret-operator-test"
		ExternalSecretKey      = "test-key"
		ExternalSecretVersion  = "test-version"
		ExternalSecret2Key     = "test-key-2"
		ExternalSecret2Version = "test-version-2"
		ExternalSecret3Key     = "test-key-3"
		ExternalSecret3Version = "test-version-3"
		SecretStoreName        = "test-externalsecret-store"
		StoreControllerName    = "test-externalsecret-ctrl"
		CredentialSecretName   = "credential-secret-external-secret"

		timeout = time.Second * 30
		// duration = time.Second * 10
		interval = time.Millisecond * 250

		StoreConfig = `
		{
			"type": "dummy",
			"auth": {
				"secretRef": {
					"name": "credential-secret-external-secret",
					"namespace": "default"
				}
			},
			"parameters": {
				"Suffix": "TestParameter"
			}
		}`
	)

	BeforeEach(func() {})

	AfterEach(func() {})

	Context("When creating ExternalSecret", func() {
		It("Should handle ExternalSecret correctly", func() {
			By("Creating a new ExternalSecret")
			ctx := context.Background()

			credentialsSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      CredentialSecretName,
					Namespace: ExternalSecretNamespace,
				},
				StringData: map[string]string{
					"operator-config.json": `{
						"Credential": "-dummyvalue"
					}`,
				},
			}

			Expect(k8sClient.Create(ctx, credentialsSecret)).Should(Succeed())

			credentialsSecretLookupKey := types.NamespacedName{Name: CredentialSecretName, Namespace: ExternalSecretNamespace}
			createdCredentialsSecret := &corev1.Secret{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, credentialsSecretLookupKey, createdCredentialsSecret)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			secretStore := &storev1alpha1.SecretStore{
				ObjectMeta: metav1.ObjectMeta{
					Name:      SecretStoreName,
					Namespace: ExternalSecretNamespace,
				},

				Spec: storev1alpha1.SecretStoreSpec{
					Controller: StoreControllerName,
					Store: runtime.RawExtension{
						Raw: []byte(StoreConfig),
					},
				},
			}

			Expect(k8sClient.Create(ctx, secretStore)).Should(Succeed())

			secretStoreLookupKey := types.NamespacedName{Name: SecretStoreName, Namespace: ExternalSecretNamespace}
			createdSecretStore := &storev1alpha1.SecretStore{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, secretStoreLookupKey, createdSecretStore)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			externalSecret := &secretsv1alpha1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ExternalSecretName,
					Namespace: ExternalSecretNamespace,
				},
				Spec: secretsv1alpha1.ExternalSecretSpec{
					StoreRef: secretsv1alpha1.StoreRef{
						Name:      SecretStoreName,
						Namespace: ExternalSecretNamespace,
					},
					Secrets: []secretsv1alpha1.Secret{
						{
							Key:     ExternalSecretKey,
							Version: ExternalSecretVersion,
						},

						{
							Key:     ExternalSecret2Key,
							Version: ExternalSecret2Version,
						},

						{
							Key:     ExternalSecret3Key,
							Version: ExternalSecret3Version,
						},
					},
				},
			}

			Expect(k8sClient.Create(ctx, externalSecret)).Should(Succeed())

			externalSecretLookupKey := types.NamespacedName{Name: ExternalSecretName, Namespace: ExternalSecretNamespace}
			createdExternalSecret := &secretsv1alpha1.ExternalSecret{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, externalSecretLookupKey, createdExternalSecret)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(len(createdExternalSecret.Spec.Secrets)).Should(BeNumerically("==", 3))

			Expect(createdExternalSecret.Spec.Secrets[0].Key).Should(Equal(ExternalSecretKey))
			Expect(createdExternalSecret.Spec.Secrets[0].Version).Should(Equal(ExternalSecretVersion))

			Expect(createdExternalSecret.Spec.Secrets[1].Key).Should(Equal(ExternalSecret2Key))
			Expect(createdExternalSecret.Spec.Secrets[1].Version).Should(Equal(ExternalSecret2Version))

			Expect(createdExternalSecret.Spec.Secrets[2].Key).Should(Equal(ExternalSecret3Key))
			Expect(createdExternalSecret.Spec.Secrets[2].Version).Should(Equal(ExternalSecret3Version))

			By("Creating a new secret with correct values")
			secret := &corev1.Secret{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, externalSecretLookupKey, secret)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			secretValue := string(secret.Data[ExternalSecretKey])
			Expect(string(secretValue)).Should(Equal("test-keytest-versionTestParameter"))

			secretValue2 := string(secret.Data[ExternalSecret2Key])
			Expect(string(secretValue2)).Should(Equal("test-key-2test-version-2TestParameter"))

			secretValue3 := string(secret.Data[ExternalSecret3Key])
			Expect(string(secretValue3)).Should(Equal("test-key-3test-version-3TestParameter"))

			By("Deleting the External Secret")
			Eventually(func() error {
				es := &secretsv1alpha1.ExternalSecret{}
				k8sClient.Get(context.Background(), externalSecretLookupKey, es)
				return k8sClient.Delete(context.Background(), es)
			}, timeout, interval).Should(Succeed())

			Eventually(func() error {
				es := &secretsv1alpha1.ExternalSecret{}
				return k8sClient.Get(context.Background(), externalSecretLookupKey, es)
			}, timeout, interval).ShouldNot(Succeed())
		})
	})

	Context("SecretStore does not exist", func() {
		ctx := context.Background()
		It("Should handle gracefully", func() {
			randomObjSafeStr, err := utils.RandomStringObjectSafe(25)
			Expect(err).To(BeNil())
			randomSecretStoreName := SecretStoreName + randomObjSafeStr
			externalSecret := &secretsv1alpha1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      randomSecretStoreName,
					Namespace: ExternalSecretNamespace,
				},
				Spec: secretsv1alpha1.ExternalSecretSpec{
					StoreRef: secretsv1alpha1.StoreRef{
						Name:      "NonExistentStore",
						Namespace: ExternalSecretNamespace,
					},
					Secrets: []secretsv1alpha1.Secret{
						{
							Key:     ExternalSecretKey,
							Version: ExternalSecretVersion,
						},
						{
							Key:     ExternalSecret2Key,
							Version: ExternalSecret3Version,
						},
					},
				},
			}

			Expect(k8sClient.Create(ctx, externalSecret)).Should(Succeed())

		})
	})

	Context("When interacting with Backend", func() {
		r := &ExternalSecretReconciler{}
		ctx := context.Background()

		It("Should Fail when a backend is uninitialized/Not ready", func() {
			randomObjSafeStr, err := utils.RandomStringObjectSafe(30)
			Expect(err).To(BeNil())
			randomSecretStoreName := SecretStoreName + randomObjSafeStr
			randomControllerName := StoreControllerName + randomObjSafeStr
			secretStore := &storev1alpha1.SecretStore{
				ObjectMeta: metav1.ObjectMeta{
					Name:      randomSecretStoreName,
					Namespace: ExternalSecretNamespace,
				},
				Spec: storev1alpha1.SecretStoreSpec{
					Controller: randomControllerName,
					Store: runtime.RawExtension{
						Raw: []byte(StoreConfig),
					},
				},
			}

			externalSecret := &secretsv1alpha1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ExternalSecretName,
					Namespace: ExternalSecretNamespace,
				},
				Spec: secretsv1alpha1.ExternalSecretSpec{
					StoreRef: secretsv1alpha1.StoreRef{
						Name:      secretStore.ObjectMeta.Name,
						Namespace: ExternalSecretNamespace,
					},
					Secrets: []secretsv1alpha1.Secret{
						{
							Key:     ExternalSecretKey,
							Version: ExternalSecretVersion,
						},
						{
							Key:     ExternalSecret2Key,
							Version: ExternalSecret3Version,
						},
					},
				},
			}
			_, err = r.backendGet(externalSecret, secretStore)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).Should(Equal("Cannot find backend:" + " " + randomControllerName))

		})

		It("Should return an error when Get() fails in the backend", func() {
			randomObjSafeStr, err := utils.RandomStringObjectSafe(30)
			Expect(err).To(BeNil())
			randomSecretStoreName := SecretStoreName + randomObjSafeStr
			randomControllerName := StoreControllerName + randomObjSafeStr
			secretStore := &storev1alpha1.SecretStore{
				ObjectMeta: metav1.ObjectMeta{
					Name:      randomSecretStoreName,
					Namespace: ExternalSecretNamespace,
				},
				Spec: storev1alpha1.SecretStoreSpec{
					Controller: randomControllerName,
					Store: runtime.RawExtension{
						Raw: []byte(StoreConfig),
					},
				},
			}

			/**
				Create the store so that the backend is intialized
			**/
			Expect(k8sClient.Create(ctx, secretStore)).Should(Succeed())

			secretStoreLookupKey := types.NamespacedName{Name: randomSecretStoreName, Namespace: ExternalSecretNamespace}
			createdSecretStore := &storev1alpha1.SecretStore{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, secretStoreLookupKey, createdSecretStore)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			externalSecret := &secretsv1alpha1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ExternalSecretName,
					Namespace: ExternalSecretNamespace,
				},
				Spec: secretsv1alpha1.ExternalSecretSpec{
					StoreRef: secretsv1alpha1.StoreRef{
						Name:      randomSecretStoreName,
						Namespace: ExternalSecretNamespace,
					},
					Secrets: []secretsv1alpha1.Secret{
						{
							Key:     "",
							Version: "",
						},
					},
				},
			}

			/**
				We need to wait for the store reconciler to intialize the backend
			**/
			Eventually(func() string {
				_, err := r.backendGet(externalSecret, secretStore)
				return err.Error()
			}, timeout, interval).Should(Equal("could not create secret due to error from backend: empty key provided"))

		})

	})

})
