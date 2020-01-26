package gomodrun_test

import (
	"os"
	"os/exec"
	"path"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dustinblackman/gomodrun"
)

const testPackage string = "github.com/dustinblackman/go-hello-world-test@v0.0.1/hello-world"
const testPackageNoGoMod string = "github.com/dustinblackman/go-hello-world-test-no-gomod@v0.0.2/hello-world-no-gomod"

var _ = Describe("pkg", func() {
	Context("GetPkgRoot", func() {
		It("it should return the current working directory that contains a go.mod", func() {
			cwd, _ := os.Getwd()
			dir, err := gomodrun.GetPkgRoot()
			Expect(err).To(BeNil())
			Expect(dir).To(Equal(cwd))
		})
	})

	Context("GetCommandVersionedPkgPath", func() {
		It("should get the binaries binaries command path", func() {
			cwd, _ := os.Getwd()
			cmdPath, err := gomodrun.GetCommandVersionedPkgPath(cwd, "hello-world")
			Expect(err).To(BeNil())
			Expect(cmdPath).To(Equal(testPackage))
		})
	})

	Context("GetCachedBin", func() {
		Context("with go.mod", func() {
			It("should return the bin path when it does not exist in cache", func() {
				cwd, _ := os.Getwd()
				goVersionOutput, _ := exec.Command("go", "version").Output()
				goVersion := strings.Split(string(goVersionOutput), " ")[2]
				os.RemoveAll(path.Join(".gomodrun", goVersion, "github.com/dustinblackman"))

				binPath, err := gomodrun.GetCachedBin(cwd, "hello-world", testPackage)
				Expect(err).To(BeNil())
				Expect(binPath).To(ContainSubstring(testPackage))
				Expect(binPath).To(ContainSubstring(".gomodrun"))
			})

			It("should return the bin path when it exists in cache", func() {
				cwd, _ := os.Getwd()
				binPath, err := gomodrun.GetCachedBin(cwd, "hello-world", testPackage)
				Expect(err).To(BeNil())
				Expect(binPath).To(ContainSubstring(testPackage))
				Expect(binPath).To(ContainSubstring(".gomodrun"))
			})
		})

		Context("without go.mod", func() {
			It("should return the bin path when it does not exist in cache", func() {
				cwd, _ := os.Getwd()
				goVersionOutput, _ := exec.Command("go", "version").Output()
				goVersion := strings.Split(string(goVersionOutput), " ")[2]
				os.RemoveAll(path.Join(".gomodrun", goVersion, "github.com/dustinblackman"))

				binPath, err := gomodrun.GetCachedBin(cwd, "hello-world-no-gomod", testPackageNoGoMod)
				Expect(err).To(BeNil())
				Expect(binPath).To(ContainSubstring(testPackageNoGoMod))
				Expect(binPath).To(ContainSubstring(".gomodrun"))
			})

			It("should return the bin path when it exists in cache", func() {
				cwd, _ := os.Getwd()
				binPath, err := gomodrun.GetCachedBin(cwd, "hello-world-no-gomod", testPackageNoGoMod)
				Expect(err).To(BeNil())
				Expect(binPath).To(ContainSubstring(testPackageNoGoMod))
				Expect(binPath).To(ContainSubstring(".gomodrun"))
			})
		})
	})
})
