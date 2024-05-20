package service

const (
	NameSql = `SELECT * FROM "name" WHERE ("name"."type" = $1 AND "name"."transaction_id" = $2) AND "name"."deleted_at" IS NULL ORDER BY created_at DESC`

	DobSql                   = `SELECT * FROM "dob" WHERE ("dob"."type" = $1 AND "dob"."transaction_id" = $2) AND "dob"."deleted_at" IS NULL ORDER BY created_at DESC`
	DobSqlAfterUpdateChanged = `UPDATE "dob" SET "updated_at"=$1,"changed"=$2 WHERE "dob"."deleted_at" IS NULL AND "type" = $3 AND "transaction_id" = $4`

	AddrSql             = `SELECT * FROM "address" WHERE ("address"."type" = $1 AND "address"."transaction_id" = $2) AND "address"."deleted_at" IS NULL ORDER BY created_at DESC`
	AddrSqlAfterUpdate  = `SELECT * FROM "address" WHERE "address"."id" = $1 AND "address"."deleted_at" IS NULL ORDER BY created_at DESC`
	AddressUpdateSqlNew = `UPDATE "address" SET "updated_at"=$1,"status"=$2,"expectation_id"=$3,"country"=$4 WHERE "address"."deleted_at" IS NULL AND "id" = $5 AND "type" = $6 AND "pa_type" = $7 AND "transaction_id" = $8`

	BankLinkDataSqlInsert = `INSERT INTO "banks" ("created_at","updated_at","deleted_at","bank_link_id","id","code","provider_id","provider_name") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	BankLinkDataSql       = `SELECT * FROM "banks" WHERE "banks"."bank_link_id" = $1 AND "banks"."deleted_at" IS NULL`
	BankDataSql           = `SELECT * FROM "banks" WHERE "banks"."deleted_at" IS NULL AND "banks"."id" = $1`

	ExptSql                = `SELECT * FROM "expectation" WHERE ("expectation"."type" = $1 AND "expectation"."transaction_id" = $2) AND "expectation"."deleted_at" IS NULL ORDER BY created_at DESC`
	ExptUpdateSqlNew       = `UPDATE "expectation" SET "updated_at"=$1,"status"=$2,"expectation_id"=$3,"data"=$4 WHERE "expectation"."deleted_at" IS NULL AND "id" = $5 AND "type" = $6 AND "pa_type" = $7 AND "transaction_id" = $8`
	ExptUpdateSqlNewNoData = `UPDATE "expectation" SET "updated_at"=$1,"status"=$2,"expectation_id"=$3 WHERE "expectation"."deleted_at" IS NULL AND "id" = $4 AND "type" = $5 AND "pa_type" = $6 AND "transaction_id" = $7`
	ExptSqlAfterUpdate     = `SELECT * FROM "expectation" WHERE "expectation"."deleted_at" IS NULL AND "id" = $1 AND "type" = $2 AND "pa_type" = $3 AND "transaction_id" = $4 AND "expectation"."id" = $5 AND "expectation"."type" = $6 AND "expectation"."pa_type" = $7 AND "expectation"."transaction_id" = $8`

	ExptMetasSql      = `SELECT * FROM "expectation_metadata" WHERE "expectation_metadata"."expectation_id" = $1 AND "expectation_metadata"."transaction_id" = $2 AND "expectation_metadata"."deleted_at" IS NULL`
	ExptMetaUpdateSql = `UPDATE "expectation_metadata" SET "updated_at"=$1,"transaction_id"=$2,"type"=$3,"data"=$4 WHERE "expectation_metadata"."deleted_at" IS NULL AND "id" = $5 AND "expectation_id" = $6`

	DocSql            = `SELECT * FROM "document" WHERE ("document"."type" = $1 AND "document"."transaction_id" = $2) AND "document"."deleted_at" IS NULL ORDER BY created_at DESC`
	DocUpdateSqlNew   = `UPDATE "document" SET "updated_at"=$1,"status"=$2,"expectation_id"=$3 WHERE "document"."deleted_at" IS NULL AND "id" = $4 AND "type" = $5 AND "pa_type" = $6 AND "transaction_id" = $7`
	DocSqlAfterUpdate = `SELECT * FROM "document" WHERE "document"."deleted_at" IS NULL AND "id" = $1 AND "type" = $2 AND "pa_type" = $3 AND "transaction_id" = $4 AND "document"."id" = $5 AND "document"."type" = $6 AND "document"."pa_type" = $7 AND "document"."transaction_id" = $8`

	DocsSql       = `SELECT * FROM "documents" WHERE "documents"."document_id" = $1 AND "documents"."deleted_at" IS NULL`
	DocsInsertSql = `INSERT INTO "documents" ("created_at","updated_at","deleted_at","document_id","id","type") VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT ("document_id","id") DO UPDATE SET "document_id"="excluded"."document_id"`
)
