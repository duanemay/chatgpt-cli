package cmd_test

import (
	"github.com/duanemay/chatgpt-cli/cmd"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sashabaranov/go-openai"
)

var _ = Describe("Embedding Flags", func() {
	It("should fail on New Flags", func() {
		embeddingFlags := cmd.NewEmbeddingFlags()
		Ω(embeddingFlags.ValidateFlags()).Error().To(HaveOccurred())
	})

	It("should validate Model", func() {
		embeddingFlags := cmd.NewEmbeddingFlags()

		embeddingFlags.ModelStr = "no-good"
		err := embeddingFlags.ValidateFlags()
		Ω(err).Error().To(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("model must be"))
		Ω(embeddingFlags.Model).To(Equal(openai.EmbeddingModel("")))

		embeddingFlags.ModelStr = string(openai.SmallEmbedding3)
		err = embeddingFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
		Ω(embeddingFlags.Model).To(Equal(openai.SmallEmbedding3))

		embeddingFlags.ModelStr = string(openai.LargeEmbedding3)
		err = embeddingFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
		Ω(embeddingFlags.Model).To(Equal(openai.LargeEmbedding3))

		embeddingFlags.ModelStr = string(openai.AdaEmbeddingV2)
		err = embeddingFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
		Ω(embeddingFlags.Model).To(Equal(openai.AdaEmbeddingV2))
	})

	It("should validate Dimensions", func() {
		embeddingFlags := cmd.NewEmbeddingFlags()
		embeddingFlags.ModelStr = string(openai.SmallEmbedding3)

		embeddingFlags.Dimensions = -1
		err := embeddingFlags.ValidateFlags()
		Ω(err).Error().To(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("dimensions must be"))

		embeddingFlags.Dimensions = 0
		err = embeddingFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())

		embeddingFlags.Dimensions = 256
		err = embeddingFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
	})
})
