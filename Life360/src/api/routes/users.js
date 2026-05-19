import { Router } from "express";
import { users, getNextUserId } from "../data.js";

const router = Router();

// GET /api/users — list all users
router.get("/", (req, res) => {
  res.json(users);
});

// GET /api/users/:id — get one user
router.get("/:id", (req, res) => {
  const id = parseInt(req.params.id, 10);
  const user = users.find((u) => u.id === id);
  if (!user) {
    return res.status(404).json({ error: "User not found" });
  }
  res.json(user);
});

// POST /api/users — create user
router.post("/", (req, res) => {
  const { name, email } = req.body;
  if (!name || !email) {
    return res.status(400).json({ error: "name and email are required" });
  }
  const newUser = {
    id: getNextUserId(),
    name,
    email,
    createdAt: new Date().toISOString(),
  };
  users.push(newUser);
  res.status(201).json(newUser);
});

// PUT /api/users/:id — update user
router.put("/:id", (req, res) => {
  const id = parseInt(req.params.id, 10);
  const index = users.findIndex((u) => u.id === id);
  if (index === -1) {
    return res.status(404).json({ error: "User not found" });
  }
  const { name, email } = req.body;
  if (!name && !email) {
    return res.status(400).json({ error: "At least one of name or email is required" });
  }
  if (name) users[index].name = name;
  if (email) users[index].email = email;
  res.json(users[index]);
});

// DELETE /api/users/:id — delete user
router.delete("/:id", (req, res) => {
  const id = parseInt(req.params.id, 10);
  const index = users.findIndex((u) => u.id === id);
  if (index === -1) {
    return res.status(404).json({ error: "User not found" });
  }
  users.splice(index, 1);
  res.status(204).send();
});

export default router;
