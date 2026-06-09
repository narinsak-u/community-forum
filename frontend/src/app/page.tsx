import Link from "next/link";
import { ArrowRight, LayoutGrid, MessageSquare, BookOpen } from "lucide-react";

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-background relative overflow-hidden">
      <video
        autoPlay
        loop
        muted
        playsInline
        className="absolute inset-0 w-full h-full object-cover"
      >
        <source src="/images/video-bg.mp4" type="video/mp4" />
      </video>
      <div className="relative z-10 max-w-[1200px] mx-auto px-6 py-20 md:py-32 space-y-20">
        <section className="space-y-6 animate-fade-up text-center">
          <div className="inline-flex items-center gap-2 px-3 py-1 border border-primary/30 bg-primary/5 rounded-sm text-[10px] uppercase tracking-[0.2em] text-primary font-mono">
            <span className="h-1.5 w-1.5 rounded-full bg-primary animate-pulse" />
            SYSTEM ONLINE // v1.0.2
          </div>

          <h1 className="heading-display text-white text-5xl md:text-7xl lg:text-8xl text-foreground leading-tight">
            The Lands <span className="text-primary">Between</span>
          </h1>

          <p className="text-sm md:text-base text-muted-foreground max-w-2xl mx-auto font-mono leading-relaxed">
            A technical forum for architects, engineers, and builders of the
            digital frontier. Share knowledge, discuss systems, and forge the
            future of technology.
          </p>

          <div className="flex items-center justify-center gap-4 pt-4">
            <Link
              href="/threads"
              className="inline-flex items-center gap-2 px-6 py-3 bg-gradient-signal hover:opacity-90 text-primary-foreground font-bold uppercase tracking-[0.15em] text-xs rounded-sm transition-all"
            >
              ENTER THE TERMINAL
              <ArrowRight className="h-4 w-4" />
            </Link>
            {/*<Link
              href="/login?redirect=/threads"
              className="inline-flex items-center gap-2 px-6 py-3 border border-border hover:border-primary/60 text-foreground font-mono text-xs uppercase tracking-[0.15em] rounded-sm transition-colors"
            >
              SIGN_IN
            </Link>*/}
          </div>
        </section>

        {/*<section className="grid md:grid-cols-3 gap-6 max-w-4xl mx-auto">
          {[
            {
              icon: LayoutGrid,
              title: "Discussions",
              desc: "Explore technical threads on architecture, systems, and engineering.",
            },
            {
              icon: MessageSquare,
              title: "Collaborate",
              desc: "Share insights, ask questions, and contribute to the collective knowledge.",
            },
            {
              icon: BookOpen,
              title: "Documentation",
              desc: "Access technical references, protocols, and system documentation.",
            },
          ].map((feature) => (
            <div
              key={feature.title}
              className="panel p-6 space-y-3 hover:border-primary/40 transition-colors"
            >
              <feature.icon className="h-5 w-5 text-primary" />
              <h3 className="text-sm font-semibold text-foreground uppercase tracking-[0.12em]">
                {feature.title}
              </h3>
              <p className="text-xs text-muted-foreground leading-relaxed">
                {feature.desc}
              </p>
            </div>
          ))}
        </section>*/}

        <footer className="text-center space-y-4 text-white">
          <div className="flex justify-center gap-6 text-[10px] uppercase tracking-[0.2em]  font-mono">
            <span>Protocol</span>
            <span>Manifesto</span>
            <span>Support</span>
          </div>
          <p className="text-[10px] uppercase tracking-[0.2em] font-mono">
            &copy; 2026 The Lands Between // Encrypted Session
          </p>
        </footer>
      </div>
    </div>
  );
}
