// Unit/integration tests for the Users REST API using supertest.

import request from 'supertest';
import { app } from '../../src/api/server.js';
import { users } from '../../src/api/data.js';

// Seed data snapshot — mirrors what data.js initialises with.
const SEED_USERS = [
  { id: 1, name: 'Alice Johnson', email: 'alice@example.com', createdAt: '2024-01-10T08:00:00.000Z' },
  { id: 2, name: 'Bob Smith',     email: 'bob@example.com',   createdAt: '2024-01-11T09:00:00.000Z' },
  { id: 3, name: 'Carol White',   email: 'carol@example.com', createdAt: '2024-01-12T10:00:00.000Z' },
];

beforeEach(() => {
  // Reset in-memory users array back to seed state between tests.
  users.splice(0, users.length, ...SEED_USERS.map((u) => ({ ...u })));
});

describe('Users API', () => {
  it('GET /api/users — returns 200 with an array of users', async () => {
    const res = await request(app).get('/api/users');
    expect(res.status).toBe(200);
    expect(Array.isArray(res.body)).toBe(true);
    expect(res.body.length).toBe(3);
    expect(res.body[0]).toMatchObject({ id: 1, name: 'Alice Johnson', email: 'alice@example.com' });
  });

  it('GET /api/users/:id — returns 200 with a single user object', async () => {
    const res = await request(app).get('/api/users/2');
    expect(res.status).toBe(200);
    expect(res.body).toMatchObject({ id: 2, name: 'Bob Smith', email: 'bob@example.com' });
  });

  it('GET /api/users/:id — returns 404 when user does not exist', async () => {
    const res = await request(app).get('/api/users/9999');
    expect(res.status).toBe(404);
    expect(res.body).toHaveProperty('error');
  });

  it('POST /api/users — 201 happy path: creates and returns new user', async () => {
    const res = await request(app)
      .post('/api/users')
      .send({ name: 'Dave Test', email: 'dave@test.com' });
    expect(res.status).toBe(201);
    expect(res.body).toMatchObject({ name: 'Dave Test', email: 'dave@test.com' });
    expect(res.body).toHaveProperty('id');
    expect(res.body).toHaveProperty('createdAt');
  });

  it('POST /api/users — 400 validation error when required fields are missing', async () => {
    const res = await request(app)
      .post('/api/users')
      .send({ name: 'NoEmail' });
    expect(res.status).toBe(400);
    expect(res.body).toMatchObject({ error: 'name and email are required' });
  });

  it('PUT /api/users/:id — 200 updates an existing user', async () => {
    const res = await request(app)
      .put('/api/users/1')
      .send({ name: 'Alice Updated' });
    expect(res.status).toBe(200);
    expect(res.body).toMatchObject({ id: 1, name: 'Alice Updated' });
  });

  it('DELETE /api/users/:id — 204 removes an existing user', async () => {
    const del = await request(app).delete('/api/users/3');
    expect(del.status).toBe(204);

    // Confirm it is gone
    const get = await request(app).get('/api/users/3');
    expect(get.status).toBe(404);
  });
});
