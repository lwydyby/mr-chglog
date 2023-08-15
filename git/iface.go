package git

type GitClient interface {
	GetTags() []*Tag
	GetMergeRequests(from, end *Tag) []*MergeRequest
	CreateTag(tag string, desc string)
	GetMRChanges(mr *MergeRequest)
	UpdateTagRelease(tagName string, desc string)
}
