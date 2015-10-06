DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users" (
	"name" varchar COLLATE "default",
	"alive" bool,
	"created_at" timestamp(6) NOT NULL,
	"updated_at" timestamp(6) NOT NULL,
	"birth_date" timestamp(6) NULL,
	"bio" varchar COLLATE "default",
	"price" numeric,
	"id" serial NOT NULL
)
WITH (OIDS=FALSE);

DROP TABLE IF EXISTS "public"."good_friends";
CREATE TABLE "public"."good_friends" (
	"first_name" varchar NOT NULL COLLATE "default",
	"last_name" varchar NOT NULL COLLATE "default",
	"id" serial NOT NULL
)
WITH (OIDS=FALSE);

DROP TABLE IF EXISTS "public"."cakes";
CREATE TABLE "public"."cakes" (
  "int_slice" int[],
  "float_slice" numeric[],
  "string_slice" varchar[],
	"id" serial NOT NULL
)
WITH (OIDS=FALSE);
