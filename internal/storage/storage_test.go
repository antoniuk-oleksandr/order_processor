package storage_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/antoniuk-oleksandr/order_processor/internal/storage"
)

var _ = Describe("Storage", Label("unit"), func() {
	Context("when creating a storage", func() {
		It("should create a storage", func() {
			s := storage.NewStorage()

			Expect(s).ToNot(BeNil(), "expected storage to be created")
		})
	})

	Context("when the user does not exist", func() {
		It("should return 0 for non-existing user", func() {
			userId := 1
			amount := 0

			s := storage.NewStorage()

			result, ok := s.Get(userId)
			Expect(result).To(Equal(amount), "expected amount to be 0 for non-existing user")
			Expect(ok).To(BeFalse(), "expected ok to be false for non-existing user")
		})
	})

	Context("when the user exists", func() {
		It("should return correct amount for existing user", func() {
			userId := 1
			amount := 100

			s := storage.NewStorage()
			s.Add(userId, amount)

			result, ok := s.Get(userId)
			Expect(result).To(Equal(amount), "expected amount to be correct for existing user")
			Expect(ok).To(BeTrue(), "expected ok to be true for existing user")
		})
	})
})
