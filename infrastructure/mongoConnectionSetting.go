package infrastructure

type MongoConnectionSetting struct {
	connectionUri string
	database      string
}

func (m MongoConnectionSetting) Database() string {
	return m.database
}

func (m MongoConnectionSetting) ConnectionUri() string {
	return m.connectionUri
}

func NewMongoConnectionSetting(connectionUri string, database string) *MongoConnectionSetting {
	return &MongoConnectionSetting{
		connectionUri: connectionUri,
		database:      database,
	}
}