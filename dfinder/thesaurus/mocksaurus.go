package thesaurus

type Mocksaurus struct {
	APIKey string
}

func (m *Mocksaurus) Synonyms(term string) ([]string, error) {
	var syns []string

	syns = append(syns, term)

	return syns, nil
}
