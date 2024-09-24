drop trigger if exists cc_agent_set_rbac_acl on logger.object_config;


create trigger object_config_set_rbac_acl
    after insert
    on logger.object_config
    for each row
execute procedure logger.tg_obj_default_rbac('logger');