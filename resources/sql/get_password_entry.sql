SELECT
    title
FROM
    passwords
WHERE
    title = $1
    AND name = $2
    AND author = $3;
