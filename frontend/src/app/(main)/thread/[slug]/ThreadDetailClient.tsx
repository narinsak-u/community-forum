"use client";

import { useState } from "react";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Skeleton } from "@/components/ui/skeleton";
import {
  ThumbsUp,
  ThumbsDown,
  Share2,
  Bookmark,
  Flag,
  Code2,
  Link2,
  Image as ImageIcon,
} from "lucide-react";
import { toast } from "sonner";
import { useThread } from "@/hooks/use-thread";
import { useCreateComment } from "@/hooks/use-comments";
import { useVoteThread } from "@/hooks/use-votes";
import { useAuthStore } from "@/stores/auth-store";
import { useRequireAuth } from "@/hooks/use-require-auth";

interface ThreadDetailClientProps {
  slug: string;
  initialThread?: any;
}

const ThreadDetailClient = ({
  slug,
  initialThread,
}: ThreadDetailClientProps) => {
  const [replyContent, setReplyContent] = useState("");
  const [replyTarget, setReplyTarget] = useState<{ id: number; author: string } | null>(null);
  const user = useAuthStore((s) => s.user);
  const { requireAuth } = useRequireAuth();

  const { data: thread, isLoading } = useThread(slug);
  const createComment = useCreateComment(slug);
  const voteThread = useVoteThread(slug);

  const currentThread = thread || initialThread;

  const author = currentThread?.author;
  const authorInitials = author?.username
    ? author.username.replace("@", "").slice(0, 2).toUpperCase()
    : "AL";

  const handleVote = (value: number) => {
    if (!requireAuth({ toast: true, description: "Authenticate to cast a signal." })) return;
    voteThread.mutate(
      { value },
      {
        onError: (err) =>
          toast.error("VOTE_FAILED", { description: err.message }),
      },
    );
  };

  const handleSubmitReply = () => {
    if (!replyContent.trim()) return;
    if (!requireAuth({ toast: true, description: "Authenticate to transmit a reply." })) return;
    createComment.mutate(
      { content: replyContent },
      {
        onSuccess: () => {
          setReplyContent("");
          toast.success("REPLY_TRANSMITTED");
        },
        onError: (err) =>
          toast.error("TRANSMIT_FAILED", { description: err.message }),
      },
    );
  };

  const formatCount = (n: number) => {
    if (n >= 1000) return (n / 1000).toFixed(1) + "K";
    return String(n);
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
      {/* Breadcrumb */}
      <nav className="flex items-center gap-2 text-[10px] uppercase tracking-[0.2em] text-muted-foreground">
        <span>Discussion</span>
        <span>/</span>
        <span>{currentThread?.tags?.[0]?.name || "Technical"}</span>
        <span>/</span>
        <span className="text-primary">
          {currentThread?.title?.slice(0, 24) || "Thread"}...
        </span>
      </nav>

      {/* Title */}
      <header className="space-y-4">
        {currentThread?.tags?.[0] && (
          <Badge className="bg-primary/15 text-primary border border-primary/40 rounded-sm uppercase text-[10px] tracking-[0.2em]">
            ● {currentThread.tags[0].name}
          </Badge>
        )}
        <h1 className="heading-display text-3xl md:text-5xl leading-tight text-foreground">
          {currentThread?.title || "Thread"}
        </h1>
        <div className="flex flex-wrap items-center justify-between gap-4 pt-2">
          <div className="flex items-center gap-3">
            <div className="h-10 w-10 rounded-full bg-gradient-signal grid place-items-center text-xs font-bold text-primary-foreground">
              {authorInitials}
            </div>
            <div>
              <div className="text-sm font-semibold text-foreground">
                {author?.username || "@unknown"}
              </div>
              <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                {currentThread?.created_at
                  ? new Date(currentThread.created_at).toLocaleDateString()
                  : ""}
              </div>
            </div>
          </div>
          <div className="flex gap-8 text-right">
            <div>
              <div className="text-lg font-display font-bold text-foreground">
                {formatCount(currentThread?.view_count || 0)}
              </div>
              <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                VIEWS
              </div>
            </div>
            <div>
              <div className="text-lg font-display font-bold text-foreground">
                {currentThread?.replies_count || 0}
              </div>
              <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                REPLIES
              </div>
            </div>
            <div>
              <div className="text-lg font-display font-bold text-primary">
                {currentThread?.upvotes != null &&
                currentThread.upvotes + currentThread.downvotes > 0
                  ? Math.round(
                      (currentThread.upvotes /
                        (currentThread.upvotes + currentThread.downvotes)) *
                        100,
                    ) + "%"
                  : "98%"}
              </div>
              <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                TRUST_INDEX
              </div>
            </div>
          </div>
        </div>
      </header>

      {/* Hero image */}
      <div className="panel scanline relative overflow-hidden aspect-[2/1]">
        <img
          src="/images/forge-hero.jpg"
          alt="System wireframe"
          className="w-full h-full object-cover"
          width={1280}
          height={640}
        />
      </div>

      {/* Body + sidebar */}
      <div className="grid lg:grid-cols-[1fr,260px] gap-8">
        <article className="space-y-6 text-sm">
          {currentThread?.content ? (
            currentThread.content.split("\n\n").map((para, i) => (
              <p key={i} className="text-foreground/85 leading-relaxed">
                {para}
              </p>
            ))
          ) : (
            <p className="text-foreground/85 leading-relaxed">
              Loading content...
            </p>
          )}

          <blockquote className="border-l-2 border-primary pl-5 py-2 italic text-primary/90">
            "Precision is not just about speed; it's about predictable outcomes
            in a chaotic data environment."
          </blockquote>

          {/* Vote bar */}
          <div className="flex items-center justify-between pt-4 border-t border-border/60">
            <div className="flex items-center gap-3">
              <Button
                variant="outline"
                size="sm"
                className="border-border hover:border-primary hover:text-primary rounded-sm"
                onClick={() => handleVote(1)}
              >
                <ThumbsUp className="h-3.5 w-3.5 mr-2" />{" "}
                {currentThread?.upvotes || 0}
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
              <button className="flex items-center gap-1.5 hover:text-primary">
                <Share2 className="h-3.5 w-3.5" /> Share
              </button>
              <button className="flex items-center gap-1.5 hover:text-primary">
                <Bookmark className="h-3.5 w-3.5" /> Bookmark
              </button>
              <button className="flex items-center gap-1.5 hover:text-primary">
                <Flag className="h-3.5 w-3.5" /> Report
              </button>
            </div>
          </div>
        </article>

        {/* References sidebar */}
        <aside className="space-y-6">
          <div className="panel p-5 space-y-4">
            <SectionLabel>References</SectionLabel>
            <ul className="space-y-2 text-sm">
              {["REF_CORE_041", "REACTIVE_HOOKS.DOC", "V2.3_LEGACY_PATCH"].map(
                (r) => (
                  <li
                    key={r}
                    className="flex items-center gap-2 text-primary hover:text-primary-glow cursor-pointer font-mono"
                  >
                    <Link2 className="h-3 w-3" /> {r}
                  </li>
                ),
              )}
            </ul>
          </div>
          <div className="panel p-5 space-y-3">
            <SectionLabel>Participants</SectionLabel>
            <div className="flex flex-wrap gap-1">
              {[...Array(5)].map((_, i) => (
                <div
                  key={i}
                  className="h-8 w-8 rounded-full bg-gradient-to-br from-primary/40 to-primary-deep/30 border border-border"
                />
              ))}
              <div className="h-8 w-8 rounded-full bg-secondary border border-border grid place-items-center text-[10px] font-mono text-primary">
                +{currentThread?.replies_count || 0}
              </div>
            </div>
          </div>
        </aside>
      </div>

      {/* Replies header */}
      <div className="flex items-center justify-between border-t border-border/60 pt-8">
        <h3 className="heading-display text-2xl text-foreground">
          {currentThread?.replies_count || 0} Replies
        </h3>
        <div className="flex gap-4 text-[10px] uppercase tracking-[0.18em]">
          <span className="text-muted-foreground">SORT_BY:</span>
          <button className="text-primary border-b border-primary">
            TOP_RATED
          </button>
          <button className="text-muted-foreground hover:text-foreground">
            LATEST
          </button>
        </div>
      </div>

      {/* Commments */}
      <div className="space-y-5">
        {currentThread?.comments?.length ? (
          currentThread.comments.map((comment) => (
            <div key={comment.id} className="panel p-5 space-y-3">
              <div className="flex items-start justify-between gap-3">
                <div className="flex items-center gap-3">
                  <div className="h-8 w-8 rounded-full bg-secondary border border-border grid place-items-center text-[10px] font-mono text-primary">
                    {comment.author?.username
                      ?.replace("@", "")
                      .slice(0, 2)
                      .toUpperCase() ?? "??"}
                  </div>
                  <div>
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-semibold text-foreground">
                        {comment.author?.username || "@unknown"}
                      </span>
                      <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                        {comment.created_at
                          ? new Date(comment.created_at).toLocaleDateString()
                          : ""}
                      </span>
                    </div>
                  </div>
                </div>
                <div className="flex items-center gap-2 text-primary text-sm font-mono">
                  +{comment.upvotes} <ThumbsUp className="h-3.5 w-3.5" />
                </div>
              </div>
              <p className="text-xs text-foreground/85 leading-relaxed">
                {comment.content}
              </p>
              {replyTarget?.id === comment.id && user && (
                <div className="ml-6 mt-4 pl-5 border-l border-primary/30 space-y-3">
                  <div className="flex items-center gap-3">
                    <div className="h-7 w-7 rounded-sm bg-gradient-signal grid place-items-center text-[9px] font-bold text-primary-foreground">
                      {user.username.replace("@", "").slice(0, 2).toUpperCase()}
                    </div>
                    <span className="text-xs uppercase tracking-[0.18em] text-muted-foreground">
                      REPLYING AS {user.username.toUpperCase()}
                    </span>
                  </div>
                  <Textarea
                    value={replyContent}
                    onChange={(e) => setReplyContent(e.target.value)}
                    placeholder={`@${replyTarget.author.replace("@", "")} `}
                    className="bg-terminal border-border min-h-[80px] font-mono text-sm focus-visible:ring-primary/40"
                  />
                  <div className="flex justify-end gap-2">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setReplyTarget(null)}
                      className="text-muted-foreground uppercase text-[10px]"
                    >
                      Cancel
                    </Button>
                    <Button
                      size="sm"
                      onClick={() => {
                        if (!requireAuth({ toast: true, description: "Authenticate to transmit a reply." })) return;
                        createComment.mutate(
                          { content: replyContent, parentId: replyTarget.id },
                          {
                            onSuccess: () => {
                              setReplyContent("");
                              setReplyTarget(null);
                              toast.success("REPLY_TRANSMITTED");
                            },
                            onError: (err) =>
                              toast.error("TRANSMIT_FAILED", { description: err.message }),
                          },
                        );
                      }}
                      disabled={createComment.isPending || !replyContent.trim()}
                      className="bg-gradient-signal hover:opacity-90 text-primary-foreground font-bold uppercase tracking-[0.18em] text-[10px] rounded-sm"
                    >
                      {createComment.isPending ? "TRANSMITTING..." : "TRANSMIT"}
                    </Button>
                  </div>
                </div>
              )}
              <div className="flex gap-4 text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                <button
                  className="hover:text-primary"
                  onClick={() => {
                    if (!requireAuth({ toast: true, description: "Authenticate to reply." })) return;
                    setReplyTarget({
                      id: comment.id,
                      author: comment.author?.username || "@unknown",
                    });
                    setReplyContent(`@${(comment.author?.username || "@unknown").replace("@", "")} `);
                  }}
                >
                  REPLY
                </button>
                {/*<button className="hover:text-primary">SHARE</button>*/}
              </div>

              {comment.replies?.map((reply) => (
                <div
                  key={reply.id}
                  className="ml-6 mt-4 pl-5 border-l border-primary/30 space-y-2"
                >
                  <div className="flex items-center gap-2">
                    <div className="h-7 w-7 rounded-full bg-gradient-signal grid place-items-center text-[9px] font-bold text-primary-foreground">
                      {reply.author?.username
                        ?.replace("@", "")
                        .slice(0, 2)
                        .toUpperCase() || "??"}
                    </div>
                    <span className="text-sm font-semibold text-foreground">
                      {reply.author?.username || "@unknown"}
                    </span>
                    <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                      {reply.created_at
                        ? new Date(reply.created_at).toLocaleDateString()
                        : ""}
                    </span>
                  </div>
                  <p className="text-sm text-foreground/85 leading-relaxed">
                    {reply.content}
                  </p>
                  <button className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground hover:text-primary">
                    REPLY
                  </button>
                </div>
              ))}
            </div>
          ))
        ) : (
          <p className="text-sm text-muted-foreground italic">
            No replies yet. Be the first to contribute.
          </p>
        )}
      </div>

      {/* Comment box */}
      {user && (
        <div className="panel p-5 space-y-3">
          <div className="flex items-center justify-between">
            <SectionLabel>Join_The_Discussion</SectionLabel>
            <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
              AS {user.username?.toUpperCase()}
            </span>
          </div>
          <Textarea
            value={replyContent}
            onChange={(e) => setReplyContent(e.target.value)}
            placeholder="ENTER_TERMINAL_INPUT..."
            className="bg-terminal border-border min-h-[120px] font-mono text-sm placeholder:text-muted-foreground/50 focus-visible:ring-primary/40"
          />
          <div className="flex items-center justify-between">
            <div className="flex gap-2 text-muted-foreground">
              <button className="hover:text-primary">
                <ImageIcon className="h-4 w-4" />
              </button>
              <button className="hover:text-primary">
                <Code2 className="h-4 w-4" />
              </button>
              <button className="hover:text-primary font-bold text-sm w-4">
                B
              </button>
            </div>
            <Button
              onClick={handleSubmitReply}
              disabled={createComment.isPending || !replyContent.trim()}
              className="bg-gradient-signal hover:opacity-90 text-primary-foreground font-bold uppercase tracking-[0.18em] text-xs rounded-sm"
            >
              {createComment.isPending ? "TRANSMITTING..." : "TRANSMIT_REPLY"}
            </Button>
          </div>
        </div>
      )}
    </div>
  );
};

export default ThreadDetailClient;
