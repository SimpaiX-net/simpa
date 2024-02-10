package drivers

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"simpaix.net/simpa/v2/engine/crypt"
	"simpaix.net/simpa/v2/engine/sessions"
)

type MongoStore struct {
	*mongo.Collection
	timeout time.Duration
	crypt   crypt.CrypterI
}

// Creates new mongo storage
func NewMongoStore(col *mongo.Collection, timeout time.Duration, crypt crypt.CrypterI) *MongoStore {
	return &MongoStore{col, timeout, crypt}
}

// Saves session back to the store. Encrypts the [session.values] object.
func (s *MongoStore) Set(sess *sessions.Session) error {
	obj_id, err := primitive.ObjectIDFromHex(sess.ID)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*s.timeout)
	defer cancel()

	raw, err := json.Marshal(sess.Values)
	if err != nil {
		return err
	}

	opts, err := json.Marshal(sess.Opts)
	if err != nil {
		return err
	}

	enc_raw, err := s.crypt.Encrypt(string(raw))
	if err != nil {
		return err
	}

	if err := s.Collection.FindOneAndUpdate(ctx, bson.M{
		"_id": obj_id,
	}, bson.M{
		"$set": bson.D{
			{Key: "values", Value: string(enc_raw)},
			{Key: "options", Value: string(opts)},
		},
	}, options.FindOneAndUpdate().SetUpsert(true)).Err(); err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	return nil
}

// Gets session back from the store. Decrypt the values object back to fit the session object's values field
func (s *MongoStore) Get(sid string) (*sessions.Session, error) {
	obj_id, err := primitive.ObjectIDFromHex(sid)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*s.timeout)
	defer cancel()

	var sess_ map[string]interface{}
	if err := s.Collection.FindOne(ctx, bson.M{
		"_id": obj_id,
	}).Decode(&sess_); err != nil {
		return nil, err
	}

	sess := new(sessions.Session)
	{
		sess.ID = sid
		sess.Opts = new(sessions.Config)
	}

	opts, ok := sess_["options"].(string)
	if !ok {
		return nil, errors.New("invalid options")
	}

	if err := json.Unmarshal([]byte(opts), sess.Opts); err != nil {
		return nil, err
	}

	vals, ok := sess_["values"].(string)
	if !ok {
		return nil, errors.New("invalid values")
	}

	vals, err = s.crypt.Decrypt(vals)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(vals), &sess.Values); err != nil {
		return nil, err
	}

	return sess, nil
}

// Purges the session: TODO; lurky xd
func (s *MongoStore) Purge(sid string) error {
	return nil
}
