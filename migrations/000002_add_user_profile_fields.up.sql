ALTER TABLE users
ADD COLUMN username TEXT,
ADD COLUMN avatar_url TEXT;

UPDATE users
SET username = regexp_replace(lower(split_part(email, '@', 1)), '[^a-z0-9_\\-\\+\\.]', '_', 'g')
WHERE username IS NULL;

UPDATE users
SET username = regexp_replace(username, '^[\\._\\+\\-]+|[\\._\\+\\-]+$', '', 'g')
WHERE username IS NOT NULL;

UPDATE users
SET username = CASE
  WHEN username IS NULL OR username = '' THEN 'user'
  WHEN length(username) < 3 THEN 'user_' || username
  ELSE username
END
WHERE username IS NOT NULL;

UPDATE users
SET username = left(username, 32)
WHERE length(username) > 32;

ALTER TABLE users
ALTER COLUMN username SET NOT NULL;
