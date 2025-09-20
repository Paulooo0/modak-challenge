CREATE TABLE notifications (
id uuid DEFAULT gen_random_uuid() NOT NULL,
user_id text NOT NULL,
type text NOT NULL,
message text NOT NULL,
created_at timestamp NOT NULL DEFAULT NOW()
);
