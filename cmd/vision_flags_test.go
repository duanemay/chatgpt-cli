package cmd_test

import (
	"github.com/duanemay/chatgpt-cli/cmd"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sashabaranov/go-openai"
)

var _ = Describe("Vision Flags", func() {
	It("should fail on New Flags", func() {
		visionFlags := cmd.NewVisionFlags()
		Ω(visionFlags.ValidateFlags()).Error().To(HaveOccurred())
	})

	It("should validate Detail", func() {
		visionFlags := cmd.NewVisionFlags()
		visionFlags.DetailStr = "no-good"
		visionFlags.Detail = ""
		err := visionFlags.ValidateFlags()
		Ω(err).Error().To(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("detail must be"))
		Ω(visionFlags.Detail).To(Equal(openai.ImageURLDetail("")))

		visionFlags.DetailStr = string(openai.ImageURLDetailAuto)
		err = visionFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
		Ω(visionFlags.Detail).To(Equal(openai.ImageURLDetailAuto))

		visionFlags.DetailStr = string(openai.ImageURLDetailAuto)
		err = visionFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
		Ω(visionFlags.Detail).To(Equal(openai.ImageURLDetailAuto))

		visionFlags.DetailStr = string(openai.ImageURLDetailAuto)
		err = visionFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
		Ω(visionFlags.Detail).To(Equal(openai.ImageURLDetailAuto))
	})
})
