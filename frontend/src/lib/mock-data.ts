export interface Author {
  id: number;
  username: string;
  avatar: string;
}

export interface Tag {
  id: number;
  name: string;
  color: string;
}

export interface ThreadItem {
  id: number;
  title: string;
  slug: string;
  content: string;
  status: string;
  view_count: number;
  upvotes: number;
  downvotes: number;
  replies_count: number;
  author: Author;
  tags: Tag[];
  created_at: string;
  recent_commenters?: Author[];
}

export interface CommentItem {
  id: number;
  content: string;
  upvotes: number;
  downvotes: number;
  author: Author;
  replies: CommentItem[];
  created_at: string;
}

export interface ThreadDetail extends ThreadItem {
  comments: CommentItem[];
}

export interface Pagination {
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
}

export interface ThreadsResponse {
  threads: ThreadItem[];
  pagination: Pagination;
}

export interface UserProfile {
  id: number;
  username: string;
  avatar: string;
  bio: string;
  role: string;
  stacks: string[];
  created_at: string;
}
