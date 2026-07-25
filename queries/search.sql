-- name: Search :many
SELECT entity_type, entity_id, user_id, title, content, tags, rank
FROM search_index
WHERE user_id = ? AND search_index MATCH ?
ORDER BY rank
LIMIT 20;
