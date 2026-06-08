# Inline Comment Reply Design

**Date:** 2026-06-08
**Status:** Approved

## Overview

Fix the broken reply flow on comments. Clicking "REPLY" on a comment opens an inline textarea below it, pre-filled with `@author`, and submits as a nested reply via the existing `parentId` API. No backend changes.

## Frontend design

Single-file change in `ThreadDetailClient.tsx`:

### State

```ts
const [replyTarget, setReplyTarget] = useState<{ id: number; author: string } | null>(null);
```

### REPLY button

Each comment's static "REPLY" button becomes:

```tsx
<button className="hover:text-primary" onClick={() => setReplyTarget({ id: comment.id, author: comment.author?.username || "@unknown" })}>
  REPLY
</button>
```

### Inline reply box

Rendered below the comment content and above its nested replies, only when `replyTarget.id === comment.id`:

```tsx
{replyTarget?.id === comment.id && (
  <div className="ml-6 mt-4 pl-5 border-l border-primary/30 space-y-3">
    <div className="flex items-center gap-3">
      <div className="h-7 w-7 rounded-sm bg-secondary border border-border grid place-items-center text-[9px] font-bold text-primary">
        {user?.username?.replace("@", "").slice(0, 2).toUpperCase()}
      </div>
      <span className="text-xs uppercase tracking-[0.18em] text-muted-foreground">
        REPLYING AS {user?.username?.toUpperCase()}
      </span>
    </div>
    <Textarea
      value={replyContent}
      onChange={(e) => setReplyContent(e.target.value)}
      placeholder={`@${replyTarget.author.replace("@", "")} `}
      className="bg-terminal border-border min-h-[80px] font-mono text-sm focus-visible:ring-primary/40"
    />
    <div className="flex justify-end gap-2">
      <Button variant="ghost" size="sm" onClick={() => setReplyTarget(null)} className="text-muted-foreground uppercase text-[10px]">
        Cancel
      </Button>
      <Button
        size="sm"
        onClick={() => {
          if (!requireAuth({ toast: true, description: "Authenticate to transmit a reply." })) return;
          createComment.mutate(
            { content: replyContent, parentId: comment.id },
            {
              onSuccess: () => {
                setReplyContent("");
                setReplyTarget(null);
                toast.success("REPLY_TRANSMITTED");
              },
              onError: (err) => toast.error("TRANSMIT_FAILED", { description: err.message }),
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

The inline box is shown only when user is authenticated (`user` truthy). If not authenticated, clicking REPLY triggers the auth toast from the existing guard in `handleSubmitReply` — but since the inline box shouldn't appear for unauthenticated users, guard the REPLY button's click handler with `requireAuth({ toast: true })` before setting `replyTarget`.

### REPLY button guard

```tsx
<button
  className="hover:text-primary"
  onClick={() => {
    if (!requireAuth({ toast: true, description: "Authenticate to reply." })) return;
    setReplyTarget({ id: comment.id, author: comment.author?.username || "@unknown" });
  }}
>
  REPLY
</button>
```

## Data flow

```
Click REPLY → auth guard → setReplyTarget → inline textarea renders
  → User types → Submit → createComment.mutate({ content, parentId })
    → POST /api/v1/threads/:slug/comments { content, parentId }
      → Backend CommentService.Create → validates parent, creates with ParentID
        → QueryClient invalidates ["thread", slug] → thread refetch → replies re-render
```

## Files changed

| File | Change |
|---|---|
| `frontend/src/app/(main)/thread/[slug]/ThreadDetailClient.tsx` | Add replyTarget state, wire REPLY buttons, inline reply box |
