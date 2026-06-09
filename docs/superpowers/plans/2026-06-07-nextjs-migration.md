# Next.js Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Migrate the community forum frontend from Vite + React 18 + React Router to Next.js 15 (App Router) + React 19 with hybrid server/client rendering.

**Architecture:** Clean replacement — replace Vite project structure with Next.js. Keep `src/` directory with `src/app/` for App Router. All existing components, hooks, stores, and lib utilities stay. Static assets move to `public/`. Server components fetch initial data from the Go API. Client components handle interactivity via TanStack Query and Zustand.

**Tech Stack:** Next.js 15 (App Router), React 19, TypeScript, Tailwind CSS, shadcn/ui, TanStack Query, Zustand, Vitest

---

### Task 1: Scaffold Next.js project

**Files:**
- Modify: `frontend/package.json`
- Modify: `frontend/tsconfig.json`
- Create: `frontend/next.config.ts`
- Modify: `frontend/.gitignore`

- [ ] **Step 1: Update package.json**

Replace the package.json contents:

```json
{
  "name": "midnight-forge",
  "private": true,
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint",
    "test": "vitest run",
    "test:watch": "vitest"
  },
  "dependencies": {
    "@hookform/resolvers": "^3.10.0",
    "@radix-ui/react-accordion": "^1.2.11",
    "@radix-ui/react-alert-dialog": "^1.1.14",
    "@radix-ui/react-aspect-ratio": "^1.1.7",
    "@radix-ui/react-avatar": "^1.1.10",
    "@radix-ui/react-checkbox": "^1.3.2",
    "@radix-ui/react-collapsible": "^1.1.11",
    "@radix-ui/react-context-menu": "^2.2.15",
    "@radix-ui/react-dialog": "^1.1.14",
    "@radix-ui/react-dropdown-menu": "^2.1.15",
    "@radix-ui/react-hover-card": "^1.1.14",
    "@radix-ui/react-label": "^2.1.7",
    "@radix-ui/react-menubar": "^1.1.15",
    "@radix-ui/react-navigation-menu": "^1.2.13",
    "@radix-ui/react-popover": "^1.1.14",
    "@radix-ui/react-progress": "^1.1.7",
    "@radix-ui/react-radio-group": "^1.3.7",
    "@radix-ui/react-scroll-area": "^1.2.9",
    "@radix-ui/react-select": "^2.2.5",
    "@radix-ui/react-separator": "^1.1.7",
    "@radix-ui/react-slider": "^1.3.5",
    "@radix-ui/react-slot": "^1.2.3",
    "@radix-ui/react-switch": "^1.2.5",
    "@radix-ui/react-tabs": "^1.1.12",
    "@radix-ui/react-toast": "^1.2.14",
    "@radix-ui/react-toggle": "^1.1.9",
    "@radix-ui/react-toggle-group": "^1.1.10",
    "@radix-ui/react-tooltip": "^1.2.7",
    "@tanstack/react-query": "^5.83.0",
    "class-variance-authority": "^0.7.1",
    "clsx": "^2.1.1",
    "cmdk": "^1.1.1",
    "date-fns": "^3.6.0",
    "embla-carousel-react": "^8.6.0",
    "input-otp": "^1.4.2",
    "lucide-react": "^0.462.0",
    "next": "^15.3.1",
    "next-themes": "^0.3.0",
    "react": "^19.1.0",
    "react-day-picker": "^8.10.1",
    "react-dom": "^19.1.0",
    "react-hook-form": "^7.61.1",
    "react-resizable-panels": "^2.1.9",
    "recharts": "^2.15.4",
    "sonner": "^1.7.4",
    "tailwind-merge": "^2.6.0",
    "tailwindcss-animate": "^1.0.7",
    "vaul": "^0.9.9",
    "zod": "^3.25.76",
    "zustand": "^5.0.12"
  },
  "devDependencies": {
    "@eslint/js": "^9.32.0",
    "@tailwindcss/typography": "^0.5.16",
    "@testing-library/jest-dom": "^6.6.0",
    "@testing-library/react": "^16.0.0",
    "@types/node": "^22.16.5",
    "@types/react": "^19.1.2",
    "@types/react-dom": "^19.1.2",
    "eslint": "^9.32.0",
    "eslint-config-next": "^15.3.1",
    "globals": "^15.15.0",
    "jsdom": "^20.0.3",
    "postcss": "^8.5.6",
    "tailwindcss": "^3.4.17",
    "typescript": "^5.8.3",
    "typescript-eslint": "^8.38.0",
    "vitest": "^3.2.4"
  }
}
```

- [ ] **Step 2: Create next.config.ts**

```ts
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  devIndicators: false,
};

export default nextConfig;
```

- [ ] **Step 3: Update tsconfig.json**

```json
{
  "compilerOptions": {
    "target": "ES2017",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": false,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "noUnusedLocals": false,
    "noUnusedParameters": false,
    "strictNullChecks": false,
    "plugins": [{ "name": "next" }],
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
```

- [ ] **Step 4: Update .gitignore**

Append to the existing `.gitignore`:

```
# Next.js
.next/
```

Also remove the `dist/` line (or keep it — doesn't hurt).

- [ ] **Step 5: Install dependencies**

Run: `bun install` (or `bun install --ignore-scripts` if postinstall scripts fail)

- [ ] **Step 6: Commit**

```bash
git add frontend/package.json frontend/next.config.ts frontend/tsconfig.json frontend/.gitignore frontend/bun.lockb
git commit -m "feat: scaffold Next.js project structure"
```

---

### Task 2: Create root layout, globals.css, and Providers

**Files:**
- Create: `frontend/src/app/globals.css`
- Create: `frontend/src/app/layout.tsx`
- Create: `frontend/src/components/providers.tsx`
- Modify: `frontend/src/index.css` (delete after moving content)
- Modify: `frontend/tailwind.config.ts`

- [ ] **Step 1: Create `src/app/globals.css`**

Content is the same as `src/index.css` but with one change: remove the `@import url(...)` line for Google Fonts (handled by `next/font` instead).

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    --background: 220 15% 6%;
    --foreground: 30 25% 92%;
    --card: 220 15% 8%;
    --card-foreground: 30 25% 92%;
    --popover: 220 18% 9%;
    --popover-foreground: 30 25% 92%;
    --primary: 18 85% 60%;
    --primary-foreground: 220 20% 8%;
    --primary-glow: 22 95% 68%;
    --primary-deep: 14 70% 45%;
    --secondary: 220 14% 12%;
    --secondary-foreground: 30 20% 88%;
    --muted: 220 14% 11%;
    --muted-foreground: 30 10% 55%;
    --accent: 18 85% 60%;
    --accent-foreground: 220 20% 8%;
    --destructive: 0 75% 55%;
    --destructive-foreground: 30 25% 95%;
    --border: 220 14% 14%;
    --input: 220 14% 12%;
    --ring: 18 85% 60%;
    --radius: 0.25rem;
    --terminal: 220 18% 5%;
    --terminal-foreground: 30 25% 88%;
    --grid-line: 18 60% 50%;
    --success: 145 65% 50%;
    --gradient-signal: linear-gradient(135deg, hsl(var(--primary)) 0%, hsl(var(--primary-glow)) 100%);
    --gradient-fade: linear-gradient(180deg, hsl(var(--background)) 0%, hsl(220 18% 4%) 100%);
    --gradient-panel: linear-gradient(180deg, hsl(220 15% 9%) 0%, hsl(220 15% 7%) 100%);
    --shadow-signal: 0 0 30px -5px hsl(var(--primary) / 0.4);
    --shadow-panel: 0 8px 40px -12px hsl(0 0% 0% / 0.6);
    --glow-text: 0 0 20px hsl(var(--primary) / 0.5);
    --sidebar-background: 220 18% 5%;
    --sidebar-foreground: 30 15% 80%;
    --sidebar-primary: 18 85% 60%;
    --sidebar-primary-foreground: 220 20% 8%;
    --sidebar-accent: 220 14% 11%;
    --sidebar-accent-foreground: 30 25% 92%;
    --sidebar-border: 220 14% 12%;
    --sidebar-ring: 18 85% 60%;
  }

  .dark {
    --background: 220 15% 6%;
    --foreground: 30 25% 92%;
  }
}

@layer base {
  * {
    @apply border-border;
  }

  html {
    color-scheme: dark;
  }

  body {
    @apply bg-background text-foreground font-mono antialiased;
    background-image:
      linear-gradient(hsl(var(--grid-line) / 0.04) 1px, transparent 1px),
      linear-gradient(90deg, hsl(var(--grid-line) / 0.04) 1px, transparent 1px);
    background-size: 48px 48px;
    background-attachment: fixed;
  }

  ::selection {
    background: hsl(var(--primary) / 0.4);
    color: hsl(var(--foreground));
  }

  ::-webkit-scrollbar { width: 10px; height: 10px; }
  ::-webkit-scrollbar-track { background: hsl(var(--background)); }
  ::-webkit-scrollbar-thumb { background: hsl(var(--border)); border-radius: 4px; }
  ::-webkit-scrollbar-thumb:hover { background: hsl(var(--primary) / 0.4); }
}

@layer components {
  .terminal-label {
    @apply font-mono text-[11px] uppercase tracking-[0.18em] text-primary;
  }

  .panel {
    @apply bg-card border border-border rounded-sm;
    background-image: var(--gradient-panel);
    box-shadow: var(--shadow-panel);
  }

  .glow-text {
    text-shadow: var(--glow-text);
  }

  .scanline::before {
    content: "";
    position: absolute;
    inset: 0;
    background: repeating-linear-gradient(
      0deg,
      transparent 0,
      transparent 2px,
      hsl(var(--foreground) / 0.015) 2px,
      hsl(var(--foreground) / 0.015) 4px
    );
    pointer-events: none;
  }

  .heading-display {
    font-family: var(--font-display), 'Space Grotesk', sans-serif;
    @apply font-bold tracking-tight;
  }
}

@layer utilities {
  .text-glow-primary {
    text-shadow: 0 0 16px hsl(var(--primary) / 0.6);
  }
  .border-signal {
    border-color: hsl(var(--primary) / 0.5);
  }
}
```

Note: `.heading-display` now references `var(--font-display)` as a fallback. The `body` still uses `font-mono` via Tailwind — we'll set the CSS variables in the layout.

- [ ] **Step 2: Create `src/app/layout.tsx`**

```tsx
import type { Metadata } from "next";
import { JetBrains_Mono, Space_Grotesk } from "next/font/google";
import { Providers } from "@/components/providers";
import "./globals.css";

const jetbrainsMono = JetBrains_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
  display: "swap",
});

const spaceGrotesk = Space_Grotesk({
  subsets: ["latin"],
  variable: "--font-display",
  display: "swap",
});

export const metadata: Metadata = {
  title: "Midnight Forge — Terminal Forum for Architects & Engineers",
  description:
    "Midnight Forge is a terminal-style technical forum for architects, engineers, and protocol designers. Discuss systems, share docs, build the network.",
  authors: [{ name: "Midnight Forge" }],
  openGraph: {
    title: "Midnight Forge — Terminal Forum",
    description: "Terminal-style technical forum for architects and engineers.",
    type: "website",
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className={`${jetbrainsMono.variable} ${spaceGrotesk.variable}`}>
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
```

- [ ] **Step 3: Create `src/components/providers.tsx`**

```tsx
"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { TooltipProvider } from "@/components/ui/tooltip";
import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { useState } from "react";

export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(() => new QueryClient());

  return (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>
        <Toaster />
        <Sonner />
        {children}
      </TooltipProvider>
    </QueryClientProvider>
  );
}
```

- [ ] **Step 4: Update `tailwind.config.ts`**

Change the `fontFamily` section to use CSS variables instead of direct font names:

```ts
fontFamily: {
  mono: ["var(--font-mono)", "JetBrains Mono", "ui-monospace", "monospace"],
  display: ["var(--font-display)", "Space Grotesk", "system-ui", "sans-serif"],
},
```

- [ ] **Step 5: Verify the app loads**

Run: `bun run dev` in the `frontend/` directory.

At this point, the dev server should start but pages will 404 (no pages exist yet). That's expected. Kill the server after confirming it starts without errors.

Expected log: `▲ Next.js ...` followed by the local server URL.

- [ ] **Step 6: Commit**

```bash
git add frontend/src/app/ frontend/src/components/providers.tsx frontend/tailwind.config.ts
git commit -m "feat: add root layout with fonts, globals, and providers"
```

---

### Task 3: Create route group layouts and navigation

**Files:**
- Create: `frontend/src/app/(auth)/layout.tsx`
- Create: `frontend/src/app/(main)/layout.tsx`
- Create: `frontend/src/components/ActiveLink.tsx`
- Modify: `frontend/src/components/forge/TopNav.tsx`
- Modify: `frontend/src/components/forge/Sidebar.tsx`
- Delete: `frontend/src/components/forge/AppLayout.tsx`
- Delete: `frontend/src/components/NavLink.tsx`

- [ ] **Step 1: Create `src/components/ActiveLink.tsx`**

```tsx
"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";

interface ActiveLinkProps extends React.ComponentPropsWithoutRef<typeof Link> {
  activeClassName?: string;
  end?: boolean;
}

export function ActiveLink({
  className,
  activeClassName,
  end = false,
  href,
  ...props
}: ActiveLinkProps) {
  const pathname = usePathname();
  const hrefStr = typeof href === "string" ? href : href.toString();
  const isActive = end ? pathname === hrefStr : pathname.startsWith(hrefStr);

  return (
    <Link
      href={href}
      className={cn(className, isActive && activeClassName)}
      {...props}
    />
  );
}
```

- [ ] **Step 2: Update `TopNav.tsx`**

Replace imports: remove `NavLink` import, add `ActiveLink` and `usePathname`. Replace all `<NavLink>` usages with `<ActiveLink>`. The component remains `"use client"`.

```tsx
"use client";

import { ActiveLink } from "@/components/ActiveLink";
import { Bell, Settings, Search, User } from "lucide-react";
import { Input } from "@/components/ui/input";

const navItems = [
  { to: "/", label: "Nexus" },
  { to: "/thread/architectural-shift", label: "Threads" },
  { to: "/profile", label: "Network" },
  { to: "/settings", label: "Terminal" },
];

export const TopNav = () => {
  return (
    <header className="sticky top-0 z-40 border-b border-border/60 bg-background/85 backdrop-blur-xl">
      <div className="flex h-16 items-center gap-6 px-6">
        <ActiveLink to="/" className="flex items-center gap-2 group">
          <div className="h-8 w-8 grid place-items-center bg-gradient-signal rounded-sm font-display font-bold text-primary-foreground">
            M
          </div>
          <span className="font-display font-bold text-lg tracking-wide text-foreground group-hover:text-primary transition-colors">
            MIDNIGHT<span className="text-primary">FORGE</span>
          </span>
        </ActiveLink>

        <nav className="hidden md:flex items-center gap-1 ml-4">
          {navItems.map((item) => (
            <ActiveLink
              key={item.to}
              to={item.to}
              end={item.to === "/"}
              className="px-3 py-1.5 text-xs uppercase tracking-[0.18em] text-muted-foreground hover:text-foreground transition-colors rounded-sm"
              activeClassName="!text-primary border-b-2 border-primary"
            >
              {item.label}
            </ActiveLink>
          ))}
        </nav>

        <div className="flex-1 max-w-md ml-auto relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-muted-foreground" />
          <Input
            placeholder="TERMINAL_SEARCH..."
            className="h-9 pl-9 bg-secondary/60 border-border/80 text-xs uppercase tracking-wider placeholder:text-muted-foreground/60 font-mono focus-visible:ring-primary/40"
          />
        </div>

        <div className="flex items-center gap-2">
          <button className="h-9 w-9 grid place-items-center text-muted-foreground hover:text-primary transition-colors relative">
            <Bell className="h-4 w-4" />
            <span className="absolute top-2 right-2 h-1.5 w-1.5 bg-primary rounded-full animate-pulse-signal" />
          </button>
          <ActiveLink to="/settings" className="h-9 w-9 grid place-items-center text-muted-foreground hover:text-primary transition-colors">
            <Settings className="h-4 w-4" />
          </ActiveLink>
          <ActiveLink to="/profile" className="h-9 w-9 grid place-items-center bg-secondary border border-border rounded-sm hover:border-primary transition-colors">
            <User className="h-4 w-4 text-primary" />
          </ActiveLink>
        </div>
      </div>
    </header>
  );
};
```

- [ ] **Step 3: Update `Sidebar.tsx`**

Replace `NavLink` import with `ActiveLink`. Replace `<NavLink>` elements with `<ActiveLink>`. Keep the rest identical.

```tsx
"use client";

import { ActiveLink } from "@/components/ActiveLink";
import { LayoutGrid, MessagesSquare, Shapes, BookOpen, HelpCircle, Archive, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { usePathname } from "next/navigation";

const items = [
  { to: "/", label: "Dashboard", icon: LayoutGrid },
  { to: "/thread/architectural-shift", label: "Discussions", icon: MessagesSquare },
  { to: "/profile", label: "Categories", icon: Shapes },
  { to: "/settings", label: "Documentation", icon: BookOpen },
];

const footerItems = [
  { to: "/settings", label: "Support", icon: HelpCircle },
  { to: "/settings", label: "Archive", icon: Archive },
];

export const Sidebar = ({ showNewEntry = false }: { showNewEntry?: boolean }) => {
  return (
    <aside className="w-64 shrink-0 border-r border-border/60 bg-sidebar/40 min-h-[calc(100vh-4rem)] flex flex-col">
      <div className="p-4 border-b border-border/60">
        <div className="flex items-center gap-3 p-2 rounded-sm bg-secondary/40">
          <div className="h-10 w-10 bg-gradient-signal grid place-items-center rounded-sm">
            <div className="h-4 w-4 bg-background/30 rotate-45" />
          </div>
          <div>
            <div className="text-sm font-bold text-foreground">Midnight Forge</div>
            <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">Technical Forum</div>
          </div>
        </div>
      </div>

      <nav className="p-3 space-y-1 flex-1">
        {items.map((item) => (
          <ActiveLink
            key={item.label}
            to={item.to}
            end={item.to === "/"}
            className="flex items-center gap-3 px-3 py-2.5 text-sm text-sidebar-foreground rounded-sm hover:bg-sidebar-accent transition-colors group"
            activeClassName="!bg-secondary !text-primary border-l-2 border-primary"
          >
            <item.icon className="h-4 w-4" />
            <span>{item.label}</span>
          </ActiveLink>
        ))}

        {showNewEntry && (
          <div className="pt-4">
            <ActiveLink to="/create">
              <Button className="w-full bg-gradient-signal hover:opacity-90 text-primary-foreground font-bold uppercase tracking-[0.15em] text-xs h-10 rounded-sm">
                <Plus className="h-4 w-4 mr-1" /> New_Entry
              </Button>
            </ActiveLink>
          </div>
        )}
      </nav>

      <div className="p-3 border-t border-border/60 space-y-1">
        {footerItems.map((item) => (
          <ActiveLink
            key={item.label}
            to={item.to}
            className="flex items-center gap-3 px-3 py-2 text-sm text-muted-foreground hover:text-foreground rounded-sm transition-colors"
          >
            <item.icon className="h-4 w-4" />
            <span>{item.label}</span>
          </ActiveLink>
        ))}
      </div>
    </aside>
  );
};
```

- [ ] **Step 4: Create `src/app/(auth)/layout.tsx`**

```tsx
export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <>{children}</>;
}
```

This is intentionally bare — no sidebar, no top nav. The Login page is a full-screen standalone view.

- [ ] **Step 5: Create `src/app/(main)/layout.tsx`**

```tsx
import { TopNav } from "@/components/forge/TopNav";
import { Sidebar } from "@/components/forge/Sidebar";

export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <TopNav />
      <div className="flex flex-1 w-full">
        <Sidebar showNewEntry />
        <main className="flex-1 min-w-0">{children}</main>
      </div>
      <footer className="border-t border-border/60 mt-auto">
        <div className="px-8 py-6 flex flex-wrap items-center justify-between gap-4 text-[11px] uppercase tracking-[0.18em] text-muted-foreground">
          <div className="flex gap-6">
            <span>Manifesto</span>
            <span>Privacy</span>
            <span>Security</span>
          </div>
          <div>© 2026 Midnight Forge // Encrypted Session</div>
        </div>
      </footer>
    </div>
  );
}
```

- [ ] **Step 6: Delete obsolete files**

Delete `src/components/forge/AppLayout.tsx` and `src/components/NavLink.tsx`.

- [ ] **Step 7: Commit**

```bash
git add frontend/src/app/\(auth\)/ frontend/src/app/\(main\)/ frontend/src/components/ActiveLink.tsx frontend/src/components/forge/TopNav.tsx frontend/src/components/forge/Sidebar.tsx
git rm frontend/src/components/forge/AppLayout.tsx frontend/src/components/NavLink.tsx
git commit -m "feat: add route group layouts and ActiveLink component"
```

---

### Task 4: Update API client for Next.js env vars

**Files:**
- Modify: `frontend/src/lib/api.ts`

- [ ] **Step 1: Update `src/lib/api.ts`**

Change `import.meta.env.VITE_API_URL` to `process.env.NEXT_PUBLIC_API_URL`:

```ts
const API_BASE = (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080") + "/api/v1";
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/lib/api.ts
git commit -m "fix: update API env var from VITE_ to NEXT_PUBLIC_"
```

---

### Task 5: Migrate Home page

**Files:**
- Create: `frontend/src/app/page.tsx`
- Create: `frontend/src/app/HomeClient.tsx`

The Home page is split into a server component (fetches data) and a client component (renders interactive UI).

- [ ] **Step 1: Create `src/app/page.tsx` (server component)**

```tsx
const API_BASE = (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080") + "/api/v1";

async function fetchJSON<T>(url: string): Promise<T | null> {
  try {
    const res = await fetch(url, { next: { revalidate: 60 } });
    if (!res.ok) return null;
    return res.json();
  } catch {
    return null;
  }
}

export default async function HomePage() {
  const [featured, trending, threads] = await Promise.all([
    fetchJSON<any>(`${API_BASE}/threads/featured`),
    fetchJSON<any>(`${API_BASE}/threads/trending`),
    fetchJSON<any>(`${API_BASE}/threads?page=1&pageSize=5&sort=latest`),
  ]);

  return (
    <HomeClient
      initialFeatured={featured}
      initialTrending={trending}
      initialThreads={threads}
    />
  );
}
```

- [ ] **Step 2: Create `src/app/HomeClient.tsx`**

This is the existing `Index.tsx` page logic, adapted to receive initial data as props. It should be a `"use client"` component.

```tsx
"use client";

import { useState } from "react";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { ActiveLink } from "@/components/ActiveLink";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ChevronRight, MessageCircle, Filter, LayoutGrid } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { useThreads, useFeaturedThread, useTrendingThreads } from "@/hooks/use-threads";
import type { ThreadItem, ThreadDetail, ThreadsResponse } from "@/lib/mock-data";

// ... rest of Index.tsx content with initial data as fallback
```

Copy the full content of `src/pages/Index.tsx` into this file. Key changes:
- Add `"use client"` at the top
- Add the `HomeClient` function component wrapping the logic
- Remove `AppLayout` import — it's provided by `(main)/layout.tsx`
- Keep `HeroImg` import — will fix in Task 11 when assets are moved

```tsx
"use client";

import { useState } from "react";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { ActiveLink } from "@/components/ActiveLink";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ChevronRight, MessageCircle, Filter, LayoutGrid } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import heroImg from "@/assets/forge-hero.jpg";
import { useThreads, useFeaturedThread, useTrendingThreads } from "@/hooks/use-threads";
import type { ThreadItem, ThreadDetail } from "@/lib/mock-data";

function timeAgo(dateStr: string): string {
  const date = new Date(dateStr);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const hours = Math.floor(diff / (1000 * 60 * 60));
  if (hours < 1) return "JUST NOW";
  if (hours < 24) return hours + "H AGO";
  const days = Math.floor(hours / 24);
  return days + "D AGO";
}

interface HomeClientProps {
  initialFeatured: ThreadDetail | null;
  initialTrending: { threads: ThreadItem[] } | null;
  initialThreads: { threads: ThreadItem[]; pagination: any } | null;
}

const HomeClient = ({ initialFeatured, initialTrending, initialThreads }: HomeClientProps) => {
  const [activeTab, setActiveTab] = useState(0);
  const [page, setPage] = useState(1);
  const [sort, setSort] = useState("latest");

  const featured = useFeaturedThread();
  const trending = useTrendingThreads();
  const threads = useThreads(page, 5, sort);

  const featuredThread = featured.data ?? initialFeatured;
  const trendingData = trending.data ?? (initialTrending ? { threads: initialTrending.threads } : undefined);
  const threadData = threads.data ?? (initialThreads ? { ...initialThreads, isEmpty: !initialThreads.threads?.length } : undefined);

  const featuredSlug = featuredThread?.slug || "architectural-shift";
  const firstTag = featuredThread?.tags?.[0]?.name || "System Core";
  const firstAuthor = featuredThread?.author?.username || "@alpha_lead";
  const firstAuthorInitials = firstAuthor.replace("@", "").slice(0, 2).toUpperCase();

  return (
    <div className="px-8 py-10 max-w-[1400px] mx-auto space-y-10">
      {/* Section: Hero — copy from Index.tsx */}
      <section className="space-y-4 animate-fade-up">
        <div className="flex items-end gap-4">
          <h1 className="heading-display text-5xl md:text-6xl text-foreground">
            The Nexus <span className="text-primary">/</span>
          </h1>
        </div>
        <p className="text-sm text-muted-foreground tracking-wide font-mono">
          Synchronized with terminal.access.level.4
        </p>
      </section>

      {/* Section: Featured + Trending grid — copy from Index.tsx */}
      <section className="grid lg:grid-cols-3 gap-6">
        {/* Featured */}
        {featured.isLoading && !initialFeatured ? (
          <Skeleton className="lg:col-span-2 h-[400px] rounded-sm" />
        ) : (
          <ActiveLink to={"/thread/" + featuredSlug} className="lg:col-span-2 group">
            <article className="panel relative overflow-hidden h-full scanline">
              <div className="relative aspect-[2/1] overflow-hidden bg-terminal">
                <img
                  src={heroImg}
                  alt={featuredThread?.title || "Featured thread"}
                  className="w-full h-full object-cover opacity-70 group-hover:opacity-90 transition-opacity"
                  width={1280}
                  height={640}
                />
                <div className="absolute inset-0 bg-gradient-to-t from-card via-card/40 to-transparent" />
              </div>
              <div className="p-6 space-y-3 -mt-16 relative">
                <div className="flex items-center gap-3">
                  <Badge className="bg-primary/15 text-primary border border-primary/30 rounded-sm uppercase text-[10px] tracking-[0.18em]">
                    {firstTag}
                  </Badge>
                  <span className="text-xs text-muted-foreground uppercase tracking-wider">
                    {featuredThread?.created_at ? timeAgo(featuredThread.created_at) : "2H AGO"}
                  </span>
                </div>
                <h2 className="heading-display text-2xl md:text-3xl text-foreground group-hover:text-primary transition-colors max-w-2xl">
                  {featuredThread?.title || "Architectural Shift"}
                </h2>
                <div className="flex items-center gap-3 pt-1">
                  <div className="h-7 w-7 rounded-full bg-gradient-signal grid place-items-center text-[10px] font-bold text-primary-foreground">
                    {firstAuthorInitials}
                  </div>
                  <span className="text-xs text-muted-foreground">
                    {firstAuthor}{featuredThread?.replies_count != null ? " · " + featuredThread.replies_count + " Replies" : " · 142 Replies"}
                  </span>
                </div>
              </div>
            </article>
          </ActiveLink>
        )}

        {/* Trending */}
        <div className="space-y-4">
          <div className="panel p-5 space-y-4">
            <div className="flex items-center gap-2 text-primary">
              <span className="h-2 w-2 rounded-full bg-primary animate-pulse-signal" />
              <span className="terminal-label">Trending Now</span>
            </div>
            {trending.isLoading && !initialTrending ? (
              <div className="space-y-3">
                {[1, 2, 3].map((i) => (
                  <Skeleton key={i} className="h-12 w-full rounded-sm" />
                ))}
              </div>
            ) : (
              <ul className="space-y-4 divide-y divide-border/60">
                {trendingData?.threads?.map((thread: ThreadItem) => (
                  <li key={thread.id} className="pt-4 first:pt-0 group cursor-pointer">
                    <ActiveLink to={"/thread/" + thread.slug}>
                      <h3 className="text-sm font-semibold text-foreground group-hover:text-primary transition-colors">
                        {thread.title}
                      </h3>
                      <p className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground mt-1">
                        {thread.upvotes} PARTICIPANTS
                      </p>
                    </ActiveLink>
                  </li>
                ))}
              </ul>
            )}
          </div>
          <div className="panel p-5">
            <p className="text-sm text-muted-foreground italic font-mono leading-relaxed">
              &ldquo;The digital architect builds not with stone, but with logic and light.&rdquo;
            </p>
            <p className="text-right terminal-label mt-3 opacity-60">NOCTURNE_OS</p>
          </div>
        </div>
      </section>

      {/* Section: Tabs + Thread list — copy from Index.tsx tabs section through pagination */}
      <section className="space-y-4">
        {/* Tabs */}
        <div className="flex items-center justify-between border-b border-border/60">
          <div className="flex gap-6">
            {["LATEST", "UNSOLVED", "MY POSTS"].map((tab, i) => (
              <button
                key={tab}
                onClick={() => {
                  setActiveTab(i);
                  setSort(i === 0 ? "latest" : "votes");
                  setPage(1);
                }}
                className={`pb-3 text-xs uppercase tracking-[0.18em] transition-colors border-b-2 ${
                  i === activeTab ? "text-primary border-primary" : "text-muted-foreground border-transparent hover:text-foreground"
                }`}
              >
                {tab}
              </button>
            ))}
          </div>
          <div className="flex gap-2 pb-2">
            <button className="h-8 w-8 grid place-items-center bg-secondary/60 border border-border rounded-sm hover:border-primary/40 transition-colors">
              <Filter className="h-3.5 w-3.5 text-muted-foreground" />
            </button>
            <button className="h-8 w-8 grid place-items-center bg-secondary/60 border border-border rounded-sm hover:border-primary/40 transition-colors">
              <LayoutGrid className="h-3.5 w-3.5 text-muted-foreground" />
            </button>
          </div>
        </div>

        {/* Thread rows */}
        <div className="space-y-3">
          {(threads.isLoading && !initialThreads) ? (
            Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={i} className="h-[88px] w-full rounded-sm" />
            ))
          ) : (
            threadData?.threads?.map((thread: ThreadItem) => {
              const tagName = thread.tags?.[0]?.name || "TECHNICAL";
              const authorName = thread.author?.username || "@unknown";
              return (
                <ActiveLink to={"/thread/" + thread.slug} key={thread.id}>
                  <article className="panel p-4 md:p-5 grid grid-cols-[64px,1fr,auto] gap-4 md:gap-6 items-center hover:border-primary/40 transition-colors group">
                    <div className="text-center">
                      <div className="text-2xl font-display font-bold text-foreground">{thread.upvotes}</div>
                      <div className="text-[10px] uppercase tracking-[0.2em] text-muted-foreground">votes</div>
                    </div>
                    <div className="space-y-2 min-w-0">
                      <div className="flex flex-wrap items-center gap-2">
                        <Badge className="bg-primary/10 text-primary border border-primary/30 rounded-sm uppercase text-[10px] tracking-[0.18em]">{tagName}</Badge>
                        <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                          Posted by <span className="text-primary/80">{authorName}</span> &middot; {timeAgo(thread.created_at)}
                        </span>
                      </div>
                      <h3 className="text-base md:text-lg font-semibold text-foreground group-hover:text-primary transition-colors line-clamp-2">
                        {thread.title}
                      </h3>
                    </div>
                    <div className="flex items-center gap-4 text-muted-foreground">
                      <div className="hidden md:flex -space-x-2">
                        {[...Array(3)].map((_, j) => (
                          <div key={j} className="h-7 w-7 rounded-full bg-secondary border-2 border-card" />
                        ))}
                      </div>
                      <div className="flex items-center gap-1.5">
                        <MessageCircle className="h-4 w-4" />
                        <span className="text-sm font-mono">{thread.replies_count}</span>
                      </div>
                      <ChevronRight className="h-4 w-4 opacity-0 group-hover:opacity-100 group-hover:translate-x-1 transition-all" />
                    </div>
                  </article>
                </ActiveLink>
              );
            })
          )}
        </div>

        {/* Pagination */}
        {threadData?.pagination && (
          <div className="flex items-center justify-center gap-2 pt-4">
            <button
              onClick={() => setPage((p: number) => Math.max(1, p - 1))}
              disabled={page <= 1}
              className="px-3 py-1.5 text-[10px] uppercase tracking-[0.18em] text-muted-foreground border border-border rounded-sm hover:border-primary/40 disabled:opacity-40"
            >
              PREV_INDEX
            </button>
            {Array.from(
              { length: Math.min(threadData.pagination.totalPages, 5) },
              (_, i) => i + 1
            ).map((p) => (
              <button
                key={p}
                onClick={() => setPage(p)}
                className={`h-8 w-10 text-xs font-mono rounded-sm border ${
                  p === page
                    ? "bg-primary text-primary-foreground border-primary"
                    : "border-border text-muted-foreground hover:border-primary/40"
                }`}
              >
                {String(p).padStart(2, "0")}
              </button>
            ))}
            {threadData.pagination.totalPages > 5 && (
              <>
                <span className="text-muted-foreground">...</span>
                <button
                  onClick={() => setPage(threadData.pagination.totalPages)}
                  className="h-8 w-10 text-xs font-mono rounded-sm border border-border text-muted-foreground hover:border-primary/40"
                >
                  {String(threadData.pagination.totalPages).padStart(2, "0")}
                </button>
              </>
            )}
            <button
              onClick={() => setPage((p: number) => Math.min(threadData.pagination.totalPages, p + 1))}
              disabled={page >= threadData.pagination.totalPages}
              className="px-3 py-1.5 text-[10px] uppercase tracking-[0.18em] text-muted-foreground border border-border rounded-sm hover:border-primary/40 disabled:opacity-40"
            >
              NEXT_INDEX
            </button>
          </div>
        )}
      </section>
    </div>
  );
};

export default HomeClient;
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/app/page.tsx frontend/src/app/HomeClient.tsx
git commit -m "feat: migrate Home page with server component data fetching"
```

---

### Task 6: Migrate Login page

**Files:**
- Create: `frontend/src/app/(auth)/login/page.tsx`

- [ ] **Step 1: Create `src/app/(auth)/login/page.tsx`**

Copy the full content from `src/pages/Login.tsx`. Key changes:
- Add `"use client"` at the top
- Replace `import { useNavigate, useSearchParams } from "react-router-dom"` with `import { useRouter, useSearchParams } from "next/navigation"`
- Replace `useNavigate()` with `useRouter()`, and `navigate("/")` with `router.push("/")`
- Remove the `AppLayout` import (not needed — (auth) layout provides minimal shell)

```tsx
"use client";

import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { ArrowRight, KeyRound, User, Mail, ShieldCheck, Fingerprint } from "lucide-react";
import { toast } from "sonner";
import { useSignin, useSignup } from "@/hooks/use-auth";

type Mode = "signin" | "signup";

const Login = () => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const initial = searchParams.get("mode") === "signup" ? "signup" : "signin";
  const [mode, setMode] = useState<Mode>(initial);
  const [persist, setPersist] = useState(false);
  const [agree, setAgree] = useState(false);
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const signin = useSignin();
  const signup = useSignup();

  const switchMode = (next: Mode) => {
    setMode(next);
    const params = new URLSearchParams(searchParams.toString());
    if (next === "signup") {
      params.set("mode", "signup");
    } else {
      params.delete("mode");
    }
    router.replace(`/login${params.toString() ? `?${params.toString()}` : ""}`);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (mode === "signup") {
      if (!agree) {
        toast.error("PROTOCOL_ACK required", {
          description: "Acknowledge the Forge Protocol to continue.",
        });
        return;
      }
      if (password !== confirmPassword) {
        toast.error("KEY_MISMATCH", {
          description: "Access codes do not match.",
        });
        return;
      }

      signup.mutate(
        { username, email, password },
        {
          onSuccess: () => {
            toast.success("IDENTITY_FORGED", {
              description: "Identity registered. You may now sign in.",
            });
            switchMode("signin");
          },
          onError: (err) => {
            toast.error("FORGE_FAILED", { description: err.message });
          },
        },
      );
    } else {
      signin.mutate(
        { login: username, password },
        {
          onSuccess: () => {
            toast.success("SESSION_INITIATED", {
              description: "Handshake complete. Routing to nexus...",
            });
            setTimeout(() => router.push("/"), 600);
          },
          onError: (err) => {
            toast.error("AUTH_FAILED", { description: err.message });
          },
        },
      );
    }
  };

  return (
    <div className="min-h-screen bg-background relative overflow-hidden">
      {/* Copy all the JSX from Login.tsx here, replacing navigate() with router.push() */}
      {/* For brevity, the full JSX from Login.tsx goes here — see the existing file */}
      <header className="relative z-10 flex items-center justify-between px-8 py-6">
        <button onClick={() => router.push("/")} className="flex items-center gap-3 group">
          <div className="h-10 w-10 grid place-items-center bg-gradient-signal rounded-sm font-display font-bold text-xl text-primary-foreground transition-transform group-hover:-rotate-6">M</div>
          <div className="font-display font-bold text-xl tracking-wide">
            MIDNIGHT <span className="text-primary">FORGE</span>
          </div>
        </button>
        {/* ... rest of Login.tsx JSX */}
      </header>
      {/* ... full content from Login.tsx */}
    </div>
  );
};

export default Login;
```

Full JSX is identical to `src/pages/Login.tsx` lines 82-347, with `navigate` replaced by `router.push`.

- [ ] **Step 2: Commit**

```bash
git add frontend/src/app/\(auth\)/login/page.tsx
git commit -m "feat: migrate Login page"
```

---

### Task 7: Migrate Thread Detail page

**Files:**
- Create: `frontend/src/app/(main)/thread/[slug]/page.tsx`
- Create: `frontend/src/app/(main)/thread/[slug]/ThreadDetailClient.tsx`

- [ ] **Step 1: Create `src/app/(main)/thread/[slug]/page.tsx`**

```tsx
const API_BASE = (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080") + "/api/v1";

async function fetchJSON<T>(url: string): Promise<T | null> {
  try {
    const res = await fetch(url, { next: { revalidate: 60 } });
    if (!res.ok) return null;
    return res.json();
  } catch {
    return null;
  }
}

export default async function ThreadPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const thread = await fetchJSON<any>(`${API_BASE}/threads/${slug}`);

  return <ThreadDetailClient slug={slug} initialThread={thread} />;
}
```

- [ ] **Step 2: Create `src/app/(main)/thread/[slug]/ThreadDetailClient.tsx`**

Copy from `src/pages/ThreadDetail.tsx` with these changes:
- Add `"use client"` at the top
- Remove `useParams` — slug comes from the server component via props
- Remove `AppLayout` wrapper (layout is provided by `(main)/layout.tsx`)
- Keep all interactive logic (comments, voting, reply form)

```tsx
"use client";

import { useState } from "react";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Skeleton } from "@/components/ui/skeleton";
import { ThumbsUp, ThumbsDown, Share2, Bookmark, Flag, Code2, Link2, Image as ImageIcon } from "lucide-react";
import { toast } from "sonner";
import heroImg from "@/assets/forge-hero.jpg";
import { useThread } from "@/hooks/use-thread";
import { useCreateComment } from "@/hooks/use-comments";
import { useVoteThread } from "@/hooks/use-votes";
import { useAuthStore } from "@/stores/auth-store";
import type { ThreadDetail } from "@/lib/mock-data";

interface ThreadDetailClientProps {
  slug: string;
  initialThread: ThreadDetail | null;
}

const ThreadDetailClient = ({ slug, initialThread }: ThreadDetailClientProps) => {
  const [replyContent, setReplyContent] = useState("");
  const user = useAuthStore((s) => s.user);

  const { data: thread } = useThread(slug);
  const createComment = useCreateComment(slug);
  const voteThread = useVoteThread(slug);

  // Use server data as fallback until client query resolves
  const displayThread = thread ?? initialThread;

  // ... rest of ThreadDetail.tsx logic and JSX
  // Copy from ThreadDetail.tsx lines 26-296, removing AppLayout wrapper
  // All the interactive logic stays identical
};
```

The full JSX is identical to `src/pages/ThreadDetail.tsx` lines 26-296, with two changes: no `AppLayout` wrapper, and data sources from props.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/app/\(main\)/thread/
git commit -m "feat: migrate Thread Detail page with server component"
```

---

### Task 8: Migrate Create Entry page

**Files:**
- Create: `frontend/src/app/(main)/create/page.tsx`

- [ ] **Step 1: Create `src/app/(main)/create/page.tsx`**

Copy from `src/pages/CreateEntry.tsx` with these changes:
- Add `"use client"` at the top
- Replace `import { useNavigate } from "react-router-dom"` with `import { useRouter } from "next/navigation"`
- Replace `navigate(...)` with `router.push(...)`
- Remove `AppLayout` wrapper

```tsx
"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Bold, Italic, Code2, List, Link2, Image as ImageIcon } from "lucide-react";
import { toast } from "sonner";
import { useCreateThread } from "@/hooks/use-thread";
```

Full JSX from `src/pages/CreateEntry.tsx` lines 12-205, minus the `AppLayout` wrapper, with `useNavigate`/`navigate` replaced by `useRouter`/`router.push`.

- [ ] **Step 2: Commit**

```bash
git add frontend/src/app/\(main\)/create/page.tsx
git commit -m "feat: migrate Create Entry page"
```

---

### Task 9: Migrate Profile pages

**Files:**
- Create: `frontend/src/app/(main)/profile/page.tsx`
- Create: `frontend/src/app/(main)/profile/[username]/page.tsx`

- [ ] **Step 1: Create `src/app/(main)/profile/page.tsx`**

This uses the logged-in user's profile from auth store. It's a client component that redirects or shows the current user's profile.

```tsx
"use client";

import { ProfileContent } from "./ProfileContent";

export default function MyProfilePage() {
  return <ProfileContent />;
}
```

- [ ] **Step 2: Create `src/app/(main)/profile/[username]/page.tsx`**

Server component that fetches the user's data:

```tsx
const API_BASE = (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080") + "/api/v1";

async function fetchJSON<T>(url: string): Promise<T | null> {
  try {
    const res = await fetch(url, { next: { revalidate: 60 } });
    if (!res.ok) return null;
    return res.json();
  } catch {
    return null;
  }
}

export default async function UserProfilePage({ params }: { params: Promise<{ username: string }> }) {
  const { username } = await params;
  const data = await fetchJSON<any>(`${API_BASE}/users/${username}`);

  return <ProfileContent username={username} initialProfile={data?.user} />;
}
```

- [ ] **Step 3: Create `src/app/(main)/profile/ProfileContent.tsx`**

Shared client component used by both profile pages. Copy from `src/pages/Profile.tsx` with changes:
- Add `"use client"`
- Remove `useParams` — username comes from props
- Replace `Link` from `react-router-dom` with Next.js `Link`
- Remove `AppLayout` wrapper

```tsx
"use client";

import { useEffect } from "react";
import { notFound, useParams } from "next/navigation";
import Link from "next/link";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Fingerprint } from "lucide-react";
import { useUserProfile, useUserThreads } from "@/hooks/use-user";
import type { UserProfile } from "@/lib/mock-data";

interface ProfileContentProps {
  username?: string;
  initialProfile?: UserProfile | null;
}

const ProfileContent = ({ username: propUsername, initialProfile }: ProfileContentProps) => {
  const params = useParams<{ username?: string }>();
  const username = propUsername || params?.username || "";

  const { data: profile, isLoading } = useUserProfile(username);
  // ... rest of Profile.tsx logic with AppLayout wrapper removed
};
```

Full JSX from `src/pages/Profile.tsx` lines 37-224, minus `AppLayout`, with `react-router-dom` imports replaced by Next.js equivalents.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/app/\(main\)/profile/
git commit -m "feat: migrate Profile pages"
```

---

### Task 10: Migrate Settings page

**Files:**
- Create: `frontend/src/app/(main)/settings/page.tsx`

- [ ] **Step 1: Create `src/app/(main)/settings/page.tsx`**

Copy from `src/pages/Settings.tsx` with these changes:
- Add `"use client"`
- Replace `useNavigate`/`useNavigate` with `useRouter`/`router.push`
- Remove `AppLayout` wrapper

```tsx
"use client";

import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { Slider } from "@/components/ui/slider";
import { Shield, Bell, Palette, Eye, Plug, CreditCard, LogOut, Terminal as TerminalIcon, RotateCcw } from "lucide-react";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { useAuthStore } from "@/stores/auth-store";
import { useSignout } from "@/hooks/use-auth";
import { useUpdateProfile } from "@/hooks/use-user";
```

Full JSX from `src/pages/Settings.tsx` lines 14-250, minus `AppLayout`, with `useNavigate`/`navigate` replaced by `useRouter`/`router.push`.

- [ ] **Step 2: Commit**

```bash
git add frontend/src/app/\(main\)/settings/page.tsx
git commit -m "feat: migrate Settings page"
```

---

### Task 11: Migrate 404 page and cleanup old Vite files

**Files:**
- Create: `frontend/src/app/not-found.tsx`
- Move: `frontend/src/assets/*` → `frontend/public/images/`
- Delete: `frontend/index.html`
- Delete: `frontend/vite.config.ts`
- Delete: `frontend/postcss.config.js`
- Delete: `frontend/tsconfig.app.json`
- Delete: `frontend/tsconfig.node.json`
- Delete: `frontend/src/pages/` (entire directory)
- Delete: `frontend/src/App.tsx`
- Delete: `frontend/src/App.css`
- Delete: `frontend/src/main.tsx`
- Delete: `frontend/src/vite-env.d.ts`
- Delete: `frontend/src/index.css`
- Delete: `frontend/dist/` (if exists)

- [ ] **Step 1: Create `src/app/not-found.tsx`**

```tsx
import Link from "next/link";

export default function NotFound() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-muted">
      <div className="text-center">
        <h1 className="mb-4 text-4xl font-bold">404</h1>
        <p className="mb-4 text-xl text-muted-foreground">Oops! Page not found</p>
        <Link href="/" className="text-primary underline hover:text-primary/90">
          Return to Home
        </Link>
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Move assets**

Move `src/assets/forge-hero.jpg` → `public/images/forge-hero.jpg`
Move `src/assets/forge-avatar.jpg` → `public/images/forge-avatar.jpg`

- [ ] **Step 3: Update image references**

In all files that import images from `@/assets/`, update to use the new path:

- `Index.tsx` → `HomeClient.tsx`: Change `import heroImg from "@/assets/forge-hero.jpg"` to `import heroImg from "@/public/images/forge-hero.jpg"` — actually, in Next.js with `public/`, images should be referenced as:
  
  ```tsx
  // Remove the import, use a path string instead
  <img src="/images/forge-hero.jpg" ... />
  ```

Wait — with Next.js, static files in `public/` are served at the root path. So update all image references:

In `HomeClient.tsx`:
```tsx
// Remove: import heroImg from "@/assets/forge-hero.jpg";

// In JSX:
<img src="/images/forge-hero.jpg" alt="..." ... />
```

In `ThreadDetailClient.tsx`:
```tsx
// Remove: import heroImg from "@/assets/forge-hero.jpg";
// Update:
<img src="/images/forge-hero.jpg" alt="..." ... />
```

In `ProfileContent.tsx`:
```tsx
// Remove: import avatarImg from "@/assets/forge-avatar.jpg";
// Update avatar fallback:
{profile?.avatar ? (
  <img src={profile.avatar} alt="Avatar" ... />
) : (
  <img src="/images/forge-avatar.jpg" alt="Node avatar" ... />
)}
```

- [ ] **Step 4: Delete obsolete files**

```bash
# Remove Vite config files
rm frontend/index.html
rm frontend/vite.config.ts
rm frontend/postcss.config.js
rm frontend/tsconfig.app.json
rm frontend/tsconfig.node.json

# Remove old pages dir
rm -rf frontend/src/pages

# Remove obsolete entry files
rm frontend/src/App.tsx
rm frontend/src/App.css
rm frontend/src/main.tsx
rm frontend/src/vite-env.d.ts
rm frontend/src/index.css

# Remove old build output
rm -rf frontend/dist
```

- [ ] **Step 5: Commit**

```bash
git add frontend/src/app/not-found.tsx frontend/public/
git rm frontend/index.html frontend/vite.config.ts frontend/postcss.config.js frontend/tsconfig.app.json frontend/tsconfig.node.json
git rm -r frontend/src/pages
git rm frontend/src/App.tsx frontend/src/App.css frontend/src/main.tsx frontend/src/vite-env.d.ts frontend/src/index.css
git rm -rf frontend/dist 2>/dev/null || true
git commit -m "feat: migrate 404 page and cleanup Vite files"
```

---

### Task 12: Verify build, lint, and tests

**Files:**
- Modify: `frontend/vitest.config.ts` (if needed)
- Verify: full build and test suite

- [ ] **Step 1: Update vitest.config.ts if needed**

Check if the current vitest config works with Next.js. If it references Vite plugins, update:

```ts
import { defineConfig } from "vitest/config";
import path from "path";

export default defineConfig({
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  test: {
    globals: true,
    environment: "jsdom",
    setupFiles: "./src/test/setup.ts",
  },
});
```

- [ ] **Step 2: Run the test suite**

```bash
cd frontend && bun test
```

Expected: All existing tests pass.

If tests fail, fix any import path issues and re-run.

- [ ] **Step 3: Run the build**

```bash
cd frontend && bun run build
```

Expected: Successful build with no TypeScript errors.

- [ ] **Step 4: Run the linter**

```bash
cd frontend && bun run lint
```

Expected: No lint errors. If `next lint` is not configured, run `npx next lint` once to set up ESLint config for Next.js.

- [ ] **Step 5: Commit final changes**

```bash
git add -A
git commit -m "chore: update vitest config and fix build issues"
```

---

### Task 13: Final verification

**Files:** N/A — full system check.

- [ ] **Step 1: Start dev server**

```bash
cd frontend && bun run dev
```

- [ ] **Step 2: Verify each route loads correctly**

Manually check:
- `http://localhost:3000/` — Home page renders
- `http://localhost:3000/login` — Login page renders
- `http://localhost:3000/login?mode=signup` — Signup mode renders
- `http://localhost:3000/thread/architectural-shift` — Thread detail loads
- `http://localhost:3000/create` — Create entry renders
- `http://localhost:3000/profile` — Profile page renders
- `http://localhost:3000/settings` — Settings renders
- `http://localhost:3000/nonexistent` — 404 page renders

- [ ] **Step 3: Done**

Migration complete.
