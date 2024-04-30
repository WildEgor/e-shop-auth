package repositories

import (
	"context"
	mongo2 "github.com/WildEgor/e-shop-auth/internal/db/mongodb"
	"github.com/WildEgor/e-shop-auth/internal/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserRepository struct {
	db *mongo2.MongoConnection
}

func NewUserRepository(
	db *mongo2.MongoConnection,
) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (ur *UserRepository) CheckExistsPhone(phone string) bool {
	filter := bson.M{"phone": phone, "status": bson.M{"$eq": models.ActiveStatus}}

	v, err := ur.db.DB().Collection(models.CollectionUsers).CountDocuments(context.TODO(), filter)
	if err != nil || v > 0 {
		return true
	}

	return false
}

func (ur *UserRepository) CheckExistsEmail(email string) bool {
	filter := bson.M{"email": email, "status": bson.M{"$eq": models.ActiveStatus}}

	v, err := ur.db.DB().Collection(models.CollectionUsers).CountDocuments(context.TODO(), filter)
	if err != nil || v > 0 {
		return true
	}

	return false
}

func (ur *UserRepository) FindByPhone(phone string) (*models.UsersModel, error) {
	filter := bson.D{{Key: "phone", Value: phone}}

	var us models.UsersModel
	if err := ur.db.DB().Collection(models.CollectionUsers).FindOne(context.TODO(), filter).Decode(&us); err != nil {
		return nil, err
	}

	return &us, nil
}

func (ur *UserRepository) FindByEmail(email string) (*models.UsersModel, error) {
	filter := bson.D{{Key: "email", Value: email}}

	var us models.UsersModel
	if err := ur.db.DB().Collection(models.CollectionUsers).FindOne(context.TODO(), filter).Decode(&us); err != nil {
		return nil, err
	}

	return &us, nil
}

func (ur *UserRepository) FindByLogin(login string) (*models.UsersModel, error) {
	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{"phone", bson.D{{"$eq", login}}}},
				bson.D{{"email", bson.D{{"$eq", login}}}},
			},
		},
	}

	var us models.UsersModel
	if err := ur.db.DB().Collection(models.CollectionUsers).FindOne(context.TODO(), filter).Decode(&us); err != nil {
		return nil, err
	}

	return &us, nil
}

func (ur *UserRepository) FindByIds(ids []string) (*[]models.UsersModel, error) {

	oids := make([]primitive.ObjectID, 0)
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}

		oids = append(oids, oid)
	}

	if len(oids) == 0 {
		return nil, errors.New("empty ids")
	}

	filter := bson.D{{"_id", bson.D{{"$in", oids}}}}

	cursor, err := ur.db.DB().Collection(models.CollectionUsers).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var users []models.UsersModel
	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	return &users, nil
}

func (ur *UserRepository) CountAllActive() (int32, error) {
	count, err := ur.db.DB().Collection(models.CollectionUsers).CountDocuments(context.TODO(), bson.M{"$eq": models.ActiveStatus})
	if err != nil {
		return 0, errors.Wrap(err, "Mongo error")
	}

	return int32(count), nil
}

func (ur *UserRepository) FindById(id string) (*models.UsersModel, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{{Key: "_id", Value: oid}}

	var us models.UsersModel
	if err := ur.db.DB().Collection(models.CollectionUsers).FindOne(context.TODO(), filter).Decode(&us); err != nil {
		return nil, err
	}

	return &us, nil
}

func (ur *UserRepository) Create(nu *models.UsersModel) (*models.UsersModel, error) {
	var checkUser models.UsersModel
	checkUser.Email = nu.Email
	checkUser.Phone = nu.Phone

	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{"phone", bson.D{{"$eq", nu.Phone}}}},
				bson.D{{"email", bson.D{{"$eq", nu.Email}}}},
			}},
	}
	count, err := ur.db.DB().Collection(models.CollectionUsers).CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, errors.Wrap(err, "Mongo error")
	}

	if count > 0 {
		err = errors.New("")
		return nil, errors.Wrap(err, "Phone or Email already taken")
	}

	us := &models.UsersModel{
		Email:        nu.Email,
		Phone:        nu.Phone,
		Password:     nu.Password,
		FirstName:    nu.FirstName,
		LastName:     nu.LastName,
		Verification: models.VerificationModel{},
		OTP:          models.OTPModel{},
		Status:       models.ActiveStatus,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	insertResult, err := ur.db.DB().Collection(models.CollectionUsers).InsertOne(context.TODO(), us)
	if err != nil {
		return nil, errors.New(`{"mail":"need uniq mail"}`)
	}

	us.Id = insertResult.InsertedID.(primitive.ObjectID)

	return us, nil
}

func (ur *UserRepository) UpdateInfo(nu *models.UsersModel) error {
	nu.UpdatedAt = time.Now().UTC()

	update := bson.D{
		{"$set",
			bson.D{
				{"first_name", nu.FirstName},
				{"last_name", nu.LastName},
				{"updated_at", nu.UpdatedAt},
			},
		},
	}

	_, err := ur.db.DB().Collection(models.CollectionUsers).UpdateByID(context.TODO(), nu.Id, update)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) UpdateContacts(nu *models.UsersModel) error {
	nu.UpdatedAt = time.Now().UTC()

	update := bson.D{
		{"$set",
			bson.D{
				{"phone", nu.Phone},
				{"email", nu.Email},
				{"updated_at", nu.UpdatedAt},
			},
		},
	}

	_, err := ur.db.DB().Collection(models.CollectionUsers).UpdateByID(context.TODO(), nu.Id, update)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) UpdatePassword(nu *models.UsersModel) error {
	nu.UpdatedAt = time.Now().UTC()

	update := bson.D{
		{"$set",
			bson.D{
				{"password", nu.Password},
				{"updated_at", nu.UpdatedAt},
			},
		},
	}

	_, err := ur.db.DB().Collection(models.CollectionUsers).UpdateByID(context.TODO(), nu.Id, update)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) UpdateOTP(id primitive.ObjectID, otp *models.OTPModel) error {
	update := bson.D{
		{"$set",
			bson.D{
				{"otp", otp},
				{"updated_at", time.Now().UTC()},
			},
		},
	}

	_, err := ur.db.DB().Collection(models.CollectionUsers).UpdateByID(context.TODO(), id, update)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) UpdateVerification(id primitive.ObjectID, v *models.VerificationModel) error {
	update := bson.D{
		{"$set",
			bson.D{
				{"verification", v},
				{"updated_at", time.Now().UTC()},
			},
		},
	}

	_, err := ur.db.DB().Collection(models.CollectionUsers).UpdateByID(context.TODO(), id, update)
	if err != nil {
		return err
	}

	return nil

}
