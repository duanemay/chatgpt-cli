package cmd_test

import (
	"github.com/duanemay/chatgpt-cli/cmd"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sashabaranov/go-openai"
)

var _ = Describe("Speech Flags", func() {
	It("should fail on New Flags", func() {
		speechFlags := cmd.NewSpeechFlags()
		Ω(speechFlags.ValidateFlags()).Error().To(HaveOccurred())
	})

	It("should validate Speed", func() {
		speechFlags := cmd.NewSpeechFlags()
		speechFlags.Speed = 0.1
		err := speechFlags.ValidateFlags()
		Ω(err).Error().To(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("speed must be"))

		speechFlags.Speed = 5.0
		err = speechFlags.ValidateFlags()
		Ω(err).Error().To(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("speed must be"))

		speechFlags.Speed = 1.0
		speechFlags.ModelStr = string(openai.TTSModel1)
		speechFlags.VoiceStr = string(openai.VoiceAlloy)
		err = speechFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
	})

	It("should validate Model", func() {
		speechFlags := cmd.NewSpeechFlags()
		speechFlags.Speed = 1.0
		speechFlags.VoiceStr = string(openai.VoiceAlloy)
		speechFlags.Model = ""

		speechFlags.ModelStr = "no-good"
		err := speechFlags.ValidateFlags()
		Ω(err).Error().To(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("model must be"))
		Ω(speechFlags.Model).To(Equal(openai.SpeechModel("")))

		speechFlags.ModelStr = string(openai.TTSModel1)
		err = speechFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
		Ω(speechFlags.Model).To(Equal(openai.TTSModel1))

		speechFlags.ModelStr = string(openai.TTSModel1HD)
		err = speechFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
		Ω(speechFlags.Model).To(Equal(openai.TTSModel1HD))

		speechFlags.ModelStr = string(openai.TTSModelCanary)
		err = speechFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
		Ω(speechFlags.Model).To(Equal(openai.TTSModelCanary))
	})

	It("should validate Voice", func() {
		speechFlags := cmd.NewSpeechFlags()
		speechFlags.Speed = 1.0
		speechFlags.ModelStr = string(openai.TTSModelCanary)

		speechFlags.VoiceStr = "no-good"
		err := speechFlags.ValidateFlags()
		Ω(err).Error().To(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("voice must be"))
		Ω(speechFlags.Voice).To(Equal(openai.SpeechVoice("")))

		speechFlags.VoiceStr = string(openai.VoiceAlloy)
		err = speechFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
		Ω(speechFlags.Voice).To(Equal(openai.VoiceAlloy))

		speechFlags.VoiceStr = string(openai.VoiceEcho)
		err = speechFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
		Ω(speechFlags.Voice).To(Equal(openai.VoiceEcho))

		speechFlags.VoiceStr = string(openai.VoiceFable)
		err = speechFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
		Ω(speechFlags.Voice).To(Equal(openai.VoiceFable))

		speechFlags.VoiceStr = string(openai.VoiceOnyx)
		err = speechFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
		Ω(speechFlags.Voice).To(Equal(openai.VoiceOnyx))

		speechFlags.VoiceStr = string(openai.VoiceNova)
		err = speechFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
		Ω(speechFlags.Voice).To(Equal(openai.VoiceNova))

		speechFlags.VoiceStr = string(openai.VoiceShimmer)
		err = speechFlags.ValidateFlags()
		Ω(err).Error().ToNot(HaveOccurred())
		Ω(speechFlags.Voice).To(Equal(openai.VoiceShimmer))
	})
})
