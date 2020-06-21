package mongo

import (
	"context"
	"fmt"
	"github.com/common-go/oauth2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type MongoIntegrationConfigurationRepository struct {
	Collection             *mongo.Collection
	OAuth2UserRepositories map[string]oauth2.OAuth2UserRepository
}

func NewMongoIntegrationConfigurationRepository(db *mongo.Database, collectionName string, oauth2UserRepositories map[string]oauth2.OAuth2UserRepository) *MongoIntegrationConfigurationRepository {
	collection := db.Collection(collectionName)
	return &MongoIntegrationConfigurationRepository{collection, oauth2UserRepositories}
}

func (s *MongoIntegrationConfigurationRepository) GetIntegrationConfiguration(ctx context.Context, id string) (*oauth2.IntegrationConfiguration, string, error) {
	var model oauth2.IntegrationConfiguration
	query := bson.M{"_id": id}
	x := s.Collection.FindOne(ctx, query)
	if x.Err() != nil {
		if strings.Compare(fmt.Sprint(x.Err()), "mongo: no documents in result") == 0 {
			return nil, "", nil
		}
		return nil, "", x.Err()
	}
	k := &model
	err := x.Decode(k)
	if err != nil {
		return nil, "", err
	}

	clientId := model.ClientId
	k.ClientId, err = s.OAuth2UserRepositories[id].GetRequestTokenOAuth(ctx, model.ClientId, model.ClientSecret)
	return k, clientId, err
}
