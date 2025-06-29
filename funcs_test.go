package database

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
)

func Test_IsQueryableContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "regular context",
			args: args{ctx: context.Background()},
			want: false,
		},
		{
			name: "queryable context",
			args: args{ctx: Context(context.Background(), nil)},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsQueryableContext(tt.args.ctx); got != tt.want {
				t.Errorf("IsQueryableContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Context(t *testing.T) {
	ctxBackground := context.Background()
	ctx := Context(ctxBackground, nil)

	if !reflect.DeepEqual(ctx, QueryableContext{Context: ctxBackground, queryable: nil}) {
		t.Errorf("Context() = %v, want %v", ctx, QueryableContext{Context: ctxBackground, queryable: nil})
	}

	if !IsQueryableContext(ctx) {
		t.Error(`IsQueryableContext() = `, IsQueryableContext(ctx), `, want `, true)
	}
}

func TestNewQueryableContextOr(t *testing.T) {
	db, _ := sql.Open("sqlite", ":memory:")
	db2, _ := sql.Open("sqlite", ":memory:")
	defer db.Close()
	defer db2.Close()

	// Case 1: Regular context should be converted to QueryableContext
	regularCtx := context.Background()
	qCtx1 := NewQueryableContextOr(regularCtx, db)

	if !IsQueryableContext(qCtx1) {
		t.Error("NewQueryableContextOr with regular context did not return a QueryableContext")
	}

	if qCtx1.Queryable() != db {
		t.Error("NewQueryableContextOr with regular context did not use the provided queryable")
	}

	// Case 2: Existing QueryableContext should be returned as is
	existingQCtx := NewQueryableContext(context.Background(), db)
	qCtx2 := NewQueryableContextOr(existingQCtx, db2) // db2 should be ignored

	if qCtx2.Queryable() != db {
		t.Error("NewQueryableContextOr with existing QueryableContext did not preserve the original queryable")
	}

	if qCtx2.Queryable() == db2 {
		t.Error("NewQueryableContextOr with existing QueryableContext incorrectly used the new queryable")
	}

	// They should be the same context
	if qCtx2 != existingQCtx {
		t.Error("NewQueryableContextOr with existing QueryableContext did not return the same context")
	}
}

func TestContextOr(t *testing.T) {
	db, _ := sql.Open("sqlite", ":memory:")
	db2, _ := sql.Open("sqlite", ":memory:")
	defer db.Close()
	defer db2.Close()

	// Case 1: Regular context should be converted to QueryableContext
	regularCtx := context.Background()
	qCtx1 := ContextOr(regularCtx, db)

	if !IsQueryableContext(qCtx1) {
		t.Error("ContextOr with regular context did not return a QueryableContext")
	}

	if qCtx1.Queryable() != db {
		t.Error("ContextOr with regular context did not use the provided queryable")
	}

	// Case 2: Existing QueryableContext should be returned as is
	existingQCtx := Context(context.Background(), db)
	qCtx2 := ContextOr(existingQCtx, db2) // db2 should be ignored

	if qCtx2.Queryable() != db {
		t.Error("ContextOr with existing QueryableContext did not preserve the original queryable")
	}

	if qCtx2.Queryable() == db2 {
		t.Error("ContextOr with existing QueryableContext incorrectly used the new queryable")
	}

	// They should be the same context
	if qCtx2 != existingQCtx {
		t.Error("ContextOr with existing QueryableContext did not return the same context")
	}
}
