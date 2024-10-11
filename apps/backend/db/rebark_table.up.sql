create table if not exists rebark (
    id text primary key,
    bark_id text not null,
    dog_id text not null, -- the ID of the dog that gave the rebark
    created_at timestamp default current_timestamp,
    foreign key (bark_id) references bark(id) on delete cascade,
    foreign key (dog_id) references dog(id) on delete cascade
);

create index if not exists rebark_bark_id on rebark(bark_id);
-- TODO: Could use one for dog_id to get all of the barks a dog gave a treat to
