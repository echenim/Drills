import { Router } from "express";
import { posts, users, getNextPostId } from "../data.js";

const router = Router();

// GET /api/posts — list all posts
router.get("/", (req, res) => {
  res.json(posts);
});

// GET /api/posts/:id — get one post
router.get("/:id", (req, res) => {
  const id = parseInt(req.params.id, 10);
  const post = posts.find((p) => p.id === id);
  if (!post) {
    return res.status(404).json({ error: "Post not found" });
  }
  res.json(post);
});

// POST /api/posts — create post
router.post("/", (req, res) => {
  const { userId, title, body } = req.body;
  if (!userId || !title || !body) {
    return res.status(400).json({ error: "userId, title, and body are required" });
  }
  const userExists = users.find((u) => u.id === parseInt(userId, 10));
  if (!userExists) {
    return res.status(404).json({ error: "User not found" });
  }
  const newPost = {
    id: getNextPostId(),
    userId: parseInt(userId, 10),
    title,
    body,
    createdAt: new Date().toISOString(),
  };
  posts.push(newPost);
  res.status(201).json(newPost);
});

// PUT /api/posts/:id — update post
router.put("/:id", (req, res) => {
  const id = parseInt(req.params.id, 10);
  const index = posts.findIndex((p) => p.id === id);
  if (index === -1) {
    return res.status(404).json({ error: "Post not found" });
  }
  const { title, body, userId } = req.body;
  if (!title && !body && !userId) {
    return res.status(400).json({ error: "At least one of title, body, or userId is required" });
  }
  if (userId !== undefined) {
    const userExists = users.find((u) => u.id === parseInt(userId, 10));
    if (!userExists) {
      return res.status(404).json({ error: "User not found" });
    }
    posts[index].userId = parseInt(userId, 10);
  }
  if (title) posts[index].title = title;
  if (body) posts[index].body = body;
  res.json(posts[index]);
});

// DELETE /api/posts/:id — delete post
router.delete("/:id", (req, res) => {
  const id = parseInt(req.params.id, 10);
  const index = posts.findIndex((p) => p.id === id);
  if (index === -1) {
    return res.status(404).json({ error: "Post not found" });
  }
  posts.splice(index, 1);
  res.status(204).send();
});

export default router;
