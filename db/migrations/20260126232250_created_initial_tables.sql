-- +goose up
CREATE SCHEMA IF NOT EXISTS "public";

CREATE TABLE "public"."categories" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "name" text NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    "updated_at" timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("id")
);

CREATE TABLE "public"."products" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "name" text NOT NULL,
    "description" text,
    "price" bigint NOT NULL,
    "created_at" timestamp  NOT NULL DEFAULT NOW(),
    "updated_at" timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("id")
);

CREATE TABLE "public"."product_images" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "product_id" uuid NOT NULL,
    "image_url" text NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("id")
);

CREATE TABLE "public"."products_category" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "product_id" uuid NOT NULL,
    "category_id" uuid NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("id")
);

ALTER TABLE "public"."products_category" ADD CONSTRAINT "fk_products_category_category_id_categories_id" FOREIGN KEY("category_id") REFERENCES "public"."categories"("id") ON DELETE CASCADE;
ALTER TABLE "public"."product_images" ADD CONSTRAINT "fk_product_images_id_products_id" FOREIGN KEY("product_id") REFERENCES "public"."products"("id");
ALTER TABLE "public"."products_category" ADD CONSTRAINT "fk_products_category_product_id_products_id" FOREIGN KEY("product_id") REFERENCES "public"."products"("id");

-- +goose down
ALTER TABLE "public"."products_category" DROP CONSTRAINT IF EXISTS "fk_products_category_product_id_products_id";
ALTER TABLE "public"."product_images" DROP CONSTRAINT IF EXISTS "fk_product_images_id_products_id";
ALTER TABLE "public"."products_category" DROP CONSTRAINT IF EXISTS "fk_products_category_category_id_categories_id";

DROP TABLE IF EXISTS "public"."products_category";
DROP TABLE IF EXISTS "public"."product_images";
DROP TABLE IF EXISTS "public"."products";
DROP TABLE IF EXISTS "public"."categories";
