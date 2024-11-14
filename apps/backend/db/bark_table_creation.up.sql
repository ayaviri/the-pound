create table if not exists bark (
    id text primary key,
    dog_id text not null,
    bark text not null,
    created_at timestamp default current_timestamp,
    foreign key (dog_id) references dog(id) on delete cascade
    -- denormalised fields
    treat_count int not null default 0,
    rebark_count int not null default 0,
    paw_count int not null default 0,
    dog_username text not null,
);

create index if not exists bark_dog_id on bark(dog_id);
