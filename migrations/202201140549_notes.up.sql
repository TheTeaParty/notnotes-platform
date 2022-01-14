create table notes
(
    id         varchar                                               not null
        constraint notes_pk
            primary key,
    name       varchar,
    content    text,
    created_at timestamp without time zone default current_timestamp not null,
    updated_at timestamp without time zone default current_timestamp not null
);

create unique index notes_id_uindex
    on notes (id);

