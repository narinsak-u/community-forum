import { AppLayout } from "@/components/forge/AppLayout";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Bold, Italic, Code2, List, Link2, Image as ImageIcon } from "lucide-react";

const sampleSource = `# DOCUMENTATION_LOG_ENTRY_402

## INFRASTRUCTURE_UPGRADE_OVERVIEW

Following the recent sector shift, all **Core-Link** protocols have been updated to v.2.4. Developers are required to audit their enclaves.

### Key Changes:
- Latency reduced by 14ms
- New encryption layer (AES-256-FORGE)
- Extended thermal headroom`;

const CreateEntry = () => {
  return (
    <AppLayout showSidebar showNewEntry>
      <div className="px-8 py-10 max-w-[1100px] mx-auto space-y-8 animate-fade-up">
        <header className="space-y-3">
          <SectionLabel>INITIALISING_PROTOCOL_V.2.4</SectionLabel>
          <h1 className="heading-display text-5xl text-foreground uppercase">CREATE NEW ENTRY</h1>
          <p className="text-sm text-muted-foreground max-w-2xl">
            Contribute to the collective intelligence. Ensure all technical documentation adheres to the Forge's architectural standards.
          </p>
        </header>

        <section className="space-y-3">
          <label className="terminal-label">TOPIC_IDENTITY</label>
          <Input
            placeholder="ENTER_UNIQUE_STRING_IDENTIFIER..."
            className="h-14 bg-secondary/60 border-border text-base font-mono uppercase tracking-wider placeholder:text-muted-foreground/40 focus-visible:ring-primary/40 rounded-sm"
          />
        </section>

        <section className="panel p-4 space-y-4">
          {/* Toolbar */}
          <div className="flex items-center justify-between border-b border-border/60 pb-3">
            <div className="flex items-center gap-1 text-muted-foreground">
              {[Bold, Italic, Code2, List, Link2, ImageIcon].map((Icon, i) => (
                <button key={i} className="h-8 w-8 grid place-items-center hover:text-primary hover:bg-secondary rounded-sm transition-colors">
                  <Icon className="h-4 w-4" />
                </button>
              ))}
            </div>
            <div className="flex items-center gap-3">
              <span className="flex items-center gap-2 text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                <span className="h-1.5 w-1.5 rounded-full bg-primary animate-pulse" />
                DRAFT_LIVE
              </span>
              <div className="flex border border-border rounded-sm overflow-hidden text-[10px] uppercase tracking-[0.18em]">
                <button className="px-3 py-1.5 bg-primary text-primary-foreground font-bold">SPLIT</button>
                <button className="px-3 py-1.5 text-muted-foreground hover:bg-secondary">FULL</button>
              </div>
            </div>
          </div>

          {/* Editor split */}
          <div className="grid md:grid-cols-2 gap-4 min-h-[400px]">
            {/* Source */}
            <div className="bg-terminal border border-border/60 rounded-sm overflow-hidden">
              <div className="flex items-center justify-between px-3 py-2 border-b border-border/60 text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                <span>SOURCE_STREAM</span>
                <span>MD_ENG_V1</span>
              </div>
              <div className="grid grid-cols-[40px,1fr] font-mono text-sm">
                <div className="border-r border-border/60 py-3 text-right pr-2 text-muted-foreground/40 text-xs leading-relaxed select-none">
                  {Array.from({ length: 16 }, (_, i) => <div key={i}>{String(i + 1).padStart(2, "0")}</div>)}
                </div>
                <Textarea
                  defaultValue={sampleSource}
                  className="bg-transparent border-0 rounded-none min-h-[360px] focus-visible:ring-0 font-mono text-sm leading-relaxed text-foreground resize-none"
                />
              </div>
            </div>

            {/* Preview */}
            <div className="bg-terminal/40 border border-border/60 rounded-sm overflow-hidden">
              <div className="flex items-center justify-between px-3 py-2 border-b border-border/60 text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                <span>RENDER_PREVIEW</span>
                <span>OUTPUT</span>
              </div>
              <div className="p-5 space-y-3 text-sm">
                <h2 className="heading-display text-xl text-foreground">DOCUMENTATION_LOG_ENTRY_402</h2>
                <h3 className="text-primary font-mono text-xs uppercase tracking-[0.18em]">INFRASTRUCTURE_UPGRADE_OVERVIEW</h3>
                <p className="text-foreground/85 leading-relaxed">
                  Following the recent sector shift, all <strong className="text-primary">Core-Link</strong> protocols have been updated to v.2.4. Developers are required to audit their enclaves.
                </p>
                <h4 className="text-foreground font-semibold pt-2">Key Changes:</h4>
                <ul className="space-y-1 text-foreground/85 list-disc list-inside marker:text-primary">
                  <li>Latency reduced by 14ms</li>
                  <li>New encryption layer (AES-256-FORGE)</li>
                  <li>Extended thermal headroom</li>
                </ul>
              </div>
            </div>
          </div>
        </section>

        <div className="flex items-center justify-between pt-2">
          <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">AUTOSAVE: ENABLED · LAST_SYNC 2s AGO</span>
          <div className="flex gap-3">
            <Button variant="outline" className="border-border text-muted-foreground hover:text-foreground rounded-sm uppercase tracking-[0.18em] text-xs">
              SAVE_DRAFT
            </Button>
            <Button className="bg-gradient-signal hover:opacity-90 text-primary-foreground font-bold uppercase tracking-[0.18em] text-xs rounded-sm">
              TRANSMIT_ENTRY
            </Button>
          </div>
        </div>
      </div>
    </AppLayout>
  );
};

export default CreateEntry;
