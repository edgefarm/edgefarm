package getrelease

type getrelease interface {
	DownloadPublic(repo string, tag string, outputDirectory string) error
}
