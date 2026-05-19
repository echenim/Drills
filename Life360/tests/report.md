# QA Phase 2 — Test Report

## Summary Table

| Test File | Tests Run | Passed | Failed |
|---|---|---|---|
| tests/unit/users-api.test.js | 7 | 7 | 0 |
| tests/unit/posts-api.test.js | 7 | 7 | 0 |
| tests/unit/components.test.jsx | 6 | 6 | 0 |
| tests/integration/full-flow.test.jsx | 3 | 3 | 0 |
| **Total** | **23** | **23** | **0** |

---

## Final `vitest run` Output

```
 RUN  v1.6.1 /Users/itachi/Dev/Drills/Life360

 ✓ tests/unit/posts-api.test.js  (7 tests) 31ms
 ✓ tests/unit/users-api.test.js  (7 tests) 28ms
 ✓ tests/unit/components.test.jsx  (6 tests) 406ms
 ✓ tests/integration/full-flow.test.jsx  (3 tests) 492ms

 Test Files  4 passed (4)
      Tests  23 passed (23)
   Start at  12:04:26
   Duration  1.60s (transform 121ms, setup 343ms, collect 791ms, tests 948ms, environment 2.09s, prepare 318ms)
```

---

## Coverage List

### Backend — Users API (tests/unit/users-api.test.js)
- GET /api/users — 200 + array of users
- GET /api/users/:id — 200 + single user object
- GET /api/users/:id — 404 when user does not exist
- POST /api/users — 201 happy path (creates + returns new user)
- POST /api/users — 400 when required fields (name, email) are missing
- PUT /api/users/:id — 200 updates and returns existing user
- DELETE /api/users/:id — 204, then GET confirms 404

### Backend — Posts API (tests/unit/posts-api.test.js)
- GET /api/posts — 200 + array of posts
- GET /api/posts/:id — 200 + single post object
- GET /api/posts/:id — 404 when post does not exist
- POST /api/posts — 201 happy path (creates post with valid userId)
- POST /api/posts — 400 when required fields missing; 404 when userId does not exist
- PUT /api/posts/:id — 200 updates and returns existing post
- DELETE /api/posts/:id — 204, then GET confirms 404

### Frontend — React Components (tests/unit/components.test.jsx)
- Renders UsersList with seed data — all user-card-{id} testids present
- Renders PostsList with mocked fetch — all post-card-{id} testids present
- Loading state — users-loading and posts-loading shown while data is in flight
- Error state — users-error from prop; posts-error when fetch rejects
- Create user via UserForm — POST called with correct payload, new card appears after refresh
- Create post via PostForm — POST called with userId, title, body; onPostCreated callback fired

### Integration — Full-stack flows (tests/integration/full-flow.test.jsx)
- Full user flow — loads seed users, creates a new user via form, new card appears in DOM
- Full post flow — loads seed posts, creates a post linked to userId=1, new post card appears
- API + UI happy path — app boots, fetches live API data, users-list renders without errors; health endpoint returns {status:"ok"}

---

## Caveats

### src/api/server.js change
The only src/ file touched. Two minimal additions:
1. The app.listen() call is now gated behind if (import.meta.url === `file://${process.argv[1]}`) so importing app in tests does not bind a TCP port.
2. A named export { app } was added alongside the existing export default app.
npm run server continues to work unchanged.

### Shared in-memory state
data.js exports mutable arrays. All test files that mutate state (POST/DELETE) call users.splice(...) / posts.splice(...) in a beforeEach to reset to seed data. Tests are therefore order-independent and isolated.

### act(...) warnings
React Testing Library emits cosmetic act(...) warnings for async state updates in a small number of tests. These are warnings only — no tests fail because of them. They arise because mocked fetch resolves synchronously outside React's scheduler; wrapping assertions in waitFor already handles timing correctly.
