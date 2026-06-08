"use client";

import { Skeleton } from "@/components/ui/skeleton";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { Button } from "@/components/ui/button";
import { ThumbsUp, ThumbsDown, Share2, Bookmark, Flag, Link2 } from "lucide-react";
import { toast } from "sonner";
import { useThread } from "@/hooks/use-thread";
import { useVoteThread } from "@/hooks/use-votes";
import { useRequireAuth } from "@/hooks/use-require-auth";
import { ThreadHeader } from "./ThreadHeader";
import { CommentSection } from "./CommentSection";

const QUOTE = "Precision is not just about speed; it's about predictable outcomes in a chaotic data environment.";

interface ThreadDetailClientProps {
  slug: string;
  initialThread?: any;
}

const ThreadDetailClient = ({ slug, initialThread }: ThreadDetailClientProps) => {
  const { data: thread, isLoading } = useThread(slug);
  const voteThread = useVoteThread(slug);
  const { requireAuth } = useRequireAuth();

  const currentThread = thread || initialThread;

  const author = currentThread?.author;
  const authorInitials = author?.username
    ? author.username.replace("@", "").slice(0, 2).toUpperCase()
    : "AL";

  const handleVote = (value: number) => {
    if (!requireAuth({ toast: true, description: "Authenticate to cast a signal." })) return;
    voteThread.mutate(
      { value },
      { onError: (err) => toast.error("VOTE_FAILED", { description: err.message }) },
    );
  };

  if (isLoading && !initialThread) {
    return (
      <div className="px-8 py-10 max-w-[1100px] mx-auto space-y-8">
        <Skeleton className="h-8 w-64" />
        <Skeleton className="h-16 w-full" />
        <Skeleton className="aspect-[2/1] w-full" />
        <Skeleton className="h-48 w-full" />
      </div>
    );
  }

  return (
    <div className="px-8 py-10 max-w-[1100px] mx-auto space-y-8 animate-fade-up">
      <ThreadHeader
        title={currentThread?.title}
        tagName={currentThread?.tags?.[0]?.name}
        authorUsername={author?.username}
        authorInitials={authorInitials}
        createdAt={currentThread?.created_at}
        viewCount={currentThread?.view_count || 0}
        repliesCount={currentThread?.replies_count || 0}
        upvotes={currentThread?.upvotes || 0}
        downvotes={currentThread?.downvotes || 0}
      />

      <div className="panel scanline relative overflow-hidden aspect-[2/1]">
        <img
          src="/images/forge-hero.jpg"
          alt="System wireframe"
          className="w-full h-full object-cover"
          width={1280}
          height={640}
        />
      </div>

      <div className="grid lg:grid-cols-[1fr,260px] gap-8">
        <article className="space-y-6 text-sm">
          {currentThread?.content ? (
            currentThread.content.split("\n\n").map((para: string, i: number) => (
              <p key={i} className="text-foreground/85 leading-relaxed">{para}</p>
            ))
          ) : (
            <p className="text-foreground/85 leading-relaxed">Loading content...</p>
          )}

          <blockquote className="border-l-2 border-primary pl-5 py-2 italic text-primary/90">
            &ldquo;{QUOTE}&rdquo;
          </blockquote>

          <div className="flex items-center justify-between pt-4 border-t border-border/60">
            <div className="flex items-center gap-3">
              <Button
                variant="outline"
                size="sm"
                className="border-border hover:border-primary hover:text-primary rounded-sm"
                onClick={() => handleVote(1)}
              >
                <ThumbsUp className="h-3.5 w-3.5 mr-2" /> {currentThread?.upvotes || 0}
              </Button>
              <Button
                variant="outline"
                size="sm"
                className="border-border rounded-sm"
                onClick={() => handleVote(-1)}
              >
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
              <div className="h-8 w-8 rounded-full bg-secondary border border-border grid place-items-center text-[10px] font-mono text-primary">
                +{currentThread?.replies_count || 0}
              </div>
            </div>
          </div>
        </aside>
      </div>

      <CommentSection
        slug={slug}
        comments={currentThread?.comments}
        repliesCount={currentThread?.replies_count || 0}
      />
    </div>
  );
};

export default ThreadDetailClient;
