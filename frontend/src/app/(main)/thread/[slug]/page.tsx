import ThreadDetailClient from "./ThreadDetailClient";
import { fetchJSON, apiUrl } from "@/lib/server-fetch";

export default async function ThreadPage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;
  const thread = await fetchJSON<any>(apiUrl(`/threads/${slug}`));
  return <ThreadDetailClient slug={slug} initialThread={thread} />;
}
