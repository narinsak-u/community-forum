import { describe, it, expect } from "vitest";
import { getInitials, timeAgo, formatCount, cn } from "@/lib/utils";

describe("getInitials", () => {
  it("returns first 2 chars uppercase from username", () => {
    expect(getInitials("@john_doe")).toBe("JO");
  });

  it("returns first 2 chars without @ prefix", () => {
    expect(getInitials("jane_smith")).toBe("JA");
  });

  it("handles undefined username", () => {
    expect(getInitials()).toBe("??");
  });

  it("handles empty username", () => {
    expect(getInitials("")).toBe("??");
  });
});

describe("formatCount", () => {
  it("formats thousands with K suffix", () => {
    expect(formatCount(1200)).toBe("1.2K");
  });

  it("returns string for small numbers", () => {
    expect(formatCount(42)).toBe("42");
  });

  it("handles zero", () => {
    expect(formatCount(0)).toBe("0");
  });
});

describe("timeAgo", () => {
  it("returns JUST NOW for very recent dates", () => {
    const now = new Date();
    expect(timeAgo(now.toISOString())).toBe("JUST NOW");
  });

  it("returns MIN AGO for dates within the hour", () => {
    const fiveMinAgo = new Date(Date.now() - 5 * 60 * 1000);
    expect(timeAgo(fiveMinAgo.toISOString())).toBe("5 MIN AGO");
  });

  it("returns H AGO for dates within the day", () => {
    const twoHoursAgo = new Date(Date.now() - 2 * 60 * 60 * 1000);
    expect(timeAgo(twoHoursAgo.toISOString())).toBe("2H AGO");
  });

  it("returns D AGO for older dates", () => {
    const threeDaysAgo = new Date(Date.now() - 3 * 24 * 60 * 60 * 1000);
    expect(timeAgo(threeDaysAgo.toISOString())).toBe("3D AGO");
  });
});

describe("cn", () => {
  it("merges tailwind classes", () => {
    expect(cn("px-2 py-1", "px-4")).toBe("py-1 px-4");
  });

  it("handles falsy values", () => {
    expect(cn("base", undefined, "extra")).toBe("base extra");
  });
});