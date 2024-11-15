create table if not exists notification (
    id text primary key,
    is_read boolean not null default false,
    to_dog_id text not null,
    type text not null,
    payload text not null,
    created_at timestamp default current_timestamp,
    foreign key (to_dog_id) references dog(id) on delete cascade
);

create index if not exists notification_to_dog_id_is_read on notification(to_dog_id, is_read);


--  _____ ____  _____    _  _____   _____ ____  ___ ____  ____ _____ ____  
-- |_   _|  _ \| ____|  / \|_   _| |_   _|  _ \|_ _/ ___|/ ___| ____|  _ \ 
--   | | | |_) |  _|   / _ \ | |     | | | |_) || | |  _| |  _|  _| | |_) |
--   | | |  _ <| |___ / ___ \| |     | | |  _ < | | |_| | |_| | |___|  _ < 
--   |_| |_| \_\_____/_/   \_\_|     |_| |_| \_\___\____|\____|_____|_| \_\
--                                                                         

-- 1) Gets the information on the affected bark
-- 2) Joins it with the dog table to get the affected bark's author info
-- 3) Filters out bark if affected by author themself
-- 4) Writes notification
create or replace function write_bark_related_notification()
returns trigger as $$
declare
    notification_type text;
begin
    notification_type := TG_ARGV[0];

    -- 4)
    insert into notification (
        id, 
        is_read, 
        to_dog_id, 
        type, 
        payload
    )
    -- 1)
    with affected_bark as (
        select
            id,
            bark,
            dog_id
        from bark
        where bark.id = NEW.bark_id
    )
    select 
        gen_random_uuid()::text,
        false,
        b.dog_id,
        notification_type,
        json_build_object(
            'from_dog_id', NEW.dog_id,
            'from_dog_username', d.username,
            'bark_id', b.id,
            'bark', b.bark
        )::text
    -- 2)
    from affected_bark join dog d on affected_bark.dog_id = d.id
    -- 3)
    where affected_bark.dog_id != NEW.dog_id;

    return NEW;
end;
$$ language plpgsql;

drop trigger if exists write_treat_notification_trigger on treat;
create trigger write_treat_notification_trigger
after insert on treat
for each row
    execute function write_bark_related_notification('treat');


--  ____  _____ ____    _    ____  _  __  _____ ____  ___ ____  ____ _____ ____  
-- |  _ \| ____| __ )  / \  |  _ \| |/ / |_   _|  _ \|_ _/ ___|/ ___| ____|  _ \ 
-- | |_) |  _| |  _ \ / _ \ | |_) | ' /    | | | |_) || | |  _| |  _|  _| | |_) |
-- |  _ <| |___| |_) / ___ \|  _ <| . \    | | |  _ < | | |_| | |_| | |___|  _ < 
-- |_| \_\_____|____/_/   \_\_| \_\_|\_\   |_| |_| \_\___\____|\____|_____|_| \_\
--                                                                               

drop trigger if exists write_rebark_notification_trigger on rebark;
create trigger write_rebark_notification_trigger
after insert on rebark
for each row
    execute function write_bark_related_notification('rebark');


--  ____   ___        __  _____ ____  ___ ____  ____ _____ ____  
-- |  _ \ / \ \      / / |_   _|  _ \|_ _/ ___|/ ___| ____|  _ \ 
-- | |_) / _ \ \ /\ / /    | | | |_) || | |  _| |  _|  _| | |_) |
-- |  __/ ___ \ V  V /     | | |  _ < | | |_| | |_| | |___|  _ < 
-- |_| /_/   \_\_/\_/      |_| |_| \_\___\____|\____|_____|_| \_\
--                                                               

-- 1) Gets the ID of the dog who authored the parent bark of the
-- newly inserted paw
-- 2) Gets the bark information of the newly inserted paw
-- 3) Joins it with the dog table to get the paw's author info
-- 4) Filters out the paw if the author has replied to themself
-- 5) Writes notification
create or replace function write_paw_notification()
returns trigger as $$
declare
    to_dog_id text;
begin
    -- 1)
    select d.id into to_dog_id 
    from dog d join bark b on d.id = b.dog_id
    where b.id = NEW.original_bark_id;

    -- 5)
    insert into notification (
        id,
        is_read,
        to_dog_id,
        type,
        payload
    )
    -- 2)
    with paw_bark as (
        select 
            id, 
            bark, 
            dog_id
        from bark 
        where bark.id = NEW.bark_id
    )
    select
        gen_random_uuid()::text,
        false,
        to_dog_id,
        'paw',
        json_build_object(
            'from_dog_id', NEW.dog_id,
            'from_dog_username', d.username,
            'bark_id', paw_bark.id,
            'bark', paw_bark.bark
        )::text
    -- 3)
    from paw_bark join dog d on paw_bark.dog_id = d.id
    -- 4)
    where to_dog_id != d.id;

    return NEW;
end;
$$ language plpgsql;

drop trigger if exists write_paw_notification_trigger on paw;
create trigger write_paw_notification_trigger
after insert on paw 
for each row
    execute function write_paw_notification();


--  _____ ___  _     _     _____        __
-- |  ___/ _ \| |   | |   / _ \ \      / /
-- | |_ | | | | |   | |  | | | \ \ /\ / / 
-- |  _|| |_| | |___| |__| |_| |\ V  V /  
-- |_|   \___/|_____|_____\___/  \_/\_/   
--                                        
--  _____ ____  ___ ____  ____ _____ ____  
-- |_   _|  _ \|_ _/ ___|/ ___| ____|  _ \ 
--   | | | |_) || | |  _| |  _|  _| | |_) |
--   | | |  _ < | | |_| | |_| | |___|  _ < 
--   |_| |_| \_\___\____|\____|_____|_| \_\
--                                         

create or replace function write_follow_notification()
returns trigger as $$
begin
    insert into notification (
        id, 
        is_read, 
        to_dog_id, 
        type, 
        payload
    )
    select 
        gen_random_uuid()::text,
        false,
        NEW.to_dog_id,
        'follow',
        json_build_object(
            'from_dog_id', NEW.from_dog_id,
            'from_dog_username', from_dog.username,
            'is_approved', NEW.is_approved
        )::text
    from dog from_dog where from_dog.id = NEW.from_dog_id;

    return NEW;
end;
$$ language plpgsql;

drop trigger if exists write_follow_notification_trigger on following;
create trigger write_follow_notification_trigger
after insert on following
for each row
    execute function write_follow_notification();
