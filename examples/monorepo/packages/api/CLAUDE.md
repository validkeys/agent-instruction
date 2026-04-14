<!-- BEGIN AGENT-INSTRUCTION -->
## Repository Structure

This is a monorepo managed with npm workspaces. All packages live in packages/* directory. Shared code belongs in packages/shared. Never duplicate code between packages.

## TypeScript Standards

Use strict TypeScript mode. Enable all strict flags in tsconfig.json. Prefer type over interface for composability. Use explicit return types on all exported functions.

- [TypeScript Configuration](tsconfig.base.json)

## Code Organization

Follow feature-based folder structure. Group related files by feature, not by type. Keep index.ts files minimal - use them only for public API exports.

## Import Standards

Use absolute imports via path aliases. Never use relative imports that go up more than one level (../..). Import order: external deps, internal packages, local files.

## Error Handling

Use custom error classes extending Error. Include error codes and context. Never throw strings. Always type catch blocks as unknown and validate error shape.

## Test Framework

Use Vitest for all tests. Configure with workspace-aware setup. Test files use .test.ts extension and live alongside implementation files.

## Test Organization

Use describe blocks to group related tests. Use test.each for parameterized tests. Keep test names descriptive using 'should' convention: 'should return error when input is invalid'.

## Test Coverage

Maintain minimum 80% code coverage. Run coverage checks in CI. Focus on edge cases and error paths. Don't write tests just to hit coverage metrics.

## Mocking Strategy

Mock external dependencies and I/O. Never mock internal functions. Use vi.mock() for module mocks. Prefer dependency injection over mocking when possible.

## Integration Tests

Write integration tests for API endpoints. Use supertest for HTTP testing. Test happy path, auth errors, validation errors, and rate limiting.

## Input Validation

Validate and sanitize all user input. Use Zod schemas for runtime validation. Never trust client-side validation alone. Validate on both client and server.

## Authentication

Use JWT tokens with short expiration (15 min access, 7 day refresh). Store tokens in httpOnly cookies. Implement refresh token rotation. Never store sensitive data in JWT payload.

## Authorization

Check permissions on every protected route. Use role-based access control (RBAC). Implement least privilege principle. Log all authorization failures for audit.

## Secrets Management

Never commit secrets to git. Use environment variables for all secrets. Use different secrets per environment. Rotate secrets regularly. Document all required env vars in .env.example.

## SQL Injection Prevention

Always use parameterized queries. Never concatenate user input into SQL strings. Use ORM query builders when possible. Validate all inputs even when using ORMs.

## XSS Prevention

Sanitize all user-generated content before rendering. Use React's default escaping. Never use dangerouslySetInnerHTML without sanitization. Set Content-Security-Policy headers.

## API Architecture

Use Express.js with TypeScript. Implement middleware chain: logging, auth, validation, rate limiting, error handling. Keep route handlers thin - delegate to service layer.

## Route Organization

Group routes by resource in separate files. Use express.Router() for each resource. Mount routers in main app.ts. Follow RESTful conventions: GET/POST/PUT/DELETE.

## Request Validation

Validate request body, query params, and path params using Zod schemas. Return 400 with detailed validation errors. Never process invalid requests.

## Response Format

Use consistent response format: { success: boolean, data?: any, error?: string }. Return appropriate HTTP status codes. Include request ID in all responses for tracing.

## Rate Limiting

Implement rate limiting per endpoint and per user. Use Redis for distributed rate limiting. Return 429 with Retry-After header. Document rate limits in API docs.

## Database Access

Use Prisma ORM for database access. Never expose Prisma models directly in API responses. Use DTO pattern to transform data. Implement connection pooling.

## Error Handling

Implement centralized error handler middleware. Log all errors with context. Never expose internal errors to clients. Return sanitized error messages.

- [Error Handler](src/middleware/errorHandler.ts)

<!-- END AGENT-INSTRUCTION -->