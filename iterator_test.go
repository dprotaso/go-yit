package yit

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"go.yaml.in/yaml/v3"
)

var _ = Describe("Iterator", func() {
	Describe("FromNode", func() {
		It("returns the root node", func() {
			node := &yaml.Node{}
			next := FromNode(node)

			item, ok := next()
			Expect(item).To(Equal(node))
			Expect(ok).To(BeTrue())

			item, ok = next()
			Expect(item).To(BeNil())
			Expect(ok).To(BeFalse())
		})
	})

	DescribeTable("RecurseNodes",
		func(yaml string, values ...*yaml.Node) {
			doc := toYAML(yaml)
			next := FromNode(doc).RecurseNodes()

			for _, value := range values {
				node, ok := next()

				Expect(ok).To(BeTrue())
				Expect(node.Kind).To(Equal(value.Kind))
				Expect(node.Value).To(Equal(value.Value))
			}

			_, ok := next()
			Expect(ok).To(BeFalse())
		},

		Entry("scalar", "a",
			docNode, scalarNode("a")),

		Entry("sequence", "[a, b, c]",
			docNode, seqNode, scalarNode("a"), scalarNode("b"), scalarNode("c")),

		Entry("map", "{a: b}",
			docNode, mapNode, scalarNode("a"), scalarNode("b")),
	)

	DescribeTable("Values",
		func(yaml string, values ...*yaml.Node) {
			doc := toYAML(yaml)
			next := FromNode(doc).
				Values(). // the root is the document node
				Values()

			for _, value := range values {
				node, ok := next()

				Expect(ok).To(BeTrue())
				Expect(node.Kind).To(Equal(value.Kind))
				Expect(node.Value).To(Equal(value.Value))
			}

			_, ok := next()
			Expect(ok).To(BeFalse())

		},

		Entry("scalar", nil /* no values */),

		Entry("sequence", "[a, b, c]",
			scalarNode("a"), scalarNode("b"), scalarNode("c")),

		Entry("map", "a: b\nc: d",
			scalarNode("a"), scalarNode("b"), scalarNode("c"), scalarNode("d"),
		),
	)

	Describe("Filter", func() {
		It("passes items through satisfying the predicate", func() {
			next := FromNode(docNode).Filter(All)
			node, ok := next()

			Expect(ok).To(BeTrue())
			Expect(node).To(Equal(docNode))
		})

		It("does not pass items that do not satisfy the predicate", func() {
			next := FromNode(docNode).Filter(None)
			_, ok := next()
			Expect(ok).To(BeFalse())
		})

		It("predicate is not invoked when there are no items", func() {
			empty := Iterator(func() (*yaml.Node, bool) {
				return nil, false
			})

			next := empty.Filter(func(*yaml.Node) bool {
				Fail("unexpected invocation of the filter")
				return true
			})

			_, ok := next()
			Expect(ok).To(BeFalse())
		})
	})

	Describe("MapKeys", func() {
		It("returns the keys of a map", func() {
			next := FromNode(toYAML("a: b\nc: d\ne: f")).
				RecurseNodes().
				Filter(WithKind(yaml.MappingNode)).
				MapKeys()

			for _, value := range []string{"a", "c", "e"} {
				node, ok := next()
				Expect(ok).To(BeTrue())
				Expect(node.Value).To(Equal(value))
			}

			_, ok := next()
			Expect(ok).To(BeFalse())
		})

		It("returns nothing for sequences", func() {
			next := FromNode(toYAML("[a, b, c, d]")).
				RecurseNodes().
				Filter(WithKind(yaml.SequenceNode)).
				MapKeys()

			_, ok := next()
			Expect(ok).To(BeFalse())
		})
	})

	Describe("MapValues", func() {
		It("returns the keys of a map", func() {
			next := FromNode(toYAML("a: b\nc: d\ne: f")).
				RecurseNodes().
				Filter(WithKind(yaml.MappingNode)).
				MapValues()

			for _, value := range []string{"b", "d", "f"} {
				node, ok := next()
				Expect(ok).To(BeTrue())
				Expect(node.Value).To(Equal(value))
			}

			_, ok := next()
			Expect(ok).To(BeFalse())
		})

		It("returns nothing for sequences", func() {
			next := FromNode(toYAML("[a, b, c, d]")).
				RecurseNodes().
				Filter(WithKind(yaml.SequenceNode)).
				MapValues()

			_, ok := next()
			Expect(ok).To(BeFalse())
		})
	})

	Describe("Iterate", func() {
		It("custom iterators can be supplied", func() {
			repeater := func(next Iterator) Iterator {
				return func() (node *yaml.Node, ok bool) {
					node, ok = next()
					if ok {
						node = scalarNode(strings.Repeat(node.Value, 2))
					}
					return
				}
			}

			next := FromNodes(scalarNode("a")).
				Iterate(repeater).
				Iterate(repeater)

			node, ok := next()
			Expect(ok).To(BeTrue())
			Expect(node.Value).To(Equal("aaaa"))
		})
	})

	Describe("ValuesForMap", func() {
		It("returns the values of a map matching the key/value predicates", func() {
			next := FromNode(toYAML("a: b\nc: d\ne: f")).
				RecurseNodes().
				Filter(WithKind(yaml.MappingNode)).
				ValuesForMap(All, func(node *yaml.Node) bool {
					return node.Value == "d"
				})

			node, ok := next()
			Expect(ok).To(BeTrue())
			Expect(node.Value).To(Equal("d"))

			_, ok = next()
			Expect(ok).To(BeFalse())
		})

		It("returns nothing for sequences", func() {
			next := FromNode(toYAML("[a, b, c, d]")).
				RecurseNodes().
				Filter(WithKind(yaml.SequenceNode)).
				ValuesForMap(All, All)

			_, ok := next()
			Expect(ok).To(BeFalse())
		})
	})

	Describe("FromIterators", func() {
		It("merges multiple iterators into a single stream", func() {
			next := FromIterators(
				FromNode(&yaml.Node{Value: "a"}),
				FromNode(&yaml.Node{Value: "b"}),
				FromNode(&yaml.Node{Value: "c"}),
			)

			for _, value := range []string{"a", "b", "c"} {
				node, ok := next()
				Expect(ok).To(BeTrue(), value+" to be present")
				Expect(node.Value).To(Equal(value))
			}

			_, ok := next()
			Expect(ok).To(BeFalse())
		})
	})
})

var mapNode = &yaml.Node{
	Kind: yaml.MappingNode,
}

var seqNode = &yaml.Node{
	Kind: yaml.SequenceNode,
}

var docNode = &yaml.Node{
	Kind: yaml.DocumentNode,
}

func scalarNode(val string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: val,
	}
}

func toYAML(s string) *yaml.Node {
	var node yaml.Node

	err := yaml.Unmarshal([]byte(s), &node)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())

	return &node
}
