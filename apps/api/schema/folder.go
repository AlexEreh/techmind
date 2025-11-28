package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Folder holds the schema definition for the Folder entity.
type Folder struct {
	ent.Schema
}

// Fields of the Folder.
func (Folder) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("company_id", uuid.UUID{}),
		field.UUID("parent_folder_id", uuid.UUID{}).
			Optional().
			Nillable(),
		field.String("name").
			NotEmpty(),
		field.Int64("size").
			Default(0).
			NonNegative(),
		field.Int("count").
			Default(0).
			NonNegative(),
	}
}

// Edges of the Folder.
func (Folder) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("company", Company.Type).
			Ref("folders").
			Field("company_id").
			Required().
			Unique(),
		edge.To("children", Folder.Type).
			From("parent").
			Field("parent_folder_id").
			Unique(),
		edge.To("documents", Document.Type),
	}
}
