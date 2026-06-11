import { describe, it, expect } from "vitest";
import { queryKeys } from "@/hooks/use-query-keys";

describe("queryKeys", () => {
  describe("auth", () => {
    it("me key", () => {
      expect(queryKeys.auth.me).toEqual(["me"]);
    });
  });

  describe("threads", () => {
    it("all key", () => {
      expect(queryKeys.threads.all).toEqual(["threads"]);
    });

    it("trending key", () => {
      expect(queryKeys.threads.trending).toEqual(["threads", "trending"]);
    });

    it("featured key", () => {
      expect(queryKeys.threads.featured).toEqual(["threads", "featured"]);
    });

    it("list key with params", () => {
      expect(queryKeys.threads.list({ page: 1, pageSize: 5, sort: "latest" })).toEqual([
        "threads",
        { page: 1, pageSize: 5, sort: "latest" },
      ]);
    });

    it("detail key with slug", () => {
      expect(queryKeys.threads.detail("hello-world")).toEqual(["threads", "hello-world"]);
    });

    it("comments key with slug", () => {
      expect(queryKeys.threads.comments("hello-world")).toEqual(["threads", "hello-world", "comments"]);
    });
  });

  describe("users", () => {
    it("profile key", () => {
      expect(queryKeys.users.profile("alice")).toEqual(["users", "alice"]);
    });

    it("threads key", () => {
      expect(queryKeys.users.threads("alice")).toEqual(["users", "alice", "threads"]);
    });
  });

  describe("votes", () => {
    it("thread vote key", () => {
      expect(queryKeys.votes.thread("hello-world")).toEqual(["votes", "thread", "hello-world"]);
    });
  });
});