# Go Echo Hexagonal Architecture Skill

## Purpose

This skill defines how an AI coding agent must work within this repository.

The project uses:

* Go
* Labstack Echo
* Hexagonal Architecture
* Dependency Injection
* Repository Pattern
* Use Case Pattern
* PostgreSQL or another relational database
* Unit and integration testing

The agent must prioritize:

* Correctness
* Maintainability
* Testability
* Security
* Clear separation of concerns
* Minimal and focused changes
* Human review before modifying or committing code

---

# Critical Workflow Rules

## Do Not Modify Code Immediately

When receiving a request, do not immediately:

* Edit source files
* Create files
* Delete files
* Rename files
* Run formatting commands that modify files
* Run code-generation tools that modify files
* Commit code
* Push code
* Create tags
* Rebase branches
* Merge branches

The first response must be an analysis and proposed implementation only.

The agent must show what it intends to change before applying any modification.

---

## Required Approval Workflow

Every coding task must follow these phases:

1. Understand
2. Inspect
3. Analyze
4. Design
5. Propose
6. Self-review
7. Wait for approval
8. Implement
9. Verify
10. Present final diff
11. Wait for commit approval
12. Commit only when explicitly requested

Do not skip phases.

---

# Phase 1: Understand the Request

Before proposing changes, identify:

* The requested behavior
* The expected result
* The affected business rule
* The likely affected layer
* Existing constraints
* Security implications
* Backward-compatibility concerns
* Testing requirements

Restate the task briefly.

Example:

```text
Task understanding:

The request is to add an endpoint that retrieves a user by ID.

Expected flow:

HTTP Request
→ Echo Handler
→ Application Use Case
→ Domain Port
→ Repository Adapter
→ Database
```

Do not modify code during this phase.

---

# Phase 2: Inspect the Existing Code

Inspect the repository before proposing implementation.

Review relevant files such as:

* Domain entities
* Domain errors
* Input ports
* Output ports
* Use cases
* HTTP handlers
* Request and response DTOs
* Repository adapters
* Database queries
* Dependency injection
* Route registration
* Middleware
* Tests
* Configuration
* Migrations

Do not assume the repository follows a generic structure.

Always follow the actual structure and conventions already used by the project.

When reporting findings, include:

* Relevant files
* Existing flow
* Existing patterns
* Potential inconsistencies
* Reusable components
* Risks

Example:

```text
Relevant files found:

- internal/domain/user.go
- internal/application/port/input/user_usecase.go
- internal/application/port/output/user_repository.go
- internal/application/service/user_service.go
- internal/adapter/input/http/user_handler.go
- internal/adapter/output/postgres/user_repository.go
```

---

# Phase 3: Analyze the Change

Analyze the request from the following perspectives.

## Functional Analysis

Check:

* Does the proposed behavior satisfy the requirement?
* What are the valid inputs?
* What are the invalid inputs?
* What happens when data is missing?
* What happens when dependencies fail?
* Are duplicate requests possible?
* Is the operation idempotent?
* Are transactions required?

## Architectural Analysis

Check:

* Does the change preserve dependency direction?
* Is business logic kept outside the HTTP handler?
* Is infrastructure isolated behind ports?
* Does the domain remain independent of Echo and database libraries?
* Are interfaces owned by the correct layer?
* Is dependency injection explicit?

## Security Analysis

Check:

* Input validation
* Authentication
* Authorization
* SQL injection
* Sensitive information exposure
* Error-message leakage
* Logging of credentials or tokens
* Rate limiting
* Replay attacks
* Path traversal
* File-upload validation
* Race conditions

## Data Analysis

Check:

* Database constraints
* Nullability
* Unique constraints
* Index requirements
* Transaction boundaries
* Concurrent updates
* Pagination
* Query performance
* Migration compatibility

## Testing Analysis

Identify required:

* Unit tests
* Handler tests
* Repository tests
* Integration tests
* Regression tests
* Failure-path tests

---

# Phase 4: Propose the Implementation

Before modifying files, provide a proposal containing:

## Summary

Explain the proposed solution in a few sentences.

## Files to Change

List each file that would be:

* Created
* Modified
* Deleted
* Renamed

Example:

```text
Files to modify:

1. internal/application/port/input/user_usecase.go
   - Add GetUserByID method.

2. internal/application/service/user_service.go
   - Implement GetUserByID business flow.

3. internal/adapter/input/http/user_handler.go
   - Add Echo handler.

4. internal/adapter/output/postgres/user_repository.go
   - Add parameterized database query.
```

## Proposed Flow

Show the request flow.

```text
Echo Handler
→ Validate path parameter
→ Call UserUseCase
→ Use case validates business rules
→ Repository port loads data
→ Handler maps result to HTTP response
```

## Proposed Interfaces

Show important interface changes.

```go
type UserUseCase interface {
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
}
```

## Proposed Code

Show representative code or a proposed diff.

The proposed code must be detailed enough for review but must not yet be written to the repository.

Use one of these formats:

* Proposed code blocks
* Unified diff
* File-by-file change preview

Prefer a unified diff when modifying existing code.

Example:

```diff
 type UserUseCase interface {
+    GetUserByID(ctx context.Context, id string) (*domain.User, error)
 }
```

---

# Phase 5: Engineering Loop Self-Review

Before asking for approval, perform an internal engineering review loop.

The loop must include:

1. Observe
2. Analyze
3. Design
4. Challenge
5. Verify
6. Refine
7. Summarize

Repeat the loop until no important unresolved issue remains.

Do not modify repository files during this loop.

---

## Step 1: Observe

Review:

* Existing code
* Existing conventions
* Current tests
* Build configuration
* Dependency direction
* Error handling
* Logging
* Database behavior

Report important observations.

---

## Step 2: Analyze

Check whether the proposed implementation:

* Solves the actual problem
* Handles edge cases
* Matches the existing architecture
* Avoids unnecessary changes
* Preserves backward compatibility
* Introduces any security risk

---

## Step 3: Design

Confirm:

* Correct layer ownership
* Correct interface placement
* Correct dependency direction
* Correct transaction boundary
* Correct request and response mapping
* Correct error model

---

## Step 4: Challenge

Act as a strict reviewer and challenge the proposal.

Ask internally:

* Can this fail under concurrent requests?
* Can the same request be submitted twice?
* Can this cause an SQL injection?
* Can this expose sensitive information?
* Can this return an incorrect HTTP status?
* Can this break an existing caller?
* Can this create a circular dependency?
* Is the interface too broad?
* Is this over-engineered?
* Is there a simpler implementation?
* Are errors being swallowed?
* Are errors being logged more than once?
* Is a transaction needed?
* Is the context propagated correctly?
* Can a goroutine leak?
* Can a database row remain partially updated?

Document discovered risks and refinements.

---

## Step 5: Verify

Before implementation, determine the commands that would be used later.

Typical commands:

```bash
go test ./...
go test -race ./...
go vet ./...
go build ./...
```

When available, also consider:

```bash
golangci-lint run
staticcheck ./...
```

Do not claim that these commands passed unless they were actually executed after implementation.

During the proposal stage, label them as planned verification commands.

---

## Step 6: Refine

Update the proposal based on issues found during self-review.

Clearly state any changes made to the initial design.

Example:

```text
Refinement:

The initial design placed validation in the handler.

After review, business-level validation was moved to the use case so that the
same rule applies to HTTP, scheduled jobs, and other adapters.
```

---

## Step 7: Summarize

Provide a final pre-implementation review summary:

```text
Self-review result:

- Architecture: Pass
- Dependency direction: Pass
- Input validation: Pass
- Error handling: Pass
- Transaction handling: Not required
- Security review: Pass
- Backward compatibility: Pass
- Test plan: Ready
- Repository files modified: No
- Commit created: No
```

Then wait for human review.

---

# Approval Gate 1: Permission to Implement

After presenting the proposal and self-review, stop.

Use wording similar to:

```text
No repository files have been modified.

Please review the proposed implementation.

After approval, I will apply the changes and run the verification loop.
```

Do not modify the code until the user explicitly approves the implementation.

Valid approvals may include:

* Approved
* Proceed
* Implement it
* Apply the changes
* Looks good, continue
* ผ่าน
* ทำต่อได้
* แก้โค้ดได้เลย

If the user requests changes, update the proposal and repeat the engineering self-review.

Do not treat a request to explain something as implementation approval.

---

# Implementation Phase

After implementation approval, apply only the approved changes.

Do not introduce unrelated refactoring.

Do not:

* Rename unrelated symbols
* Reformat unrelated files
* Upgrade dependencies without approval
* Change public APIs unnecessarily
* Rewrite working modules
* Add speculative abstractions
* Remove existing behavior without approval

Keep the diff minimal and focused.

---

# Hexagonal Architecture Rules

## Dependency Direction

Dependencies must point inward.

```text
Adapters → Application → Domain
```

The domain must not depend on:

* Echo
* SQL drivers
* ORM libraries
* HTTP libraries
* Framework-specific DTOs
* Configuration libraries
* Logging implementations

The application layer may depend on the domain.

Adapters may depend on application ports and domain types.

---

## Domain Layer

The domain layer contains:

* Entities
* Value objects
* Domain services
* Domain errors
* Business invariants
* Domain behavior

The domain layer must not contain:

* Echo context
* HTTP status codes
* SQL queries
* Database models tied to a driver
* JSON binding logic
* Environment-variable access
* Framework-specific annotations

Prefer behavior-rich domain models when appropriate.

Example:

```go
package domain

import "errors"

var ErrInvalidEmail = errors.New("invalid email")

type User struct {
	ID    string
	Email string
}

func NewUser(id, email string) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}

	return &User{
		ID:    id,
		Email: email,
	}, nil
}
```

---

## Application Layer

The application layer contains:

* Input ports
* Output ports
* Use cases
* Application services
* Transaction orchestration
* Business workflow coordination

Use cases must not depend on Echo.

Use cases should accept `context.Context` when performing I/O or invoking ports.

Example:

```go
type CreateUserUseCase interface {
	Execute(ctx context.Context, input CreateUserInput) (*domain.User, error)
}
```

Use-case inputs should represent application intent, not raw HTTP requests.

---

## Input Ports

Input ports define operations exposed by the application.

Keep interfaces focused.

Prefer:

```go
type GetUserUseCase interface {
	Execute(ctx context.Context, id string) (*domain.User, error)
}
```

Avoid overly broad interfaces:

```go
type UserService interface {
	Create(...)
	Update(...)
	Delete(...)
	Get(...)
	List(...)
	Export(...)
	Import(...)
}
```

Use interface segregation when operations have different consumers or dependencies.

---

## Output Ports

Output ports define capabilities required by the application.

Examples:

* Repository
* Message publisher
* File storage
* Email sender
* External API client
* Clock
* ID generator
* Transaction manager

Example:

```go
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Save(ctx context.Context, user *domain.User) error
}
```

Output ports should use domain or application types rather than adapter-specific models.

---

## HTTP Adapter

Echo handlers are input adapters.

Handlers are responsible for:

* Reading request data
* Binding request DTOs
* Basic syntactic validation
* Calling an input port
* Mapping application results
* Mapping errors to HTTP responses

Handlers must not contain:

* SQL queries
* Direct database access
* Complex business rules
* Transaction management
* Infrastructure-specific retry logic

Keep handlers thin.

Example:

```go
func (h *UserHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_USER_ID",
			Message: "user id is required",
		})
	}

	user, err := h.getUserUseCase.Execute(c.Request().Context(), id)
	if err != nil {
		return h.errorMapper.ToResponse(c, err)
	}

	return c.JSON(http.StatusOK, mapUserResponse(user))
}
```

Always pass the request context:

```go
ctx := c.Request().Context()
```

Do not replace it with `context.Background()`.

---

## Request and Response DTOs

HTTP DTOs belong to the HTTP adapter.

Do not reuse HTTP request DTOs as domain entities.

Example:

```go
type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}
```

Map the request into an application input:

```go
input := application.CreateUserInput{
	Email: req.Email,
}
```

Map domain or application results into response DTOs.

Do not expose internal database fields automatically.

---

## Repository Adapter

Repository implementations belong to output adapters.

Repositories are responsible for:

* Query execution
* Persistence mapping
* Converting database rows into domain objects
* Translating database-specific errors
* Respecting context cancellation

Always use parameterized queries.

Correct:

```go
row := r.db.QueryRowContext(
	ctx,
	`SELECT id, email FROM users WHERE id = $1`,
	id,
)
```

Incorrect:

```go
query := "SELECT id, email FROM users WHERE id = '" + id + "'"
```

Do not expose driver-specific errors directly to handlers.

Translate them into stable application or domain errors.

---

## Database Models

When database structure differs from domain structure, use separate persistence models.

Example:

```go
type userRow struct {
	ID        string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
```

Map explicitly:

```go
func mapUserRow(row userRow) *domain.User {
	return &domain.User{
		ID:    row.ID,
		Email: row.Email,
	}
}
```

Avoid spreading database tags throughout domain entities.

---

## Transactions

Transactions must be controlled at the application boundary when a use case performs multiple related writes.

A repository should not silently create independent transactions for operations that must be atomic together.

Use an explicit transaction abstraction when needed.

Example:

```go
type TransactionManager interface {
	WithinTransaction(
		ctx context.Context,
		fn func(ctx context.Context) error,
	) error
}
```

Within the transaction:

* Use the transaction-bound context or repository
* Return errors to trigger rollback
* Do not swallow errors
* Avoid external network calls unless the consistency tradeoff is understood

Document whether external side effects require:

* Outbox pattern
* Eventual consistency
* Compensation
* Idempotency key

---

# Echo-Specific Rules

## Route Registration

Group routes by bounded context or feature.

Example:

```go
func RegisterUserRoutes(group *echo.Group, handler *UserHandler) {
	users := group.Group("/users")

	users.POST("", handler.Create)
	users.GET("/:id", handler.GetByID)
}
```

Avoid constructing repositories or use cases inside route registration.

Dependencies must be assembled in the composition root.

---

## Middleware

Middleware may handle cross-cutting transport concerns such as:

* Authentication
* Request IDs
* Structured logging
* Panic recovery
* CORS
* Rate limiting
* Security headers

Middleware must not contain feature-specific business logic.

Do not trust client-provided identity values without verification.

Store authenticated identity in a typed and documented context mechanism.

---

## Error Responses

Return a consistent error structure.

Example:

```go
type ErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"requestId,omitempty"`
}
```

Do not expose:

* Stack traces
* SQL statements
* Database credentials
* Internal hostnames
* Raw driver errors
* Authentication secrets
* Tokens

Use stable error codes for clients.

Example:

```json
{
  "code": "USER_NOT_FOUND",
  "message": "user not found"
}
```

---

## HTTP Status Mapping

Use appropriate status codes.

Common mappings:

* Invalid syntax or request format: `400 Bad Request`
* Authentication missing or invalid: `401 Unauthorized`
* Authenticated but forbidden: `403 Forbidden`
* Resource not found: `404 Not Found`
* Duplicate or state conflict: `409 Conflict`
* Semantic validation failure: `422 Unprocessable Entity`
* Rate limit exceeded: `429 Too Many Requests`
* Unexpected internal failure: `500 Internal Server Error`
* Dependency unavailable: `503 Service Unavailable`

Do not return `200 OK` for failed operations unless an existing API contract explicitly requires it.

---

# Go Coding Rules

## Formatting

All modified Go files must be formatted with:

```bash
gofmt
```

When import management is available, use:

```bash
goimports
```

Do not reformat unrelated files.

---

## Error Handling

Always handle returned errors.

Avoid:

```go
result, _ := service.Execute(ctx)
```

Wrap errors only when adding useful context:

```go
return fmt.Errorf("load user %s: %w", id, err)
```

Use `errors.Is` and `errors.As` for error inspection.

Do not compare errors by their message text.

Correct:

```go
if errors.Is(err, domain.ErrUserNotFound) {
	// map error
}
```

Incorrect:

```go
if err.Error() == "user not found" {
	// map error
}
```

Avoid logging and returning the same error at every layer.

Log errors at a clear boundary where sufficient context exists.

---

## Context

Pass `context.Context` as the first parameter.

```go
func (s *Service) Execute(
	ctx context.Context,
	input Input,
) (*Output, error)
```

Do not store context in structs.

Do not use `context.Background()` inside request processing to bypass cancellation.

Respect cancellation in:

* Database queries
* External API calls
* Message publishing
* Long-running loops

---

## Constructors

Use explicit constructors.

```go
func NewUserService(
	repository output.UserRepository,
	clock output.Clock,
) *UserService {
	return &UserService{
		repository: repository,
		clock:      clock,
	}
}
```

Validate required dependencies when the project convention expects it.

Do not use package-level mutable dependencies.

---

## Interfaces

Define interfaces near the consumer when consistent with the project style.

Keep interfaces minimal.

Do not create an interface only for the sake of having an interface.

Interfaces are useful when:

* Defining application ports
* Supporting multiple adapters
* Enabling meaningful test doubles
* Isolating external dependencies

---

## Naming

Use clear names.

Prefer:

* `userRepository`
* `createUserUseCase`
* `FindByID`
* `CreateUserInput`
* `ErrUserNotFound`

Avoid vague names:

* `data`
* `obj`
* `manager`
* `helper`
* `util`
* `process`
* `doAction`

Use initialisms consistently:

* `ID`
* `HTTP`
* `URL`
* `API`
* `DTO`

Examples:

* `userID`
* `HTTPClient`
* `apiURL`

---

## Nil Handling

Avoid returning typed nil values through interfaces.

Validate nil dependencies where needed.

Return empty slices instead of nil only when required by the API contract.

Do not add defensive nil checks everywhere without understanding whether nil is valid.

---

## Concurrency

When adding goroutines:

* Define ownership
* Handle cancellation
* Handle errors
* Avoid data races
* Avoid unbounded goroutine creation
* Avoid blocked channels
* Ensure goroutines terminate
* Verify with the race detector

Do not add concurrency solely for perceived performance improvement.

---

# Validation Rules

Separate validation into:

## Transport Validation

Examples:

* Missing JSON field
* Invalid UUID format
* Invalid query parameter type
* Malformed date

This may occur in the HTTP adapter.

## Business Validation

Examples:

* User cannot approve their own request
* Account must be active
* Order cannot be cancelled after settlement
* Duplicate policy is not allowed

This belongs in the domain or application layer.

Do not rely only on frontend validation.

---

# Logging Rules

Use structured logging.

Include useful context such as:

* Request ID
* Trace ID
* User ID when safe
* Operation
* Resource ID
* Duration

Do not log:

* Passwords
* Access tokens
* Refresh tokens
* OTP values
* Authorization headers
* Private keys
* Full payment-card data
* Sensitive personal information

Avoid logging the same error at every layer.

---

# Testing Rules

## Unit Tests

Add or update unit tests for:

* Use-case success
* Validation failure
* Repository error
* Not-found behavior
* Conflict behavior
* Authorization failure
* Edge cases

Prefer table-driven tests where they improve clarity.

Example:

```go
func TestGetUserService_Execute(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		setupMock func(*mockUserRepository)
		wantErr   error
	}{
		{
			name:   "returns user successfully",
			userID: "user-1",
		},
		{
			name:    "returns not found",
			userID:  "missing-user",
			wantErr: domain.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			// Act
			// Assert
		})
	}
}
```

---

## Handler Tests

Handler tests should verify:

* Request parsing
* Validation
* Status codes
* Response structure
* Error mapping
* Use-case invocation

Do not retest all business rules through the handler.

---

## Repository Tests

Repository tests should verify:

* Query behavior
* Row mapping
* Database error mapping
* Not-found behavior
* Constraint conflicts
* Transaction behavior

Do not claim a repository is tested only because its mock was used in a service test.

---

## Test Quality

Tests must:

* Be deterministic
* Avoid dependence on execution order
* Avoid real external services in unit tests
* Assert important behavior
* Avoid testing implementation details unnecessarily
* Clean up created resources
* Avoid arbitrary sleeps

---

# Post-Implementation Engineering Loop

After implementation, perform the engineering loop again.

## 1. Observe

Review:

* Actual diff
* Changed files
* New dependencies
* Generated files
* Test changes

## 2. Validate Scope

Confirm:

* Only approved files changed
* No unrelated formatting occurred
* No accidental generated files were added
* No secrets were introduced

## 3. Run Formatting

Run on modified Go files:

```bash
gofmt
```

Use `goimports` if available.

## 4. Run Tests

Run the most relevant focused tests first.

Example:

```bash
go test ./internal/application/service/...
go test ./internal/adapter/input/http/...
go test ./internal/adapter/output/postgres/...
```

Then run the complete suite:

```bash
go test ./...
```

## 5. Run Static Analysis

Run:

```bash
go vet ./...
```

When configured:

```bash
golangci-lint run
staticcheck ./...
```

## 6. Run Race Detection

When the change affects concurrency, shared state, caching, goroutines, channels, or request parallelism:

```bash
go test -race ./...
```

## 7. Build

Run:

```bash
go build ./...
```

## 8. Review Failure Paths

Manually review:

* Invalid input
* Not found
* Duplicate data
* Database failure
* Timeout
* Context cancellation
* Unauthorized access
* Forbidden access
* Partial failure
* Concurrent execution

## 9. Review Architecture

Confirm:

* Domain does not import Echo
* Domain does not import database libraries
* Handlers contain no business logic
* Use cases depend only on ports
* Repositories implement output ports
* Dependency injection remains explicit
* No circular imports were introduced

## 10. Review Security

Confirm:

* Queries are parameterized
* Inputs are validated
* Authorization is enforced
* Sensitive data is not logged
* Internal errors are not leaked
* New endpoints have suitable middleware
* Rate limiting is considered for sensitive operations

## 11. Review Git Diff

Inspect:

```bash
git status --short
git diff --stat
git diff
```

Also check staged changes separately when relevant:

```bash
git diff --cached
```

Do not stage files automatically unless explicitly approved.

## 12. Refine

Fix issues discovered by the verification loop.

Repeat tests and analysis until the implementation is stable.

---

# Required Post-Implementation Report

After implementation and verification, report:

## Changes Made

List changed files and their purpose.

## Verification Results

Report actual command results.

Example:

```text
Verification:

- gofmt: Passed
- go test ./internal/application/service/...: Passed
- go test ./...: Passed
- go vet ./...: Passed
- go test -race ./...: Passed
- go build ./...: Passed
```

If a command was not run, explicitly say:

```text
- golangci-lint: Not run because the repository has no configuration.
```

Never claim a check passed if it was not executed.

## Issues Found and Fixed

Explain issues found during self-review.

## Remaining Risks

State unresolved concerns honestly.

## Git Status

Always state:

```text
- Files modified: Yes
- Files staged: No
- Commit created: No
- Changes pushed: No
```

Then show the final diff or summarize it file by file.

---

# Approval Gate 2: Permission to Commit

After implementation and verification, stop before committing.

Use wording similar to:

```text
The implementation is complete and has passed the checks listed above.

No commit has been created.

Please review the final diff. If approved, explicitly instruct me to create the commit.
```

Do not commit based only on implementation approval.

Implementation approval and commit approval are separate approvals.

---

# Commit Rules

Create a commit only when the user explicitly requests it after reviewing the implementation.

Valid commit instructions include:

* Commit this
* Create the commit
* Approved, commit it
* Commit with this message
* ผ่านแล้ว commit ได้
* approve ให้ commit

Before committing:

1. Run `git status --short`
2. Review `git diff`
3. Confirm only intended files are included
4. Stage only approved files
5. Create one focused commit unless instructed otherwise

Do not use:

```bash
git add .
```

Prefer staging explicit files:

```bash
git add internal/application/service/user_service.go
git add internal/application/service/user_service_test.go
```

Do not include:

* Secrets
* `.env`
* Local IDE files
* Temporary files
* Build artifacts
* Unrelated changes
* Debug logs

---

## Commit Message

Use the repository's existing commit convention.

If no convention exists, use Conventional Commits.

Examples:

```text
feat(user): add get user by id use case
```

```text
fix(auth): prevent OTP verification bypass
```

```text
refactor(order): isolate repository mapping
```

```text
test(user): add user service failure cases
```

Before committing, show the proposed commit message.

If the user provides an exact commit message, use it unless it is unsafe or clearly incorrect.

---

# Push Rules

A commit instruction does not automatically authorize pushing.

Do not run:

```bash
git push
```

unless the user explicitly asks to push.

Before pushing, report:

* Branch name
* Commit hash
* Remote destination

Do not force push unless explicitly requested and the consequences have been explained.

---

# Destructive Git Operations

Never run these without explicit permission:

```bash
git reset --hard
git clean -fd
git checkout -- .
git restore .
git rebase
git push --force
git push --force-with-lease
git branch -D
```

Do not overwrite or discard user changes.

If unrelated uncommitted changes exist, preserve them.

---

# Response Template for a New Task

Use the following structure before implementation:

```text
## Task Understanding

[Brief understanding of the request]

## Existing Code Findings

[Relevant architecture, files, and current behavior]

## Proposed Changes

1. [File and change]
2. [File and change]
3. [File and change]

## Proposed Flow

[Request and dependency flow]

## Proposed Code or Diff

[Code preview or unified diff]

## Engineering Loop Review

### Architecture
[Result]

### Correctness
[Result]

### Edge Cases
[Result]

### Security
[Result]

### Database and Transactions
[Result]

### Testing Plan
[Result]

### Planned Verification Commands

- gofmt
- go test ./...
- go vet ./...
- go build ./...

## Self-Review Result

- Architecture: Pass / Needs revision
- Correctness: Pass / Needs revision
- Security: Pass / Needs revision
- Testing plan: Ready / Incomplete
- Repository files modified: No
- Commit created: No

No repository files have been modified.

Please review the proposal. After approval, I will apply the changes and run the verification loop.
```

---

# Response Template After Implementation

```text
## Implementation Summary

[What was implemented]

## Files Changed

1. [File]
   - [Change]

2. [File]
   - [Change]

## Verification Results

- gofmt: Passed
- go test ./...: Passed
- go vet ./...: Passed
- go build ./...: Passed

## Engineering Loop Findings

[Issues found, refinements, and final assessment]

## Remaining Risks

[Known limitations or "None identified"]

## Git Status

- Files modified: Yes
- Files staged: No
- Commit created: No
- Changes pushed: No

The implementation is ready for review.

No commit has been created. Please explicitly approve the final diff before committing.
```

---

# Final Behavioral Requirements

The agent must always:

* Inspect before proposing
* Propose before modifying
* Show code before applying it
* Perform a self-review engineering loop
* Wait for implementation approval
* Apply only approved changes
* Run appropriate verification
* Report actual results honestly
* Show the final diff
* Wait for separate commit approval
* Commit only after explicit approval
* Push only after explicit push approval

The agent must never:

* Modify code immediately after receiving a task
* Hide implementation details before review
* Claim tests passed without running them
* Commit during the initial implementation
* Push automatically
* Discard unrelated user changes
* Mix unrelated refactoring into a task
* Bypass architecture boundaries for convenience
* Put business logic inside Echo handlers
* Put Echo or database dependencies inside the domain

# Sensitive Data Policy

The agent must never:

- Store secrets in memory outside the current task.
- Copy sensitive data into documentation.
- Include secrets in commit messages.
- Include secrets in pull request descriptions.
- Echo secrets unless explicitly requested.
- Persist credentials into generated files.

Sensitive data includes but is not limited to:

- API Keys
- Access Tokens
- Refresh Tokens
- JWT Tokens
- OAuth Credentials
- Client Secrets
- Database Passwords
- Connection Strings
- SSH Private Keys
- AWS Credentials
- GCP Credentials
- Azure Credentials
- Encryption Keys
- OTP Codes
- Session Cookies
- Authorization Headers
- Personal Information (PII)
- Production URLs containing credentials

When sensitive data is encountered:

- Redact it when displaying examples.
- Replace with placeholders such as:
  - <API_KEY>
  - <ACCESS_TOKEN>
  - <PASSWORD>
  - <JWT_TOKEN>

Never memorize, persist, or reuse sensitive values from previous conversations.
Treat all secrets as ephemeral.

# Git Commit Convention

Follow Conventional Commits.

Allowed commit types:

- feat
- fix
- refactor
- test
- docs
- style
- perf
- build
- ci
- chore
- revert

Format:

<type>(<scope>): <summary>

Examples:

feat(auth): add OTP verification endpoint

fix(auth): prevent OTP verification bypass

refactor(repository): simplify user query mapping

test(user): add service unit tests

docs(api): update authentication guide

perf(cache): optimize policy lookup

Summary rules:

- imperative mood
- lowercase
- under 72 characters
- no period at the end
- describe why when appropriate

Before committing:

1. Show the proposed commit message.
2. Wait for user approval.
3. Stage only approved files.
4. Create exactly one focused commit unless instructed otherwise.

# Branch Policy

Never switch branches automatically.

Never create a branch unless explicitly requested.

Never merge branches automatically.

Never rebase automatically.

Never add or upgrade dependencies without approval.

If a new dependency is required:

- explain why
- compare alternatives
- estimate impact
- wait for approval