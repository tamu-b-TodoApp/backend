-- Create "companies" table
CREATE TABLE "public"."companies" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "name" character varying(255) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "company_members" table
CREATE TABLE "public"."company_members" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "company_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_company_members_company_id" to table: "company_members"
CREATE INDEX "idx_company_members_company_id" ON "public"."company_members" ("company_id");
-- Create index "idx_company_members_user_id" to table: "company_members"
CREATE INDEX "idx_company_members_user_id" ON "public"."company_members" ("user_id");
-- Create "team_members" table
CREATE TABLE "public"."team_members" (
  "team_id" bigint NOT NULL,
  "company_member_id" bigint NOT NULL,
  PRIMARY KEY ("team_id", "company_member_id")
);
-- Create "teams" table
CREATE TABLE "public"."teams" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "company_id" bigint NOT NULL,
  "name" character varying(255) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_teams_company_id" to table: "teams"
CREATE INDEX "idx_teams_company_id" ON "public"."teams" ("company_id");
-- Create "todos" table
CREATE TABLE "public"."todos" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "team_id" bigint NOT NULL,
  "parent_id" bigint NULL,
  "assignee_id" bigint NULL,
  "title" character varying(255) NOT NULL,
  "description" text NOT NULL,
  "status" character varying(20) NOT NULL DEFAULT 'not_started',
  "due_date" date NULL,
  "story_points" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_todos_team_id" to table: "todos"
CREATE INDEX "idx_todos_team_id" ON "public"."todos" ("team_id");
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "email" character varying(255) NOT NULL,
  "password" character varying(255) NOT NULL,
  "email_verified_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_email" to table: "users"
CREATE INDEX "idx_users_email" ON "public"."users" ("email");
