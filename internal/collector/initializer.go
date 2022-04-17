package collector

type Article struct {
	Link    string
	Title   string
	Summary string
	Author  Author
	Date    string
	Texts   []string
	Urls    []string
}

type Author struct {
	Name string
	Link string
}
