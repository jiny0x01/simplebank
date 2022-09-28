-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;


/* 특정 params만 업데이트 할 때 */ 
/* 1안.CASE 사용 */
-- name: UpdateUser :one
/*
UPDATE users
SET
    hashed_password = CASE
        WHEN @set_hahsed_password::boolean = TRUE THEN @hashed_password
        ELSE hashed_password
    END,
    full_name = CASE
        WHEN @set_full_name::boolean = TRUE THEN @full_name
        ELSE full_name 
    END,
    email = CASE
        WHEN @set_email::boolean = TRUE THEN @email
        ELSE email 
    END
WHERE
    username = @username
RETURNING *;
*/

/* 2안. COALSECE와 sqlc.narg()사용 */
/* sqlc.narg()는 파라미터에 null이 들어갈 수 있음*/
/* where절에선 검색해야하니 sqlc.narg()가 아닌 sqlc.arg() */

-- name: UpdateUser :one
UPDATE users
SET
    hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
    full_name = COALESCE(sqlc.narg(full_name), full_name),
    email = COALESCE(sqlc.narg(email), email)
WHERE
    username = sqlc.arg(username) 
RETURNING *;

