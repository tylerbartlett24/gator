-- name: GetPostsForUser :many
SELECT posts.*, feeds.name AS feed_name
FROM feed_follows
INNER JOIN posts ON feed_follows.feed_id = posts.feed_id
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1
LIMIT $2;