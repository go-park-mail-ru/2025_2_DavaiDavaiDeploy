SELECT EXISTS(
    SELECT 1 FROM user_table 
    WHERE vkid = $1
);