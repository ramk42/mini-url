package normalizer

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNormalizer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Normalizer Suite")
}

var _ = Describe("NormalizeURL", func() {
	DescribeTable("Normalization cases",
		func(input, expected string) {
			result, err := NormalizeURL(input)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(expected))
		},
		Entry("Path segment handling", "http://a.com/foo/./bar/../baz", "http://a.com/foo/baz"),
		Entry("Removal of duplicate slashes", "http://a.com//foo///bar", "http://a.com/foo/bar"),
		Entry("Complex path resolution", "http://a.com/../a/b/../c/./d.html", "http://a.com/a/c/d.html"),
		Entry("Removal of default port", "http://example.com:80/path", "http://example.com/path"),
		Entry("Uppercase scheme and host", "HTTP://EXAMPLE.COM", "http://example.com"),
		Entry("Query parameter sorting", "http://example.com?b=2&a=1", "http://example.com?a=1&b=2"),
	)
})
