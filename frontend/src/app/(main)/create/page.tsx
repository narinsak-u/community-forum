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

const CreateEntry = () => {
  const router = useRouter();
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [tagsInput, setTagsInput] = useState("");

  const createThread = useCreateThread();

  const handleSaveDraft = () => {
    if (title.length < 5) {
      toast.error("VALIDATION_FAILED", { description: "Title must be at least 5 characters" });
      return;
    }
    createThread.mutate(
      {
        title,
        content: content || "# New Entry\n\nStart writing here...",
        tags: tagsInput
          .split(",")
          .map((t) => t.trim())
          .filter(Boolean),
        status: "draft",
      },
      {
        onSuccess: (thread) => {
          toast.success("DRAFT_SAVED", { description: "Entry saved as draft" });
          router.push("/thread/" + thread.slug);
        },
        onError: (err) => toast.error("SAVE_FAILED", { description: err.message }),
      },
    );
  };

  const handlePublish = () => {
    if (title.length < 5) {
      toast.error("VALIDATION_FAILED", { description: "Title must be at least 5 characters" });
      return;
    }
    if (title.length > 255) {
      toast.error("VALIDATION_FAILED", { description: "Title must be at most 255 characters" });
      return;
    }
    if (content.length < 10) {
      toast.error("VALIDATION_FAILED", { description: "Content must be at least 10 characters" });
      return;
    }
    createThread.mutate(
      {
        title,
        content,
        tags: tagsInput
          .split(",")
          .map((t) => t.trim())
          .filter(Boolean),
        status: "published",
      },
      {
        onSuccess: (thread) => {
          toast.success("ENTRY_TRANSMITTED", { description: "Your entry is now live on the nexus" });
          router.push("/thread/" + thread.slug);
        },
        onError: (err) => toast.error("TRANSMIT_FAILED", { description: err.message }),
      },
    );
  };

  const previewTitle = title || "NEW_ENTRY";

  return (
    <div className="px-8 py-10 max-w-[1100px] mx-auto space-y-8 animate-fade-up">
      <header className="space-y-3">
        <SectionLabel>INITIALISING_PROTOCOL_V.2.4</SectionLabel>
        <h1 className="heading-display text-5xl text-foreground uppercase">CREATE NEW ENTRY</h1>
        <p className="text-sm text-muted-foreground max-w-2xl">
          Contribute to the collective intelligence. Ensure all documentation adheres to the Forge's architectural standards.
        </p>
      </header>

      <section className="space-y-3">
        <label className="terminal-label">TOPIC_IDENTITY</label>
        <Input
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="ENTER_UNIQUE_STRING_IDENTIFIER..."
          className="h-14 bg-secondary/60 border-border text-base font-mono uppercase tracking-wider placeholder:text-muted-foreground/40 focus-visible:ring-primary/40 rounded-sm"
        />
        <div className="flex items-center gap-2">
          <label className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">TAGS (comma-separated):</label>
          <Input
            value={tagsInput}
            onChange={(e) => setTagsInput(e.target.value)}
            placeholder="technical, announcement, design"
            className="h-8 bg-secondary/60 border-border text-xs font-mono tracking-wider placeholder:text-muted-foreground/40 focus-visible:ring-primary/40 rounded-sm w-64"
          />
        </div>
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
                value={content}
                onChange={(e) => setContent(e.target.value)}
                placeholder="Begin composing your entry..."
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
            <div className="p-5 space-y-3 text-sm overflow-auto max-h-[400px]">
              <h2 className="heading-display text-xl text-foreground">
                {previewTitle.replace(/\s+/g, "_").toUpperCase()}
              </h2>
              {content ? (
                content.split("\n").map((line, i) => {
                  if (line.startsWith("# ")) return <h2 key={i} className="heading-display text-xl text-foreground">{line.slice(2)}</h2>;
                  if (line.startsWith("## ")) return <h3 key={i} className="text-primary font-mono text-xs uppercase tracking-[0.18em]">{line.slice(3)}</h3>;
                  if (line.startsWith("### ")) return <h4 key={i} className="text-foreground font-semibold pt-2">{line.slice(4)}</h4>;
                  if (line.startsWith("- ")) return <li key={i} className="text-foreground/85 list-disc list-inside marker:text-primary">{line.slice(2)}</li>;
                  return <p key={i} className="text-foreground/85 leading-relaxed">{line}</p>;
                })
              ) : (
                <p className="text-foreground/50 italic">Begin typing in the source panel to see a preview...</p>
              )}
            </div>
          </div>
        </div>
      </section>

      <div className="flex items-center justify-between pt-2">
        <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
          {createThread.isPending ? "PROCESSING..." : "AUTOSAVE: ENABLED · READY"}
        </span>
        <div className="flex gap-3">
          <Button
            variant="outline"
            onClick={handleSaveDraft}
            disabled={createThread.isPending}
            className="border-border text-muted-foreground hover:text-foreground rounded-sm uppercase tracking-[0.18em] text-xs"
          >
            SAVE_DRAFT
          </Button>
          <Button
            onClick={handlePublish}
            disabled={createThread.isPending}
            className="bg-gradient-signal hover:opacity-90 text-primary-foreground font-bold uppercase tracking-[0.18em] text-xs rounded-sm"
          >
            {createThread.isPending ? "TRANSMITTING..." : "TRANSMIT_ENTRY"}
          </Button>
        </div>
      </div>
    </div>
  );
};

export default CreateEntry;