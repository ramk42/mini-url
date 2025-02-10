package usecase

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ramk42/mini-url/internal/apperr"
	"github.com/ramk42/mini-url/internal/model"
	"github.com/ramk42/mini-url/internal/port/mock_port"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Url use case", func() {
	const (
		validURL = "https://example.com/path?query=param"
	)
	var (
		ctx         context.Context
		mockCtrl    *gomock.Controller
		urlRepoMock *mock_port.MockURLRepository
		uc          *URLShortener
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockCtrl = gomock.NewController(GinkgoT())
		urlRepoMock = mock_port.NewMockURLRepository(mockCtrl)
		uc = NewURLShortener(urlRepoMock, "example.com")
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("ShortenURL", func() {
		Context("Nominal case", func() {
			It("should generate valid short URL", func() {
				urlRepoMock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(model.URL{
					Slug: "slug",
				}, nil)

				// Execution
				result, err := uc.ShortenURL(ctx, model.URL{Original: validURL}, 0)

				// Assertions
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal("example.com/slug"))
			})
		})
		Context("Slug collision handling", func() {
			It("should retry after collision", func() {
				// Mock configuration with 2 collisions
				urlRepoMock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(model.URL{}, apperr.ErrURLConflict).Times(2)
				urlRepoMock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(model.URL{Slug: "hummus"}, nil)

				_, err := uc.ShortenURL(ctx, model.URL{Original: validURL}, 0)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should fail after 3 attempts", func() {
				urlRepoMock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(model.URL{Slug: "pizza"}, apperr.ErrURLConflict).Times(3)
				_, err := uc.ShortenURL(ctx, model.URL{Original: validURL}, 0)
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Context("Resolve", func() {
		It("should return the original URL on success", func() {
			urlRepoMock.EXPECT().UpdateClicks(ctx, gomock.Any()).Return(model.URL{}, nil)

			_, err := uc.Resolve(ctx, "hummus")
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return an error if UpdateClicks fails", func() {
			urlRepoMock.EXPECT().UpdateClicks(ctx, gomock.Any()).Return(model.URL{}, errors.New("not found"))

			result, err := uc.Resolve(ctx, "hummus")
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeEmpty())
		})
	})
})
