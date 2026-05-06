import { useState } from "react";
import { AppLayout } from "@/components/forge/AppLayout";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { NavLink } from "@/components/NavLink";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ChevronRight, MessageCircle, Filter, LayoutGrid } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import heroImg from "@/assets/forge-hero.jpg";
import { useThreads, useFeaturedThread, useTrendingThreads } from "@/hooks/use-threads";

const tabs = ["LATEST", "UNSOLVED", "MY POSTS"];

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

const Index = () => {
  const [activeTab, setActiveTab] = useState(0);
  const [page, setPage] = useState(1);
  const [sort, setSort] = useState("latest");

  const featured = useFeaturedThread();
  const trending = useTrendingThreads();
  const threads = useThreads(page, 5, sort);

  const featuredThread = featured.data;
  const trendingData = trending.data;
  const threadData = threads.data;

  const featuredSlug = featuredThread?.slug || "architectural-shift";
  const firstTag = featuredThread?.tags?.[0]?.name || "System Core";
  const firstAuthor = featuredThread?.author?.username || "@alpha_lead";
  const firstAuthorInitials = firstAuthor.replace("@", "").slice(0, 2).toUpperCase();

  return (
    <AppLayout showSidebar showNewEntry>
      <div className="px-8 py-10 max-w-[1400px] mx-auto space-y-10">
        {/* Hero */}
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

        {/* Featured + trending grid */}
        <section className="grid lg:grid-cols-3 gap-6">
          {/* Featured */}
          {featured.isLoading ? (
            <Skeleton className="lg:col-span-2 h-[400px] rounded-sm" />
          ) : (
            <NavLink to={"/thread/" + featuredSlug} className="lg:col-span-2 group">
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
            </NavLink>
          )}

          {/* Trending */}
          <div className="space-y-4">
            <div className="panel p-5 space-y-4">
              <div className="flex items-center gap-2 text-primary">
                <span className="h-2 w-2 rounded-full bg-primary animate-pulse-signal" />
                <span className="terminal-label">Trending Now</span>
              </div>
              {trending.isLoading ? (
                <div className="space-y-3">
                  {[1, 2, 3].map((i) => (
                    <Skeleton key={i} className="h-12 w-full rounded-sm" />
                  ))}
                </div>
              ) : (
                <ul className="space-y-4 divide-y divide-border/60">
                  {trendingData?.threads?.map((thread) => (
                    <li key={thread.id} className="pt-4 first:pt-0 group cursor-pointer">
                      <NavLink to={"/thread/" + thread.slug}>
                        <h3 className="text-sm font-semibold text-foreground group-hover:text-primary transition-colors">
                          {thread.title}
                        </h3>
                        <p className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground mt-1">
                          {thread.upvotes} PARTICIPANTS
                        </p>
                      </NavLink>
                    </li>
                  ))}
                </ul>
              )}
            </div>
            <div className="panel p-5">
              <p className="text-sm text-muted-foreground italic font-mono leading-relaxed">
                "The digital architect builds not with stone, but with logic and light."
              </p>
              <p className="text-right terminal-label mt-3 opacity-60">NOCTURNE_OS</p>
            </div>
          </div>
        </section>

        {/* Tabs */}
        <section className="space-y-4">
          <div className="flex items-center justify-between border-b border-border/60">
            <div className="flex gap-6">
              {tabs.map((tab, i) => (
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
            {threads.isLoading ? (
              Array.from({ length: 5 }).map((_, i) => (
                <Skeleton key={i} className="h-[88px] w-full rounded-sm" />
              ))
            ) : (
              threadData?.threads?.map((thread) => {
                const tagName = thread.tags?.[0]?.name || "TECHNICAL";
                const authorName = thread.author?.username || "@unknown";
                return (
                  <NavLink to={"/thread/" + thread.slug} key={thread.id}>
                    <article className="panel p-4 md:p-5 grid grid-cols-[64px,1fr,auto] gap-4 md:gap-6 items-center hover:border-primary/40 transition-colors group">
                      <div className="text-center">
                        <div className="text-2xl font-display font-bold text-foreground">{thread.upvotes}</div>
                        <div className="text-[10px] uppercase tracking-[0.2em] text-muted-foreground">votes</div>
                      </div>
                      <div className="space-y-2 min-w-0">
                        <div className="flex flex-wrap items-center gap-2">
                          <Badge className="bg-primary/10 text-primary border border-primary/30 rounded-sm uppercase text-[10px] tracking-[0.18em]">{tagName}</Badge>
                          <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                            Posted by <span className="text-primary/80">{authorName}</span> · {timeAgo(thread.created_at)}
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
                  </NavLink>
                );
              })
            )}
          </div>

          {/* Pagination */}
          {threadData?.pagination && (
            <div className="flex items-center justify-center gap-2 pt-4">
              <button
                onClick={() => setPage((p) => Math.max(1, p - 1))}
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
                onClick={() => setPage((p) => Math.min(threadData.pagination.totalPages, p + 1))}
                disabled={page >= threadData.pagination.totalPages}
                className="px-3 py-1.5 text-[10px] uppercase tracking-[0.18em] text-muted-foreground border border-border rounded-sm hover:border-primary/40 disabled:opacity-40"
              >
                NEXT_INDEX
              </button>
            </div>
          )}
        </section>
      </div>
    </AppLayout>
  );
};

export default Index;
