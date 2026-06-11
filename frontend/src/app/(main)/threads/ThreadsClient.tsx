"use client";

import Image from "next/image";
import { useState, startTransition, useMemo } from "react";
import { ActiveLink } from "@/components/ActiveLink";
import { Badge } from "@/components/ui/badge";
import { ChevronRight, MessageCircle } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { ThreadTabs } from "@/components/ThreadTabs";
import {
  useThreads,
  useFeaturedThread,
  useTrendingThreads,
  useMyThreads,
} from "@/hooks/use-threads";
import { useMe } from "@/hooks/use-auth";
import type { ThreadItem, ThreadDetail } from "@/lib/mock-data";
import { timeAgo } from "@/lib/utils";

interface ThreadsClientProps {
  initialFeatured: ThreadDetail | null;
  initialTrending: { threads: ThreadItem[] } | null;
  initialThreads: { threads: ThreadItem[]; pagination: any } | null;
}

const ThreadsClient = ({
  initialFeatured,
  initialTrending,
  initialThreads,
}: ThreadsClientProps) => {
  const [activeTab, setActiveTab] = useState("latest");
  const [page, setPage] = useState(1);
  const [sort, setSort] = useState("latest");

  const { data: currentUser } = useMe();
  const featured = useFeaturedThread(initialFeatured ?? undefined);
  const trending = useTrendingThreads(
    initialTrending ? { threads: initialTrending.threads } : undefined,
  );
  const initialThreadsData = initialThreads
    ? ({ threads: initialThreads.threads, pagination: initialThreads.pagination } as const)
    : undefined;
  const allThreads = useThreads(page, 5, sort, {
    initialData: initialThreadsData as any,
    enabled: activeTab !== "my-posts",
  });
  const myThreads = useMyThreads(currentUser?.username ?? "", page, 5, {
    enabled: activeTab === "my-posts",
  });
  const threads = activeTab === "my-posts" ? myThreads : allThreads;

  const tabs = useMemo(() => {
    const items: { label: string; value: string; sort: string }[] = [
      { label: "LATEST", value: "latest", sort: "latest" },
      { label: "Most Popular", value: "popular", sort: "votes" },
    ];
    if (currentUser) {
      items.push({ label: "MY POSTS", value: "my-posts", sort: "votes" });
    }
    return items;
  }, [currentUser]);

  const handleTabChange = (tab: (typeof tabs)[number]) => {
    startTransition(() => {
      setActiveTab(tab.value);
      setSort(tab.sort);
      setPage(1);
    });
  };

  const featuredThread = featured.data ?? initialFeatured;
  const trendingData = trending.data;
  const threadData = activeTab !== "my-posts" ? threads.data : threads.data;

  const featuredSlug = featuredThread?.slug || "architectural-shift";
  const firstTag = featuredThread?.tags?.[0]?.name || "System Core";
  const firstAuthor = featuredThread?.author?.username || "@alpha_lead";
  const firstAuthorInitials = firstAuthor
    .replace("@", "")
    .slice(0, 2)
    .toUpperCase();

  return (
    <div className="px-8 py-10 max-w-[1400px] mx-auto space-y-10">
      <section className="space-y-4 animate-fade-up">
        <div className="flex items-end gap-4">
          <h1 className="heading-display text-5xl md:text-6xl text-foreground">
            The Nexus <span className="text-primary">/</span>
          </h1>
        </div>
        <p className="text-sm text-muted-foreground tracking-wide font-mono">
          Synchronized with terminal
        </p>
      </section>

      <section className="grid lg:grid-cols-3 gap-6">
        {featured.isLoading && !initialFeatured ? (
          <Skeleton className="lg:col-span-2 h-[400px] rounded-sm" />
        ) : (
          <ActiveLink
            href={"/thread/" + featuredSlug}
            className="lg:col-span-2 group"
          >
            <article className="panel relative overflow-hidden h-full scanline">
              <div className="relative aspect-[2/1] overflow-hidden bg-terminal">
                <Image
                  src="/images/forge-hero.jpg"
                  alt={featuredThread?.title || "Featured thread"}
                  fill
                  className="object-cover opacity-70 group-hover:opacity-90 transition-opacity"
                  sizes="(max-width: 1200px) 100vw, 800px"
                  priority
                />
                <div className="absolute inset-0 bg-gradient-to-t from-card via-card/40 to-transparent" />
              </div>
              <div className="p-6 space-y-3 -mt-16 relative">
                <div className="flex items-center gap-3">
                  <Badge className="bg-primary/15 text-primary border border-primary/30 rounded-sm uppercase text-[10px] tracking-[0.18em]">
                    {firstTag}
                  </Badge>
                  <span className="text-xs text-muted-foreground uppercase tracking-wider">
                    {featuredThread?.created_at
                      ? timeAgo(featuredThread.created_at)
                      : "2H AGO"}
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
                    {firstAuthor}
                    {featuredThread?.replies_count != null
                      ? " · " + featuredThread.replies_count + " Replies"
                      : " · 142 Replies"}
                  </span>
                </div>
              </div>
            </article>
          </ActiveLink>
        )}

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
                  <li
                    key={thread.id}
                    className="pt-4 first:pt-0 group cursor-pointer"
                  >
                    <ActiveLink href={"/thread/" + thread.slug}>
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
              &ldquo;The digital architect builds not with stone, but with logic
              and light.&rdquo;
            </p>
            <p className="text-right terminal-label mt-3 opacity-60">
              NOCTURNE_OS
            </p>
          </div>
        </div>
      </section>

      <section className="space-y-4">
        <ThreadTabs
          tabs={tabs}
          activeTab={activeTab}
          onTabChange={handleTabChange}
        />

        {/* thread items*/}
        <div className="space-y-3">
          {threads.isLoading ? (
            Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={i} className="h-[88px] w-full rounded-sm" />
            ))
          ) : !threadData?.threads?.length ? (
            <p className="text-sm text-muted-foreground text-center py-12">
              {activeTab === "my-posts"
                ? "You haven't posted yet"
                : "No threads yet"}
            </p>
          ) : (
            threadData?.threads?.map((thread: ThreadItem) => {
              const tagName = thread.tags?.[0]?.name || "TECHNICAL";
              const authorName = thread.author?.username || "@unknown";
              const href =
                thread.status === "draft"
                  ? `/thread/${thread.slug}/edit`
                  : `/thread/${thread.slug}`;
              return (
                <ActiveLink href={href} key={thread.id}>
                  <article className="panel p-4 md:p-5 grid grid-cols-[64px,1fr,auto] gap-4 md:gap-6 items-center hover:border-primary/40 transition-colors group">
                    <div className="text-center">
                      <div className="text-2xl font-display font-bold text-foreground">
                        {thread.upvotes}
                      </div>
                      <div className="text-[10px] uppercase tracking-[0.2em] text-muted-foreground">
                        votes
                      </div>
                    </div>
                    <div className="space-y-2 min-w-0">
                      <div className="flex flex-wrap items-center gap-2">
                        <Badge className="bg-primary/10 text-primary border border-primary/30 rounded-sm uppercase text-[10px] tracking-[0.18em]">
                          {tagName}
                        </Badge>
                        {thread.status === "draft" && (
                          <Badge className="bg-amber-500/10 text-amber-400 border border-amber-500/30 rounded-sm uppercase text-[10px] tracking-[0.18em]">
                            DRAFTED
                          </Badge>
                        )}
                        <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                          Posted by{" "}
                          <span className="text-primary/80">{authorName}</span>{" "}
                          &middot; {timeAgo(thread.created_at)}
                        </span>
                      </div>
                      <h3 className="text-base md:text-lg font-semibold text-foreground group-hover:text-primary transition-colors line-clamp-2">
                        {thread.title}
                      </h3>
                    </div>
                    <div className="flex items-center gap-4 text-muted-foreground">
                      <div className="hidden md:flex -space-x-2">
                        {thread.recent_commenters
                          ?.slice(0, 3)
                          .map((commenter) => (
                            <div
                              key={commenter.id}
                              className="h-7 w-7 rounded-full bg-secondary border-2 border-card overflow-hidden"
                              title={commenter.username}
                            >
                              {commenter.avatar ? (
                                <Image
                                  src={commenter.avatar}
                                  alt={commenter.username}
                                  width={28}
                                  height={28}
                                  className="object-cover w-full h-full"
                                />
                              ) : (
                                <div className="w-full h-full grid place-items-center text-[10px] font-bold text-muted-foreground">
                                  {commenter.username
                                    .replace("@", "")
                                    .slice(0, 2)
                                    .toUpperCase()}
                                </div>
                              )}
                            </div>
                          ))}
                      </div>
                      <div className="flex items-center gap-1.5">
                        <MessageCircle className="h-4 w-4" />
                        <span className="text-sm font-mono">
                          {thread.replies_count}
                        </span>
                      </div>
                      <ChevronRight className="h-4 w-4 opacity-0 group-hover:opacity-100 group-hover:translate-x-1 transition-all" />
                    </div>
                  </article>
                </ActiveLink>
              );
            })
          )}
        </div>

        {threadData?.pagination && (
          <div className="flex items-center justify-center gap-2 pt-4">
            <button
              onClick={() =>
                startTransition(() =>
                  setPage((p: number) => Math.max(1, p - 1)),
                )
              }
              disabled={page <= 1}
              className="px-3 py-1.5 text-[10px] uppercase tracking-[0.18em] text-muted-foreground border border-border rounded-sm hover:border-primary/40 disabled:opacity-40"
            >
              PREV_INDEX
            </button>
            {Array.from(
              { length: Math.min(threadData.pagination.totalPages, 5) },
              (_, i) => i + 1,
            ).map((p) => (
              <button
                key={p}
                onClick={() => startTransition(() => setPage(p))}
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
                  onClick={() =>
                    startTransition(() =>
                      setPage(threadData.pagination.totalPages),
                    )
                  }
                  className="h-8 w-10 text-xs font-mono rounded-sm border border-border text-muted-foreground hover:border-primary/40"
                >
                  {String(threadData.pagination.totalPages).padStart(2, "0")}
                </button>
              </>
            )}
            <button
              onClick={() =>
                startTransition(() =>
                  setPage((p: number) =>
                    Math.min(threadData.pagination.totalPages, p + 1),
                  ),
                )
              }
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

export default ThreadsClient;
