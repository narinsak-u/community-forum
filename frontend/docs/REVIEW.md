# Codebase Review & Improvement Plan

This review was conducted using specialized React, Next.js, and Shadcn/ui skills. It identifies critical runtime errors, performance bottlenecks, and architectural improvements.

## 1. Critical Runtime Errors

### 1.1 Hydration Mismatch (`TopNav.tsx`)
**Problem:** The theme toggle renders a different icon on the server (defaulting to Sun) than on the client (if the user has Dark mode saved).
**Impact:** Hydration failed error, layout shift, and potential broken interactivity.
**Solution:** 
- Use the `mounted` state correctly to defer rendering of the theme-dependent icon until after hydration.
- Alternatively, render a placeholder with fixed dimensions to prevent layout shift.

### 1.2 State Update During Render (`SessionLoader.tsx` / `Sidebar.tsx`)
**Problem:** `SessionLoader` updates the global `authStore` in a way that triggers a re-render of `Sidebar` (which subscribes to the store) during a parent render cycle.
**Impact:** `Uncaught Error: Cannot update a component (Sidebar) while rendering a different component (SessionLoader).`
**Solution:**
- Ensure `setUser` is only called inside `useEffect`.
- If `useQuery` data is already available, consider initializing the store outside the React render cycle if possible, or use `useSyncExternalStore` for the auth state.
- Refactor `SessionLoader` to avoid unnecessary re-renders.

---

## 2. React 19 & Next.js Best Practices

### 2.1 React 19 API Migration
**Problem:** The project uses React 19 but still employs `forwardRef` and `useContext` patterns.
**Impact:** Deprecation warnings and slightly more verbose code.
**Solution:**
- Replace `forwardRef` with direct `ref` props.
- Replace `useContext(Context)` with `use(Context)`.
- **Target Files:** Most `src/components/ui/` components.

### 2.2 Metadata & SEO
**Problem:** Multiple pages (`profile`, `create`, `thread/[slug]`) are missing Next.js metadata.
**Impact:** Poor SEO and social sharing.
**Solution:** Add static or dynamic `generateMetadata` exports to all route segments.

### 2.3 Image Optimization
**Problem:** Standard `<img>` tags are used in `HomeClient`, `ProfileContent`, etc.
**Impact:** Poor performance (LCP), no automatic WebP/AVIF optimization.
**Solution:** Replace all `<img>` with `next/image`'s `Image` component.

---

## 3. Data Fetching (TanStack Query)

### 3.1 Query Key Factory
**Problem:** Query keys are hardcoded as strings like `["me"]` or `["threads"]`.
**Impact:** Harder to manage cache invalidation, prone to typos.
**Solution:** Implement a central query key factory in `src/hooks/use-query-keys.ts`.

### 3.2 Error Handling
**Problem:** API errors are handled locally in components.
**Solution:** Implement a global error handler in `QueryClient` (via `defaultOptions.mutations.onError`) to show toast notifications for failed operations.

---

## 4. UI & Component Architecture

### 4.1 Accessibility (A11y)
**Findings:**
- Missing `htmlFor` on labels in `login` and `settings` pages.
- Buttons missing explicit `type="button"` or `type="submit"`.
- Interactive controls missing accessible labels (aria-labels).
**Solution:** Audit and fix all `react-doctor` identified accessibility warnings.

### 4.2 Stable Keys
**Problem:** Array indices are used as `key` props in lists.
**Impact:** Incorrect re-rendering and loss of state when lists are filtered/sorted.
**Solution:** Use unique `id` or `slug` from data objects.

### 4.3 Context Provider Optimization
**Problem:** `form.tsx`, `carousel.tsx`, etc., pass objects directly to Context Providers without `useMemo`.
**Impact:** Unnecessary re-renders of all child components on every parent update.
**Solution:** Wrap provider `value` in `useMemo`.

---

## 5. Implementation Roadmap

### Phase 1: Stability (Immediate)
- [ ] Fix `TopNav` hydration mismatch.
- [ ] Fix `SessionLoader` state update error.
- [ ] Fix array index keys in critical paths (`ThreadDetail`, `Home`).

### Phase 2: React 19 Modernization
- [ ] Batch update `ui/` components to remove `forwardRef`.
- [ ] Update `useContext` to `use()`.

### Phase 3: Performance & SEO
- [ ] Convert `<img>` to `next/image`.
- [ ] Add Metadata to all pages.
- [ ] Memoize Context Providers.

### Phase 4: A11y & Cleanup
- [ ] Fix missing labels and button types.
- [ ] Implement Query Key Factory.
- [ ] Resolve all remaining `react-doctor` warnings.
