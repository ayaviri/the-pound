create table if not exists following (
    id text primary key,
    from_dog_id text not null,
    to_dog_id text not null,
    is_approved boolean not null, 
    created_at timestamp default current_timestamp,
    foreign key (from_dog_id) references dog(id) on delete cascade,
    foreign key (to_dog_id) references dog(id) on delete cascade,
    unique(from_dog_id, to_dog_id)
);

create index if not exists following_from_dog_id on following(from_dog_id);
create index if not exists following_to_dog_id on following(to_dog_id);
