-- Drop key_pool first (it has foreign key to key_link)
DROP TABLE IF EXISTS key_pool;

-- Then drop key_link
DROP TABLE IF EXISTS key_link;

