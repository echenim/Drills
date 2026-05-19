// Component tests using React Testing Library + mocked fetch.

import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import App from '../../src/App.jsx';
import UsersList from '../../src/components/UsersList.jsx';
import PostsList from '../../src/components/PostsList.jsx';
import UserForm from '../../src/components/UserForm.jsx';
import PostForm from '../../src/components/PostForm.jsx';

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

const MOCK_USERS = [
  { id: 1, name: 'Alice Johnson', email: 'alice@example.com', createdAt: '2024-01-10T08:00:00.000Z' },
  { id: 2, name: 'Bob Smith',     email: 'bob@example.com',   createdAt: '2024-01-11T09:00:00.000Z' },
];

const MOCK_POSTS = [
  { id: 1, userId: 1, title: 'Hello World', body: 'My first post!', createdAt: '2024-01-15T08:30:00.000Z' },
  { id: 2, userId: 2, title: 'React Tips',  body: 'Some React tips.', createdAt: '2024-01-16T09:00:00.000Z' },
];

/** Returns a Response-like mock that resolves to the given body. */
function makeFetchResponse(body, status = 200) {
  return Promise.resolve({
    ok: status >= 200 && status < 300,
    status,
    json: () => Promise.resolve(body),
  });
}

// ---------------------------------------------------------------------------
// Restore fetch after each test
// ---------------------------------------------------------------------------
afterEach(() => {
  vi.restoreAllMocks();
});

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe('React Components', () => {
  it('renders users list — displays fetched users in the DOM', async () => {
    render(
      <UsersList users={MOCK_USERS} loading={false} error={null} />
    );

    expect(screen.getByTestId('users-list')).toBeInTheDocument();
    expect(screen.getByTestId('user-card-1')).toBeInTheDocument();
    expect(screen.getByTestId('user-card-2')).toBeInTheDocument();
    expect(screen.getByText('Alice Johnson')).toBeInTheDocument();
    expect(screen.getByText('Bob Smith')).toBeInTheDocument();
  });

  it('renders posts list — displays fetched posts in the DOM', async () => {
    // PostsList fetches /api/posts internally via useEffect.
    global.fetch = vi.fn(() => makeFetchResponse(MOCK_POSTS));

    render(<PostsList users={MOCK_USERS} version={0} />);

    await waitFor(() => {
      expect(screen.getByTestId('posts-list')).toBeInTheDocument();
    });

    expect(screen.getByTestId('post-card-1')).toBeInTheDocument();
    expect(screen.getByTestId('post-card-2')).toBeInTheDocument();
    expect(screen.getByText('Hello World')).toBeInTheDocument();
    expect(screen.getByText('React Tips')).toBeInTheDocument();
  });

  it('handles loading state — shows a loading indicator while data is in flight', async () => {
    // UsersList receives loading=true directly.
    render(<UsersList users={[]} loading={true} error={null} />);
    expect(screen.getByTestId('users-loading')).toBeInTheDocument();

    // PostsList shows posts-loading while fetch is pending.
    let resolveFetch;
    global.fetch = vi.fn(
      () =>
        new Promise((resolve) => {
          resolveFetch = resolve;
        })
    );

    render(<PostsList users={[]} version={0} />);
    expect(screen.getByTestId('posts-loading')).toBeInTheDocument();

    // Clean up — resolve so the component doesn't leak.
    resolveFetch({ ok: true, status: 200, json: () => Promise.resolve([]) });
  });

  it('handles error state — shows an error message when the API call fails', async () => {
    // UsersList with error prop set.
    render(<UsersList users={[]} loading={false} error="Couldn't load users." />);
    expect(screen.getByTestId('users-error')).toBeInTheDocument();
    expect(screen.getByText("Couldn't load users.")).toBeInTheDocument();

    // PostsList when fetch rejects.
    global.fetch = vi.fn(() => Promise.reject(new Error('Network error')));
    render(<PostsList users={[]} version={0} />);

    await waitFor(() => {
      expect(screen.getByTestId('posts-error')).toBeInTheDocument();
    });
  });

  it('creates a new user — form submit calls POST and updates the list', async () => {
    const user = userEvent.setup();

    // App fetches /api/users on mount, then again after creation.
    const newUser = { id: 3, name: 'Dave Test', email: 'dave@test.com', createdAt: new Date().toISOString() };
    global.fetch = vi
      .fn()
      // Initial load
      .mockImplementationOnce(() => makeFetchResponse(MOCK_USERS))
      // POST create
      .mockImplementationOnce(() => makeFetchResponse(newUser, 201))
      // Refresh after creation
      .mockImplementationOnce(() => makeFetchResponse([...MOCK_USERS, newUser]));

    render(<App />);

    // Wait for initial load to finish.
    await waitFor(() => expect(screen.getByTestId('users-list')).toBeInTheDocument());

    // Fill in the form.
    await user.type(screen.getByTestId('user-name-input'), 'Dave Test');
    await user.type(screen.getByTestId('user-email-input'), 'dave@test.com');
    await user.click(screen.getByTestId('user-submit-btn'));

    // After refresh, new user card should appear.
    await waitFor(() => {
      expect(screen.getByTestId('user-card-3')).toBeInTheDocument();
    });

    // POST was called with correct payload.
    const postCall = global.fetch.mock.calls.find(
      ([url, opts]) => opts && opts.method === 'POST'
    );
    expect(postCall).toBeDefined();
    const body = JSON.parse(postCall[1].body);
    expect(body).toMatchObject({ name: 'Dave Test', email: 'dave@test.com' });
  });

  it('creates a new post — form submit calls POST with userId and updates the list', async () => {
    const user = userEvent.setup();

    const newPost = {
      id: 3, userId: 1, title: 'Brand New Post', body: 'Fresh content here.',
      createdAt: new Date().toISOString(),
    };

    // PostForm doesn't fetch; it receives users as a prop.
    // We render PostForm directly.
    global.fetch = vi
      .fn()
      // POST /api/posts
      .mockImplementationOnce(() => makeFetchResponse(newPost, 201));

    const onPostCreated = vi.fn();

    render(
      <PostForm users={MOCK_USERS} usersLoading={false} onPostCreated={onPostCreated} />
    );

    // Select author, fill title and body.
    await user.selectOptions(screen.getByTestId('post-user-select'), '1');
    await user.type(screen.getByTestId('post-title-input'), 'Brand New Post');
    await user.type(screen.getByTestId('post-body-input'), 'Fresh content here.');
    await user.click(screen.getByTestId('post-submit-btn'));

    await waitFor(() => expect(onPostCreated).toHaveBeenCalledTimes(1));

    // POST was called with correct payload.
    const postCall = global.fetch.mock.calls[0];
    expect(postCall[1].method).toBe('POST');
    const body = JSON.parse(postCall[1].body);
    expect(body).toMatchObject({ userId: 1, title: 'Brand New Post', body: 'Fresh content here.' });
  });
});
