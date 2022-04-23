package sqlx

type connMeta interface {
	getInstance() string
}
