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

const testPackage string = "github.com/dustinblackman/go-hello-world-test@v0.0.2/hello-world"
const testPackageNoGoMod string = "github.com/dustinblackman/go-hello-world-test-no-gomod@v0.0.2/hello-world-no-gomod"

var _ = Describe("pkg", func() {
	cwd, _ := os.Getwd()
	goVersionOutput, _ := exec.Command("go", "version").Output()
	goVersion := strings.Split(string(goVersionOutput), " ")[2]

	Context("GetPkgRoot", func() {
		It("it should return the current working directory that contains a go.mod", func() {
			dir, err := gomodrun.GetPkgRoot()
			Expect(err).To(BeNil())
			Expect(dir).To(Equal(cwd))
		})
	})

	Context("GetCommandVersionedPkgPath", func() {
		It("should get the binaries binaries command path", func() {
			cmdPath, err := gomodrun.GetCommandVersionedPkgPath(cwd, "hello-world")
			Expect(err).To(BeNil())
			Expect(cmdPath).To(Equal(testPackage))
		})
	})

	Context("GetCachedBin", func() {
		Context("with go.mod", func() {
			It("should return the bin path when it does not exist in cache", func() {
				os.RemoveAll(path.Join(".gomodrun", goVersion, "github.com/dustinblackman"))
				binPath, err := gomodrun.GetCachedBin(cwd, "hello-world", testPackage)
				Expect(err).To(BeNil())
				Expect(binPath).To(ContainSubstring(testPackage))
				Expect(binPath).To(ContainSubstring(".gomodrun"))
			})

			It("should return the bin path when it exists in cache", func() {
				binPath, err := gomodrun.GetCachedBin(cwd, "hello-world", testPackage)
				Expect(err).To(BeNil())
				Expect(binPath).To(ContainSubstring(testPackage))
				Expect(binPath).To(ContainSubstring(".gomodrun"))
			})
		})

		Context("without go.mod", func() {
			It("should return the bin path when it does not exist in cache", func() {
				os.RemoveAll(path.Join(".gomodrun", goVersion, "github.com/dustinblackman"))
				binPath, err := gomodrun.GetCachedBin(cwd, "hello-world-no-gomod", testPackageNoGoMod)
				Expect(err).To(BeNil())
				Expect(binPath).To(ContainSubstring(testPackageNoGoMod))
				Expect(binPath).To(ContainSubstring(".gomodrun"))
			})

			It("should return the bin path when it exists in cache", func() {
				binPath, err := gomodrun.GetCachedBin(cwd, "hello-world-no-gomod", testPackageNoGoMod)
				Expect(err).To(BeNil())
				Expect(binPath).To(ContainSubstring(testPackageNoGoMod))
				Expect(binPath).To(ContainSubstring(".gomodrun"))
			})
		})
	})

	Context("Run", func() {
		It("should return exit code 0 when binary exits with 0", func() {
			exitCode, err := gomodrun.Run("hello-world", []string{}, gomodrun.Options{})
			Expect(err).To(BeNil())
			Expect(exitCode).To(Equal(0))
		})

		It("should return exit code 1 when the binary exists with 1", func() {
			exitCode, err := gomodrun.Run("hello-world", []string{"1"}, gomodrun.Options{})
			Expect(err).To(BeNil())
			Expect(exitCode).To(Equal(1))
		})

		It("should return an error when the binary does not exist", func() {
			exitCode, err := gomodrun.Run("not-real", []string{}, gomodrun.Options{})
			Expect(err).ToNot(BeNil())
			Expect(exitCode).To(Equal(-1))
		})
	})
})
