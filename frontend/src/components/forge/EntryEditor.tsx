"use client";

import { useMemo } from "react";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Bold,
  Italic,
  Code2,
  List,
  Link2,
  Image as ImageIcon,
} from "lucide-react";

interface EntryEditorProps {
  title: string;
  onTitleChange: (value: string) => void;
  content: string;
  onContentChange: (value: string) => void;
  tagsInput: string;
  onTagsChange: (value: string) => void;
  viewMode: "split" | "full";
  onViewModeChange: (mode: "split" | "full") => void;
  footer?: React.ReactNode;
}

export function EntryEditor({
  title,
  onTitleChange,
  content,
  onContentChange,
  tagsInput,
  onTagsChange,
  viewMode,
  onViewModeChange,
  footer,
}: EntryEditorProps) {
  const previewTitle = useMemo(
    () => (title || "NEW_ENTRY").replace(/\s+/g, "_").toUpperCase(),
    [title],
  );

  return (
    <>
      <section className="space-y-3">
        <label htmlFor="title-input" className="terminal-label">
          TOPIC_IDENTITY
        </label>
        <Input
          id="title-input"
          value={title}
          onChange={(e) => onTitleChange(e.target.value)}
          placeholder="ENTER_UNIQUE_STRING_IDENTIFIER..."
          className="h-14 bg-secondary/60 border-border text-base font-mono uppercase tracking-wider placeholder:text-muted-foreground/40 focus-visible:ring-primary/40 rounded-sm"
        />
        <div className="flex items-center gap-2">
          <label
            htmlFor="tags-input"
            className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground"
          >
            TAGS (comma-separated):
          </label>
          <Input
            id="tags-input"
            value={tagsInput}
            onChange={(e) => onTagsChange(e.target.value)}
            placeholder="technical, announcement, design"
            className="h-8 bg-secondary/60 border-border text-xs font-mono tracking-wider placeholder:text-muted-foreground/40 focus-visible:ring-primary/40 rounded-sm w-96"
          />
        </div>
      </section>

      <section className="panel p-4 space-y-4">
        <div className="flex items-center justify-between border-b border-border/60 pb-3">
          <div className="flex items-center gap-1 text-muted-foreground">
            {[Bold, Italic, Code2, List, Link2, ImageIcon].map((Icon, i) => (
              <button
                key={i}
                type="button"
                className="h-8 w-8 grid place-items-center hover:text-primary hover:bg-secondary rounded-sm transition-colors"
              >
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
              <button
                type="button"
                onClick={() => onViewModeChange("split")}
                className={`px-3 py-1.5 transition-colors ${
                  viewMode === "split"
                    ? "bg-primary text-primary-foreground font-bold"
                    : "text-muted-foreground hover:bg-secondary"
                }`}
              >
                SPLIT
              </button>
              <button
                type="button"
                onClick={() => onViewModeChange("full")}
                className={`px-3 py-1.5 transition-colors ${
                  viewMode === "full"
                    ? "bg-primary text-primary-foreground font-bold"
                    : "text-muted-foreground hover:bg-secondary"
                }`}
              >
                FULL
              </button>
            </div>
          </div>
        </div>

        <div
          className={`grid ${viewMode === "split" ? "md:grid-cols-2" : "md:grid-cols-1"} gap-4 min-h-[400px]`}
        >
          <div className="bg-terminal border border-border/60 rounded-sm overflow-hidden">
            <div className="flex items-center justify-between px-3 py-2 border-b border-border/60 text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
              <span>SOURCE_STREAM</span>
              <span>INPUT</span>
            </div>
            <div className="grid grid-cols-[40px,1fr] font-mono text-sm">
              <div className="border-r border-border/60 py-3 text-right pr-2 text-muted-foreground/40 text-xs leading-relaxed select-none">
                {Array.from({ length: 16 }, (_, i) => (
                  <div key={i}>{String(i + 1).padStart(2, "0")}</div>
                ))}
              </div>
              <Textarea
                value={content}
                onChange={(e) => onContentChange(e.target.value)}
                placeholder="Begin composing your entry..."
                className="bg-transparent border-0 rounded-none min-h-[360px] focus-visible:ring-0 font-mono text-sm leading-relaxed text-foreground resize-none"
              />
            </div>
          </div>

          <div className="bg-terminal/40 border border-border/60 rounded-sm overflow-hidden">
            <div className="flex items-center justify-between px-3 py-2 border-b border-border/60 text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
              <span>RENDER_PREVIEW</span>
              <span>OUTPUT</span>
            </div>
            <div className="p-5 space-y-3 text-sm overflow-auto max-h-[400px]">
              <h2 className="heading-display text-xl text-foreground">
                {previewTitle}
              </h2>
              {content ? (
                content.split("\n").map((line, i) => {
                  if (line.startsWith("# "))
                    return (
                      <h2
                        key={i}
                        className="heading-display text-xl text-foreground"
                      >
                        {line.slice(2)}
                      </h2>
                    );
                  if (line.startsWith("## "))
                    return (
                      <h3
                        key={i}
                        className="text-primary font-mono text-xs uppercase tracking-[0.18em]"
                      >
                        {line.slice(3)}
                      </h3>
                    );
                  if (line.startsWith("### "))
                    return (
                      <h4
                        key={i}
                        className="text-foreground font-semibold pt-2"
                      >
                        {line.slice(4)}
                      </h4>
                    );
                  if (line.startsWith("- "))
                    return (
                      <li
                        key={i}
                        className="text-foreground/85 list-disc list-inside marker:text-primary"
                      >
                        {line.slice(2)}
                      </li>
                    );
                  return (
                    <p key={i} className="text-foreground/85 leading-relaxed">
                      {line}
                    </p>
                  );
                })
              ) : (
                <p className="text-foreground/50 italic">
                  Begin typing in the source panel to see a preview...
                </p>
              )}
            </div>
          </div>
        </div>
      </section>

      {footer && (
        <div className="flex items-center justify-between pt-2">{footer}</div>
      )}
    </>
  );
}
