package gomodrun

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/otiai10/copy"
)

var _ = Describe("tidy", func() {
	var tempDir string
	var goVersion string

	BeforeSuite(func() {
		var err error
		goVersion, err = getGoVersion()
		if err != nil {
			panic(err)
		}

		tempDir, err = ioutil.TempDir(os.TempDir(), "gomodrun-tidy")
		if err != nil {
			panic(err)
		}

		//nolint:dogsled // Test file, don't need any of the extra values
		_, filename, _, _ := runtime.Caller(0)
		err = copy.Copy(path.Join(path.Dir(filename), "./tests/alternative-tools-dir"), tempDir)
		if err != nil {
			panic(err)
		}

		bins := []string{
			// Bins to keep
			"github.com/dustinblackman/go-hello-world-test@v0.0.2/hello-world",
			// Bins to drop
			"github.com/dustinblackman/go-hello-world-test@v0.0.1/hello-world",
		}

		emptyFolders := []string{
			"github.com/dustinblackman/go-hello-world-test-no-gomod@v0.0.1",
		}

		for _, binPath := range bins {
			fullPath := path.Join(tempDir, ".gomodrun", goVersion, binPath, path.Base(binPath))
			err = os.MkdirAll(path.Dir(fullPath), 0750)
			if err != nil {
				panic(err)
			}

			emptyFile, emptyErr := os.Create(fullPath)
			if emptyErr != nil {
				panic(err)
			}

			err = emptyFile.Close()
			if err != nil {
				panic(err)
			}
		}

		for _, folderPath := range emptyFolders {
			fullPath := path.Join(tempDir, ".gomodrun", goVersion, folderPath)
			err = os.MkdirAll(path.Dir(fullPath), 0750)
			if err != nil {
				panic(err)
			}
		}
	})

	AfterSuite(func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			panic(err)
		}
	})

	It("should clean outdated binaries and empty folders from .gomodrun", func() {
		baseDir := path.Join(tempDir, ".gomodrun", goVersion)

		err := Tidy(tempDir)
		Expect(err).To(BeNil())

		_, existsErr := os.Stat(path.Join(baseDir, "github.com/dustinblackman/go-hello-world-test@v0.0.2/hello-world/hello-world"))
		Expect(os.IsNotExist(existsErr)).To(BeFalse())

		_, existsErr = os.Stat(path.Join(baseDir, "github.com/dustinblackman/go-hello-world-test@v0.0.1/hello-world/hello-world"))
		Expect(os.IsNotExist(existsErr)).To(BeTrue())

		_, existsErr = os.Stat(path.Join(baseDir, "github.com/dustinblackman/go-hello-world-test@v0.0.1/hello-world"))
		Expect(os.IsNotExist(existsErr)).To(BeTrue())

		_, existsErr = os.Stat(path.Join(baseDir, "github.com/dustinblackman/go-hello-world-test@v0.0.1"))
		Expect(os.IsNotExist(existsErr)).To(BeTrue())

		_, existsErr = os.Stat(path.Join(baseDir, "github.com/dustinblackman/go-hello-world-test-no-gomod@v0.0.1"))
		Expect(os.IsNotExist(existsErr)).To(BeTrue())

	})
})
