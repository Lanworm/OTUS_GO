create table c_user (
    id uuid not null default gen_random_uuid() PRIMARY KEY,
    first_name varchar(255) not null,
    last_name varchar(255) not null,
    created_at timestamp without time zone default now(),
    updated_at timestamp without time zone
)
