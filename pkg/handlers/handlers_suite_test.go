package handlers

import (
  . "github.com/onsi/ginkgo/v2"
  . "github.com/onsi/gomega"
  "testing"
)

func TestHandlers(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "Handlers Suite")
}