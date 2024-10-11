create table if not exists notification (
    id text primary key,
    is_read boolean not null default false,
    for_dog_id text not null,
    type text not null,
    payload text not null,
    created_at timestamp default current_timestamp,
    foreign key (for_dog_id) references dog(id) on delete cascade
);

create index if not exists notification_for_dog_id_is_read on notification(for_dog_id, is_read);
