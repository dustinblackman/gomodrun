package gomodrun_test

import (
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dustinblackman/gomodrun"
)

const (
	testPackage        string = "github.com/dustinblackman/go-hello-world-test@v0.0.2/hello-world"
	testPackageNoGoMod string = "github.com/dustinblackman/go-hello-world-test-no-gomod@v0.0.2/hello-world-no-gomod"
)

func formatForOS(pkgPath string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(pkgPath, "/", "\\")
	}

	return pkgPath
}

var _ = Describe("pkg", func() {
	cwd, _ := os.Getwd()
	goVersionOutput, _ := exec.Command("go", "version").Output()
	goVersion := strings.Split(string(goVersionOutput), " ")[2]

	Context("GetPkgRoot", func() {
		It("should return the current working directory that contains a go.mod", func() {
			dir, err := gomodrun.GetPkgRoot()
			Expect(err).To(BeNil())
			Expect(dir).To(Equal(cwd))
		})

		It("should return an error when it can not find a go.mod", func() {
			err := os.Chdir("../")
			if err != nil {
				panic(err)
			}

			dir, err := gomodrun.GetPkgRoot()
			Expect(err).ToNot(BeNil())
			Expect(dir).To(Equal(""))

			err = os.Chdir(cwd)
			if err != nil {
				panic(err)
			}
		})
	})

	Context("GetCommandVersionedPkgPath", func() {
		It("should get the binaries binaries command path", func() {
			cmdPath, err := gomodrun.GetCommandVersionedPkgPath(cwd, "hello-world")
			Expect(err).To(BeNil())
			Expect(cmdPath).To(Equal(testPackage))
		})

		It("should get the binaries binaries command path with suffix includes .exe", func() {
			cmdPath, err := gomodrun.GetCommandVersionedPkgPath(cwd, "hello-world.exe")
			Expect(err).To(BeNil())
			Expect(cmdPath).To(Equal(testPackage))
		})

		It("should throw an error when it cant find tools imports", func() {
			cmdPath, err := gomodrun.GetCommandVersionedPkgPath(path.Join(cwd, "../../"), "hello-world")
			Expect(err).ToNot(BeNil())
			Expect(cmdPath).To(Equal(""))
		})

		It("should throw an error when it cant find specified bin in imports", func() {
			cmdPath, err := gomodrun.GetCommandVersionedPkgPath(cwd, "not-real")
			Expect(err).ToNot(BeNil())
			Expect(strings.Contains(err.Error(), "cant find bin not-real in tools file")).To(BeTrue())
			Expect(cmdPath).To(Equal(""))
		})

		It("should throw an error when go.mod cant be found", func() {
			cmdPath, err := gomodrun.GetCommandVersionedPkgPath(path.Join(cwd, "./tests/missing-go-mod"), "hello-world")
			Expect(err).ToNot(BeNil())
			Expect(cmdPath).To(Equal(""))
		})

		It("should throw an error when go.mod is corrupted", func() {
			cmdPath, err := gomodrun.GetCommandVersionedPkgPath(path.Join(cwd, "./tests/corrupted-go-mod"), "hello-world")
			Expect(err).ToNot(BeNil())
			Expect(cmdPath).To(Equal(""))
		})

		It("should throw an error when go.mod is missing the requested dependency", func() {
			cmdPath, err := gomodrun.GetCommandVersionedPkgPath(path.Join(cwd, "./tests/incomplete-go-mod"), "hello-world")
			Expect(err).ToNot(BeNil())
			Expect(strings.Contains(err.Error(), "cant find require")).To(BeTrue())
			Expect(cmdPath).To(Equal(""))
		})
	})

	Context("GetCachedBin", func() {
		Context("with go.mod", func() {
			It("should return the bin path when it does not exist in cache", func() {
				err := os.RemoveAll(path.Join(".gomodrun", goVersion, "github.com/dustinblackman"))
				if err != nil {
					panic(err)
				}

				binPath, err := gomodrun.GetCachedBin(cwd, "hello-world", testPackage)
				Expect(err).To(BeNil())
				Expect(binPath).To(ContainSubstring(formatForOS(testPackage)))
				Expect(binPath).To(ContainSubstring(".gomodrun"))
			})

			It("should return the bin path when it exists in cache", func() {
				binPath, err := gomodrun.GetCachedBin(cwd, "hello-world", testPackage)
				Expect(err).To(BeNil())
				Expect(binPath).To(ContainSubstring(formatForOS(testPackage)))
				Expect(binPath).To(ContainSubstring(".gomodrun"))
			})
		})

		Context("without go.mod", func() {
			It("should return the bin path when it does not exist in cache", func() {
				err := os.RemoveAll(path.Join(".gomodrun", goVersion, "github.com/dustinblackman"))
				if err != nil {
					panic(err)
				}

				binPath, err := gomodrun.GetCachedBin(cwd, "hello-world-no-gomod", testPackageNoGoMod)
				Expect(err).To(BeNil())
				Expect(binPath).To(ContainSubstring(formatForOS(testPackageNoGoMod)))
				Expect(binPath).To(ContainSubstring(".gomodrun"))
			})

			It("should return the bin path when it exists in cache", func() {
				binPath, err := gomodrun.GetCachedBin(cwd, "hello-world-no-gomod", testPackageNoGoMod)
				Expect(err).To(BeNil())
				Expect(binPath).To(ContainSubstring(formatForOS(testPackageNoGoMod)))
				Expect(binPath).To(ContainSubstring(".gomodrun"))
			})
		})
	})

	Context("Run", func() {
		It("should return exit code 0 when binary exits with 0", func() {
			exitCode, err := gomodrun.Run("hello-world", []string{}, &gomodrun.Options{})
			Expect(err).To(BeNil())
			Expect(exitCode).To(Equal(0))
		})

		It("should return exit code 1 when the binary exists with 1", func() {
			exitCode, err := gomodrun.Run("hello-world", []string{"1"}, &gomodrun.Options{})
			Expect(err).To(BeNil())
			Expect(exitCode).To(Equal(1))
		})

		It("should return an error when the binary does not exist", func() {
			exitCode, err := gomodrun.Run("not-real", []string{}, &gomodrun.Options{})
			Expect(err).ToNot(BeNil())
			Expect(exitCode).To(Equal(-1))
		})

		Context("Alternative tools directory", func() {
			options := &gomodrun.Options{
				PkgRoot: path.Join(cwd, "./tests/alternative-tools-dir"),
			}

			It("should return exit code 0 when binary exits with 0", func() {
				exitCode, err := gomodrun.Run("hello-world", []string{}, options)
				Expect(err).To(BeNil())
				Expect(exitCode).To(Equal(0))
			})

			It("should return exit code 1 when the binary exists with 1", func() {
				exitCode, err := gomodrun.Run("hello-world", []string{"1"}, options)
				Expect(err).To(BeNil())
				Expect(exitCode).To(Equal(1))
			})

			It("should return an error when the binary does not exist", func() {
				exitCode, err := gomodrun.Run("not-real", []string{}, options)
				Expect(err).ToNot(BeNil())
				Expect(exitCode).To(Equal(-1))
			})
		})
	})
})
