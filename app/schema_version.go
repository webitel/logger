package app

import (
	"context"
	errors "github.com/webitel/engine/model"
	"github.com/webitel/logger/model"
	proto "github.com/webitel/protos/logger"
)

func (a *App) GetFlowVersions(ctx context.Context, in *proto.SearchSchemaVersionRequest) (*proto.SchemaVersions, errors.AppError) {
	var (
		res       proto.SchemaVersions
		searchOpt *model.SearchOptions
	)

	searchOpt = ExtractSearchOptions(in)
	modelVersions, err := a.storage.SchemaVersion().Search(
		ctx,
		searchOpt,
		model.Filter{
			Column:         "schema_version.schema_id",
			Value:          in.GetSchemaId(),
			ComparisonType: model.Equal,
		},
	)

	versions, err := convertSchemaVersionModelToMessageBulk(modelVersions)
	if err != nil {
		return nil, err
	}
	res.Page = int32(searchOpt.Page)
	if err != nil {
		if IsErrNoRows(err) {
			return &res, nil
		} else {
			return nil, err
		}
	}

	if len(versions)-1 == searchOpt.Size {
		res.Next = true
		res.Items = versions[0 : len(versions)-1]
	} else {
		res.Items = versions
	}

	return &res, nil

}

func convertSchemaVersionModelToMessageBulk(in []*model.SchemaVersion) ([]*proto.SchemaVersion, errors.AppError) {
	if in == nil {
		return nil, errors.NewInternalError("app.schema_version.convert_schema_version_model.check_args.error", "nothing to convert")
	}
	var out []*proto.SchemaVersion

	for _, version := range in {
		res := &proto.SchemaVersion{
			Id:       int32(version.Id),
			SchemaId: int32(version.SchemaId),
			//CreatedOn: version.CreatedOn.UnixMilli(),
			CreatedBy: &proto.Lookup{
				Id:   version.CreatedBy.Id.Int32(),
				Name: version.CreatedBy.Name.String(),
			},
			State:   string(version.ObjectData),
			Version: int32(version.Version),
			Note:    version.Note.String(),
		}

		if !version.CreatedOn.IsZero() {
			res.CreatedOn = version.CreatedOn.UnixMilli()
		}
		out = append(out, res)
	}
	return out, nil
}
