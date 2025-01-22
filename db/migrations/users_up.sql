CREATE TYPE if not exists users_state  AS ENUM ('REGISTERED', 'PENDING_VERIFICATION', 'VERIFIED', 'DELETED');

CREATE TABLE public.users (
	id bigserial NOT NULL,
	username varchar NOT NULL,
	"password" text NOT NULL,
	state public."users_state" DEFAULT 'REGISTERED'::users_state NOT NULL,
	"name" varchar(50) NOT NULL,
  id_card_no text,
  id_card_file text,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (id),
	CONSTRAINT users_username_unique UNIQUE (username)
);
CREATE UNIQUE INDEX idx_password ON public.users USING btree (password);
CREATE UNIQUE INDEX idx_username ON public.users USING btree (username);