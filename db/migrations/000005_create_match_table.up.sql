CREATE TABLE IF NOT EXISTS match(
         id serial PRIMARY KEY,
         user_id BIGINT NOT NULL,
         cat_id BIGINT NOT NULL,
         match_cat_id BIGINT NOT NULL,
         match_user_id BIGINT NOT NULL,
         message text not null,

         is_approved boolean not null default false,
         is_rejected boolean not null default false,

         created_at timestamptz NOT NULL DEFAULT now(),
         updated_at timestamptz NOT NULL DEFAULT now(),
         deleted_at timestamptz DEFAULT NULL,

        CONSTRAINT fk_users_match foreign key(user_id) references public.users(id),
        CONSTRAINT fk_users_match2 foreign key(match_user_id) references public.users(id),
        CONSTRAINT fk_cats_match foreign key(cat_id) references public.cat(id),
        CONSTRAINT fk_cats_match2 foreign key(match_cat_id) references public.cat(id)

);

create index if not exists idx_match_user_id on match(user_id); 
create index if not exists idx_match_match_user_id on match(match_user_id); 