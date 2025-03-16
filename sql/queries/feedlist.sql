-- name: GetFeeds :many
SELECT feeds.name, feeds.url, users.name
FROM feeds
LEFT JOIN users ON feeds.user_id = users.id;