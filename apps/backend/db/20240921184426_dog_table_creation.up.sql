 create table if not exists dog (
    id text primary key,
    username text not null,
    password_hash text not null,
    created_at timestamp default current_timestamp
);

create unique index if not exists dog_id on dog(id);