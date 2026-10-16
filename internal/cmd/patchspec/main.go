// Command patchspec rewrites Together AI's OpenAPI document into a shape
// oapi-codegen can generate from. It reads the downloaded spec, applies each
// patch to the YAML tree and writes the result next to it. The downloaded
// spec stays untouched, so every change lives in this file.
package main

import (
	"bytes"
	"context"
	"log/slog"
	"os"

	"github.com/psyb0t/ctxerrors"
	"gopkg.in/yaml.v3"
)

const (
	specIn       = "api/openapi.yaml"
	specOut      = "api/openapi.patched.yaml"
	specFileMode = 0o600
	yamlIndent   = 2

	keyDiscriminator = "discriminator"
	keyMapping       = "mapping"
	keyOneOf         = "oneOf"
	keyAnyOf         = "anyOf"
	keyType          = "type"
	keyItems         = "items"
	keyEnum          = "enum"
	typeArray        = "array"
)

// patch changes the tree in place and returns how many spots it touched.
type patch struct {
	name  string
	apply func(root *yaml.Node) int
}

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "patch spec", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	raw, err := os.ReadFile(specIn)
	if err != nil {
		return ctxerrors.Wrap(err, "read spec")
	}

	var root yaml.Node
	if err := yaml.Unmarshal(raw, &root); err != nil {
		return ctxerrors.Wrap(err, "parse spec")
	}

	patches := []patch{
		{name: "unmapped_discriminators", apply: dropUnmappedDiscriminators},
		{name: "union_array_item_enums", apply: dropUnionArrayItemEnums},
	}

	for _, p := range patches {
		slog.InfoContext(ctx, "patched spec",
			"patch", p.name, "count", p.apply(&root))
	}

	var out bytes.Buffer

	enc := yaml.NewEncoder(&out)
	enc.SetIndent(yamlIndent)

	if err := enc.Encode(&root); err != nil {
		return ctxerrors.Wrap(err, "encode spec")
	}

	if err := enc.Close(); err != nil {
		return ctxerrors.Wrap(err, "close encoder")
	}

	if err := os.WriteFile(specOut, out.Bytes(), specFileMode); err != nil {
		return ctxerrors.Wrap(err, "write spec")
	}

	return nil
}

// walk calls fn on every mapping node in the tree.
func walk(n *yaml.Node, fn func(m *yaml.Node)) {
	if n.Kind == yaml.MappingNode {
		fn(n)
	}

	for _, c := range n.Content {
		walk(c, fn)
	}
}

// get returns the value of key in mapping m, or nil.
func get(m *yaml.Node, key string) *yaml.Node {
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			return m.Content[i+1]
		}
	}

	return nil
}

// del removes key from mapping m and reports whether it was there.
func del(m *yaml.Node, key string) bool {
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			m.Content = append(m.Content[:i], m.Content[i+2:]...)

			return true
		}
	}

	return false
}

// dropUnmappedDiscriminators removes every discriminator without a mapping.
// oapi-codegen refuses a discriminator over inline oneOf members it cannot
// name ("not all schemas were mapped"). The union still decodes through its
// As* and From* methods.
func dropUnmappedDiscriminators(root *yaml.Node) int {
	count := 0

	walk(root, func(m *yaml.Node) {
		d := get(m, keyDiscriminator)
		if d == nil || d.Kind != yaml.MappingNode || get(d, keyMapping) != nil {
			return
		}

		if del(m, keyDiscriminator) {
			count++
		}
	})

	return count
}

// dropUnionArrayItemEnums removes the inline enum from the items of an array
// that is a oneOf or anyOf member. oapi-codegen names that enum and the
// member's array type the same and stops with "duplicate typename". The items
// become plain strings; the allowed values stay in the description.
func dropUnionArrayItemEnums(root *yaml.Node) int {
	count := 0

	walk(root, func(m *yaml.Node) {
		for _, key := range []string{keyOneOf, keyAnyOf} {
			members := get(m, key)
			if members == nil || members.Kind != yaml.SequenceNode {
				continue
			}

			for _, member := range members.Content {
				if dropArrayItemEnum(member) {
					count++
				}
			}
		}
	})

	return count
}

// dropArrayItemEnum removes the enum from member's items when member is an
// array schema, and reports whether it did.
func dropArrayItemEnum(member *yaml.Node) bool {
	if member.Kind != yaml.MappingNode {
		return false
	}

	t := get(member, keyType)
	if t == nil || t.Value != typeArray {
		return false
	}

	items := get(member, keyItems)

	return items != nil && items.Kind == yaml.MappingNode && del(items, keyEnum)
}
