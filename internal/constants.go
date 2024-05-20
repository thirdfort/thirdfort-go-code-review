package internal

import "time"

// PA/CE statuses
const (
	StatusOpen      = "open"
	StatusClosed    = "closed"
	StatusSubmitted = "submitted"
	StatusPending   = "pending"
	StatusComplete  = "complete"
	StatusCompleted = "completed" // used by apps
	StatusAccepted  = "accepted"
	StatusCancelled = "cancelled" // used by apps
	StatusFulfilled = "fulfilled"
	StatusError     = "error"
	StatusRejected  = "rejected"
	StatusAborted   = "aborted"
)

// statuses we use for FE
const (
	StatusNotStarted = "not_started" // pending
	StatusInProgress = "in_progress" // open
	StatusInReview   = "in_review"   // submitted
	// StatusCompleted  = "completed"
	// StatusCancelled  = "cancelled"
)

// date and time related variables
const (
	TimeYear      = 356 * 24 * time.Hour // Helper for calculating age
	DateLayout    = "2006-01-02"
	ActorCacheTTL = 15 * time.Minute
)

// Misc constants
const (
	HeaderKey      = "headers"
	OnfidoAppIDKey = "app_id"
	OnfidoTokenKey = "onfido_token"
)

// mime types
const (
	ContentTypeTextPlain = "text/plain"
	ContentTypeJson      = "application/json"
)

// Tables with embedded tables
const (
	DbTableDocument    = "document"
	DbTableExpectation = "expectation"
)

// Task types
const (
	TypePersonalInformation       = "personal_information"
	TypePersonalInformationName   = "pi_name"
	TypePersonalInformationDob    = "pi_dob"
	TypeAddress                   = "address"
	TypeOriginalIdentity          = "original_id"
	TypeOriginalIdentityDocs      = "oi_onfido_capture_documents"
	TypeOriginalIdentityVideo     = "oi_onfido_capture_video"
	TypeStandardIdentity          = "standard_id"
	TypeStandardIdentityNfc       = "si_nfc"
	TypeStandardIdentityIproovBio = "si_iproov_biometrics"
	TypeStandardIdentityIdentity  = "si_documents_identity"
	TypeSofPurchaser              = "sof_purchaser"
	TypeSofPurchaserProperty      = "sof_purchaser_property"
	TypeSofPurchaserQuestionnaire = "sof_purchaser_questionnaire"
	TypeSofGiftor                 = "sof_giftor"
	TypeSofGiftorGift             = "sof_giftor_gift"
	TypeSofGiftorQuestionnaire    = "sof_giftor_questionnaire"
	TypeBankLink                  = "bank_link"
	TypeDocumentsPoa              = "documents_poa"
	TypeDocumentsPoo              = "documents_poo"
	TypeDocumentsBankStatement    = "documents_bank_statement"
	TypeDocumentsDivorce          = "documents_divorce"
	TypeDocumentsInheritance      = "documents_inheritance"
	TypeDocumentsMortgage         = "documents_mortgage"
	TypeDocumentsSaleAssets       = "documents_sale_assets"
	TypeDocumentsSavings          = "documents_savings"
)

// Reason codes
const (
	CodeDatabaseAddrOrPii    = "database_address_or_pii"
	CodeDatabaseDob          = "database_dob"
	CodeDocumentDob          = "document_dob"
	CodeDocumentAuthenticity = "document_authenticity"
	CodeDocumentImageQuality = "document_image_quality"
	CodeDocumentExpiry       = "document_expiry"
	CodeDocumentUnsupported  = "document_unsupported"
	CodeVideoAuthenticity    = "video_authenticity"
	CodeVideoAndDocument     = "video_and_document"
	CodeDocumentName         = "document_name"
)

// Fail Flow Reasons
const (
	FFVerificationDob                  = "verification:dob"
	FFComparisonDob                    = "comparison:dob"
	FFAuthenticityOriginalDocument     = "authenticity:original-document"
	FFAuthenticityFonts                = "authenticity:fonts"
	FFAuthenticityPictureFaceIntegrity = "authenticity:picture-face-integrity"
	FFAuthenticitySecurity             = "authenticity:security"
	FFConsistencyDocumentType          = "consistency:document-type"
	FFConsistencyDocumentNumbers       = "consistency:document-numbers"
	FFIntegrityDocumentQuality         = "integrity:document quality"
	FFIntegrityImageQuality            = "integrity:image-quality"
	FFValidationDocumentNumbers        = "validation:document-numbers"
	FFValidationDocumentExpiration     = "validation:document-expiration"
	FFValidationExpiryDate             = "validation:expiry-date"
	FFIntegritySupportedDocument       = "integrity:supported-document"
	FFAuthenticityLiveliness           = "authenticity:liveliness"
	FFAuthenticitySpoofing             = "authenticity:spoofing"
	FFIntegrityImage                   = "integrity:image"
	FFComparisonFace                   = "comparison:face"
)

// Recommendation Slugs
const (
	SlugFootPrint  = "address+dob+name|documents:poa"
	SlugDob        = "dob"
	SlugName       = "name"
	SlugOcd        = "onfido:capture:documents"
	SlugOcv        = "onfido:capture:video"
	SlugOriginalId = "onfido:capture:documents+onfido:capture:video"
)

// Expectations
const (
	ExptDocumentBankStatement          = "documents:bank-statement"
	ExptDocumentDivorce                = "documents:divorce"
	ExptDocumentIdentity               = "documents:identity"
	ExptDocumentInheritance            = "documents:inheritance"
	ExptDocumentMortgage               = "documents:mortgage"
	ExptDocumentPoa                    = "documents:poa"
	ExptDocumentPoo                    = "documents:poo"
	ExptDocumentSaleAssets             = "documents:sale-assets"
	ExptDocumentSavings                = "documents:savings"
	ExptIproovBiometrics               = "iproov:biometrics"
	ExptName                           = "name"
	ExptNationality                    = "nationality"
	ExptNfc                            = "nfc"
	ExptNiNumber                       = "ni-number"
	ExptOnfidoCaptureVideo             = "onfido:capture:video"
	ExptOnfidoCaptureDocuments         = "onfido:capture:documents"
	ExptPropertySof                    = "property"
	ExptGiftorSof                      = "giftor"
	ExptPropertyQuestionnaireSof       = "questionnaire:sof"
	ExptPropertyQuestionnaireSofGiftor = "questionnaire:sof-giftor"
	ExptTruelayerCode                  = "truelayer:code"
	ExptDob                            = "dob"
	ExptAddress                        = "address"
)

// Banking redirect-uri
const (
	BankRedirectUri = "https://www.thirdfort.com/banking-redirect/"
)
