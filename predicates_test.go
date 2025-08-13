package yit

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"go.yaml.in/yaml/v4"
)

var _ = Describe("Predicates", func() {
	Describe("WithKind", func() {
		predicate := WithKind(yaml.ScalarNode)

		It("returns true when nodes match the supplied Kind", func() {
			Expect(predicate(&yaml.Node{Kind: yaml.ScalarNode})).To(BeTrue())
		})

		It("returns false when nodes don't match the supplied Kind", func() {
			Expect(predicate(&yaml.Node{Kind: yaml.MappingNode})).To(BeFalse())
		})
	})

	DescribeTable("Truths",
		func(op func(p ...Predicate) Predicate, a, b, expected bool) {
			actual := op(
				func(node *yaml.Node) bool {
					return a
				}, func(node *yaml.Node) bool {
					return b
				},
			)

			Expect(actual(nil)).To(Equal(expected))
		},
		Entry("T | T = T", Union, true, true, true),
		Entry("T | F = T", Union, true, false, true),
		Entry("F | T = T", Union, false, true, true),
		Entry("F | F = F", Union, false, false, false),
		Entry("T & T = T", Intersect, true, true, true),
		Entry("T & F = F", Intersect, true, false, false),
		Entry("F & T = F", Intersect, false, true, false),
		Entry("F & F = F", Intersect, false, false, false),
	)

	Describe("WithShortTag", func() {
		predicate := WithShortTag("booooo")

		It("returns true when nodes match the tag", func() {
			Expect(predicate(&yaml.Node{Tag: "booooo"})).To(BeTrue())
		})

		It("returns false when nodes do not match the tag", func() {
			Expect(predicate(&yaml.Node{Tag: "not boooo"})).To(BeFalse())
		})
	})

	Describe("WithMapKey", func() {
		predicate := WithMapKey("a")

		It("returns true when the map has a specific key", func() {
			node := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: "a"},
				{Kind: yaml.ScalarNode, Value: "b"},
			}}

			result := predicate(node)
			Expect(result).To(BeTrue())
		})

		It("returns false when the map doesn't have a specific key", func() {
			node := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: "c"},
				{Kind: yaml.ScalarNode, Value: "d"},
			}}

			result := predicate(node)
			Expect(result).To(BeFalse())
		})

		It("returns false when the node isn't a map", func() {
			node := &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: "a"},
			}}

			result := predicate(node)
			Expect(result).To(BeFalse())
		})
	})

	Describe("WithMapValues", func() {
		predicate := WithMapValue("b")

		It("returns true when the map has a specific key", func() {
			node := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: "a"},
				{Kind: yaml.ScalarNode, Value: "b"},
			}}

			result := predicate(node)
			Expect(result).To(BeTrue())
		})

		It("returns false when the map doesn't have a specific key", func() {
			node := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: "c"},
				{Kind: yaml.ScalarNode, Value: "d"},
			}}

			result := predicate(node)
			Expect(result).To(BeFalse())
		})

		It("returns false when the node isn't a map", func() {
			node := &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: "b"},
			}}

			result := predicate(node)
			Expect(result).To(BeFalse())
		})
	})

	Describe("WithMapKeyValue", func() {
		predicate := WithMapKeyValue(
			// Key Predicate
			WithStringValue("a"),

			// Value Predicate
			WithStringValue("b"),
		)

		It("returns true when the map has a specific key value pair", func() {
			node := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: "a"},
				{Kind: yaml.ScalarNode, Value: "b"},
			}}

			result := predicate(node)
			Expect(result).To(BeTrue())
		})

		It("returns false when the map doesn't have a specific key value pair", func() {
			node := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: "a"},
				{Kind: yaml.ScalarNode, Value: "c"},
			}}

			result := predicate(node)
			Expect(result).To(BeFalse())
		})

		It("returns false when the node isn't a map", func() {
			node := &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: "a"},
				{Kind: yaml.ScalarNode, Value: "b"},
			}}

			result := predicate(node)
			Expect(result).To(BeFalse())
		})
	})

	Describe("WithPrefix", func() {
		predicate := WithPrefix("pre")

		It("returns true when the node's value has a prefix", func() {
			node := &yaml.Node{Value: "prefix"}
			result := predicate(node)
			Expect(result).To(BeTrue())
		})

		It("returns false when the node's value does not have a prefix", func() {
			node := &yaml.Node{Value: "postfix"}
			result := predicate(node)
			Expect(result).To(BeFalse())
		})
	})

	Describe("WithSuffix", func() {
		predicate := WithSuffix("fix")

		It("returns true when the node's value has a prefix", func() {
			node := &yaml.Node{Value: "prefix"}
			result := predicate(node)
			Expect(result).To(BeTrue())
		})

		It("returns false when the node's value does not have a prefix", func() {
			node := &yaml.Node{Value: "fixpost"}
			result := predicate(node)
			Expect(result).To(BeFalse())
		})
	})

	Describe("Negate", func() {
		It("reverses the result of a predicate", func() {
			Expect(Negate(All)(nil)).To(BeFalse())
			Expect(Negate(None)(nil)).To(BeTrue())
		})
	})
})
