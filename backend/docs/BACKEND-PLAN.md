# Features

- Authentication with username and password, signup with username and password
- create post with title, description, tag, authorId, text (markdown text), status (draft/published)
- comment with text, userId (commentor), postId
- reply with text, authorId, commentId
- upvote with keeping userId
- downvote with keeping userId


# User
- userId
- name
- description
- tag (optional)
- status (optional)
- staks (optional)


# Post 
- postId
- text
- tag
- authorId
- reply (store all replies)
- comment (store all comments)
- upvote (count of upvote)
- downvote (count of downvote)
- views (count of visits)
- replies_count

# Comment
- commentId
- text
- authorId
- postId
- upvote (count of upvote)
- downvote (count of downvote)


# Data fetching
- get featured post with most vote and replies
- get 3 of popular poasts 
- get all post with pagination (pagesize = 5), ordered by desc
- get post by id    
- get posts by userId
 