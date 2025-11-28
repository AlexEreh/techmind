package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// DocumentTag holds the schema definition for the DocumentTag entity.
type DocumentTag struct {
	ent.Schema
}

// Fields of the DocumentTag.
func (DocumentTag) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("document_id", uuid.UUID{}),
		field.UUID("tag_id", uuid.UUID{}),
	}
}

// Edges of the DocumentTag.
func (DocumentTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("document", Document.Type).
			Ref("document_tags"),
		edge.From("tag", Tag.Type).
			Ref("document_tags"),
	}
}
