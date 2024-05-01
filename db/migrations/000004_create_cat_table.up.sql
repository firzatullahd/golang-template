CREATE TABLE IF NOT EXISTS cat(
        id serial PRIMARY KEY,
        user_id BIGINT NOT NULL,

        name text not null,
        image_urls text[] not null,
        sex cat_sex not null,
        race cat_race not null,

        created_at timestamptz NOT NULL DEFAULT now(),
        updated_at timestamptz NOT NULL DEFAULT now(),
        deleted_at timestamptz DEFAULT NULL,

        CONSTRAINT fk_users_cats foreign key(user_id) references public.users(id)
);