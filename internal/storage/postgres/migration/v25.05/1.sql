alter table logger.object_config
    drop constraint object_config_wbt_class_id_fk;

alter table logger.object_config
    add constraint object_config_wbt_class_id_fk
        foreign key (object_id) references directory.wbt_class
            on update cascade on delete cascade;

alter table logger.log
    alter column record_id type text using record_id::text;

