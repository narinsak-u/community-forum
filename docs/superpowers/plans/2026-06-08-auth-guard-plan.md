# Auth Guard Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add auth guard layer — redirect protected pages to `/login`, toast for protected actions, auth-aware sidebar.

**Architecture:** A single `useRequireAuth` hook handles both redirect (pages) and toast (actions) modes. Sidebar conditionally renders user info + logout or login link based on `useAuthStore.isAuthenticated`.

**Tech Stack:** Zustand auth-store, TanStack Query mutations, sonner toast, Next.js router

---

## File Structure

| File | Status | Responsibility |
|---|---|---|
| `frontend/src/hooks/use-require-auth.ts` | Create | `useRequireAuth` hook — redirect + toast |
| `frontend/src/app/(main)/profile/page.tsx` | Modify | Add auth guard useEffect |
| `frontend/src/app/(main)/settings/page.tsx` | Modify | Add auth guard useEffect |
| `frontend/src/app/(main)/create/page.tsx` | Modify | Add auth guard useEffect |
| `frontend/src/app/(main)/thread/[slug]/ThreadDetailClient.tsx` | Modify | Guard vote + reply handlers |
| `frontend/src/components/forge/Sidebar.tsx` | Modify | Auth-aware footer |

### Task 1: Create `useRequireAuth` hook

**Files:**
- Create: `frontend/src/hooks/use-require-auth.ts`

- [ ] **Step 1: Write the hook**

```ts
import { useCallback } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { useAuthStore } from "@/stores/auth-store";

interface RequireAuthOptions {
  redirect?: string;
  toast?: boolean;
  description?: string;
}

export function useRequireAuth() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const router = useRouter();

  const requireAuth = useCallback(
    (opts?: RequireAuthOptions): boolean => {
      if (isAuthenticated) return true;

      if (opts?.toast) {
        toast("SIGN_IN_REQUIRED", {
          description: opts.description || "Authenticate to continue.",
          action: {
            label: "LOGIN",
            onClick: () => router.push("/login"),
          },
        });
      } else {
        const redirectPath = opts?.redirect || "/";
        router.push(`/login?redirect=${encodeURIComponent(redirectPath)}`);
      }

      return false;
    },
    [isAuthenticated, router],
  );

  return { requireAuth };
}
```

- [ ] **Step 2: Verify file compiles**

Run: `cd frontend && npx next build`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add frontend/src/hooks/use-require-auth.ts
git commit -m "feat: add useRequireAuth hook"
```

### Task 2: Guard profile, settings, create pages

**Files:**
- Modify: `frontend/src/app/(main)/profile/page.tsx`
- Modify: `frontend/src/app/(main)/settings/page.tsx`
- Modify: `frontend/src/app/(main)/create/page.tsx`

- [ ] **Step 1: Guard profile page**

```tsx
"use client";

import { useEffect } from "react";
import { useRequireAuth } from "@/hooks/use-require-auth";
import { ProfileContent } from "./ProfileContent";

export default function MyProfilePage() {
  const { requireAuth } = useRequireAuth();

  useEffect(() => {
    requireAuth({ redirect: "/profile" });
  }, [requireAuth]);

  return <ProfileContent />;
}
```

- [ ] **Step 2: Guard settings page**

Add import and useEffect. The page is already `"use client"` and has `useAuthStore` — add `useRequireAuth` import and the same useEffect pattern. After the first `useState` line, add:

```tsx
import { useEffect } from "react";
import { useRequireAuth } from "@/hooks/use-require-auth";
```

Inside the component body, after the `useState` lines:

```tsx
const { requireAuth } = useRequireAuth();

useEffect(() => {
  requireAuth({ redirect: "/settings" });
}, [requireAuth]);
```

- [ ] **Step 3: Guard create page**

`frontend/src/app/(main)/create/page.tsx` — read the current file first, then add `useEffect` + `useRequireAuth` guard. If it's already `"use client"`, add the same pattern as settings.

- [ ] **Step 4: Build to verify**

Run: `cd frontend && npx next build`
Expected: No errors

- [ ] **Step 5: Commit**

```bash
git add frontend/src/app/\(main\)/profile/page.tsx frontend/src/app/\(main\)/settings/page.tsx frontend/src/app/\(main\)/create/page.tsx
git commit -m "feat: guard profile, settings, create pages with auth redirect"
```

### Task 3: Guard vote and reply actions in ThreadDetailClient

**Files:**
- Modify: `frontend/src/app/(main)/thread/[slug]/ThreadDetailClient.tsx`

- [ ] **Step 1: Add import**

After existing imports:

```tsx
import { useRequireAuth } from "@/hooks/use-require-auth";
```

- [ ] **Step 2: Add hook at top of component**

After `const user = useAuthStore((s) => s.user);`:

```tsx
const { requireAuth } = useRequireAuth();
```

- [ ] **Step 3: Guard the vote handler**

Find `handleVote` function. Add auth check at the top:

```tsx
const handleVote = (value: number) => {
  if (!requireAuth({ toast: true, description: "Authenticate to cast a signal." })) return;
  voteThread.mutate({ value }, { ... });
};
```

- [ ] **Step 4: Guard the reply/comment handler**

Find the submit handler (the one that calls `createComment.mutate`). Add auth check before mutate:

```tsx
const handleSubmitReply = () => {
  if (!replyContent.trim()) return;
  if (!requireAuth({ toast: true, description: "Authenticate to transmit a reply." })) return;
  createComment.mutate({ content: replyContent }, { ... });
};
```

- [ ] **Step 5: Build to verify**

Run: `cd frontend && npx next build`
Expected: No errors

- [ ] **Step 6: Commit**

```bash
git add "frontend/src/app/\(main\)/thread/\[slug\]/ThreadDetailClient.tsx"
git commit -m "feat: guard vote and reply actions with auth toast"
```

### Task 4: Auth-aware sidebar footer

**Files:**
- Modify: `frontend/src/components/forge/Sidebar.tsx`

- [ ] **Step 1: Update imports**

Replace the current footerItems and add new imports:

```tsx
import { LayoutGrid, MessagesSquare, Shapes, BookOpen, Plus, ChevronLeft, ChevronRight, LogIn, LogOut } from "lucide-react";
```

Remove `HelpCircle`, `Archive` from the imports (they're no longer in footerItems).

- [ ] **Step 2: Remove static footerItems**

Delete the `footerItems` array — it will be replaced with auth-aware content.

- [ ] **Step 3: Add auth imports**

```tsx
import { useAuthStore } from "@/stores/auth-store";
import { useSignout } from "@/hooks/use-auth";
```

- [ ] **Step 4: Add auth state and signout to component**

After the `toggle` function:

```tsx
const user = useAuthStore((s) => s.user);
const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
const signout = useSignout();
```

- [ ] **Step 5: Add router and signout handler**

Add import:

```tsx
import { useRouter } from "next/navigation";
```

Add at top of component body:

```tsx
const router = useRouter();
```

Add this signout handler function after the `toggle` function:

```tsx
const handleSignout = () => {
  signout.mutate(undefined, {
    onSuccess: () => router.push("/"),
  });
};
```

- [ ] **Step 6: Replace footer div**

Replace the existing `<div className="p-3 border-t border-border/60 space-y-1">` block with:

```tsx
      <div className="p-3 border-t border-border/60 space-y-1">
        {isAuthenticated && user ? (
          <>
            <ActiveLink
              href="/profile"
              className={cn(
                "flex items-center gap-3 px-3 py-2 text-sm text-sidebar-foreground hover:text-foreground rounded-sm transition-colors",
                collapsed && "justify-center px-0",
              )}
            >
              <div className="h-6 w-6 rounded-sm bg-secondary border border-border grid place-items-center text-[10px] font-bold text-primary shrink-0">
                {user.username.slice(0, 2).toUpperCase()}
              </div>
              {!collapsed && (
                <span className="truncate">{user.username}</span>
              )}
            </ActiveLink>
            {collapsed ? (
              <Tooltip delayDuration={300}>
                <TooltipTrigger asChild>
                  <button
                    onClick={handleSignout}
                    className="flex items-center justify-center px-3 py-2 w-full text-sm text-muted-foreground hover:text-destructive rounded-sm transition-colors"
                  >
                    <LogOut className="h-4 w-4 shrink-0" />
                  </button>
                </TooltipTrigger>
                <TooltipContent side="right" className="text-xs">Log Out</TooltipContent>
              </Tooltip>
            ) : (
              <button
                onClick={handleSignout}
                className="flex items-center gap-3 px-3 py-2 w-full text-sm text-muted-foreground hover:text-destructive rounded-sm transition-colors"
              >
                <LogOut className="h-4 w-4 shrink-0" />
                <span>Log Out</span>
              </button>
            )}
          </>
        ) : (
          <ActiveLink
            href="/login?redirect=/"
            className={cn(
              "flex items-center gap-3 px-3 py-2 text-sm text-muted-foreground hover:text-foreground rounded-sm transition-colors",
              collapsed && "justify-center px-0",
            )}
          >
            <LogIn className="h-4 w-4 shrink-0" />
            {!collapsed && <span>Log In</span>}
          </ActiveLink>
        )}
      </div>
```

- [ ] **Step 7: Build to verify**

Run: `cd frontend && npx next build`
Expected: No errors

- [ ] **Step 8: Commit**

```bash
git add frontend/src/components/forge/Sidebar.tsx
git commit -m "feat: add auth-aware sidebar footer with login/logout"
```

### Task 5: Final verification

- [ ] **Step 1: Full build**

Run: `cd frontend && npx next build`
Expected: All routes compile successfully

- [ ] **Step 2: Run lint**

Run: `cd frontend && npx next lint`
Expected: No errors

- [ ] **Step 3: Final commit**

```bash
git add -A
git commit -m "feat: complete auth guard implementation"
```
