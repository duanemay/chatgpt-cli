package cmd_test

import (
	"github.com/duanemay/chatgpt-cli/cmd"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sashabaranov/go-openai"
)

var _ = Describe("Image Flags", func() {
	It("should fail on New Flags", func() {
		imageFlags := cmd.NewImageFlags()
		Ω(imageFlags.ValidateFlags()).Error().To(HaveOccurred())
	})

	It("should validate Size", func() {
		imageFlags := cmd.NewImageFlags()
		imageFlags.NumberImages = 1
		err := imageFlags.ValidateFlags()
		Ω(err).Error().To(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("Size must be one of"))

		imageFlags.Size = openai.CreateImageSize256x256
		err = imageFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())

		imageFlags.Size = openai.CreateImageSize512x512
		err = imageFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())

		imageFlags.Size = openai.CreateImageSize1024x1024
		err = imageFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
	})

	It("should validate Number of Images", func() {
		imageFlags := cmd.NewImageFlags()
		imageFlags.Size = openai.CreateImageSize1024x1024

		err := imageFlags.ValidateFlags()
		Ω(err).Error().To(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("NumberImages must be between 1 and 10"))

		imageFlags.NumberImages = 11
		err = imageFlags.ValidateFlags()
		Ω(err).Error().To(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("NumberImages must be between 1 and 10"))

		imageFlags.NumberImages = 1
		err = imageFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())

		imageFlags.NumberImages = 10
		err = imageFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
	})
})
