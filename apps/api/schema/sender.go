package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Sender holds the schema definition for the Sender entity.
type Sender struct {
	ent.Schema
}

// Fields of the Sender.
func (Sender) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("name").
			NotEmpty(),
		field.String("email").
			Optional().
			Nillable(),
	}
}

// Edges of the Sender.
func (Sender) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("documents", Document.Type),
	}
}
