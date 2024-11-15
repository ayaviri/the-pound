create extension if not exists ltree;
create table if not exists thread (
    bark_id text primary key,
    thread_path ltree not null,
    root_bark_id text not null,
    created_at timestamp default current_timestamp,
    foreign key (bark_id) references bark(id) on delete cascade,
    foreign key (root_bark_id) references bark(id) on delete cascade
);
create index if not exists thread_thread_path on thread using GIST (thread_path);

-- 1) Gets the thread path of the newly inserted paw's parent bark
-- 2) Creates a thread for the parent if it did not exist
-- 3) Creates a thread for the newly inserted paw
create or replace function update_thread_path()
returns trigger as $$
declare
    parent_thread ltree;
    root_bark_id text;
begin
    -- 1)
    select t.thread_path, t.root_bark_id into parent_thread, root_bark_id
    from thread t where t.bark_id = NEW.original_bark_id;

    -- 2)
    if parent_thread is null then
        parent_thread = NEW.original_bark_id::text::ltree;
        root_bark_id = NEW.original_bark_id;

        insert into thread (bark_id, thread_path, root_bark_id)
        values (
            NEW.original_bark_id, 
            parent_thread, 
            root_bark_id
        );
    end if;

    -- 3)
    insert into thread (bark_id, thread_path, root_bark_id)
    values (
        NEW.bark_id,
        parent_thread || NEW.bark_id::text::ltree,
        root_bark_id
    );

    return NEW;
end;
$$ language plpgsql;

drop trigger if exists set_thread_path on paw;
create trigger set_thread_path
after insert on paw
for each row 
    execute function update_thread_path();
