package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// CompanyUser holds the schema definition for the CompanyUser entity.
type CompanyUser struct {
	ent.Schema
}

// Fields of the CompanyUser.
func (CompanyUser) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("company_id", uuid.UUID{}),
		field.Int("role"),
		field.Time("added_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the CompanyUser.
func (CompanyUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("company_users").
			Field("user_id").
			Required().
			Unique(),
		edge.From("company", Company.Type).
			Ref("company_users").
			Field("company_id").
			Required().
			Unique(),
	}
}
