-- name: GetCertificateFromPK :one
SELECT * FROM certificate
WHERE pk = ? LIMIT 1;

-- name: GetCertificateFromSubjectCommonName :one
SELECT * FROM certificate
WHERE subject_common_name = ? LIMIT 1;

-- name: ListCertificates :many
SELECT * FROM certificate
ORDER BY subject_common_name;

-- name: ListCertificatesByExpiration :many
SELECT * FROM certificate
ORDER BY not_after;

-- name: CreateCertificate :one
INSERT INTO certificate(
                        user_email, subject_common_name, subject_organization, subject_organization_unit, subject_country, subject_locality, subject_province,
                        issuer_id, serial_number, not_before, not_after
) VALUES ( ?,?,?,?,?,?,?,?,?,?,?)
RETURNING *;

-- name: UpdateCertificate :exec
UPDATE certificate
set user_email=?, subject_common_name=?, subject_organization=?, subject_organization_unit=?, subject_country=?, subject_locality=?, subject_province=?,
    issuer_id=?, serial_number=?, not_before=?, not_after=?
WHERE pk = ?;

-- name: DeleteCertificate :exec
DELETE FROM certificate
WHERE pk = ?;
