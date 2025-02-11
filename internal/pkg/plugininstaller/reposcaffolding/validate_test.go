package reposcaffolding_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/reposcaffolding"
)

var _ = Describe("Validate func", func() {
	var (
		rawOpts configmanager.RawOptions
	)
	When("reposcaffolding option is not valid", func() {
		BeforeEach(func() {
			rawOpts = configmanager.RawOptions{
				"not_exist": true,
			}
		})
		It("should return err", func() {
			_, err := reposcaffolding.Validate(rawOpts)
			Expect(err).Should(HaveOccurred())
		})
	})
	When("reposcaffolding option is valid", func() {
		BeforeEach(func() {
			rawOpts = configmanager.RawOptions{
				"sourceRepo": map[string]string{
					"owner":    "test_user",
					"repo":     "test_repo",
					"repoType": "github",
				},
				"destinationRepo": map[string]string{
					"owner":    "dst_user",
					"repo":     "dst_repo",
					"repoType": "github",
				},
			}
		})
		It("should return noraml", func() {
			_, err := reposcaffolding.Validate(rawOpts)
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
	})
})
