package service

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"myMod/entity"
	"myMod/repository"
)

const (
	TITLE       = "Video Title"
	DESCRIPTION = "Video Description"
	URL         = "https://youtu.be/JgW-i2QjgHQ"
	FIRSTNAME   = "John"
	LASTNAME    = "Doe"
	EMAIL       = "jdoe@mail.com"
)

var testVideo entity.Video = entity.Video{
	Title:       TITLE,
	Description: DESCRIPTION,
	URL:         URL,
	Author: entity.Person{
		FirstName: FIRSTNAME,
		LastName:  LASTNAME,
		Email:     EMAIL,
	},
}

var _ = ginkgo.Describe("Video Service", func() {

	var (
		videoRepository repository.VideoRepository
		videoService    VideoService
		//videoController controller.VideoController
	)

	ginkgo.BeforeSuite(func() {
		videoRepository = repository.NewVideoRepository()
		videoService = NewVideoService(videoRepository)
		//videoController = controller.NewVideoController(videoService)

	})

	ginkgo.Describe("Fetching all existing videos", func() {

		ginkgo.Context("If there is a video in the database", func() {

			ginkgo.BeforeEach(func() {
				videoService.Create(testVideo)
				//videoController.Create()
			})

			ginkgo.It("should return at least one element", func() {
				videoList := videoService.FindAll()

				gomega.Expect(videoList).ShouldNot(gomega.BeEmpty())
			})

			ginkgo.It("should map the fields correctly", func() {
				firstVideo := videoService.FindAll()[0]

				gomega.Expect(firstVideo.Title).Should(gomega.Equal(TITLE))
				gomega.Expect(firstVideo.Description).Should(gomega.Equal(DESCRIPTION))
				gomega.Expect(firstVideo.URL).Should(gomega.Equal(URL))
				gomega.Expect(firstVideo.Author.FirstName).Should(gomega.Equal(FIRSTNAME))
				gomega.Expect(firstVideo.Author.LastName).Should(gomega.Equal(LASTNAME))
				gomega.Expect(firstVideo.Author.Email).Should(gomega.Equal(EMAIL))
			})

			/*ginkgo.AfterEach(func() {
				video := videoService.FindAll()[0]
				videoService.Delete(video)
			})*/

		})

		ginkgo.Context("If there are no videos in the database", func() {

			ginkgo.It("should return an empty list", func() {
				videos := videoService.FindAll()

				gomega.Expect(videos).Should(gomega.BeEmpty())
			})

		})
	})

	ginkgo.AfterSuite(func() {
		//videoRepository.CloseDB()
	})
})
