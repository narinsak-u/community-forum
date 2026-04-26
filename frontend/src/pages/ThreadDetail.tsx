import { AppLayout } from "@/components/forge/AppLayout";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { ThumbsUp, ThumbsDown, Share2, Bookmark, Flag, Code2, Link2, Image as ImageIcon, Bold, Italic, List } from "lucide-react";
import heroImg from "@/assets/forge-hero.jpg";

const replies = [
  {
    user: "@void_iterator",
    badge: "SENIOR_DEV",
    time: "1 HOUR AGO",
    score: 82,
    body: "This transition is overdue. The imperative bloat in v2.3 was becoming unmanageable for large-scale telemetry dashboards. My only concern is the overhead on the primary event bus—have we stress-tested the reactive hydration with >10k concurrent nodes?",
    children: [
      {
        user: "@alpha_lead",
        badge: "OP",
        time: "45 MINS AGO",
        body: "@void_iterator Yes, internal tests on the 'KRYPTON' cluster showed a 14% memory reduction compared to v2.3 under high load. The event bus is now sharded by default.",
      },
    ],
  },
  {
    user: "@null_pointer",
    time: "35 MINS AGO",
    score: 12,
    body: "Will there be an automated migration script for legacy v2.2 legacy_hooks? Or is it a manual rewrite?",
  },
];

const ThreadDetail = () => {
  return (
    <AppLayout showSidebar showNewEntry>
      <div className="px-8 py-10 max-w-[1100px] mx-auto space-y-8 animate-fade-up">
        {/* Breadcrumb */}
        <nav className="flex items-center gap-2 text-[10px] uppercase tracking-[0.2em] text-muted-foreground">
          <span>Discussion</span><span>/</span>
          <span>Technical</span><span>/</span>
          <span className="text-primary">Architectural Shift...</span>
        </nav>

        {/* Title */}
        <header className="space-y-4">
          <Badge className="bg-primary/15 text-primary border border-primary/40 rounded-sm uppercase text-[10px] tracking-[0.2em]">● System Core</Badge>
          <h1 className="heading-display text-3xl md:text-5xl leading-tight text-foreground">
            Architectural Shift: Transitioning to{" "}
            <span className="text-primary">Reactive Components</span> in v2.4
          </h1>
          <div className="flex flex-wrap items-center justify-between gap-4 pt-2">
            <div className="flex items-center gap-3">
              <div className="h-10 w-10 rounded-full bg-gradient-signal grid place-items-center text-xs font-bold text-primary-foreground">AL</div>
              <div>
                <div className="text-sm font-semibold text-foreground">@alpha_lead</div>
                <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">ARCHITECT_ELITE · 2 hours ago</div>
              </div>
            </div>
            <div className="flex gap-8 text-right">
              <div>
                <div className="text-lg font-display font-bold text-foreground">52.4K</div>
                <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">VIEWS</div>
              </div>
              <div>
                <div className="text-lg font-display font-bold text-foreground">142</div>
                <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">REPLIES</div>
              </div>
              <div>
                <div className="text-lg font-display font-bold text-primary">98%</div>
                <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">TRUST_INDEX</div>
              </div>
            </div>
          </div>
        </header>

        {/* Hero image */}
        <div className="panel scanline relative overflow-hidden aspect-[2/1]">
          <img src={heroImg} alt="System wireframe" className="w-full h-full object-cover" width={1280} height={640} />
        </div>

        {/* Body + sidebar */}
        <div className="grid lg:grid-cols-[1fr,260px] gap-8">
          <article className="space-y-6">
            <p className="text-foreground/90 leading-relaxed">
              <span className="text-primary text-3xl font-display float-left mr-2 leading-none">T</span>
              he transition to version 2.4 marks a fundamental departure from our previous imperative component lifecycle. We are moving towards a fully reactive, state-driven architecture that prioritizes deterministic rendering over performance-heavy mutations.
            </p>
            <p className="text-foreground/80 leading-relaxed">
              This shift necessitates a re-evaluation of how we handle global state injections. In the new reactive model, components no longer "fetch" their dependencies; they are "hydrated" by the system's core event bus through a series of prioritized streams.
            </p>

            <blockquote className="border-l-2 border-primary pl-5 py-2 italic text-primary/90">
              "Precision is not just about speed; it's about predictable outcomes in a chaotic data environment."
            </blockquote>

            <h2 className="heading-display text-xl text-foreground pt-2">Implementation Protocol</h2>
            <p className="text-foreground/80 leading-relaxed">
              Below is a demonstration of the new component registration syntax. Notice the lack of manual mounting hooks. The architecture now infers mounting priority based on the <code className="px-1.5 py-0.5 bg-secondary text-primary rounded-sm text-sm">data-stream-id</code>.
            </p>

            <div className="panel bg-terminal/80 overflow-hidden">
              <div className="flex items-center justify-between px-4 py-2 border-b border-border/60 text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                <span>component_manifest.js</span>
                <span>UTF-8</span>
              </div>
              <pre className="p-5 text-sm leading-relaxed font-mono text-foreground/90 overflow-x-auto">
{`export const ReactiveGrid = (props) => {
  const { state, dispatch } = useSystemStream(props.id);

  // v2.4 Reactive Bridge
  useEffect(() => {
    dispatch({ type: 'INIT_SYNC', payload: Date.now() });
  }, [state.isReady]);

  return <div className="neon-grid">{state.nodes}</div>;
};`}
              </pre>
            </div>

            {/* Vote bar */}
            <div className="flex items-center justify-between pt-4 border-t border-border/60">
              <div className="flex items-center gap-3">
                <Button variant="outline" size="sm" className="border-border hover:border-primary hover:text-primary rounded-sm">
                  <ThumbsUp className="h-3.5 w-3.5 mr-2" /> 1,482
                </Button>
                <Button variant="outline" size="sm" className="border-border rounded-sm">
                  <ThumbsDown className="h-3.5 w-3.5" />
                </Button>
              </div>
              <div className="flex items-center gap-4 text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                <button className="flex items-center gap-1.5 hover:text-primary"><Share2 className="h-3.5 w-3.5" /> Share</button>
                <button className="flex items-center gap-1.5 hover:text-primary"><Bookmark className="h-3.5 w-3.5" /> Bookmark</button>
                <button className="flex items-center gap-1.5 hover:text-primary"><Flag className="h-3.5 w-3.5" /> Report</button>
              </div>
            </div>
          </article>

          {/* References sidebar */}
          <aside className="space-y-6">
            <div className="panel p-5 space-y-4">
              <SectionLabel>References</SectionLabel>
              <ul className="space-y-2 text-sm">
                {["REF_CORE_041", "REACTIVE_HOOKS.DOC", "V2.3_LEGACY_PATCH"].map((r) => (
                  <li key={r} className="flex items-center gap-2 text-primary hover:text-primary-glow cursor-pointer font-mono">
                    <Link2 className="h-3 w-3" /> {r}
                  </li>
                ))}
              </ul>
            </div>
            <div className="panel p-5 space-y-3">
              <SectionLabel>Participants</SectionLabel>
              <div className="flex flex-wrap gap-1">
                {[...Array(5)].map((_, i) => (
                  <div key={i} className="h-8 w-8 rounded-full bg-gradient-to-br from-primary/40 to-primary-deep/30 border border-border" />
                ))}
                <div className="h-8 w-8 rounded-full bg-secondary border border-border grid place-items-center text-[10px] font-mono text-primary">+138</div>
              </div>
            </div>
          </aside>
        </div>

        {/* Replies header */}
        <div className="flex items-center justify-between border-t border-border/60 pt-8">
          <h3 className="heading-display text-2xl text-foreground">142 Replies</h3>
          <div className="flex gap-4 text-[10px] uppercase tracking-[0.18em]">
            <span className="text-muted-foreground">SORT_BY:</span>
            <button className="text-primary border-b border-primary">TOP_RATED</button>
            <button className="text-muted-foreground hover:text-foreground">LATEST</button>
          </div>
        </div>

        {/* Replies */}
        <div className="space-y-5">
          {replies.map((r, i) => (
            <div key={i} className="panel p-5 space-y-3">
              <div className="flex items-start justify-between gap-3">
                <div className="flex items-center gap-3">
                  <div className="h-8 w-8 rounded-full bg-secondary border border-border grid place-items-center text-[10px] font-mono text-primary">
                    {r.user.slice(1, 3).toUpperCase()}
                  </div>
                  <div>
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-semibold text-foreground">{r.user}</span>
                      <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">{r.time}</span>
                      {r.badge && <Badge className="bg-primary/10 text-primary border border-primary/30 rounded-sm uppercase text-[9px] tracking-[0.18em]">{r.badge}</Badge>}
                    </div>
                  </div>
                </div>
                <div className="flex items-center gap-2 text-primary text-sm font-mono">
                  +{r.score} <ThumbsUp className="h-3.5 w-3.5" />
                </div>
              </div>
              <p className="text-sm text-foreground/85 leading-relaxed">{r.body}</p>
              <div className="flex gap-4 text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                <button className="hover:text-primary">REPLY</button>
                <button className="hover:text-primary">SHARE</button>
              </div>

              {r.children?.map((c, j) => (
                <div key={j} className="ml-6 mt-4 pl-5 border-l border-primary/30 space-y-2">
                  <div className="flex items-center gap-2">
                    <div className="h-7 w-7 rounded-full bg-gradient-signal grid place-items-center text-[9px] font-bold text-primary-foreground">
                      {c.user.slice(1, 3).toUpperCase()}
                    </div>
                    <span className="text-sm font-semibold text-foreground">{c.user}</span>
                    <Badge className="bg-primary/10 text-primary border border-primary/30 rounded-sm uppercase text-[9px] tracking-[0.18em]">{c.badge}</Badge>
                    <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">{c.time}</span>
                  </div>
                  <p className="text-sm text-foreground/85 leading-relaxed">{c.body}</p>
                  <button className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground hover:text-primary">REPLY</button>
                </div>
              ))}
            </div>
          ))}
        </div>

        {/* Reply box */}
        <div className="panel p-5 space-y-3">
          <div className="flex items-center justify-between">
            <SectionLabel>Join_The_Discussion</SectionLabel>
            <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">AS @SYSTEM_ADMIN</span>
          </div>
          <Textarea
            placeholder="ENTER_TERMINAL_INPUT..."
            className="bg-terminal border-border min-h-[120px] font-mono text-sm placeholder:text-muted-foreground/50 focus-visible:ring-primary/40"
          />
          <div className="flex items-center justify-between">
            <div className="flex gap-2 text-muted-foreground">
              <button className="hover:text-primary"><ImageIcon className="h-4 w-4" /></button>
              <button className="hover:text-primary"><Code2 className="h-4 w-4" /></button>
              <button className="hover:text-primary font-bold text-sm w-4">B</button>
            </div>
            <Button className="bg-gradient-signal hover:opacity-90 text-primary-foreground font-bold uppercase tracking-[0.18em] text-xs rounded-sm">
              TRANSMIT_REPLY
            </Button>
          </div>
        </div>
      </div>
    </AppLayout>
  );
};

export default ThreadDetail;
