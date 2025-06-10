// Package certs provides functions related to certificates.
// This file, certs_db, is for functions to discover, retrieve, and store cert information in the database.
package certs

import (
	"context"
	"database/sql"
	certsSQL "ph.certs.com/clm_main/certs/sql"
	"ph.certs.com/clm_main/sqlite"
)

// GetServerCert gets the certificate of the requested server, along with the signing chain of the certificate.

func insertIntoDB(certs []*JSONCertificate, userEmail string) error {
	var err error
	for _, cert := range certs {
		err = saveCertToDB(cert, userEmail)
		if err != nil {
			continue
		}
	}
	return err
}
func saveCertToDB(cert *JSONCertificate, userEmail string) error {
	ctx := context.Background()
	certificateParams := certsSQL.CreateCertificateParams{
		UserEmail:                 sql.NullString{String: userEmail, Valid: true},
		SubjectCommonName:         sql.NullString{String: cert.Subject.CommonName, Valid: true},
		SubjectOrganization:       FromStringArrayToNullString(cert.Subject.Organization),
		SubjectOrganizationalUnit: FromStringArrayToNullString(cert.Subject.OrganizationalUnit),
		SubjectCountry:            FromStringArrayToNullString(cert.Subject.Country),
		SubjectLocality:           FromStringArrayToNullString(cert.Subject.Locality),
		SubjectProvince:           FromStringArrayToNullString(cert.Subject.Province),
		IssuerCommonName:          sql.NullString{String: cert.Issuer.CommonName, Valid: true},
		IssuerOrganization:        FromStringArrayToNullString(cert.Issuer.Organization),
		IssuerOrganizationalUnit:  FromStringArrayToNullString(cert.Issuer.OrganizationalUnit),
		IssuerCountry:             FromStringArrayToNullString(cert.Issuer.Country),
		SerialNumber:              sql.NullString{String: cert.SerialNumber, Valid: true},
		NotBefore:                 sql.NullString{String: cert.NotBefore.String(), Valid: true},
		NotAfter:                  sql.NullString{String: cert.NotAfter.String(), Valid: true},
	}
	_, err := sqlite.CertsQueryCental.CreateCertificate(ctx, certificateParams)
	if err != nil {
		return err
	} else {
		return nil
	}
}
func FromStringArrayToNullString(data []string) sql.NullString {
	if len(data) == 0 {
		return sql.NullString{String: "", Valid: false}
	} else {
		return sql.NullString{String: data[0], Valid: true}
	}
}

func certAlreadyInServer(serverName string, requesterUserEmail string) bool {
	ctx := context.Background()
	commonNameAndEmail := certsSQL.GetCertificateFromSubjectCommonNameAndUserEmailParams{
		SubjectCommonName: sql.NullString{String: serverName, Valid: true},
		UserEmail:         requesterUserEmail,
	}
	_, err := sqlite.CertsQueryCental.GetCertificateFromSubjectCommonNameAndUserEmail(ctx, commonNameAndEmail)
	if err != nil {
		// err could be 'no rows in result set'.
		return false
	} else {
		return true
	}

}
