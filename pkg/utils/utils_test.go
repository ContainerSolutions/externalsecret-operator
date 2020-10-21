package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	utils "github.com/containersolutions/externalsecret-operator/pkg/utils"
)

var _ = Describe("Utils", func() {

	BeforeEach(func() {})

	AfterEach(func() {})
	Context("Should generate random base64 encoded string", func() {
		It("Should succeed", func() {
			_, err := utils.RandomStringObjectSafe(40)
			Expect(err).To(BeNil())
		})

		It("Should not be nil", func() {
			str, err := utils.RandomStringObjectSafe(40)
			Expect(err).To(BeNil())
			Expect(str).ToNot(BeNil())
		})
	})

	Context("Should generate random bytes", func() {
		It("Should succeed", func() {
			_, err := utils.RandomBytes(40)
			Expect(err).To(BeNil())
		})

		It("Should not be nil", func() {
			str, err := utils.RandomBytes(40)
			Expect(err).To(BeNil())
			Expect(str).ToNot(BeNil())
		})
	})

	Context("Should generate int64", func() {
		It("Should succeed", func() {
			_, err := utils.RandomInt()
			Expect(err).To(BeNil())
		})

		It("Should not be nil", func() {
			ranInt64, err := utils.RandomInt()
			Expect(err).To(BeNil())
			Expect(ranInt64).ToNot(BeNil())
		})
	})

})
