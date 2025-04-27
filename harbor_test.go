package ginkgo_harbor_test

import (
	"fmt"

	"github.com/goharbor/tracker/ginkgo_harbor/envs"
	"github.com/goharbor/tracker/ginkgo_harbor/lib"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type HarborConInfo struct {
	Host     string
	Username string
	Password string
	ProxyURL string
}

var _ = Describe("Harbor", func() {
	var conInfo *HarborConInfo
	var rootURI string
	var harborEnv *envs.HarborEnvironment
	BeforeEach(func() {
		conInfo = &HarborConInfo{
			Host:     "zdj-dev.local",
			Username: "admin",
			Password: "Harbor12345",
			ProxyURL: "",
		}
		rootURI = fmt.Sprintf("https://%s", conInfo.Host)
		harborEnv = envs.NewHarborEnvironment(conInfo.Host, conInfo.Password, conInfo.ProxyURL)
	})

	It("has a system info", func() {
		Expect(conInfo).NotTo(BeNil())
		Expect(conInfo.Host).NotTo(BeEmpty())
		Expect(conInfo.Username).NotTo(BeEmpty())
		Expect(conInfo.Password).NotTo(BeEmpty())

		sys := lib.NewSystemUtil(rootURI, conInfo.Host, harborEnv.HTTPClient)
		Expect(sys).NotTo(BeNil())
		err := sys.GetSystemInfo()
		Expect(err).NotTo(HaveOccurred())
	})

	It("can create testing project", func() {
		Expect(conInfo).NotTo(BeNil())
		Expect(conInfo.Host).NotTo(BeEmpty())
		Expect(conInfo.Username).NotTo(BeEmpty())
		Expect(conInfo.Password).NotTo(BeEmpty())

		pro := lib.NewProjectUtil(rootURI, harborEnv.HTTPClient)
		err := pro.CreateProject(harborEnv.TestingProject, false)
		Expect(err).NotTo(HaveOccurred())

		By("create testing user", func() {

			usr := lib.NewUserUtil(rootURI, harborEnv.HTTPClient)
			img := lib.NewImageUtil(rootURI, harborEnv.HTTPClient)
			err := usr.CreateUser(harborEnv.Account, harborEnv.Password)
			Expect(err).NotTo(HaveOccurred())

			By("assign user001 as developer", func() {
				err := pro.AssignRole(harborEnv.TestingProject, harborEnv.Account)
				Expect(err).NotTo(HaveOccurred())
				dc := lib.DockerClient{}

				By("pushing image with this user")
				err = dc.PushImage(harborEnv)
				Expect(err).NotTo(HaveOccurred())

				By("scan the image")
				err = img.ScanArtifact(harborEnv.TestingProject, harborEnv.ImageName, harborEnv.ImageDigest)
				Expect(err).NotTo(HaveOccurred())

				By("pulling image with this user")
				err = dc.PullImage(harborEnv)
				Expect(err).NotTo(HaveOccurred())

				By("revoke user's role from project")
				err = pro.RevokeRole(harborEnv.TestingProject, harborEnv.Account)
				Expect(err).NotTo(HaveOccurred())

				By("pull image with this user, it should fail")
				err = dc.PullImage(harborEnv)
				Expect(err).To(HaveOccurred())

				By("delete testing repository")
				img := lib.NewImageUtil(rootURI, harborEnv.HTTPClient)
				err = img.DeleteRepo(harborEnv.TestingProject, harborEnv.ImageName)
				Expect(err).NotTo(HaveOccurred())

			})
			By("Delete testing user")
			err = usr.DeleteUser(harborEnv.Account)
			Expect(err).NotTo(HaveOccurred())

		})

		By("Delete testing project")
		err = pro.DeleteProject(harborEnv.TestingProject)
		Expect(err).NotTo(HaveOccurred())
	})

})
