CREATE TABLE IF NOT EXISTS "estimations" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "total_price" INT NOT NULL,
    "estimated_delivery_time_in_minutes" INT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)