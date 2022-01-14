alter table note_tags
    add constraint note_tags_note_tags_id_fk
        foreign key (note_id) references note_tags
            on delete cascade;

alter table note_tags
    add constraint note_tags_tags_id_fk
        foreign key (tag_id) references tags
            on delete cascade;

