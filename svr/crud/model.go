package crud

import (
	"reimirno.com/golinks/pkg/pb"
	"reimirno.com/golinks/pkg/types"
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
