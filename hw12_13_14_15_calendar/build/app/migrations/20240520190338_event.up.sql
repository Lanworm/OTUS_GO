create table c_event
(
    id             uuid                        not null default gen_random_uuid() PRIMARY KEY,
    title          varchar(255)                not null,
    start_datetime timestamp without time zone not null,
    end_datetime   timestamp without time zone not null,
    description    text,
    user_id        uuid                        not null,
    remind_before  int,
    check (
        end_datetime > start_datetime
    )
)
