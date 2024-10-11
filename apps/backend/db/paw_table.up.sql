create table if not exists paw (
    id text primary key,
    original_bark_id text not null, -- ID of bark that was replied to
    bark_id text not null unique, -- ID of the paw in the bark table
    dog_id text not null,
    created_at timestamp default current_timestamp,
    foreign key (original_bark_id) references bark(id) on delete cascade,
    foreign key (bark_id) references bark(id) on delete cascade,
    foreign key (dog_id) references dog(id) on delete cascade
);
