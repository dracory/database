# Proposal: Add Context Timeout Helpers

**Date:** 2025-06-29  
**Status:** Proposed  
**Type:** Enhancement  
**Priority:** Medium

## Description
While the module supports context for cancellation, it could benefit from helper functions for common timeout scenarios to reduce boilerplate and enforce best practices.

## Current Limitations
- Manual context timeout setup required
- Inconsistent timeout handling across codebase
- Easy to forget timeout handling
- No standard timeouts for different operations

## Proposed Changes
1. Add helper functions for common timeout scenarios:
   ```go
   // WithTimeout creates a new context with the specified timeout
   func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc)
   
   // WithDefaultTimeout creates a new context with a default timeout
   func WithDefaultTimeout(parent context.Context) (context.Context, context.CancelFunc)
   
   // WithOperationTimeout creates a new context with a timeout based on operation type
   func WithOperationTimeout(parent context.Context, opType OperationType) (context.Context, context.CancelFunc)
   ```

2. Define standard operation types:
   ```go
   type OperationType string
   
   const (
       OpRead   OperationType = "read"
       OpWrite  OperationType = "write"
       OpBatch  OperationType = "batch"
       OpAdmin  OperationType = "admin"
   )
   ```

3. Add configuration for default timeouts:
   ```go
   type TimeoutConfig struct {
       Default  time.Duration
       Read     time.Duration
       Write    time.Duration
       Batch    time.Duration
       Admin    time.Duration
   }
   ```

## Example Usage
```go
// Current way
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
rows, err := db.QueryContext(ctx, "SELECT * FROM users")

// Proposed way
ctx, cancel := database.WithDefaultTimeout(ctx)
defer cancel()
rows, err := db.QueryContext(ctx, "SELECT * FROM users")

// Or with operation-specific timeout
ctx, cancel := database.WithOperationTimeout(ctx, database.OpRead)
defer cancel()
rows, err := db.QueryContext(ctx, "SELECT * FROM users")
```

## Expected Benefits
- More consistent timeout handling
- Reduced boilerplate code
- Better defaults for common scenarios
- Easier to enforce best practices
- More maintainable code

## Potential Drawbacks
- Additional API surface
- Need to document proper usage
- Potential for misuse if timeouts are too aggressive
- May hide important context handling details

## Implementation Considerations
- Backward compatibility must be maintained
- Sensible defaults should be provided
- Configuration should be flexible
- Clear documentation is essential
- Consider adding metrics for timeouts

## Next Steps
1. Gather feedback on the proposed API
2. Finalize the timeout configuration structure
3. Implement the helper functions
4. Add tests for the new functionality
5. Update documentation with examples and best practices
