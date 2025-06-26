package updater

type Manifest struct {
	Version     string `json:"version"`
	ObjectKey   string `json:"objectKey"`
	ChecksumSHA string `json:"checksum"`
}
