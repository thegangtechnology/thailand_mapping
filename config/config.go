package config

const (
	UserKey             = "User"
	HTTPStatusCheckFail = 300
)

const TransactionGetter = "db_trx"

const False = "false"
const True = "true"

func GetImageType() []string {
	return []string{
		"image/*",
		"image/_*",
		"image/apng",
		"image/avif",
		"image/bmp",
		"image/gif",
		"image/jpeg",
		"image/png",
		"image/svg+xml",
		"image/tiff",
		"image/webp",
		"image/x-ico",
	}
}
