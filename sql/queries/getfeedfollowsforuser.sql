-- name: GetFeedFollowsForUser :many
SELECT users.name AS username, feeds.name AS name
FROM feed_follows
INNER JOIN feeds ON feed_id = feeds.id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE users.id = $1;