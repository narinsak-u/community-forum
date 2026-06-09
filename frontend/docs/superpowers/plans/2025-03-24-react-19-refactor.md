# React 19 UI Component Refactor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor all UI components in `src/components/ui/` to React 19 standards, removing `React.forwardRef`, using `React.use()` for context, and memoizing context providers.

**Architecture:** 
- Remove `React.forwardRef` and pass `ref` as a prop directly to components.
- Replace `useContext(C)` with `React.use(C)`.
- Wrap context provider `value` props in `useMemo` to prevent "Unstable context provider value" warnings.
- Preserve `displayName` and all existing logic/styling.

**Tech Stack:** React 19, TypeScript, Tailwind CSS, shadcn/ui.

---

### Task 1: Finish Sidebar and fix existing Providers

**Files:**
- Modify: `src/components/ui/sidebar.tsx`
- Modify: `src/components/ui/form.tsx`
- Modify: `src/components/ui/carousel.tsx`
- Modify: `src/components/ui/toggle-group.tsx`
- Modify: `src/components/ui/chart.tsx`

- [ ] **Step 1: Memoize context values in form.tsx**
Wrap the `id` value in `useMemo` in `FormItem`.
- [ ] **Step 2: Memoize context values in carousel.tsx**
Wrap the context object in `useMemo` in `Carousel`.
- [ ] **Step 3: Refactor toggle-group.tsx**
Remove `forwardRef`, pass `ref` as prop, and memoize context value.
- [ ] **Step 4: Memoize context values in chart.tsx**
Wrap the `{ config }` object in `useMemo` in `ChartContainer`.
- [ ] **Step 5: Final check of sidebar.tsx**
Ensure no `forwardRef` remain and all context usage uses `React.use()`.

### Task 2: Refactor Group 1 (Accordion, Alert Dialog, Alert, Avatar)

**Files:**
- Modify: `src/components/ui/accordion.tsx`
- Modify: `src/components/ui/alert-dialog.tsx`
- Modify: `src/components/ui/alert.tsx`
- Modify: `src/components/ui/avatar.tsx`

- [ ] **Step 1: Refactor accordion.tsx**
- [ ] **Step 2: Refactor alert-dialog.tsx**
- [ ] **Step 3: Refactor alert.tsx**
- [ ] **Step 4: Refactor avatar.tsx**

### Task 3: Refactor Group 2 (Breadcrumb, Card, Checkbox, Command)

**Files:**
- Modify: `src/components/ui/breadcrumb.tsx`
- Modify: `src/components/ui/card.tsx`
- Modify: `src/components/ui/checkbox.tsx`
- Modify: `src/components/ui/command.tsx`

- [ ] **Step 1: Refactor breadcrumb.tsx**
- [ ] **Step 2: Refactor card.tsx**
- [ ] **Step 3: Refactor checkbox.tsx**
- [ ] **Step 4: Refactor command.tsx**

### Task 4: Refactor Group 3 (Context Menu, Dialog, Drawer, Dropdown Menu)

**Files:**
- Modify: `src/components/ui/context-menu.tsx`
- Modify: `src/components/ui/dialog.tsx`
- Modify: `src/components/ui/drawer.tsx`
- Modify: `src/components/ui/dropdown-menu.tsx`

- [ ] **Step 1: Refactor context-menu.tsx**
- [ ] **Step 2: Refactor dialog.tsx**
- [ ] **Step 3: Refactor drawer.tsx**
- [ ] **Step 4: Refactor dropdown-menu.tsx**

### Task 5: Refactor Group 4 (Hover Card, Menubar, Navigation Menu, Pagination)

**Files:**
- Modify: `src/components/ui/hover-card.tsx`
- Modify: `src/components/ui/menubar.tsx`
- Modify: `src/components/ui/navigation-menu.tsx`
- Modify: `src/components/ui/pagination.tsx`

- [ ] **Step 1: Refactor hover-card.tsx**
- [ ] **Step 2: Refactor menubar.tsx**
- [ ] **Step 3: Refactor navigation-menu.tsx**
- [ ] **Step 4: Refactor pagination.tsx**

### Task 6: Refactor Group 5 (Popover, Progress, Radio Group, Scroll Area)

**Files:**
- Modify: `src/components/ui/popover.tsx`
- Modify: `src/components/ui/progress.tsx`
- Modify: `src/components/ui/radio-group.tsx`
- Modify: `src/components/ui/scroll-area.tsx`

- [ ] **Step 1: Refactor popover.tsx**
- [ ] **Step 2: Refactor progress.tsx**
- [ ] **Step 3: Refactor radio-group.tsx**
- [ ] **Step 4: Refactor scroll-area.tsx**

### Task 7: Refactor Group 6 (Select, Separator, Sheet, Slider)

**Files:**
- Modify: `src/components/ui/select.tsx`
- Modify: `src/components/ui/separator.tsx`
- Modify: `src/components/ui/sheet.tsx`
- Modify: `src/components/ui/slider.tsx`

- [ ] **Step 1: Refactor select.tsx**
- [ ] **Step 2: Refactor separator.tsx**
- [ ] **Step 3: Refactor sheet.tsx**
- [ ] **Step 4: Refactor slider.tsx**

### Task 8: Refactor Group 7 (Switch, Table, Tabs, Textarea, Toast)

**Files:**
- Modify: `src/components/ui/switch.tsx`
- Modify: `src/components/ui/table.tsx`
- Modify: `src/components/ui/tabs.tsx`
- Modify: `src/components/ui/textarea.tsx`
- Modify: `src/components/ui/toast.tsx`

- [ ] **Step 1: Refactor switch.tsx**
- [ ] **Step 2: Refactor table.tsx**
- [ ] **Step 3: Refactor tabs.tsx**
- [ ] **Step 4: Refactor textarea.tsx**
- [ ] **Step 5: Refactor toast.tsx**

### Task 9: Refactor Group 8 (Toggle, Tooltip)

**Files:**
- Modify: `src/components/ui/toggle.tsx`
- Modify: `src/components/ui/tooltip.tsx`

- [ ] **Step 1: Refactor toggle.tsx**
- [ ] **Step 2: Refactor tooltip.tsx**

### Task 10: Final Verification

- [ ] **Step 1: Run linting**
Run: `bun run lint`
Expected: No errors.

- [ ] **Step 2: Run react-doctor**
Run: `npx -y react-doctor@latest . --verbose --diff`
Expected: High score, no React 19 compatibility issues.
