package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	secretsv1alpha1 "github.com/containersolutions/externalsecret-operator/apis/secrets/v1alpha1"
)

const ExternalSecretNamespace = "default"

var _ = Describe("ExternalsecretController", func() {
	var (
		ExternalSecretName    = "externalsecret-operator-test"
		ExternalSecretKey     = "test-key"
		ExternalSecretVersion = "test-version"
		ExternalSecretBackend = "test-backend"
		// SecretName            = "test-secret"

		timeout = time.Second * 30
		// duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	BeforeEach(func() {})

	AfterEach(func() {})

	Context("When creating ExternalSecret", func() {
		It("Should handle ExternalSecret correctly", func() {
			By("Creating a new ExternalSecret")
			ctx := context.Background()
			externalSecret := &secretsv1alpha1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ExternalSecretName,
					Namespace: ExternalSecretNamespace,
				},
				Spec: secretsv1alpha1.ExternalSecretSpec{
					Backend: ExternalSecretBackend,
					Key:     ExternalSecretKey,
					Version: ExternalSecretVersion,
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

			Expect(createdExternalSecret.Spec.Backend).Should(Equal("test-backend"))
			Expect(createdExternalSecret.Spec.Version).Should(Equal("test-version"))
			Expect(createdExternalSecret.Spec.Key).Should(Equal("test-key"))

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

			Expect(string(secretValue)).Should(Equal("test-keytest-version-ohlord"))

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

	Context("Errors", func() {
		r := &ExternalSecretReconciler{}

		It("Should Fail when nil externalsecret is passed", func() {
			_, err := r.newSecretForCR(nil)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).Should(Equal("externalsecret is nil"))

		})

		It("Should Fail when nil externalsecret is passed", func() {
			_, err := r.backendGet(nil)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).Should(Equal("externalsecret is nil"))

		})

		It("Should Fail when nil an invalid backend is passed", func() {

			externalSecret := &secretsv1alpha1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ExternalSecretName,
					Namespace: ExternalSecretNamespace,
				},
				Spec: secretsv1alpha1.ExternalSecretSpec{
					Backend: "invalid backend",
					Key:     ExternalSecretKey,
					Version: ExternalSecretVersion,
				},
			}
			_, err := r.backendGet(externalSecret)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).Should(Equal("Cannot find backend: invalid backend"))

		})

		It("Should fail when backend with empty key is passed", func() {
			externalSecret := &secretsv1alpha1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ExternalSecretName,
					Namespace: ExternalSecretNamespace,
				},
				Spec: secretsv1alpha1.ExternalSecretSpec{
					Backend: ExternalSecretBackend,
					Key:     "",
					Version: ExternalSecretVersion,
				},
			}

			_, err := r.newSecretForCR(externalSecret)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).Should(Equal("could not create secret due to error from backend: empty key provided"))

		})

		It("Should return an error when Get() fails in the backend", func() {

			externalSecret := &secretsv1alpha1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ExternalSecretName,
					Namespace: ExternalSecretNamespace,
				},
				Spec: secretsv1alpha1.ExternalSecretSpec{
					Backend: ExternalSecretBackend,
					Key:     "",
					Version: "",
				},
			}
			_, err := r.backendGet(externalSecret)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).Should(Equal("could not create secret due to error from backend: empty key provided"))

		})

	})

})
