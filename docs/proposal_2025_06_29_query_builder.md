# Proposal: Implement Query Builder (Low Priority)

**Date:** 2025-06-29  
**Status:** Proposed  
**Type:** Feature  
**Priority:** Low  
**Note:** Consider using existing SQL builders instead (see below)

## Description
The current implementation requires raw SQL strings for all database operations, which can be error-prone and harder to maintain. This proposal suggests adding a fluent query builder interface to improve developer experience and reduce SQL injection risks.

## Current Limitations
- Raw SQL strings are error-prone
- No compile-time validation of queries
- Harder to compose queries dynamically
- Increased risk of SQL injection
- No type safety for query parameters

## Recommended SQL Builders
Instead of implementing our own query builder, consider these established SQL builders for Go:

- [squirrel](https://github.com/Masterminds/squirrel) - Fluent SQL generator for Go
- [goqu](https://github.com/doug-martin/goqu) - Expressive SQL builder and query library
- [dbr](https://github.com/gocraft/dbr) - Additions to Go's database/sql for building queries
- [jet](https://github.com/go-jet/jet) - Type-safe SQL builder with code generation
- [go-sqlbuilder](https://github.com/huandu/go-sqlbuilder) - A flexible and powerful SQL builder
- [sqlair](https://github.com/canonical/sqlair) - Type-safe, runtime-agnostic query builder
- [sq](https://github.com/bokwoon95/go-structured-query) - Type-safe SQL query builder and struct mapper

### When to Consider This Proposal
This proposal should only be considered if:
1. None of the existing solutions meet specific requirements
2. There's a need for tight integration with this package
3. The maintenance burden of an external dependency is too high

## Proposed Changes
1. Introduce a new `QueryBuilder` interface:
   ```go
   type QueryBuilder interface {
       Select(columns ...string) SelectQuery
       Insert(table string) InsertQuery
       Update(table string) UpdateQuery
       Delete() DeleteQuery
   }
   ```

2. Implement builder interfaces for different query types:
   - `SelectQuery` with methods like `From()`, `Where()`, `Join()`, `GroupBy()`, etc.
   - `InsertQuery` with methods like `Columns()`, `Values()`, `Returning()`
   - `UpdateQuery` with methods like `Set()`, `Where()`
   - `DeleteQuery` with methods like `From()`, `Where()`

3. Add support for parameterized queries
4. Include type-safe where conditions
5. Add support for common SQL functions and operators

## Example Usage
```go
// Example of building a query
query := db.Select("id", "name", "email")
    .From("users")
    .Where("status = ?", "active")
    .OrderBy("created_at DESC")
    .Limit(10)

// Execute the query
rows, err := query.QueryContext(ctx)
```

## Expected Benefits
- Reduced SQL injection risks
- Better code completion and IDE support
- Easier query composition
- Improved code readability
- Type safety for query parameters
- Easier refactoring

## Potential Drawbacks
- Additional dependency or code to maintain
- Learning curve for new team members
- Potential performance overhead for complex queries
- May not cover all SQL features
- Additional abstraction layer to debug

## Alternatives Considered
1. **Raw SQL with Templates**: Using Go templates for SQL
   - Pros: More flexible, closer to SQL
   - Cons: Still error-prone, no type safety

2. **ORM Integration**: Full ORM like GORM
   - Pros: More features, active development
   - Cons: Heavier, more opinionated, steeper learning curve

## Implementation Consideration
Given the maturity and feature completeness of existing SQL builders, we strongly discourage implementing our own solution. The maintenance burden and feature parity with existing solutions would be significant.

## Recommendation
We recommend using one of the established SQL builders listed above. These libraries provide:

### Key Benefits
- Type-safe query construction
- Protection against SQL injection
- Fluent, chainable API
- Support for complex queries
- Active maintenance and community support
- Comprehensive documentation and examples
- Battle-tested in production environments

### Implementation Benefits
- No need to maintain custom query building code
- Immediate access to advanced features
- Better performance through optimization
- Community support and contributions
- Regular updates and security fixes

## Next Steps
1. Evaluate the recommended SQL builders for your specific use case
2. Consider factors like:
   - API design and ergonomics
   - Performance characteristics
   - Feature completeness
   - Community activity and support
   - Documentation quality
3. Prototype with the most promising candidates
4. Make an informed decision based on your requirements
