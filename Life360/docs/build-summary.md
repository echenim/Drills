# Life360 — Build Summary

A small full-stack app with a REST API for **users** and **posts**, a React UI that consumes it, and a green test suite. Built by a coordinated team of three Sonnet sub-agents (Backend, Frontend, QA) with this session as the orchestrator.

---

## What was built

```
Life360/
├── package.json               single project, all deps + scripts
├── vite.config.js             vite dev server on :3000, proxies /api/* → :3001, vitest config
├── index.html                 React mount point
├── src/
│   ├── main.jsx               React 18 entry
│   ├── App.jsx                tab toggle (Users / Posts), shared user state
│   ├── styles.css             utility CSS — cards grid, tabs, forms, status states
│   ├── api/
│   │   ├── server.js          Express app on :3001 (gated listen, named export `app` for tests)
│   │   ├── data.js            in-memory store, 3 seed users + 5 seed posts, id helpers
│   │   └── routes/
│   │       ├── users.js       full CRUD: GET/POST/PUT/DELETE, 400/404 error shapes
│   │       └── posts.js       full CRUD; POST validates referenced userId
│   └── components/
│       ├── UsersList.jsx      renders user cards from props (loading/error/empty)
│       ├── UserForm.jsx       controlled form → POST /api/users → triggers parent re-fetch
│       ├── PostsList.jsx      fetches /api/posts on mount + on version bump
│       └── PostForm.jsx       controlled form → POST /api/posts (userId from <select>)
├── tests/
│   ├── setup.js               vitest setup, imports jest-dom matchers
│   ├── README.md              test layout + run instructions
│   ├── report.md              QA pass/fail report (23/23 passing)
│   ├── unit/
│   │   ├── users-api.test.js     7 tests, supertest against the express app
│   │   ├── posts-api.test.js     7 tests, supertest against the express app
│   │   └── components.test.jsx   6 tests, React Testing Library + mocked fetch
│   └── integration/
│       └── full-flow.test.jsx    3 tests, real <App /> with supertest-bridged fetch
└── docs/
    └── build-summary.md       this file
```

### REST API contract (port 3001, proxied through Vite)

| Method | Path                | Body                              | Response                               |
|--------|---------------------|-----------------------------------|----------------------------------------|
| GET    | /api/health         | —                                 | `{status:"ok"}`                        |
| GET    | /api/users          | —                                 | `[{id,name,email,createdAt}]`          |
| GET    | /api/users/:id      | —                                 | user, or `404 {error}`                 |
| POST   | /api/users          | `{name,email}` (required)         | `201` user, or `400 {error}`           |
| PUT    | /api/users/:id      | `{name?,email?}`                  | user, or `404 {error}`                 |
| DELETE | /api/users/:id      | —                                 | `204` no body                          |
| GET    | /api/posts          | —                                 | `[{id,userId,title,body,createdAt}]`   |
| GET    | /api/posts/:id      | —                                 | post, or `404 {error}`                 |
| POST   | /api/posts          | `{userId,title,body}` (required)  | `201` post; `400` missing; `404` userId|
| PUT    | /api/posts/:id      | `{title?,body?,userId?}`          | post, or `404 {error}`                 |
| DELETE | /api/posts/:id      | —                                 | `204` no body                          |

Errors are always `{error: "<message>"}`. IDs are integers; `createdAt` is ISO 8601 UTC.

---

## Key decisions

- **Single project, two processes.** One `package.json` for the whole repo. Express serves `:3001`, Vite serves `:3000` with a `/api/*` proxy — so the React code calls `fetch('/api/users')` directly with no env vars or absolute URLs. `concurrently` runs both with `npm start`.
- **In-memory store, not a real DB.** The brief asked for a working app and a test report — adding SQLite/Postgres would have been scope creep. `data.js` exports the arrays so tests can splice them back to seed in `beforeEach`.
- **ESM throughout.** `"type": "module"` everywhere. No mixed `require`/`import`.
- **Vitest, not Jest.** Vitest reuses the Vite config, runs JSX out of the box, and handles both jsdom (component tests) and node (supertest tests) in one suite.
- **Testable Express app.** `src/api/server.js` exports the `app` and only calls `app.listen()` when run as the entrypoint (gated via `import.meta.url === \`file://${process.argv[1]}\``). This lets supertest hit the app in-process — no port juggling, no flaky teardown.
- **Integration tests via supertest-bridged `fetch`.** Rather than starting a real server inside vitest, `tests/integration/full-flow.test.jsx` swaps `global.fetch` with an adapter that funnels `/api/*` calls through `request(app)`. Keeps the suite hermetic and fast (~1.6s end to end).
- **Stable selectors via `data-testid`.** Every interactive element and status row has a testid (`tab-users`, `user-card-{id}`, `post-form-error`, etc.). Tests don't grep on copy.
- **Plain `useState` / `useEffect`.** No router, no state library, no auth — just the slice the brief asked for. App-level state owns the `users` array (so PostForm's `<select>` and PostsList's author lookup share the same source of truth).

---

## How to run it

From the project root: `/Users/itachi/Dev/Drills/Life360`.

### One-time
```
npm install         # already done in this session
```

### Run the app (both processes together)
```
npm start
```
Open **http://localhost:3000**. Backend is on `http://localhost:3001`.

### Run them separately (handy for debugging)
```
npm run server      # Express on :3001
npm run dev         # Vite on :3000 (proxies /api/* to :3001)
```

### Run the tests
```
npm test            # vitest run, 23 tests, ~1.6s
npm run test:watch  # watch mode
```

### Production build (frontend only)
```
npm run build       # outputs to dist/
```

---

## Verification status (this session)

- `npm test` — **4 files, 23 tests, all passing.** See `tests/report.md` for the breakdown.
- Backend smoke tests via Vite proxy at `:3000` — all green:
  - `GET /api/health` → `{status:"ok"}`
  - `GET /api/users` and `/api/posts` → seed data returned
  - `POST /api/users` → 201 + new user (id 4)
  - `POST /api/posts` → 201 + new post (id 6)
- React shell at `http://localhost:3000/` returns HTTP 200, `/src/main.jsx` is served by Vite.

Both servers were left running at the end of the session. To stop them: `pkill -f "node src/api/server.js"` and `pkill -f vite`.

---

## Known caveats

- **Shared in-memory state.** Restarting the API resets all data to the seed. Created users/posts disappear on restart — this is intentional for a demo, not a bug.
- **No DELETE in the UI.** v1 ships create + read only; the API supports DELETE and tests cover it, but no button is wired up. Add one to `UsersList`/`PostsList` if needed.
- **Cosmetic `act(...)` warnings** in component tests — async state updates inside mocked fetch callbacks. No test failures; cleaning these up is a future polish item.
- **`nextUserId` / `nextPostId` only go up.** Tests reset the arrays each `beforeEach`, but the id counters are not reset; integration tests assert the exact next ids (4, 6) the seed will produce. Adding a `resetCounters()` helper is the right move if more integration tests get added.
