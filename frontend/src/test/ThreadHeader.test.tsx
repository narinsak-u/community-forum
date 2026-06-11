// @vitest-environment jsdom

import React from "react";
import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { ThreadHeader } from "@/app/(main)/thread/[slug]/ThreadHeader";

vi.mock("next/link", () => ({
  default: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}));

vi.mock("lucide-react", () => ({
  Pencil: () => <span data-testid="pencil-icon" />,
}));

describe("ThreadHeader", () => {
  const baseThread = {
    title: "Test Thread Title",
    tags: [{ name: "Technical" }],
    author: { id: 1, username: "@test_user" },
    created_at: "2026-01-15T00:00:00Z",
    view_count: 1200,
    replies_count: 42,
    upvotes: 95,
    downvotes: 5,
  };

  it("renders thread title", () => {
    render(<ThreadHeader thread={baseThread} authorInitials="TU" />);
    expect(screen.getByText("Test Thread Title")).toBeDefined();
  });

  it("renders tag name", () => {
    render(<ThreadHeader thread={baseThread} authorInitials="TU" />);
    expect(screen.getByText("● Technical")).toBeDefined();
  });

  it("renders author info", () => {
    render(<ThreadHeader thread={baseThread} authorInitials="TU" />);
    expect(screen.getByText("@test_user")).toBeDefined();
    expect(screen.getByText("TU")).toBeDefined();
  });

  it("renders stats with formatted view count", () => {
    render(<ThreadHeader thread={baseThread} authorInitials="TU" />);
    expect(screen.getByText("1.2K")).toBeDefined();
    expect(screen.getByText("VIEWS")).toBeDefined();
  });

  it("renders replies count", () => {
    render(<ThreadHeader thread={baseThread} authorInitials="TU" />);
    expect(screen.getByText("42")).toBeDefined();
    expect(screen.getByText("REPLIES")).toBeDefined();
  });

  it("renders trust index", () => {
    render(<ThreadHeader thread={baseThread} authorInitials="TU" />);
    expect(screen.getByText("TRUST_INDEX")).toBeDefined();
    expect(screen.getByText("95%")).toBeDefined();
  });

  it("shows edit link when current user is author", () => {
    render(
      <ThreadHeader
        thread={baseThread}
        slug="test-thread"
        currentUserId={1}
        authorInitials="TU"
      />,
    );
    expect(screen.getByText("EDIT")).toBeDefined();
  });

  it("does not show edit link for different user", () => {
    render(
      <ThreadHeader
        thread={baseThread}
        slug="test-thread"
        currentUserId={999}
        authorInitials="TU"
      />,
    );
    expect(screen.queryByText("EDIT")).toBeNull();
  });

  it("uses fallbacks when thread is empty", () => {
    render(
      <ThreadHeader
        thread={{ view_count: 0, replies_count: 0, upvotes: 0, downvotes: 0 }}
        authorInitials="??"
      />,
    );
    expect(screen.getByText("Thread")).toBeDefined();
    expect(screen.getByText("@unknown")).toBeDefined();
  });
});