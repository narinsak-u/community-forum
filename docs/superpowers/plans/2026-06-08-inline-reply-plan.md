# Inline Comment Reply Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Wire up REPLY buttons on comments to open an inline reply box pre-filled with `@author`, submitting as nested replies via existing `parentId` API.

**Architecture:** Single-file change in ThreadDetailClient.tsx. A `replyTarget` state tracks which comment is being replied to. When set, an inline textarea renders below that comment. Submit calls `createComment.mutate({ content, parentId })`. No backend changes.

**Tech Stack:** Next.js 15, React 19, TanStack Query, Zustand auth store, sonner toast

---

## File Structure

| File | Status | Responsibility |
|---|---|---|
| `frontend/src/app/(main)/thread/[slug]/ThreadDetailClient.tsx` | Modify | Add replyTarget state, wire REPLY buttons, inline reply box |

### Task 1: Add inline comment reply

**Files:**
- Modify: `frontend/src/app/(main)/thread/[slug]/ThreadDetailClient.tsx`

- [ ] **Step 1: Add replyTarget state**

After the existing `const [replyContent, setReplyContent] = useState("");` (line 34), add:

```tsx
const [replyTarget, setReplyTarget] = useState<{ id: number; author: string } | null>(null);
```

- [ ] **Step 2: Wire up the REPLY buttons on comments**

Find the existing static "REPLY" button (around line 316):

```tsx
<button className="hover:text-primary">REPLY</button>
```

Replace with:

```tsx
<button
  className="hover:text-primary"
  onClick={() => {
    if (!requireAuth({ toast: true, description: "Authenticate to reply." })) return;
    setReplyTarget({
      id: comment.id,
      author: comment.author?.username || "@unknown",
    });
    setReplyContent(`@${(comment.author?.username || "@unknown").replace("@", "")} `);
  }}
>
  REPLY
</button>
```

- [ ] **Step 3: Add inline reply box below each comment**

After the `</p>` that shows `comment.content` (around line 314) and before the `<div className="flex gap-4...">` that contains the REPLY button, add:

```tsx
              {replyTarget?.id === comment.id && user && (
                <div className="ml-6 mt-4 pl-5 border-l border-primary/30 space-y-3">
                  <div className="flex items-center gap-3">
                    <div className="h-7 w-7 rounded-sm bg-gradient-signal grid place-items-center text-[9px] font-bold text-primary-foreground">
                      {user.username.replace("@", "").slice(0, 2).toUpperCase()}
                    </div>
                    <span className="text-xs uppercase tracking-[0.18em] text-muted-foreground">
                      REPLYING AS {user.username.toUpperCase()}
                    </span>
                  </div>
                  <Textarea
                    value={replyContent}
                    onChange={(e) => setReplyContent(e.target.value)}
                    placeholder={`@${replyTarget.author.replace("@", "")} `}
                    className="bg-terminal border-border min-h-[80px] font-mono text-sm focus-visible:ring-primary/40"
                  />
                  <div className="flex justify-end gap-2">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setReplyTarget(null)}
                      className="text-muted-foreground uppercase text-[10px]"
                    >
                      Cancel
                    </Button>
                    <Button
                      size="sm"
                      onClick={() => {
                        if (!requireAuth({ toast: true, description: "Authenticate to transmit a reply." })) return;
                        createComment.mutate(
                          { content: replyContent, parentId: replyTarget.id },
                          {
                            onSuccess: () => {
                              setReplyContent("");
                              setReplyTarget(null);
                              toast.success("REPLY_TRANSMITTED");
                            },
                            onError: (err) =>
                              toast.error("TRANSMIT_FAILED", { description: err.message }),
                          },
                        );
                      }}
                      disabled={createComment.isPending || !replyContent.trim()}
                      className="bg-gradient-signal hover:opacity-90 text-primary-foreground font-bold uppercase tracking-[0.18em] text-[10px] rounded-sm"
                    >
                      {createComment.isPending ? "TRANSMITTING..." : "TRANSMIT"}
                    </Button>
                  </div>
                </div>
              )}
```

- [ ] **Step 4: Build to verify**

Run: `cd frontend && npx next build`
Expected: No errors, all routes compile

- [ ] **Step 5: Commit**

```bash
git add "frontend/src/app/(main)/thread/\[slug\]/ThreadDetailClient.tsx"
git commit -m "feat: add inline comment reply with @author prefill"
```
