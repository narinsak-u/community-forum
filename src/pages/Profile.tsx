import { AppLayout } from "@/components/forge/AppLayout";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { NavLink } from "@/components/NavLink";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Fingerprint, Terminal, Shield, Archive, FileText, FileText as FileIcon, Globe, Lock, Plus } from "lucide-react";
import avatarImg from "@/assets/forge-avatar.jpg";

const navItems = [
  { label: "Identity", icon: Fingerprint, active: true },
  { label: "Contributions", icon: Terminal },
  { label: "Trust Metrics", icon: Shield },
  { label: "Archived Vault", icon: Archive },
  { label: "System Logs", icon: FileText },
];

const vaultItems = [
  { icon: FileIcon, tag: "HARDWARE_LOGIC", title: "Neural Interface Schema v2", desc: "Complete architectural breakdown of the quantum-resistant neural gateway used..." },
  { icon: Globe, tag: "NETWORK_MAP", title: "Deep-Web Relay Topography", desc: "Visualizing the hidden relays within the forge-network to ensure 99.9% protocol..." },
  { icon: Lock, tag: "SECURITY_AUDIT", title: "Vault Entropy Benchmarks", desc: "Comparative analysis of random seed generators across the four vault sectors." },
];

const Profile = () => {
  return (
    <AppLayout showSidebar={false}>
      <div className="grid lg:grid-cols-[260px,1fr] gap-0 min-h-[calc(100vh-4rem)]">
        {/* Side rail */}
        <aside className="border-r border-border/60 bg-sidebar/40 p-6 space-y-6">
          <div className="flex items-center gap-3 p-2">
            <div className="h-10 w-10 bg-secondary border border-border rounded-sm grid place-items-center text-primary text-xs">
              <Fingerprint className="h-4 w-4" />
            </div>
            <div>
              <div className="text-sm font-bold">NODE_8829</div>
              <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">Senior Architect</div>
            </div>
          </div>
          <nav className="space-y-1">
            {navItems.map((item) => (
              <button
                key={item.label}
                className={`w-full flex items-center gap-3 px-3 py-2.5 text-sm rounded-sm transition-colors ${
                  item.active ? "bg-secondary text-primary border-l-2 border-primary" : "text-muted-foreground hover:bg-secondary/50 hover:text-foreground"
                }`}
              >
                <item.icon className="h-4 w-4" />
                {item.label}
              </button>
            ))}
          </nav>
          <Button className="w-full bg-gradient-signal hover:opacity-90 text-primary-foreground font-bold uppercase tracking-[0.18em] text-xs h-10 rounded-sm">
            INITIATE_UPLINK
          </Button>
        </aside>

        {/* Content */}
        <div className="p-8 lg:p-10 space-y-10 animate-fade-up">
          {/* Header */}
          <section className="grid md:grid-cols-[120px,1fr,auto] gap-6 items-start">
            <div className="h-28 w-28 panel overflow-hidden">
              <img src={avatarImg} alt="Node avatar" className="w-full h-full object-cover" width={512} height={512} />
            </div>
            <div className="space-y-3">
              <h1 className="heading-display text-5xl text-foreground">NODE_8829</h1>
              <div className="text-xs uppercase tracking-[0.2em] text-primary font-mono">
                SENIOR ARCHITECT // SYSTEMS INTEGRATOR
              </div>
              <div className="flex flex-wrap gap-2 pt-1">
                <Badge className="bg-primary/10 text-primary border border-primary/30 rounded-sm uppercase text-[10px] tracking-[0.18em]">● CORE_CONTRIBUTOR</Badge>
                <Badge className="bg-primary/10 text-primary border border-primary/30 rounded-sm uppercase text-[10px] tracking-[0.18em]">● SECURITY_AUDITOR</Badge>
              </div>
            </div>
            <div className="panel p-5 grid grid-cols-3 gap-4 min-w-[300px]">
              {[
                { v: "1.2k", l: "CONTRIBUTIONS" },
                { v: "98%", l: "TRUST INDEX" },
                { v: "L_04", l: "ACCESS LEVEL" },
              ].map((s) => (
                <div key={s.l}>
                  <div className="text-2xl font-display font-bold text-foreground">{s.v}</div>
                  <div className="text-[9px] uppercase tracking-[0.18em] text-muted-foreground">{s.l}</div>
                </div>
              ))}
            </div>
          </section>

          {/* Two-column body */}
          <section className="grid lg:grid-cols-[1fr,1.4fr] gap-10">
            {/* Identity column */}
            <div className="space-y-8">
              <div className="space-y-4">
                <SectionLabel>SYSTEM IDENTITY</SectionLabel>
                <p className="text-sm text-foreground/85 leading-relaxed">
                  Specializing in decentralized infrastructure and encrypted communication protocols. Former lead dev for Project Wraith. Currently optimizing neural-link throughput for the ARCHITECT_FORUM core.
                </p>
              </div>

              <div className="panel p-5 space-y-3">
                <div className="text-[10px] uppercase tracking-[0.18em] text-primary">TECHNICAL_STACK</div>
                <div className="flex flex-wrap gap-2">
                  {["Rust", "Solidity", "Go", "Post-Quantum Cryptography"].map((t) => (
                    <span key={t} className="px-2.5 py-1 bg-secondary border border-border text-xs text-foreground rounded-sm">{t}</span>
                  ))}
                </div>
              </div>

              <div className="space-y-4">
                <SectionLabel>CONTRIBUTION STREAM</SectionLabel>
                {[
                  { title: "Fragmented Node Recovery Protocol", time: "2H_AGO", tags: ["PROTOCOL", "CRITICAL"] },
                  { title: "Vault Security Patch v8.4.2", time: "1D_AGO", tags: ["SECURITY", "VAULT"] },
                ].map((c) => (
                  <div key={c.title} className="flex justify-between items-start gap-4 border-b border-border/60 pb-3">
                    <div className="space-y-1.5">
                      <div className="text-sm font-semibold text-foreground">{c.title}</div>
                      <div className="flex gap-1.5">
                        {c.tags.map((tg) => (
                          <span key={tg} className="px-1.5 py-0.5 bg-secondary border border-border text-[9px] uppercase tracking-[0.18em] text-muted-foreground rounded-sm">{tg}</span>
                        ))}
                      </div>
                    </div>
                    <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground shrink-0">{c.time}</span>
                  </div>
                ))}
              </div>
            </div>

            {/* Vault grid */}
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <SectionLabel>DIGITAL VAULT</SectionLabel>
                <button className="text-[10px] uppercase tracking-[0.18em] text-primary hover:text-primary-glow">ACCESS_ALL_ARCHIVES</button>
              </div>
              <div className="grid sm:grid-cols-2 gap-4">
                {vaultItems.map((v) => (
                  <div key={v.title} className="panel p-5 space-y-3 hover:border-primary/40 transition-colors group cursor-pointer">
                    <div className="aspect-[3/2] bg-terminal/80 border border-border/60 rounded-sm relative overflow-hidden">
                      <div className="absolute inset-0 grid place-items-center text-primary/30 group-hover:text-primary/60 transition-colors">
                        <v.icon className="h-12 w-12" />
                      </div>
                      <div
                        className="absolute inset-0 opacity-30"
                        style={{
                          backgroundImage: 'radial-gradient(circle at 30% 30%, hsl(var(--primary) / 0.4), transparent 60%)',
                        }}
                      />
                    </div>
                    <div className="text-[10px] uppercase tracking-[0.2em] text-primary">{v.tag}</div>
                    <div className="text-base font-bold text-foreground">{v.title}</div>
                    <div className="text-xs text-muted-foreground line-clamp-2">{v.desc}</div>
                  </div>
                ))}
                <button className="panel p-5 grid place-items-center min-h-[280px] border-dashed hover:border-primary/40 transition-colors group">
                  <div className="text-center space-y-2">
                    <div className="h-12 w-12 mx-auto bg-secondary border border-border rounded-sm grid place-items-center text-muted-foreground group-hover:text-primary group-hover:border-primary/40 transition-colors">
                      <Plus className="h-5 w-5" />
                    </div>
                    <div className="text-[10px] uppercase tracking-[0.2em] text-muted-foreground">ADD_TO_VAULT</div>
                  </div>
                </button>
              </div>
            </div>
          </section>
        </div>
      </div>
    </AppLayout>
  );
};

export default Profile;
