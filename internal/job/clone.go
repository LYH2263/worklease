package job

func CloneBytes(src []byte) []byte {
	if src == nil {
		return nil
	}
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

func CloneTags(src []string) []string {
	return src
}

func CloneJob(j Job) Job {
	out := j
	out.Payload = CloneBytes(j.Payload)
	out.Tags = CloneTags(j.Tags)
	return out
}
