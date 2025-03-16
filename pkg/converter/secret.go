package converter

import (
	"github.com/romanp1989/gophkeeper/internal/server/domain"
	"github.com/romanp1989/gophkeeper/pkg/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProtoToType конвертирует объект protobuf SecretType в объект модели данных SecretType
func ProtoToType(pbType proto.SecretType) domain.SecretType {
	switch pbType {
	case proto.SecretType_SECRET_TYPE_CREDENTIALS:
		return domain.CredSecret
	case proto.SecretType_SECRET_TYPE_TEXT:
		return domain.TextSecret
	case proto.SecretType_SECRET_TYPE_BLOB:
		return domain.BlobSecret
	case proto.SecretType_SECRET_TYPE_CARD:
		return domain.CardSecret
	default:
		return domain.UnknownSecret
	}
}

// TypeToProto конвертирует объект модели данных SecretType в объект protobuf SecretType
func TypeToProto(sType string) proto.SecretType {
	switch sType {
	case string(domain.CredSecret):
		return proto.SecretType_SECRET_TYPE_CREDENTIALS
	case string(domain.TextSecret):
		return proto.SecretType_SECRET_TYPE_TEXT
	case string(domain.BlobSecret):
		return proto.SecretType_SECRET_TYPE_BLOB
	case string(domain.CardSecret):
		return proto.SecretType_SECRET_TYPE_CARD
	default:
		return proto.SecretType_SECRET_TYPE_UNSPECIFIED
	}
}

// SecretToProto конвертирует объект модели данных Secret в объект protobuf Secret
func SecretToProto(secret *domain.Secret) *proto.Secret {
	return &proto.Secret{
		Id:         secret.ID,
		Title:      secret.Title,
		Metadata:   secret.Metadata,
		Payload:    secret.Payload,
		SecretType: TypeToProto(secret.SecretType),
		CreatedAt:  timestamppb.New(secret.CreatedAt),
		UpdatedAt:  timestamppb.New(secret.UpdatedAt),
	}
}

// ProtoToSecret конвертирует объект protobuf Secret в объект Secret модели данных
func ProtoToSecret(pbSecret *proto.Secret) *domain.Secret {
	return &domain.Secret{
		ID:         pbSecret.Id,
		Title:      pbSecret.Title,
		Metadata:   pbSecret.Metadata,
		SecretType: string(ProtoToType(pbSecret.SecretType)),
		Payload:    pbSecret.Payload,
		CreatedAt:  pbSecret.CreatedAt.AsTime(),
		UpdatedAt:  pbSecret.UpdatedAt.AsTime(),
	}
}

// ProtoToSecrets конвертирует список объекто protobuf Secret в список объектов Secret модели данных
func ProtoToSecrets(pbSecrets []*proto.Secret) []*domain.Secret {
	var secrets []*domain.Secret
	for _, s := range pbSecrets {
		secrets = append(secrets, ProtoToSecret(s))
	}
	return secrets
}

// SecretsToProto конвертирует список объекто модели данных Secret в список объектов Secret protobuf
func SecretsToProto(secrets []*domain.Secret) []*proto.Secret {
	var pbSecrets []*proto.Secret
	for _, s := range secrets {
		pbSecrets = append(pbSecrets, SecretToProto(s))
	}
	return pbSecrets
}
