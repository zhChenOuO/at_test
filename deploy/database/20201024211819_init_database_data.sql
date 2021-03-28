-- +goose Up
DROP TABLE IF EXISTS "identity_accounts";

CREATE TABLE "identity_accounts" (
    "id" serial NOT NULL PRIMARY KEY,
    "name" TEXT UNIQUE,
    "email" TEXT UNIQUE,
    "phone" TEXT,
    "phone_area_code" TEXT,
    "password" TEXT,
    "phone_verify_status" INT,
    "send_phone_verify_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "email_verify_status" INT,
    "send_email_verify_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "accept_language" TEXT,
    "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX phone_and_phone_area_code ON identity_accounts (phone, phone_area_code);

COMMENT ON COLUMN "identity_accounts"."name" IS '姓名';

COMMENT ON COLUMN "identity_accounts"."email" IS 'email';

COMMENT ON COLUMN "identity_accounts"."phone" IS '電話號碼';

COMMENT ON COLUMN "identity_accounts"."phone_area_code" IS '電話號碼區域碼';

COMMENT ON COLUMN "identity_accounts"."accept_language" IS '語系';

COMMENT ON COLUMN "identity_accounts"."send_phone_verify_at" IS '發送簡訊驗證時間';

COMMENT ON COLUMN "identity_accounts"."send_email_verify_at" IS '發送email驗證時間';