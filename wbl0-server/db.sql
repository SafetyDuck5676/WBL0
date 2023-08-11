-- Adminer 4.8.1 PostgreSQL 15.3 (Debian 15.3-1.pgdg120+1) dump

\connect "postgres";

CREATE SEQUENCE deliveries_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."deliveries" (
    "id" integer DEFAULT nextval('deliveries_id_seq') NOT NULL,
    "name" character varying(255) NOT NULL,
    "phone" character varying(255) NOT NULL,
    "zip" character varying(255) NOT NULL,
    "city" character varying(255) NOT NULL,
    "address" character varying(255) NOT NULL,
    "region" character varying(255) NOT NULL,
    "email" character varying(255) NOT NULL,
    CONSTRAINT "deliveries_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

INSERT INTO "deliveries" ("id", "name", "phone", "zip", "city", "address", "region", "email") VALUES
(1,	'Test Testov',	'+9720000000',	'2639809',	'Kiryat Mozkin',	'Ploshad Mira 15',	'Kraiot',	'test@gmail.com');

CREATE TABLE "public"."item_order" (
    "item_id" integer NOT NULL,
    "order_id" integer NOT NULL
) WITH (oids = false);

INSERT INTO "item_order" ("item_id", "order_id") VALUES
(1,	1),
(2,	1);

CREATE SEQUENCE items_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."items" (
    "id" integer DEFAULT nextval('items_id_seq') NOT NULL,
    "chrt_id" integer NOT NULL,
    "rid" character varying(255) NOT NULL,
    "name" character varying(255) NOT NULL,
    "sale" integer NOT NULL,
    "size" integer NOT NULL,
    "nm_id" integer NOT NULL,
    "brand" character varying NOT NULL,
    "status" integer NOT NULL,
    "total_price" integer NOT NULL,
    "price" integer NOT NULL,
    CONSTRAINT "items_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

INSERT INTO "items" ("id", "chrt_id", "rid", "name", "sale", "size", "nm_id", "brand", "status", "total_price", "price") VALUES
(1,	9934930,	'ab4219087a764ae0btest',	'Mascaras',	30,	0,	2389212,	'Vivienne Sabo',	202,	317,	453),
(2,	9934930,	'ab4219087a764ae0btest',	'Mascaras',	30,	0,	2389212,	'Vivienne Sabo',	202,	317,	453);

CREATE SEQUENCE orders_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;

CREATE TABLE "public"."orders" (
    "id" bigint DEFAULT nextval('orders_id_seq') NOT NULL,
    "order_uid" character varying(255) NOT NULL,
    "track_number" character varying(255) NOT NULL,
    "entry" character varying(255) NOT NULL,
    "locale" character varying(255) NOT NULL,
    "internal_signature" character varying(255) NOT NULL,
    "customer_id" character varying(255) NOT NULL,
    "delivery_service" character varying(255) NOT NULL,
    "shardkey" integer NOT NULL,
    "sm_id" integer NOT NULL,
    "date_created" timestamp NOT NULL,
    "oof_shard" integer NOT NULL,
    "payment_id" integer NOT NULL,
    "delivery_id" integer NOT NULL,
    CONSTRAINT "orders_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

INSERT INTO "orders" ("id", "order_uid", "track_number", "entry", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard", "payment_id", "delivery_id") VALUES
(1,	'b563feb7b2b84b6test',	'WBILMTESTTRACK',	'WBIL',	'en',	' ',	'test',	'meest',	9,	99,	'2023-08-02 18:31:20.785145',	1,	1,	1);

CREATE SEQUENCE payments_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."payments" (
    "id" integer DEFAULT nextval('payments_id_seq') NOT NULL,
    "transaction" character varying(255) NOT NULL,
    "request_id" character varying(255) NOT NULL,
    "currency" character varying(255) NOT NULL,
    "provider" character varying(255) NOT NULL,
    "amount" integer NOT NULL,
    "payment_dt" integer NOT NULL,
    "bank" character varying(255) NOT NULL,
    "delivery_cost" integer NOT NULL,
    "goods_total" integer NOT NULL,
    "custom_fee" integer NOT NULL,
    CONSTRAINT "payments_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

INSERT INTO "payments" ("id", "transaction", "request_id", "currency", "provider", "amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee") VALUES
(1,	'b563feb7b2b84b6test',	' ',	'USD',	'wbpay',	1817,	1637907727,	'alpha',	1500,	317,	0);

ALTER TABLE ONLY "public"."item_order" ADD CONSTRAINT "item_order_item_id_fkey" FOREIGN KEY (item_id) REFERENCES items(id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."item_order" ADD CONSTRAINT "item_order_order_id_fkey" FOREIGN KEY (order_id) REFERENCES orders(id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."orders" ADD CONSTRAINT "orders_delivery_id_fkey" FOREIGN KEY (delivery_id) REFERENCES deliveries(id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."orders" ADD CONSTRAINT "orders_payment_id_fkey" FOREIGN KEY (payment_id) REFERENCES payments(id) NOT DEFERRABLE;

-- 2023-08-02 22:29:25.438854+00