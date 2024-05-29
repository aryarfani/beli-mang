CREATE TABLE IF NOT EXISTS "orders" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "user_id" UUID NOT NULL,
    "merchant_id" UUID NOT NULL,
    "estimation_id" UUID NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("merchant_id") REFERENCES "merchants" ("id");