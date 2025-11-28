package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Document holds the schema definition for the Document entity.
type Document struct {
	ent.Schema
}

// Fields of the Document.
func (Document) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("company_id", uuid.UUID{}),
		field.UUID("folder_id", uuid.UUID{}).
			Optional().
			Nillable(),
		field.String("name"),
		field.String("file_path").
			NotEmpty(),
		field.String("preview_file_path").
			Optional().
			Nillable(),
		field.Int64("file_size").
			Positive(),
		field.String("mime_type").
			NotEmpty(),
		field.String("checksum").
			NotEmpty(),
		field.UUID("sender_id", uuid.UUID{}).
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Document.
func (Document) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("company", Company.Type).
			Ref("documents").
			Field("company_id").
			Required().
			Unique(),
		edge.From("folder", Folder.Type).
			Ref("documents").
			Field("folder_id").
			Unique(),
		edge.From("sender", Sender.Type).
			Ref("documents").
			Field("sender_id").
			Unique(),
		edge.To("document_tags", DocumentTag.Type),
	}
}
