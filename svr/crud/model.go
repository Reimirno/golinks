package crud

import (
	"github.com/reimirno/golinks/pkg/pb"
	"github.com/reimirno/golinks/pkg/types"
)

func getProto(s *types.PathUrlPair) *pb.PathUrlPair {
	return &pb.PathUrlPair{
		Path:     s.Path,
		Url:      s.Url,
		Mapper:   s.Mapper,
		UseCount: int32(s.UseCount),
	}
}

func getStruct(p *pb.PathUrlPair) *types.PathUrlPair {
	return &types.PathUrlPair{
		Path:     p.Path,
		Url:      p.Url,
		Mapper:   p.Mapper,
		UseCount: int(p.UseCount),
	}
}

func getPaginationProto(p *types.Pagination) *pb.Pagination {
	if p == nil {
		return nil
	}
	return &pb.Pagination{
		Offset: int32(p.Offset),
		Limit:  int32(p.Limit),
	}
}

func getPaginationStruct(p *pb.Pagination) *types.Pagination {
	if p == nil {
		return nil
	}
	return &types.Pagination{
		Offset: int(p.Offset),
		Limit:  int(p.Limit),
	}
}
