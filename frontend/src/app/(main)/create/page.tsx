"use client";

import dynamic from "next/dynamic";

const CreateEntryInner = dynamic(() => import("./CreateEntryInner"), {
  ssr: false,
  loading: () => (
    <div className="px-8 py-10 max-w-[1100px] mx-auto space-y-8">
      <div className="space-y-3">
        <div className="h-5 w-48 bg-secondary/60 animate-pulse rounded-sm" />
        <div className="h-12 w-full bg-secondary/60 animate-pulse rounded-sm" />
      </div>
      <div className="h-[400px] bg-secondary/60 animate-pulse rounded-sm" />
    </div>
  ),
});

export default function CreateEntry() {
  return <CreateEntryInner />;
}
