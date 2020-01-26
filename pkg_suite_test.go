package gomodrun_test

import (
	"testing"

	"github.com/novln/macchiato"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGomodrun(t *testing.T) {
	RegisterFailHandler(Fail)
	macchiato.RunSpecs(t, "gomodrun")
}
