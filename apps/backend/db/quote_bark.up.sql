create table if not exists quote_bark (
    id text primary key,
    bark_id text not null,
    dog_id text not null, -- the ID of the dog that gave the quote bark
    created_at timestamp default current_timestamp,
    quote_bark text not null,
    foreign key (bark_id) references bark(id) on delete cascade
    foreign key (dog_id) references dog(id) on delete cascade
);
