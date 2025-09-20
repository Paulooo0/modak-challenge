CREATE INDEX idx_notifications_user_type_time
		ON notifications(user_id, type, created_at);
