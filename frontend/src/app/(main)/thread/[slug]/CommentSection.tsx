"use client";

import { useState } from "react";
import { SectionLabel } from "@/components/forge/SectionLabel";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { ThumbsUp, ThumbsDown, Image as ImageIcon, Code2 } from "lucide-react";
import { toast } from "sonner";
import { useCreateComment } from "@/hooks/use-comments";
import { useVoteComment } from "@/hooks/use-votes";
import { useMe } from "@/hooks/use-auth";
import { useRequireAuth } from "@/hooks/use-require-auth";
import type { CommentItem } from "@/lib/mock-data";
import { getInitials } from "@/lib/utils";

interface NestedCommentProps {
  comment: CommentItem;
  threadSlug: string;
}

function NestedComment({ comment, threadSlug }: NestedCommentProps) {
  const voteComment = useVoteComment(threadSlug);
  const { requireAuth } = useRequireAuth();

  const handleVote = (value: number) => {
    if (!requireAuth({ toast: true, description: "Authenticate to cast a signal." })) return;
    voteComment.mutate({ commentId: comment.id, value });
  };

  return (
    <div className="ml-6 mt-4 pl-5 border-l border-primary/30 space-y-2">
      <div className="flex items-center gap-2">
        <div className="h-7 w-7 rounded-full bg-gradient-signal grid place-items-center text-[9px] font-bold text-primary-foreground">
          {getInitials(comment.author?.username)}
        </div>
        <span className="text-sm font-semibold text-foreground">
          {comment.author?.username || "@unknown"}
        </span>
        <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
          {comment.created_at ? new Date(comment.created_at).toLocaleDateString() : ""}
        </span>
      </div>
      <p className="text-sm text-foreground/85 leading-relaxed">{comment.content}</p>
      <div className="flex items-center gap-1">
        <Button
          variant="outline"
          size="sm"
          className="rounded-sm h-6 px-1.5 text-[9px] border-border hover:border-primary hover:text-primary"
          onClick={() => handleVote(1)}
        >
          <ThumbsUp className="h-2.5 w-2.5 mr-0.5" /> {comment.upvotes}
        </Button>
        <Button
          variant="outline"
          size="sm"
          className="rounded-sm h-6 px-1.5 text-[9px] border-border hover:border-primary hover:text-primary"
          onClick={() => handleVote(-1)}
        >
          <ThumbsDown className="h-2.5 w-2.5" />
        </Button>
      </div>
    </div>
  );
}

interface CommentFormProps {
  slug: string;
  onSuccess: () => void;
}

function CommentForm({ slug, onSuccess }: CommentFormProps) {
  const [content, setContent] = useState("");
  const { data: user } = useMe();
  const createComment = useCreateComment(slug);
  const { requireAuth } = useRequireAuth();

  const handleSubmit = () => {
    if (!content.trim()) return;
    if (!requireAuth({ toast: true, description: "Authenticate to transmit a reply." })) return;
    createComment.mutate(
      { content },
      {
        onSuccess: () => {
          setContent("");
          toast.success("REPLY_TRANSMITTED");
          onSuccess();
        },
        onError: (err) => toast.error("TRANSMIT_FAILED", { description: err.message }),
      },
    );
  };

  if (!user) return null;

  return (
    <div className="panel p-5 space-y-3">
      <div className="flex items-center justify-between">
        <SectionLabel>Join_The_Discussion</SectionLabel>
        <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
          AS {user.username?.toUpperCase()}
        </span>
      </div>
      <Textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="ENTER_TERMINAL_INPUT..."
        className="bg-terminal border-border min-h-[120px] font-mono text-sm placeholder:text-muted-foreground/50 focus-visible:ring-primary/40"
      />
      <div className="flex items-center justify-between">
        <div className="flex gap-2 text-muted-foreground">
          <button className="hover:text-primary"><ImageIcon className="h-4 w-4" /></button>
          <button className="hover:text-primary"><Code2 className="h-4 w-4" /></button>
          <button className="hover:text-primary font-bold text-sm w-4">B</button>
        </div>
        <Button
          onClick={handleSubmit}
          disabled={createComment.isPending || !content.trim()}
          className="bg-gradient-signal hover:opacity-90 text-primary-foreground font-bold uppercase tracking-[0.18em] text-xs rounded-sm"
        >
          {createComment.isPending ? "TRANSMITTING..." : "TRANSMIT_REPLY"}
        </Button>
      </div>
    </div>
  );
}

interface CommentSectionProps {
  slug: string;
  comments?: CommentItem[];
  repliesCount: number;
}

export function CommentSection({ slug, comments, repliesCount }: CommentSectionProps) {
  const [replyTarget, setReplyTarget] = useState<{ id: number; author: string } | null>(null);
  const [replyContent, setReplyContent] = useState("");
  const { data: user } = useMe();
  const { requireAuth } = useRequireAuth();
  const createComment = useCreateComment(slug);

  const handleReplySubmit = () => {
    if (!replyContent.trim() || !replyTarget) return;
    if (!requireAuth({ toast: true, description: "Authenticate to transmit a reply." })) return;
    createComment.mutate(
      { content: replyContent, parentId: replyTarget.id },
      {
        onSuccess: () => {
          setReplyContent("");
          setReplyTarget(null);
          toast.success("REPLY_TRANSMITTED");
        },
        onError: (err) => toast.error("TRANSMIT_FAILED", { description: err.message }),
      },
    );
  };

  return (
    <>
      <div className="flex items-center justify-between border-t border-border/60 pt-8">
        <h3 className="heading-display text-2xl text-foreground">
          {repliesCount} Replies
        </h3>
        <div className="flex gap-4 text-[10px] uppercase tracking-[0.18em]">
          <span className="text-muted-foreground">SORT_BY:</span>
          <button className="text-primary border-b border-primary">TOP_RATED</button>
          <button className="text-muted-foreground hover:text-foreground">LATEST</button>
        </div>
      </div>

      <div className="space-y-5">
        {comments?.length ? (
          comments.map((comment) => (
            <div key={comment.id} className="panel p-5 space-y-3">
              <div className="flex items-start justify-between gap-3">
                <div className="flex items-center gap-3">
                  <div className="h-8 w-8 rounded-full bg-secondary border border-border grid place-items-center text-[10px] font-mono text-primary">
                    {getInitials(comment.author?.username)}
                  </div>
                  <div>
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-semibold text-foreground">
                        {comment.author?.username || "@unknown"}
                      </span>
                      <span className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
                        {comment.created_at ? new Date(comment.created_at).toLocaleDateString() : ""}
                      </span>
                    </div>
                  </div>
                </div>
                <VoteButtons commentId={comment.id} upvotes={comment.upvotes} slug={slug} />
              </div>
              <p className="text-xs text-foreground/85 leading-relaxed">{comment.content}</p>

              {replyTarget?.id === comment.id && user && (
                <div className="ml-6 mt-4 pl-5 border-l border-primary/30 space-y-3">
                  <div className="flex items-center gap-3">
                    <div className="h-7 w-7 rounded-sm bg-gradient-signal grid place-items-center text-[9px] font-bold text-primary-foreground">
                      {getInitials(user.username)}
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
                      onClick={() => { setReplyContent(""); setReplyTarget(null); }}
                      className="text-muted-foreground uppercase text-[10px]"
                    >
                      Cancel
                    </Button>
                    <Button
                      size="sm"
                      onClick={handleReplySubmit}
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
                    setReplyTarget({ id: comment.id, author: comment.author?.username || "@unknown" });
                    setReplyContent(`@${(comment.author?.username || "@unknown").replace("@", "")} `);
                  }}
                >
                  REPLY
                </button>
              </div>

              {comment.replies?.map((reply) => (
                <NestedComment key={reply.id} comment={reply} threadSlug={slug} />
              ))}
            </div>
          ))
        ) : (
          <p className="text-sm text-muted-foreground italic">
            No replies yet. Be the first to contribute.
          </p>
        )}
      </div>

      <CommentForm slug={slug} onSuccess={() => {}} />
    </>
  );
}

function VoteButtons({ commentId, upvotes, slug }: { commentId: number; upvotes: number; slug: string }) {
  const voteComment = useVoteComment(slug);
  const { requireAuth } = useRequireAuth();

  const handleVote = (value: number) => {
    if (!requireAuth({ toast: true, description: "Authenticate to cast a signal." })) return;
    voteComment.mutate({ commentId, value });
  };

  return (
    <div className="flex items-center gap-1.5">
      <Button
        variant="outline"
        size="sm"
        className="rounded-sm h-7 px-2 text-[10px] border-border hover:border-primary hover:text-primary"
        onClick={() => handleVote(1)}
      >
        <ThumbsUp className="h-3 w-3 mr-1" /> {upvotes}
      </Button>
      <Button
        variant="outline"
        size="sm"
        className="rounded-sm h-7 px-2 text-[10px] border-border hover:border-primary hover:text-primary"
        onClick={() => handleVote(-1)}
      >
        <ThumbsDown className="h-3 w-3" />
      </Button>
    </div>
  );
}
