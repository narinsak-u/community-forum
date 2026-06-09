"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { EntryEditor } from "@/components/forge/EntryEditor";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { useRequireAuth } from "@/hooks/use-require-auth";
import { useCreateThread } from "@/hooks/use-thread";

export default function CreateEntryInner() {
  const router = useRouter();
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [tagsInput, setTagsInput] = useState("");
  const [viewMode, setViewMode] = useState<"split" | "full">("split");
  const { requireAuth } = useRequireAuth();

  useEffect(() => {
    requireAuth({ redirect: "/create" });
  }, [requireAuth]);

  const createThread = useCreateThread();

  const handleSaveDraft = () => {
    if (title.length < 5) {
      toast.error("VALIDATION_FAILED", {
        description: "Title must be at least 5 characters",
      });
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
        onError: (err) =>
          toast.error("SAVE_FAILED", { description: err.message }),
      },
    );
  };

  const handlePublish = () => {
    if (title.length < 5) {
      toast.error("VALIDATION_FAILED", {
        description: "Title must be at least 5 characters",
      });
      return;
    }
    if (title.length > 255) {
      toast.error("VALIDATION_FAILED", {
        description: "Title must be at most 255 characters",
      });
      return;
    }
    if (content.length < 10) {
      toast.error("VALIDATION_FAILED", {
        description: "Content must be at least 10 characters",
      });
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
          toast.success("ENTRY_TRANSMITTED", {
            description: "Your entry is now live on the nexus",
          });
          router.push("/thread/" + thread.slug);
        },
        onError: (err) =>
          toast.error("TRANSMIT_FAILED", { description: err.message }),
      },
    );
  };

  return (
    <div className="px-8 py-10 max-w-[1100px] mx-auto space-y-8 animate-fade-up">
      <header className="space-y-3">
        <SectionLabel>INITIALISING_PROTOCOL_V.2.4</SectionLabel>
        <h1 className="heading-display text-5xl text-foreground uppercase">
          CREATE NEW ENTRY
        </h1>
        <p className="text-sm text-muted-foreground max-w-2xl">
          Contribute to the collective intelligence. Ensure all documentation
          adheres to the Forge's architectural standards.
        </p>
      </header>

      <EntryEditor
        title={title}
        onTitleChange={setTitle}
        content={content}
        onContentChange={setContent}
        tagsInput={tagsInput}
        onTagsChange={setTagsInput}
        viewMode={viewMode}
        onViewModeChange={setViewMode}
        footer={
          <>
            <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
              {createThread.isPending
                ? "PROCESSING..."
                : "AUTOSAVE: ENABLED · READY"}
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
          </>
        }
      />
    </div>
  );
}
