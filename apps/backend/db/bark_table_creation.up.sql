create table if not exists bark (
    id text primary key,
    dog_id text not null,
    bark text not null,
    created_at timestamp default current_timestamp,
    foreign key (dog_id) references dog(id) on delete cascade
);

create index if not exists bark_dog_id on bark(dog_id);
