package blockmania

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bitset", func() {

	var b bitset

	BeforeEach(func() {
		b = bitset{}
	})

	Describe("clone", func() {
		var cloned *bitset

		BeforeEach(func() {
			b.cms = []uint64{1, 2, 3}
			b.prs = []uint64{4, 5}

			cloned = b.clone()
		})

		It("should clone the original bitset", func() {
			Expect(cloned.cms).To(Equal(b.cms))
			Expect(cloned.prs).To(Equal(b.prs))
		})
	})

	Describe("commitCount", func() {
		BeforeEach(func() {
			b.cms = []uint64{1, 2}
		})

		It("should return the right count", func() {
			actual := b.commitCount()
			expected := 2
			Expect(actual).To(Equal(expected))
		})
	})

	Describe("hasCommit", func() {
		BeforeEach(func() {
			b.cms = []uint64{1, 2, 3, 4}
		})

		Context("when the commit exists", func() {
			It("should return true", func() {
				actual := b.hasCommit(0)
				expected := true
				Expect(actual).To(Equal(expected))
			})
		})

		Context("when the commit does not exists", func() {
			It("should return false", func() {
				actual := b.hasCommit(1)
				expected := false
				Expect(actual).To(Equal(expected))
			})
		})
	})

	Describe("hasPrepare", func() {
		BeforeEach(func() {
			b.prs = []uint64{1, 2, 3, 4}
		})

		Context("when the prepare exists", func() {
			It("should return true", func() {
				actual := b.hasPrepare(0)
				expected := true
				Expect(actual).To(Equal(expected))
			})
		})

		Context("when the prepare does not exists", func() {
			It("should return false", func() {
				actual := b.hasPrepare(1)
				expected := false
				Expect(actual).To(Equal(expected))
			})
		})
	})

	Describe("prepareCount", func() {
		BeforeEach(func() {
			b.prs = []uint64{1, 2}
		})

		It("should return the right count", func() {
			actual := b.prepareCount()
			expected := 2
			Expect(actual).To(Equal(expected))
		})
	})

	Describe("setCommit", func() {
		BeforeEach(func() {
			b.cms = []uint64{1}
			b.setCommit(1)
		})

		It("should set the commit", func() {
			actual := b.cms
			expected := []uint64{3}
			Expect(actual).To(Equal(expected))
		})
	})

	Describe("setPrepare", func() {
		BeforeEach(func() {
			b.prs = []uint64{1}
			b.setPrepare(1)
		})

		It("should set the prepare", func() {
			actual := b.prs
			expected := []uint64{3}
			Expect(actual).To(Equal(expected))
		})
	})

	Describe("newBitset", func() {
		var newed *bitset

		BeforeEach(func() {
			newed = newBitset(423)
		})

		It("should create a new bitset", func() {
			Expect(newed.cms).To(Equal([]uint64{0, 0, 0, 0, 0, 0, 0}))
			Expect(newed.prs).To(Equal([]uint64{0, 0, 0, 0, 0, 0, 0}))
		})
	})
})
