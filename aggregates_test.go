package yit

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Aggregates", func() {
	Describe("AnyMatch", func() {
		It("returns true if any element matches the predicate", func() {
			result := FromNode(docNode).AnyMatch(All)
			Expect(result).To(BeTrue())
		})

		It("returns false if no element matches the predicate", func() {
			result := FromNode(docNode).AnyMatch(None)
			Expect(result).To(BeFalse())
		})
	})

	Describe("AllMatch", func() {
		It("returns true if all elements matches the predicate", func() {
			result := FromNodes(
				scalarNode("a"),
				scalarNode("a"),
			).AllMatch(WithValue("a"))

			Expect(result).To(BeTrue())
		})

		It("returns false if any element does not matches the predicate", func() {
			result := FromNodes(
				scalarNode("a"),
				scalarNode("b"),
			).AllMatch(WithValue("a"))

			Expect(result).To(BeFalse())
		})
	})

	Describe("ToArray", func() {
		It("adds all the iterated elements to an array", func() {
			nodes := []*yaml.Node{
				{Value: "a"},
				{Value: "b"},
				{Value: "c"},
			}

			result := FromNodes(nodes...).ToArray()
			Expect(result).To(Equal(nodes))
		})
	})
})
