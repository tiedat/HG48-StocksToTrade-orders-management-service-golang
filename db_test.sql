-- public.orders definition

-- Drop table

-- DROP TABLE public.orders;

-- name: create_orders_table
CREATE TABLE public.orders (
	id serial NOT NULL,
	recurly_uid varchar NULL,
	product_id int4 NULL,
	email varchar NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	subscription_uid varchar NULL,
	state varchar NULL DEFAULT 'active'::character varying,
	CONSTRAINT orders_pkey PRIMARY KEY (id)

);

-- name: create_email_index
CREATE INDEX index_orders_on_email ON public.orders USING btree (email);

-- name: create_orders_data
INSERT INTO orders (email, created_at, updated_at) VALUES
  ( 'foo@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP ),
  ( 'bar@example.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP );
