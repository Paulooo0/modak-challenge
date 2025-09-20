-- name: CreateNotification :one
INSERT INTO notifications (user_id, type, message)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CountNotificationsInTimeWindow :one
SELECT COUNT(*) as total
FROM notifications
WHERE user_id = $1
  AND type = $2
  AND created_at >= $3;