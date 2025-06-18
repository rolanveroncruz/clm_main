-- name: GetCertificateFromPK :one
SELECT * FROM certificate
WHERE pk = ? LIMIT 1;

-- name: GetCertificateFromSubjectCommonName :one
SELECT * FROM certificate
WHERE subject_common_name = ? LIMIT 1;

-- name: GetCertificateFromSubjectCommonNameAndUserEmail :one
SELECT * FROM certificate
WHERE subject_common_name = ? AND user_email=? LIMIT 1;

-- name: GetCertificatesFromUserEmail :many
SELECT * FROM certificate
WHERE user_email=?
ORDER BY sql_time_stamp;


-- name: ListCertificates :many
SELECT * FROM certificate
ORDER BY subject_common_name;

-- name: ListCertificatesByExpiration :many
SELECT * FROM certificate
ORDER BY not_after;

-- name: CreateCertificate :one
INSERT INTO certificate(user_email, subject_common_name, subject_organization, subject_organizational_unit, subject_country, subject_locality,
                        subject_province, issuer_common_name, issuer_organization, issuer_organizational_unit, issuer_country, serial_number,
                        signature_algorithm, public_key_algorithm, public_key_size, not_before, not_after, requested_server)
VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
RETURNING *;

-- name: UpdateCertificate :exec
UPDATE certificate
set user_email=?, subject_common_name=?, subject_organization=?, subject_organizational_unit=?, subject_country=?, subject_locality=?, subject_province=?,
    issuer_common_name=?, issuer_country=?, issuer_organization=?, issuer_organizational_unit=?, serial_number=?, not_before=?, not_after=?, requested_server=?,
    signature_algorithm=?, public_key_algorithm=?, public_key_size=?
WHERE pk = ?;

-- name: DeleteCertificate :exec
DELETE FROM certificate
WHERE pk = ?;
