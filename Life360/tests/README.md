# Test Suite — Life360 App

## Layout

```
tests/
  setup.js                      # Vitest setup: imports jest-dom matchers
  unit/
    users-api.test.js           # Users REST API contract tests (supertest)
    posts-api.test.js           # Posts REST API contract tests (supertest)
    components.test.jsx         # React component tests (React Testing Library)
  integration/
    full-flow.test.jsx          # End-to-end UI + API flows
```

## How to run

```bash
# Single run (CI-friendly)
npm test

# Watch mode (dev)
npm run test:watch
```

The backend server (`src/api/server.js`, port 3001) must be running separately when executing API or integration tests:

```bash
npm run server   # in one terminal
npm test         # in another
```

## Phase contract

All test bodies in phase 1 are `it.todo(...)` placeholders — they are discoverable by Vitest but produce no assertions. Phase 2 will replace each placeholder with real `supertest` calls (API tests) or `render` + `screen` queries (component/integration tests).

Environment: jsdom (configured in `vite.config.js`). Globals (`describe`, `it`, `expect`) are enabled — no imports needed in test files.
