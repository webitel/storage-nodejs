
CREATE OR REPLACE FUNCTION file_increment_profile_size()
  RETURNS trigger AS $file_increment_profile_size$
BEGIN
  update file_backend_profiles p
      SET data_size = data_size + f.sum_size,
      data_count = data_count + f.count_files
      from (
             select
               profile_id,
               sum(size) * 0.000001 as sum_size,
               count(*) as count_files
             from tg_data
             group by profile_id
           ) as f
      where p.id = f.profile_id;
  RETURN NULL;
END;
$file_increment_profile_size$
LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION file_decrement_profile_size()
  RETURNS trigger AS $file_decrement_profile_size$
BEGIN
  update file_backend_profiles p
      SET data_size = data_size - f.sum_size,
        data_count = data_count - f.count_files
      from (
             select
               profile_id,
               sum(size) * 0.000001 as sum_size,
               count(*) as count_files
             from tg_data
             group by profile_id
           ) as f
      where p.id = f.profile_id;
  RETURN NULL;
END;
$file_decrement_profile_size$
LANGUAGE plpgsql;



CREATE TRIGGER tg_file_change_profile_size_increment
  AFTER INSERT ON files
  REFERENCING NEW TABLE AS tg_data
  FOR EACH STATEMENT
  EXECUTE PROCEDURE file_increment_profile_size();

CREATE TRIGGER tg_file_decrement_profile_size
  AFTER DELETE ON files
  REFERENCING OLD TABLE AS tg_data
  FOR EACH STATEMENT
  EXECUTE PROCEDURE file_decrement_profile_size();