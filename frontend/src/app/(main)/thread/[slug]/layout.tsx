import type { Metadata } from "next";
import { fetchJSON, apiUrl } from "@/lib/server-fetch";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ slug: string }>;
}): Promise<Metadata> {
  const { slug } = await params;
  try {
    const thread = await fetchJSON<any>(apiUrl(`/threads/${slug}`));
    return {
      title: `${thread?.title || "Thread"} | Midnight Forge`,
      description:
        thread?.content?.slice(0, 160) || "View this thread on Midnight Forge.",
      openGraph: {
        title: thread?.title || "Thread",
        description:
          thread?.content?.slice(0, 160) ||
          "View this thread on Midnight Forge.",
        type: "article",
      },
    };
  } catch {
    return {
      title: "Thread | Midnight Forge",
      description: "View this thread on Midnight Forge.",
    };
  }
}

export default function ThreadLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return children;
}
