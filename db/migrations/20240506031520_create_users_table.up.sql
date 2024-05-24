CREATE TABLE IF NOT EXISTS "users" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "username" TEXT UNIQUE NOT NULL,
    "email" TEXT NOT NULL,
    "password" TEXT NULL,
    "role" TEXT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);