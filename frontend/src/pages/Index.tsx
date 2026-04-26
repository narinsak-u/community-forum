import { AppLayout } from "@/components/forge/AppLayout";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { NavLink } from "@/components/NavLink";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ChevronRight, MessageCircle, ArrowUp, Filter, LayoutGrid, Plus } from "lucide-react";
import heroImg from "@/assets/forge-hero.jpg";

const trending = [
  { title: "State management in Rust-based UI", meta: "89 PARTICIPANTS", tag: "ANNOUNCEMENT" },
  { title: "Legacy code archive is now live for v1 members", meta: "ANNOUNCEMENT", tag: "" },
  { title: "Shader optimization techniques for WebGL", meta: "DESIGN DOCS", tag: "" },
];

const threads = [
  { votes: 24, tag: "DESIGN", author: "@void_strider", time: "4H AGO", title: "Best practices for CSS containment in complex dashboards?", comments: 12 },
  { votes: 56, tag: "ANNOUNCEMENTS", author: "@system_admin", time: "12H AGO", title: "Server Maintenance: Scheduled migration to the new edge nodes", comments: 8 },
  { votes: 8, tag: "TECHNICAL", author: "@codex_null", time: "1D AGO", title: "Optimizing garbage collection for high-throughput node.js streams", comments: 45 },
  { votes: 112, tag: "THE VOID", author: "@neural_link", time: "2D AGO", title: "Has anyone successfully implemented eye-tracking navigation in VR?", comments: 301 },
  { votes: 38, tag: "TECHNICAL", author: "@quantum_byte", time: "3D AGO", title: "Edge-cached event streams: predictable hydration patterns explored", comments: 22 },
];

const Index = () => {
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
          <NavLink to="/thread/architectural-shift" className="lg:col-span-2 group">
            <article className="panel relative overflow-hidden h-full scanline">
              <div className="relative aspect-[2/1] overflow-hidden bg-terminal">
                <img
                  src={heroImg}
                  alt="Architectural shift wireframe"
                  className="w-full h-full object-cover opacity-70 group-hover:opacity-90 transition-opacity"
                  width={1280}
                  height={640}
                />
                <div className="absolute inset-0 bg-gradient-to-t from-card via-card/40 to-transparent" />
              </div>
              <div className="p-6 space-y-3 -mt-16 relative">
                <div className="flex items-center gap-3">
                  <Badge className="bg-primary/15 text-primary border border-primary/30 rounded-sm uppercase text-[10px] tracking-[0.18em]">System Core</Badge>
                  <span className="text-xs text-muted-foreground uppercase tracking-wider">2H ago</span>
                </div>
                <h2 className="heading-display text-2xl md:text-3xl text-foreground group-hover:text-primary transition-colors max-w-2xl">
                  Architectural Shift: Transitioning to Reactive Components in v2.4
                </h2>
                <div className="flex items-center gap-3 pt-1">
                  <div className="h-7 w-7 rounded-full bg-gradient-signal grid place-items-center text-[10px] font-bold text-primary-foreground">AL</div>
                  <span className="text-xs text-muted-foreground">@alpha_lead · 142 Replies</span>
                </div>
              </div>
            </article>
          </NavLink>

          {/* Trending */}
          <div className="space-y-4">
            <div className="panel p-5 space-y-4">
              <div className="flex items-center gap-2 text-primary">
                <span className="h-2 w-2 rounded-full bg-primary animate-pulse-signal" />
                <span className="terminal-label">Trending Now</span>
              </div>
              <ul className="space-y-4 divide-y divide-border/60">
                {trending.map((t, i) => (
                  <li key={i} className="pt-4 first:pt-0 group cursor-pointer">
                    <h3 className="text-sm font-semibold text-foreground group-hover:text-primary transition-colors">
                      {t.title}
                    </h3>
                    <p className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground mt-1">{t.meta}</p>
                  </li>
                ))}
              </ul>
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
              {["LATEST", "UNSOLVED", "MY POSTS"].map((tab, i) => (
                <button
                  key={tab}
                  className={`pb-3 text-xs uppercase tracking-[0.18em] transition-colors border-b-2 ${
                    i === 0 ? "text-primary border-primary" : "text-muted-foreground border-transparent hover:text-foreground"
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
            {threads.map((t, i) => (
              <NavLink to="/thread/architectural-shift" key={i}>
                <article className="panel p-4 md:p-5 grid grid-cols-[64px,1fr,auto] gap-4 md:gap-6 items-center hover:border-primary/40 transition-colors group">
                  <div className="text-center">
                    <div className="text-2xl font-display font-bold text-foreground">{t.votes}</div>
                    <div className="text-[10px] uppercase tracking-[0.2em] text-muted-foreground">votes</div>
                  </div>
                  <div className="space-y-2 min-w-0">
                    <div className="flex flex-wrap items-center gap-2">
                      <Badge className="bg-primary/10 text-primary border border-primary/30 rounded-sm uppercase text-[10px] tracking-[0.18em]">{t.tag}</Badge>
                      <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                        Posted by <span className="text-primary/80">{t.author}</span> · {t.time}
                      </span>
                    </div>
                    <h3 className="text-base md:text-lg font-semibold text-foreground group-hover:text-primary transition-colors line-clamp-2">
                      {t.title}
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
                      <span className="text-sm font-mono">{t.comments}</span>
                    </div>
                    <ChevronRight className="h-4 w-4 opacity-0 group-hover:opacity-100 group-hover:translate-x-1 transition-all" />
                  </div>
                </article>
              </NavLink>
            ))}
          </div>

          {/* Pagination */}
          <div className="flex items-center justify-center gap-2 pt-4">
            <button className="px-3 py-1.5 text-[10px] uppercase tracking-[0.18em] text-muted-foreground border border-border rounded-sm hover:border-primary/40">
              PREV_INDEX
            </button>
            {["01", "02", "03", "...", "48"].map((p, i) => (
              <button
                key={p}
                className={`h-8 w-10 text-xs font-mono rounded-sm border ${
                  i === 0 ? "bg-primary text-primary-foreground border-primary" : "border-border text-muted-foreground hover:border-primary/40"
                }`}
              >
                {p}
              </button>
            ))}
            <button className="px-3 py-1.5 text-[10px] uppercase tracking-[0.18em] text-muted-foreground border border-border rounded-sm hover:border-primary/40">
              NEXT_INDEX
            </button>
          </div>
        </section>
      </div>
    </AppLayout>
  );
};

export default Index;
