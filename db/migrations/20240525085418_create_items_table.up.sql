CREATE TABLE IF NOT EXISTS "items" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "name" TEXT NOT NULL,
    "category" TEXT NOT NULL,
    "price" INT NOT NULL,
    "image_url" TEXT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);