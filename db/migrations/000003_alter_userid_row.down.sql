ALTER TABLE notifications 
ALTER COLUMN user_id TYPE TEXT USING user_id::text;