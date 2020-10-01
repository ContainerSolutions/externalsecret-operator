package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	utils "github.com/containersolutions/externalsecret-operator/pkg/utils"
)

var _ = Describe("Utils", func() {

	BeforeEach(func() {})

	AfterEach(func() {})
	Context("Should generate random string", func() {
		It("Should succeed", func() {
			Expect(len(utils.RandomString(40))).Should(BeEquivalentTo(40))
		})

		It("Should not be nil", func() {
			actual := utils.RandomString(40)
			Expect(actual).ToNot(BeNil())
		})

		// Expect(actual).To(B)
	})

})
