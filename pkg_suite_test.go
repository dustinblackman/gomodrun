package gomodrun_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGomodrun(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoModRun")
}
