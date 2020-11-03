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
		ExternalSecretName          = "externalsecret-operator-test"
		ExternalSecretKey           = "test-key"
		ExternalSecretVersion       = "test-version"
		ExternalSecret2Key          = "test-key-2"
		ExternalSecret2Version      = "test-version-2"
		ExternalSecret3Key          = "test-key-3"
		ExternalSecret3Version      = "test-version-3"
		ExternalSecretKeyError      = "ErroredKey"
		ExternalSecretKeyUpdate     = "test-key-update"
		ExternalSecretVersionUpdate = "test-version-update"
		SecretStoreName             = "test-externalsecret-store"
		StoreControllerName         = "test-externalsecret-ctrl"
		CredentialSecretName        = "credential-secret-external-secret"
		TargetName                  = "test-secret-target"

		timeout  = time.Second * 30
		duration = time.Second * 5
		interval = time.Millisecond * 250

		StoreConfig = `
		{
			"type": "dummy",
			"auth": {
				"secretRef": {
					"name": "credential-secret-external-secret"
				}
			},
			"parameters": {
				"Suffix": "TestParameter"
			}
		}`
	)

	BeforeEach(func() {})

	AfterEach(func() {})

	Context("Given an ExternalSecret", func() {
		It("Should handle ExternalSecret correctly", func() {
			By("Creating a new ExternalSecret")
			ctx := context.Background()

			credentialsSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      CredentialSecretName,
					Namespace: ExternalSecretNamespace,
				},
				StringData: map[string]string{
					"credentials.json": `{
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
					StoreRef: secretsv1alpha1.ExternalSecretStoreRef{
						Name: SecretStoreName,
					},
					Target: secretsv1alpha1.ExternalSecretTarget{
						Template: runtime.RawExtension{
							Raw: []byte(`{}`),
						},
					},
					Data: []secretsv1alpha1.ExternalSecretData{
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

			Expect(len(createdExternalSecret.Spec.Data)).Should(BeNumerically("==", 3))

			Expect(createdExternalSecret.Spec.Data[0].Key).Should(Equal(ExternalSecretKey))
			Expect(createdExternalSecret.Spec.Data[0].Version).Should(Equal(ExternalSecretVersion))

			Expect(createdExternalSecret.Spec.Data[1].Key).Should(Equal(ExternalSecret2Key))
			Expect(createdExternalSecret.Spec.Data[1].Version).Should(Equal(ExternalSecret2Version))

			Expect(createdExternalSecret.Spec.Data[2].Key).Should(Equal(ExternalSecret3Key))
			Expect(createdExternalSecret.Spec.Data[2].Version).Should(Equal(ExternalSecret3Version))

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
			Expect(secretValue).Should(Equal("test-keytest-versionTestParameter"))

			secretValue2 := string(secret.Data[ExternalSecret2Key])
			Expect(secretValue2).Should(Equal("test-key-2test-version-2TestParameter"))

			secretValue3 := string(secret.Data[ExternalSecret3Key])
			Expect(secretValue3).Should(Equal("test-key-3test-version-3TestParameter"))

			By("Updating the Secret if it already exists")
			updatedSecrets := []secretsv1alpha1.ExternalSecretData{
				{
					Key:     ExternalSecretKeyUpdate,
					Version: ExternalSecretVersionUpdate,
				},
			}

			createdExternalSecret.Spec.Data = updatedSecrets

			Expect(k8sClient.Update(ctx, createdExternalSecret)).Should(Succeed())

			updatedExternalSecret := &secretsv1alpha1.ExternalSecret{}
			Eventually(func() []secretsv1alpha1.ExternalSecretData {
				err := k8sClient.Get(ctx, externalSecretLookupKey, updatedExternalSecret)
				if err != nil {
					return []secretsv1alpha1.ExternalSecretData{}
				}
				return updatedExternalSecret.Spec.Data
			}, timeout, interval).Should(Equal(updatedSecrets))

			updatedSecret := &corev1.Secret{}
			Eventually(func() string {
				err := k8sClient.Get(ctx, externalSecretLookupKey, updatedSecret)
				if err != nil {
					return ""
				}
				return string(updatedSecret.Data[ExternalSecretKeyUpdate])
			}, timeout, interval).Should(Equal("test-key-updatetest-version-updateTestParameter"))

			By("Deleting the External Secret")
			Eventually(func() error {
				es := &secretsv1alpha1.ExternalSecret{}
				k8sClient.Get(context.Background(), externalSecretLookupKey, es)
				return k8sClient.Delete(ctx, es)
			}, timeout, interval).Should(Succeed())

			Eventually(func() error {
				es := &secretsv1alpha1.ExternalSecret{}
				return k8sClient.Get(ctx, externalSecretLookupKey, es)
			}, timeout, interval).ShouldNot(Succeed())
		})

	})

	Context("When target Name is provided", func() {
		It("Should use it as name for the secret resource", func() {
			ctx := context.Background()

			randomObjSafeStr, err := utils.RandomStringObjectSafe(32)
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

			Expect(k8sClient.Create(ctx, secretStore)).Should(Succeed())

			randomObjSafeStr, err = utils.RandomStringObjectSafe(21)
			Expect(err).To(BeNil())

			randomExternalSecretName := ExternalSecretName + randomObjSafeStr
			externalSecret := &secretsv1alpha1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      randomExternalSecretName,
					Namespace: ExternalSecretNamespace,
				},
				Spec: secretsv1alpha1.ExternalSecretSpec{
					StoreRef: secretsv1alpha1.ExternalSecretStoreRef{
						Name: secretStore.ObjectMeta.Name,
					},
					Target: secretsv1alpha1.ExternalSecretTarget{
						Name: TargetName,
						Template: runtime.RawExtension{
							Raw: []byte(`{}`),
						},
					},
					Data: []secretsv1alpha1.ExternalSecretData{
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

			secretLookupKey := types.NamespacedName{Name: TargetName, Namespace: ExternalSecretNamespace}
			secret := &corev1.Secret{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, secretLookupKey, secret)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(secret.ObjectMeta.Name).Should(Equal(TargetName))
		})
	})

	Context("SecretStore does not exist", func() {
		ctx := context.Background()
		It("Should return an error", func() {
			randomObjSafeStr, err := utils.RandomStringObjectSafe(25)
			Expect(err).To(BeNil())

			randomExternalSecretName := ExternalSecretName + randomObjSafeStr
			externalSecret := &secretsv1alpha1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      randomExternalSecretName,
					Namespace: ExternalSecretNamespace,
				},
				Spec: secretsv1alpha1.ExternalSecretSpec{
					StoreRef: secretsv1alpha1.ExternalSecretStoreRef{
						Name: "NonExistentStore",
					},
					Target: secretsv1alpha1.ExternalSecretTarget{
						Template: runtime.RawExtension{
							Raw: []byte(`{}`),
						},
					},
					Data: []secretsv1alpha1.ExternalSecretData{
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

			secretLookupKey := types.NamespacedName{Name: randomExternalSecretName, Namespace: ExternalSecretNamespace}
			secret := &corev1.Secret{}
			Consistently(func() error {
				return k8sClient.Get(ctx, secretLookupKey, secret)
			}).ShouldNot(Succeed())

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

			randomObjSafeStr, err = utils.RandomStringObjectSafe(21)
			Expect(err).To(BeNil())

			randomExternalSecretName := ExternalSecretName + randomObjSafeStr
			externalSecret := &secretsv1alpha1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      randomExternalSecretName,
					Namespace: ExternalSecretNamespace,
				},
				Spec: secretsv1alpha1.ExternalSecretSpec{
					StoreRef: secretsv1alpha1.ExternalSecretStoreRef{
						Name: secretStore.ObjectMeta.Name,
					},
					Target: secretsv1alpha1.ExternalSecretTarget{
						Template: runtime.RawExtension{
							Raw: []byte(`{}`),
						},
					},
					Data: []secretsv1alpha1.ExternalSecretData{
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

			randomObjSafeStr, err = utils.RandomStringObjectSafe(37)
			Expect(err).To(BeNil())

			randomExternalSecretName := ExternalSecretName + randomObjSafeStr
			externalSecret := &secretsv1alpha1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      randomExternalSecretName,
					Namespace: ExternalSecretNamespace,
				},
				Spec: secretsv1alpha1.ExternalSecretSpec{
					StoreRef: secretsv1alpha1.ExternalSecretStoreRef{
						Name: randomSecretStoreName,
					},
					Target: secretsv1alpha1.ExternalSecretTarget{
						Template: runtime.RawExtension{
							Raw: []byte(`{}`),
						},
					},
					Data: []secretsv1alpha1.ExternalSecretData{
						{
							Key:     ExternalSecretKeyError,
							Version: "",
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, externalSecret)).Should(Succeed())

			secretLookupKey := types.NamespacedName{Name: randomExternalSecretName, Namespace: ExternalSecretNamespace}
			secret := &corev1.Secret{}
			Consistently(func() error {
				return k8sClient.Get(ctx, secretLookupKey, secret)
			}, duration, interval).ShouldNot(Succeed())

			/**
				We need to wait for the store reconciler to intialize the backend
			**/
			Eventually(func() string {
				_, err := r.backendGet(externalSecret, secretStore)
				return err.Error()
			}, timeout, interval).Should(Equal("could not create secret due to error from backend: Mocked error"))

		})

	})

	Context("Given a refreshInterval", func() {
		r := &ExternalSecretReconciler{}
		It("Should fail if the refreshInterval is invalid", func() {
			ctx := context.Background()
			randomObjSafeStr, err := utils.RandomStringObjectSafe(12)
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

			randomObjSafeStr, err = utils.RandomStringObjectSafe(38)
			Expect(err).To(BeNil())

			randomExternalSecretName := ExternalSecretName + randomObjSafeStr
			externalSecret := &secretsv1alpha1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      randomExternalSecretName,
					Namespace: ExternalSecretNamespace,
				},
				Spec: secretsv1alpha1.ExternalSecretSpec{
					RefreshInterval: "h",
					StoreRef: secretsv1alpha1.ExternalSecretStoreRef{
						Name: secretStore.ObjectMeta.Name,
					},
					Target: secretsv1alpha1.ExternalSecretTarget{
						Template: runtime.RawExtension{
							Raw: []byte(`{}`),
						},
					},
					Data: []secretsv1alpha1.ExternalSecretData{
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

			secretLookupKey := types.NamespacedName{Name: randomExternalSecretName, Namespace: ExternalSecretNamespace}
			secret := &corev1.Secret{}
			Consistently(func() error {
				return k8sClient.Get(ctx, secretLookupKey, secret)
			}, duration, interval).ShouldNot(Succeed())

		})

		It("Should parse it correctly", func() {
			refreshInterval, err := r.parseRefreshInterval("4h")
			Expect(err).To(BeNil())
			Expect(refreshInterval.Hours()).To(Equal(4.0))
		})

		It("Should return an error when passed refreshInterval is invalid", func() {
			_, err := r.parseRefreshInterval("h")
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("time: invalid duration \"h\""))
		})

		It("Should return an error if refreshInterval is has unit missing", func() {
			_, err := r.parseRefreshInterval("1")
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("time: missing unit in duration \"1\""))
		})

		It("Should return return a default interval when refreshInvterval is empty", func() {
			refreshInterval, err := r.parseRefreshInterval("")
			Expect(err).To(BeNil())
			Expect(refreshInterval.Hours()).To(Equal(1.0))
		})
	})

})
