# Agent Teams — Master Reference Guide

> A grounded, practical reference for designing and operating Claude Code agent teams and subagent definitions.
>
> **Primary sources** (always re-check these — features evolve):
> - https://code.claude.com/docs/en/agent-teams
> - https://code.claude.com/docs/en/sub-agents
> - https://code.claude.com/docs/en/hooks
> - https://code.claude.com/docs/en/settings

---

## 0. Read this first — two distinct concepts

Claude Code uses two terms that sound similar and overlap on disk, but solve different problems. Get this distinction right and the rest of the guide makes sense.

| Concept | What it is | Where it lives | Communication model |
|---|---|---|---|
| **Subagent** | A specialized worker spawned **inside one session**. Returns a single summary to the caller. | `.claude/agents/*.md`, `~/.claude/agents/*.md`, plugin, or `--agents` CLI. **GA.** | Caller → subagent → caller. No peer-to-peer. |
| **Agent team** | Multiple **separate Claude Code sessions** (lead + teammates) coordinating via a shared task list and direct messaging. | Runtime state under `~/.claude/teams/{team-name}/` and `~/.claude/tasks/{team-name}/`. **Experimental** — gated by `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`. | Teammate ↔ teammate ↔ lead, plus shared task list. |

**The key crossover:** a subagent definition file (with YAML frontmatter) can be reused as an agent team teammate. Define the role once, use it both ways. (Caveat: when a definition runs as a teammate, `skills` and `mcpServers` from frontmatter are **not** applied — teammates load skills/MCP from project + user settings instead.)

**Decision tree:**
```
Need parallel work?
├── No  → main conversation, possibly with skills
└── Yes → Do workers need to talk to each other?
         ├── No  → Subagents (cheaper, simpler, GA)
         └── Yes → Agent team (experimental, higher token cost, more capability)
```

---

## 1. The frontmatter spec (canonical)

Every agent file is Markdown with YAML frontmatter. Only `name` and `description` are required. Source of truth: the sub-agents docs.

| Field | Required | Purpose | Notes |
|---|---|---|---|
| `name` | yes | Unique identifier (lowercase + hyphens). | Used for @-mention, `--agent`, and `Agent(name)` permission rules. |
| `description` | yes | When Claude should delegate to this agent. | **Most important field.** This is what the orchestrator reads to decide routing. Be specific about triggers. Phrases like "use proactively" encourage delegation. |
| `tools` | no | Allowlist of tools. Inherits all if omitted. | Comma-separated. Use `Agent(worker, researcher)` to allowlist which subagent types can be spawned (only meaningful when this agent is the main thread via `--agent`). |
| `disallowedTools` | no | Denylist; applied before `tools`. | Use to inherit everything except a few (e.g. `Write, Edit`). |
| `model` | no | `sonnet`, `opus`, `haiku`, full ID, or `inherit`. | Default `inherit`. Resolution order: env var `CLAUDE_CODE_SUBAGENT_MODEL` > per-invocation > frontmatter > parent. |
| `permissionMode` | no | `default`, `acceptEdits`, `auto`, `dontAsk`, `bypassPermissions`, `plan`. | Parent's `bypassPermissions` / `acceptEdits` / `auto` override the child's. |
| `maxTurns` | no | Hard cap on agentic turns. | Useful for read-only scanners and bounded research. |
| `skills` | no | Skill names to **preload into context** at startup. | Full skill body is injected. Subagents do **not** inherit skills from parent. **Not applied when used as a teammate.** |
| `mcpServers` | no | List of inline server defs or string references to existing servers. | Inline servers connect at start, disconnect at end. **Not applied when used as a teammate.** |
| `hooks` | no | Lifecycle hooks scoped to this agent. | `PreToolUse`, `PostToolUse`, `Stop` (auto-converted to `SubagentStop`). |
| `memory` | no | `user`, `project`, or `local`. | Enables a persistent memory directory that survives across conversations. Auto-enables Read/Write/Edit. |
| `background` | no | `true` to always run as a background task. | Default `false`. Permissions are pre-approved at spawn; unapproved tools auto-deny. |
| `effort` | no | `low`, `medium`, `high`, `xhigh`, `max`. | Available levels depend on model. |
| `isolation` | no | `worktree` to run in a temp git worktree. | Auto-cleaned if no changes; otherwise the path/branch is returned. |
| `color` | no | `red`, `blue`, `green`, `yellow`, `purple`, `orange`, `pink`, `cyan`. | Display only. |
| `initialPrompt` | no | Auto-submitted as the first user turn when this agent runs as the **main session** agent (`--agent`). | Commands and skills are processed. |

**The body** of the file becomes the agent's system prompt. Subagents receive only this prompt + minimal env details — **not** the full Claude Code system prompt. (Forks are the exception: see §6.)

---

## 2. Scopes and precedence

Agent definitions can come from five places. When names collide, the higher-priority source wins.

| Priority | Location | Scope |
|---|---|---|
| 1 (highest) | Managed settings (`.claude/agents/` inside the managed settings dir) | Org-wide |
| 2 | `--agents '{...}'` JSON on the CLI | One session |
| 3 | `<repo>/.claude/agents/` | This project |
| 4 | `~/.claude/agents/` | All your projects |
| 5 (lowest) | Plugin's `agents/` directory | Where the plugin is enabled |

Project-scoped agents are discovered by walking up from the cwd. Directories added with `--add-dir` are **not** scanned for agents. Plugin agents cannot use `hooks`, `mcpServers`, or `permissionMode` (security restriction).

**Where to put new agents:**
- Project-specific reviewer / engineer roles → `.claude/agents/` (commit them — collaborative).
- Personal patterns you reuse across projects → `~/.claude/agents/`.
- Experimental / one-off → `--agents` JSON.
- Distribute to others → plugin.

---

## 3. Anatomy of a high-quality agent definition

The frontmatter routes the work; the body shapes the work. A great agent body is **scoped, opinionated, and falsifiable**. Every section answers a question the agent (or a colleague reading it later) will actually ask.

### Recommended skeleton

```markdown
---
name: <kebab-case-name>
description: <when-to-use-this; include trigger phrases>
tools: <allowlist>            # optional
model: <alias|inherit>        # optional
---

# Agent: <Title-Case Name>

## Role
One sentence. Who this agent is in the team.

## Mission
2–3 sentences. What success looks like.

## Responsibilities
Bulleted, scoped, concrete. Group related concerns.

## Non-Responsibilities
What this agent must NOT do. This is where role drift dies.

## Source of Truth
What artifacts this agent treats as authoritative (files, schemas, plan docs).

## Operating Rules
Numbered, terse, enforceable. Each rule should be checkable.

## Expected Inputs
What the agent expects to receive when invoked.

## Expected Outputs
Concrete deliverables. Files written, comments posted, summaries returned.

## Quality Bar
Pass/fail criteria — not aspirations.

## Escalation Rules
Who to escalate to and when. Includes "escalate to user" gates.

## Testing Expectations
How this agent's work is verified.

## Definition of Done
The unambiguous checklist. The agent should be able to self-audit against this.
```

### What separates "fine" from "load-bearing"

| Good signal | Bad signal |
|---|---|
| `description` lists concrete triggers ("use after Go code changes touching `internal/repository/`") | `description` reads like marketing ("an excellent code reviewer") |
| Non-responsibilities are explicit | Only responsibilities are listed (role drift is inevitable) |
| Operating rules cite *why* | Rules are bare imperatives the agent can't argue with |
| Definition of Done is checkable from outside | DoD is "high quality" |
| Escalation paths are named (PM / Reviewer / user) | Agent has no escape hatch — it will guess |

---

## 4. Patterns that work

### 4.1 The PM / IC / Reviewer triad
Three roles, clear seams. Used in this repo's `.claude/agents/`.

- **PM** — plans, sequences, owns DoD. No code.
- **IC(s)** — implements the plan (`backend-engineer`, `frontend-engineer`).
- **Reviewer** — independent quality gate, blocking vs. non-blocking comments. Doesn't implement.

Strength: each role has a distinct success criterion that doesn't overlap with the others, so escalations are unambiguous.
Use when: the work spans multiple components and you need delivery discipline, not just throughput.

### 4.2 Parallel specialists (the "different lenses" pattern)
Multiple agents review or research the **same target** through different filters.

Example: PR review with three reviewers — security, performance, test coverage. Each is told to *not* drift into the others' lanes.

Strength: avoids the single-reviewer drift toward one issue type at a time.
Use when: thorough coverage of a fixed surface matters more than coordination.

### 4.3 Adversarial debate (competing hypotheses)
N agents each take a different theory; their explicit job is to disprove the others.

Example from official docs: "Spawn 5 teammates to investigate different hypotheses. Have them talk to each other to try to disprove each other's theories, like a scientific debate."

Strength: counters the anchoring bias of sequential investigation.
Use when: root cause is unclear and the cost of premature commitment is high.

### 4.4 Cross-layer coordination
One teammate per layer (frontend / backend / tests / migrations) on a single feature.

Strength: each owns its files; conflicts are minimized; contracts must be explicit.
Use when: a feature genuinely spans layers and the layers are file-disjoint.

### 4.5 Sequential pipeline (chained subagents)
Subagent A → main → Subagent B → main → ...

Example: "Use the code-reviewer subagent to find performance issues, then use the optimizer subagent to fix them."

Strength: each stage produces a clean handoff; main agent stays in control.
Use when: stages are dependent, output of one is input to the next.

### 4.6 Forked exploration (experimental)
Fork the current conversation to try a side task or alternative approach without losing the prompt cache.

Use when: a named subagent would need too much background to be useful, or you want to try several approaches in parallel from the same starting point.

---

## 5. Sizing, sequencing, and steering

From the official best practices, calibrated:

- **Team size:** start with 3–5 teammates. Three focused often beats five scattered.
- **Tasks per teammate:** aim for 5–6. Fewer = idle teammates; more = lost context.
- **Task size:** self-contained units producing a clear deliverable (a function, a test file, a review report). Too small → coordination overhead exceeds benefit. Too large → no check-ins, wasted runs.
- **Avoid file conflicts:** assign disjoint file sets per teammate. Same-file edits = overwrites.
- **Wait for teammates:** the lead sometimes implements instead of waiting. If you see this, tell the lead to wait.
- **Start with research and review:** the safest first agent-team task. No write conflicts, clear value.
- **Monitor and steer:** check teammate panes, redirect dead ends, synthesize as you go.
- **Plan approval for risky changes:** spawn teammates with "require plan approval before they make any changes." The teammate stays in plan mode until the lead approves.

**Token cost is real.** Each teammate is a full Claude Code instance with its own context window. Cost scales linearly with teammate count. Use a single session or subagents for routine tasks.

---

## 6. Forks vs. named subagents (when each is right)

|  | Fork | Named subagent |
|---|---|---|
| Context | Full conversation history inherited | Fresh context with the spawn prompt |
| System prompt + tools | Same as main session | From the agent definition |
| Model | Same as main session | From definition's `model` field |
| Permissions | Prompts surface in your terminal | Pre-approved at launch, then auto-deny |
| Prompt cache | Shared with main session (cheap first request) | Separate cache |
| Can spawn more forks? | No | N/A |

Heuristic: **fork** when the situation is too rich to brief; **subagent** when the role is reusable and the prompt is self-contained.

---

## 7. Hooks: enforce rules without trusting the agent

Hooks are the deterministic backstop. The harness runs them, not Claude.

### Inside the agent definition (`hooks:` frontmatter)
- `PreToolUse` (matcher: tool name) — gate or transform tool calls. Exit 2 to block.
- `PostToolUse` (matcher: tool name) — run linters, formatters, validators after edits.
- `Stop` — runs when the agent finishes. Auto-converts to `SubagentStop` when invoked as a subagent.

### In `settings.json` (project- or user-level)
- `SubagentStart` / `SubagentStop` (matcher: agent type) — wire setup/teardown around specific roles.
- `TaskCreated` / `TaskCompleted` (agent teams only) — gate task lifecycle. Exit 2 to block + send feedback.
- `TeammateIdle` (agent teams only) — exit 2 to keep a teammate working.

**Use hooks for things the agent cannot be trusted to enforce on itself**: lint clean before a commit, no `INSERT/UPDATE/DELETE` from a read-only DB agent, every PR has a reviewer assigned, etc.

---

## 8. Building a new agent team — checklist

Run this every time you stand up a team. The order matters.

### Discovery
- [ ] What is the *one* outcome that justifies a team? (If you can't name it, don't spawn one.)
- [ ] Can the work be split into 3–5 file-disjoint pieces? If not, prefer a single session or subagents.
- [ ] Are inter-agent communication and shared state genuinely needed? If not, subagents are cheaper.
- [ ] What is the smallest team that can ship this? Default: 3.

### Roles
- [ ] Each role has a one-sentence mission and a non-overlapping success criterion.
- [ ] Each role has explicit non-responsibilities.
- [ ] Escalation paths are named: who escalates to whom, who escalates to the user.
- [ ] Decide which roles are reusable (define as agent files) vs. one-off (use spawn prompts).

### Files
- [ ] Each reusable role gets a file under `.claude/agents/` (or `~/.claude/agents/` for personal).
- [ ] Frontmatter: `name`, `description` (rich, with triggers), and a tool allowlist if the role is read-only or scoped.
- [ ] Body follows the §3 skeleton.
- [ ] If destructive tools aren't needed, restrict them. (E.g., reviewers should not have `Edit` / `Write`.)

### Contracts
- [ ] Inter-role contracts exist *before* implementation: API shapes, file ownership, message formats.
- [ ] Source-of-truth artifacts are listed for each role.
- [ ] PM owns the plan doc; engineers own implementation; reviewer owns approval.

### Operations
- [ ] Local-dev parity is defined as a deliverable, not an afterthought.
- [ ] Hooks are configured for the rules you can't trust the agents to follow (lint, secret scanning, etc.).
- [ ] Permission settings pre-approve common ops to reduce interrupts (especially with teammates — their prompts bubble up to the lead).

### Verification
- [ ] DoD for each task is checkable by a non-author.
- [ ] At least one E2E smoke run is defined per milestone.
- [ ] The reviewer role explicitly verifies acceptance criteria, not just code style.

### Shutdown
- [ ] You know how to clean up: tell the **lead** to clean up the team; never have a teammate do it.
- [ ] If using tmux split panes, you know how to find and kill orphaned sessions (`tmux ls` / `tmux kill-session`).

---

## 9. Common failure modes and how to avoid them

| Failure | Symptom | Fix |
|---|---|---|
| Vague `description` field | Wrong agent gets delegated to, or none does | Rewrite with concrete triggers and "use proactively" cues |
| Role drift | Reviewer starts coding; PM starts implementing | Explicit non-responsibilities + escalation rules |
| File conflicts in agent teams | Overwrites, lost work | Assign disjoint file sets at planning time |
| Lead implements instead of waiting | Slow, sequential outcome | Tell the lead to wait; size tasks to be claimable |
| Lead shuts down too early | "Done" while tasks remain | Tell it to keep going; check the shared task list |
| Token cost balloons | High burn, marginal value | Drop to subagents or a single session |
| Stuck task (lag) | Dependent tasks blocked | Manually mark done or nudge the teammate |
| Orphaned tmux sessions | Old panes lingering | `tmux ls` + `tmux kill-session -t <name>` |
| Resume breaks teammates | Lead messages dead teammates after `/resume` | Spawn fresh teammates; in-process teammates don't survive resume (known limitation) |
| Plugin agent ignores `hooks`/`mcpServers`/`permissionMode` | Configured but not enforced | Copy the file into `.claude/agents/` (these fields are stripped from plugin agents for security) |

---

## 10. Limitations to design around (as of these docs)

- **Experimental flag required** for agent teams: `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`.
- **Version floor:** Claude Code v2.1.32+ for agent teams; v2.1.117+ for forked subagents (`CLAUDE_CODE_FORK_SUBAGENT=1`).
- **No nested teams** — teammates can't spawn their own teams.
- **One team per session** — clean up before starting another.
- **Lead is fixed** — can't promote a teammate.
- **Permissions set at spawn** — can't pre-set per-teammate modes; change after spawning.
- **Subagents cannot spawn subagents** — chain through main, or use forks/skills.
- **In-process teammates don't survive `/resume` or `/rewind`.**
- **Split-pane mode requires tmux or iTerm2** — not available in VS Code's integrated terminal, Windows Terminal, or Ghostty.
- **`--add-dir` directories aren't scanned for agent definitions.**

---

## 11. Templates

### 11.1 Read-only reviewer (subagent)
```markdown
---
name: <something>-reviewer
description: <When to use proactively. List concrete triggers — file globs, change types, PR contexts.>
tools: Read, Grep, Glob, Bash
model: inherit
---

# Agent: <Name>

## Role
Independent review for <domain>.

## Mission
Catch <class of defects> before merge. Block on real issues; suggest on style.

## Responsibilities
- ...

## Non-Responsibilities
- Does not implement fixes unless explicitly asked.

## Operating Rules
1. Read the full diff, not snippets.
2. Two comment classes, labeled BLOCKING vs NON-BLOCKING.
3. Evidence over opinion.

## Definition of Done
- Verdict issued: APPROVED / APPROVED WITH NITS / REQUEST CHANGES / BLOCKED.
- Acceptance criteria explicitly verified.
```

### 11.2 Bounded researcher (subagent)
```markdown
---
name: <topic>-researcher
description: Investigates <topic> and returns a structured summary. Use when verbose output would flood main context.
tools: Read, Grep, Glob, WebFetch, WebSearch
maxTurns: 20
model: haiku
---

# Agent: <Name>

## Role
Read-only investigator.

## Operating Rules
1. Search broadly, read narrowly.
2. Return findings as: facts, sources, open questions.
3. Do not modify any files.
```

### 11.3 Implementer (subagent or teammate)
```markdown
---
name: <component>-engineer
description: Implements <component> changes. Owns <file paths>. Use after PM brief is ready.
# tools: omitted -> inherits all
model: inherit
---

# Agent: <Name>

## Role
The implementer for <component>.

## Operating Rules
1. TDD where practical.
2. Real dependencies for integration tests; no mock-only "integration" tests.
3. Verification before completion: build + test + manual smoke.
4. Never edit files outside <component> without coordination.

## Definition of Done
- Lint, build, tests green (cite evidence).
- Reviewer approved.
- Acceptance criteria met.
```

### 11.4 Spawn prompt for an agent team (one-shot, no file)
```text
Create an agent team to <outcome>. Spawn <N> teammates:
- <name1>: <one-line mission, files owned, success criterion>
- <name2>: <...>
- <name3>: <...>

Coordination:
- <name1> publishes the contract before <name2> consumes it.
- All teammates wait for the lead to approve plans before writing code.

Hard rules:
- No teammate edits files outside its assigned set.
- Reviewer must sign off before any task is marked done.
- Use Sonnet for each teammate.
```

---

## 12. Quick decision aids

### Should I create a new agent definition?
- I keep writing the same instructions for the same kind of work → **yes**, file under the right scope.
- This is a one-off task → **no**, use the main conversation or `--agents` JSON.
- I want to share with my team → **yes**, project scope (`.claude/agents/`), commit it.
- I want it everywhere on my machine → **yes**, user scope (`~/.claude/agents/`).

### Should I spawn an agent team?
- Workers need to talk to each other → **yes**.
- Workers don't need to talk; only results matter → **no**, use subagents.
- Iterative back-and-forth with the user → **no**, main conversation.
- Token budget is tight → **no**, use a single session.

### Should I use `tools` allowlist or `disallowedTools`?
- "This agent is read-only" → **allowlist**: `Read, Grep, Glob, Bash`.
- "This agent does almost everything except writing files" → **denylist**: `Write, Edit`.
- "Don't touch MCP at all" → allowlist (MCP tools are inherited by default).

### Should I set `model`?
- Bulk read-only / search → `haiku`.
- Implementation / review of non-trivial changes → `inherit` (or `sonnet` explicitly).
- Architectural reasoning, complex planning → `opus`.

---

## 13. Pointers (always re-check before relying on details)

- Agent teams: https://code.claude.com/docs/en/agent-teams
- Subagents: https://code.claude.com/docs/en/sub-agents
- Hooks: https://code.claude.com/docs/en/hooks
- Settings: https://code.claude.com/docs/en/settings
- Permissions: https://code.claude.com/docs/en/permissions
- Plugins: https://code.claude.com/docs/en/plugins
- Skills: https://code.claude.com/docs/en/skills
- MCP: https://code.claude.com/docs/en/mcp
- Costs (incl. agent team token costs): https://code.claude.com/docs/en/costs

> Features evolve. When something in this guide disagrees with the current docs, **the docs win** — update this guide.
