package cmd

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Common", func() {
	Describe("setupOpenAIClient", func() {
		Context("when the API key is not set", func() {
			It("should exit with an error", func() {
				_, err := setupOpenAIClient("")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
