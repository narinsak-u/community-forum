"use client";

import { useEffect, useState, use } from "react";
import { useRouter } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";
import { useThread, useUpdateThread } from "@/hooks/use-thread";
import { useRequireAuth } from "@/hooks/use-require-auth";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { EntryEditor } from "@/components/forge/EntryEditor";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

export default function EditEntryPage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = use(params);

  return <EditEntryInner slug={slug} />;
}

function EditEntryInner({ slug }: { slug: string }) {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { data: thread, isLoading } = useThread(slug);
  const updateThread = useUpdateThread(slug);
  const { requireAuth } = useRequireAuth();

  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [tagsInput, setTagsInput] = useState("");
  const [viewMode, setViewMode] = useState<"split" | "full">("split");

  useEffect(() => {
    requireAuth({ redirect: `/thread/${slug}` });
  }, [requireAuth, slug]);

  useEffect(() => {
    if (thread) {
      setTitle(thread.title);
      setContent(thread.content);
      setTagsInput(thread.tags?.map((t) => t.name).join(", ") || "");
    }
  }, [thread]);

  const isDraft = thread?.status === "draft";

  const validate = () => {
    if (title.length < 5) {
      toast.error("VALIDATION_FAILED", {
        description: "Title must be at least 5 characters",
      });
      return false;
    }
    return true;
  };

  const handleSaveDraft = () => {
    if (!validate()) return;
    updateThread.mutate(
      {
        title,
        content,
        tags: tagsInput
          .split(",")
          .map((t) => t.trim())
          .filter(Boolean),
        status: "draft",
      },
      {
        onSuccess: (updated) => {
          toast.success("DRAFT_SAVED", { description: "Entry saved as draft" });
          queryClient.invalidateQueries({ queryKey: ["users"] });
          router.push(`/thread/${updated.slug}/edit`);
        },
        onError: (err) =>
          toast.error("SAVE_FAILED", { description: err.message }),
      },
    );
  };

  const handlePublish = () => {
    if (!validate()) return;
    if (content.length < 10) {
      toast.error("VALIDATION_FAILED", {
        description: "Content must be at least 10 characters",
      });
      return;
    }
    updateThread.mutate(
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
        onSuccess: (updated) => {
          toast.success("ENTRY_TRANSMITTED", {
            description: "Your entry is now live on the nexus",
          });
          queryClient.invalidateQueries({ queryKey: ["users"] });
          router.push(`/thread/${updated.slug}`);
        },
        onError: (err) =>
          toast.error("TRANSMIT_FAILED", { description: err.message }),
      },
    );
  };

  const handleSaveChanges = () => {
    if (!validate()) return;
    updateThread.mutate(
      {
        title,
        content,
        tags: tagsInput
          .split(",")
          .map((t) => t.trim())
          .filter(Boolean),
      },
      {
        onSuccess: (updated) => {
          toast.success("ENTRY_UPDATED", {
            description: "Your changes have been saved",
          });
          queryClient.invalidateQueries({ queryKey: ["users"] });
          router.push(`/thread/${updated.slug}`);
        },
        onError: (err) =>
          toast.error("UPDATE_FAILED", { description: err.message }),
      },
    );
  };

  if (isLoading) {
    return (
      <div className="px-8 py-10 max-w-[1100px] mx-auto space-y-8">
        <div className="h-5 w-48 bg-secondary/60 animate-pulse rounded-sm" />
        <div className="h-12 w-full bg-secondary/60 animate-pulse rounded-sm" />
        <div className="h-[400px] bg-secondary/60 animate-pulse rounded-sm" />
      </div>
    );
  }

  return (
    <div className="px-8 py-10 max-w-[1100px] mx-auto space-y-8 animate-fade-up">
      <header className="space-y-3">
        <SectionLabel>MODIFYING_PROTOCOL_V.2.4</SectionLabel>
        <h1 className="heading-display text-5xl text-foreground uppercase">
          EDIT ENTRY
        </h1>
        <p className="text-sm text-muted-foreground max-w-2xl">
          Revise your contribution to the collective intelligence.
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
              {updateThread.isPending
                ? "PROCESSING..."
                : "AUTOSAVE: ENABLED · READY"}
            </span>
            <div className="flex gap-3">
              <Button
                size="sm"
                variant="outline"
                onClick={() => router.push(isDraft ? "/threads" : `/thread/${slug}`)}
                className="border-border text-muted-foreground hover:text-foreground rounded-sm uppercase tracking-[0.18em] text-xs"
              >
                CANCEL
              </Button>
              <Button
                size="sm"
                variant="outline"
                onClick={handleSaveDraft}
                disabled={updateThread.isPending}
                className="border-amber-500/30 text-amber-400 hover:text-amber-300 rounded-sm uppercase tracking-[0.18em] text-xs border"
              >
                {updateThread.isPending ? "SAVING..." : "DRAFT"}
              </Button>
              <Button
                size="sm"
                onClick={isDraft ? handlePublish : handleSaveChanges}
                disabled={updateThread.isPending}
                className="bg-gradient-signal hover:opacity-90 text-primary-foreground font-bold uppercase tracking-[0.18em] text-xs rounded-sm"
              >
                {updateThread.isPending
                  ? "PROCESSING..."
                  : isDraft
                    ? "TRANSMIT_ENTRY"
                    : "SAVE_CHANGES"}
              </Button>
            </div>
          </>
        }
      />
    </div>
  );
}
