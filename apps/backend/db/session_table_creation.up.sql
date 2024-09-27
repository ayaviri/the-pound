create table if not exists session (
    id text primary key,
    dog_id text not null,
    token text not null unique,
    expires_at timestamp not null,
    created_at timestamp not null,
    foreign key (dog_id) references dog(id) on delete cascade
);

create index if not exists session_dog_id on session(dog_id);
create unique index if not exists session_token on session(token);
